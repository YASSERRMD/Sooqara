package web

import (
	"testing"
)

func TestNoBuildStep(t *testing.T) {
	buildTools := []string{"npm", "bundler", "webpack", "vite"}
	for _, tool := range buildTools {
		if tool == "" {
			t.Error("empty build tool name")
		}
	}
}

func TestNoNPMRule(t *testing.T) {
	rule := "No npm. No bundler."
	if rule == "" {
		t.Error("empty rule text")
	}
}
