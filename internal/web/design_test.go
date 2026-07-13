package web

import (
	"testing"
)

func TestDesignDirectionMinimal(t *testing.T) {
	directions := []string{"minimal", "editorial", "generous whitespace", "one accent colour"}
	for _, d := range directions {
		if d == "" {
			t.Error("empty design direction")
		}
	}
}

func TestDesignDirectionExcludes(t *testing.T) {
	excluded := []string{"gradients", "glassmorphism", "card shadows"}
	for _, e := range excluded {
		if e == "" {
			t.Error("empty exclusion")
		}
	}
}
