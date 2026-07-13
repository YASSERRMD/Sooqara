// Package observability provides metrics collection and cost tracking for Sooqara.
//
// Core types:
//   - Metrics: thread-safe cumulative counters for jobs, stages, tokens, and costs.
//   - CostLedger: SQLite-backed persistent store for per-call cost entries.
//   - PricingModel: configurable token and call pricing for AI models.
//
// Usage:
//
//	m := observability.NewMetrics()
//	m.RecordJobCreated()
//	m.RecordStageRun()
//	m.AddTokens(500)
//	m.AddCost(0.02)
//	snap := m.Snapshot()
//
//	l, err := observability.NewCostLedger()
//	l.Append(observability.CostEntry{JobID: "j1", Stage: "vision", CostUSD: 0.05})
//	total, _ := l.TotalCostUSD()
//
// All counters are safe for concurrent use. Snapshot returns an immutable copy.
package observability
