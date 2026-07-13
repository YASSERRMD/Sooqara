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

	artifact, err := Analyse(context.Background(), fake, s, "job-1", sourcePath, 4)
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

func TestAnalyseFenceWrappedJSON(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	db, _ := store.Open(dbPath)
	defer db.Close()
	blob := store.NewFilesystemBlob(filepath.Join(tmpDir, "storage"))
	s := store.NewStore(db, blob)

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
			if callCount == 1 {
				return FenceWrappedJSON(analysis), nil
			}
			// Second call (repair) should succeed with plain JSON
			return ValidJSON(analysis), nil
		},
	}

	_, err := Analyse(context.Background(), fake, s, "job-2", sourcePath, 4)
	if err != nil {
		t.Fatalf("Analyse failed: %v", err)
	}
	if callCount != 2 {
		t.Errorf("expected 2 calls (initial + repair), got %d", callCount)
	}
}

func TestAnalyseMalformedNeverRepairs(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	db, _ := store.Open(dbPath)
	defer db.Close()
	blob := store.NewFilesystemBlob(filepath.Join(tmpDir, "storage"))
	s := store.NewStore(db, blob)

	sourcePath := filepath.Join(tmpDir, "source.jpg")
	os.WriteFile(sourcePath, []byte("fake image"), 0644)

	fake := &FakeProvider{
		ChatFn: func(ctx context.Context, req provider.ChatRequest) (provider.ChatResponse, error) {
			return MalformedJSON(), nil
		},
	}

	_, err := Analyse(context.Background(), fake, s, "job-3", sourcePath, 4)
	if err == nil {
		t.Fatal("expected error for malformed JSON")
	}
}

func TestAnalyseTooFewSettings(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	db, _ := store.Open(dbPath)
	defer db.Close()
	blob := store.NewFilesystemBlob(filepath.Join(tmpDir, "storage"))
	s := store.NewStore(db, blob)

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

	_, err := Analyse(context.Background(), fake, s, "job-4", sourcePath, 4)
	if err != nil {
		t.Fatalf("Analyse failed: %v", err)
	}
	if callCount != 2 {
		t.Errorf("expected 2 calls (initial + repair), got %d", callCount)
	}
}

func TestStripFences(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{"```\n{\"key\":1}\n```", `{"key":1}`},
		{"```json\n{\"key\":1}\n```", `{"key":1}`},
		{"{\"key\":1}", `{"key":1}`},
		{"  ```\n{\"key\":1}\n```  ", `{"key":1}`},
	}
	for _, tt := range tests {
		got := stripFences(tt.input)
		if got != tt.want {
			t.Errorf("stripFences(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestValidateDominantColourHex(t *testing.T) {
	a := &ProductAnalysis{
		ProductName:      "Test",
		Category:         "Cat",
		ShapeDescription: "Desc",
	}
	hex := "not-a-hex"
	a.DominantColourHex = &hex

	if err := validateAnalysis(a, 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.DominantColourHex != nil {
		t.Error("expected nil for invalid hex, got non-nil")
	}
}

func TestValidateValidHex(t *testing.T) {
	a := &ProductAnalysis{
		ProductName:      "Test",
		Category:         "Cat",
		ShapeDescription: "Desc",
	}
	hex := "#FF5733"
	a.DominantColourHex = &hex

	if err := validateAnalysis(a, 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.DominantColourHex == nil || *a.DominantColourHex != "#FF5733" {
		t.Error("expected valid hex to be preserved")
	}
}
