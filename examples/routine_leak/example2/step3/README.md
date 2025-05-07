# STEP3 - Replacing the pproof dump with a Routine Observer

In this example, we implement an observer that prints the number of running routine:
```go

type leakingRoutineObserver struct{}

func (l leakingRoutineObserver) RunningRoutineCount(count int) {
	fmt.Println("running routine count:", count)
}
...
...
```

Next, to ensure the observer is automatically notified of all running routines, we need to register the 
observer and start the monitor.

```go
func main() {
    async.Manager(async.WithSnapshottingInterval(500 * time.Millisecond)).Monitor().Start()
    async.Manager().AddObserver(&leakingRoutineObserver{})
	...
	...
}
```

The output will be something like this: 
```
running routine count: 14
running routine count: 17
running routine count: 19
running routine count: 21
running routine count: 23
running routine count: 28
running routine count: 30
running routine count: 31
...
```

We can see we have a routine leak, but we still have no idea of what is the routine leak.