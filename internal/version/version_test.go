package version

import (
	"strings"
	"testing"
)

func TestInfoReturnsDev(t *testing.T) {
	info := Info()
	if info != "dev" {
		t.Errorf("expected 'dev', got %q", info)
	}
}

func TestBuildInfoContainsVersion(t *testing.T) {
	info := BuildInfo()
	if !strings.Contains(info, "sooqara") {
		t.Error("expected BuildInfo to contain 'sooqara'")
	}
}

func TestBuildInfoContainsCommit(t *testing.T) {
	info := BuildInfo()
	if !strings.Contains(info, "commit") {
		t.Error("expected BuildInfo to contain 'commit'")
	}
}

func TestBuildInfoContainsBuilt(t *testing.T) {
	info := BuildInfo()
	if !strings.Contains(info, "built") {
		t.Error("expected BuildInfo to contain 'built'")
	}
}

func TestCustomVersion(t *testing.T) {
	orig := Version
	Version = "v1.0.0"
	defer func() { Version = orig }()
	if Info() != "v1.0.0" {
		t.Error("expected custom version")
	}
}
