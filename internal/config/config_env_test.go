package config

import (
	"os"
	"testing"
)

func TestGetEnvOrDefaultReturnsValue(t *testing.T) {
	os.Setenv("TEST_VAR", "test-value")
	defer os.Clearenv()

	got := getEnvOrDefault("TEST_VAR", "default")
	if got != "test-value" {
		t.Errorf("getEnvOrDefault = %s, want test-value", got)
	}
}

func TestGetEnvOrDefaultReturnsDefault(t *testing.T) {
	os.Clearenv()
	got := getEnvOrDefault("NONEXISTENT_VAR", "fallback")
	if got != "fallback" {
		t.Errorf("getEnvOrDefault = %s, want fallback", got)
	}
}
