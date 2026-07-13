package journal

import (
	"context"
	"testing"
	"time"
)

func TestJournalBasics(t *testing.T) {
	db, journal, err := InitDB(t)
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}
	defer db.Close()

	req := "POST /v1/chat/completions model=agnes-2.0-flash"
	a := NewActivityFromRequest("agnes.text", req, "ok", 150*time.Millisecond, "job-1", "agnes-2.0-flash")

	if err := journal.Record(context.Background(), a); err != nil {
		t.Fatalf("Record failed: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	kind := "agnes.text"
	results, err := journal.Query(context.Background(), Filter{Kind: &kind, Limit: 10})
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Kind != "agnes.text" {
		t.Errorf("kind = %s, want agnes.text", results[0].Kind)
	}
	if results[0].Outcome != "ok" {
		t.Errorf("outcome = %s, want ok", results[0].Outcome)
	}
	if results[0].LatencyMs != 150 {
		t.Errorf("latency_ms = %d, want 150", results[0].LatencyMs)
	}
}
