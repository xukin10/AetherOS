package runtime

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"testing"
)

func TestRuntimeBootstrap(t *testing.T) {
	rt := NewRuntime()
	if rt.Status() != Created {
		t.Fatalf("expected status %s, got %s", Created, rt.Status())
	}

	if err := rt.Bootstrap(context.Background()); err != nil {
		t.Fatalf("bootstrap runtime: %v", err)
	}

	if rt.Status() != Initialized {
		t.Fatalf("expected status %s, got %s", Initialized, rt.Status())
	}
}

func TestRuntimeLifecycleTransitions(t *testing.T) {
	rt := NewRuntime()
	events := []string{}
	if err := rt.Register(&testEngine{name: "engine-a", events: &events}); err != nil {
		t.Fatalf("register engine: %v", err)
	}

	if err := rt.Bootstrap(context.Background()); err != nil {
		t.Fatalf("bootstrap runtime: %v", err)
	}
	if rt.Status() != Initialized {
		t.Fatalf("expected status %s, got %s", Initialized, rt.Status())
	}

	if err := rt.Start(context.Background()); err != nil {
		t.Fatalf("start runtime: %v", err)
	}
	if rt.Status() != Running {
		t.Fatalf("expected status %s, got %s", Running, rt.Status())
	}

	if err := rt.Stop(context.Background()); err != nil {
		t.Fatalf("stop runtime: %v", err)
	}
	if rt.Status() != Stopped {
		t.Fatalf("expected status %s, got %s", Stopped, rt.Status())
	}

	expected := []string{"start:engine-a", "stop:engine-a"}
	if !reflect.DeepEqual(events, expected) {
		t.Fatalf("expected events %v, got %v", expected, events)
	}
}

func TestRuntimeStartUsesDependencyOrder(t *testing.T) {
	rt := NewRuntime()
	events := []string{}

	for _, engine := range []Engine{
		&testEngine{name: "api", dependencies: []string{"model"}, events: &events},
		&testEngine{name: "storage", events: &events},
		&testEngine{name: "model", dependencies: []string{"storage"}, events: &events},
	} {
		if err := rt.Register(engine); err != nil {
			t.Fatalf("register engine: %v", err)
		}
	}

	if err := rt.Start(context.Background()); err != nil {
		t.Fatalf("start runtime: %v", err)
	}

	expected := []string{"start:storage", "start:model", "start:api"}
	if !reflect.DeepEqual(events, expected) {
		t.Fatalf("expected events %v, got %v", expected, events)
	}
}

func TestRuntimeBootstrapFailureSetsFailed(t *testing.T) {
	rt := NewRuntime()
	events := []string{}
	if err := rt.Register(&testEngine{name: "api", dependencies: []string{"storage"}, events: &events}); err != nil {
		t.Fatalf("register engine: %v", err)
	}

	if err := rt.Bootstrap(context.Background()); !errors.Is(err, ErrDependencyMissing) {
		t.Fatalf("expected %v, got %v", ErrDependencyMissing, err)
	}
	if rt.Status() != Failed {
		t.Fatalf("expected status %s, got %s", Failed, rt.Status())
	}
}

func TestLifecycleStateMachineValidTransitions(t *testing.T) {
	lifecycle := newLifecycleStateMachine()
	for _, next := range []Lifecycle{Initialized, Starting, Running, Stopping, Stopped} {
		if err := lifecycle.Transition(next); err != nil {
			t.Fatalf("transition to %s: %v", next, err)
		}
	}
	if lifecycle.Status() != Stopped {
		t.Fatalf("expected status %s, got %s", Stopped, lifecycle.Status())
	}
}

func TestLifecycleStateMachineInvalidTransition(t *testing.T) {
	lifecycle := newLifecycleStateMachine()
	if err := lifecycle.Transition(Running); !errors.Is(err, ErrInvalidLifecycleTransition) {
		t.Fatalf("expected %v, got %v", ErrInvalidLifecycleTransition, err)
	}
	if lifecycle.Status() != Created {
		t.Fatalf("expected status %s, got %s", Created, lifecycle.Status())
	}
}

func TestRuntimeRegisterAfterBootstrapFails(t *testing.T) {
	rt := NewRuntime()
	if err := rt.Bootstrap(context.Background()); err != nil {
		t.Fatalf("bootstrap runtime: %v", err)
	}

	events := []string{}
	if err := rt.Register(&testEngine{name: "late", events: &events}); !errors.Is(err, ErrInvalidLifecycleTransition) {
		t.Fatalf("expected %v, got %v", ErrInvalidLifecycleTransition, err)
	}
}

func TestRuntimeConcurrentStart(t *testing.T) {
	rt := NewRuntime()
	engine := &countingEngine{name: "engine-a"}
	if err := rt.Register(engine); err != nil {
		t.Fatalf("register engine: %v", err)
	}

	const total = 20
	var wg sync.WaitGroup
	errs := make(chan error, total)
	for i := 0; i < total; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errs <- rt.Start(context.Background())
		}()
	}
	wg.Wait()
	close(errs)

	for err := range errs {
		if err != nil {
			t.Fatalf("start runtime concurrently: %v", err)
		}
	}
	if engine.Starts() != 1 {
		t.Fatalf("expected engine to start once, got %d", engine.Starts())
	}
}

type countingEngine struct {
	mu     sync.Mutex
	name   string
	starts int
}

func (e *countingEngine) Name() string {
	return e.name
}

func (e *countingEngine) Dependencies() []string {
	return nil
}

func (e *countingEngine) Start(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	e.starts++
	return nil
}

func (e *countingEngine) Stop(ctx context.Context) error {
	return ctx.Err()
}

func (e *countingEngine) Health(ctx context.Context) error {
	return ctx.Err()
}

func (e *countingEngine) Starts() int {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.starts
}
