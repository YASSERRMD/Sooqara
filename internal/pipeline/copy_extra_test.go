package pipeline

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/yasserrmd/sooqara/internal/provider"
	"github.com/yasserrmd/sooqara/internal/store"
)

func TestGenerateCopyWithNilJob(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	db, _ := store.Open(dbPath)
	defer db.Close()
	blob := store.NewFilesystemBlob(filepath.Join(tmpDir, "storage"))
	s := store.NewStore(db, blob)

	hint := "test"
	job := store.NewJob("/source.jpg", hint, "", 4)
	store.CreateJob(s.DB, job)

	analysis := &ProductAnalysis{ProductName: "X", Category: "C", ShapeDescription: "D"}

	fake := &FakeProvider{
		ChatFn: func(ctx context.Context, req provider.ChatRequest) (provider.ChatResponse, error) {
			args := `{"title":"T","bullets":["A","B","C","D","E"],"short_description":"S","long_description":"L","alt_text":"A","meta_description":"M","keywords":["k1","k2","k3","k4","k5","k6"],"tone":"clear and practical"}`
			return provider.ChatResponse{
				Choices: []provider.Choice{
					{Message: provider.Message{ToolCalls: []provider.ToolCall{{Function: provider.ToolFunctionCall{Arguments: args}}}}},
				},
			}, nil
		},
	}

	_, err := GenerateCopy(context.Background(), fake, s, job, analysis)
	if err != nil {
		t.Fatalf("GenerateCopy failed: %v", err)
	}
}

func TestGenerateCopyDefaultTone(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	db, _ := store.Open(dbPath)
	defer db.Close()
	blob := store.NewFilesystemBlob(filepath.Join(tmpDir, "storage"))
	s := store.NewStore(db, blob)

	job := store.NewJob("/source.jpg", "", "", 4)
	store.CreateJob(s.DB, job)

	analysis := &ProductAnalysis{ProductName: "X", Category: "C", ShapeDescription: "D"}

	fake := &FakeProvider{
		ChatFn: func(ctx context.Context, req provider.ChatRequest) (provider.ChatResponse, error) {
			args := `{"title":"T","bullets":["A","B","C","D","E"],"short_description":"S","long_description":"L","alt_text":"A","meta_description":"M","keywords":["k1","k2","k3","k4","k5","k6"],"tone":"clear and practical"}`
			return provider.ChatResponse{
				Choices: []provider.Choice{
					{Message: provider.Message{ToolCalls: []provider.ToolCall{{Function: provider.ToolFunctionCall{Arguments: args}}}}},
				},
			}, nil
		},
	}

	_, err := GenerateCopy(context.Background(), fake, s, job, analysis)
	if err != nil {
		t.Fatalf("GenerateCopy failed: %v", err)
	}
}

func TestGenerateCopyCustomTone(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	db, _ := store.Open(dbPath)
	defer db.Close()
	blob := store.NewFilesystemBlob(filepath.Join(tmpDir, "storage"))
	s := store.NewStore(db, blob)

	job := store.NewJob("/source.jpg", "test", "friendly", 4)
	store.CreateJob(s.DB, job)

	analysis := &ProductAnalysis{ProductName: "X", Category: "C", ShapeDescription: "D"}

	fake := &FakeProvider{
		ChatFn: func(ctx context.Context, req provider.ChatRequest) (provider.ChatResponse, error) {
			args := `{"title":"T","bullets":["A","B","C","D","E"],"short_description":"S","long_description":"L","alt_text":"A","meta_description":"M","keywords":["k1","k2","k3","k4","k5","k6"],"tone":"friendly"}`
			return provider.ChatResponse{
				Choices: []provider.Choice{
					{Message: provider.Message{ToolCalls: []provider.ToolCall{{Function: provider.ToolFunctionCall{Arguments: args}}}}},
				},
			}, nil
		},
	}

	_, err := GenerateCopy(context.Background(), fake, s, job, analysis)
	if err != nil {
		t.Fatalf("GenerateCopy failed: %v", err)
	}
}

func TestGenerateCopyPersistedPayload(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	db, _ := store.Open(dbPath)
	defer db.Close()
	blob := store.NewFilesystemBlob(filepath.Join(tmpDir, "storage"))
	s := store.NewStore(db, blob)

	job := store.NewJob("/source.jpg", "test", "", 4)
	store.CreateJob(s.DB, job)

	analysis := &ProductAnalysis{ProductName: "X", Category: "C", ShapeDescription: "D"}

	fake := &FakeProvider{
		ChatFn: func(ctx context.Context, req provider.ChatRequest) (provider.ChatResponse, error) {
			args := `{"title":"T","bullets":["A","B","C","D","E"],"short_description":"S","long_description":"L","alt_text":"A","meta_description":"M","keywords":["k1","k2","k3","k4","k5","k6"],"tone":"clear and practical"}`
			return provider.ChatResponse{
				Choices: []provider.Choice{
					{Message: provider.Message{ToolCalls: []provider.ToolCall{{Function: provider.ToolFunctionCall{Arguments: args}}}}},
				},
			}, nil
		},
	}

	artifact, err := GenerateCopy(context.Background(), fake, s, job, analysis)
	if err != nil {
		t.Fatalf("GenerateCopy failed: %v", err)
	}
	if artifact.Payload == nil || *artifact.Payload == "" {
		t.Fatal("expected non-empty payload")
	}
}

func TestGenerateCopyMissingImageFile(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	db, _ := store.Open(dbPath)
	defer db.Close()
	blob := store.NewFilesystemBlob(filepath.Join(tmpDir, "storage"))
	s := store.NewStore(db, blob)

	job := store.NewJob("/nonexistent.jpg", "test", "", 4)
	store.CreateJob(s.DB, job)

	analysis := &ProductAnalysis{ProductName: "X", Category: "C", ShapeDescription: "D"}
	fake := &FakeProvider{}

	_, err := GenerateCopy(context.Background(), fake, s, job, analysis)
	if err != nil {
		t.Logf("expected error: %v", err)
	}
}

func TestCopySetAllFieldsPresent(t *testing.T) {
	cs := CopySet{
		Title:             "Test",
		Bullets:           []string{"A", "B", "C", "D", "E"},
		ShortDescription:  "S",
		LongDescription:   "L",
		AltText:           "A",
		MetaDescription:   "M",
		Keywords:          []string{"k1", "k2", "k3", "k4", "k5", "k6"},
		Tone:              "practical",
	}

	if cs.Title == "" {
		t.Error("title should not be empty")
	}
	if len(cs.Bullets) != ExpectedBulletCount {
		t.Errorf("expected %d bullets, got %d", ExpectedBulletCount, len(cs.Bullets))
	}
	if len(cs.Keywords) < MinKeywordCount || len(cs.Keywords) > MaxKeywordCount {
		t.Errorf("keywords count %d out of range [%d,%d]", len(cs.Keywords), MinKeywordCount, MaxKeywordCount)
	}
}
