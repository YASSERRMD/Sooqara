package release

import (
	"fmt"
	"strings"
)

// GenerateReleaseNote creates a GitHub-style release note from PackageInfo.
func GenerateReleaseNote(pkg PackageInfo) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# Release %s\n\n", pkg.Version))
	sb.WriteString(fmt.Sprintf("**Platform:** %s/%s\n\n", pkg.Platform, pkg.Arch))
	sb.WriteString(fmt.Sprintf("**Description:** %s\n\n", pkg.Description))

	if len(pkg.Artifacts) > 0 {
		sb.WriteString("## Artifacts\n\n")
		for _, a := range pkg.Artifacts {
			checksum := ""
			if pkg.Checksums != nil {
				checksum = pkg.Checksums[a]
			}
			sb.WriteString(fmt.Sprintf("- `%s` (SHA256: `%s`)\n", a, checksum))
		}
		sb.WriteString("\n")
	}

	sb.WriteString(fmt.Sprintf("Generated: %s\n", pkg.BuildTime.Format("2006-01-02")))
	return sb.String()
}
