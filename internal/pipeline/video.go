package pipeline

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/yasserrmd/sooqara/internal/provider"
	"github.com/yasserrmd/sooqara/internal/store"
)

// VideoConfig holds video generation parameters.
type VideoConfig struct {
	Height    int
	Width     int
	NumFrames int
	FrameRate int
	Deadline  time.Duration
	PollStart time.Duration
	PollCeil  time.Duration
}

// DefaultVideoConfig returns standard video parameters.
func DefaultVideoConfig() VideoConfig {
	return VideoConfig{
		Height:    768,
		Width:     1152,
		NumFrames: 121,
		FrameRate: 24,
		Deadline:  15 * time.Minute,
		PollStart: 10 * time.Second,
		PollCeil:  60 * time.Second,
	}
}

// CreateVideoJob creates a video job and persists the video_id.
func CreateVideoJob(ctx context.Context, p provider.Provider, s *store.Store, job *store.Job, analysis *ProductAnalysis) (*store.Artifact, error) {
	bestImage, err := getBestImageVariant(s, job.ID)
	if err != nil {
		return nil, fmt.Errorf("find best image: %w", err)
	}

	prompt := fmt.Sprintf("%s, %s, %s",
		analysis.ProductName, StyleV1, MotionV1,
	)

	req := provider.VideoRequest{
		Model:     "agnes-video-v2.0",
		Prompt:    prompt,
		Height:    768,
		Width:     1152,
		NumFrames: 121,
		FrameRate: 24,
	}

	vj, err := p.CreateVideo(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("create video: %w", err)
	}

	a := store.NewArtifact(job.ID, store.ArtifactVideo, 0)
	a.Payload = strPtr(fmt.Sprintf(`{"video_id":"%s"}`, vj.VideoID))
	if err := store.CreateArtifact(s.DB, a); err != nil {
		return nil, fmt.Errorf("persist video_id: %w", err)
	}

	if err := store.Transition(s.DB, job.ID, store.StateImaging, store.StateVideoing); err != nil {
		return nil, fmt.Errorf("transition to videoing: %w", err)
	}

	return a, nil
}

func getBestImageVariant(s *store.Store, jobID string) (*store.Artifact, error) {
	arts, err := store.GetArtifactsByJob(s.DB, jobID)
	if err != nil {
		return nil, err
	}
	for _, a := range arts {
		if a.Kind == store.ArtifactImage && a.Seq == 0 {
			return a, nil
		}
	}
	return nil, fmt.Errorf("no image variant seq 0 found")
}

// VideoPoller scans for videoing jobs and polls for completion.
type VideoPoller struct {
	p      provider.Provider
	s      *store.Store
	config VideoConfig
	stop   chan struct{}
}

// NewVideoPoller creates a new poller.
func NewVideoPoller(p provider.Provider, s *store.Store, cfg VideoConfig) *VideoPoller {
	return &VideoPoller{p: p, s: s, config: cfg, stop: make(chan struct{})}
}

// Run starts the poll loop.
func (vp *VideoPoller) Run(ctx context.Context) {
	ticker := time.NewTicker(vp.config.PollStart)
	defer ticker.Stop()

	for {
		select {
		case <-vp.stop:
			return
		case <-ctx.Done():
			return
		case <-ticker.C:
			vp.pollOnce(ctx)
			// Reset ticker to start interval
			ticker.Reset(vp.config.PollStart)
		}
	}
}

func (vp *VideoPoller) pollOnce(ctx context.Context) {
	arts, err := vp.s.DB.Query(`SELECT id, payload FROM artifacts WHERE kind = 'video' AND path IS NULL`)
	if err != nil {
		return
	}
	defer arts.Close()

	for arts.Next() {
		var artID, payload string
		if err := arts.Scan(&artID, &payload); err != nil {
			continue
		}

		var videoID string
		fmt.Sscanf(payload, `{"video_id":"%s"}`, &videoID)
		if videoID == "" {
			continue
		}

		result, err := vp.p.PollVideo(ctx, videoID)
		if err != nil {
			continue
		}

		switch result.Status {
		case "done":
			// Download video
			vp.resolveVideo(ctx, artID, result.URL)
		case "failed":
			w := "video generation failed"
			vp.markWarning(artID, &w)
		}
	}
}

func (vp *VideoPoller) resolveVideo(ctx context.Context, artID, url string) {
	resp, err := httpGet(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)
	path, _ := vp.s.Blob.Put(ctx, "video.mp4", bytes.NewReader(data))

	vp.s.DB.Exec("UPDATE artifacts SET path = ? WHERE id = ?", path, artID)
}

func (vp *VideoPoller) markWarning(artID string, w *string) {
	vp.s.DB.Exec("UPDATE artifacts SET payload = ? WHERE id = ?", *w, artID)
}

func (vp *VideoPoller) Stop() { close(vp.stop) }

func httpGet(url string) (*http.Response, error) { return nil, nil }
