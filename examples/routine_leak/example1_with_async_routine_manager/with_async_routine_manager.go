package main

import (
	"context"
	"github.com/openshift-online/async-routine"
	"github.com/openshift-online/async-routine/opid"
	"log/slog"
	"time"
)

var _ async.RoutinesObserver = (*exampleRoutineObserver)(nil)

func main() {
	slog.Info("Program started")

	// Setup the AsyncRoutineManager
	async.Manager(
		// Take a snapshot every 2 seconds
		async.WithSnapshottingInterval(2 * time.Second)).
		// Start the routine monitor
		Monitor().Start()

	// Add our custom observer to the list of routine observers
	_ = async.Manager().AddObserver(&exampleRoutineObserver{})

	for i := 0; i < 10; i++ {
		foo(opid.NewContext())
	}

	// Wait enough time to have some routine started
	time.Sleep(4 * time.Second)
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
	async.NewAsyncRoutine("parent go routine", ctx,
		func() {
			parentGoroutine(ctx)
		}).
		Timebox(2 * time.Second).
		Run()
	slog.Info("bar() ended",
		"opid", opid.FromContext(ctx))
}

func parentGoroutine(ctx context.Context) {
	slog.Info("parentGoroutine() started",
		"opid", opid.FromContext(ctx))

	async.NewAsyncRoutine("stuck in select", ctx, stuckInSelect).
		Timebox(2 * time.Second).
		Run()
	time.Sleep(500 * time.Millisecond)

	slog.Info("parentGoroutine() ended",
		"opid", opid.FromContext(ctx))
}

func stuckInSelect() {
	slog.Info("parentGoroutine() started")
	select {}
	slog.Info("parentGoroutine() ended")
}

type exampleRoutineObserver struct{}

func (e exampleRoutineObserver) RoutineStarted(routine async.AsyncRoutine) {
	slog.Info("Routine started",
		"name", routine.Name(),
		"opid", routine.OpId(),
		"parent-opid", routine.OriginatorOpId(),
	)
}

func (e exampleRoutineObserver) RoutineFinished(routine async.AsyncRoutine) {
	slog.Info("Routine finished",
		"name", routine.Name(),
		"opid", routine.OpId(),
		"parent-opid", routine.OriginatorOpId(),
	)
}

func (e exampleRoutineObserver) RoutineExceededTimebox(routine async.AsyncRoutine) {
	slog.Warn("Routine exceeded timebox",
		"name", routine.Name(),
		"opid", routine.OpId(),
		"parent-opid", routine.OriginatorOpId(),
		"startedAt", routine.StartedAt(),
	)
}

func (e exampleRoutineObserver) RunningRoutineCount(count int) {
	// nothing to do in this example
}

func (e exampleRoutineObserver) RunningRoutineByNameCount(name string, count int) {
	// nothing to do in this example
}
