package orchestrator

import (
	"testing"
	"time"
)

func TestStageAnalyseTimeout(t *testing.T) {
	s := &analyseStage{}
	if s.Timeout() != 120*time.Second {
		t.Errorf("analyse timeout = %v, want 120s", s.Timeout())
	}
}

func TestStageCopyTimeout(t *testing.T) {
	s := &copyStage{}
	if s.Timeout() != 120*time.Second {
		t.Errorf("copy timeout = %v, want 120s", s.Timeout())
	}
}

func TestStageImageTimeout(t *testing.T) {
	s := &imageStage{}
	if s.Timeout() != 180*time.Second {
		t.Errorf("image timeout = %v, want 180s", s.Timeout())
	}
}

func TestStageVideoTimeout(t *testing.T) {
	s := &videoStage{}
	if s.Timeout() != 60*time.Second {
		t.Errorf("video timeout = %v, want 60s", s.Timeout())
	}
}
