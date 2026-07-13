package store

import "time"

// ArtifactKind represents the type of an artifact.
type ArtifactKind string

const (
	ArtifactAnalysis ArtifactKind = "analysis"
	ArtifactCopy     ArtifactKind = "copy"
	ArtifactImage    ArtifactKind = "image"
	ArtifactVideo    ArtifactKind = "video"
)

// Artifact represents a generated artifact for a job.
type Artifact struct {
	ID        string       `db:"id"`
	JobID     string       `db:"job_id"`
	Kind      ArtifactKind `db:"kind"`
	Seq       int          `db:"seq"`
	Path      *string      `db:"path"`
	Payload   *string      `db:"payload"`
	Seed      *int64       `db:"seed"`
	Prompt    *string      `db:"prompt"`
	StyleVer  *string      `db:"style_ver"`
	CreatedAt int64        `db:"created_at"`
}

// NewArtifact creates a new artifact record.
func NewArtifact(jobID string, kind ArtifactKind, seq int) *Artifact {
	return &Artifact{
		JobID:     jobID,
		Kind:      kind,
		Seq:       seq,
		CreatedAt: time.Now().UnixMilli(),
	}
}
