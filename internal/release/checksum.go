package release

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

// SHA256Checksum computes the SHA-256 hex digest of a file.
func SHA256Checksum(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// SHA256FromBytes computes the SHA-256 hex digest of in-memory bytes.
func SHA256FromBytes(data []byte) string {
	h := sha256.Sum256(data)
	return fmt.Sprintf("%x", h[:])
}
