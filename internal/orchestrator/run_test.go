package orchestrator

import (
	"context"
	"testing"
	"time"

	"github.com/yasserrmd/sooqara/internal/store"
)

func TestNewWithZeroConfig(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := tmpDir + "/test.db"
	db, _ := store.Open(dbPath)
	defer db.Close()
	blob := store.NewFilesystemBlob(tmpDir + "/storage")
	s := store.NewStore(db, blob)

	o := New(Config{}, s, nil, nil)
	if o.cfg.WorkerCount != 3 {
		t.Errorf("default worker count = %d, want 3", o.cfg.WorkerCount)
	}
}

func TestRunStopsOnCancel(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := tmpDir + "/test.db"
	db, _ := store.Open(dbPath)
	defer db.Close()
	blob := store.NewFilesystemBlob(tmpDir + "/storage")
	s := store.NewStore(db, blob)

	cfg := Config{WorkerCount: 1, PollInterval: 100 * time.Millisecond}
	o := New(cfg, s, nil, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	done := make(chan error)
	go func() {
		done <- o.Run(ctx)
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Logf("Run returned: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Run did not return in time")
	}
}
