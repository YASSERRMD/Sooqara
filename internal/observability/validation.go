package observability

import (
	"errors"
	"strings"
	"time"
)

// ErrInvalidCostEntry is returned when a CostEntry fails validation.
var ErrInvalidCostEntry = errors.New("invalid cost entry")

// ValidateCostEntry checks that a CostEntry has required fields populated.
func ValidateCostEntry(e CostEntry) error {
	var errs []error
	if strings.TrimSpace(e.JobID) == "" {
		errs = append(errs, errors.New("job_id is required"))
	}
	if strings.TrimSpace(e.Stage) == "" {
		errs = append(errs, errors.New("stage is required"))
	}
	if strings.TrimSpace(e.Model) == "" {
		errs = append(errs, errors.New("model is required"))
	}
	if e.CostUSD < 0 {
		errs = append(errs, errors.New("cost_usd must be non-negative"))
	}
	if e.Tokens < 0 {
		errs = append(errs, errors.New("tokens must be non-negative"))
	}
	if e.Timestamp.IsZero() {
		errs = append(errs, errors.New("timestamp is required"))
	}
	if len(errs) > 0 {
		return errors.Join(append([]error{ErrInvalidCostEntry}, errs...)...)
	}
	return nil
}
