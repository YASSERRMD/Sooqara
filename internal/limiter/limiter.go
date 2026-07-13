package limiter

import (
	"sync"
	"time"
)

// Limiter is a token-bucket rate limiter enforcing a configurable RPM budget.
type Limiter struct {
	rate       time.Duration
	burst      int
	clock      Clock
	mu         sync.Mutex
	tokens     float64
	lastRefill time.Time
}

// New creates a new token-bucket limiter.
// rpm is requests per minute, burst is the maximum token bucket size,
// clock provides time abstraction.
func New(rpm int, burst int, clock Clock) *Limiter {
	if clock == nil {
		clock = realClock{}
	}
	rate := time.Minute / time.Duration(rpm)
	return &Limiter{
		rate:       rate,
		burst:      burst,
		clock:      clock,
		tokens:     float64(burst),
		lastRefill: clock.Now(),
	}
}
