# STEP2 - integrating the AsyncRoutineManager

In this step, we integrate the `AsyncRoutineManager` to handle routine management.
In fact, incorporating the `AsyncRoutineManager` is almost as simple as replacing the `go` keyword with a call
to the `NewAsyncRoutine` function.
Running this code produces the same behavior as in STEP 1, but now we have established the foundation to 
fully leverage the capabilities of the `AsyncRoutineManager`.