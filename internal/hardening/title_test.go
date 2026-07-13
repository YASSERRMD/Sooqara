package hardening

import "testing"

func TestValidateTitleOK(t *testing.T) {
	if err := ValidateTitle("My Listing"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateTitleEmpty(t *testing.T) {
	err := ValidateTitle("")
	if err == nil {
		t.Fatal("expected error for empty title")
	}
}

func TestValidateTitleWhitespaceOnly(t *testing.T) {
	err := ValidateTitle("   ")
	if err == nil {
		t.Fatal("expected error for whitespace-only title")
	}
}

func TestValidateTitleTooLong(t *testing.T) {
	long := ""
	for i := 0; i < MaxTitleLength+1; i++ {
		long += "x"
	}
	err := ValidateTitle(long)
	if err == nil {
		t.Fatal("expected error for overly long title")
	}
}

func TestValidateTitleMaxLength(t *testing.T) {
	long := ""
	for i := 0; i < MaxTitleLength; i++ {
		long += "x"
	}
	if err := ValidateTitle(long); err != nil {
		t.Errorf("expected OK for max-length title, got: %v", err)
	}
}

func TestValidateTitleWithControlChars(t *testing.T) {
	err := ValidateTitle("Hello\x00World")
	if err != nil {
		t.Errorf("expected OK after sanitization, got: %v", err)
	}
}
