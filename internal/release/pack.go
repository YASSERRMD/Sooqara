package release

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// PackTarball creates a .tar.gz archive from a list of files.
func PackTarball(archivePath, baseDir string, files []string) error {
	f, err := os.Create(archivePath)
	if err != nil {
		return fmt.Errorf("create archive: %w", err)
	}
	defer f.Close()

	gw := gzip.NewWriter(f)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	for _, file := range files {
		fullPath := filepath.Join(baseDir, file)
		finfo, err := os.Stat(fullPath)
		if err != nil {
			return fmt.Errorf("stat %s: %w", file, err)
		}

		header, err := tar.FileInfoHeader(finfo, "")
		if err != nil {
			return fmt.Errorf("header for %s: %w", file, err)
		}
		header.Name = file

		if err := tw.WriteHeader(header); err != nil {
			return fmt.Errorf("write header for %s: %w", file, err)
		}

		if finfo.Mode().IsRegular() {
			data, err := os.Open(fullPath)
			if err != nil {
				return fmt.Errorf("open %s: %w", file, err)
			}
			if _, err := io.Copy(tw, data); err != nil {
				data.Close()
				return fmt.Errorf("copy %s: %w", file, err)
			}
			data.Close()
		}
	}
	return nil
}

// ListFiles returns all regular files under baseDir matching patterns.
func ListFiles(baseDir string, patterns ...string) ([]string, error) {
	var matches []string
	for _, pat := range patterns {
		found, err := filepath.Glob(filepath.Join(baseDir, pat))
		if err != nil {
			return nil, err
		}
		for _, f := range found {
			info, err := os.Stat(f)
			if err == nil && !info.IsDir() {
				rel, _ := filepath.Rel(baseDir, f)
				matches = append(matches, rel)
			}
		}
	}
	return matches, nil
}

// SanitizeFilename removes path separators and trims unsafe characters.
func SanitizeFilename(name string) string {
	return strings.ReplaceAll(name, string(filepath.Separator), "_")
}
