package release

import (
	"fmt"
	"os"
	"path/filepath"
)

// PrepareReleaseDir creates a staging directory for release artifacts.
func PrepareReleaseDir(base, version string) (string, error) {
	dir := filepath.Join(base, "dist", version)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("create release dir: %w", err)
	}
	return dir, nil
}

// CopyArtifact copies a single file into the release directory.
func CopyArtifact(src, destDir string) (string, error) {
	data, err := os.ReadFile(src)
	if err != nil {
		return "", fmt.Errorf("read artifact %s: %w", src, err)
	}
	name := filepath.Base(src)
	dest := filepath.Join(destDir, name)
	if err := os.WriteFile(dest, data, 0644); err != nil {
		return "", fmt.Errorf("write artifact %s: %w", dest, err)
	}
	return dest, nil
}

// ListArtifacts returns all files in a directory.
func ListArtifacts(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var files []string
	for _, e := range entries {
		if !e.IsDir() {
			files = append(files, e.Name())
		}
	}
	return files, nil
}
