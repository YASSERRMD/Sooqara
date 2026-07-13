package pipeline

import (
	"testing"
)

func TestStyleV1HasNoEmDashes(t *testing.T) {
	if emDashRe.MatchString(StyleV1) {
		t.Error("STYLE_V1 should not contain em dashes")
	}
}

func TestMotionV1HasNoEmDashes(t *testing.T) {
	if emDashRe.MatchString(MotionV1) {
		t.Error("MOTION_V1 should not contain em dashes")
	}
}

func TestStyleV1Versioned(t *testing.T) {
	if StyleV1 == "" {
		t.Error("STYLE_V1 must not be empty")
	}
}

func TestMotionV1Versioned(t *testing.T) {
	if MotionV1 == "" {
		t.Error("MOTION_V1 must not be empty")
	}
}
