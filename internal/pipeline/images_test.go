package pipeline

import (
	"context"
	"fmt"
	"testing"

	"github.com/yasserrmd/sooqara/internal/provider"
	"github.com/yasserrmd/sooqara/internal/store"
)

func ptrInt64(i int64) *int64 { return &i }
	tmpDir := t.TempDir()
	dbPath := tmpDir + "/test.db"
	db, _ := store.Open(dbPath)
	defer db.Close()
	blob := store.NewFilesystemBlob(tmpDir + "/storage")
	s := store.NewStore(db, blob)

	job := store.NewJob(tmpDir+"/source.jpg", "test", "", 2)
	job.Seed = ptrInt64(100)
	store.CreateJob(s.DB, job)

	analysis := &ProductAnalysis{
		ProductName:        "Shoe",
		Category:           "Footwear",
		Materials:          []string{"leather"},
		Colours:            []string{"brown"},
		ShapeDescription:   "sneaker",
		DistinguishingFeatures: []string{"lace-up"},
		SuggestedLifestyleSettings: []string{"studio", "outdoor"},
	}

	fake := &FakeProvider{
		GenerateImg: func(ctx context.Context, req provider.ImageRequest) (provider.ImageResponse, error) {
			return provider.ImageResponse{
				Images: []provider.Image{{URL: "http://example.com/img.png", Seed: req.Seed}},
			}, nil
		},
	}

	arts, err := GenerateVariants(context.Background(), fake, s, job, analysis)
	if err != nil {
		t.Fatalf("GenerateVariants failed: %v", err)
	}
	if len(arts) != 2 {
		t.Errorf("got %d artifacts, want 2", len(arts))
	}
}

func TestGenerateVariantsPartialFailure(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := tmpDir + "/test.db"
	db, _ := store.Open(dbPath)
	defer db.Close()
	blob := store.NewFilesystemBlob(tmpDir + "/storage")
	s := store.NewStore(db, blob)

	job := store.NewJob(tmpDir+"/source.jpg", "test", "", 3)
	job.Seed = ptrInt64(200)
	store.CreateJob(s.DB, job)

	analysis := &ProductAnalysis{
		ProductName:        "Watch",
		Category:           "Accessory",
		Materials:          []string{"steel"},
		Colours:            []string{"silver"},
		ShapeDescription:   "round",
		DistinguishingFeatures: []string{"leather strap"},
		SuggestedLifestyleSettings: []string{"office", "gym", "dinner"},
	}

	callCount := 0
	fake := &FakeProvider{
		GenerateImg: func(ctx context.Context, req provider.ImageRequest) (provider.ImageResponse, error) {
			callCount++
			if callCount == 2 {
				return provider.ImageResponse{}, fmt.Errorf("provider error")
			}
			return provider.ImageResponse{
				Images: []provider.Image{{URL: "http://example.com/img.png", Seed: req.Seed}},
			}, nil
		},
	}

	arts, err := GenerateVariants(context.Background(), fake, s, job, analysis)
	if err != nil {
		t.Fatalf("GenerateVariants failed: %v", err)
	}
	if len(arts) != 2 {
		t.Errorf("expected 2 artifacts (partial success), got %d", len(arts))
	}
}

func TestGenerateVariantsAllFail(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := tmpDir + "/test.db"
	db, _ := store.Open(dbPath)
	defer db.Close()
	blob := store.NewFilesystemBlob(tmpDir + "/storage")
	s := store.NewStore(db, blob)

	job := store.NewJob(tmpDir+"/source.jpg", "test", "", 2)
	job.Seed = ptrInt64(300)
	store.CreateJob(s.DB, job)

	analysis := &ProductAnalysis{ProductName: "X", Category: "C", ShapeDescription: "D"}

	fake := &FakeProvider{
		GenerateImg: func(ctx context.Context, req provider.ImageRequest) (provider.ImageResponse, error) {
			return provider.ImageResponse{}, fmt.Errorf("always fails")
		},
	}

	_, err := GenerateVariants(context.Background(), fake, s, job, analysis)
	if err == nil {
		t.Fatal("expected error when all variants fail")
	}
}

func TestBuildPrompt(t *testing.T) {
	analysis := &ProductAnalysis{
		ProductName:        "Red Shoe",
		Materials:          []string{"leather"},
		Colours:            []string{"red"},
		ShapeDescription:   "sneaker",
		DistinguishingFeatures: []string{"white sole"},
	}
	prompt := buildPrompt(analysis, "beach sunset")
	expected := "Red Shoe, leather, red, sneaker, white sole placed in beach sunset, " + StyleV1
	if prompt != expected {
		t.Errorf("prompt = %q, want %q", prompt, expected)
	}
}

func TestGetSettingOutOfRange(t *testing.T) {
	analysis := &ProductAnalysis{SuggestedLifestyleSettings: []string{"only one"}}
	setting := getSetting(analysis, 5)
	if setting != "neutral studio background" {
		t.Errorf("got %q, want neutral studio background", setting)
	}
}

func TestJoinStringsEmpty(t *testing.T) {
	if joinStrings(nil) != "" {
		t.Error("expected empty string for nil slice")
	}
	if joinStrings([]string{}) != "" {
		t.Error("expected empty string for empty slice")
	}
}

func TestJoinStringsSingle(t *testing.T) {
	got := joinStrings([]string{"a"})
	if got != "a" {
		t.Errorf("got %q, want a", got)
	}
}

func TestJoinStringsMultiple(t *testing.T) {
	got := joinStrings([]string{"a", "b", "c"})
	if got != "a, b, c" {
		t.Errorf("got %q, want a, b, c", got)
	}
}
