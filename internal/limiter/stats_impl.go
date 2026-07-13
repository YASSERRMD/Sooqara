package limiter

// Stats returns the current limiter statistics.
func (l *Limiter) Stats() Stats {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.refillLocked()
	return Stats{
		AvailableTokens: l.tokens,
	}
}
