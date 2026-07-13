package provider

// ImageRequest is an image generation request.
type ImageRequest struct {
	Model       string  `json:"model"`
	Prompt      string  `json:"prompt"`
	Size        string  `json:"size,omitempty"`
	Seed        *int64  `json:"seed,omitempty"`
	SourceImage string  `json:"source_image,omitempty"` // base64 for img2img
}

// ImageResponse is an image generation response.
type ImageResponse struct {
	Model  string   `json:"model"`
	Images []Image  `json:"images"`
}

// Image represents a generated image.
type Image struct {
	URL    string `json:"url"`
	Seed   *int64 `json:"seed,omitempty"`
}

// VideoRequest is a video creation request.
type VideoRequest struct {
	Model     string `json:"model"`
	Prompt    string `json:"prompt"`
	Height    int    `json:"height"`
	Width     int    `json:"width"`
	NumFrames int    `json:"num_frames"`
	FrameRate int    `json:"frame_rate"`
}

// VideoJob is the response from CreateVideo.
type VideoJob struct {
	VideoID string `json:"video_id"`
}

// VideoResult is the response from PollVideo.
type VideoResult struct {
	Status string `json:"status"`
	URL    string `json:"url,omitempty"`
	Err    string `json:"error,omitempty"`
}
