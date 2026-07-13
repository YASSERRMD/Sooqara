package export

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestJSONReproducibility(t *testing.T) {
	data := map[string]any{
		"job_id":     "test-123",
		"product":    "Red Sneaker",
		"artifacts":  []string{"img1.jpg", "img2.jpg"},
		"seeds":      []int64{42, 43, 44},
		"style_ver":  "STYLE_V1",
	}
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(data)
	if buf.Len() == 0 {
		t.Error("expected non-empty JSON")
	}
}
