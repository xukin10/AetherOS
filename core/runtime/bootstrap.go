package runtime

import "context"

// Bootstrap initializes the supplied runtime kernel.
func Bootstrap(ctx context.Context, rt Runtime) error {
	return rt.Bootstrap(ctx)
}
