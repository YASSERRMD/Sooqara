package web

import (
	"testing"
)

func TestGoldenFileFormats(t *testing.T) {
	formats := []string{"html", "css", "javascript"}
	for _, f := range formats {
		if f == "" {
			t.Error("empty format")
		}
	}
}

func TestTemplateExtensions(t *testing.T) {
	exts := []string{".html", ".tmpl"}
	for _, ext := range exts {
		if ext == "" {
			t.Error("empty extension")
		}
	}
}
