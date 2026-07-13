package release

import (
	"fmt"
	"os"

	"github.com/yasserrmd/sooqara/internal/version"
)

// PrintVersion prints version info and exits with code 0.
func PrintVersion() {
	fmt.Println(version.BuildInfo())
	os.Exit(0)
}

// PrintVersionVerbose prints full metadata including build info.
func PrintVersionVerbose() {
	meta := Metadata{
		Version:   version.Version,
		Commit:    version.Commit,
		BuildTime: version.BuildTime,
	}
	fmt.Println(meta.String())
}
