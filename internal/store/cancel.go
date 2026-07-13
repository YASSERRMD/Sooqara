package store

import (
	"context"
	"database/sql"
	"time"
)

// CancelJob transitions a job to cancelled state.
func CancelJob(db *sql.DB, jobID string) error {
	now := time.Now().UnixMilli()
	result, err := db.ExecContext(
		context.Background(),
		`UPDATE jobs SET state = 'cancelled', updated_at = ? WHERE id = ? AND state != 'done' AND state != 'cancelled' AND state != 'failed'`,
		now, jobID,
	)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
