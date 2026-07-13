package web

import (
	"testing"
)

func TestNoEmDashesInUI(t *testing.T) {
	uiTexts := []string{
		"Upload a photo",
		"Processing your listing",
		"Download your export",
	}
	for _, text := range uiTexts {
		if containsEmDash(text) {
			t.Errorf("UI text contains em dash: %s", text)
		}
	}
}

func containsEmDash(s string) bool {
	for _, r := range s {
		if r == '—' || r == '–' {
			return true
		}
	}
	return false
}
