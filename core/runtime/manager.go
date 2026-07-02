package runtime

import (
	"context"
	"fmt"
	"sort"
	"sync"
)

// EngineManager stores engines and coordinates dependency-aware lifecycle order.
type EngineManager struct {
	registryMu sync.Mutex
	runMu      sync.Mutex
	engines    map[string]Engine
	order      []Engine
	started    []Engine
}

// NewEngineManager creates an empty engine manager.
func NewEngineManager() *EngineManager {
	return &EngineManager{
		engines: make(map[string]Engine),
	}
}

// Register stores an engine by name.
func (m *EngineManager) Register(engine Engine) error {
	if engine == nil {
		return ErrNilEngine
	}

	name := engine.Name()
	if name == "" {
		return ErrEmptyEngineName
	}

	m.registryMu.Lock()
	defer m.registryMu.Unlock()

	if _, exists := m.engines[name]; exists {
		return fmt.Errorf("%w: %s", ErrEngineExists, name)
	}

	m.engines[name] = engine
	m.order = nil
	return nil
}

// Bootstrap validates dependencies and computes startup order.
func (m *EngineManager) Bootstrap() error {
	m.registryMu.Lock()
	defer m.registryMu.Unlock()

	order, err := m.resolveOrderLocked()
	if err != nil {
		return err
	}

	m.order = order
	return nil
}

// Start starts engines in dependency order.
func (m *EngineManager) Start(ctx context.Context) error {
	m.runMu.Lock()
	defer m.runMu.Unlock()

	m.registryMu.Lock()
	defer m.registryMu.Unlock()

	if m.order == nil {
		order, err := m.resolveOrderLocked()
		if err != nil {
			return err
		}
		m.order = order
	}

	m.started = nil
	for _, engine := range m.order {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := engine.Start(ctx); err != nil {
			return err
		}
		m.started = append(m.started, engine)
	}

	return nil
}

// Stop stops started engines in reverse startup order.
func (m *EngineManager) Stop(ctx context.Context) error {
	m.runMu.Lock()
	defer m.runMu.Unlock()

	for i := len(m.started) - 1; i >= 0; i-- {
		if err := ctx.Err(); err != nil {
			return err
		}
		engine := m.started[i]
		if err := engine.Stop(ctx); err != nil {
			return err
		}
	}

	m.started = nil
	return nil
}

// Engines returns a snapshot of registered engines by name.
func (m *EngineManager) Engines() map[string]Engine {
	m.registryMu.Lock()
	defer m.registryMu.Unlock()

	engines := make(map[string]Engine, len(m.engines))
	for name, engine := range m.engines {
		engines[name] = engine
	}
	return engines
}

// Engine returns a registered engine by name.
func (m *EngineManager) Engine(name string) (Engine, error) {
	m.registryMu.Lock()
	defer m.registryMu.Unlock()

	engine, ok := m.engines[name]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrEngineNotFound, name)
	}
	return engine, nil
}

// StartOrder returns the dependency-resolved startup order.
func (m *EngineManager) StartOrder() ([]string, error) {
	m.registryMu.Lock()
	defer m.registryMu.Unlock()

	if m.order == nil {
		order, err := m.resolveOrderLocked()
		if err != nil {
			return nil, err
		}
		m.order = order
	}

	names := make([]string, 0, len(m.order))
	for _, engine := range m.order {
		names = append(names, engine.Name())
	}
	return names, nil
}

func (m *EngineManager) resolveOrderLocked() ([]Engine, error) {
	visited := make(map[string]bool, len(m.engines))
	visiting := make(map[string]bool, len(m.engines))
	order := make([]Engine, 0, len(m.engines))

	var visit func(Engine) error
	visit = func(engine Engine) error {
		name := engine.Name()
		if visiting[name] {
			return fmt.Errorf("%w: %s", ErrDependencyCycle, name)
		}
		if visited[name] {
			return nil
		}

		visiting[name] = true
		dependencies := append([]string(nil), engine.Dependencies()...)
		sort.Strings(dependencies)
		for _, dependency := range dependencies {
			dependencyEngine, ok := m.engines[dependency]
			if !ok {
				return fmt.Errorf("%w: %s requires %s", ErrDependencyMissing, name, dependency)
			}
			if err := visit(dependencyEngine); err != nil {
				return err
			}
		}
		visiting[name] = false
		visited[name] = true
		order = append(order, engine)
		return nil
	}

	names := make([]string, 0, len(m.engines))
	for name := range m.engines {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		if err := visit(m.engines[name]); err != nil {
			return nil, err
		}
	}

	return order, nil
}
