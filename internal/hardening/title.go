package hardening

import "fmt"

const MaxTitleLength = 200

// ValidateTitle checks that a title is non-empty, within length limits, and free of control characters.
func ValidateTitle(title string) error {
	cleaned := SanitizeTitle(title)
	if cleaned == "" {
		return fmt.Errorf("title must not be empty")
	}
	if len(cleaned) > MaxTitleLength {
		return fmt.Errorf("title must be at most %d characters, got %d", MaxTitleLength, len(cleaned))
	}
	return nil
}
