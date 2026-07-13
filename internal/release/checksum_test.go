package release

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSHA256FromBytesKnown(t *testing.T) {
	data := []byte("hello")
	hash := SHA256FromBytes(data)
	// Known SHA-256 of "hello"
	expected := "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
	if hash != expected {
		t.Errorf("expected %s, got %s", expected, hash)
	}
}

func TestSHA256FromBytesEmpty(t *testing.T) {
	hash := SHA256FromBytes([]byte{})
	expected := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	if hash != expected {
		t.Errorf("expected empty SHA256, got %s", hash)
	}
}

func TestSHA256ChecksumFile(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "test.txt")
	os.WriteFile(path, []byte("checksum test"), 0644)

	hash, err := SHA256Checksum(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(hash) != 64 {
		t.Errorf("expected 64-char hex digest, got %d chars", len(hash))
	}
	if !strings.Contains(hash, "a") {
		// Just verify it's a non-trivial hash
		t.Log("hash:", hash)
	}
}

func TestSHA256ChecksumMissingFile(t *testing.T) {
	tmp := t.TempDir()
	_, err := SHA256Checksum(filepath.Join(tmp, "missing.txt"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
