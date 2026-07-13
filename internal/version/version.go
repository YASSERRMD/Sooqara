// Package version provides build-time version information for Sooqara.
package version

import "fmt"

// These are set via -ldflags at build time.
var (
	Version   = "dev"
	Commit    = "none"
	BuildTime = "unknown"
)

// Info returns a human-readable version string.
func Info() string {
	return Version
}

// BuildInfo returns a multi-line build info string.
func BuildInfo() string {
	return fmt.Sprintf("sooqara %s\ncommit %s\nbuilt %s", Version, Commit, BuildTime)
}
