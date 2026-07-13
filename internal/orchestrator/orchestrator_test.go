package orchestrator

import (
	"context"
	"testing"
	"time"

	"github.com/yasserrmd/sooqara/internal/store"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.WorkerCount != 3 {
		t.Errorf("WorkerCount = %d, want 3", cfg.WorkerCount)
	}
	if cfg.GracePeriod != 60*time.Second {
		t.Errorf("GracePeriod = %v, want 60s", cfg.GracePeriod)
	}
}

func TestStageTimeouts(t *testing.T) {
	stages := []Stage{
		&analyseStage{}, &copyStage{}, &imageStage{},
		&videoStage{}, &assembleStage{}, &noopStage{},
	}
	for _, s := range stages {
		if s.Timeout() <= 0 {
			t.Errorf("%s timeout = %v, want positive", s.Name(), s.Timeout())
		}
	}
}

func TestStageNames(t *testing.T) {
	names := make(map[string]bool)
	for _, name := range []string{"analyse", "copy", "image", "video", "assemble", "noop"} {
		names[name] = false
	}
	stages := []Stage{&analyseStage{}, &copyStage{}, &imageStage{}, &videoStage{}, &assembleStage{}, &noopStage{}}
	seen := make(map[string]bool)
	for _, s := range stages {
		seen[s.Name()] = true
	}
	expected := map[string]bool{
		"analyse": true, "copy": true, "image": true,
		"video": true, "assemble": true, "noop": true,
	}
	for name := range expected {
		if !seen[name] {
			t.Errorf("missing stage name: %s", name)
		}
	}
}

func TestStageForState(t *testing.T) {
	tests := []struct {
		state    store.State
		expected string
	}{
		{store.StateQueued, "analyse"},
		{store.StateAnalysing, "copy"},
		{store.StateCopywriting, "image"},
		{store.StateImaging, "video"},
		{store.StateAssembling, "assemble"},
		{store.StateDone, "noop"},
	}
	for _, tt := range tests {
		s := stageFor(tt.state)
		if s.Name() != tt.expected {
			t.Errorf("stageFor(%s) = %s, want %s", tt.state, s.Name(), tt.expected)
		}
	}
}

func TestStopWorkers(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := tmpDir + "/test.db"
	db, _ := store.Open(dbPath)
	defer db.Close()
	blob := store.NewFilesystemBlob(tmpDir + "/storage")
	s := store.NewStore(db, blob)

	cfg := DefaultConfig()
	o := New(cfg, s, nil, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Start a goroutine that stops after a brief delay
	go func() {
		time.Sleep(100 * time.Millisecond)
		o.stopWorkers(ctx)
	}()

	// Give it time to complete
	time.Sleep(200 * time.Millisecond)
}
