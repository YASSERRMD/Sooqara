package pipeline

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/yasserrmd/sooqara/internal/provider"
	"github.com/yasserrmd/sooqara/internal/store"
)

func TestAnalyseValidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	db, _ := store.Open(dbPath)
	defer db.Close()
	blob := store.NewFilesystemBlob(filepath.Join(tmpDir, "storage"))
	s := store.NewStore(db, blob)

	hint := "test"
	job := store.NewJob("/source.jpg", hint, "", 4)
	store.CreateJob(db, job)

	sourcePath := filepath.Join(tmpDir, "source.jpg")
	os.WriteFile(sourcePath, []byte("fake image"), 0644)

	analysis := ProductAnalysis{
		ProductName:        "Red Sneaker",
		Category:           "Footwear",
		Materials:          []string{"leather", "rubber"},
		Colours:            []string{"red"},
		ShapeDescription:   "low-top sneaker",
		DistinguishingFeatures: []string{"white sole", "lace-up"},
		SuggestedLifestyleSettings: []string{"studio white", "beach sunset", "urban street", "mountain trail"},
		Confidence:         0.9,
	}

	fake := &FakeProvider{
		ChatFn: func(ctx context.Context, req provider.ChatRequest) (provider.ChatResponse, error) {
			return ValidJSON(analysis), nil
		},
	}

	artifact, err := Analyse(context.Background(), fake, s, job.ID, sourcePath, 4)
	if err != nil {
		t.Fatalf("Analyse failed: %v", err)
	}
	if artifact == nil {
		t.Fatal("expected artifact")
	}
	if artifact.Kind != store.ArtifactAnalysis {
		t.Errorf("kind = %s, want %s", artifact.Kind, store.ArtifactAnalysis)
	}
}

func TestAnalyseMalformedNeverRepairs(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	db, _ := store.Open(dbPath)
	defer db.Close()
	blob := store.NewFilesystemBlob(filepath.Join(tmpDir, "storage"))
	s := store.NewStore(db, blob)

	hint := "test"
	job := store.NewJob("/source.jpg", hint, "", 4)
	store.CreateJob(db, job)

	sourcePath := filepath.Join(tmpDir, "source.jpg")
	os.WriteFile(sourcePath, []byte("fake image"), 0644)

	fake := &FakeProvider{
		ChatFn: func(ctx context.Context, req provider.ChatRequest) (provider.ChatResponse, error) {
			return MalformedJSON(), nil
		},
	}

	_, err := Analyse(context.Background(), fake, s, job.ID, sourcePath, 4)
	if err == nil {
		t.Fatal("expected error for malformed JSON")
	}
}
