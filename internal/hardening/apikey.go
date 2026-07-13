package hardening

import (
	"fmt"
	"strings"
)

const MinAPIKeyLength = 16
const MaxAPIKeyLength = 512

// ValidateAPIKey checks that an API key meets length and character requirements.
func ValidateAPIKey(key string) error {
	key = strings.TrimSpace(key)
	if key == "" {
		return fmt.Errorf("API key must not be empty")
	}
	if len(key) < MinAPIKeyLength {
		return fmt.Errorf("API key must be at least %d characters, got %d", MinAPIKeyLength, len(key))
	}
	if len(key) > MaxAPIKeyLength {
		return fmt.Errorf("API key must be at most %d characters, got %d", MaxAPIKeyLength, len(key))
	}
	return nil
}
