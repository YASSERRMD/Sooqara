package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// ClaimNext atomically claims one job in the given states.
// Uses a single UPDATE ... RETURNING to prevent races.
func ClaimNext(db *sql.DB, states []State) (*Job, error) {
	if len(states) == 0 {
		return nil, nil
	}

	placeholders := make([]string, len(states))
	args := make([]any, len(states))
	for i, s := range states {
		placeholders[i] = "?"
		args[i] = string(s)
	}

	query := fmt.Sprintf(
		`UPDATE jobs SET state = 'analysing' WHERE id IN (
			SELECT id FROM jobs WHERE state IN (%s) ORDER BY created_at ASC LIMIT 1
		) RETURNING id, created_at, updated_at, state, source_image_path,
		            product_hint, tone, variant_count, seed, warning, error`,
		strings.Join(placeholders, ","),
	)

	var job Job
	var productHint, warning, errorStr sql.NullString
	var seed sql.NullInt64

	err := db.QueryRow(query, args...).Scan(
		&job.ID, &job.CreatedAt, &job.UpdatedAt, &job.State,
		&job.SourceImagePath, &productHint, &job.Tone,
		&job.VariantCount, &seed, &warning, &errorStr,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("claim next: %w", err)
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

	return &job, nil
}
