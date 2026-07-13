package pipeline

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/yasserrmd/sooqara/internal/provider"
	"github.com/yasserrmd/sooqara/internal/store"
)

func TestGenerateCopyNoToolCallFallback(t *testing.T) {
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
			return provider.ChatResponse{
				Choices: []provider.Choice{
					{Message: provider.Message{Role: "assistant", Content: `{"title":"JSON Title","bullets":["A","B","C","D","E"],"short_description":"S","long_description":"L","alt_text":"A","meta_description":"M","keywords":["k1","k2","k3","k4","k5","k6"],"tone":"practical"}`}},
				},
			}, nil
		},
	}

	_, err := GenerateCopy(context.Background(), fake, s, job, analysis)
	if err != nil {
		t.Fatalf("GenerateCopy failed: %v", err)
	}
}
