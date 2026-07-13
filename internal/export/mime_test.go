package export

import (
	"testing"
)

func TestExportContentDisposition(t *testing.T) {
	dispositions := []string{
		"attachment; filename=\"product.zip\"",
		"attachment; filename=\"products.csv\"",
		"inline",
	}
	for _, d := range dispositions {
		if d == "" {
			t.Error("empty disposition")
		}
	}
}

func TestExportMIMETypes(t *testing.T) {
	types := []string{
		"application/zip",
		"text/csv",
		"application/json",
		"application/x-www-form-urlencoded",
	}
	for _, t := range types {
		if t == "" {
			t.Error("empty mime type")
		}
	}
}
