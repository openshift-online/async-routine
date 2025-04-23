package metrics

import (
	"fmt"
	"sort"

	"github.com/openshift-online/async-routine"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var _ async.RoutinesObserver = (*metricObserver)(nil)

type metricObserver struct {
	registerer            prometheus.Registerer
	runningRoutines       prometheus.Gauge
	runningRoutinesByName *prometheus.GaugeVec
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
	m.runningRoutinesByName.
		With(prometheus.Labels{"routine_name": routine.Name(), "data": mapToString(routine.GetData())}).
		Inc()
}

func (m *metricObserver) RoutineFinished(routine async.AsyncRoutine) {
	m.runningRoutinesByName.
		With(prometheus.Labels{"routine_name": routine.Name(), "data": mapToString(routine.GetData())}).
		Dec()
}

func (m *metricObserver) RoutineExceededTimebox(routine async.AsyncRoutine) {
}

func (m *metricObserver) RunningRoutineCount(count int) {
	m.runningRoutines.Set(float64(count))
}

func (m *metricObserver) RunningRoutineByNameCount(name string, count int) {
}

// MetricOption defines options for the metrics observer.
type MetricOption func(*metricObserver)

// WithRegisterer configures the Prometheus registry where to register metrics.
func WithRegisterer(r prometheus.Registerer) MetricOption {
	return func(observer *metricObserver) {
		observer.registerer = r
	}
}

// NewMetricObserver returns an observer which tracks metrics for the async
// routines.
func NewMetricObserver(opts ...MetricOption) async.RoutinesObserver {
	observer := &metricObserver{
		registerer: prometheus.DefaultRegisterer,
	}

	for _, opt := range opts {
		opt(observer)
	}

	observer.runningRoutines = promauto.With(observer.registerer).NewGauge(
		prometheus.GaugeOpts{
			Name: "async_routine_manager_routines",
			Help: "Number of running routines.",
		},
	)
	observer.runningRoutinesByName = promauto.With(observer.registerer).NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "async_routine_manager_routine_instances",
			Help: "Number of running instances of a given routine.",
		},
		[]string{"routine_name", "data"},
	)

	return observer
}
