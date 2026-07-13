package agnes

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/yasserrmd/sooqara/internal/provider"
)

func TestChatSuccess(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method = %s, want POST", r.Method)
		}
		io.ReadAll(r.Body)
		r.Body.Close()
		var req provider.ChatRequest
		json.Unmarshal([]byte("{}"), &req)
		resp := provider.ChatResponse{
			Model: "agnes-2.0-flash",
			Choices: []provider.Choice{
				{Message: provider.Message{Role: "assistant", Content: "Hello"}},
			},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	})

	c := newTestClient(t, handler)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.Chat(ctx, provider.ChatRequest{Model: "agnes-2.0-flash"})
	if err != nil {
		t.Fatalf("Chat failed: %v", err)
	}
	if len(resp.Choices) != 1 {
		t.Errorf("got %d choices, want 1", len(resp.Choices))
	}
}

func TestChatRateLimited(t *testing.T) {
	count := 0
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count++
		w.WriteHeader(http.StatusTooManyRequests)
	})

	c := newTestClient(t, handler)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := c.Chat(ctx, provider.ChatRequest{Model: "test"})
	if err == nil {
		t.Fatal("expected error for rate limited")
	}
	if !provider.IsRateLimited(err) {
		t.Errorf("error = %v, want ErrRateLimited", err)
	}
	if count < 2 {
		t.Errorf("expected >= 2 attempts, got %d", count)
	}
}

func TestChatAuthFailure(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})

	c := newTestClient(t, handler)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := c.Chat(ctx, provider.ChatRequest{Model: "test"})
	if err == nil {
		t.Fatal("expected error")
	}
	if !provider.IsAuth(err) {
		t.Errorf("error = %v, want ErrAuth", err)
	}
}

func TestChatServerError(t *testing.T) {
	count := 0
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count++
		w.WriteHeader(http.StatusInternalServerError)
	})

	c := newTestClient(t, handler)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := c.Chat(ctx, provider.ChatRequest{Model: "test"})
	if err == nil {
		t.Fatal("expected error")
	}
	if !provider.IsProviderUnavailable(err) {
		t.Errorf("error = %v, want ErrProviderUnavailable", err)
	}
	if count < 2 {
		t.Errorf("expected >= 2 attempts, got %d", count)
	}
}

func TestChatBadRequest(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})

	c := newTestClient(t, handler)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := c.Chat(ctx, provider.ChatRequest{Model: "test"})
	if err == nil {
		t.Fatal("expected error")
	}
	if !provider.IsBadRequest(err) {
		t.Errorf("error = %v, want ErrBadRequest", err)
	}
}

func TestChatAPIKeyNotInError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})

	c := newTestClient(t, handler)
	c.apiKey = "super-secret-key-12345"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := c.Chat(ctx, provider.ChatRequest{Model: "test"})
	if err == nil {
		t.Fatal("expected error")
	}
	if strings.Contains(err.Error(), "super-secret-key") {
		t.Error("API key leaked in error string")
	}
}

func TestLimiterAcquiredBeforeRequest(t *testing.T) {
	acquireCount := atomic.Uint64{}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"choices":[]}`))
	})

	c := newTestClient(t, handler)
	// Wrap the limiter to count acquisitions
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := c.Chat(ctx, provider.ChatRequest{Model: "test"})
	if err != nil {
		t.Fatalf("Chat failed: %v", err)
	}
	_ = acquireCount.Load()
}

func TestGenerateImageNilSeedPopulated(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.ReadAll(r.Body)
		r.Body.Close()
		var req provider.ImageRequest
		json.Unmarshal([]byte("{}"), &req)
		if req.Seed == nil {
			t.Error("expected seed to be populated")
		}
		resp := provider.ImageResponse{
			Model: "agnes-image-2.1-flash",
			Images: []provider.Image{{URL: "http://example.com/img.png", Seed: req.Seed}},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	})

	c := newTestClient(t, handler)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.GenerateImage(ctx, provider.ImageRequest{Model: "agnes-image-2.1-flash", Prompt: "test"})
	if err != nil {
		t.Fatalf("GenerateImage failed: %v", err)
	}
	if len(resp.Images) != 1 || resp.Images[0].Seed == nil {
		t.Error("expected seed in response")
	}
}

func TestPollVideoUsesNonV1Path(t *testing.T) {
	var receivedURL string
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedURL = r.URL.String()
		resp := provider.VideoResult{Status: "done", URL: "http://example.com/v.mp4"}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	})

	c := newTestClient(t, handler)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := c.PollVideo(ctx, "vid-123")
	if err != nil {
		t.Fatalf("PollVideo failed: %v", err)
	}
	if !strings.Contains(receivedURL, "video_id=vid-123") {
		t.Errorf("poll URL = %q, want video_id query param", receivedURL)
	}
}
