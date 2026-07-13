package pipeline

import (
	"testing"
)

func TestEmDashReMatches(t *testing.T) {
	tests := []struct {
		input    string
		hasDash  bool
	}{
		{"Hello — World", true},
		{"Hello – World", true},
		{"Hello World", false},
		{"A—B", true},
	}
	for _, tt := range tests {
		matches := emDashRe.MatchString(tt.input)
		if matches != tt.hasDash {
			t.Errorf("emDashRe.Match(%q) = %v, want %v", tt.input, matches, tt.hasDash)
		}
	}
}

func TestEnforceConstraintsNoWarnings(t *testing.T) {
	cs := &CopySet{
		Title:             "Clean Title",
		Bullets:           []string{"A", "B", "C", "D", "E"},
		ShortDescription:  "Short",
		LongDescription:   "Long",
		AltText:           "Alt",
		MetaDescription:   "Meta",
		Keywords:          []string{"k1", "k2", "k3", "k4", "k5", "k6"},
		Tone:              "practical",
	}
	var warning string
	err := enforceConstraints(cs, &warning)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if warning != "" {
		t.Errorf("expected no warning, got %q", warning)
	}
}

func TestEnforceConstraintsTitleTruncation(t *testing.T) {
	cs := &CopySet{
		Title:             "A Very Long Title That Exceeds The Maximum Allowed Character Limit For ECommerce Listings",
		Bullets:           []string{"A", "B", "C", "D", "E"},
		ShortDescription:  "Short",
		LongDescription:   "Long",
		AltText:           "Alt",
		MetaDescription:   "Meta",
		Keywords:          []string{"k1", "k2", "k3", "k4", "k5", "k6"},
		Tone:              "practical",
	}
	var warning string
	enforceConstraints(cs, &warning)
	if len(cs.Title) > MaxTitleLen {
		t.Errorf("title len = %d, want <= %d", len(cs.Title), MaxTitleLen)
	}
	if warning == "" {
		t.Error("expected warning for truncation")
	}
}
