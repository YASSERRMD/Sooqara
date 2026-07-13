package hardening

import (
	"strings"
	"testing"
)

func TestValidateCopyLengthOK(t *testing.T) {
	text := strings.Repeat("a", MaxCopyLength)
	if err := ValidateCopyLength(text); err != nil {
		t.Errorf("expected OK, got: %v", err)
	}
}

func TestValidateCopyLengthExceeded(t *testing.T) {
	text := strings.Repeat("a", MaxCopyLength+1)
	err := ValidateCopyLength(text)
	if err == nil {
		t.Fatal("expected error for oversized copy")
	}
}

func TestValidateDescriptionOK(t *testing.T) {
	text := strings.Repeat("a", MaxDescriptionLen)
	if err := ValidateDescriptionLength(text); err != nil {
		t.Errorf("expected OK, got: %v", err)
	}
}

func TestValidateDescriptionExceeded(t *testing.T) {
	text := strings.Repeat("a", MaxDescriptionLen+1)
	err := ValidateDescriptionLength(text)
	if err == nil {
		t.Fatal("expected error for oversized description")
	}
}
