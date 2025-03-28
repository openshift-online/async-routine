package metrics

import (
	"fmt"
	"sort"

	"github.com/openshift-online/async-routine"
	"github.com/prometheus/client_golang/prometheus"
)

var _ async.RoutinesObserver = (*metricObserver)(nil)

type metricObserver struct{}

var (
	runningRoutines = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "async_routine_manager_routines",
			Help: "Number of running routines.",
		},
	)

	runningRoutinesByName = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "async_routine_manager_routine_instances",
			Help: "Number of running instances of a given routine.",
		},
		[]string{"routine_name", "data"},
	)
)

func init() {
	prometheus.MustRegister(runningRoutines)
	prometheus.MustRegister(runningRoutinesByName)
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

func (m *metricObserver) RoutineStarted(routine async.AsyncRoutine) {
	runningRoutinesByName.
		With(prometheus.Labels{"routine_name": routine.Name(), "data": mapToString(routine.GetData())}).
		Inc()
}

func (m *metricObserver) RoutineFinished(routine async.AsyncRoutine) {
	runningRoutinesByName.
		With(prometheus.Labels{"routine_name": routine.Name(), "data": mapToString(routine.GetData())}).
		Dec()
}

func (m *metricObserver) RoutineExceededTimebox(routine async.AsyncRoutine) {
}

func (m *metricObserver) RunningRoutineCount(count int) {
	runningRoutines.Set(float64(count))
}

func (m *metricObserver) RunningRoutineByNameCount(name string, count int) {
}

func NewMetricObserver() async.RoutinesObserver {
	return &metricObserver{}
}
