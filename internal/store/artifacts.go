package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// CreateArtifact inserts an artifact for a job.
func CreateArtifact(db *sql.DB, a *Artifact) error {
	a.ID = newUUID()
	a.CreatedAt = ctxTime()

	path, prompt, payload, styleVer := artifactStrings(a)

	_, err := db.ExecContext(
		context.Background(),
		`INSERT INTO artifacts (id, job_id, kind, seq, path, payload, seed, prompt, style_ver, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		a.ID, a.JobID, string(a.Kind), a.Seq, path, payload, a.Seed, prompt, styleVer, a.CreatedAt,
	)
	return err
}

// GetArtifactsByJob returns all artifacts for a job, ordered by kind then seq.
func GetArtifactsByJob(db *sql.DB, jobID string) ([]*Artifact, error) {
	rows, err := db.QueryContext(
		context.Background(),
		`SELECT id, job_id, kind, seq, path, payload, seed, prompt, style_ver, created_at
		 FROM artifacts WHERE job_id = ? ORDER BY kind, seq`, jobID,
	)
	if err != nil {
		return nil, fmt.Errorf("get artifacts: %w", err)
	}
	defer rows.Close()

	var artifacts []*Artifact
	for rows.Next() {
		a := &Artifact{}
		var path, prompt, payload, styleVer sql.NullString
		var seed sql.NullInt64
		if err := rows.Scan(&a.ID, &a.JobID, (*string)(&a.Kind), &a.Seq, &path, &payload, &seed, &prompt, &styleVer, &a.CreatedAt); err != nil {
			return nil, err
		}
		if path.Valid {
			a.Path = &path.String
		}
		if prompt.Valid {
			a.Prompt = &prompt.String
		}
		if payload.Valid {
			a.Payload = &payload.String
		}
		if styleVer.Valid {
			a.StyleVer = &styleVer.String
		}
		if seed.Valid {
			s := seed.Int64
			a.Seed = &s
		}
		artifacts = append(artifacts, a)
	}

	return artifacts, rows.Err()
}

func artifactStrings(a *Artifact) (path, prompt, payload, styleVer *string) {
	if a.Path != nil {
		v := *a.Path
		path = &v
	}
	if a.Prompt != nil {
		v := *a.Prompt
		prompt = &v
	}
	if a.Payload != nil {
		v := *a.Payload
		payload = &v
	}
	if a.StyleVer != nil {
		v := *a.StyleVer
		styleVer = &v
	}
	return path, prompt, payload, styleVer
}

func ctxTime() int64 {
	return time.Now().UnixMilli()
}
