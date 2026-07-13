package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// Transition atomically transitions a job from one state to another.
// It fails if the job is not currently in the 'from' state.
func Transition(db *sql.DB, jobID string, from, to State) error {
	if !isValidTransition(from, to) {
		return fmt.Errorf("invalid transition: %s -> %s", from, to)
	}

	now := time.Now().UnixMilli()
	result, err := db.ExecContext(
		context.Background(),
		`UPDATE jobs SET state = ?, updated_at = ? WHERE id = ? AND state = ?`,
		to, now, jobID, from,
	)
	if err != nil {
		return fmt.Errorf("transition %s -> %s: %w", from, to, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("transition %s -> %s failed: job %s is in state %s, not %s",
			to, from, jobID, from, to)
	}

	return nil
}

func isValidTransition(from, to State) bool {
	valid, ok := ValidTransitions[from]
	if !ok {
		return false
	}
	for _, s := range valid {
		if s == to {
			return true
		}
	}
	return false
}
