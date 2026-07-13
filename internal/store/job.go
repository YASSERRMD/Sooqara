package store

import "time"

// Job represents an e-commerce listing job.
type Job struct {
	ID              string    `db:"id"`
	CreatedAt       int64     `db:"created_at"`
	UpdatedAt       int64     `db:"updated_at"`
	State           State     `db:"state"`
	SourceImagePath string    `db:"source_image_path"`
	ProductHint     *string   `db:"product_hint"`
	Tone            string    `db:"tone"`
	VariantCount    int       `db:"variant_count"`
	Seed            *int64    `db:"seed"`
	Warning         *string   `db:"warning"`
	Error           *string   `db:"error"`
}

// NewJob creates a new job in queued state.
func NewJob(sourceImagePath, productHint, tone string, variantCount int) *Job {
	now := time.Now().UnixMilli()
	hint := productHint
	if hint == "" {
		hint = "general"
	}
	if tone == "" {
		tone = "clear and practical"
	}
	return &Job{
		CreatedAt:       now,
		UpdatedAt:       now,
		State:           StateQueued,
		SourceImagePath: sourceImagePath,
		ProductHint:     &hint,
		Tone:            tone,
		VariantCount:    variantCount,
	}
}
