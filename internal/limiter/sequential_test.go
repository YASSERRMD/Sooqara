package limiter

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestSequentialTokensReleaseOnePerSecond(t *testing.T) {
	clock := newFakeClock(time.Time{})
	l := New(60, 1, clock) // 1 token/sec, burst 1

	ctx := context.Background()

	if err := l.Acquire(ctx); err != nil {
		t.Fatalf("first Acquire failed: %v", err)
	}

	done1 := make(chan struct{})
	go func() {
		l.Acquire(ctx)
		close(done1)
	}()
	clock.advance(time.Second)
	<-done1

	done2 := make(chan struct{})
	go func() {
		l.Acquire(ctx)
		close(done2)
	}()
	clock.advance(time.Second)
	<-done2
}
