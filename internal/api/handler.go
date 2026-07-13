package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yasserrmd/sooqara/internal/limiter"
	"github.com/yasserrmd/sooqara/internal/provider"
	"github.com/yasserrmd/sooqara/internal/store"
)

// Handler serves HTTP API endpoints.
type Handler struct {
	store   *store.Store
	prov    provider.Provider
	limiter *limiter.Limiter
}

// NewHandler creates a new API handler.
func NewHandler(s *store.Store, p provider.Provider, lm *limiter.Limiter) *Handler {
	return &Handler{store: s, prov: p, limiter: lm}
}

// Routes registers all API routes on the given mux.
func (h *Handler) Routes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/jobs", h.CreateJob)
	mux.HandleFunc("GET /api/jobs/{id}", h.GetJob)
	mux.HandleFunc("GET /api/jobs/{id}/events", h.JobEvents)
	mux.HandleFunc("POST /api/jobs/{id}/cancel", h.CancelJob)
	mux.HandleFunc("POST /api/jobs/{id}/regenerate", h.Regenerate)
	mux.HandleFunc("GET /api/jobs", h.ListJobs)
	mux.HandleFunc("GET /healthz", h.Healthz)
	mux.HandleFunc("GET /api/artifacts/{id}/raw", h.RawArtifact)
	mux.HandleFunc("GET /metrics", h.Metrics)
}

// CreateJob handles POST /api/jobs.
func (h *Handler) CreateJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"job_id": "new-job-id"})
}

// GetJob handles GET /api/jobs/{id}.
func (h *Handler) GetJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": "job-id"})
}

// JobEvents handles GET /api/jobs/{id}/events (SSE).
func (h *Handler) JobEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
}

// CancelJob handles POST /api/jobs/{id}/cancel.
func (h *Handler) CancelJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "cancelled"})
}

// Regenerate handles POST /api/jobs/{id}/regenerate.
func (h *Handler) Regenerate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "regenerated"})
}

// ListJobs handles GET /api/jobs.
func (h *Handler) ListJobs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]string{})
}

// Healthz handles GET /healthz.
func (h *Handler) Healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status": "ok",
		"time":   time.Now().UTC().Format(time.RFC3339),
	})
}

// RawArtifact handles GET /api/artifacts/{id}/raw.
func (h *Handler) RawArtifact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "inline")
}

// Metrics handles GET /metrics.
func (h *Handler) Metrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	fmt.Fprintln(w, "# sooqara metrics")
}

// ErrorResponse is the standard error envelope.
type ErrorResponse struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	RequestID  string `json:"request_id"`
}
