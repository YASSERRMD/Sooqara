package observability

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

// CostLedger stores cost entries in an SQLite database.
type CostLedger struct {
	db  *sql.DB
	mu  sync.Mutex
	idx int64
}

// NewCostLedger opens an in-memory SQLite database and creates the costs table.
func NewCostLedger() (*CostLedger, error) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return nil, fmt.Errorf("open ledger db: %w", err)
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS costs (
			id          TEXT PRIMARY KEY,
			job_id      TEXT NOT NULL,
			stage       TEXT NOT NULL,
			model       TEXT NOT NULL,
			tokens      INTEGER DEFAULT 0,
			cost_usd    REAL NOT NULL,
			timestamp   TEXT NOT NULL
		);
	`)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("create costs table: %w", err)
	}
	return &CostLedger{db: db}, nil
}

// Close flushes and closes the underlying database.
func (l *CostLedger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.db.Close()
}

// Append records a cost entry.
func (l *CostLedger) Append(entry CostEntry) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.idx++
	id := entry.ID
	if id == "" {
		id = fmt.Sprintf("cost-%d", l.idx)
	}
	_, err := l.db.Exec(
		"INSERT INTO costs (id, job_id, stage, model, tokens, cost_usd, timestamp) VALUES (?, ?, ?, ?, ?, ?, ?)",
		id, entry.JobID, entry.Stage, entry.Model, entry.Tokens, entry.CostUSD, entry.Timestamp.Format(time.RFC3339),
	)
	return err
}

// TotalCostUSD returns the sum of all recorded costs.
func (l *CostLedger) TotalCostUSD() (float64, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	var total sql.NullFloat64
	err := l.db.QueryRow("SELECT SUM(cost_usd) FROM costs").Scan(&total)
	if err != nil {
		return 0, err
	}
	if !total.Valid {
		return 0, nil
	}
	return total.Float64, nil
}

// EntriesByJob returns cost entries for a given job ID.
func (l *CostLedger) EntriesByJob(jobID string) ([]CostEntry, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	rows, err := l.db.Query("SELECT id, job_id, stage, model, tokens, cost_usd, timestamp FROM costs WHERE job_id = ?", jobID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var entries []CostEntry
	for rows.Next() {
		var e CostEntry
		var ts string
		if err := rows.Scan(&e.ID, &e.JobID, &e.Stage, &e.Model, &e.Tokens, &e.CostUSD, &ts); err != nil {
			return nil, err
		}
		e.Timestamp, _ = time.Parse(time.RFC3339, ts)
		entries = append(entries, e)
	}
	return entries, nil
}
