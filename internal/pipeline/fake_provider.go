package pipeline

import (
	"context"
	"encoding/json"

	"github.com/yasserrmd/sooqara/internal/provider"
)

// FakeProvider returns pre-programmed responses for testing.
type FakeProvider struct {
	ChatFn      func(ctx context.Context, req provider.ChatRequest) (provider.ChatResponse, error)
	GenerateImg func(ctx context.Context, req provider.ImageRequest) (provider.ImageResponse, error)
	CreateVid   func(ctx context.Context, req provider.VideoRequest) (provider.VideoJob, error)
	PollVid     func(ctx context.Context, videoID string) (provider.VideoResult, error)
}

func (f *FakeProvider) Chat(ctx context.Context, req provider.ChatRequest) (provider.ChatResponse, error) {
	if f.ChatFn != nil {
		return f.ChatFn(ctx, req)
	}
	return provider.ChatResponse{}, nil
}

func (f *FakeProvider) GenerateImage(ctx context.Context, req provider.ImageRequest) (provider.ImageResponse, error) {
	if f.GenerateImg != nil {
		return f.GenerateImg(ctx, req)
	}
	return provider.ImageResponse{}, nil
}

func (f *FakeProvider) CreateVideo(ctx context.Context, req provider.VideoRequest) (provider.VideoJob, error) {
	if f.CreateVid != nil {
		return f.CreateVid(ctx, req)
	}
	return provider.VideoJob{}, nil
}

func (f *FakeProvider) PollVideo(ctx context.Context, videoID string) (provider.VideoResult, error) {
	if f.PollVid != nil {
		return f.PollVid(ctx, videoID)
	}
	return provider.VideoResult{}, nil
}

func (f *FakeProvider) Name() string { return "fake" }

// ValidJSON returns a ChatResponse with valid analysis JSON.
func ValidJSON(analysis ProductAnalysis) provider.ChatResponse {
	data, _ := jsonMarshal(analysis)
	return provider.ChatResponse{
		Model: "agnes-2.0-flash",
		Choices: []provider.Choice{
			{Message: provider.Message{Role: "assistant", Content: data}},
		},
	}
}

// FenceWrappedJSON wraps JSON in markdown fences (simulating model behavior).
func FenceWrappedJSON(analysis ProductAnalysis) provider.ChatResponse {
	data, _ := jsonMarshal(analysis)
	return provider.ChatResponse{
		Model: "agnes-2.0-flash",
		Choices: []provider.Choice{
			{Message: provider.Message{Role: "assistant", Content: "```\n" + data + "\n```"}},
		},
	}
}

// MalformedJSON returns a response with invalid JSON.
func MalformedJSON() provider.ChatResponse {
	return provider.ChatResponse{
		Model: "agnes-2.0-flash",
		Choices: []provider.Choice{
			{Message: provider.Message{Role: "assistant", Content: "{broken json!!!"}},
		},
	}
}

// TooFewSettings returns valid JSON but with insufficient lifestyle settings.
func TooFewSettings(analysis ProductAnalysis) provider.ChatResponse {
	analysis.SuggestedLifestyleSettings = []string{"setting1"}
	data, _ := jsonMarshal(analysis)
	return provider.ChatResponse{
		Model: "agnes-2.0-flash",
		Choices: []provider.Choice{
			{Message: provider.Message{Role: "assistant", Content: data}},
		},
	}
}

func jsonMarshal(v any) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
