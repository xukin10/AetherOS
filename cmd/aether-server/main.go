package main

import (
	"context"
	"fmt"
	"os"

	"github.com/xukin10/AetherOS/core/runtime"
)

func main() {
	rt := runtime.NewRuntime()
	if err := rt.Bootstrap(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "bootstrap runtime: %v\n", err)
		os.Exit(1)
	}
}
