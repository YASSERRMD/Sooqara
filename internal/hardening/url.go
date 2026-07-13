package hardening

import (
	"fmt"
	"net/url"
	"strings"
)

const MaxURLLength = 2048

// ValidateURL checks that a string is a well-formed HTTP or HTTPS URL.
func ValidateURL(raw string) error {
	if len(raw) > MaxURLLength {
		return fmt.Errorf("URL must be at most %d characters", MaxURLLength)
	}
	parsed, err := url.ParseRequestURI(raw)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return fmt.Errorf("URL scheme must be http or https, got %q", parsed.Scheme)
	}
	if parsed.Host == "" {
		return fmt.Errorf("URL must have a host")
	}
	return nil
}

// SanitizeURL trims whitespace from a URL string.
func SanitizeURL(raw string) string {
	return strings.TrimSpace(raw)
}
