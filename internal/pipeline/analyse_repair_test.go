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

func TestAnalyseFenceWrappedJSON(t *testing.T) {
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
		ProductName:        "Blue Mug",
		Category:           "Kitchen",
		Materials:          []string{"ceramic"},
		Colours:            []string{"blue"},
		ShapeDescription:   "cylindrical mug",
		DistinguishingFeatures: []string{"handle on right"},
		SuggestedLifestyleSettings: []string{"kitchen counter", "coffee shop", "outdoor picnic", "desk workspace"},
		Confidence:         0.85,
	}

	callCount := 0
	fake := &FakeProvider{
		ChatFn: func(ctx context.Context, req provider.ChatRequest) (provider.ChatResponse, error) {
			callCount++
			data, _ := json.Marshal(analysis)
			if callCount == 1 {
				return provider.ChatResponse{
					Model: "agnes-2.0-flash",
					Choices: []provider.Choice{
						{Message: provider.Message{Role: "assistant", Content: "```json\n" + string(data) + "\n``` garbage extra"}},
					},
				}, nil
			}
			return ValidJSON(analysis), nil
		},
	}

	_, err := Analyse(context.Background(), fake, s, job.ID, sourcePath, 4)
	if err != nil {
		t.Fatalf("Analyse failed: %v", err)
	}
	if callCount != 2 {
		t.Errorf("expected 2 calls (initial + repair), got %d", callCount)
	}
}

func TestAnalyseTooFewSettings(t *testing.T) {
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
		ProductName:        "Watch",
		Category:           "Accessories",
		Materials:          []string{"steel"},
		Colours:            []string{"silver"},
		ShapeDescription:   "round face",
		DistinguishingFeatures: []string{"leather strap"},
		SuggestedLifestyleSettings: []string{"office"},
		Confidence:         0.7,
	}

	callCount := 0
	fake := &FakeProvider{
		ChatFn: func(ctx context.Context, req provider.ChatRequest) (provider.ChatResponse, error) {
			callCount++
			if callCount == 1 {
				return TooFewSettings(analysis), nil
			}
			analysis.SuggestedLifestyleSettings = []string{"office", "boardroom", "casual friday", "evening gala"}
			return ValidJSON(analysis), nil
		},
	}

	_, err := Analyse(context.Background(), fake, s, job.ID, sourcePath, 4)
	if err != nil {
		t.Fatalf("Analyse failed: %v", err)
	}
	if callCount != 2 {
		t.Errorf("expected 2 calls (initial + repair), got %d", callCount)
	}
}
