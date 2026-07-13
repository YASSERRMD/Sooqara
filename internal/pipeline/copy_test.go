package pipeline

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/yasserrmd/sooqara/internal/provider"
	"github.com/yasserrmd/sooqara/internal/store"
)

func TestGenerateCopyToolCallValid(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	db, _ := store.Open(dbPath)
	defer db.Close()
	blob := store.NewFilesystemBlob(filepath.Join(tmpDir, "storage"))
	s := store.NewStore(db, blob)

	hint := "test"
	job := store.NewJob("/source.jpg", hint, "", 4)
	store.CreateJob(s.DB, job)

	analysis := &ProductAnalysis{
		ProductName:        "Red Sneaker",
		Category:           "Footwear",
		Materials:          []string{"leather"},
		Colours:            []string{"red"},
		ShapeDescription:   "low-top",
		DistinguishingFeatures: []string{"white sole"},
	}

	fake := &FakeProvider{
		ChatFn: func(ctx context.Context, req provider.ChatRequest) (provider.ChatResponse, error) {
			args := `{"title":"Red Leather Sneakers","bullets":["Bullet 1","Bullet 2","Bullet 3","Bullet 4","Bullet 5"],"short_description":"Great shoes.","long_description":"Long description.","alt_text":"Red sneakers","meta_description":"Buy red sneakers.","keywords":["shoes","red","leather","sneakers","footwear","casual"],"tone":"clear and practical"}`
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
	if artifact == nil {
		t.Fatal("expected artifact")
	}
	if artifact.Kind != store.ArtifactCopy {
		t.Errorf("kind = %s, want copy", artifact.Kind)
	}
	if job.Warning == nil {
		t.Log("no warning set (expected for banned phrase in title)")
	}
}

func TestGenerateCopyEmDashStripped(t *testing.T) {
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
			args := `{"title":"Test\u2014Product","bullets":["A\u2014B","C","D","E","F"],"short_description":"Short\u2014desc","long_description":"Long\u2014desc","alt_text":"Alt\u2014text","meta_description":"Meta\u2014desc","keywords":["kw1","kw2","kw3","kw4","kw5","kw6"],"tone":"practical"}`
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

func TestGenerateCopyBannedPhrase(t *testing.T) {
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
			args := `{"title":"Revolutionary Product","bullets":["A","B","C","D","E"],"short_description":"Short desc","long_description":"Long desc","alt_text":"Alt text","meta_description":"Meta desc","keywords":["kw1","kw2","kw3","kw4","kw5","kw6"],"tone":"practical"}`
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

func TestTruncateAtWord(t *testing.T) {
	tests := []struct {
		input    string
		max      int
		expected string
	}{
		{"Hello World Foo", 5, "Hello"},
		{"Hello World Foo", 20, "Hello World Foo"},
		{"NoSpacesHere", 5, "NoSpa"},
	}
	for _, tt := range tests {
		got := truncateAtWord(tt.input, tt.max)
		if got != tt.expected {
			t.Errorf("truncateAtWord(%q, %d) = %q, want %q", tt.input, tt.max, got, tt.expected)
		}
	}
}

func TestBannedPhrases(t *testing.T) {
	for _, bp := range bannedPhrases {
		if bp == "" {
			t.Error("banned phrase is empty")
		}
	}
	if len(bannedPhrases) < 5 {
		t.Errorf("expected at least 5 banned phrases, got %d", len(bannedPhrases))
	}
}

func TestCopySetConstants(t *testing.T) {
	if MaxTitleLen != 60 {
		t.Errorf("MaxTitleLen = %d, want 60", MaxTitleLen)
	}
	if ExpectedBulletCount != 5 {
		t.Errorf("ExpectedBulletCount = %d, want 5", ExpectedBulletCount)
	}
	if DefaultTone != "clear and practical" {
		t.Errorf("DefaultTone = %q, want clear and practical", DefaultTone)
	}
}

func TestGenerateCopyJSONFallback(t *testing.T) {
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

	callCount := 0
	fake := &FakeProvider{
		ChatFn: func(ctx context.Context, req provider.ChatRequest) (provider.ChatResponse, error) {
			callCount++
			if callCount == 1 {
				return provider.ChatResponse{
					Choices: []provider.Choice{
						{Message: provider.Message{Role: "assistant", Content: "{broken"}},
					},
				}, nil
			}
			args := `{"title":"Fixed Title","bullets":["A","B","C","D","E"],"short_description":"Short","long_description":"Long","alt_text":"Alt","meta_description":"Meta","keywords":["kw1","kw2","kw3","kw4","kw5","kw6"],"tone":"practical"}`
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
	if callCount != 2 {
		t.Errorf("expected 2 calls, got %d", callCount)
	}
}
