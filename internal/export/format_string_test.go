package export

import (
	"testing"
)

func TestExportFormatString(t *testing.T) {
	tests := []struct {
		format Format
		want   string
	}{
		{FormatZIP, "zip"},
		{FormatShopify, "shopify"},
		{FormatWooCommerce, "woocommerce"},
		{FormatAmazon, "amazon"},
		{FormatJSON, "json"},
	}
	for _, tt := range tests {
		if string(tt.format) != tt.want {
			t.Errorf("Format(%s) = %q, want %q", tt.format, tt.format, tt.want)
		}
	}
}
