package runtime_test

import (
	"testing"

	"github.com/xukin10/AetherOS/core/runtime"
)

func TestRuntimeBootstrap(t *testing.T) {
	r := runtime.Bootstrap()

	if r == nil {
		t.Fatal("runtime is nil")
	}
}
