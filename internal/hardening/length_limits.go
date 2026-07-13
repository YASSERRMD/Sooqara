package hardening

import "fmt"

const (
	MaxCopyLength     = 5000
	MaxDescriptionLen = 10000
	MaxTagName        = 50
	MaxTagsPerListing = 20
)

// ValidateCopyLength checks that generated copy text stays within limits.
func ValidateCopyLength(text string) error {
	if len(text) > MaxCopyLength {
		return fmt.Errorf("copy text exceeds %d characters (got %d)", MaxCopyLength, len(text))
	}
	return nil
}

// ValidateDescriptionLength checks that a product description stays within limits.
func ValidateDescriptionLength(text string) error {
	if len(text) > MaxDescriptionLen {
		return fmt.Errorf("description exceeds %d characters (got %d)", MaxDescriptionLen, len(text))
	}
	return nil
}
