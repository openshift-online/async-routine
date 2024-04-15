package async

import (
	"time"
)

type AsyncRoutineMonitor interface {
	startMonitoring()
	IsSnapshottingEnabled() bool
	IsStarted() bool
	Start()
	Stop()
}

var _ AsyncRoutineMonitor = (*asyncRoutineManager)(nil)

func (arm *asyncRoutineManager) IsSnapshottingEnabled() bool {
	return arm.snapshottingToggle()
}

func (arm *asyncRoutineManager) startMonitoring() {
	NewAsyncRoutine("async-routine-monitor", arm.ctx, func() {
		ticker := time.NewTicker(arm.snapshottingInterval)
		defer ticker.Stop()

		for {
			select {
			case <-arm.monitorStopChannel:
				arm.monitorStarted = false
				arm.monitorWaitingGroup.Done()
				return
			case <-ticker.C:
				if arm.IsEnabled() && arm.IsSnapshottingEnabled() {
					arm.snapshot()
				}
			}
		}
	}).Run()
}

func (arm *asyncRoutineManager) IsStarted() bool {
	return arm.monitorStarted
}

func (arm *asyncRoutineManager) Stop() {
	arm.monitorLock.Lock()
	defer arm.monitorLock.Unlock()
	if arm.IsStarted() {
		arm.monitorStopChannel <- true
	}
	arm.monitorWaitingGroup.Wait()
}

func (arm *asyncRoutineManager) Start() {
	arm.monitorLock.Lock()
	defer arm.monitorLock.Unlock()
	if !arm.IsStarted() {
		arm.monitorStarted = true
		arm.startMonitoring()
		arm.monitorWaitingGroup.Add(1)
	}
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
