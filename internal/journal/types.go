package journal

import "time"

// Activity represents a single journal entry.
type Activity struct {
	ID          string
	Ts          int64
	JobID       *string
	Kind        string
	Model       *string
	RequestHash string
	LatencyMs   int64
	Outcome     string
	Detail      string
}

// Filter allows querying journal entries.
type Filter struct {
	JobID  *string
	Kind   *string
	From   *time.Time
	To     *time.Time
	Limit  int
	Offset int
}
