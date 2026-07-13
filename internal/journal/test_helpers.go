package journal

import (
	"database/sql"
	"path/filepath"
	"testing"
)

// InitDB creates a temp SQLite DB and returns *sql.DB and *Journal.
func InitDB(t *testing.T) (*sql.DB, *Journal, error) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, nil, err
	}

	for _, pragma := range []string{
		"PRAGMA journal_mode=WAL",
		"PRAGMA busy_timeout=5000",
	} {
		if _, err := db.Exec(pragma); err != nil {
			db.Close()
			return nil, nil, err
		}
	}

	if err := CreateSchema(db); err != nil {
		db.Close()
		return nil, nil, err
	}

	journal := New(db, 256)
	return db, journal, nil
}
