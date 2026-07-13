package export

import (
	"bytes"
	"encoding/csv"
	"testing"
)

func TestCSVQuoting(t *testing.T) {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	w.Write([]string{"field with, comma", `field with "quotes"`})
	w.Flush()
	if buf.Len() == 0 {
		t.Error("expected non-empty CSV")
	}
}

func TestCSVNewlines(t *testing.T) {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	w.Write([]string{"line1\nline2"})
	w.Flush()
	if buf.Len() == 0 {
		t.Error("expected non-empty CSV")
	}
}
