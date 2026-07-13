package release

import (
	"fmt"
	"os"
	"path/filepath"
)

// WriteChecksumsFile generates a SHA256SUMS file from artifacts in a directory.
func WriteChecksumsFile(artifactsDir, outputPath string) error {
	files, err := ListArtifacts(artifactsDir)
	if err != nil {
		return fmt.Errorf("list artifacts: %w", err)
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("create checksums file: %w", err)
	}
	defer f.Close()

	for _, file := range files {
		fullPath := filepath.Join(artifactsDir, file)
		hash, err := SHA256Checksum(fullPath)
		if err != nil {
			return fmt.Errorf("checksum %s: %w", file, err)
		}
		if _, err := fmt.Fprintf(f, "%s  %s\n", hash, file); err != nil {
			return fmt.Errorf("write checksum line: %w", err)
		}
	}
	return nil
}
