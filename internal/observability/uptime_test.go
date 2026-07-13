package observability

import (
	"testing"
	"time"
)

func TestMetricsUptimeIncreasing(t *testing.T) {
	m := NewMetrics()
	snap1 := m.Snapshot()
	time.Sleep(50 * time.Millisecond)
	snap2 := m.Snapshot()
	if snap2.Uptime <= snap1.Uptime {
		t.Error("expected uptime to increase over time")
	}
}

func TestMetricsAllCountersZeroInitially(t *testing.T) {
	m := NewMetrics()
	snap := m.Snapshot()
	counters := []struct {
		name  string
		value int64
	}{
		{"JobsCreated", snap.JobsCreated},
		{"JobsCompleted", snap.JobsCompleted},
		{"JobsFailed", snap.JobsFailed},
		{"StagesRun", snap.StagesRun},
		{"ErrorsCount", snap.ErrorsCount},
		{"TotalTokens", snap.TotalTokens},
	}
	for _, c := range counters {
		if c.value != 0 {
			t.Errorf("expected %s to be 0 initially, got %d", c.name, c.value)
		}
	}
}
