package journal

import (
	"context"
	"database/sql"
)

// Query returns journal entries matching the filter.
func (j *Journal) Query(ctx context.Context, f Filter) ([]Activity, error) {
	query := "SELECT id, ts, job_id, kind, model, request_hash, latency_ms, outcome, detail FROM activities WHERE 1=1"
	args := []interface{}{}

	if f.JobID != nil {
		query += " AND job_id = ?"
		args = append(args, *f.JobID)
	}
	if f.Kind != nil {
		query += " AND kind = ?"
		args = append(args, *f.Kind)
	}
	if f.From != nil {
		query += " AND ts >= ?"
		args = append(args, f.From.UnixMilli())
	}
	if f.To != nil {
		query += " AND ts <= ?"
		args = append(args, f.To.UnixMilli())
	}

	query += " ORDER BY ts DESC"

	if f.Limit > 0 {
		query += " LIMIT ?"
		args = append(args, f.Limit)
	}
	if f.Offset > 0 {
		query += " OFFSET ?"
		args = append(args, f.Offset)
	}

	rows, err := j.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Activity
	for rows.Next() {
		var a Activity
		var jobID sql.NullString
		var model sql.NullString
		if err := rows.Scan(&a.ID, &a.Ts, &jobID, &a.Kind, &model, &a.RequestHash, &a.LatencyMs, &a.Outcome, &a.Detail); err != nil {
			return nil, err
		}
		if jobID.Valid {
			a.JobID = &jobID.String
		}
		if model.Valid {
			a.Model = &model.String
		}
		results = append(results, a)
	}
	return results, rows.Err()
}
