package observability

import "testing"

func TestCostLedgerCloseIdempotent(t *testing.T) {
	l, err := NewCostLedger()
	if err != nil {
		t.Fatal(err)
	}
	if err := l.Close(); err != nil {
		t.Fatalf("first close failed: %v", err)
	}
	// Second close should not panic
	if err := l.Close(); err != nil {
		t.Logf("second close returned error (acceptable): %v", err)
	}
}

func TestCostLedgerEntriesByJobNotFound(t *testing.T) {
	l, err := NewCostLedger()
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	entries, err := l.EntriesByJob("nonexistent-job")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}
