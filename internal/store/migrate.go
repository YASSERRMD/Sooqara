package store

import (
	"database/sql"
	"embed"
	"fmt"
	"sort"
	"time"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// migrate runs any pending migrations against the database.
func migrate(db *sql.DB) error {
	entries, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("read migrations: %w", err)
	}

	// Sort entries by name (which encodes version)
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for _, entry := range entries {
		version := parseVersion(entry.Name())
		if version == 0 {
			continue
		}

		// Check if already applied
		var applied int
		err := db.QueryRow("SELECT 1 FROM schema_migrations WHERE version = ?", version).Scan(&applied)
		if err == nil {
			continue // already applied
		}

		sqlBytes, err := migrationsFS.ReadFile("migrations/" + entry.Name())
		if err != nil {
			return fmt.Errorf("read migration %s: %w", entry.Name(), err)
		}

		if _, err := db.Exec(string(sqlBytes)); err != nil {
			return fmt.Errorf("apply migration %s: %w", entry.Name(), err)
		}

		if _, err := db.Exec(
			"INSERT INTO schema_migrations (version, applied_at) VALUES (?, ?)",
			version, time.Now().UnixMilli(),
		); err != nil {
			return fmt.Errorf("record migration %s: %w", entry.Name(), err)
		}
	}

	return nil
}

func parseVersion(filename string) int {
	var v int
	fmt.Sscanf(filename, "%d_", &v)
	return v
}
