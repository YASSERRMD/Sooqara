package release

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

// CheckEnvironment validates that the build environment is suitable for release.
func CheckEnvironment() error {
	goVer := runtime.Version()
	if !strings.HasPrefix(goVer, "go1.") {
		return fmt.Errorf("unsupported Go version: %s", goVer)
	}

	home := os.Getenv("HOME")
	if home == "" {
		return fmt.Errorf("HOME environment variable not set")
	}

	key := os.Getenv("AGNES_API_KEY")
	if key == "" {
		return fmt.Errorf("AGNES_API_KEY not set; required for release builds")
	}

	return nil
}
