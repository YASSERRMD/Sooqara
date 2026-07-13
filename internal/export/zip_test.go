package export

import (
	"archive/zip"
	"bytes"
	"io"
	"testing"
)

func TestZipStream(t *testing.T) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	f, _ := zw.Create("test.txt")
	f.Write([]byte("hello"))
	zw.Close()
	if buf.Len() == 0 {
		t.Error("expected non-empty ZIP")
	}
}

func TestZipRead(t *testing.T) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	f, _ := zw.Create("test.txt")
	f.Write([]byte("hello world"))
	zw.Close()

 zr, _ := zip.NewReader(&buf, int64(buf.Len()))
 if len(zr.File) == 0 {
  t.Fatal("expected at least 1 file in ZIP")
 }
}
