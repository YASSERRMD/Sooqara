package pipeline

import (
	"context"
	"testing"

	"github.com/yasserrmd/sooqara/internal/provider"
)

func TestFakeProviderAllMethods(t *testing.T) {
	fake := &FakeProvider{}

	// Chat
	resp, err := fake.Chat(context.Background(), provider.ChatRequest{})
	if err != nil {
		t.Errorf("Chat error: %v", err)
	}
	_ = resp

	// GenerateImage
	resp2, err := fake.GenerateImage(context.Background(), provider.ImageRequest{})
	if err != nil {
		t.Errorf("GenerateImage error: %v", err)
	}
	_ = resp2

	// CreateVideo
	resp3, err := fake.CreateVideo(context.Background(), provider.VideoRequest{})
	if err != nil {
		t.Errorf("CreateVideo error: %v", err)
	}
	_ = resp3

	// PollVideo
	resp4, err := fake.PollVideo(context.Background(), "test")
	if err != nil {
		t.Errorf("PollVideo error: %v", err)
	}
	_ = resp4

	// Name
	if fake.Name() != "fake" {
		t.Errorf("Name() = %s, want fake", fake.Name())
	}
}
