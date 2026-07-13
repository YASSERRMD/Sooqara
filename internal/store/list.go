package store

import (
	"context"
	"database/sql"
	"fmt"
)

// ListJobs returns paginated jobs filtered by state.
func ListJobs(db *sql.DB, state *State, limit, offset int) ([]*Job, error) {
	query := "SELECT id, created_at, updated_at, state, source_image_path, product_hint, tone, variant_count, seed, warning, error FROM jobs WHERE 1=1"
	args := []any{}

	if state != nil {
		query += " AND state = ?"
		args = append(args, string(*state))
	}

	query += " ORDER BY created_at DESC"
	if limit > 0 {
		query += " LIMIT ?"
		args = append(args, limit)
	}
	if offset > 0 {
		query += " OFFSET ?"
		args = append(args, offset)
	}

	rows, err := db.QueryContext(context.Background(), query, args...)
	if err != nil {
		return nil, fmt.Errorf("list jobs: %w", err)
	}
	defer rows.Close()

	var jobs []*Job
	for rows.Next() {
		var job Job
		var productHint, warning, errorStr sql.NullString
		var seed sql.NullInt64
		if err := rows.Scan(
			&job.ID, &job.CreatedAt, &job.UpdatedAt, &job.State,
			&job.SourceImagePath, &productHint, &job.Tone,
			&job.VariantCount, &seed, &warning, &errorStr,
		); err != nil {
			return nil, err
		}
		if productHint.Valid {
			job.ProductHint = &productHint.String
		}
		if warning.Valid {
			job.Warning = &warning.String
		}
		if errorStr.Valid {
			job.Error = &errorStr.String
		}
		if seed.Valid {
			s := seed.Int64
			job.Seed = &s
		}
		jobs = append(jobs, &job)
	}

	return jobs, rows.Err()
}
