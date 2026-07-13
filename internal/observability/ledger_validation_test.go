package observability

import (
	"testing"
	"time"
)

func TestCostLedgerAppendWithValidation(t *testing.T) {
	l, err := NewCostLedger()
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	valid := CostEntry{
		ID:        "c-1",
		JobID:     "job-1",
		Stage:     "vision",
		Model:     "agnes-v1",
		Tokens:    100,
		CostUSD:   0.05,
		Timestamp: time.Now().UTC(),
	}
	if err := ValidateCostEntry(valid); err != nil {
		t.Fatalf("valid entry failed validation: %v", err)
	}

	if err := l.Append(valid); err != nil {
		t.Fatalf("failed to append valid entry: %v", err)
	}

	total, err := l.TotalCostUSD()
	if err != nil {
		t.Fatal(err)
	}
	if total != 0.05 {
		t.Errorf("expected 0.05, got %f", total)
	}
}
