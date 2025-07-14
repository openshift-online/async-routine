package async

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

// Event types for routine lifecycle and snapshotting
type routineLifeCycleEventType = int
type routineSnapshottingEventType = int

const (
	routineStarted routineLifeCycleEventType = iota
	routineEnded
	routineTimeboxExceeded
)

const (
	routineCount routineSnapshottingEventType = iota
	routineByNameCount
)

// DefaultAsyncObserverBufferSize defines the default buffer size for the event channel
// used by async routine observers. This value determines how many events can be queued
// before send operations block or, if not guaranteed, events are dropped.
const DefaultAsyncObserverBufferSize = 16

type asyncRoutineObserverOption func(o *asyncRoutineObserver)

// WithObserverTimeout sets a timeout for the observer routine.
func WithObserverTimeout(timeout time.Duration) func(o *asyncRoutineObserver) {
	return func(o *asyncRoutineObserver) {
		o.timeout = &timeout
	}
}

// WithObserverRoutineData adds custom data to the observer routine.
func WithObserverRoutineData(key string, value string) func(o *asyncRoutineObserver) {
	return func(o *asyncRoutineObserver) {
		if o.observerNotifierRoutineData == nil {
			o.observerNotifierRoutineData = make(routineData)
		}
		o.observerNotifierRoutineData[key] = value
	}
}

// WithChannelSize sets the buffer size for the async observer's event channel.
// Use this as an option when creating the async observer.
func WithChannelSize(size int) func(o *asyncRoutineObserver) {
	return func(o *asyncRoutineObserver) {
		if o.channel != nil {
			close(o.channel)
		}
		o.channel = make(chan asyncRoutineEvent, size)
	}
}

// WithGuaranteedDelivery enables guaranteed delivery mode for the async observer's events.
// When enabled, the observer will not drop events if the channel is full, but will block until space is available.
// Use this as an option when creating the observer.
func WithGuaranteedDelivery() func(o *asyncRoutineObserver) {
	return func(o *asyncRoutineObserver) {
		o.guaranteedDelivery = true
	}
}

func NewAsyncRoutineObserver(delegate RoutinesObserver, options ...asyncRoutineObserverOption) RoutinesObserver {
	asyncObserver := asyncRoutineObserver{
		delegate:           delegate,
		guaranteedDelivery: false,
	}

	asyncObserver.channelIsOpen.Store(true)

	for _, option := range options {
		option(&asyncObserver)
	}

	if asyncObserver.channel == nil {
		asyncObserver.channel = make(chan asyncRoutineEvent, DefaultAsyncObserverBufferSize)
	}

	asyncObserver.startObserving()
	return &asyncObserver
}

var _ RoutinesObserver = (*asyncRoutineObserver)(nil)

type asyncRoutineObserver struct {
	delegate                    RoutinesObserver
	channel                     chan asyncRoutineEvent
	channelIsOpen               atomic.Bool
	timeout                     *time.Duration
	observerNotifierRoutineData routineData
	closeOnce                   sync.Once
	guaranteedDelivery          bool
}

func (a *asyncRoutineObserver) RoutineStarted(routine AsyncRoutine) {
	a.notify(&routineLifecycleEvent{
		evtType: routineStarted,
		routine: routine,
	})
}

func (a *asyncRoutineObserver) RoutineFinished(routine AsyncRoutine) {
	a.notify(&routineLifecycleEvent{
		evtType: routineEnded,
		routine: routine,
	})
}

func (a *asyncRoutineObserver) RoutineExceededTimebox(routine AsyncRoutine) {
	a.notify(&routineLifecycleEvent{
		evtType: routineTimeboxExceeded,
		routine: routine,
	})
}

func (a *asyncRoutineObserver) RunningRoutineCount(count int) {
	a.notify(&routineCountEvent{
		count: count,
	})
}

func (a *asyncRoutineObserver) RunningRoutineByNameCount(name string, count int) {
	a.notify(&routineByNameCountEvent{
		routineName: name,
		count:       count,
	})
}

func (a *asyncRoutineObserver) manageLifecycleEvent(evt *routineLifecycleEvent) {
	switch evt.evtType {
	case routineStarted:
		a.delegate.RoutineStarted(evt.routine)
	case routineEnded:
		a.delegate.RoutineFinished(evt.routine)
	case routineTimeboxExceeded:
		a.delegate.RoutineExceededTimebox(evt.routine)
	}
}

func (a *asyncRoutineObserver) startObserving() {
	routine := NewAsyncRoutine("async-observer-notifier", context.Background(), func() {
		for evt := range a.channel {
			switch evt.(type) {
			case *routineLifecycleEvent:
				a.manageLifecycleEvent(evt.(*routineLifecycleEvent))
			case *routineCountEvent:
				a.delegate.RunningRoutineCount(evt.(*routineCountEvent).count)
			case *routineByNameCountEvent:
				evt := evt.(*routineByNameCountEvent)
				a.delegate.RunningRoutineByNameCount(evt.routineName, evt.count)
			}
		}
	})

	for key, value := range a.observerNotifierRoutineData {
		routine = routine.WithData(key, value)
	}

	if a.timeout != nil {
		routine = routine.Timebox(*a.timeout)
	}

	routine.Run()
}

func (a *asyncRoutineObserver) stopObserving() {
	a.closeOnce.Do(func() {
		a.channelIsOpen.Store(false)
		close(a.channel)
	})
}

// notify sends an event to the observer channel, non-blocking if closed.
func (a *asyncRoutineObserver) notify(evt asyncRoutineEvent) {
	if a.guaranteedDelivery {
		if a.channelIsOpen.Load() {
			a.channel <- evt
		}
		// channel is closed. Event is lost. This happens when `stopObserving` is called.
		return
	}

	select {
	case a.channel <- evt:
	default:
		// Event dropped
	}
}

// --- Event types and interfaces ---

type asyncRoutineEvent interface {
	getEventType() int
}

var _ asyncRoutineEvent = (*routineLifecycleEvent)(nil)
var _ asyncRoutineEvent = (*routineCountEvent)(nil)
var _ asyncRoutineEvent = (*routineByNameCountEvent)(nil)

type routineLifecycleEvent struct {
	evtType routineLifeCycleEventType
	routine AsyncRoutine
}

func (r *routineLifecycleEvent) getEventType() int {
	return r.evtType
}

type routineCountEvent struct {
	count int
}

func (r routineCountEvent) getEventType() int {
	return routineCount
}

type routineByNameCountEvent struct {
	routineName string
	count       int
}

func (r routineByNameCountEvent) getEventType() int {
	return routineByNameCount
}
