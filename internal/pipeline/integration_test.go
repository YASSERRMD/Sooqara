package pipeline

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/yasserrmd/sooqara/internal/provider"
	"github.com/yasserrmd/sooqara/internal/store"
)

func TestEmptyFakeProviderDoesNotPanic(t *testing.T) {
	fake := &FakeProvider{}
	resp, err := fake.Chat(context.Background(), provider.ChatRequest{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Choices) != 0 {
		t.Error("expected empty choices")
	}
}

func TestFakeProviderName(t *testing.T) {
	fake := &FakeProvider{}
	if fake.Name() != "fake" {
		t.Errorf("Name() = %s, want fake", fake.Name())
	}
}

func TestValidJSONProducesCorrectResponse(t *testing.T) {
	a := ProductAnalysis{ProductName: "X", Category: "Y", ShapeDescription: "Z"}
	resp := ValidJSON(a)
	if len(resp.Choices) != 1 {
		t.Fatal("expected 1 choice")
	}
	if resp.Choices[0].Message.Role != "assistant" {
		t.Error("expected assistant role")
	}
}

func TestMalformedJSONProducesBrokenResponse(t *testing.T) {
	resp := MalformedJSON()
	if len(resp.Choices) != 1 {
		t.Fatal("expected 1 choice")
	}
	content := resp.Choices[0].Message.Content
	if content != "{broken json!!!" {
		t.Errorf("got %q, want {broken json!!!", content)
	}
}

func TestTooFewSettingsReducesLifestyleSettings(t *testing.T) {
	a := ProductAnalysis{
		ProductName:        "W",
		Category:           "C",
		ShapeDescription:   "D",
		SuggestedLifestyleSettings: []string{"one", "two", "three", "four"},
	}
	TooFewSettings(a)
	if len(a.SuggestedLifestyleSettings) >= 4 {
		t.Error("expected settings reduced to 1")
	}
}

func TestNewJobDefaults(t *testing.T) {
	job := store.NewJob("/img.jpg", "", "", 2)
	if job.Tone != "clear and practical" {
		t.Errorf("tone = %s, want clear and practical", job.Tone)
	}
	if job.VariantCount != 2 {
		t.Errorf("variant_count = %d, want 2", job.VariantCount)
	}
	if job.State != store.StateQueued {
		t.Errorf("state = %s, want %s", job.State, store.StateQueued)
	}
}

func TestNextStateChain(t *testing.T) {
	states := []store.State{
		store.StateQueued,
		store.StateAnalysing,
		store.StateCopywriting,
		store.StateImaging,
		store.StateVideoing,
		store.StateAssembling,
		store.StateDone,
	}
	for i := 0; i < len(states)-1; i++ {
		got := store.NextState(states[i])
		if got != states[i+1] {
			t.Errorf("NextState(%s) = %s, want %s", states[i], got, states[i+1])
		}
	}
}

func TestValidTransitions(t *testing.T) {
	// Verify cancelled is reachable from queued
	transitions := store.ValidTransitions[store.StateQueued]
	found := false
	for _, s := range transitions {
		if s == store.StateCancelled {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected cancelled in queued transitions")
	}
}

func TestInvalidTransitions(t *testing.T) {
	transitions := store.ValidTransitions[store.StateDone]
	for _, s := range transitions {
		if s == store.StateAnalysing {
			t.Error("done -> analysing should not be valid")
		}
	}
}

func TestStoreOpenAndClose(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	db, err := store.Open(dbPath)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	if err := db.Close(); err != nil {
		t.Fatalf("Close failed: %v", err)
	}
}

func TestFilesystemBlobGetNonExistent(t *testing.T) {
	tmpDir := t.TempDir()
	blob := store.NewFilesystemBlob(tmpDir)
	_, err := blob.Get(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestFilesystemBlobURL(t *testing.T) {
	blob := store.NewFilesystemBlob("/storage")
	url := blob.URL("/storage/abc123/test.jpg")
	if url != "/api/artifacts/test.jpg" {
		t.Errorf("URL = %s, want /api/artifacts/test.jpg", url)
	}
}

func TestListJobsEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	db, _ := store.Open(dbPath)
	defer db.Close()

	jobs, err := store.ListJobs(db, nil, 10, 0)
	if err != nil {
		t.Fatalf("ListJobs failed: %v", err)
	}
	if len(jobs) != 0 {
		t.Errorf("expected 0 jobs, got %d", len(jobs))
	}
}

func TestCancelTerminalJobFails(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	db, _ := store.Open(dbPath)
	defer db.Close()
	blob := store.NewFilesystemBlob(filepath.Join(tmpDir, "storage"))
	s := store.NewStore(db, blob)

	hint := "test"
	job := store.NewJob("/source.jpg", hint, "", 2)
	store.CreateJob(s.DB, job)

	store.Transition(s.DB, job.ID, store.StateQueued, store.StateAnalysing)
	store.Transition(s.DB, job.ID, store.StateAnalysing, store.StateDone)

	err := store.CancelJob(s.DB, job.ID)
	if err == nil {
		t.Fatal("expected error cancelling done job")
	}
}

func TestSourceImageReadError(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	db, _ := store.Open(dbPath)
	defer db.Close()
	blob := store.NewFilesystemBlob(filepath.Join(tmpDir, "storage"))
	s := store.NewStore(db, blob)

	fake := &FakeProvider{}
	_, err := Analyse(context.Background(), fake, s, "job-1", "/nonexistent/image.jpg", 4)
	if err == nil {
		t.Fatal("expected error for nonexistent image")
	}
}

func TestProductAnalysisJSONRoundTrip(t *testing.T) {
	a := ProductAnalysis{
		ProductName:        "Test Product",
		Category:           "TestCategory",
		Materials:          []string{"material1", "material2"},
		Colours:            []string{"red", "blue"},
		ShapeDescription:   "round shape",
		DistinguishingFeatures: []string{"feature1"},
		SuggestedLifestyleSettings: []string{"home", "office", "outdoor", "travel"},
		Confidence:         0.85,
	}
	hex := "#FF5733"
	a.DominantColourHex = &hex

	data, err := json.Marshal(a)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var decoded ProductAnalysis
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if decoded.ProductName != a.ProductName {
		t.Errorf("productName round-trip: got %q, want %q", decoded.ProductName, a.ProductName)
	}
}
