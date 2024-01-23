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

var routines = cmap.New[*asyncRoutine]()
var observers = cmap.New[RoutinesObserver]()

// AddObserver adds a new RoutineObserver to the list of observers.
// Assigns and returns an observer ID to the RoutineObserver
func AddObserver(observer RoutinesObserver) string {
	uid := uuid.New().String()
	observers.Set(uid, observer)
	return uid
}

// RemoveObserver removes the given RoutineObserver from the list of observers
func RemoveObserver(observerId string) {
	observers.Remove(observerId)
}

func sendEvent(eventSource func(observer RoutinesObserver)) {
	for _, observer := range observers.Items() {
		eventSource(observer)
	}
}

func startMonitoring() {
	go func() {
		for {
			time.Sleep(routineMonitoringDelay)
			runningThreads := 0
			runningThreadByName := map[string]int{}
			for id, thread := range routines.Items() {
				switch {
				case thread.hasExceededTimebox():
					sendEvent(func(observer RoutinesObserver) {
						thread.status = RoutineStatusExceededTimebox
						observer.RoutineExceededTimebox(thread)
					})
				case thread.isFinished():
					sendEvent(func(observer RoutinesObserver) {
						observer.RoutineFinished(thread)
					})
					routines.Remove(id)
				}

				if thread.isRunning() {
					runningThreads++
					count := runningThreadByName[thread.name]
					count++
					runningThreadByName[thread.name] = count
				}
			}

			sendEvent(func(observer RoutinesObserver) {
				observer.RunningRoutineCount(runningThreads)
				for name, count := range runningThreadByName {
					observer.RunningRoutineByNameCount(name, count)
				}
			})
		}
	}()
}

func init() {
	startMonitoring()
}
