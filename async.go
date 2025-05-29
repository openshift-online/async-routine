package async

import (
	"math"
	"time"
)

const (
	// DefaultRoutineSnapshottingInterval defines how often the routine manager checks routine status
	DefaultRoutineSnapshottingInterval = 30 * time.Second
	DefaultObserverTimeout             = time.Duration(math.MaxInt64)
)

// A RoutinesObserver is an object that observes the status of the executions of routines.
// The interface includes methods for notifying when a routine starts, finishes, times out,
// and for getting the count of running routines.
//
//go:generate mockgen -source=async.go -package=async -destination=mock_routine_observer.go
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

type routineEventType int

const (
	routineStarted routineEventType = iota
	routineEnded
	routineTimeboxExceeded
	takeSnapshot
)

type routineEvent struct {
	Type     routineEventType
	routine  AsyncRoutine
	snapshot Snapshot
}

func newRoutineEvent(eventType routineEventType) routineEvent {
	return routineEvent{
		Type: eventType,
	}
}

func routineStartedEvent() routineEvent {
	return newRoutineEvent(routineStarted)
}

func routineFinishedEvent() routineEvent {
	return newRoutineEvent(routineEnded)
}

func routineTimeboxExceededEvent() routineEvent {
	return newRoutineEvent(routineTimeboxExceeded)
}

func takeSnapshotEvent() routineEvent {
	return newRoutineEvent(takeSnapshot)
}
