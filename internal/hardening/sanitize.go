// Package hardening provides input validation and security checks for Sooqara.
package hardening

import "strings"

// SanitizeTitle removes control characters and trims whitespace from a title.
func SanitizeTitle(s string) string {
	var cleaned strings.Builder
	for _, r := range s {
		if (r >= ' ' && r <= '~') || r == '\n' || r == '\t' {
			cleaned.WriteRune(r)
		}
	}
	return strings.TrimSpace(cleaned.String())
}
