package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthzReturns200(t *testing.T) {
	h := &Handler{}
	req := httptest.NewRequest("GET", "/healthz", nil)
	w := httptest.NewRecorder()
	h.Healthz(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", w.Code)
	}
	if w.Body.Len() == 0 {
		t.Error("expected non-empty body")
	}
}

func TestMetricsReturnsPlainText(t *testing.T) {
	h := &Handler{}
	req := httptest.NewRequest("GET", "/metrics", nil)
	w := httptest.NewRecorder()
	h.Metrics(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", w.Code)
	}
	ct := w.Header().Get("Content-Type")
	if ct != "text/plain; version=0.0.4" {
		t.Errorf("Content-Type = %q, want text/plain; version=0.0.4", ct)
	}
}

func TestRawArtifactSetsHeaders(t *testing.T) {
	h := &Handler{}
	req := httptest.NewRequest("GET", "/api/artifacts/test123/raw", nil)
	w := httptest.NewRecorder()
	h.RawArtifact(w, req)
	cd := w.Header().Get("Content-Disposition")
	if cd != "inline" {
		t.Errorf("Content-Disposition = %q, want inline", cd)
	}
}

func TestErrorResponseStructure(t *testing.T) {
	errResp := ErrorResponse{
		Code:      "bad_request",
		Message:   "invalid input",
		RequestID: "req-123",
	}
	data, _ := json.Marshal(errResp)
	if len(data) == 0 {
		t.Error("expected non-empty JSON")
	}
}
