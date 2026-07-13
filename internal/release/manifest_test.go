package release

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestWriteAndReadManifest(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "manifest.json")

	pkg := PackageInfo{
		Name:      "sooqara",
		Version:   "v1.0.0",
		Platform:  "darwin",
		Arch:      "arm64",
		BuildTime: time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
		Description: "test",
	}

	if err := WriteManifest(path, pkg); err != nil {
		t.Fatalf("WriteManifest: %v", err)
	}

	loaded, err := ReadManifest(path)
	if err != nil {
		t.Fatalf("ReadManifest: %v", err)
	}

	if loaded.Name != pkg.Name {
		t.Errorf("Name: got %q, want %q", loaded.Name, pkg.Name)
	}
	if loaded.Version != pkg.Version {
		t.Errorf("Version: got %q, want %q", loaded.Version, pkg.Version)
	}
	if loaded.Platform != pkg.Platform {
		t.Errorf("Platform: got %q, want %q", loaded.Platform, pkg.Platform)
	}
	if loaded.Arch != pkg.Arch {
		t.Errorf("Arch: got %q, want %q", loaded.Arch, pkg.Arch)
	}
}

func TestReadManifestMissing(t *testing.T) {
	tmp := t.TempDir()
	_, err := ReadManifest(filepath.Join(tmp, "missing.json"))
	if err == nil {
		t.Fatal("expected error for missing manifest")
	}
}

func TestWriteManifestInvalidPath(t *testing.T) {
	err := WriteManifest("/nonexistent/dir/manifest.json", PackageInfo{})
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}

func TestManifestFileExists(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "test.json")
	pkg := PackageInfo{Name: "test"}
	WriteManifest(path, pkg)

	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Size() == 0 {
		t.Error("expected non-empty manifest file")
	}
}
