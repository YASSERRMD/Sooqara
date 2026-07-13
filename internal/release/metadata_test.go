package release

import (
	"strings"
	"testing"
)

func TestMetadataString(t *testing.T) {
	m := Metadata{
		Version:   "v1.0.0",
		Commit:    "abc123",
		BuildTime: "2025-01-01T00:00:00Z",
		GoVersion: "go1.23",
		Platform:  "linux/amd64",
	}
	out := m.String()
	expected := []string{"v1.0.0", "abc123", "2025-01-01T00:00:00Z", "go1.23", "linux/amd64"}
	for _, exp := range expected {
		if !strings.Contains(out, exp) {
			t.Errorf("expected %q in metadata string", exp)
		}
	}
}

func TestMetadataEmpty(t *testing.T) {
	m := Metadata{}
	out := m.String()
	if !strings.Contains(out, "Release Metadata") {
		t.Error("expected header in empty metadata")
	}
}
