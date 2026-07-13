package store

import (
	"database/sql"
)

// Store combines the database and blob storage into a single interface.
type Store struct {
	DB   *sql.DB
	Blob Blob
}

// NewStore creates a new Store from an open database and blob backend.
func NewStore(db *sql.DB, blob Blob) *Store {
	return &Store{DB: db, Blob: blob}
}

// Close shuts down the store.
func (s *Store) Close() error {
	return s.DB.Close()
}
