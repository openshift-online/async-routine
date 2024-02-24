package async

import (
	"time"

	"github.com/google/uuid"
	cmap "github.com/orcaman/concurrent-map/v2"
)

const (
	// routineMonitoringDelay defines how often the routine manager checks routine status
	routineMonitoringDelay = 30 * time.Second
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

type AsyncRoutineManager interface {
	AddObserver(observer RoutinesObserver) string
	RemoveObserver(observerId string)
	Run(routine AsyncRoutine, routines ...AsyncRoutine)

	notify(eventSource func(observer RoutinesObserver))
}

type asyncRoutineManager struct {
	routines  cmap.ConcurrentMap[string, AsyncRoutine]
	observers cmap.ConcurrentMap[string, RoutinesObserver]
}

// AddObserver adds a new RoutineObserver to the list of observers.
// Assigns and returns an observer ID to the RoutineObserver
func (arm *asyncRoutineManager) AddObserver(observer RoutinesObserver) string {
	uid := uuid.New().String()
	arm.observers.Set(uid, observer)
	return uid
}

// RemoveObserver removes the given RoutineObserver from the list of observers
func (arm *asyncRoutineManager) RemoveObserver(observerId string) {
	arm.observers.Remove(observerId)
}

func (arm *asyncRoutineManager) notify(eventSource func(observer RoutinesObserver)) {
	for _, observer := range arm.observers.Items() {
		eventSource(observer)
	}
}

func (arm *asyncRoutineManager) Run(routine AsyncRoutine, routines ...AsyncRoutine) {
	for _, r := range append(routines, routine) {
		arm.routines.Set(uuid.New().String(), r)
		r.run(arm)
	}
}

func (arm *asyncRoutineManager) startMonitoring() {
	go func() {
		ticker := time.NewTicker(routineMonitoringDelay)
		defer ticker.Stop()

		for {
			<-ticker.C
			arm.snapshot()
		}
	}()
}

func (arm *asyncRoutineManager) snapshot() {
	runningThreads := 0
	runningThreadByName := map[string]int{}
	for monitorItem := range arm.routines.IterBuffered() {
		id := monitorItem.Key
		thread := monitorItem.Val
		switch {
		case thread.hasExceededTimebox():
			arm.notify(func(observer RoutinesObserver) {
				observer.RoutineExceededTimebox(thread)
			})
		case thread.isFinished():
			arm.routines.Remove(id)
		}

		if thread.isRunning() {
			runningThreads++
			count := runningThreadByName[thread.Name()]
			count++
			runningThreadByName[thread.Name()] = count
		}
	}

	arm.notify(func(observer RoutinesObserver) {
		observer.RunningRoutineCount(runningThreads)
		for name, count := range runningThreadByName {
			observer.RunningRoutineByNameCount(name, count)
		}
	})
}
