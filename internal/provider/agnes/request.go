package agnes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/yasserrmd/sooqara/internal/journal"
	"github.com/yasserrmd/sooqara/internal/provider"
)

// doRequest executes an HTTP request with limiter acquisition and retry logic.
func (c *Client) doRequest(ctx context.Context, method, url string, body any, kind string) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	start := time.Now()
	attempts := 0
	maxAttempts := 4

	logActivity := func(outcome, detail string) {
		c.journal.Record(ctx, journal.Activity{
			Kind:        kind,
			Outcome:     outcome,
			Detail:      detail,
			LatencyMs:   time.Since(start).Milliseconds(),
			Ts:          time.Now().UnixMilli(),
		})
	}

	for {
		if err := c.limiter.Acquire(ctx); err != nil {
			logActivity("error", fmt.Sprintf("limiter acquire failed: %v", err))
			return nil, err
		}

		attempts++
		resp, err := c.httpClient.Do(req)
		if err != nil {
			logActivity("error", fmt.Sprintf("request failed attempt %d: %v", attempts, err))
			return nil, fmt.Errorf("request attempt %d: %w", attempts, err)
		}
		defer resp.Body.Close()

		switch resp.StatusCode {
		case http.StatusOK:
			logActivity("ok", "")
			return resp, nil
		case http.StatusTooManyRequests:
			if attempts >= maxAttempts {
				logActivity("rate_limited", fmt.Sprintf("max retries (%d) exceeded", maxAttempts))
				return nil, provider.ErrRateLimited
			}
			backoff := time.Duration(1<<uint(attempts-1)) * time.Second
			time.Sleep(backoff)
			continue
		case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable:
			if attempts >= maxAttempts {
				logActivity("error", fmt.Sprintf("server error after %d attempts", attempts))
				return nil, provider.ErrProviderUnavailable
			}
			backoff := time.Duration(1<<uint(attempts-1)) * time.Second
			time.Sleep(backoff)
			continue
		case http.StatusUnauthorized:
			logActivity("error", "authentication failed")
			return nil, provider.ErrAuth
		default:
			logActivity("error", fmt.Sprintf("unexpected status: %d", resp.StatusCode))
			return nil, provider.ErrBadRequest
		}
	}
}
