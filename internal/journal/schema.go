package journal

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

// CreateSchema creates the journal tables if they don't exist.
func CreateSchema(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS activities(
			id            TEXT PRIMARY KEY,
			ts            INTEGER NOT NULL,
			job_id        TEXT,
			kind          TEXT NOT NULL,
			model         TEXT,
			request_hash  TEXT NOT NULL,
			latency_ms    INTEGER,
			outcome       TEXT NOT NULL,
			detail        TEXT
		);
		CREATE INDEX IF NOT EXISTS idx_act_job ON activities(job_id, ts);
		CREATE INDEX IF NOT EXISTS idx_act_kind ON activities(kind, ts);
	`)
	return err
}
