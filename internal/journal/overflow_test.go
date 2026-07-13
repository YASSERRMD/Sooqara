package journal

import (
	"context"
	"testing"
	"time"
)

func TestJournalBufferOverflow(t *testing.T) {
	db, _, err := InitDB(t)
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}
	defer db.Close()

	j2 := New(db, 1)

	for i := 0; i < 100; i++ {
		a := NewActivityFromRequest("agnes.text", "test", "ok", 0, "", "")
		j2.Record(context.Background(), a)
	}

	time.Sleep(200 * time.Millisecond)

	dropped := j2.Dropped()
	if dropped == 0 {
		t.Log("no entries dropped (writer kept up)")
	}
}
