package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// CreateJob inserts a new job into the database.
func CreateJob(db *sql.DB, job *Job) error {
	job.ID = newUUID()
	job.CreatedAt = time.Now().UnixMilli()
	job.UpdatedAt = job.CreatedAt

	productHint := ""
	if job.ProductHint != nil {
		productHint = *job.ProductHint
	}

	_, err := db.ExecContext(
		context.Background(),
		`INSERT INTO jobs (id, created_at, updated_at, state, source_image_path,
		 product_hint, tone, variant_count, seed, warning, error)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		job.ID, job.CreatedAt, job.UpdatedAt, job.State, job.SourceImagePath,
		productHint, job.Tone, job.VariantCount, job.Seed, job.Warning, job.Error,
	)
	return err
}

// GetJob retrieves a job by ID.
func GetJob(db *sql.DB, id string) (*Job, error) {
	var job Job
	var productHint, warning, errorStr sql.NullString
	var seed sql.NullInt64

	err := db.QueryRow(
		`SELECT id, created_at, updated_at, state, source_image_path,
		        product_hint, tone, variant_count, seed, warning, error
		 FROM jobs WHERE id = ?`, id,
	).Scan(
		&job.ID, &job.CreatedAt, &job.UpdatedAt, &job.State,
		&job.SourceImagePath, &productHint, &job.Tone,
		&job.VariantCount, &seed, &warning, &errorStr,
	)
	if err != nil {
		return nil, fmt.Errorf("get job %s: %w", id, err)
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
