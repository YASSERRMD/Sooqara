package journal

import (
	"strings"
	"testing"
)

func TestTruncateDetail(t *testing.T) {
	long := strings.Repeat("x", 200)
	shortened := truncateDetail(long)
	if len(shortened) != maxDetailLen {
		t.Errorf("truncateDetail(%d chars) = %d chars, want %d", len(long), len(shortened), maxDetailLen)
	}
}
