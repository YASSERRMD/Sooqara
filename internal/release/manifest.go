package release

import (
	"encoding/json"
	"fmt"
	"os"
)

// WriteManifest serializes PackageInfo to a JSON manifest file.
func WriteManifest(path string, pkg PackageInfo) error {
	data, err := json.MarshalIndent(pkg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal manifest: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write manifest: %w", err)
	}
	return nil
}

// ReadManifest loads a PackageInfo from a JSON manifest file.
func ReadManifest(path string) (PackageInfo, error) {
	var pkg PackageInfo
	data, err := os.ReadFile(path)
	if err != nil {
		return pkg, fmt.Errorf("read manifest: %w", err)
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return pkg, fmt.Errorf("unmarshal manifest: %w", err)
	}
	return pkg, nil
}
