package observability

import (
	"encoding/json"
	"testing"
)

func TestMetricsSnapshotJSON(t *testing.T) {
	m := NewMetrics()
	m.RecordJobCreated()
	m.RecordJobCompleted()
	m.AddTokens(100)
	m.AddCost(0.50)

	snap := m.Snapshot()
	data, err := json.Marshal(snap)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var decoded MetricsSnapshot
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if decoded.JobsCreated != 1 {
		t.Error("expected JobsCreated=1")
	}
	if decoded.JobsCompleted != 1 {
		t.Error("expected JobsCompleted=1")
	}
	if decoded.TotalTokens != 100 {
		t.Errorf("expected TotalTokens=100, got %d", decoded.TotalTokens)
	}
	if decoded.TotalCostUSD != 0.50 {
		t.Errorf("expected TotalCostUSD=0.50, got %f", decoded.TotalCostUSD)
	}
}

func TestMetricsSnapshotJSONKeys(t *testing.T) {
	m := NewMetrics()
	snap := m.Snapshot()
	data, err := json.Marshal(snap)
	if err != nil {
		t.Fatal(err)
	}
	expectedKeys := []string{
		"jobs_created", "jobs_completed", "jobs_failed",
		"stages_run", "errors_count", "total_tokens",
		"total_cost_usd", "uptime",
	}
	for _, key := range expectedKeys {
		if !containsString(string(data), `"`+key+`"`) {
			t.Errorf("expected JSON key %q not found in %s", key, string(data))
		}
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && indexOf(s, substr) >= 0
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
