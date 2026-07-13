package journal

import (
	"context"
	"testing"
	"time"
)

func TestQueryFilter(t *testing.T) {
	db, journal, err := InitDB(t)
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}
	defer db.Close()

	a1 := NewActivityFromRequest("agnes.text", "chat request", "ok", 50*time.Millisecond, "job-1", "agnes-2.0-flash")
	a2 := NewActivityFromRequest("agnes.image", "image request", "ok", 200*time.Millisecond, "job-2", "agnes-image-2.1-flash")
	journal.Record(context.Background(), a1)
	journal.Record(context.Background(), a2)

	time.Sleep(100 * time.Millisecond)

	kind := "agnes.text"
	results, err := journal.Query(context.Background(), Filter{Kind: &kind, Limit: 10})
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("filtered query: expected 1 result, got %d", len(results))
	}
	if results[0].Kind != "agnes.text" {
		t.Errorf("got kind %s, want agnes.text", results[0].Kind)
	}
}
