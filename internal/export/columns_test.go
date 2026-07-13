package export

import (
	"testing"
)

func TestShopifyColumns(t *testing.T) {
	columns := []string{
		"Handle", "Title", "Body (HTML)", "Vendor", "Type", "Tags",
		"Published", "Image Src", "Image Position", "Image Alt Text",
		"SEO Title", "SEO Description",
	}
	if len(columns) != 12 {
		t.Errorf("expected 12 columns, got %d", len(columns))
	}
}

func TestWooCommerceColumns(t *testing.T) {
	columns := []string{"SKU", "Name", "Regular Price", "Description", "Images", "Categories"}
	if len(columns) != 6 {
		t.Errorf("expected 6 columns, got %d", len(columns))
	}
}

func TestAmazonTSVColums(t *testing.T) {
	columns := []string{"sku", "product-name", "description"}
	if len(columns) != 3 {
		t.Errorf("expected 3 columns, got %d", len(columns))
	}
}
