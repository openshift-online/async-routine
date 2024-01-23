package async

import (
	"context"
	"sync"
	"time"

	. "github.com/onsi/ginkgo/v2/dsl/core"
	. "github.com/onsi/gomega"
)

var _ = Describe("AsyncRoutine", func() {
	It("Run Async Routine", func() {

		routineRan := false
		var wg sync.WaitGroup
		wg.Add(1)
		routine := asyncRoutine{
			name:           "testRoutine",
			routine:        func() { routineRan = true; wg.Done() },
			createdAt:      time.Now().UTC(),
			status:         RoutineStatusCreated,
			ctx:            context.Background(),
			originatorOpId: "12345678",
		}

		routine.run()
		wg.Wait()
		Expect(routineRan).To(BeTrue())
		Expect(routine.StartedAt()).ToNot(BeNil())
		Expect(routine.FinishedAt()).ToNot(BeNil())
		Expect(routine.FinishedAt().After(*routine.StartedAt())).To(BeTrue())
		Expect(routine.Status()).To(Equal(RoutineStatusFinished))
	})
})
