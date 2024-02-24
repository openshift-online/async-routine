package async

import "time"

type AsyncRoutineMonitor interface {
	startMonitoring()
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
