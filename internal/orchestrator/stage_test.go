package orchestrator

import (
	"context"
	"testing"

	"github.com/yasserrmd/sooqara/internal/store"
)

func TestNoopStageReturnsNil(t *testing.T) {
	s := &noopStage{}
	err := s.Run(context.Background(), nil, nil, &store.Job{})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestAssembleStageReturnsNil(t *testing.T) {
	s := &assembleStage{}
	err := s.Run(context.Background(), nil, nil, &store.Job{})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestStageErrorHandling(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := tmpDir + "/test.db"
	db, _ := store.Open(dbPath)
	defer db.Close()
	blob := store.NewFilesystemBlob(tmpDir + "/storage")
	s := store.NewStore(db, blob)

	hint := "test"
	job := store.NewJob(tmpDir+"/src.jpg", hint, "", 2)
	store.CreateJob(s.DB, job)

	// Transition to failed state
	errStr := "stage failed"
	s.DB.Exec("UPDATE jobs SET error = ?, state = 'failed' WHERE id = ?", errStr, job.ID)

	var gotErr string
	s.DB.QueryRow("SELECT error FROM jobs WHERE id = ?", job.ID).Scan(&gotErr)
	if gotErr != errStr {
		t.Errorf("error = %q, want %q", gotErr, errStr)
	}
}
