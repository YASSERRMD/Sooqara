package limiter

import (
	"context"
	"testing"
	"time"
)

func TestAcquireReturnsContextError(t *testing.T) {
	clock := newFakeClock(time.Time{})
	l := New(60, 1, clock)

	ctx, cancel := context.WithCancel(context.Background())
	l.Acquire(ctx)
	cancel()

	err := l.Acquire(ctx)
	if err == nil {
		t.Fatal("expected cancellation error, got nil")
	}
}
