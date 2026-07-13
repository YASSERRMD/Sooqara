// Package observability provides metrics collection and cost tracking for Sooqara.
package observability

import (
	"sync"
	"time"
)

// Metrics holds cumulative counters for the orchestrator lifecycle.
type Metrics struct {
	mu            sync.Mutex
	jobsCreated   int64
	jobsCompleted int64
	jobsFailed    int64
	stagesRun     int64
	errorsCount   int64
	totalTokens   int64
	totalCostUSD  float64
	startTime     time.Time
}

// NewMetrics creates a Metrics instance initialized with the current time.
func NewMetrics() *Metrics {
	return &Metrics{startTime: time.Now()}
}

// RecordJobCreated increments the jobs-created counter.
func (m *Metrics) RecordJobCreated() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.jobsCreated++
}

// RecordJobCompleted increments the jobs-completed counter.
func (m *Metrics) RecordJobCompleted() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.jobsCompleted++
}

// RecordJobFailed increments the jobs-failed counter.
func (m *Metrics) RecordJobFailed() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.jobsFailed++
}

// RecordStageRun increments the stages-run counter.
func (m *Metrics) RecordStageRun() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.stagesRun++
}

// RecordError increments the errors counter.
func (m *Metrics) RecordError() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.errorsCount++
}

// AddTokens adds tokens consumed by a provider call.
func (m *Metrics) AddTokens(n int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.totalTokens += n
}

// AddCost adds USD spent on a provider call.
func (m *Metrics) AddCost(usd float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.totalCostUSD += usd
}

// Snapshot returns a point-in-time copy of all counters.
func (m *Metrics) Snapshot() MetricsSnapshot {
	m.mu.Lock()
	defer m.mu.Unlock()
	return MetricsSnapshot{
		JobsCreated:   m.jobsCreated,
		JobsCompleted: m.jobsCompleted,
		JobsFailed:    m.jobsFailed,
		StagesRun:     m.stagesRun,
		ErrorsCount:   m.errorsCount,
		TotalTokens:   m.totalTokens,
		TotalCostUSD:  m.totalCostUSD,
		Uptime:        time.Since(m.startTime),
	}
}

// MetricsSnapshot is an immutable copy of Metrics counters.
type MetricsSnapshot struct {
	JobsCreated   int64         `json:"jobs_created"`
	JobsCompleted int64         `json:"jobs_completed"`
	JobsFailed    int64         `json:"jobs_failed"`
	StagesRun     int64         `json:"stages_run"`
	ErrorsCount   int64         `json:"errors_count"`
	TotalTokens   int64         `json:"total_tokens"`
	TotalCostUSD  float64       `json:"total_cost_usd"`
	Uptime        time.Duration `json:"uptime"`
}
