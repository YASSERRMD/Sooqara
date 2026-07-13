package orchestrator

import (
	"context"
	"testing"

	"github.com/yasserrmd/sooqara/internal/store"
)

func TestClaimNextReturnsNilWhenNoJobs(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := tmpDir + "/test.db"
	db, _ := store.Open(dbPath)
	defer db.Close()

	job, err := store.ClaimNext(db, []store.State{store.StateQueued})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if job != nil {
		t.Error("expected nil job when no jobs exist")
	}
}

func TestClaimNextWithEmptyStates(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := tmpDir + "/test.db"
	db, _ := store.Open(dbPath)
	defer db.Close()

	job, err := store.ClaimNext(db, []store.State{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if job != nil {
		t.Error("expected nil job with empty states")
	}
}

func TestListJobsWithStateFilter(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := tmpDir + "/test.db"
	db, _ := store.Open(dbPath)
	defer db.Close()

	hint := "test"
	job := store.NewJob("/src.jpg", hint, "", 2)
	store.CreateJob(db, job)

	// Filter by queued state
	state := store.StateQueued
	jobs, err := store.ListJobs(db, &state, 10, 0)
	if err != nil {
		t.Fatalf("ListJobs failed: %v", err)
	}
	if len(jobs) != 1 {
		t.Errorf("got %d jobs, want 1", len(jobs))
	}
}

func TestCancelJobReturnsErrorOnUnknownJob(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := tmpDir + "/test.db"
	db, _ := store.Open(dbPath)
	defer db.Close()

	err := store.CancelJob(db, "nonexistent-id")
	if err == nil {
		t.Fatal("expected error for unknown job")
	}
}

func TestTransitionSameStateFails(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := tmpDir + "/test.db"
	db, _ := store.Open(dbPath)
	defer db.Close()

	hint := "test"
	job := store.NewJob("/src.jpg", hint, "", 2)
	store.CreateJob(db, job)

	err := store.Transition(db, job.ID, store.StateQueued, store.StateQueued)
	if err == nil {
		t.Fatal("expected error transitioning to same state")
	}
}
