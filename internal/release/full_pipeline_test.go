package release

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFullReleasePipeline(t *testing.T) {
	os.Setenv("AGNES_API_KEY", "test-key-12345678901234567890")
	defer os.Unsetenv("AGNES_API_KEY")

	tmp := t.TempDir()

	// Stage 1: Prepare
	dir, err := PrepareReleaseDir(tmp, "v1.0.0")
	if err != nil {
		t.Fatalf("PrepareReleaseDir: %v", err)
	}

	// Stage 2: Write manifest
	pkg := PackageInfo{
		Name:        "sooqara",
		Version:     "v1.0.0",
		Platform:    "linux",
		Arch:        "amd64",
		Description: "full pipeline test",
	}
	if err := WriteManifest(filepath.Join(dir, "manifest.json"), pkg); err != nil {
		t.Fatalf("WriteManifest: %v", err)
	}

	// Stage 3: List
	files, err := ListArtifacts(dir)
	if err != nil {
		t.Fatalf("ListArtifacts: %v", err)
	}
	if len(files) != 1 {
		t.Errorf("expected 1 artifact, got %d", len(files))
	}

	// Stage 4: Checksums
	if err := WriteChecksumsFile(dir, filepath.Join(dir, "SHA256SUMS")); err != nil {
		t.Fatalf("WriteChecksumsFile: %v", err)
	}

	// Stage 5: Release note
	note := GenerateReleaseNote(pkg)
	if !strings.Contains(note, "v1.0.0") {
		t.Error("expected version in release note")
	}

	// Stage 6: Build flags
	flags := GenerateBuildFlags("v1.0.0", "test")
	if !strings.Contains(flags, "v1.0.0") {
		t.Error("expected version in build flags")
	}

	// Stage 7: Environment
	if err := CheckEnvironment(); err != nil {
		t.Errorf("CheckEnvironment: %v", err)
	}
}
