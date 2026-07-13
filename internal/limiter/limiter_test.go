package limiter

import (
	"context"
	"testing"
	"time"
)

func TestTokenRefill(t *testing.T) {
	clock := newFakeClock(time.Time{})
	l := New(60, 3, clock) // 1 token/sec, burst 3

	if l.Stats().AvailableTokens != 3 {
		t.Errorf("initial tokens = %f, want 3", l.Stats().AvailableTokens)
	}

	ctx := context.Background()
	if err := l.Acquire(ctx); err != nil {
		t.Fatalf("Acquire failed: %v", err)
	}
	if l.Stats().AvailableTokens != 2 {
		t.Errorf("after acquire tokens = %f, want 2", l.Stats().AvailableTokens)
	}

	clock.advance(time.Second)
	if l.Stats().AvailableTokens != 3 {
		t.Errorf("after 1s tokens = %f, want 3 (refilled)", l.Stats().AvailableTokens)
	}

	clock.advance(10 * time.Second)
	if l.Stats().AvailableTokens != 3 {
		t.Errorf("after 11s tokens = %f, want 3 (capped at burst)", l.Stats().AvailableTokens)
	}
}
