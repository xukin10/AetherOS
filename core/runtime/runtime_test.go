package runtime

import (
	"context"
	"testing"
)

func TestNewRuntime(t *testing.T) {
	rt := New()
	if rt.State() != StateCreated {
		t.Fatalf("expected StateCreated, got %s", rt.State())
	}
}

func TestRegisterEngine(t *testing.T) {
	rt := New()
	e := NewEngine("test", nil, nil)
	if err := rt.Register(e); err != nil {
		t.Fatalf("Register failed: %v", err)
	}
}

func TestRegisterDuplicate(t *testing.T) {
	rt := New()
	e := NewEngine("dup", nil, nil)
	_ = rt.Register(e)
	err := rt.Register(e)
	if err != ErrEngineAlreadyRegistered {
		t.Fatalf("expected ErrEngineAlreadyRegistered, got %v", err)
	}
}

func TestStartRuntime(t *testing.T) {
	rt := New()
	started := false
	e := NewEngine("test",
		func(_ context.Context) error {
			started = true
			return nil
		},
		nil,
	)
	_ = rt.Register(e)

	ctx := context.Background()
	if err := rt.Start(ctx); err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	if !started {
		t.Fatal("expected engine to be started")
	}
	if rt.State() != StateRunning {
		t.Fatalf("expected StateRunning, got %s", rt.State())
	}
}

func TestStopRuntime(t *testing.T) {
	rt := New()
	stopped := false
	e := NewEngine("test",
		nil,
		func(_ context.Context) error {
			stopped = true
			return nil
		},
	)
	_ = rt.Register(e)

	ctx := context.Background()
	_ = rt.Start(ctx)
	if err := rt.Stop(ctx); err != nil {
		t.Fatalf("Stop failed: %v", err)
	}
	if !stopped {
		t.Fatal("expected engine to be stopped")
	}
	if rt.State() != StateStopped {
		t.Fatalf("expected StateStopped, got %s", rt.State())
	}
}

func TestStopWithoutStart(t *testing.T) {
	rt := New()
	ctx := context.Background()
	err := rt.Stop(ctx)
	if err != ErrRuntimeNotStarted {
		t.Fatalf("expected ErrRuntimeNotStarted, got %v", err)
	}
}

func TestDoubleStart(t *testing.T) {
	rt := New()
	_ = rt.Register(NewEngine("e", nil, nil))
	ctx := context.Background()
	_ = rt.Start(ctx)
	err := rt.Start(ctx)
	if err != ErrRuntimeAlreadyStarted {
		t.Fatalf("expected ErrRuntimeAlreadyStarted, got %v", err)
	}
}

func TestEngineStartOrder(t *testing.T) {
	rt := New()
	var order []string

	e1 := NewEngine("a",
		func(_ context.Context) error { order = append(order, "a"); return nil },
		nil,
	)
	e2 := NewEngine("b",
		func(_ context.Context) error { order = append(order, "b"); return nil },
		nil,
	)
	_ = rt.Register(e1)
	_ = rt.Register(e2)

	ctx := context.Background()
	_ = rt.Start(ctx)

	if len(order) != 2 || order[0] != "a" || order[1] != "b" {
		t.Fatalf("expected start order [a b], got %v", order)
	}
}

func TestEngineStopOrder(t *testing.T) {
	rt := New()
	var order []string

	e1 := NewEngine("a",
		nil,
		func(_ context.Context) error { order = append(order, "a"); return nil },
	)
	e2 := NewEngine("b",
		nil,
		func(_ context.Context) error { order = append(order, "b"); return nil },
	)
	_ = rt.Register(e1)
	_ = rt.Register(e2)

	ctx := context.Background()
	_ = rt.Start(ctx)
	_ = rt.Stop(ctx)

	if len(order) != 2 || order[0] != "b" || order[1] != "a" {
		t.Fatalf("expected stop order [b a], got %v", order)
	}
}

func TestBootstrapBuild(t *testing.T) {
	e := NewEngine("boot", nil, nil)
	rt, err := Quick(e)
	if err != nil {
		t.Fatalf("Quick failed: %v", err)
	}
	if rt == nil {
		t.Fatal("expected non-nil runtime")
	}
	ctx := context.Background()
	if err := rt.Start(ctx); err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	if err := rt.Stop(ctx); err != nil {
		t.Fatalf("Stop failed: %v", err)
	}
}
