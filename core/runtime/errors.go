package runtime

// Error is a stable runtime error value.
type Error string

// Error returns the error message.
func (e Error) Error() string {
	return string(e)
}

const (
	// ErrNilEngine indicates that a nil engine was registered.
	ErrNilEngine Error = "runtime: nil engine"
	// ErrEmptyEngineName indicates that an engine returned an empty name.
	ErrEmptyEngineName Error = "runtime: empty engine name"
	// ErrEngineExists indicates that two engines share the same name.
	ErrEngineExists Error = "runtime: engine already exists"
	// ErrEngineNotFound indicates that a requested engine was not registered.
	ErrEngineNotFound Error = "runtime: engine not found"
	// ErrDependencyCycle indicates that engine dependencies contain a cycle.
	ErrDependencyCycle Error = "runtime: dependency cycle"
	// ErrDependencyMissing indicates that an engine dependency was not registered.
	ErrDependencyMissing Error = "runtime: dependency missing"
	// ErrInvalidLifecycleTransition indicates that a lifecycle transition is not allowed.
	ErrInvalidLifecycleTransition Error = "runtime: invalid lifecycle transition"
)
