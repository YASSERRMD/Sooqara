package release

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPrepareReleaseDir(t *testing.T) {
	tmp := t.TempDir()
	dir, err := PrepareReleaseDir(tmp, "v1.0.0")
	if err != nil {
		t.Fatalf("PrepareReleaseDir: %v", err)
	}
	expected := filepath.Join(tmp, "dist", "v1.0.0")
	if dir != expected {
		t.Errorf("got %q, want %q", dir, expected)
	}
	info, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("directory not created: %v", err)
	}
	if !info.IsDir() {
		t.Error("expected directory")
	}
}

func TestCopyArtifact(t *testing.T) {
	tmp := t.TempDir()
	src := filepath.Join(tmp, "src.bin")
	os.WriteFile(src, []byte("binary"), 0644)

	destDir := filepath.Join(tmp, "dest")
	os.MkdirAll(destDir, 0755)

	dest, err := CopyArtifact(src, destDir)
	if err != nil {
		t.Fatalf("CopyArtifact: %v", err)
	}
	if filepath.Base(dest) != "src.bin" {
		t.Errorf("expected basename src.bin, got %s", filepath.Base(dest))
	}
}

func TestCopyArtifactMissing(t *testing.T) {
	tmp := t.TempDir()
	_, err := CopyArtifact(filepath.Join(tmp, "missing"), tmp)
	if err == nil {
		t.Fatal("expected error for missing source")
	}
}

func TestListArtifacts(t *testing.T) {
	tmp := t.TempDir()
	os.WriteFile(filepath.Join(tmp, "a.txt"), []byte("a"), 0644)
	os.WriteFile(filepath.Join(tmp, "b.txt"), []byte("b"), 0644)
	os.MkdirAll(filepath.Join(tmp, "sub"), 0755)

	files, err := ListArtifacts(tmp)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 2 {
		t.Errorf("expected 2 files, got %d", len(files))
	}
}

func TestListArtifactsEmpty(t *testing.T) {
	tmp := t.TempDir()
	files, err := ListArtifacts(tmp)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 0 {
		t.Errorf("expected 0 files, got %d", len(files))
	}
}
