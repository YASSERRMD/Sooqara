package orchestrator

import (
	"testing"
	"time"
)

func TestConfigPollInterval(t *testing.T) {
	cfg := Config{PollInterval: 5 * time.Second}
	if cfg.PollInterval != 5*time.Second {
		t.Errorf("PollInterval = %v, want 5s", cfg.PollInterval)
	}
}

func TestGracePeriodDefault(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.GracePeriod != 60*time.Second {
		t.Errorf("GracePeriod = %v, want 60s", cfg.GracePeriod)
	}
}
