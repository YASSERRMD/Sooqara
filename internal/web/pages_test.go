package web

import (
	"testing"
)

func TestPageCount(t *testing.T) {
	pages := []string{"/", "/jobs/{id}", "/jobs"}
	if len(pages) != 3 {
		t.Errorf("expected 3 pages, got %d", len(pages))
	}
}
