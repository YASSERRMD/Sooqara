package web

import (
	"testing"
)

func TestHTMXAttributes(t *testing.T) {
	attrs := []string{"hx-get", "hx-post", "hx-trigger", "hx-target", "hx-swap"}
	for _, attr := range attrs {
		if attr == "" {
			t.Error("empty HTMX attribute")
		}
	}
	if len(attrs) != 5 {
		t.Errorf("expected 5 attributes, got %d", len(attrs))
	}
}
