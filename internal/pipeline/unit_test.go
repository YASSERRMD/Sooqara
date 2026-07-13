package pipeline

import (
	"testing"
)

func TestStripFences(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{"```\n{\"key\":1}\n```", `{"key":1}`},
		{"```json\n{\"key\":1}\n```", `{"key":1}`},
		{"{\"key\":1}", `{"key":1}`},
		{"  ```\n{\"key\":1}\n```  ", `{"key":1}`},
		{"no fences here", "no fences here"},
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

func TestValidateMissingProductName(t *testing.T) {
	a := &ProductAnalysis{
		Category:         "Cat",
		ShapeDescription: "Desc",
	}
	if err := validateAnalysis(a, 1); err == nil {
		t.Fatal("expected error for missing product_name")
	}
}

func TestValidateMissingCategory(t *testing.T) {
	a := &ProductAnalysis{
		ProductName:      "Test",
		ShapeDescription: "Desc",
	}
	if err := validateAnalysis(a, 1); err == nil {
		t.Fatal("expected error for missing category")
	}
}

func TestValidateMissingShapeDescription(t *testing.T) {
	a := &ProductAnalysis{
		ProductName: "Test",
		Category:    "Cat",
	}
	if err := validateAnalysis(a, 1); err == nil {
		t.Fatal("expected error for missing shape_description")
	}
}
