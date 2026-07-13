package limiter

// refillLocked adds tokens based on elapsed time. Caller must hold l.mu.
func (l *Limiter) refillLocked() {
	now := l.clock.Now()
	elapsed := now.Sub(l.lastRefill)
	if elapsed <= 0 {
		return
	}
	newTokens := float64(elapsed) / float64(l.rate)
	l.tokens += newTokens
	if l.tokens > float64(l.burst) {
		l.tokens = float64(l.burst)
	}
	l.lastRefill = now
}
