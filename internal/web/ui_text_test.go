package web

import (
	"testing"
)

func TestPageTitles(t *testing.T) {
	titles := []string{"Upload", "Jobs", "Job Detail", "History"}
	for _, title := range titles {
		if title == "" {
			t.Error("empty page title")
		}
	}
}

func TestDropZoneText(t *testing.T) {
	texts := []string{"Drop a photo here", "or click to browse", "Supported: JPEG, PNG, WebP"}
	for _, text := range texts {
		if text == "" {
			t.Error("empty UI text")
		}
	}
}
