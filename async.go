package async

import (
	"time"
)

const (
	// DefaultRoutineSnapshottingInterval defines how often the routine manager checks routine status
	DefaultRoutineSnapshottingInterval = 30 * time.Second
)

// A RoutinesObserver is an object that observes the status of the executions of routines.
// The interface includes methods for notifying when a routine starts, finishes, times out,
// and for getting the count of running routines.
type RoutinesObserver interface {
	RoutineStarted(routine AsyncRoutine)

	RoutineFinished(routine AsyncRoutine)

	RoutineExceededTimebox(routine AsyncRoutine)

	// RunningRoutineCount is called to notify the observer about the total number of managed routines that are
	// currently running
	RunningRoutineCount(count int)

	// RunningRoutineByNameCount is called to notify the observer about how many routines with a given name are
	// currently running
	RunningRoutineByNameCount(name string, count int)
}

