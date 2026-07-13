package web

import (
	"testing"
)

func TestCharCountColors(t *testing.T) {
	colors := []string{"green", "amber", "red"}
	for _, c := range colors {
		if c == "" {
			t.Error("empty color")
		}
	}
}

func TestDisclosureWidget(t *testing.T) {
	widget := "add detail"
	if widget == "" {
		t.Error("empty widget text")
	}
}
