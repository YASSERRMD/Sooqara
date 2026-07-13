package observability

import "testing"

func TestNewMetricsStartTime(t *testing.T) {
	m := NewMetrics()
	snap := m.Snapshot()
	if snap.JobsCreated != 0 {
		t.Errorf("expected 0 jobs created, got %d", snap.JobsCreated)
	}
	if snap.Uptime <= 0 {
		t.Error("expected positive uptime")
	}
}

func TestRecordJobCreated(t *testing.T) {
	m := NewMetrics()
	m.RecordJobCreated()
	snap := m.Snapshot()
	if snap.JobsCreated != 1 {
		t.Errorf("expected 1, got %d", snap.JobsCreated)
	}
}

func TestRecordJobCompleted(t *testing.T) {
	m := NewMetrics()
	m.RecordJobCompleted()
	snap := m.Snapshot()
	if snap.JobsCompleted != 1 {
		t.Errorf("expected 1, got %d", snap.JobsCompleted)
	}
}

func TestRecordJobFailed(t *testing.T) {
	m := NewMetrics()
	m.RecordJobFailed()
	snap := m.Snapshot()
	if snap.JobsFailed != 1 {
		t.Errorf("expected 1, got %d", snap.JobsFailed)
	}
}

func TestRecordStageRun(t *testing.T) {
	m := NewMetrics()
	m.RecordStageRun()
	snap := m.Snapshot()
	if snap.StagesRun != 1 {
		t.Errorf("expected 1, got %d", snap.StagesRun)
	}
}

func TestRecordError(t *testing.T) {
	m := NewMetrics()
	m.RecordError()
	snap := m.Snapshot()
	if snap.ErrorsCount != 1 {
		t.Errorf("expected 1, got %d", snap.ErrorsCount)
	}
}

func TestAddTokens(t *testing.T) {
	m := NewMetrics()
	m.AddTokens(100)
	m.AddTokens(50)
	snap := m.Snapshot()
	if snap.TotalTokens != 150 {
		t.Errorf("expected 150 tokens, got %d", snap.TotalTokens)
	}
}

func TestAddCost(t *testing.T) {
	m := NewMetrics()
	m.AddCost(0.05)
	m.AddCost(0.10)
	snap := m.Snapshot()
	if snap.TotalCostUSD != 0.15 {
		t.Errorf("expected 0.15 USD, got %f", snap.TotalCostUSD)
	}
}
