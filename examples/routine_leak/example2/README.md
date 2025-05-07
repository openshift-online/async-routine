# Routine Leak Detection with AsyncRoutineManager
This demonstration illustrates how the `AsyncRoutineManager` help identifying routine leaks and pinpoint their source.

Each folder, named step1 to stepN, adds a progressive integration of the `AsyncRoutineManager`, starting from the 
naked code of `step1` to the full integration in the last step.

The application we are going to use for this demonstration is very simple:

1. The `main` function starts the `doJob` go routine and then, every 4 seconds, prints the `pproof` goroutine dump
2. `doJob` runs indefinitely and every 50 milliseconds invokes the `foo` function
3. `foo` just invokes `bar` wich in turn starts the `parentGoRoutine` routine
4. `parentGoRoutine` starts the `doInterestingStuff` routine passing a random number between 0 and 99
5. `doInterestingStuff` hangs indefinitely if the received value is a multiple of 4, otherwise exits