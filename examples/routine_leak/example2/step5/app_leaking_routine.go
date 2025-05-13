package main

import (
	"context"
	"fmt"
	"github.com/openshift-online/async-routine"
	"github.com/openshift-online/async-routine/examples/routine_leak/example2/data"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func main() {
	async.Manager(async.WithSnapshottingInterval(500 * time.Millisecond)).Monitor().Start()
	async.Manager().AddObserver(&leakingRoutineObserver{})
	for {
		for i := 0; i < 10; i++ {
			url := data.Websites[rand.Intn(len(data.Websites))]
			async.NewAsyncRoutine("do-job", context.Background(),
				func() {
					doJob(url)
				}).
				Timebox(5*time.Second).
				WithData("url", url).
				Run()
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func doJob(url string) {
	resultChan := make(chan int64)
	async.NewAsyncRoutine(
		"get-website-size",
		context.Background(),
		func() {
			getResponseSize(url, resultChan)
		}).Run()
	size := <-resultChan
	// do something fun with the size - here we just avoid the compilation error
	size = size
}

// getResponseSize fetches the given URL and sends the response size (in bytes) to the provided channel.
// Returns an error if the site does not exist or the request fails.
func getResponseSize(url string, ch chan<- int64) error {
	// Perform the HTTP request
	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("site unreachable: %w", err)
	}
	defer res.Body.Close()

	// Handle HTTP errors
	if res.StatusCode != http.StatusOK {
		switch res.StatusCode {
		case http.StatusNotFound:
			return fmt.Errorf("site does not exist (404 Not Found)")
		default:
			return fmt.Errorf("invalid HTTP response: %d %s", res.StatusCode, http.StatusText(res.StatusCode))
		}
	}

	// Try to use Content-Length header if available
	if contentLength := res.Header.Get("Content-Length"); contentLength != "" {
		size, err := strconv.ParseInt(contentLength, 10, 64)
		if err == nil && size > 0 {
			ch <- size
			return nil
		}
	}

	// Calculate the size by reading the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	ch <- int64(len(body))
	return nil
}
