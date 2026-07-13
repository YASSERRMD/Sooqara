package limiter

import "context"

// Acquire blocks until a token is available or ctx is cancelled.
// On cancellation it returns ctx.Err() and consumes nothing.
func (l *Limiter) Acquire(ctx context.Context) error {
	for {
		l.mu.Lock()
		l.refillLocked()
		if l.tokens >= 1 {
			l.tokens--
			l.mu.Unlock()
			return nil
		}
		l.mu.Unlock()

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			l.clock.Sleep(l.rate)
		}
	}
}
