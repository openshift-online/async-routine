# STEP 5 - Adding some context 

The `AsyncRoutineManager` gives the ability to add arbitrary data to each routine so that each execution can be tied to a set of data.
We will leverage this ability to restrict even more our view to the leaked routine: we will timebox 
the routine to 1 second of execution and we add the received value as routine data.

```go
	async.NewAsyncRoutine("do-job", context.Background(),
        func() {
            getWebsiteResponseSize(url)
        }).
        Timebox(5*time.Second).
        WithData("url", url).
        Run()
    time.Sleep(500 * time.Millisecond)
```

Now the output will be similar to:
```
...
routine: [count: 6]
routine: [count: 6]
leaked routine: [name: do-job started-at: 2025-05-14 09:54:54.415189 +0000 UTC data: map[url:https://fool-bad.com]]
routine: [count: 4]
leaked routine: [name: do-job started-at: 2025-05-14 09:54:54.415189 +0000 UTC data: map[url:https://fool-bad.com]]
leaked routine: [name: do-job started-at: 2025-05-14 09:54:54.915739 +0000 UTC data: map[url:https://weather-bad.com]]
routine: [count: 4]
...
```
Analyzing the data, we observe that the routine leaks occur when there is an error with the website url.
Inspecting the code we see that `getWebsiteResponseSize` calls `getResponseSize` but ignores any error:
```go
    async.NewAsyncRoutine(
		"get-website-size",
		context.Background(),
		func() {
			getResponseSize(url, resultChan)
		}).Run()

    ...
    ...
    func getResponseSize(url string, ch chan<- int64) error {
```

Let's fix the code:
```go
    async.NewAsyncRoutine(
		"get-website-size",
		context.Background(),
        func() {
            err := getResponseSize(url, resultChan)
            if err != nil {
                resultChan <- -1
                fmt.Println("error fetching website size:", err)
            }
		}).Run()

    ...
    ...
    func getResponseSize(url string, ch chan<- int64) error {
```

Now the output is:
```go
routine: [count: 1]
Error getting response size: site unreachable: Get "https://duolingo-bad.com": dial tcp: lookup duolingo-bad.com: no such host
routine: [count: 1]
routine: [count: 1]
Error getting response size: site unreachable: Get "https://trulia-bad.com": dial tcp: lookup trulia-bad.com: no such host
routine: [count: 1]
routine: [count: 1]
routine: [count: 3]
routine: [count: 3]
routine: [count: 3]
routine: [count: 5]
routine: [count: 1]
routine: [count: 1]
routine: [count: 3]
routine: [count: 1]
routine: [count: 3]
routine: [count: 1]
routine: [count: 1]
routine: [count: 3]
routine: [count: 1]
routine: [count: 1]
routine: [count: 3]
routine: [count: 3]
routine: [count: 1]
routine: [count: 1]
routine: [count: 1]
routine: [count: 1]
```

and we don't have leaked routines anymore: we can say that because we don't have any routine exceeding the timebox
anymore and the routine count always drops to 1 (the only running routine: the `async-routine-monitor` one) 