package release

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunReleaseStages(t *testing.T) {
	os.Setenv("AGNES_API_KEY", "test-key-12345678901234567890")
	defer os.Unsetenv("AGNES_API_KEY")

	tmp := t.TempDir()
	err := RunRelease(tmp, "v1.0.0", "abc123")
	if err != nil {
		t.Fatalf("RunRelease failed: %v", err)
	}

	manifestPath := filepath.Join(tmp, "dist", "v1.0.0", "manifest.json")
	info, err := os.Stat(manifestPath)
	if err != nil {
		t.Fatalf("manifest not created: %v", err)
	}
	if info.Size() == 0 {
		t.Error("expected non-empty manifest")
	}
}

func TestRunReleaseMissingAPIKey(t *testing.T) {
	os.Unsetenv("AGNES_API_KEY")
	tmp := t.TempDir()
	err := RunRelease(tmp, "v1.0.0", "abc")
	if err == nil {
		t.Fatal("expected error without API key")
	}
}
