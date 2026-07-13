package observability

import (
	"testing"
	"time"
)

func TestNewCostLedger(t *testing.T) {
	l, err := NewCostLedger()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer l.Close()
}

func TestAppendAndTotal(t *testing.T) {
	l, err := NewCostLedger()
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	now := time.Now().UTC()
	err = l.Append(CostEntry{
		JobID:     "job-1",
		Stage:     "vision",
		Model:     "agnes-v1",
		Tokens:    100,
		CostUSD:   0.05,
		Timestamp: now,
	})
	if err != nil {
		t.Fatal(err)
	}

	total, err := l.TotalCostUSD()
	if err != nil {
		t.Fatal(err)
	}
	if total != 0.05 {
		t.Errorf("expected 0.05, got %f", total)
	}
}

func TestEntriesByJob(t *testing.T) {
	l, err := NewCostLedger()
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	now := time.Now().UTC()
	_ = l.Append(CostEntry{JobID: "job-1", Stage: "vision", Model: "m1", Tokens: 50, CostUSD: 0.02, Timestamp: now})
	_ = l.Append(CostEntry{JobID: "job-2", Stage: "copy", Model: "m2", Tokens: 100, CostUSD: 0.03, Timestamp: now})

	entries, err := l.EntriesByJob("job-1")
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Stage != "vision" {
		t.Errorf("expected stage vision, got %s", entries[0].Stage)
	}
}

func TestTotalCostEmpty(t *testing.T) {
	l, err := NewCostLedger()
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	total, err := l.TotalCostUSD()
	if err != nil {
		t.Fatal(err)
	}
	if total != 0 {
		t.Errorf("expected 0 for empty ledger, got %f", total)
	}
}
