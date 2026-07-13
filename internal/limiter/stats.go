package limiter

// Stats holds limiter statistics reported by Stats().
type Stats struct {
	AvailableTokens float64
	TotalAcquired   uint64
	CumulativeWait  int64
}
