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
...
running routine count: 10
Routine count for do-job: 9
Routine count for async-routine-monitor: 1
running routine count: 11
Routine count for do-job: 10
Routine count for async-routine-monitor: 1
running routine count: 11
Routine count for do-job: 10
Routine count for async-routine-monitor: 1
running routine count: 11
Routine count for do-job: 10
Routine count for async-routine-monitor: 1
running routine count: 11
Routine count for do-job: 10
Routine count for async-routine-monitor: 1
running routine count: 11
Routine count for do-job: 10
Routine count for async-routine-monitor: 1
running routine count: 11
Routine count for do-job: 10
Routine count for async-routine-monitor: 1
running routine count: 11
Routine count for do-job: 10
Routine count for async-routine-monitor: 1
running routine count: 11
Routine count for do-job: 10
Routine count for async-routine-monitor: 1
```
The `async-routine-montor` is the routine the manager started when we started the monitor.
Now it is clear: we are leaking the `do-job` routine. However, we don't have any context, 
so it is hard to understand why and when the routine is leaked.