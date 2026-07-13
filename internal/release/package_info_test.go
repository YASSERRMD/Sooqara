package release

import (
	"testing"
	"time"
)

func TestPackageInfoFields(t *testing.T) {
	p := PackageInfo{
		Name:      "sooqara",
		Version:   "v1.0.0",
		Platform:  "linux",
		Arch:      "amd64",
		BuildTime: time.Now(),
		Description: "E-commerce listing factory",
	}
	if p.Name != "sooqara" {
		t.Error("expected Name sooqara")
	}
	if p.Version != "v1.0.0" {
		t.Error("expected Version v1.0.0")
	}
	if p.Platform != "linux" {
		t.Error("expected Platform linux")
	}
	if p.Arch != "amd64" {
		t.Error("expected Arch amd64")
	}
	if p.Description != "E-commerce listing factory" {
		t.Error("expected Description")
	}
}

func TestPackageInfoEmpty(t *testing.T) {
	p := PackageInfo{}
	if p.Name != "" {
		t.Error("expected empty Name")
	}
}

func TestPackageInfoArtifactsNil(t *testing.T) {
	p := PackageInfo{}
	if p.Artifacts != nil && len(p.Artifacts) != 0 {
		t.Error("expected nil or empty Artifacts")
	}
}

func TestPackageInfoChecksumsNil(t *testing.T) {
	p := PackageInfo{}
	if p.Checksums != nil && len(p.Checksums) != 0 {
		t.Error("expected nil or empty Checksums")
	}
}
