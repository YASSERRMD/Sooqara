package web

import (
	"testing"
)

func TestEmbedFS(t *testing.T) {
	entries, err := fs.ReadDir("templates")
	if err != nil {
		t.Logf("ReadDir returned: %v", err)
	}
	_ = entries
}
