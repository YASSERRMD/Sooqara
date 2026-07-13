package hardening

import (
	"fmt"
	"strings"
)

// ValidateTags checks tag count and individual tag lengths.
func ValidateTags(tags []string) error {
	if len(tags) > MaxTagsPerListing {
		return fmt.Errorf("too many tags: maximum %d, got %d", MaxTagsPerListing, len(tags))
	}
	for i, tag := range tags {
		cleaned := strings.TrimSpace(tag)
		if len(cleaned) > MaxTagName {
			return fmt.Errorf("tag %d exceeds %d characters", i+1, MaxTagName)
		}
		if cleaned == "" {
			return fmt.Errorf("tag %d must not be empty", i+1)
		}
	}
	return nil
}
