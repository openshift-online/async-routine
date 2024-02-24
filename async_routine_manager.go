package async

import (
	"github.com/google/uuid"
	cmap "github.com/orcaman/concurrent-map/v2"
)

type AsyncRoutineManager interface {
	AddObserver(observer RoutinesObserver) string
	RemoveObserver(observerId string)
	IsEnabled() bool
	Run(routine AsyncRoutine, routines ...AsyncRoutine)
	notify(eventSource func(observer RoutinesObserver))
	monitor() AsyncRoutineMonitor
}

var _ AsyncRoutineManager = (*asyncRoutineManager)(nil)

type Toggle func() bool

type asyncRoutineManager struct {
	managerToggle Toggle
	routines      cmap.ConcurrentMap[string, AsyncRoutine]
	observers     cmap.ConcurrentMap[string, RoutinesObserver]
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

func (arm *asyncRoutineManager) notify(eventSource func(observer RoutinesObserver)) {
	for _, observer := range arm.observers.Items() {
		eventSource(observer)
	}
}

func (arm *asyncRoutineManager) monitor() AsyncRoutineMonitor {
	return arm
}

func (arm *asyncRoutineManager) Run(routine AsyncRoutine, routines ...AsyncRoutine) {
	for _, r := range append(routines, routine) {
		if arm.IsEnabled() {
			arm.routines.Set(uuid.New().String(), r)
		}
		r.run(arm)
	}
}
