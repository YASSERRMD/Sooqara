package release

import (
	"os"
	"strings"
	"testing"
)

func TestGenerateBuildFlagsContainsVersion(t *testing.T) {
	flags := GenerateBuildFlags("v1.0.0", "abc123")
	if !strings.Contains(flags, "v1.0.0") {
		t.Error("expected flags to contain version")
	}
}

func TestGenerateBuildFlagsContainsCommit(t *testing.T) {
	flags := GenerateBuildFlags("v1.0.0", "abc123")
	if !strings.Contains(flags, "abc123") {
		t.Error("expected flags to contain commit")
	}
}

func TestGenerateBuildFlagsContainsBuildTime(t *testing.T) {
	flags := GenerateBuildFlags("v1.0.0", "abc123")
	if !strings.Contains(flags, "BuildTime") {
		t.Error("expected flags to contain BuildTime")
	}
}

func TestGenerateBuildFlagsContainsPackagePath(t *testing.T) {
	flags := GenerateBuildFlags("v1.0.0", "abc123")
	if !strings.Contains(flags, "github.com/yasserrmd/sooqara/internal/version") {
		t.Error("expected flags to contain package path")
	}
}

func TestDefaultFlagsUsesDev(t *testing.T) {
	os.Unsetenv("VERSION")
	os.Unsetenv("COMMIT")
	flags := DefaultFlags()
	if !strings.Contains(flags, "dev") {
		t.Error("expected DefaultFlags to contain 'dev'")
	}
}

func TestDefaultFlagsUsesEnv(t *testing.T) {
	os.Setenv("VERSION", "v2.0.0")
	os.Setenv("COMMIT", "deadbeef")
	defer func() {
		os.Unsetenv("VERSION")
		os.Unsetenv("COMMIT")
	}()
	flags := DefaultFlags()
	if !strings.Contains(flags, "v2.0.0") {
		t.Error("expected DefaultFlags to use VERSION env")
	}
	if !strings.Contains(flags, "deadbeef") {
		t.Error("expected DefaultFlags to use COMMIT env")
	}
}
