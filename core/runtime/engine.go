package runtime

import "context"

// Engine is a runtime-managed component with lifecycle hooks.
type Engine interface {
	Name() string
	Dependencies() []string
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Health(ctx context.Context) error
}
