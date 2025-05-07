package main

import (
	"fmt"
	"github.com/openshift-online/async-routine"
)

var _ async.RoutinesObserver = (*leakingRoutineObserver)(nil)

type leakingRoutineObserver struct{}

func (l leakingRoutineObserver) RoutineStarted(routine async.AsyncRoutine) {}

func (l leakingRoutineObserver) RoutineFinished(routine async.AsyncRoutine) {}

func (l leakingRoutineObserver) RoutineExceededTimebox(routine async.AsyncRoutine) {
	fmt.Printf("leaked routine: [name: %s started-at: %v data: %s]\n", routine.Name(), routine.StartedAt(), routine.GetData())
}

func (l leakingRoutineObserver) RunningRoutineCount(count int) {
}

func (l leakingRoutineObserver) RunningRoutineByNameCount(name string, count int) {
}
