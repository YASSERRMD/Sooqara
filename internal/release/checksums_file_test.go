package release

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteChecksumsFile(t *testing.T) {
	tmp := t.TempDir()
	os.WriteFile(filepath.Join(tmp, "a.bin"), []byte("aaa"), 0644)
	os.WriteFile(filepath.Join(tmp, "b.bin"), []byte("bbb"), 0644)

	out := filepath.Join(tmp, "SHA256SUMS")
	if err := WriteChecksumsFile(tmp, out); err != nil {
		t.Fatalf("WriteChecksumsFile: %v", err)
	}

	data, _ := os.ReadFile(out)
	content := string(data)
	if !strings.Contains(content, "a.bin") {
		t.Error("expected a.bin in checksums file")
	}
	if !strings.Contains(content, "b.bin") {
		t.Error("expected b.bin in checksums file")
	}
	if !strings.Contains(content, "  ") {
		t.Error("expected two-space separator")
	}
}

func TestWriteChecksumsFileEmptyDir(t *testing.T) {
	tmp := t.TempDir()
	out := filepath.Join(tmp, "SHA256SUMS")
	if err := WriteChecksumsFile(tmp, out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(out)
	if len(data) != 0 {
		t.Error("expected empty file for empty directory")
	}
}
