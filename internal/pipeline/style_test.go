package pipeline

import (
	"strings"
	"testing"
)

func TestStyleV1IsNotEmpty(t *testing.T) {
	if StyleV1 == "" {
		t.Error("STYLE_V1 should not be empty")
	}
}

func TestMotionV1IsNotEmpty(t *testing.T) {
	if MotionV1 == "" {
		t.Error("MOTION_V1 should not be empty")
	}
}

func TestStyleV1ContainsKeyTerms(t *testing.T) {
	keyTerms := []string{"natural daylight", "shallow depth of field", "photographic realism", "no text", "no watermark"}
	for _, kt := range keyTerms {
		if !strings.Contains(StyleV1, kt) {
			t.Errorf("STYLE_V1 missing term: %s", kt)
		}
	}
}

func TestMotionV1ContainsKeyTerms(t *testing.T) {
	keyTerms := []string{"orbit", "smooth camera motion", "loops seamlessly"}
	for _, kt := range keyTerms {
		if !strings.Contains(MotionV1, kt) {
			t.Errorf("MOTION_V1 missing term: %s", kt)
		}
	}
}
