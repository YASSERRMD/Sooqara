package api

import (
	"bytes"
	"io"
	"testing"
)

func TestIOReadAllEmpty(t *testing.T) {
	data, err := io.ReadAll(bytes.NewReader(nil))
	if err != nil {
		t.Fatalf("ReadAll failed: %v", err)
	}
	if len(data) != 0 {
		t.Errorf("expected empty data, got %d bytes", len(data))
	}
}

func TestIOReadAllWithData(t *testing.T) {
	data, err := io.ReadAll(bytes.NewReader([]byte("hello")))
	if err != nil {
		t.Fatalf("ReadAll failed: %v", err)
	}
	if string(data) != "hello" {
		t.Errorf("got %q, want hello", string(data))
	}
}
