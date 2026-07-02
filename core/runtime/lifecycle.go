package runtime

// State represents the lifecycle state of an engine or runtime.
type State int

const (
	StateCreated State = iota
	StateRunning
	StateStopped
)

func (s State) String() string {
	switch s {
	case StateCreated:
		return "created"
	case StateRunning:
		return "running"
	case StateStopped:
		return "stopped"
	default:
		return "unknown"
	}
}
