package limiter

import (
	"context"
	"testing"
	"time"
)

func TestStatsReportsCorrectTokens(t *testing.T) {
	clock := newFakeClock(time.Time{})
	l := New(60, 5, clock)

	stats := l.Stats()
	if stats.AvailableTokens != 5 {
		t.Errorf("initial AvailableTokens = %f, want 5", stats.AvailableTokens)
	}

	ctx := context.Background()
	l.Acquire(ctx)
	l.Acquire(ctx)
	l.Acquire(ctx)

	stats = l.Stats()
	if stats.AvailableTokens != 2 {
		t.Errorf("after 3 acquires AvailableTokens = %f, want 2", stats.AvailableTokens)
	}
}
