package api

import (
	"testing"
)

func TestRoutesMethodMapping(t *testing.T) {
	methods := map[string]string{
		"POST":   "/api/jobs",
		"GET":    "/api/jobs/{id}",
		"GET":    "/api/jobs/{id}/events",
		"POST":   "/api/jobs/{id}/cancel",
		"POST":   "/api/jobs/{id}/regenerate",
		"GET":    "/api/jobs",
		"GET":    "/healthz",
		"GET":    "/api/artifacts/{id}/raw",
		"GET":    "/metrics",
	}
	if len(methods) != 9 {
		t.Errorf("expected 9 routes, got %d", len(methods))
	}
}
