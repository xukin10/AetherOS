package runtime

import (
	"fmt"
	"sync"
)

// Lifecycle represents the runtime kernel lifecycle state.
type Lifecycle string

const (
	// Created is the initial state before bootstrap.
	Created Lifecycle = "Created"
	// Initialized means the runtime has validated and ordered engines.
	Initialized Lifecycle = "Initialized"
	// Starting means the runtime is starting registered engines.
	Starting Lifecycle = "Starting"
	// Running means all registered engines started successfully.
	Running Lifecycle = "Running"
	// Stopping means the runtime is stopping started engines.
	Stopping Lifecycle = "Stopping"
	// Stopped means all started engines have been stopped.
	Stopped Lifecycle = "Stopped"
	// Failed means bootstrap, start, or stop failed.
	Failed Lifecycle = "Failed"
)

type lifecycleStateMachine struct {
	mu    sync.RWMutex
	state Lifecycle
}

func newLifecycleStateMachine() *lifecycleStateMachine {
	return &lifecycleStateMachine{state: Created}
}

func (m *lifecycleStateMachine) Status() Lifecycle {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.state
}

func (m *lifecycleStateMachine) Transition(next Lifecycle) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.state == next {
		return nil
	}
	if !canTransition(m.state, next) {
		return fmt.Errorf("%w: %s -> %s", ErrInvalidLifecycleTransition, m.state, next)
	}

	m.state = next
	return nil
}

func canTransition(current, next Lifecycle) bool {
	if next == Failed {
		return current != Stopped
	}

	switch current {
	case Created:
		return next == Initialized
	case Initialized:
		return next == Starting
	case Starting:
		return next == Running
	case Running:
		return next == Stopping
	case Stopping:
		return next == Stopped
	default:
		return false
	}
}
