package agnes

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/yasserrmd/sooqara/internal/provider"
)

// Chat sends a chat completion request.
func (c *Client) Chat(ctx context.Context, req provider.ChatRequest) (provider.ChatResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	resp, err := c.doRequest(ctx, "POST", c.baseURL+"/chat/completions", req, "agnes.text")
	if err != nil {
		return provider.ChatResponse{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return provider.ChatResponse{}, fmt.Errorf("read response: %w", err)
	}

	var result provider.ChatResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return provider.ChatResponse{}, fmt.Errorf("unmarshal response: %w", err)
	}

	return result, nil
}
