package async

import (
	"context"
	. "github.com/onsi/ginkgo/v2"
	"sync"
	"time"

	"go.uber.org/mock/gomock"

	. "github.com/onsi/gomega"
)

var _ = Describe("AsyncRoutine", func() {
	DescribeTable("Run Async Routine", func(asyncObserver bool) {
		mockCtrl := gomock.NewController(GinkgoT())
		routineRan := false
		var wg sync.WaitGroup
		wg.Add(1)

		originatorOpId := "12345678"
		observer := NewMockRoutinesObserver(mockCtrl)

		manager := newAsyncRoutineManager()

		if asyncObserver {
			_ = manager.AddObserver(NewAsyncRoutineObserver(observer))
		} else {
			_ = manager.AddObserver(observer)
		}

		observer.EXPECT().RoutineStarted(gomock.Any()).AnyTimes()
		observer.EXPECT().RoutineFinished(gomock.Any()).AnyTimes().Do(
			func(r AsyncRoutine) {
				if r.OriginatorOpId() == originatorOpId {
					wg.Done()
				}
			})
		routine := asyncRoutine{
			name:           "testRoutine",
			routine:        func() { routineRan = true },
			createdAt:      time.Now().UTC(),
			status:         RoutineStatusCreated,
			ctx:            context.Background(),
			originatorOpId: originatorOpId,
		}
		routine.run(manager)
		wg.Wait()
		Expect(routineRan).To(BeTrue())
		Expect(routine.StartedAt()).ToNot(BeNil())
		Expect(routine.FinishedAt()).ToNot(BeNil())
		Expect(routine.FinishedAt().After(*routine.StartedAt())).To(BeTrue())
		Expect(routine.Status()).To(Equal(RoutineStatusFinished))
	},
		Entry("with sync observer", false),
		Entry("with async observer", true),
	)
})
