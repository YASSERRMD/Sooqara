package agnes

import (
	"time"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"

	"github.com/yasserrmd/sooqara/internal/provider"
)

// GenerateImage sends an image generation request.
func (c *Client) GenerateImage(ctx context.Context, req provider.ImageRequest) (provider.ImageResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 180*time.Second)
	defer cancel()

	// Generate seed if nil
	if req.Seed == nil {
		s := rand.Int63n(1<<30)
		req.Seed = &s
	}

	resp, err := c.doRequest(ctx, "POST", c.baseURL+"/images/generations", req, "agnes.image")
	if err != nil {
		return provider.ImageResponse{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return provider.ImageResponse{}, fmt.Errorf("read response: %w", err)
	}

	var result provider.ImageResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return provider.ImageResponse{}, fmt.Errorf("unmarshal response: %w", err)
	}

	// Ensure seeds are propagated back
	for i := range result.Images {
		if result.Images[i].Seed == nil {
			s := *req.Seed
			result.Images[i].Seed = &s
		}
	}

	return result, nil
}
