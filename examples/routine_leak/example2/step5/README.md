# STEP 5 - Adding some context 

The `AsyncRoutineManager` gives the ability to add arbitrary data to each routine so that each execution can be tied to a set of data.
We will leverage this ability to restrict even more our view to the leaked routine: we will timebox 
the routine to 1 second of execution and we add the received value as routine data.

```go
    async.
		NewAsyncRoutine("run-command", context.Background(), func() { doInterestingStuff(data) }).
		WithData("data", fmt.Sprintf("%d", data)).
		Timebox(1 * time.Second).
		Run()
```

Now the output will be similar to:
```
leaked routine: [name: run-command started-at: 2025-05-07 15:35:58.218663 +0000 UTC data: map[data:52]]
leaked routine: [name: run-command started-at: 2025-05-07 15:35:59.038886 +0000 UTC data: map[data:56]]
leaked routine: [name: run-command started-at: 2025-05-07 15:35:58.218663 +0000 UTC data: map[data:52]]
leaked routine: [name: run-command started-at: 2025-05-07 15:35:58.987592 +0000 UTC data: map[data:72]]
leaked routine: [name: run-command started-at: 2025-05-07 15:35:58.834007 +0000 UTC data: map[data:16]]
leaked routine: [name: run-command started-at: 2025-05-07 15:35:58.885224 +0000 UTC data: map[data:84]]
leaked routine: [name: run-command started-at: 2025-05-07 15:35:58.680921 +0000 UTC data: map[data:44]]
leaked routine: [name: run-command started-at: 2025-05-07 15:35:59.038886 +0000 UTC data: map[data:56]]
```
Analyzing the data, we observe that the routine leaks occur only when the data is even.
Upon closer inspection, it becomes evident that the leaks happen specifically when the data is divisible by 4.