// Package observability provides metrics collection and cost tracking for Sooqara.
package observability

import "time"

// CostEntry records a single AI-API cost event.
type CostEntry struct {
	ID        string    `json:"id"`
	JobID     string    `json:"job_id"`
	Stage     string    `json:"stage"`
	Model     string    `json:"model"`
	Tokens    int64     `json:"tokens,omitempty"`
	CostUSD   float64   `json:"cost_usd"`
	Timestamp time.Time `json:"timestamp"`
}
