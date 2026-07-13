package config

import (
	"os"
	"testing"
)

func TestGetEnvOrDefaultIntInvalidValue(t *testing.T) {
	os.Setenv("SOOQARA_RPM", "not-a-number")
	defer os.Clearenv()

	c := Load()
	if c.RPM != 18 {
		t.Errorf("RPM = %d, want default 18 for invalid value", c.RPM)
	}
}

func TestGetEnvOrDefaultIntEmptyValue(t *testing.T) {
	os.Setenv("SOOQARA_RPM", "")
	defer os.Clearenv()

	c := Load()
	if c.RPM != 18 {
		t.Errorf("RPM = %d, want default 18 for empty value", c.RPM)
	}
}
