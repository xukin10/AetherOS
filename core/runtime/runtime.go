package runtime

import "context"

// Runtime manages the lifecycle of registered engines.
// It provides sequential Start and Stop across all engines in registration order.
type Runtime struct {
	manager *Manager
	state   State
}

// New creates an empty Runtime in the Created state.
func New() *Runtime {
	return &Runtime{
		manager: NewManager(),
		state:   StateCreated,
	}
}

// Register adds an engine to the runtime.
// Returns ErrEngineAlreadyRegistered if an engine with the same name exists.
func (r *Runtime) Register(engine *Engine) error {
	for _, e := range r.manager.All() {
		if e.Name() == engine.Name() {
			return ErrEngineAlreadyRegistered
		}
	}
	r.manager.Add(engine)
	return nil
}

// State returns the current runtime state.
func (r *Runtime) State() State {
	return r.state
}

// Start transitions the runtime to Running and calls Start on each engine in registration order.
// If any engine fails to start, already-started engines are stopped and
// the runtime returns to the Stopped state with the originating error.
func (r *Runtime) Start(ctx context.Context) error {
	if r.state == StateRunning {
		return ErrRuntimeAlreadyStarted
	}

	r.state = StateRunning

	for _, engine := range r.manager.All() {
		if engine.Start == nil {
			continue
		}
		if err := engine.Start(ctx); err != nil {
			r.stopEngines(ctx)
			return err
		}
	}

	return nil
}

// Stop transitions the runtime to Stopped and calls Stop on each engine in reverse registration order.
func (r *Runtime) Stop(ctx context.Context) error {
	if r.state != StateRunning {
		return ErrRuntimeNotStarted
	}

	r.stopEngines(ctx)
	r.state = StateStopped
	return nil
}

// stopEngines calls Stop on all engines in reverse registration order, best-effort.
func (r *Runtime) stopEngines(ctx context.Context) {
	for i := r.manager.Len() - 1; i >= 0; i-- {
		engine := r.manager.Get(i)
		if engine == nil || engine.Stop == nil {
			continue
		}
		_ = engine.Stop(ctx) // best-effort during rollback
	}
}
