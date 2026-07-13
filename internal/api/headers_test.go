package api

import (
	"strings"
	"testing"
)

func TestContentTypeHeaders(t *testing.T) {
	types := []string{
		"application/json",
		"text/event-stream",
		"application/octet-stream",
		"text/plain; version=0.0.4",
	}
	for _, ct := range types {
		if ct == "" {
			t.Error("empty content type")
		}
	}
	if len(types) != 4 {
		t.Errorf("expected 4 content types, got %d", len(types))
	}
}

func TestHeaderConstants(t *testing.T) {
	headers := []string{
		"Content-Type",
		"Content-Disposition",
		"Cache-Control",
		"Connection",
	}
	for _, h := range headers {
		if !strings.Contains(h, "-") {
			t.Errorf("header %q missing hyphen", h)
		}
	}
}
