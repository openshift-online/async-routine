package async

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	cmap "github.com/orcaman/concurrent-map/v2"
)

var routineManager AsyncRoutineManager

type AsyncRoutineManager interface {
	AddObserver(observer RoutinesObserver) string
	RemoveObserver(observerId string)
	IsEnabled() bool
	GetSnapshot() Snapshot
	notifyAll(src AsyncRoutine, evt routineEvent)
	Monitor() AsyncRoutineMonitor
	register(routine AsyncRoutine)
	deregister(routine AsyncRoutine)
}

type Toggle func() bool

type asyncRoutineManager struct {
	ctx                  context.Context
	managerToggle        Toggle
	snapshottingToggle   Toggle
	snapshottingInterval time.Duration
	routines             cmap.ConcurrentMap[string, AsyncRoutine]
	observers            cmap.ConcurrentMap[string, *observerProxy]

	monitorLock sync.Mutex // user to sync the `Start` and `Stop` methods that are used to start the
	// snapshotting routine
	monitorStopChannel  chan bool // Used to notify the snapshotting routine that it should stop running
	monitorStarted      bool
	monitorWaitingGroup sync.WaitGroup // Used by the `Stop` method to wait for the snapshotting routine to end
}

func (arm *asyncRoutineManager) IsEnabled() bool {
	return arm.managerToggle()
}

// AddObserver adds a new RoutineObserver to the list of observers.
// Assigns and returns an observer ID to the RoutineObserver
func (arm *asyncRoutineManager) AddObserver(observer RoutinesObserver) string {
	return arm.AddObserverWithTimeout(observer, DefaultObserverTimeout)
}

// AddObserverWithTimeout registers a new RoutinesObserver with the asyncRoutineManager,
// associating it with a unique identifier and a specified timeout duration.
// The function returns the unique ID assigned to the observer.
func (arm *asyncRoutineManager) AddObserverWithTimeout(observer RoutinesObserver, timeout time.Duration) string {
	uid := uuid.New().String()
	proxy := newObserverProxy(uid, observer, arm, timeout)
	arm.observers.Set(uid, proxy)
	proxy.startObserving()
	return uid
}

// RemoveObserver removes the given RoutineObserver from the list of observers
func (arm *asyncRoutineManager) RemoveObserver(observerId string) {
	observer, ok := arm.observers.Get(observerId)
	if !ok {
		return
	}
	observer.stopObserving()
	arm.observers.Remove(observerId)
}

func (arm *asyncRoutineManager) GetSnapshot() Snapshot {
	snapshot := newSnapshot()
	for runningRoutines := range arm.routines.IterBuffered() {
		snapshot.registerRoutine(runningRoutines.Val)
	}
	return snapshot
}

// notifyAll notifies all the observers of the event evt received from the routine src
func (arm *asyncRoutineManager) notifyAll(src AsyncRoutine, evt routineEvent) {
	for _, observer := range arm.observers.Items() {
		observer.notify(src, evt)
	}
}

func (arm *asyncRoutineManager) Monitor() AsyncRoutineMonitor {
	return arm
}

func (arm *asyncRoutineManager) register(routine AsyncRoutine) {
	arm.routines.Set(routine.id(), routine)
}

func (arm *asyncRoutineManager) deregister(routine AsyncRoutine) {
	arm.routines.Remove(routine.id())
}

type AsyncManagerOption func(mgr *asyncRoutineManager)

func WithSnapshottingInterval(interval time.Duration) AsyncManagerOption {
	return func(mgr *asyncRoutineManager) {
		mgr.snapshottingInterval = interval
	}
}

func WithSnapshottingToggle(toggle Toggle) AsyncManagerOption {
	return func(mgr *asyncRoutineManager) {
		mgr.snapshottingToggle = toggle
	}
}

func WithManagerToggle(toggle Toggle) AsyncManagerOption {
	return func(mgr *asyncRoutineManager) {
		mgr.managerToggle = toggle
	}
}

func WithContext(ctx context.Context) AsyncManagerOption {
	return func(mgr *asyncRoutineManager) {
		mgr.ctx = ctx
	}
}

var lock sync.RWMutex

func newAsyncRoutineManager(options ...AsyncManagerOption) AsyncRoutineManager {
	mgr := &asyncRoutineManager{
		routines:             cmap.New[AsyncRoutine](),
		observers:            cmap.New[*observerProxy](),
		snapshottingInterval: DefaultRoutineSnapshottingInterval,
		ctx:                  context.Background(),
		managerToggle:        func() bool { return true }, // manager is enabled by default
		snapshottingToggle:   func() bool { return true }, // snapshotting is enabled by default
		monitorStopChannel:   make(chan bool),
	}

	for _, option := range options {
		option(mgr)
	}

	return mgr
}

func Manager(options ...AsyncManagerOption) AsyncRoutineManager {
	if routineManager == nil {
		lock.Lock()
		defer lock.Unlock()
		if routineManager == nil {
			routineManager = newAsyncRoutineManager(options...)
			return routineManager
		}
	}

	for _, option := range options {
		option(routineManager.(*asyncRoutineManager))
	}
	return routineManager
}
