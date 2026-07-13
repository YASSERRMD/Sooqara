package export

import (
	"bytes"
	"testing"
)

func TestWriteZIP(t *testing.T) {
	e := NewExporter(FormatZIP)
	var buf bytes.Buffer
	err := e.WriteZIP(&buf, "Test Product", []string{}, `{}`)
	if err != nil {
		t.Fatalf("WriteZIP failed: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty ZIP")
	}
}

func TestWriteShopifyCSV(t *testing.T) {
	e := NewExporter(FormatShopify)
	var buf bytes.Buffer
	err := e.WriteShopifyCSV(&buf, "Test", "Title", "Body", []string{"http://img"}, []string{"alt"}, []string{"kw1", "kw2"})
	if err != nil {
		t.Fatalf("WriteShopifyCSV failed: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty CSV")
	}
}

func TestWriteWooCommerceCSV(t *testing.T) {
	e := NewExporter(FormatWooCommerce)
	var buf bytes.Buffer
	err := e.WriteWooCommerceCSV(&buf, "Test", "Title", "Body", []string{"http://img"})
	if err != nil {
		t.Fatalf("WriteWooCommerceCSV failed: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty CSV")
	}
}

func TestWriteAmazonTSV(t *testing.T) {
	e := NewExporter(FormatAmazon)
	var buf bytes.Buffer
	err := e.WriteAmazonTSV(&buf, "Test", "Title", "Body")
	if err != nil {
		t.Fatalf("WriteAmazonTSV failed: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty TSV")
	}
}

func TestWriteJSON(t *testing.T) {
	e := NewExporter(FormatJSON)
	var buf bytes.Buffer
	data := map[string]any{"key": "value"}
	err := e.WriteJSON(&buf, data)
	if err != nil {
		t.Fatalf("WriteJSON failed: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty JSON")
	}
}

func TestSlugify(t *testing.T) {
	tests := []struct{ in, want string }{
		{"Hello World", "hello-world"},
		{"UPPER CASE", "upper case"},
		{"Mixed Case", "mixed case"},
	}
	for _, tt := range tests {
		got := slugify(tt.in)
		if got != tt.want {
			t.Errorf("slugify(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestExportFormats(t *testing.T) {
	formats := []Format{FormatZIP, FormatShopify, FormatWooCommerce, FormatAmazon, FormatJSON}
	for _, f := range formats {
		if f == "" {
			t.Error("empty format")
		}
	}
	if len(formats) != 5 {
		t.Errorf("expected 5 formats, got %d", len(formats))
	}
}
