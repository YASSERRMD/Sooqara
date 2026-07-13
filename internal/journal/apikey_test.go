package journal

import (
	"strings"
	"testing"
)

func TestNeverLogAPIKeyInHash(t *testing.T) {
	req := `{"api_key":"secret-key-123","prompt":"test"}`
	a := NewActivityFromRequest("agnes.text", req, "ok", 0, "", "")

	if strings.Contains(a.RequestHash, "secret-key") {
		t.Error("request hash contains API key material")
	}
}
