package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestPrometheusMetrics(t *testing.T) {
	r := prometheus.NewPedanticRegistry()

	mo := NewMetricObserver(WithRegisterer(r))

	for _, tc := range []struct {
		name     string
		updateFn func()
		expCount int
	}{
		{
			name:     "init",
			updateFn: func() {},
			expCount: 1,
		},
		{
			name:     "set routine count to zero",
			updateFn: func() { mo.RunningRoutineCount(0) },
			expCount: 1,
		},
		{
			name: "set routine count to 1",
			updateFn: func() {
				mo.RunningRoutineCount(1)
			},
			expCount: 1,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tc.updateFn()

			n, err := testutil.GatherAndCount(r)
			if err != nil {
				t.Fatal(err)
			}
			if n != tc.expCount {
				t.Errorf("expected %d metrics, got %d", tc.expCount, n)
			}

			problems, err := testutil.GatherAndLint(r)
			if err != nil {
				t.Fatal(err)
			}

			for _, p := range problems {
				t.Errorf("found linting issue: %s: %s", p.Metric, p.Text)
			}
		})
	}
}
