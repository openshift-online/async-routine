package metrics

import (
	"bytes"

	. "github.com/onsi/ginkgo/v2/dsl/core"
	. "github.com/onsi/ginkgo/v2/dsl/table"
	. "github.com/onsi/gomega"
	"github.com/openshift-online/async-routine"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Metrics Observer", func() {
	DescribeTable("Updates metrics",
		func(updateFn func(async.RoutinesObserver), expected string) {
			registry := prometheus.NewPedanticRegistry()
			observer := NewMetricObserver(WithRegisterer(registry))

			updateFn(observer)

			err := testutil.GatherAndCompare(registry, bytes.NewBufferString(expected))
			Expect(err).NotTo(HaveOccurred())

			problems, err := testutil.GatherAndLint(registry)
			Expect(err).NotTo(HaveOccurred())
			Expect(problems).To(BeEmpty())
		},
		Entry("from initial state",
			func(async.RoutinesObserver) {},
			`# HELP async_routine_manager_routines Number of running routines.
# TYPE async_routine_manager_routines gauge
async_routine_manager_routines 0
`,
		),
		Entry("when 1 routine is started",
			func(observer async.RoutinesObserver) {
				ctrl := gomock.NewController(GinkgoT())

				routine := async.NewMockAsyncRoutine(ctrl)
				routine.EXPECT().
					Name().
					Return("test").
					Times(1)
				routine.EXPECT().
					GetData().
					Return(nil).
					Times(1)

				observer.RoutineStarted(routine)
				observer.RunningRoutineCount(1)
			},
			`# HELP async_routine_manager_routine_instances Number of running instances of a given routine.
# TYPE async_routine_manager_routine_instances gauge
async_routine_manager_routine_instances{data="",routine_name="test"} 1
# HELP async_routine_manager_routines Number of running routines.
# TYPE async_routine_manager_routines gauge
async_routine_manager_routines 1
`),
		Entry("when 2 routines (1 with data) are started",
			func(observer async.RoutinesObserver) {
				ctrl := gomock.NewController(GinkgoT())

				routine := async.NewMockAsyncRoutine(ctrl)
				routine.EXPECT().
					Name().
					Return("test").
					Times(1)
				routine.EXPECT().
					GetData().
					Return(nil).
					Times(1)
				observer.RoutineStarted(routine)

				routine2 := async.NewMockAsyncRoutine(ctrl)
				routine2.EXPECT().
					Name().
					Return("test2")
				routine2.EXPECT().
					GetData().
					Return(map[string]string{"foo": "bar"}).
					Times(1)

				observer.RoutineStarted(routine2)
				observer.RunningRoutineCount(2)
			},
			`# HELP async_routine_manager_routine_instances Number of running instances of a given routine.
# TYPE async_routine_manager_routine_instances gauge
async_routine_manager_routine_instances{data="",routine_name="test"} 1
async_routine_manager_routine_instances{data="foo=bar",routine_name="test2"} 1
# HELP async_routine_manager_routines Number of running routines.
# TYPE async_routine_manager_routines gauge
async_routine_manager_routines 2
`,
		),
		Entry("when 1 routine is started then stopped",
			func(observer async.RoutinesObserver) {
				ctrl := gomock.NewController(GinkgoT())

				routine := async.NewMockAsyncRoutine(ctrl)
				routine.EXPECT().
					Name().
					Return("test").
					Times(2)
				routine.EXPECT().
					GetData().
					Return(nil).
					Times(2)

				observer.RoutineStarted(routine)
				observer.RoutineFinished(routine)
				observer.RunningRoutineCount(0)
			},
			`# HELP async_routine_manager_routine_instances Number of running instances of a given routine.
# TYPE async_routine_manager_routine_instances gauge
async_routine_manager_routine_instances{data="",routine_name="test"} 0
# HELP async_routine_manager_routines Number of running routines.
# TYPE async_routine_manager_routines gauge
async_routine_manager_routines 0
`,
		),
	)
})
