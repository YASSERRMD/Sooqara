package api

import (
	"encoding/json"
	"testing"
)

func TestJSONResponseFormat(t *testing.T) {
	data, err := json.Marshal(map[string]string{"job_id": "test-123"})
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty JSON")
	}
}

func TestJSONDecodeEmptyBody(t *testing.T) {
	var m map[string]string
	err := json.Unmarshal([]byte{}, &m)
	if err == nil {
		t.Error("expected error for empty body")
	}
}
