package main

import (
	"context"
	"fmt"
	"github.com/openshift-online/async-routine/opid"
	"log/slog"
	"os"
	"runtime/pprof"
	"time"
)

func main() {
	slog.Info("Program started")

	for i := 0; i < 10; i++ {
		foo(opid.NewContext())
	}

	// Wait enough time to have some routine started
	time.Sleep(4 * time.Second)

	pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
	fmt.Scanln()
}

func foo(ctx context.Context) {
	slog.Info("foo() started",
		"opid", opid.FromContext(ctx))
	bar(ctx)
	slog.Info("foo() ended",
		"opid", opid.FromContext(ctx))
}

func bar(ctx context.Context) {
	slog.Info("bar() started",
		"opid", opid.FromContext(ctx))
	go parentGoroutine(ctx)
	slog.Info("bar() ended",
		"opid", opid.FromContext(ctx))
}

func parentGoroutine(ctx context.Context) {
	slog.Info("parentGoroutine() started",
		"opid", opid.FromContext(ctx))

	go stuckInSelect()
	time.Sleep(500 * time.Millisecond)
	slog.Info("parentGoroutine() ended",
		"opid", opid.FromContext(ctx))
}

func stuckInSelect() {
	slog.Info("stuckInSelect() started")
	select {}
	slog.Info("stuckInSelect() ended")
}
