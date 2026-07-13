package observability

import (
	"strings"
	"testing"
	"time"
)

func TestValidateCostEntryValid(t *testing.T) {
	e := CostEntry{
		ID:        "c-1",
		JobID:     "job-1",
		Stage:     "vision",
		Model:     "agnes-v1",
		Tokens:    100,
		CostUSD:   0.05,
		Timestamp: time.Now().UTC(),
	}
	if err := ValidateCostEntry(e); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestValidateCostEntryMissingJobID(t *testing.T) {
	e := CostEntry{
		JobID:     "",
		Stage:     "vision",
		Model:     "agnes-v1",
		CostUSD:   0.05,
		Timestamp: time.Now().UTC(),
	}
	err := ValidateCostEntry(e)
	if err == nil {
		t.Fatal("expected error for empty job_id")
	}
	if !errorsContains(err, "job_id is required") {
		t.Errorf("expected 'job_id is required' in error, got: %v", err)
	}
}

func TestValidateCostEntryNegativeCost(t *testing.T) {
	e := CostEntry{
		JobID:     "job-1",
		Stage:     "vision",
		Model:     "agnes-v1",
		CostUSD:   -1.0,
		Timestamp: time.Now().UTC(),
	}
	err := ValidateCostEntry(e)
	if err == nil {
		t.Fatal("expected error for negative cost")
	}
	if !errorsContains(err, "cost_usd must be non-negative") {
		t.Errorf("expected 'cost_usd must be non-negative' in error, got: %v", err)
	}
}

func TestValidateCostEntryMultipleErrors(t *testing.T) {
	e := CostEntry{
		JobID:     "",
		Stage:     "",
		Model:     "",
		CostUSD:   -1,
		Tokens:    -10,
		Timestamp: time.Time{},
	}
	err := ValidateCostEntry(e)
	if err == nil {
		t.Fatal("expected error for multiple invalid fields")
	}
}

func errorsContains(err error, s string) bool {
	return err != nil && strings.Contains(err.Error(), s)
}
