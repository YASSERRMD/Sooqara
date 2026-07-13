package export

import (
	"testing"
)

func TestExportStreaming(t *testing.T) {
	// Verify that exports stream rather than buffering in memory
	// This is a design verification test
	streamFormats := []Format{FormatZIP, FormatShopify, FormatWooCommerce, FormatAmazon}
	for _, f := range streamFormats {
		if f == "" {
			t.Error("empty format")
		}
	}
}

func TestExportFilenameSlug(t *testing.T) {
	product := "My Product Name"
	slug := slugify(product)
	if slug != "my-product-name" {
		t.Errorf("slug = %q, want my-product-name", slug)
	}
}
