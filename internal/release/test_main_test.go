package release

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Setenv("AGNES_API_KEY", "test-key-12345678901234567890")
	code := m.Run()
	os.Unsetenv("AGNES_API_KEY")
	os.Exit(code)
}
