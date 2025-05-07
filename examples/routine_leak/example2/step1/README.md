# STEP1 - the leaking code

In this step, we will use the native `go` keyword to start go routines.
The code is very simple:

1. The `main` function starts the `doJob` go routine and then, every 4 seconds, prints the `pproof` goroutine dump
2. `doJob` runs indefinitely and every 50 milliseconds invokes the `foo` function
3. `foo` just invokes `bar` wich in turn starts the `parentGoRoutine` routine
4. `parentGoRoutine` starts the `doInterestingStuff` routine passing a random number between 0 and 99
5. `doInterestingStuff` hangs indefinitely if the received value is a multiple of 4, otherwise exits

Running the application we can see how the total number of goroutine keeps increasing:
```
goroutine profile: total 32
20 @ 0x100e9caa8 0x100e7d520 0x100ef588c 0x100ea42b4
#	0x100ef588b	main.doInterestingStuff+0x3b
...
...
goroutine profile: total 50
38 @ 0x100e9caa8 0x100e7d520 0x100ef588c 0x100ea42b4
#	0x100ef588b	main.doInterestingStuff+0x3b
```

We can see that the leaked routine os the `doInterestingStuff` but we have no idea of the context of execution.