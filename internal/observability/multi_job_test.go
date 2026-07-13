package observability

import (
	"testing"
	"time"
)

func TestCostLedgerMultiJobIsolation(t *testing.T) {
	l, err := NewCostLedger()
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	now := time.Now().UTC()
	_ = l.Append(CostEntry{JobID: "j-a", Stage: "vision", Model: "m1", Tokens: 100, CostUSD: 0.10, Timestamp: now})
	_ = l.Append(CostEntry{JobID: "j-b", Stage: "copy", Model: "m2", Tokens: 200, CostUSD: 0.20, Timestamp: now})
	_ = l.Append(CostEntry{JobID: "j-a", Stage: "images", Model: "m3", Tokens: 50, CostUSD: 0.05, Timestamp: now})

	a, err := l.EntriesByJob("j-a")
	if err != nil {
		t.Fatal(err)
	}
	if len(a) != 2 {
		t.Fatalf("expected 2 entries for j-a, got %d", len(a))
	}

	b, err := l.EntriesByJob("j-b")
	if err != nil {
		t.Fatal(err)
	}
	if len(b) != 1 {
		t.Fatalf("expected 1 entry for j-b, got %d", len(b))
	}

	total, err := l.TotalCostUSD()
	if err != nil {
		t.Fatal(err)
	}
	if total != 0.35 {
		t.Errorf("expected total 0.35, got %f", total)
	}
}
