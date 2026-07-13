package release

import (
	"strings"
	"testing"
)

func TestFormatChangelogBasic(t *testing.T) {
	entries := []ChangelogEntry{
		{Type: "feat", Message: "add vision analysis"},
		{Type: "fix", Message: "fix limiter deadlock"},
	}
	out := FormatChangelog(entries)
	if !strings.Contains(out, "[feat]") {
		t.Error("expected [feat] in output")
	}
	if !strings.Contains(out, "[fix]") {
		t.Error("expected [fix] in output")
	}
	if !strings.Contains(out, "add vision analysis") {
		t.Error("expected message in output")
	}
}

func TestFormatChangelogEmpty(t *testing.T) {
	out := FormatChangelog(nil)
	if !strings.Contains(out, "## Changes") {
		t.Error("expected header in empty changelog")
	}
}

func TestCategorizeEntries(t *testing.T) {
	entries := []ChangelogEntry{
		{Type: "feat", Message: "a"},
		{Type: "fix", Message: "b"},
		{Type: "feat", Message: "c"},
	}
	cat := CategorizeEntries(entries)
	if len(cat["feat"]) != 2 {
		t.Errorf("expected 2 feat entries, got %d", len(cat["feat"]))
	}
	if len(cat["fix"]) != 1 {
		t.Errorf("expected 1 fix entry, got %d", len(cat["fix"]))
	}
}

func TestCategorizeEntriesEmpty(t *testing.T) {
	cat := CategorizeEntries(nil)
	if cat == nil {
		t.Fatal("expected non-nil map")
	}
	if len(cat) != 0 {
		t.Errorf("expected empty map, got %d entries", len(cat))
	}
}
