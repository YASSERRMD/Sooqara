package journal

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

// Init opens or creates the SQLite database and runs schema migrations.
func Init(dbPath string) (*sql.DB, *Journal, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, nil, fmt.Errorf("open database: %w", err)
	}

	for _, pragma := range []string{
		"PRAGMA journal_mode=WAL",
		"PRAGMA busy_timeout=5000",
		"PRAGMA foreign_keys=ON",
		"PRAGMA synchronous=NORMAL",
	} {
		if _, err := db.Exec(pragma); err != nil {
			db.Close()
			return nil, nil, fmt.Errorf("exec pragma: %w", err)
		}
	}

	if err := CreateSchema(db); err != nil {
		db.Close()
		return nil, nil, fmt.Errorf("create schema: %w", err)
	}

	journal := New(db, 256)
	return db, journal, nil
}
