package metrics

import (
	"fmt"
	"sort"

	"github.com/prometheus/client_golang/prometheus"

	"gitlab.cee.redhat.com/service/uhc-clusters-service/pkg/async"
)

var _ async.RoutinesObserver = (*metricObserver)(nil)

type metricObserver struct{}

var (
	runningManagedRoutinesCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "async_routine_manager_routines_total",
			Help: "The total number of running manager routines.",
		},
	)

	runningManagedRoutinesByNameCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "async_routine_manager_routines_instances_count",
			Help: "The total number of running instance of a given routine.",
		},
		[]string{"routine_name", "data"},
	)
)

func init() {
	prometheus.MustRegister(runningManagedRoutinesCount)
	prometheus.MustRegister(runningManagedRoutinesByNameCount)
}

func mapToString(m map[string]string) string {
	if len(m) == 0 {
		return ""
	}

	// We sort all the keys to ensure the `data` field in the managed async routine metrics consistently
	// has keys in the same order.
	var keys = make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	dataString := ""
	for _, k := range keys {
		dataString = fmt.Sprintf("%s,%s=%s", dataString, k, m[k])
	}
	return dataString[1:]
}

func (m *metricObserver) routineData(routine async.AsyncRoutine) map[string]string {
	routineData := routine.GetData()

	// we always want to have the routine opid and its originator opid
	routineData["opid"] = routine.OpId()
	routineData["parent_opid"] = routine.OriginatorOpId()

	return routineData
}

func (m *metricObserver) RoutineStarted(routine async.AsyncRoutine) {
	runningManagedRoutinesByNameCount.
		With(prometheus.Labels{"routine_name": routine.Name(), "data": mapToString(m.routineData(routine))}).
		Inc()
}

func (m *metricObserver) RoutineFinished(routine async.AsyncRoutine) {
	runningManagedRoutinesByNameCount.
		With(prometheus.Labels{"routine_name": routine.Name(), "data": mapToString(m.routineData(routine))}).
		Dec()
}

func (m *metricObserver) RoutineExceededTimebox(routine async.AsyncRoutine) {
}

func (m *metricObserver) RunningRoutineCount(count int) {
	runningManagedRoutinesCount.Set(float64(count))
}

func (m *metricObserver) RunningRoutineByNameCount(name string, count int) {
}

func NewMetricObserver() async.RoutinesObserver {
	return &metricObserver{}
}
