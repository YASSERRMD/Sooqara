package agnes

import (
	"net/http"
	"time"

	"github.com/yasserrmd/sooqara/internal/journal"
	"github.com/yasserrmd/sooqara/internal/limiter"
)

// Journaler records activities.
type Journaler interface {
	Record(ctx context.Context, a journal.Activity) error
}

// Client implements provider.Provider for Agnes AI.
type Client struct {
	baseURL    string
	pollURL    string
	apiKey     string
	httpClient *http.Client
	limiter    *limiter.Limiter
	journal    Journaler
	timeout    time.Duration
}

// New creates a new Agnes provider client.
func New(baseURL, pollURL, apiKey string, lm *limiter.Limiter, j *journal.Journal) *Client {
	return &Client{
		baseURL:    baseURL,
		pollURL:    pollURL,
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 120 * time.Second},
		limiter:    lm,
		journal:    j,
		timeout:    120 * time.Second,
	}
}

// Name returns the provider name.
func (c *Client) Name() string { return "agnes" }
