package release

import (
	"strings"
	"testing"
	"time"
)

func TestGenerateReleaseNoteBasic(t *testing.T) {
	pkg := PackageInfo{
		Version:     "v1.0.0",
		Platform:    "linux",
		Arch:        "amd64",
		Description: "test release",
		BuildTime:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	note := GenerateReleaseNote(pkg)
	if !strings.Contains(note, "v1.0.0") {
		t.Error("expected version in note")
	}
	if !strings.Contains(note, "linux/amd64") {
		t.Error("expected platform in note")
	}
	if !strings.Contains(note, "test release") {
		t.Error("expected description in note")
	}
}

func TestGenerateReleaseNoteWithArtifacts(t *testing.T) {
	pkg := PackageInfo{
		Version:     "v2.0.0",
		Platform:    "darwin",
		Arch:        "arm64",
		Description: "with artifacts",
		Artifacts:   []string{"sooqara.tar.gz", "checksums.sha256"},
		Checksums:   map[string]string{"sooqara.tar.gz": "abc123"},
		BuildTime:   time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
	}
	note := GenerateReleaseNote(pkg)
	if !strings.Contains(note, "sooqara.tar.gz") {
		t.Error("expected artifact name in note")
	}
	if !strings.Contains(note, "abc123") {
		t.Error("expected checksum in note")
	}
	if !strings.Contains(note, "## Artifacts") {
		t.Error("expected artifacts section")
	}
}

func TestGenerateReleaseNoteEmpty(t *testing.T) {
	pkg := PackageInfo{}
	note := GenerateReleaseNote(pkg)
	if !strings.Contains(note, "# Release") {
		t.Error("expected header even for empty package")
	}
}
