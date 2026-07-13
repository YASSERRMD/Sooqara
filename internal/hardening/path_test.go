package hardening

import "testing"

func TestSafePathAllowed(t *testing.T) {
	if !SafePath("/data", "/data/files/img.png") {
		t.Error("expected /data/files/img.png to be safe")
	}
}

func TestSafePathDotDot(t *testing.T) {
	if SafePath("/data", "/data/../etc/passwd") {
		t.Error("expected /data/../etc/passwd to be unsafe")
	}
}

func TestSafePathEqualBase(t *testing.T) {
	if !SafePath("/data", "/data") {
		t.Error("expected base path itself to be safe")
	}
}

func TestSafePathSubdir(t *testing.T) {
	if !SafePath("/data", "/data/sub/deep/file.txt") {
		t.Error("expected deep subpath to be safe")
	}
}

func TestSafePathCleanEscapes(t *testing.T) {
	if SafePath("/data", "/data/sub/../../etc/shadow") {
		t.Error("expected cleaned escape attempt to be rejected")
	}
}
