package journal

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

// NewActivityFromRequest creates an Activity from a request hash and metadata.
func NewActivityFromRequest(kind, request, outcome string, latency time.Duration, jobID, model string) Activity {
	ts := time.Now().UnixMilli()
	hash := sanitizeHash(request)
	detail := truncateDetail(request)

	a := Activity{
		ID:          randomID(),
		Ts:          ts,
		Kind:        kind,
		RequestHash: hash,
		LatencyMs:   int64(latency.Milliseconds()),
		Outcome:     outcome,
		Detail:      detail,
	}

	if jobID != "" {
		a.JobID = &jobID
	}
	if model != "" {
		a.Model = &model
	}

	return a
}

// randomID generates a short random hex ID.
func randomID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}
