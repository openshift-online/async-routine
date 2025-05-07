package main

import (
	"context"
	"github.com/openshift-online/async-routine"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"
)

func main() {
	async.
		NewAsyncRoutine("main-job", context.Background(), doJob).
		Run()

	for {
		time.Sleep(4 * time.Second)
		_ = pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
	}

}

func doJob() {
	for {
		foo()
		time.Sleep(50 * time.Millisecond)
	}
}

func foo() {
	bar()
}

func bar() {
	async.
		NewAsyncRoutine("parent-go-routine", context.Background(), parentGoroutine).
		Run()
}

func parentGoroutine() {
	async.
		NewAsyncRoutine("run-command", context.Background(), func() { doInterestingStuff(rand.Intn(100)) }).
		Run()
	time.Sleep(500 * time.Millisecond)
}

func doInterestingStuff(value int) {
	if value%4 == 0 {
		select {}
	}
}
