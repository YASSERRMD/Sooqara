package release

import (
	"fmt"
	"strings"
)

// ChangelogEntry represents a single changelog line.
type ChangelogEntry struct {
	Type    string // "feat", "fix", "docs", "chore"
	Message string
}

// FormatChangelog renders a slice of entries into a markdown changelog section.
func FormatChangelog(entries []ChangelogEntry) string {
	var sb strings.Builder
	sb.WriteString("## Changes\n\n")
	for _, e := range entries {
		sb.WriteString(fmt.Sprintf("- [%s] %s\n", e.Type, e.Message))
	}
	return sb.String()
}

// CategorizeEntries groups entries by type.
func CategorizeEntries(entries []ChangelogEntry) map[string][]string {
	cat := make(map[string][]string)
	for _, e := range entries {
		cat[e.Type] = append(cat[e.Type], e.Message)
	}
	return cat
}
