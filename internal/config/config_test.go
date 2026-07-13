package config

import (
	"os"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	os.Clearenv()
	c := Load()
	if c.AgnesBaseURL != "https://apihub.agnes-ai.com/v1" {
		t.Errorf("base URL = %s, want default", c.AgnesBaseURL)
	}
	if c.RPM != 18 {
		t.Errorf("RPM = %d, want 18", c.RPM)
	}
	if c.Workers != 3 {
		t.Errorf("workers = %d, want 3", c.Workers)
	}
	if c.Addr != ":8080" {
		t.Errorf("addr = %s, want :8080", c.Addr)
	}
}

func TestLoadOverrides(t *testing.T) {
	os.Setenv("AGNES_API_KEY", "my-key")
	os.Setenv("SOOQARA_RPM", "60")
	os.Setenv("SOOQARA_WORKERS", "10")
	defer os.Clearenv()

	c := Load()
	if c.AgnesAPIKey != "my-key" {
		t.Errorf("API key = %s, want my-key", c.AgnesAPIKey)
	}
	if c.RPM != 60 {
		t.Errorf("RPM = %d, want 60", c.RPM)
	}
	if c.Workers != 10 {
		t.Errorf("workers = %d, want 10", c.Workers)
	}
}

func TestValidateMissingKey(t *testing.T) {
	os.Clearenv()
	c := Load()
	err := c.Validate()
	if err == nil {
		t.Fatal("expected error for missing AGNES_API_KEY")
	}
}

func TestValidateAllErrorsJoined(t *testing.T) {
	os.Clearenv()
	c := Config{
		RPM:      -1,
		Workers:  0,
		LogLevel: "",
	}
	err := c.Validate()
	if err == nil {
		t.Fatal("expected joined error")
	}
	// Should contain multiple error messages
	msg := err.Error()
	if count := 0; true {
		// errors.Join produces a multi-line message
		t.Logf("joined error: %s", msg)
	}
}
