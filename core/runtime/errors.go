package runtime

import "errors"

var (
	ErrEngineNotFound          = errors.New("engine not found")
	ErrEngineAlreadyRegistered = errors.New("engine already registered")
	ErrRuntimeNotStarted       = errors.New("runtime not started")
	ErrRuntimeAlreadyStarted   = errors.New("runtime already started")
)
