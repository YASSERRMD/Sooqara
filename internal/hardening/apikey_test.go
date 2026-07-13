package hardening

import (
	"strings"
	"testing"
)

func TestValidateAPIKeyOk(t *testing.T) {
	key := strings.Repeat("a", 32)
	if err := ValidateAPIKey(key); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateAPIKeyEmpty(t *testing.T) {
	err := ValidateAPIKey("")
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestValidateAPIKeyTooShort(t *testing.T) {
	key := strings.Repeat("a", 15)
	err := ValidateAPIKey(key)
	if err == nil {
		t.Fatal("expected error for short key")
	}
}

func TestValidateAPIKeyMinLength(t *testing.T) {
	key := strings.Repeat("a", MinAPIKeyLength)
	if err := ValidateAPIKey(key); err != nil {
		t.Errorf("expected OK for min-length key, got: %v", err)
	}
}

func TestValidateAPIKeyMaxLength(t *testing.T) {
	key := strings.Repeat("a", MaxAPIKeyLength)
	if err := ValidateAPIKey(key); err != nil {
		t.Errorf("expected OK for max-length key, got: %v", err)
	}
}

func TestValidateAPIKeyTooLong(t *testing.T) {
	key := strings.Repeat("a", MaxAPIKeyLength+1)
	err := ValidateAPIKey(key)
	if err == nil {
		t.Fatal("expected error for too-long key")
	}
}

func TestValidateAPIKeyWhitespace(t *testing.T) {
	key := strings.Repeat("a", 32)
	if err := ValidateAPIKey("  " + key + "  "); err != nil {
		t.Errorf("expected OK after trimming, got: %v", err)
	}
}
