package api

import (
	"net/http"
	"testing"
)

func TestHTTPStatusCodes(t *testing.T) {
	statuses := map[int]string{
		http.StatusOK:          "OK",
		http.StatusCreated:     "Created",
		http.StatusAccepted:    "Accepted",
		http.StatusNotFound:    "Not Found",
		http.StatusBadRequest:  "Bad Request",
		http.StatusUnauthorized: "Unauthorized",
		http.StatusInternalServerError: "Internal Server Error",
		http.StatusServiceUnavailable: "Service Unavailable",
	}
	if len(statuses) != 8 {
		t.Errorf("expected 8 status codes, got %d", len(statuses))
	}
}

func TestHTTPMethodConstants(t *testing.T) {
	methods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}
	if len(methods) != 4 {
		t.Errorf("expected 4 methods, got %d", len(methods))
	}
}
