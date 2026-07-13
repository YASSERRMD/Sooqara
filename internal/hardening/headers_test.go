package hardening

import (
	"strings"
	"testing"
)

func TestValidateHeaderNameAllowed(t *testing.T) {
	names := []string{"content-type", "Accept", "X-REQUEST-ID", "x-correlation-id"}
	for _, name := range names {
		if err := ValidateHeaderName(name); err != nil {
			t.Errorf("ValidateHeaderName(%q) = %v, want nil", name, err)
		}
	}
}

func TestValidateHeaderNameDisallowed(t *testing.T) {
	err := ValidateHeaderName("cookie")
	if err == nil {
		t.Fatal("expected error for cookie header")
	}
}

func TestValidateHeaderNameEmpty(t *testing.T) {
	err := ValidateHeaderName("")
	if err == nil {
		t.Fatal("expected error for empty header name")
	}
}

func TestHeaderNotAllowedError(t *testing.T) {
	err := &HeaderNotAllowedError{name: "x-custom"}
	msg := err.Error()
	if !strings.Contains(msg, "x-custom") {
		t.Errorf("expected error message to contain 'x-custom', got %q", msg)
	}
}

func TestValidateHeaderNameWhitespace(t *testing.T) {
	err := ValidateHeaderName("  content-type  ")
	if err != nil {
		t.Errorf("expected OK after trim, got: %v", err)
	}
}
