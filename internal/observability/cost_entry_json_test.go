package observability

import (
	"encoding/json"
	"testing"
	"time"
)

func TestCostEntryJSONRoundTrip(t *testing.T) {
	original := CostEntry{
		ID:        "c-1",
		JobID:     "job-1",
		Stage:     "vision",
		Model:     "agnes-v1",
		Tokens:    1000,
		CostUSD:   0.05,
		Timestamp: time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC),
	}
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatal(err)
	}

	var decoded CostEntry
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}

	if decoded.ID != original.ID {
		t.Errorf("ID mismatch: got %q, want %q", decoded.ID, original.ID)
	}
	if decoded.JobID != original.JobID {
		t.Errorf("JobID mismatch: got %q, want %q", decoded.JobID, original.JobID)
	}
	if decoded.Stage != original.Stage {
		t.Errorf("Stage mismatch: got %q, want %q", decoded.Stage, original.Stage)
	}
	if decoded.Model != original.Model {
		t.Errorf("Model mismatch: got %q, want %q", decoded.Model, original.Model)
	}
	if decoded.Tokens != original.Tokens {
		t.Errorf("Tokens mismatch: got %d, want %d", decoded.Tokens, original.Tokens)
	}
	if decoded.CostUSD != original.CostUSD {
		t.Errorf("CostUSD mismatch: got %f, want %f", decoded.CostUSD, original.CostUSD)
	}
}

func TestCostEntryEmptyFields(t *testing.T) {
	e := CostEntry{}
	data, err := json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty JSON for empty CostEntry")
	}
}
