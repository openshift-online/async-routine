package async

import (
	"context"
	"sync"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"gitlab.cee.redhat.com/service/uhc-clusters-service/pkg/opid"
)

var _ RoutinesObserver = (*testAsyncObserver)(nil)

type testAsyncObserver struct {
	startedRoutines           []string
	finishedRoutines          []string
	timedoutRoutines          []string
	runningRoutineCount       int
	runningRoutineByNameCount map[string]int
	originatorCtx             context.Context
	wg                        sync.WaitGroup
}

func (t *testAsyncObserver) RoutineStarted(routine AsyncRoutine) {
	Expect(routine.OriginatorOpId()).To(Equal(opid.FromContext(t.originatorCtx)))
	Expect(routine.OpId()).ToNot(BeEmpty())
	t.startedRoutines = append(t.startedRoutines, routine.Name())
}

func (t *testAsyncObserver) RoutineFinished(routine AsyncRoutine) {
	t.finishedRoutines = append(t.finishedRoutines, routine.Name())
	t.wg.Done()
}

func (t *testAsyncObserver) RoutineExceededTimebox(routine AsyncRoutine) {
	t.timedoutRoutines = append(t.timedoutRoutines, routine.Name())
}

func (t *testAsyncObserver) RunningRoutineCount(count int) {
	t.runningRoutineCount = count
}

func (t *testAsyncObserver) RunningRoutineByNameCount(name string, count int) {
	t.runningRoutineByNameCount[name] = count
}

var _ = Describe("Async Routine Monitor", func() {
	It("Track async routine execution", func() {

		ctx := opid.NewContext()
		testObserver := testAsyncObserver{
			runningRoutineByNameCount: map[string]int{},
			originatorCtx:             ctx,
		}
		observerId := AddObserver(&testObserver)
		defer RemoveObserver(observerId)

		testObserver.wg.Add(3)
		NewAsyncRoutine("count up to 9", ctx, func() {
			for i := 0; i < 10; i++ {
			}
		}).Run()
		NewAsyncRoutine("count up to 9", ctx, func() {
			for i := 0; i < 10; i++ {
			}
		}).Run()
		NewAsyncRoutine("count up to 4", ctx, func() {
			for i := 0; i < 5; i++ {
			}
		}).Run()

		testObserver.wg.Wait()
		Expect(testObserver.startedRoutines).To(ConsistOf("count up to 9", "count up to 9", "count up to 4"))
		Expect(testObserver.finishedRoutines).To(ConsistOf("count up to 9", "count up to 9", "count up to 4"))
	})
})
