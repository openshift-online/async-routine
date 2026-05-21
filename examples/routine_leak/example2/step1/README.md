# STEP1 - the leaking code

In this step, we will use the native `go` keyword to start go routines.
The code is very simple:

1. The `main` function repeatedly starts cycles where it processes multiple websites.
   For each cycle, randomly selects 10 website URLs from a predefined list.
   For each selected URL, asynchronously invokes the `getWebsiteResponseSize`
2. The `getWebsiteResponseSize` asynchronously calls the `getResponseSize` and do some fun stuff with the result
3. The `getResponseSize` contacts the site and returns the response size

Running the application we can see how the total number of goroutine keeps increasing:
```
goroutine profile: total 18
14 @ 0x102d71768 0x102d361a8 0x102d70a10 0x102db7a38 0x102db830c 0x102db82fd 0x102e451d8 0x102e4d974 0x102e87430 0x102dcd820 0x102e87610 0x102e84dac 0x102e8ad44 0x102e8ad45 0x102eb83b4 0x102d9d940 0x102ee5b18 0x102ee5aed 0x102ee6148 0x102ef6b50 0x102ef60c8 0x102d79aa4
#	0x102d70a0f	internal/poll.runtime_pollWait+0x9f	
...
...
goroutine profile: total 33
26 @ 0x102d71768 0x102d361a8 0x102d70a10 0x102db7a38 0x102db830c 0x102db82fd 0x102e451d8 0x102e4d974 0x102e87430 0x102dcd820 0x102e87610 0x102e84dac 0x102e8ad44 0x102e8ad45 0x102eb83b4 0x102d9d940 0x102ee5b18 0x102ee5aed 0x102ee6148 0x102ef6b50 0x102ef60c8 0x102d79aa4
#	0x102d70a0f	internal/poll.runtime_pollWait+0x9f
```

Understanding which routine we are leaking is not immediate.
