package main

import (
	"fmt"
	"github.com/openshift-online/async-routine"
)

var _ async.RoutinesObserver = (*leakingRoutineObserver)(nil)

type leakingRoutineObserver struct{}

func (l leakingRoutineObserver) RoutineStarted(routine async.AsyncRoutine) {}

func (l leakingRoutineObserver) RoutineFinished(routine async.AsyncRoutine) {}

func (l leakingRoutineObserver) RoutineExceededTimebox(routine async.AsyncRoutine) {}

func (l leakingRoutineObserver) RunningRoutineCount(count int) {
	fmt.Println("running routine count:", count)
}

func (l leakingRoutineObserver) RunningRoutineByNameCount(name string, count int) {
	fmt.Printf("Routine count for %s: %d\n", name, count)
}
