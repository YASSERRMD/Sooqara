// Package release provides build and packaging utilities for Sooqara.
package release

import (
	"fmt"
	"os"
	"time"
)

// GenerateBuildFlags produces -ldflags string for the current build.
func GenerateBuildFlags(version, commit string) string {
	buildTime := time.Now().UTC().Format(time.RFC3339)
	return fmt.Sprintf(
		"-X github.com/yasserrmd/sooqara/internal/version.Version=%s "+
			"-X github.com/yasserrmd/sooqara/internal/version.Commit=%s "+
			"-X github.com/yasserrmd/sooqara/internal/version.BuildTime=%s",
		version, commit, buildTime,
	)
}

// DefaultFlags reads VERSION and COMMIT from the environment and returns ldflags.
func DefaultFlags() string {
	return GenerateBuildFlags(
		getenvOr("VERSION", "dev"),
		getenvOr("COMMIT", "none"),
	)
}

func getenvOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
