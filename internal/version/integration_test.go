package version

import (
	"strings"
	"testing"
)

func TestBuildInfoMultiline(t *testing.T) {
	info := BuildInfo()
	lines := strings.Split(info, "\n")
	if len(lines) < 3 {
		t.Errorf("expected at least 3 lines, got %d", len(lines))
	}
}

func TestVersionDefaults(t *testing.T) {
	if Version != "dev" {
		t.Errorf("expected default Version='dev', got %q", Version)
	}
	if Commit != "none" {
		t.Errorf("expected default Commit='none', got %q", Commit)
	}
	if BuildTime != "unknown" {
		t.Errorf("expected default BuildTime='unknown', got %q", BuildTime)
	}
}

func TestInfoEqualsVersion(t *testing.T) {
	if Info() != Version {
		t.Error("Info() should equal Version")
	}
}
