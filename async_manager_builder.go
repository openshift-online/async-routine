package async

import cmap "github.com/orcaman/concurrent-map/v2"

type AsyncRoutineManagerBuilder struct {
	routineManager *asyncRoutineManager
}

func NewAsyncManagerBuilder() *AsyncRoutineManagerBuilder {
	routineManager := &asyncRoutineManager{
		routines:  cmap.New[AsyncRoutine](),
		observers: cmap.New[RoutinesObserver](),
	}

	routineManager.managerToggle = func() bool { return true }      // by default the manager is enabled
	routineManager.snapshottingToggle = func() bool { return true } // by default snapshotting is enabled

	return &AsyncRoutineManagerBuilder{
		routineManager: routineManager,
	}
}

func (b *AsyncRoutineManagerBuilder) WithManagerToggle(toggle Toggle) *AsyncRoutineManagerBuilder {
	b.routineManager.managerToggle = toggle
	return b
}

func (b *AsyncRoutineManagerBuilder) WithSnapshottingToggle(toggle Toggle) *AsyncRoutineManagerBuilder {
	b.routineManager.snapshottingToggle = toggle
	return b
}

func (b *AsyncRoutineManagerBuilder) Build() AsyncRoutineManager {
	b.routineManager.monitor().startMonitoring()
	return b.routineManager
}
