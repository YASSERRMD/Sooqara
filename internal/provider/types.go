package provider

import "context"

// Provider is the interface for AI model calls.
type Provider interface {
	Chat(ctx context.Context, req ChatRequest) (ChatResponse, error)
	GenerateImage(ctx context.Context, req ImageRequest) (ImageResponse, error)
	CreateVideo(ctx context.Context, req VideoRequest) (VideoJob, error)
	PollVideo(ctx context.Context, videoID string) (VideoResult, error)
	Name() string
}

// ChatRequest is a chat completion request.
type ChatRequest struct {
	Model     string       `json:"model"`
	Messages  []ChatMessage `json:"messages"`
	Temperature *float64   `json:"temperature,omitempty"`
	MaxTokens    *int      `json:"max_tokens,omitempty"`
	Stream     bool        `json:"stream,omitempty"`
	Tools      []Tool      `json:"tools,omitempty"`
}

// ChatMessage represents a single message in a conversation.
type ChatMessage struct {
	Role    string        `json:"role"`
	Content interface{}   `json:"content"`
}

// Tool represents a function tool for structured output.
type Tool struct {
	Type     string    `json:"type"`
	Function FunctionDef `json:"function"`
}

// FunctionDef defines a callable function.
type FunctionDef struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Parameters  map[string]any    `json:"parameters"`
}

// ChatResponse is a chat completion response.
type ChatResponse struct {
	Model    string    `json:"model"`
	Choices  []Choice  `json:"choices"`
}

// Choice is a single choice in a chat response.
type Choice struct {
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// Message is the response message.
type Message struct {
	Role     string      `json:"role"`
	Content  string      `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

// ToolCall represents a tool invocation from the model.
type ToolCall struct {
	ID       string           `json:"id"`
	Type     string           `json:"type"`
	Function ToolFunctionCall `json:"function"`
}

// ToolFunctionCall is the function call details.
type ToolFunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}
