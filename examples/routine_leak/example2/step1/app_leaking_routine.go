package main

import (
	"math/rand"
	"os"
	"runtime/pprof"
	"time"
)

func main() {
	go doJob()

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
	go parentGoroutine()
}

func parentGoroutine() {
	go doInterestingStuff(rand.Intn(100))
	time.Sleep(500 * time.Millisecond)
}

func doInterestingStuff(value int) {
	if value%4 == 0 {
		select {}
	}
}
