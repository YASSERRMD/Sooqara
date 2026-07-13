package agnes

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yasserrmd/sooqara/internal/limiter"
)

// stubJournal discards all records.
type stubJournal struct{}

func (s *stubJournal) Record(_ context.Context, _ journal.Activity) error { return nil }

// newTestClient creates a Client backed by an httptest server.
func newTestClient(t *testing.T, handler http.Handler) *Client {
	t.Helper()
	ts := httptest.NewServer(handler)
	t.Cleanup(func() { ts.Close() })

	cl := limiter.New(60, 3, &fakeClock{})

	return &Client{
		baseURL:    ts.URL,
		pollURL:    ts.URL + "/poll",
		apiKey:     "test-key",
		httpClient: ts.Client(),
		limiter:    cl,
		journal:    &stubJournal{},
	}
}

// fakeClock is a minimal clock for tests.
type fakeClock struct{}

func (f *fakeClock) Now() time.Time                        { return time.Now() }
func (f *fakeClock) Sleep(d time.Duration)                 { time.Sleep(d) }
