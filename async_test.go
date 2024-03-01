package async

import (
	"sync"
	"time"

	"go.uber.org/mock/gomock"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"gitlab.cee.redhat.com/service/uhc-clusters-service/pkg/opid"
)

var _ = Describe("Async Routine Monitor", Ordered, func() {
	var mockCtrl *gomock.Controller
	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
	})
	It("Track async routine execution", func() {
		executionLog := map[string][]string{}

		logData := func(key string, value string) {
			executionLog[key] = append(executionLog[key], value)
		}

		var wg sync.WaitGroup
		wg.Add(3)
		ctx := opid.NewContext()

		observer := NewMockRoutinesObserver(mockCtrl)
		// Although the gomock matchers could be used to match params, utilizing them in asynchronous routines will lead
		// to panics in ginkgo when the expectation is not met, making it challenging to pinpoint the cause of failure.
		// To enhance clarity in error identification, I opt to manually verify them using gomega.
		observer.EXPECT().
			RoutineStarted(gomock.Any()).
			AnyTimes().
			Do(func(r AsyncRoutine) { logData("RoutineStarted", r.Name()) })
		observer.EXPECT().
			RoutineFinished(gomock.Any()).
			AnyTimes().
			Do(func(r AsyncRoutine) { logData("RoutineFinished", r.Name()); wg.Done() })

		manager := NewAsyncManagerBuilder().Build()

		observerId := manager.AddObserver(observer)
		defer manager.RemoveObserver(observerId)

		NewAsyncRoutine("count up to 9", ctx, func() {
			for i := 0; i < 10; i++ {
			}
		}).Run(manager)
		NewAsyncRoutine("count up to 9", ctx, func() {
			for i := 0; i < 10; i++ {
			}
		}).Run(manager)
		NewAsyncRoutine("count up to 4", ctx, func() {
			for i := 0; i < 5; i++ {
			}
		}).Run(manager)

		wg.Wait()
		Expect(executionLog["RoutineStarted"]).
			To(ConsistOf("count up to 9", "count up to 9", "count up to 4"))
		Expect(executionLog["RoutineFinished"]).
			To(ConsistOf("count up to 9", "count up to 9", "count up to 4"))
	})
	It("Snapshotting", func() {

		var wg sync.WaitGroup
		wg.Add(3)

		callCount := map[string]int{}
		lock := sync.Mutex{}
		methodCalled := func(name string) {
			lock.Lock()
			defer lock.Unlock()
			callCount[name] = callCount[name] + 1
		}

		ctx := opid.NewContext()
		observer := NewMockRoutinesObserver(mockCtrl)

		// Although the gomock Times method could be used for counting, utilizing it in asynchronous routines will lead
		// to panics in ginkgo when the expectation is not met, making it challenging to pinpoint the cause of failure.
		// To enhance clarity in error identification, I opt to manually count the calls and later verify them
		// using gomega.
		observer.EXPECT().
			RoutineStarted(gomock.Any()).
			AnyTimes().
			Do(func(r any) { methodCalled("RoutineStarted") })
		observer.EXPECT().
			RoutineFinished(gomock.Any()).
			AnyTimes().
			Do(func(r any) { methodCalled("RoutineFinished"); wg.Done() })
		observer.EXPECT().
			RunningRoutineCount(gomock.Any()).
			AnyTimes().
			Do(func(c int) { methodCalled("RunningRoutineCount") })
		observer.EXPECT().
			RunningRoutineByNameCount(gomock.Any(), gomock.Any()).
			AnyTimes().
			Do(func(name string, count int) { methodCalled("RunningRoutineByNameCount") })

		manager := NewAsyncManagerBuilder().
			WithSnapshottingInterval(time.Second).
			Build()
		observerId := manager.AddObserver(observer)
		defer manager.RemoveObserver(observerId)

		Expect(manager.monitor().IsSnapshottingEnabled()).To(BeTrue())

		r1 := NewAsyncRoutine("count up to 4 - 1", ctx, func() {
			for i := 0; i < 5; i++ {
				time.Sleep(time.Second)
			}
		}).Build()
		r2 := NewAsyncRoutine("count up to 4 - 2", ctx, func() {
			for i := 0; i < 5; i++ {
				time.Sleep(time.Second)
			}
		}).Build()
		r3 := NewAsyncRoutine("count up to 4 - 3", ctx, func() {
			for i := 0; i < 5; i++ {
				time.Sleep(time.Second)
			}
		}).Build()

		manager.Run(r1, r2, r3)

		wg.Wait()

		Expect(callCount["RoutineStarted"]).To(Equal(3))
		Expect(callCount["RoutineFinished"]).To(Equal(3))

		// The test is set up with the monitor configured to run every second, while the routine sleeps for 5 seconds.
		// Ideally, I anticipate the `RunningRoutineCount` to be invoked approximately 4 or 5 times.
		// However, to stay on the safe side, I've set the validation to ensure that `RunningRoutineCount`
		// is called at least 3 times or more.
		Expect(callCount["RunningRoutineCount"]).To(BeNumerically(">", 2))

		// The `RunningRoutineByNameCount` is invoked for each running routine whenever the monitor is triggered.
		// Given the presence of 3 defined routines in the test, plus the one managing the monitor itself,
		// it is anticipated that `RunningRoutineByNameCount` should be called more than 8 times
		// if the monitor is triggered more than 2 times.
		Expect(callCount["RunningRoutineByNameCount"]).To(BeNumerically(">", 8))
	})
})
