package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestPrometheusMetrics(t *testing.T) {
	problems, err := testutil.GatherAndLint(prometheus.DefaultGatherer)
	if err != nil {
		t.Fatal(err)
	}

	for _, p := range problems {
		t.Errorf("found linting issue: %s: %s", p.Metric, p.Text)
	}
}
