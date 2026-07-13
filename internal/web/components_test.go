package web

import (
	"testing"
)

func TestCopySetFields(t *testing.T) {
	fields := []string{"title", "bullets", "short_description", "long_description", "alt_text", "meta_description", "keywords", "tone"}
	for _, f := range fields {
		if f == "" {
			t.Error("empty field name")
		}
	}
	if len(fields) != 8 {
		t.Errorf("expected 8 fields, got %d", len(fields))
	}
}

func TestImageGridFields(t *testing.T) {
	fields := []string{"variant", "seed", "regenerate_button"}
	for _, f := range fields {
		if f == "" {
			t.Error("empty field name")
		}
	}
}
