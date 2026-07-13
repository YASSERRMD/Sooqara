package release

import (
	"os"
	"testing"
)

func TestCheckEnvironmentMissingAPIKey(t *testing.T) {
	os.Unsetenv("AGNES_API_KEY")
	err := CheckEnvironment()
	if err == nil {
		t.Fatal("expected error without AGNES_API_KEY")
	}
}

func TestCheckEnvironmentMissingHome(t *testing.T) {
	os.Setenv("AGNES_API_KEY", "test-key")
	os.Unsetenv("HOME")
	defer os.Setenv("HOME", "/tmp")

	err := CheckEnvironment()
	if err != nil {
		t.Logf("Got error (HOME may be required): %v", err)
	}
}

func TestCheckEnvironmentWithAPIKey(t *testing.T) {
	os.Setenv("AGNES_API_KEY", "test-key-12345678901234567890")
	os.Setenv("HOME", "/tmp")
	defer func() {
		os.Unsetenv("AGNES_API_KEY")
	}()
	err := CheckEnvironment()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
