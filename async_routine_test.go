package async

import (
	"context"
	"sync"
	"time"

	"go.uber.org/mock/gomock"

	. "github.com/onsi/ginkgo/v2/dsl/core"
	. "github.com/onsi/gomega"
)

var _ = Describe("AsyncRoutine", func() {
	It("Run Async Routine", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		routineRan := false
		var wg sync.WaitGroup
		wg.Add(1)

		observer := NewMockRoutinesObserver(mockCtrl)

		observerId := Manager().AddObserver(observer)
		defer Manager().RemoveObserver(observerId)

		observer.EXPECT().RoutineStarted(gomock.Any()).AnyTimes()
		observer.EXPECT().RoutineFinished(gomock.Any()).AnyTimes().Do(func(r any) { wg.Done() })

		routine := asyncRoutine{
			name:           "testRoutine",
			routine:        func() { routineRan = true },
			createdAt:      time.Now().UTC(),
			status:         RoutineStatusCreated,
			ctx:            context.Background(),
			originatorOpId: "12345678",
		}
		Manager().run(&routine)
		wg.Wait()
		Expect(routineRan).To(BeTrue())
		Expect(routine.StartedAt()).ToNot(BeNil())
		Expect(routine.FinishedAt()).ToNot(BeNil())
		Expect(routine.FinishedAt().After(*routine.StartedAt())).To(BeTrue())
		Expect(routine.Status()).To(Equal(RoutineStatusFinished))
	})
})
