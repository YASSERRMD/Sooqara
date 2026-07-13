package release

import (
	"strings"
	"testing"
)

func TestMetadataZeroValue(t *testing.T) {
	var m Metadata
	out := m.String()
	if out == "" {
		t.Error("expected non-empty output for zero-value Metadata")
	}
}

func TestFormatChangelogPreservesOrder(t *testing.T) {
	entries := []ChangelogEntry{
		{Type: "feat", Message: "first"},
		{Type: "fix", Message: "second"},
		{Type: "docs", Message: "third"},
	}
	out := FormatChangelog(entries)
	posFirst := strings.Index(out, "first")
	posSecond := strings.Index(out, "second")
	posThird := strings.Index(out, "third")
	if posFirst >= posSecond || posSecond >= posThird {
		t.Error("expected entries to preserve insertion order")
	}
}
