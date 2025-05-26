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
running routine count: 2
running routine count: 3
running routine count: 3
running routine count: 2
running routine count: 4
running routine count: 4
running routine count: 2
running routine count: 4
running routine count: 4
running routine count: 2
running routine count: 3
running routine count: 5
running routine count: 3
running routine count: 3
running routine count: 5
running routine count: 5
running routine count: 3
running routine count: 4
running routine count: 4
running routine count: 4
running routine count: 5
running routine count: 5
running routine count: 7
```

We can see we have a routine leak, but we still have no idea of what is the routine leak.