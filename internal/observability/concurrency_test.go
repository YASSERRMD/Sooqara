package observability

import "sync"

func TestMetricsConcurrencySafety(t *testing.T) {
	m := NewMetrics()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			m.RecordJobCreated()
			m.RecordJobCompleted()
			m.RecordJobFailed()
			m.RecordStageRun()
			m.RecordError()
			m.AddTokens(1)
			m.AddCost(0.001)
		}()
	}
	wg.Wait()
	snap := m.Snapshot()
	if snap.JobsCreated != 100 {
		t.Errorf("expected 100, got %d", snap.JobsCreated)
	}
	if snap.JobsCompleted != 100 {
		t.Errorf("expected 100, got %d", snap.JobsCompleted)
	}
	if snap.JobsFailed != 100 {
		t.Errorf("expected 100, got %d", snap.JobsFailed)
	}
	if snap.StagesRun != 100 {
		t.Errorf("expected 100, got %d", snap.StagesRun)
	}
	if snap.ErrorsCount != 100 {
		t.Errorf("expected 100, got %d", snap.ErrorsCount)
	}
	if snap.TotalTokens != 100 {
		t.Errorf("expected 100, got %d", snap.TotalTokens)
	}
}
