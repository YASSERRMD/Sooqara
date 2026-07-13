package orchestrator

import (
	"testing"
	"time"
)

func TestWorkerCountDefault(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.WorkerCount != 3 {
		t.Errorf("WorkerCount = %d, want 3", cfg.WorkerCount)
	}
}

func TestPollIntervalDefault(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.PollInterval != 2*time.Second {
		t.Errorf("PollInterval = %v, want 2s", cfg.PollInterval)
	}
}

func TestTimeoutValues(t *testing.T) {
	stages := []struct {
		stage Stage
		min   time.Duration
	}{
		{&analyseStage{}, 60 * time.Second},
		{&copyStage{}, 60 * time.Second},
		{&imageStage{}, 120 * time.Second},
		{&videoStage{}, 30 * time.Second},
		{&assembleStage{}, 30 * time.Second},
	}
	for _, tt := range stages {
		if tt.stage.Timeout() < tt.min {
			t.Errorf("%s timeout = %v, want >= %v", tt.stage.Name(), tt.stage.Timeout(), tt.min)
		}
	}
}
