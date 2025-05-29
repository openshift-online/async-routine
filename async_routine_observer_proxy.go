package async

import (
	"context"
	"time"
)

// observerProxy acts as an intermediary between the AsyncRoutineManager and a RoutinesObserver.
// It receives routine events via a channel and dispatches them to the observer's callback methods.
// The proxy manages event notification asynchronously and can enforce a timeout on the observer's lifecycle.
type observerProxy struct {
	manager    AsyncRoutineManager
	observerId string
	observer   RoutinesObserver
	channel    chan routineEvent
	timeout    time.Duration
}

// newObserverProxy creates and initializes a new observerProxy instance.
// It sets up an asynchronous routine that listens for routine events on the proxy's channel
// and forwards them to the appropriate methods of the provided RoutinesObserver.
//
// Parameters:
//   - observerId: a unique identifier for the observer instance.
//   - observer: the RoutinesObserver to be notified of routine events.
//   - manager: the AsyncRoutineManager responsible for managing routines.
//   - timeout: the duration after which the observer routine is considered 'exceeding the timebox'.
//
// Returns:
//   - A pointer to the initialized observerProxy.
func newObserverProxy(observerId string, observer RoutinesObserver, manager AsyncRoutineManager, timeout time.Duration) *observerProxy {
	proxy := &observerProxy{
		manager:    manager,
		observerId: observerId,
		observer:   observer,
		channel:    make(chan routineEvent),
		timeout:    timeout,
	}

	return proxy
}

func (proxy *observerProxy) startObserving() {
	NewAsyncRoutine("async-observer-notifier", context.Background(), func() {
		for evt := range proxy.channel {
			switch evt.Type {
			case routineStarted:
				proxy.observer.RoutineStarted(evt.routine)
			case routineEnded:
				proxy.observer.RoutineFinished(evt.routine)
			case routineTimeboxExceeded:
				proxy.observer.RoutineExceededTimebox(evt.routine)
			case takeSnapshot:
				proxy.observer.RunningRoutineCount(evt.snapshot.GetTotalRoutineCount())
				for _, routineName := range evt.snapshot.GetRunningRoutinesNames() {
					proxy.observer.RunningRoutineByNameCount(routineName, evt.snapshot.GetRunningRoutinesCount(routineName))
				}
			}
		}
	}).
		Timebox(proxy.timeout).
		WithData("observer-id", proxy.observerId).
		Run()
}

func (proxy *observerProxy) stopObserving() {
	close(proxy.channel)
}

// notify sends a routine event to the observerProxy's channel.
// Depending on the event type, it either forwards the routine information
// or triggers a snapshot retrieval from the manager.
func (proxy *observerProxy) notify(routine AsyncRoutine, evt routineEvent) {
	switch evt.Type {
	case routineStarted, routineEnded, routineTimeboxExceeded:
		proxy.channel <- routineEvent{
			Type:    evt.Type,
			routine: routine,
		}
	case takeSnapshot:
		proxy.channel <- routineEvent{
			Type:     takeSnapshot,
			snapshot: proxy.manager.GetSnapshot(),
		}
	}
}
