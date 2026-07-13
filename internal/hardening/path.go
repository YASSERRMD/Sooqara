package hardening

import (
	"path/filepath"
	"strings"
)

// SafePath ensures a file path does not escape the base directory.
func SafePath(base, candidate string) bool {
	rel, err := filepath.Rel(base, filepath.Clean(candidate))
	if err != nil {
		return false
	}
	return rel != ".." && !filepath.IsAbs(rel) && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}
