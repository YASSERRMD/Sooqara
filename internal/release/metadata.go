package release

import "fmt"

// Metadata holds release artifact information.
type Metadata struct {
	Version   string
	Commit    string
	BuildTime string
	GoVersion string
	Platform  string
}

// String returns a formatted metadata block.
func (m Metadata) String() string {
	return fmt.Sprintf(`Release Metadata
================
Version:  %s
Commit:   %s
Built:    %s
Go:       %s
Platform: %s`, m.Version, m.Commit, m.BuildTime, m.GoVersion, m.Platform)
}
