package pipeline

import (
	"testing"
	"time"
)

func TestDefaultVideoConfigValues(t *testing.T) {
	cfg := DefaultVideoConfig()
	if cfg.Height != 768 {
		t.Errorf("Height = %d, want 768", cfg.Height)
	}
	if cfg.Width != 1152 {
		t.Errorf("Width = %d, want 1152", cfg.Width)
	}
	if cfg.NumFrames != 121 {
		t.Errorf("NumFrames = %d, want 121", cfg.NumFrames)
	}
	if cfg.FrameRate != 24 {
		t.Errorf("FrameRate = %d, want 24", cfg.FrameRate)
	}
	if cfg.Deadline != 15*time.Minute {
		t.Errorf("Deadline = %v, want 15m", cfg.Deadline)
	}
}
