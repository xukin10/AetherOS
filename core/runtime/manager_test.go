package runtime

import (
	"context"
	"errors"
	"reflect"
	"strconv"
	"sync"
	"testing"
)

type testEngine struct {
	name         string
	dependencies []string
	events       *[]string
}

func (e *testEngine) Name() string {
	return e.name
}

func (e *testEngine) Dependencies() []string {
	return e.dependencies
}

func (e *testEngine) Start(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	*e.events = append(*e.events, "start:"+e.name)
	return nil
}

func (e *testEngine) Stop(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	*e.events = append(*e.events, "stop:"+e.name)
	return nil
}

func (e *testEngine) Health(ctx context.Context) error {
	return ctx.Err()
}

func TestEngineManagerRegister(t *testing.T) {
	manager := NewEngineManager()
	events := []string{}
	engine := &testEngine{name: "engine-a", events: &events}

	if err := manager.Register(engine); err != nil {
		t.Fatalf("register engine: %v", err)
	}

	engines := manager.Engines()
	if len(engines) != 1 {
		t.Fatalf("expected 1 engine, got %d", len(engines))
	}
	if engines["engine-a"] != engine {
		t.Fatalf("registered engine mismatch")
	}
}

func TestEngineManagerRegisterDuplicate(t *testing.T) {
	manager := NewEngineManager()
	events := []string{}

	if err := manager.Register(&testEngine{name: "engine-a", events: &events}); err != nil {
		t.Fatalf("register engine: %v", err)
	}
	if err := manager.Register(&testEngine{name: "engine-a", events: &events}); !errors.Is(err, ErrEngineExists) {
		t.Fatalf("expected %v, got %v", ErrEngineExists, err)
	}
}

func TestEngineManagerEngineNotFound(t *testing.T) {
	manager := NewEngineManager()
	if _, err := manager.Engine("missing"); !errors.Is(err, ErrEngineNotFound) {
		t.Fatalf("expected %v, got %v", ErrEngineNotFound, err)
	}
}

func TestEngineManagerSimpleDAGStartOrder(t *testing.T) {
	manager := NewEngineManager()
	events := []string{}

	for _, engine := range []Engine{
		&testEngine{name: "api", dependencies: []string{"model"}, events: &events},
		&testEngine{name: "storage", events: &events},
		&testEngine{name: "model", dependencies: []string{"storage"}, events: &events},
	} {
		if err := manager.Register(engine); err != nil {
			t.Fatalf("register engine: %v", err)
		}
	}

	if err := manager.Bootstrap(); err != nil {
		t.Fatalf("bootstrap manager: %v", err)
	}
	if err := manager.Start(context.Background()); err != nil {
		t.Fatalf("start manager: %v", err)
	}

	expected := []string{"start:storage", "start:model", "start:api"}
	if !reflect.DeepEqual(events, expected) {
		t.Fatalf("expected events %v, got %v", expected, events)
	}
}

func TestEngineManagerComplexDAGStartOrder(t *testing.T) {
	manager := NewEngineManager()
	events := []string{}

	for _, engine := range []Engine{
		&testEngine{name: "api", dependencies: []string{"storage", "model"}, events: &events},
		&testEngine{name: "model", dependencies: []string{"config", "storage"}, events: &events},
		&testEngine{name: "storage", dependencies: []string{"config"}, events: &events},
		&testEngine{name: "config", dependencies: []string{"foundation"}, events: &events},
		&testEngine{name: "foundation", events: &events},
	} {
		if err := manager.Register(engine); err != nil {
			t.Fatalf("register engine: %v", err)
		}
	}

	order, err := manager.StartOrder()
	if err != nil {
		t.Fatalf("resolve start order: %v", err)
	}

	expected := []string{"foundation", "config", "storage", "model", "api"}
	if !reflect.DeepEqual(order, expected) {
		t.Fatalf("expected order %v, got %v", expected, order)
	}
}

func TestEngineManagerStopOrder(t *testing.T) {
	manager := NewEngineManager()
	events := []string{}

	for _, engine := range []Engine{
		&testEngine{name: "api", dependencies: []string{"model"}, events: &events},
		&testEngine{name: "storage", events: &events},
		&testEngine{name: "model", dependencies: []string{"storage"}, events: &events},
	} {
		if err := manager.Register(engine); err != nil {
			t.Fatalf("register engine: %v", err)
		}
	}

	if err := manager.Start(context.Background()); err != nil {
		t.Fatalf("start manager: %v", err)
	}
	events = nil
	if err := manager.Stop(context.Background()); err != nil {
		t.Fatalf("stop manager: %v", err)
	}

	expected := []string{"stop:api", "stop:model", "stop:storage"}
	if !reflect.DeepEqual(events, expected) {
		t.Fatalf("expected events %v, got %v", expected, events)
	}
}

func TestEngineManagerBootstrapMissingDependency(t *testing.T) {
	manager := NewEngineManager()
	events := []string{}

	if err := manager.Register(&testEngine{
		name:         "api",
		dependencies: []string{"storage"},
		events:       &events,
	}); err != nil {
		t.Fatalf("register engine: %v", err)
	}

	if err := manager.Bootstrap(); !errors.Is(err, ErrDependencyMissing) {
		t.Fatalf("expected %v, got %v", ErrDependencyMissing, err)
	}
}

func TestEngineManagerBootstrapDependencyCycle(t *testing.T) {
	manager := NewEngineManager()
	events := []string{}

	for _, engine := range []Engine{
		&testEngine{name: "api", dependencies: []string{"model"}, events: &events},
		&testEngine{name: "model", dependencies: []string{"storage"}, events: &events},
		&testEngine{name: "storage", dependencies: []string{"api"}, events: &events},
	} {
		if err := manager.Register(engine); err != nil {
			t.Fatalf("register engine: %v", err)
		}
	}

	if err := manager.Bootstrap(); !errors.Is(err, ErrDependencyCycle) {
		t.Fatalf("expected %v, got %v", ErrDependencyCycle, err)
	}
}

func TestEngineManagerConcurrentRegister(t *testing.T) {
	manager := NewEngineManager()
	events := []string{}
	const total = 100

	var wg sync.WaitGroup
	errs := make(chan error, total)
	for i := 0; i < total; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			errs <- manager.Register(&testEngine{
				name:   "engine-" + strconv.Itoa(i),
				events: &events,
			})
		}(i)
	}
	wg.Wait()
	close(errs)

	for err := range errs {
		if err != nil {
			t.Fatalf("register engine concurrently: %v", err)
		}
	}

	order, err := manager.StartOrder()
	if err != nil {
		t.Fatalf("resolve start order: %v", err)
	}
	if len(order) != total {
		t.Fatalf("expected %d engines, got %d", total, len(order))
	}
}
