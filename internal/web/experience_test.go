package web

import (
	"strings"
	"testing"
)

func TestNoExperienceWord(t *testing.T) {
	uiTexts := []string{
		"Upload your product photo",
		"View your listings",
		"Download exports",
	}
	for _, text := range uiTexts {
		if containsExperience(text) {
			t.Errorf("UI text contains 'experience': %s", text)
		}
	}
}

func containsExperience(s string) bool {
	lower := strings.ToLower(s)
	return strings.Contains(lower, "experience")
}
