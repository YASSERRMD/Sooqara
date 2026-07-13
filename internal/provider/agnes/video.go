package agnes

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/yasserrmd/sooqara/internal/provider"
)

// CreateVideo sends a video creation request.
func (c *Client) CreateVideo(ctx context.Context, req provider.VideoRequest) (provider.VideoJob, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	resp, err := c.doRequest(ctx, "POST", c.baseURL+"/videos", req, "agnes.video.create")
	if err != nil {
		return provider.VideoJob{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return provider.VideoJob{}, fmt.Errorf("read response: %w", err)
	}

	var result provider.VideoJob
	if err := json.Unmarshal(body, &result); err != nil {
		return provider.VideoJob{}, fmt.Errorf("unmarshal response: %w", err)
	}

	return result, nil
}

// PollVideo polls the status of a video job.
func (c *Client) PollVideo(ctx context.Context, videoID string) (provider.VideoResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	params := url.Values{}
	params.Set("video_id", videoID)
	pollURL := c.pollURL + "?" + params.Encode()

	resp, err := c.doRequest(ctx, "GET", pollURL, nil, "agnes.video.poll")
	if err != nil {
		return provider.VideoResult{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return provider.VideoResult{}, fmt.Errorf("read response: %w", err)
	}

	var result provider.VideoResult
	if err := json.Unmarshal(body, &result); err != nil {
		return provider.VideoResult{}, fmt.Errorf("unmarshal response: %w", err)
	}

	return result, nil
}
