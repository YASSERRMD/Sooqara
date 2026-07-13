package release

import "time"

// PackageInfo holds metadata for a release package.
type PackageInfo struct {
	Name        string
	Version     string
	Platform    string
	Arch        string
	BuildTime   time.Time
	Artifacts   []string
	Checksums   map[string]string
	Description string
}
