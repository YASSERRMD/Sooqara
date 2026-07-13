package store

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func newTestDB(t *testing.T) *Store {
	t.Helper()
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	blob := NewFilesystemBlob(filepath.Join(tmpDir, "storage"))
	return NewStore(db, blob)
}

func TestCreateAndGetJob(t *testing.T) {
	s := newTestDB(t)
	hint := "running shoes"
	job := NewJob("/path/to/photo.jpg", hint, "friendly", 4)

	if err := CreateJob(s.DB, job); err != nil {
		t.Fatalf("CreateJob failed: %v", err)
	}
	if job.ID == "" {
		t.Error("expected job ID to be set")
	}

	got, err := GetJob(s.DB, job.ID)
	if err != nil {
		t.Fatalf("GetJob failed: %v", err)
	}
	if got.State != StateQueued {
		t.Errorf("state = %s, want %s", got.State, StateQueued)
	}
	if got.Tone != "friendly" {
		t.Errorf("tone = %s, want friendly", got.Tone)
	}
}

func TestTransitionSuccess(t *testing.T) {
	s := newTestDB(t)
	hint := "test"
	job := NewJob("/path.jpg", hint, "", 2)
	CreateJob(s.DB, job)

	if err := Transition(s.DB, job.ID, StateQueued, StateAnalysing); err != nil {
		t.Fatalf("Transition failed: %v", err)
	}

	got, _ := GetJob(s.DB, job.ID)
	if got.State != StateAnalysing {
		t.Errorf("state = %s, want %s", got.State, StateAnalysing)
	}
}

func TestTransitionWrongStateFails(t *testing.T) {
	s := newTestDB(t)
	hint := "test"
	job := NewJob("/path.jpg", hint, "", 2)
	CreateJob(s.DB, job)

	// Try to transition from wrong state
	err := Transition(s.DB, job.ID, StateCopywriting, StateImaging)
	if err == nil {
		t.Fatal("expected error for invalid transition")
	}
}

func TestClaimNextReturnsOne(t *testing.T) {
	s := newTestDB(t)

	hint1, hint2 := "a", "b"
	job1 := NewJob("/path1.jpg", hint1, "", 2)
	job2 := NewJob("/path2.jpg", hint2, "", 2)
	CreateJob(s.DB, job1)
	CreateJob(s.DB, job2)

	claimed, err := ClaimNext(s.DB, []State{StateQueued})
	if err != nil {
		t.Fatalf("ClaimNext failed: %v", err)
	}
	if claimed == nil {
		t.Fatal("expected a claimed job")
	}
	if claimed.State != StateAnalysing {
		t.Errorf("claimed state = %s, want %s", claimed.State, StateAnalysing)
	}
}

func TestClaimNextRace(t *testing.T) {
	s := newTestDB(t)

	hint := "race"
	for i := 0; i < 5; i++ {
		j := NewJob("/path.jpg", hint, "", 2)
		CreateJob(s.DB, j)
	}

	claimed := make(chan *Job, 5)
	for i := 0; i < 5; i++ {
		go func() {
			j, _ := ClaimNext(s.DB, []State{StateAnalysing})
			claimed <- j
		}()
	}

	count := 0
	for i := 0; i < 5; i++ {
		j := <-claimed
		if j != nil {
			count++
		}
	}
	if count != 0 {
		t.Errorf("expected 0 claims (all already analysing), got %d", count)
	}
}

func TestCancelJob(t *testing.T) {
	s := newTestDB(t)
	hint := "cancel"
	job := NewJob("/path.jpg", hint, "", 2)
	CreateJob(s.DB, job)

	if err := CancelJob(s.DB, job.ID); err != nil {
		t.Fatalf("CancelJob failed: %v", err)
	}

	got, _ := GetJob(s.DB, job.ID)
	if got.State != StateCancelled {
		t.Errorf("state = %s, want %s", got.State, StateCancelled)
	}
}

func TestListJobs(t *testing.T) {
	s := newTestDB(t)

	for i := 0; i < 3; i++ {
		hint := "list"
		job := NewJob("/path.jpg", hint, "", 2)
		CreateJob(s.DB, job)
	}

	jobs, err := ListJobs(s.DB, nil, 10, 0)
	if err != nil {
		t.Fatalf("ListJobs failed: %v", err)
	}
	if len(jobs) != 3 {
		t.Errorf("got %d jobs, want 3", len(jobs))
	}
}

func TestFilesystemBlobPutAndGet(t *testing.T) {
	tmpDir := t.TempDir()
	blob := NewFilesystemBlob(tmpDir)

	data := []byte("hello world")
	path, err := blob.Put(context.Background(), "test.txt", bytes.NewReader(data))
	if err != nil {
		t.Fatalf("Put failed: %v", err)
	}
	if path == "" {
		t.Error("expected non-empty path")
	}

	rc, err := blob.Get(context.Background(), path)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	defer rc.Close()

	readData, err := io.ReadAll(rc)
	if err != nil {
		t.Fatalf("ReadAll failed: %v", err)
	}
	if string(readData) != "hello world" {
		t.Errorf("got %q, want hello world", string(readData))
	}
}

func TestFilesystemBlobDeduplication(t *testing.T) {
	tmpDir := t.TempDir()
	blob := NewFilesystemBlob(tmpDir)

	data := []byte("same content")
	path1, _ := blob.Put(context.Background(), "a.txt", bytes.NewReader(data))
	path2, _ := blob.Put(context.Background(), "b.txt", bytes.NewReader(data))

	if path1 != path2 {
		t.Errorf("expected same path for dedup, got %q and %q", path1, path2)
	}
}

func TestFilesystemBlobURL(t *testing.T) {
	blob := NewFilesystemBlob("/storage")
	url := blob.URL("abc123/test.jpg")
	if url != "/api/artifacts/test.jpg" {
		t.Errorf("URL = %s, want /api/artifacts/test.jpg", url)
	}
}
