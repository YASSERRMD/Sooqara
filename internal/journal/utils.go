package journal

import (
	"crypto/sha256"
	"fmt"
)

const maxDetailLen = 120

// truncateDetail limits detail to maxDetailLen characters.
func truncateDetail(s string) string {
	if len(s) > maxDetailLen {
		return s[:maxDetailLen]
	}
	return s
}

// sanitizeHash computes a short SHA-256 of the canonical request string.
func sanitizeHash(req string) string {
	h := sha256.Sum256([]byte(req))
	return fmt.Sprintf("%x", h[:8])
}
