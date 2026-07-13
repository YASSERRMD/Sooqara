package agnes

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yasserrmd/sooqara/internal/journal"
	"github.com/yasserrmd/sooqara/internal/limiter"
)

// newTestClient creates a Client backed by an httptest server.
func newTestClient(t *testing.T, handler http.Handler) *Client {
	t.Helper()
	ts := httptest.NewServer(handler)
	t.Cleanup(func() { ts.Close() })

	cl := limiter.New(60, 3, &fakeClock{})
	j := &journal.Journal{}

	return &Client{
		baseURL:   ts.URL,
		pollURL:   ts.URL + "/poll",
		apiKey:    "test-key",
		httpClient: ts.Client(),
		limiter:   cl,
		journal:   j,
	}
}

// fakeClock is a minimal clock for tests.
type fakeClock struct{}

func (f *fakeClock) Now() time.Time                        { return time.Now() }
func (f *fakeClock) Sleep(d time.Duration)                 { time.Sleep(d) }
