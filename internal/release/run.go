package release

import (
	"fmt"
	"path/filepath"
	"time"
)

// RunRelease prepares, packages, and generates artifacts for a release.
func RunRelease(baseDir, version, commit string) error {
	if err := CheckEnvironment(); err != nil {
		return fmt.Errorf("environment check: %w", err)
	}

	releaseDir, err := PrepareReleaseDir(baseDir, version)
	if err != nil {
		return fmt.Errorf("prepare release dir: %w", err)
	}

	pkg := PackageInfo{
		Name:        "sooqara",
		Version:     version,
		Platform:    "auto",
		Arch:        "auto",
		BuildTime:   time.Now(),
		Description: "E-commerce listing factory",
	}

	manifestPath := filepath.Join(releaseDir, "manifest.json")
	if err := WriteManifest(manifestPath, pkg); err != nil {
		return fmt.Errorf("write manifest: %w", err)
	}

	artifacts, _ := ListArtifacts(releaseDir)
	pkg.Artifacts = artifacts

	_ = commit

	return nil
}
