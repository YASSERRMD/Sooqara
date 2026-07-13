package api

import (
	"fmt"
	"testing"
)

func TestFormatString(t *testing.T) {
	s := fmt.Sprintf("job-%d", 123)
	if s != "job-123" {
		t.Errorf("got %q, want job-123", s)
	}
}

func TestFormatMultiArg(t *testing.T) {
	s := fmt.Sprintf("%s: %d errors", "test", 5)
	if s != "test: 5 errors" {
		t.Errorf("got %q, want test: 5 errors", s)
	}
}
