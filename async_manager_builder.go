package async

import cmap "github.com/orcaman/concurrent-map/v2"

type AsyncRoutineManagerBuilder struct {
}

func NewAsyncManagerBuilder() *AsyncRoutineManagerBuilder {
	return &AsyncRoutineManagerBuilder{}
}

func (b *AsyncRoutineManagerBuilder) Build() AsyncRoutineManager {
	manager := &asyncRoutineManager{
		routines:  cmap.New[AsyncRoutine](),
		observers: cmap.New[RoutinesObserver](),
	}
	manager.startMonitoring()
	return manager
}
