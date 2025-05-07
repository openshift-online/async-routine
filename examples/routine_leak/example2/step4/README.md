# STEP 4 - Identify which routine is being leaked

In this step, we extend the observer by implementing a method to track the number of active instances
for each routine type at any given time:

```go
type leakingRoutineObserver struct{}
...
...
func (l leakingRoutineObserver) RunningRoutineByNameCount(name string, count int) {
	fmt.Printf("Routine count for %s: %d\n", name, count)
}
...
...
```

The output will be similar to this:
```
running routine count: 30
Routine count for run-command: 18
Routine count for async-routine-monitor: 1
Routine count for parent-go-routine: 10
Routine count for main-job: 1
running routine count: 31
Routine count for run-command: 20
Routine count for parent-go-routine: 9
Routine count for async-routine-monitor: 1
Routine count for main-job: 1
running routine count: 37
Routine count for run-command: 25
Routine count for parent-go-routine: 10
Routine count for async-routine-monitor: 1
Routine count for main-job: 1
```

Now it is clear: we are leaking the `run-command` routine. However, we don't have any context, 
so it is hard to understand why and when the routine is leaked.