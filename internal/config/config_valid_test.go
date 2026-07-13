package config

import (
	"os"
	"testing"
)

func TestValidateValidConfig(t *testing.T) {
	os.Setenv("AGNES_API_KEY", "test-key")
	os.Setenv("SOOQARA_LOG_LEVEL", "debug")
	defer os.Clearenv()

	c := Load()
	if err := c.Validate(); err != nil {
		t.Fatalf("unexpected error for valid config: %v", err)
	}
}
