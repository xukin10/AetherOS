package runtime

import (
	"context"
)

type Runtime struct {
	manager   *Manager
	lifecycle Lifecycle
}

func NewRuntime() *Runtime {
	return &Runtime{
		manager:   NewManager(),
		lifecycle: Created,
	}
}

func (r *Runtime) Register(e Engine) {
	r.manager.Register(e)
}

func (r *Runtime) Start(ctx context.Context) error {
	r.lifecycle = Running

	for _, e := range r.manager.Engines() {
		if err := e.Start(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (r *Runtime) Stop(ctx context.Context) error {
	r.lifecycle = Stopped

	engines := r.manager.Engines()
	for i := len(engines) - 1; i >= 0; i-- {
		_ = engines[i].Stop(ctx)
	}

	return nil
}

func (r *Runtime) State() Lifecycle {
	return r.lifecycle
}