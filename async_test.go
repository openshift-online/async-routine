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
		manager := newAsyncRoutineManager()
		defer manager.Monitor().Stop()

		executionLog := map[string][]string{}

		var mu sync.Mutex

		logData := func(key string, value string) {
			// I use a mutex here instead of a synced map since we need to sync writing in both the map and the
			// slice contained into the map
			mu.Lock()
			defer mu.Unlock()
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
			Do(func(r AsyncRoutine) {
				if r.OriginatorOpId() != opid.FromContext(ctx) {
					// we care only of our routines
					return
				}
				logData("RoutineStarted", r.Name())
			})
		observer.EXPECT().
			RoutineFinished(gomock.Any()).
			AnyTimes().
			Do(func(r AsyncRoutine) {
				if r.OriginatorOpId() != opid.FromContext(ctx) {
					// we care only of our routines
					return
				}
				logData("RoutineFinished", r.Name())
				wg.Done()
			})

		observerId := manager.AddObserver(observer)
		manager.Monitor().Start()
		defer manager.RemoveObserver(observerId)

		NewAsyncRoutine("count up to 9", ctx, func() {
			for i := 0; i < 10; i++ {
			}
		}).withRoutineManager(manager).
			Run()
		NewAsyncRoutine("count up to 9", ctx, func() {
			for i := 0; i < 10; i++ {
			}
		}).withRoutineManager(manager).
			Run()
		NewAsyncRoutine("count up to 4", ctx, func() {
			for i := 0; i < 5; i++ {
			}
		}).withRoutineManager(manager).
			Run()

		wg.Wait()
		Expect(executionLog["RoutineStarted"]).
			To(ConsistOf("count up to 9", "count up to 9", "count up to 4"))
		Expect(executionLog["RoutineFinished"]).
			To(ConsistOf("count up to 9", "count up to 9", "count up to 4"))
	})
	It("Snapshotting", func() {

		manager := newAsyncRoutineManager(WithSnapshottingInterval(time.Second))
		Expect(manager.Monitor().IsSnapshottingEnabled()).To(BeTrue())

		defer manager.Monitor().Stop()

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
			Do(func(r any) {
				if r.(AsyncRoutine).OriginatorOpId() != opid.FromContext(ctx) {
					// we care only of our routines
					return
				}
				methodCalled("RoutineStarted")
			})
		observer.EXPECT().
			RoutineFinished(gomock.Any()).
			AnyTimes().
			Do(func(r any) {
				if r.(AsyncRoutine).OriginatorOpId() != opid.FromContext(ctx) {
					// we care only of our routines
					return
				}
				methodCalled("RoutineFinished")
				wg.Done()
			})
		observer.EXPECT().
			RoutineExceededTimebox(gomock.Any()).
			AnyTimes().
			Do(func(r AsyncRoutine) {
				if r.OriginatorOpId() != opid.FromContext(ctx) {
					// we care only of our routines
					return
				}
				methodCalled("RoutineExceededTimebox")
			})
		observer.EXPECT().
			RunningRoutineCount(gomock.Any()).
			AnyTimes().
			Do(func(c int) { methodCalled("RunningRoutineCount") })
		observer.EXPECT().
			RunningRoutineByNameCount(gomock.Any(), gomock.Any()).
			AnyTimes().
			Do(func(name string, count int) { methodCalled("RunningRoutineByNameCount") })

		_ = manager.AddObserver(observer)
		manager.Monitor().Start()

		NewAsyncRoutine("count up to 4 - 1", ctx, func() {
			for i := 0; i < 5; i++ {
				time.Sleep(time.Second)
			}
		}).
			withRoutineManager(manager).
			Run()
		NewAsyncRoutine("count up to 4 - 2", ctx, func() {
			for i := 0; i < 5; i++ {
				time.Sleep(time.Second)
			}
		}).
			withRoutineManager(manager).
			Timebox(time.Second). // We want this routine to timeout
			Run()
		NewAsyncRoutine("count up to 4 - 3", ctx, func() {
			for i := 0; i < 5; i++ {
				time.Sleep(time.Second)
			}
		}).
			withRoutineManager(manager).
			Run()

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

		// We have one routine that will exceed the timebox. The `RoutineExceededTimebox` will notify that that routine
		// has exceeded the timebox at each iteration. Since the monitor will run every second, I would expect this to
		// happen 4 times, however due to load on the server, could happen that the monitor doesn't run exactly every
		// second so, to stay on the safe side, I check that the method is called at least once.
		Expect(callCount["RoutineExceededTimebox"]).ToNot(Equal(0))

		snapshot := manager.GetSnapshot()
		Expect(snapshot.GetTotalRoutineCount()).To(Equal(1))
		manager.Monitor().Stop()
		snapshot = manager.GetSnapshot()
		Expect(snapshot.GetTotalRoutineCount()).To(Equal(0))
	})
})
