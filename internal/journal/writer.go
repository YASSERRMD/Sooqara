package journal

import (
	"context"
)

func (j *Journal) writer() {
	defer j.wg.Done()
	for a := range j.ch {
		if j.db == nil {
			continue
		}
		if err := j.insert(a); err != nil {
			return
		}
	}
}

func (j *Journal) insert(a Activity) error {
	jobID := (*string)(nil)
	if a.JobID != nil {
		jobID = a.JobID
	}
	model := (*string)(nil)
	if a.Model != nil {
		model = a.Model
	}

	_, err := j.db.ExecContext(
		context.Background(),
		`INSERT INTO activities (id, ts, job_id, kind, model, request_hash, latency_ms, outcome, detail)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		a.ID, a.Ts, jobID, a.Kind, model, a.RequestHash, a.LatencyMs, a.Outcome, truncateDetail(a.Detail),
	)
	return err
}
