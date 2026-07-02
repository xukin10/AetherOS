package runtime

// Manager stores engines and provides ordered iteration.
// Engines are stored in registration order.
type Manager struct {
	engines []*Engine
}

// NewManager creates an empty engine manager.
func NewManager() *Manager {
	return &Manager{engines: make([]*Engine, 0)}
}

// Add appends an engine to the managed list.
func (m *Manager) Add(engine *Engine) {
	m.engines = append(m.engines, engine)
}

// All returns a copy of all registered engines in insertion order.
func (m *Manager) All() []*Engine {
	result := make([]*Engine, len(m.engines))
	copy(result, m.engines)
	return result
}

// Len returns the number of registered engines.
func (m *Manager) Len() int {
	return len(m.engines)
}

// Get returns the engine at index i. Returns nil if out of range.
func (m *Manager) Get(i int) *Engine {
	if i < 0 || i >= len(m.engines) {
		return nil
	}
	return m.engines[i]
}
