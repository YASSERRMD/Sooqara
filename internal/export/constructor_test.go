package export

import (
	"testing"
)

func TestExporterCreation(t *testing.T) {
	formats := []Format{FormatZIP, FormatShopify, FormatWooCommerce, FormatAmazon, FormatJSON}
	for _, f := range formats {
		e := NewExporter(f)
		if e.Format != f {
			t.Errorf("Exporter.Format = %s, want %s", e.Format, f)
		}
	}
}
