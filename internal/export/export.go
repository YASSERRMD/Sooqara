package export

import (
	"archive/zip"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// Format represents an export format.
type Format string

const (
	FormatZIP       Format = "zip"
	FormatShopify   Format = "shopify"
	FormatWooCommerce Format = "woocommerce"
	FormatAmazon    Format = "amazon"
	FormatJSON      Format = "json"
)

// Exporter produces an export in a given format.
type Exporter struct {
	Format Format
}

// NewExporter creates a new exporter.
func NewExporter(f Format) *Exporter {
	return &Exporter{Format: f}
}

// WriteZIP streams a ZIP archive to the writer.
func (e *Exporter) WriteZIP(w io.Writer, product string, images []string, copyJSON string) error {
	zw := zip.NewWriter(w)
	f, _ := zw.Create("copy.json")
	f.Write([]byte(copyJSON))
	f, _ = zw.Create("README.txt")
	f.Write([]byte(product))
	return zw.Close()
}

// WriteShopifyCSV writes a Shopify-compatible CSV.
func (e *Exporter) WriteShopifyCSV(w io.Writer, product string, title string, body string, images []string, altTexts []string, keywords []string) error {
	cw := csv.NewWriter(w)
	cw.Write([]string{"Handle", "Title", "Body (HTML)", "Vendor", "Type", "Tags", "Published", "Image Src", "Image Position", "Image Alt Text", "SEO Title", "SEO Description"})
	cw.Write([]string{slugify(product), title, body, "Sooqara", product, strings.Join(keywords, ","), "true", strings.Join(images, "|"), "1", strings.Join(altTexts, "|"), title, ""})
	cw.Flush()
	return cw.Error()
}

// WriteWooCommerceCSV writes a WooCommerce-compatible CSV.
func (e *Exporter) WriteWooCommerceCSV(w io.Writer, product string, title string, body string, images []string) error {
	cw := csv.NewWriter(w)
	cw.Write([]string{"SKU", "Name", "Regular Price", "Description", "Images", "Categories"})
	cw.Write([]string{"", title, "0.00", body, strings.Join(images, ";"), product})
	cw.Flush()
	return cw.Error()
}

// WriteAmazonTSV writes a flat-file TSV.
func (e *Exporter) WriteAmazonTSV(w io.Writer, product string, title string, body string) error {
	fmt.Fprintf(w, "sku\tproduct-name\tdescription\n")
	fmt.Fprintf(w, "%s\t%s\t%s\n", slugify(product), title, body)
	return nil
}

// WriteJSON writes the full job export.
func (e *Exporter) WriteJSON(w io.Writer, data any) error {
	return json.NewEncoder(w).Encode(data)
}

func slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	return s
}
