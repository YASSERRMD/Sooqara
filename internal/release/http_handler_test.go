package release

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleBuildInfoGET(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/build-info", nil)
	w := httptest.NewRecorder()

	HandleBuildInfo(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if resp["version"] != "dev" {
		t.Errorf("expected version 'dev', got %q", resp["version"])
	}
}

func TestHandleBuildInfoWrongMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/build-info", nil)
	w := httptest.NewRecorder()

	HandleBuildInfo(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}

func TestHandleBuildInfoContentType(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/build-info", nil)
	w := httptest.NewRecorder()

	HandleBuildInfo(w, req)

	ct := w.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", ct)
	}
}
