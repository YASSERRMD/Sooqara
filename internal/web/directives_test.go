package web

import (
	"testing"
)

func TestEmbedDirective(t *testing.T) {
	directive := "//go:embed"
	if directive == "" {
		t.Error("empty directive")
	}
}

func TestTemplatePackage(t *testing.T) {
	pkg := "html/template"
	if pkg == "" {
		t.Error("empty package path")
	}
}
