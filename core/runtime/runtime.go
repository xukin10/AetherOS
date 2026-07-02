package runtime

import (
	"context"
	"fmt"
	"sync"
)

// Runtime is the public kernel API for registering and orchestrating engines.
type Runtime interface {
	Register(Engine) error
	Bootstrap(ctx context.Context) error
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Status() Lifecycle
}

type kernel struct {
	opMu      sync.Mutex
	manager   *EngineManager
	lifecycle *lifecycleStateMachine
}

// NewRuntime creates a runtime kernel in Created state.
func NewRuntime() Runtime {
	return &kernel{
		manager:   NewEngineManager(),
		lifecycle: newLifecycleStateMachine(),
	}
}

func (k *kernel) Register(engine Engine) error {
	k.opMu.Lock()
	defer k.opMu.Unlock()

	if status := k.lifecycle.Status(); status != Created {
		return fmt.Errorf("%w: register from %s", ErrInvalidLifecycleTransition, status)
	}

	return k.manager.Register(engine)
}

func (k *kernel) Bootstrap(ctx context.Context) error {
	k.opMu.Lock()
	defer k.opMu.Unlock()

	if err := ctx.Err(); err != nil {
		_ = k.lifecycle.Transition(Failed)
		return err
	}

	switch status := k.lifecycle.Status(); status {
	case Created:
	case Initialized:
		return nil
	default:
		return fmt.Errorf("%w: bootstrap from %s", ErrInvalidLifecycleTransition, status)
	}

	if err := k.manager.Bootstrap(); err != nil {
		_ = k.lifecycle.Transition(Failed)
		return err
	}

	if err := k.lifecycle.Transition(Initialized); err != nil {
		return err
	}
	return nil
}

func (k *kernel) Start(ctx context.Context) error {
	k.opMu.Lock()
	defer k.opMu.Unlock()

	if err := ctx.Err(); err != nil {
		_ = k.lifecycle.Transition(Failed)
		return err
	}

	switch status := k.lifecycle.Status(); status {
	case Created:
		if err := k.manager.Bootstrap(); err != nil {
			_ = k.lifecycle.Transition(Failed)
			return err
		}
		if err := k.lifecycle.Transition(Initialized); err != nil {
			return err
		}
	case Initialized:
	case Running:
		return nil
	default:
		return fmt.Errorf("%w: start from %s", ErrInvalidLifecycleTransition, status)
	}

	if err := k.lifecycle.Transition(Starting); err != nil {
		return err
	}
	if err := k.manager.Start(ctx); err != nil {
		_ = k.manager.Stop(ctx)
		_ = k.lifecycle.Transition(Failed)
		return err
	}

	return k.lifecycle.Transition(Running)
}

func (k *kernel) Stop(ctx context.Context) error {
	k.opMu.Lock()
	defer k.opMu.Unlock()

	if err := ctx.Err(); err != nil {
		_ = k.lifecycle.Transition(Failed)
		return err
	}

	switch status := k.lifecycle.Status(); status {
	case Running:
	case Stopped:
		return nil
	default:
		return fmt.Errorf("%w: stop from %s", ErrInvalidLifecycleTransition, status)
	}

	if err := k.lifecycle.Transition(Stopping); err != nil {
		return err
	}
	if err := k.manager.Stop(ctx); err != nil {
		_ = k.lifecycle.Transition(Failed)
		return err
	}

	return k.lifecycle.Transition(Stopped)
}

func (k *kernel) Status() Lifecycle {
	return k.lifecycle.Status()
}
