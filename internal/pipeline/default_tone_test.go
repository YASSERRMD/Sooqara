package pipeline

import (
	"testing"
)

func TestDefaultToneIsClearAndPractical(t *testing.T) {
	if DefaultTone != "clear and practical" {
		t.Errorf("DefaultTone = %q, want clear and practical", DefaultTone)
	}
}
