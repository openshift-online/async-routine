package main

import (
	"context"
	"fmt"
	"github.com/openshift-online/async-routine"
	"math/rand"
	"time"
)

func main() {
	async.Manager(async.WithSnapshottingInterval(500 * time.Millisecond)).Monitor().Start()
	async.Manager().AddObserver(&leakingRoutineObserver{})

	async.
		NewAsyncRoutine("main-job", context.Background(), doJob).
		Run()

	// wait for enter to be pressed
	_, _ = fmt.Scanln()
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
	data := rand.Intn(100)
	async.
		NewAsyncRoutine("run-command", context.Background(), func() { doInterestingStuff(data) }).
		WithData("data", fmt.Sprintf("%d", data)).
		Timebox(1 * time.Second).
		Run()
	time.Sleep(500 * time.Millisecond)
}

func doInterestingStuff(value int) {
	if value%4 == 0 {
		select {}
	}
}
