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
	notify(eventSource func(observer RoutinesObserver))
	Monitor() AsyncRoutineMonitor
	run(routine AsyncRoutine)
}

type Toggle func() bool

type asyncRoutineManager struct {
	ctx                  context.Context
	managerToggle        Toggle
	snapshottingToggle   Toggle
	snapshottingInterval time.Duration
	routines             cmap.ConcurrentMap[string, AsyncRoutine]
	observers            cmap.ConcurrentMap[string, RoutinesObserver]

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
	uid := uuid.New().String()
	arm.observers.Set(uid, observer)
	return uid
}

// RemoveObserver removes the given RoutineObserver from the list of observers
func (arm *asyncRoutineManager) RemoveObserver(observerId string) {
	arm.observers.Remove(observerId)
}

func (arm *asyncRoutineManager) GetSnapshot() Snapshot {
	snapshot := newSnapshot()
	for runningRoutines := range arm.routines.IterBuffered() {
		snapshot.registerRoutine(runningRoutines.Val)
	}
	return snapshot
}

func (arm *asyncRoutineManager) notify(eventSource func(observer RoutinesObserver)) {
	for _, observer := range arm.observers.Items() {
		eventSource(observer)
	}
}

func (arm *asyncRoutineManager) Monitor() AsyncRoutineMonitor {
	return arm
}

func (arm *asyncRoutineManager) run(routine AsyncRoutine) {
	if arm.IsEnabled() {
		arm.routines.Set(uuid.New().String(), routine)
	}
	routine.run(arm)
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

func Manager(options ...AsyncManagerOption) AsyncRoutineManager {
	if routineManager == nil {
		lock.Lock()
		defer lock.Unlock()
		if routineManager == nil {
			mgr := &asyncRoutineManager{
				routines:             cmap.New[AsyncRoutine](),
				observers:            cmap.New[RoutinesObserver](),
				snapshottingInterval: DefaultRoutineSnapshottingInterval,
				ctx:                  context.Background(),
				managerToggle:        func() bool { return true }, // manager is enabled by default
				snapshottingToggle:   func() bool { return true }, // snapshotting is enabled by default
				monitorStopChannel:   make(chan bool),
			}
			routineManager = mgr
		}
	}

	for _, option := range options {
		option(routineManager.(*asyncRoutineManager))
	}
	return routineManager
}
