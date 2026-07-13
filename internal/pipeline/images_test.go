package pipeline

import (
	"context"
	"testing"

	"github.com/yasserrmd/sooqara/internal/store"
)

func ptrInt64(i int64) *int64 { return &i }

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

func TestSeedLockedVariants(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := tmpDir + "/test.db"
	db, _ := store.Open(dbPath)
	defer db.Close()

	job := store.NewJob(tmpDir+"/src.jpg", "test", "", 3)
	job.Seed = ptrInt64(42)
	store.CreateJob(db, job)

	// Verify seeds are canonical_seed + i
	for i := 0; i < 3; i++ {
		expectedSeed := int64(42 + i)
		if *job.Seed+int64(i) != expectedSeed {
			t.Errorf("seed for variant %d = %d, want %d", i, *job.Seed+int64(i), expectedSeed)
		}
	}
}

func TestVariantCountZero(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := tmpDir + "/test.db"
	db, _ := store.Open(dbPath)
	defer db.Close()
	blob := store.NewFilesystemBlob(tmpDir + "/storage")
	s := store.NewStore(db, blob)

	job := store.NewJob(tmpDir+"/src.jpg", "test", "", 0)
	store.CreateJob(s.DB, job)

	analysis := &ProductAnalysis{ProductName: "X", Category: "C", ShapeDescription: "D"}
	fake := &FakeProvider{}

	arts, err := GenerateVariants(context.Background(), fake, s, job, analysis)
	if err != nil {
		t.Fatalf("GenerateVariants failed: %v", err)
	}
	if len(arts) != 0 {
		t.Errorf("expected 0 artifacts for variant_count=0, got %d", len(arts))
	}
}
