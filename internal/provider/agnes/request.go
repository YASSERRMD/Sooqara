package agnes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
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

	for {
		if err := c.limiter.Acquire(ctx); err != nil {
			c.journal.Record(ctx, journal.Activity{
				Kind:        kind,
				Outcome:     "error",
				Detail:      fmt.Sprintf("limiter acquire failed: %v", err),
				LatencyMs:   time.Since(start).Milliseconds(),
			})
			return nil, err
		}

		attempts++
		resp, err := c.httpClient.Do(req)
		if err != nil {
			c.journal.Record(ctx, journal.Activity{
				Kind:        kind,
				Outcome:     "error",
				Detail:      fmt.Sprintf("request failed attempt %d: %v", attempts, err),
				LatencyMs:   time.Since(start).Milliseconds(),
			})
			return nil, fmt.Errorf("request attempt %d: %w", attempts, err)
		}
		defer resp.Body.Close()

		switch resp.StatusCode {
		case http.StatusOK:
			c.journal.Record(ctx, journal.Activity{
				Kind:        kind,
				Outcome:     "ok",
				LatencyMs:   time.Since(start).Milliseconds(),
			})
			return resp, nil
		case http.StatusTooManyRequests:
			if attempts >= maxAttempts {
				c.journal.Record(ctx, journal.Activity{
					Kind:        kind,
					Outcome:     "rate_limited",
					Detail:      fmt.Sprintf("max retries (%d) exceeded", maxAttempts),
					LatencyMs:   time.Since(start).Milliseconds(),
				})
				return nil, provider.ErrRateLimited
			}
			backoff := time.Duration(1<<uint(attempts-1)) * time.Second
			time.Sleep(backoff)
			continue
		case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable:
			if attempts >= maxAttempts {
				c.journal.Record(ctx, journal.Activity{
					Kind:        kind,
					Outcome:     "error",
					Detail:      fmt.Sprintf("server error after %d attempts", attempts),
					LatencyMs:   time.Since(start).Milliseconds(),
				})
				return nil, provider.ErrProviderUnavailable
			}
			backoff := time.Duration(1<<uint(attempts-1)) * time.Second
			time.Sleep(backoff)
			continue
		case http.StatusUnauthorized:
			c.journal.Record(ctx, journal.Activity{
				Kind:        kind,
				Outcome:     "error",
				Detail:      "authentication failed",
				LatencyMs:   time.Since(start).Milliseconds(),
			})
			return nil, provider.ErrAuth
		default:
			c.journal.Record(ctx, journal.Activity{
				Kind:        kind,
				Outcome:     "error",
				Detail:      fmt.Sprintf("unexpected status: %d", resp.StatusCode),
				LatencyMs:   time.Since(start).Milliseconds(),
			})
			return nil, provider.ErrBadRequest
		}
	}
}
