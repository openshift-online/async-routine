package async

import (
	"context"

	. "github.com/onsi/ginkgo/v2/dsl/core"
	. "github.com/onsi/gomega"
)

var _ = Describe("AsyncRoutineManager", func() {

	routine1Name := "AsyncRoutineManagerTest:test1"
	routine2Name := "AsyncRoutineManagerTest:test2"

	testFunctionFactory := func(quitChannel chan bool) func() {
		return func() {
			<-quitChannel
		}
	}

	It("Get Async Routine Snapshot", func() {
		routine1QuitChannel := make(chan bool)
		routine2QuitChannel := make(chan bool)
		routine3QuitChannel := make(chan bool)

		defer func() {
			routine1QuitChannel <- true
			routine2QuitChannel <- true
			routine3QuitChannel <- true
		}()

		NewAsyncRoutine(routine1Name, context.Background(), testFunctionFactory(routine1QuitChannel)).Run()

		NewAsyncRoutine(routine2Name, context.Background(), testFunctionFactory(routine2QuitChannel)).Run()

		NewAsyncRoutine(routine2Name, context.Background(), testFunctionFactory(routine3QuitChannel)).Run()

		snapshot := Manager().GetSnapshot()
		Expect(snapshot.GetTotalRoutineCount()).To(BeNumerically(">=", 3))
		Expect(snapshot.GetRunningRoutinesCount(routine1Name)).To(Equal(1))
		Expect(snapshot.GetRunningRoutinesCount(routine2Name)).To(Equal(2))
	})
})
