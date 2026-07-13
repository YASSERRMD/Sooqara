package pipeline

import (
	"fmt"
	"regexp"
)

var hexColorRe = regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`)

// validateAnalysis checks the analysis results for correctness.
func validateAnalysis(a *ProductAnalysis, variantCount int) error {
	if a.ProductName == "" {
		return fmt.Errorf("product_name is empty")
	}
	if a.Category == "" {
		return fmt.Errorf("category is empty")
	}
	if a.ShapeDescription == "" {
		return fmt.Errorf("shape_description is empty")
	}

	// Validate dominant colour hex
	if a.DominantColourHex != nil && !hexColorRe.MatchString(*a.DominantColourHex) {
		a.DominantColourHex = nil
	}

	return nil
}
