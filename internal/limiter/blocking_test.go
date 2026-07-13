package limiter

import (
	"context"
	"testing"
	"time"
)

func TestAcquireBlocksUntilTokenAvailable(t *testing.T) {
	clock := newFakeClock(time.Time{})
	l := New(60, 1, clock) // 1 token/sec, burst 1

	ctx := context.Background()
	if err := l.Acquire(ctx); err != nil {
		t.Fatalf("first Acquire failed: %v", err)
	}

	done := make(chan struct{})
	go func() {
		l.Acquire(ctx)
		close(done)
	}()

	select {
	case <-done:
		t.Fatal("Acquire returned immediately, expected block")
	case <-time.After(50 * time.Millisecond):
		// Expected: still blocked
	}

	clock.advance(time.Second)
	<-done
}
