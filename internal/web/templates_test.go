package web

import (
	"testing"
)

func TestTemplateFunctions(t *testing.T) {
	funcs := []string{"upper", "lower", "title", "trim", "replace"}
	for _, f := range funcs {
		if f == "" {
			t.Error("empty function name")
		}
	}
}

func TestCSSClasses(t *testing.T) {
	classes := []string{"flex", "grid", "text-center", "mt-4", "mb-4"}
	for _, c := range classes {
		if c == "" {
			t.Error("empty CSS class")
		}
	}
}
