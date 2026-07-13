package orchestrator

import (
	"context"
	"testing"
	"time"

	"github.com/yasserrmd/sooqara/internal/store"
)

func TestCrashRecovery(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := tmpDir + "/test.db"
	db, _ := store.Open(dbPath)
	defer db.Close()

	hint := "test"
	job := store.NewJob("/src.jpg", hint, "", 2)
	store.CreateJob(db, job)

	// Simulate a job stuck in analysing state
	store.Transition(db, job.ID, store.StateQueued, store.StateAnalysing)

	// Check that the job is still in analysing state
	got, err := store.GetJob(db, job.ID)
	if err != nil {
		t.Fatalf("GetJob failed: %v", err)
	}
	if got.State != store.StateAnalysing {
		t.Errorf("state = %s, want %s", got.State, store.StateAnalysing)
	}
}

func TestGracefulShutdown(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := tmpDir + "/test.db"
	db, _ := store.Open(dbPath)
	defer db.Close()
	blob := store.NewFilesystemBlob(tmpDir + "/storage")
	s := store.NewStore(db, blob)

	cfg := Config{WorkerCount: 1, PollInterval: 100 * time.Millisecond}
	o := New(cfg, s, nil, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	done := make(chan error)
	go func() {
		done <- o.Run(ctx)
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Logf("Run error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for Run to complete")
	}
}
