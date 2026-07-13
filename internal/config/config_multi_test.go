package config

import (
	"os"
	"strings"
	"testing"
)

func TestValidateMultipleErrorsJoined(t *testing.T) {
	os.Clearenv()
	c := Config{
		RPM:      -1,
		Workers:  0,
		LogLevel: "",
	}
	err := c.Validate()
	if err == nil {
		t.Fatal("expected joined error for multiple invalid fields")
	}
	msg := err.Error()
	if !strings.Contains(msg, "AGNES_API_KEY") {
		t.Error("missing AGNES_API_KEY in joined error")
	}
	if !strings.Contains(msg, "SOOQARA_RPM") {
		t.Error("missing SOOQARA_RPM in joined error")
	}
	if !strings.Contains(msg, "SOOQARA_WORKERS") {
		t.Error("missing SOOQARA_WORKERS in joined error")
	}
}
