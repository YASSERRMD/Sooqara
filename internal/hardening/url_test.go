package hardening

import "testing"

func TestValidateURLOk(t *testing.T) {
	tests := []string{
		"https://example.com/path?query=1",
		"http://localhost:8080/api",
		"https://api.agnes.ai/v1/chat",
	}
	for _, raw := range tests {
		if err := ValidateURL(raw); err != nil {
			t.Errorf("ValidateURL(%q) = %v, want nil", raw, err)
		}
	}
}

func TestValidateURLBadScheme(t *testing.T) {
	err := ValidateURL("ftp://example.com")
	if err == nil {
		t.Fatal("expected error for ftp scheme")
	}
}

func TestValidateURLNoHost(t *testing.T) {
	err := ValidateURL("https:///path")
	if err == nil {
		t.Fatal("expected error for missing host")
	}
}

func TestValidateURLMalformed(t *testing.T) {
	err := ValidateURL("not a url")
	if err == nil {
		t.Fatal("expected error for malformed URL")
	}
}

func TestSanitizeURLTrims(t *testing.T) {
	got := SanitizeURL("  https://example.com  ")
	if got != "https://example.com" {
		t.Errorf("expected 'https://example.com', got %q", got)
	}
}

func TestValidateURLTooLong(t *testing.T) {
	long := "https://example.com/"
	for len(long) < MaxURLLength+10 {
		long += "x"
	}
	err := ValidateURL(long)
	if err == nil {
		t.Fatal("expected error for too-long URL")
	}
}
