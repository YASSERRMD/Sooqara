package store

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// FilesystemBlob stores files on disk using content-addressed keys.
type FilesystemBlob struct {
	storageDir string
}

// NewFilesystemBlob creates a new filesystem blob store.
func NewFilesystemBlob(storageDir string) *FilesystemBlob {
	return &FilesystemBlob{storageDir: storageDir}
}

// Put writes content to the storage directory with a content-addressed key.
func (fb *FilesystemBlob) Put(_ context.Context, originalName string, r io.Reader) (string, error) {
	h := sha256.New()
	tee := io.TeeReader(r, h)
	data, err := io.ReadAll(tee)
	if err != nil {
		return "", fmt.Errorf("read content: %w", err)
	}

	hashHex := fmt.Sprintf("%x", h.Sum(nil)[:8])
	dir := filepath.Join(fb.storageDir, hashHex)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("create dir: %w", err)
	}

	key := filepath.Join(dir, originalName)
	if _, err := os.Stat(key); err == nil {
		return key, nil // deduplicated
	}

	if err := os.WriteFile(key, data, 0644); err != nil {
		return "", fmt.Errorf("write file: %w", err)
	}

	return key, nil
}

// Get reads a file from the storage directory.
func (fb *FilesystemBlob) Get(_ context.Context, key string) (io.ReadCloser, error) {
	return os.Open(key)
}

// URL returns a URL for the given key.
func (fb *FilesystemBlob) URL(key string) string {
	return "/api/artifacts/" + filepath.Base(key)
}
