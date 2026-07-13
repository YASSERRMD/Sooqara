package observability

import (
	"sync"
	"testing"
	"time"
)

func TestCostLedgerConcurrency(t *testing.T) {
	l, err := NewCostLedger()
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			entry := CostEntry{
				JobID:     "job-concurrent",
				Stage:     "vision",
				Model:     "agnes-v1",
				Tokens:    int64(id),
				CostUSD:   0.01,
				Timestamp: time.Now().UTC(),
			}
			if err := l.Append(entry); err != nil {
				t.Errorf("append failed: %v", err)
			}
		}(i)
	}
	wg.Wait()

	total, err := l.TotalCostUSD()
	if err != nil {
		t.Fatal(err)
	}
	expected := 50 * 0.01
	if total != expected {
		t.Errorf("expected %f, got %f", expected, total)
	}
}
