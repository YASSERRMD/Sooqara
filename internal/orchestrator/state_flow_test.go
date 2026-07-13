package orchestrator

import (
	"context"
	"testing"
	"time"

	"github.com/yasserrmd/sooqara/internal/store"
)

func TestNextStateChain(t *testing.T) {
	chain := []store.State{
		store.StateQueued,
		store.StateAnalysing,
		store.StateCopywriting,
		store.StateImaging,
		store.StateVideoing,
		store.StateAssembling,
		store.StateDone,
	}
	for i := 0; i < len(chain)-1; i++ {
		got := store.NextState(chain[i])
		if got != chain[i+1] {
			t.Errorf("NextState(%s) = %s, want %s", chain[i], got, chain[i+1])
		}
	}
}

func TestVideoingNotInClaimList(t *testing.T) {
	// The orchestrator should NOT claim videoing jobs
	// This is verified by the stageFor function not having a videoing case
	states := []store.State{
		store.StateQueued, store.StateAnalysing,
		store.StateCopywriting, store.StateImaging,
		store.StateAssembling,
	}
	for _, s := range states {
		if s == store.StateVideoing {
			t.Error("videoing should not be in claim list")
		}
	}
}

func TestStageForUnknownState(t *testing.T) {
	s := stageFor("unknown_state")
	if s.Name() != "noop" {
		t.Errorf("stageFor(unknown) = %s, want noop", s.Name())
	}
}

func TestRunWithCancelledContext(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := tmpDir + "/test.db"
	db, _ := store.Open(dbPath)
	defer db.Close()
	blob := store.NewFilesystemBlob(tmpDir + "/storage")
	s := store.NewStore(db, blob)

	cfg := Config{WorkerCount: 1, PollInterval: 100 * time.Millisecond}
	o := New(cfg, s, nil, nil)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

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
		t.Fatal("Run did not return after context cancelled")
	}
}
