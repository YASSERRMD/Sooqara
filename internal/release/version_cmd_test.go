package release

import (
	"strings"
	"testing"

	"github.com/yasserrmd/sooqara/internal/version"
)

func TestPrintVersionVerboseContainsVersion(t *testing.T) {
	meta := Metadata{
		Version:   version.Version,
		Commit:    version.Commit,
		BuildTime: version.BuildTime,
	}
	out := meta.String()
	if !strings.Contains(out, version.Version) {
		t.Error("expected version in verbose output")
	}
}

func TestPrintVersionVerboseContainsHeader(t *testing.T) {
	meta := Metadata{}
	out := meta.String()
	if !strings.Contains(out, "Release Metadata") {
		t.Error("expected header in metadata output")
	}
}
