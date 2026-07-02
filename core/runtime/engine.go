package runtime

import "context"

// Engine is the minimal unit of work in the runtime.
// Each engine is identified by Name and has Start/Stop lifecycle functions.
type Engine struct {
	name  string
	Start func(ctx context.Context) error
	Stop  func(ctx context.Context) error
}

// NewEngine creates a named engine with the given start and stop functions.
// Both functions may be nil, in which case the engine is a no-op placeholder.
func NewEngine(name string, start, stop func(ctx context.Context) error) *Engine {
	return &Engine{
		name:  name,
		Start: start,
		Stop:  stop,
	}
}

// Name returns the engine's unique identifier within a runtime.
func (e *Engine) Name() string {
	return e.name
}
