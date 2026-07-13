package hardening

import (
	"strings"
	"testing"
)

func TestAllValidatorsBasicFlow(t *testing.T) {
	if err := ValidateTitle("Test Product"); err != nil {
		t.Errorf("title validation failed: %v", err)
	}
	if err := ValidateURL("https://example.com"); err != nil {
		t.Errorf("URL validation failed: %v", err)
	}
	if err := ValidateAPIKey(strings.Repeat("a", 32)); err != nil {
		t.Errorf("API key validation failed: %v", err)
	}
	if err := ValidateTags([]string{"tag1", "tag2"}); err != nil {
		t.Errorf("tags validation failed: %v", err)
	}
	if err := ValidateCopyLength("short text"); err != nil {
		t.Errorf("copy length validation failed: %v", err)
	}
	if err := ValidateHeaderName("content-type"); err != nil {
		t.Errorf("header validation failed: %v", err)
	}
}

func TestAllValidatorsErrorFormats(t *testing.T) {
	tests := []struct {
		name   string
		getErr func() error
		substr string
	}{
		{"empty title", func() error { return ValidateTitle("") }, "title must not be empty"},
		{"bad URL", func() error { return ValidateURL("mailto:x") }, "http or https"},
		{"short key", func() error { return ValidateAPIKey("abc") }, "at least"},
		{"too many tags", func() error { return ValidateTags(make([]string, MaxTagsPerListing+1)) }, "too many tags"},
		{"copy too long", func() error { return ValidateCopyLength(strings.Repeat("x", MaxCopyLength+1)) }, "exceeds"},
		{"bad header", func() error { return ValidateHeaderName("cookie") }, "not allowed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.getErr()
			if err == nil {
				t.Fatalf("expected non-nil error for %s", tt.name)
			}
			if !strings.Contains(err.Error(), tt.substr) {
				t.Errorf("error for %s should contain %q, got %q", tt.name, tt.substr, err.Error())
			}
		})
	}
}
