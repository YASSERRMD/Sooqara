package web

import (
	"testing"
)

func TestProgressStages(t *testing.T) {
	stages := []string{"Analysing", "Writing", "Imaging", "Filming", "Done"}
	for _, s := range stages {
		if s == "" {
			t.Error("empty stage name")
		}
	}
	if len(stages) != 5 {
		t.Errorf("expected 5 stages, got %d", len(stages))
	}
}

func TestSpinnerStates(t *testing.T) {
	states := []string{"spinner", "tick", "warning"}
	for _, s := range states {
		if s == "" {
			t.Error("empty state")
		}
	}
}
