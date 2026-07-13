package pipeline

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/yasserrmd/sooqara/internal/provider"
	"github.com/yasserrmd/sooqara/internal/store"
)

// Analyse runs vision analysis on the source image.
func Analyse(ctx context.Context, p provider.Provider, s *store.Store, jobID, sourceImagePath string, variantCount int) (*store.Artifact, error) {
	systemPrompt := `Return ONLY a JSON object. No prose. No markdown fences. No commentary.

Target schema:
{
  "product_name": "string",
  "category": "string",
  "materials": ["string"],
  "colours": ["string"],
  "dominant_colour_hex": "#RRGGBB",
  "shape_description": "string",
  "distinguishing_features": ["string"],
  "photo_quality": { "lighting": "...", "background": "...", "angle": "..." },
  "suggested_lifestyle_settings": ["string", ...],
  "confidence": 0.0
}`

	imgData, err := os.ReadFile(sourceImagePath)
	if err != nil {
		return nil, fmt.Errorf("read source image: %w", err)
	}
	base64Img := base64.StdEncoding.EncodeToString(imgData)

	contentParts := []map[string]any{
		{"type": "text", "text": systemPrompt},
		{
			"type": "image_url",
			"image_url": map[string]string{
				"url": "data:image/jpeg;base64," + base64Img,
			},
		},
	}

	req := provider.ChatRequest{
		Model:    "agnes-2.0-flash",
		Messages: []provider.ChatMessage{{Role: "user", Content: contentParts}},
	}

	analysis, err := runWithRepair(ctx, p, req, variantCount, systemPrompt)
	if err != nil {
		return nil, err
	}

	if err := validateAnalysis(analysis, variantCount); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	payload, _ := json.Marshal(analysis)
	a := store.NewArtifact(jobID, store.ArtifactAnalysis, 0)
	payloadStr := string(payload)
	a.Payload = &payloadStr
	if err := store.CreateArtifact(s.DB, a); err != nil {
		return nil, fmt.Errorf("persist analysis artifact: %w", err)
	}

	return a, nil
}

func runWithRepair(ctx context.Context, p provider.Provider, req provider.ChatRequest, variantCount int, systemPrompt string) (*ProductAnalysis, error) {
	for attempt := 0; attempt <= 1; attempt++ {
		resp, err := p.Chat(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("chat call (attempt %d): %w", attempt+1, err)
		}

		if len(resp.Choices) == 0 {
			return nil, fmt.Errorf("empty chat response")
		}

		raw := resp.Choices[0].Message.Content
		raw = stripFences(raw)

		var analysis ProductAnalysis
		if err := json.Unmarshal([]byte(raw), &analysis); err != nil {
			if attempt == 0 {
				req = provider.ChatRequest{
					Model: "agnes-2.0-flash",
					Messages: []provider.ChatMessage{
						{Role: "system", Content: systemPrompt},
						{Role: "user", Content: fmt.Sprintf("Previous output was malformed:\n%s\n\nFix it and return valid JSON only.", raw)},
					},
				}
				continue
			}
			return nil, fmt.Errorf("parse analysis: %w", err)
		}

		if len(analysis.SuggestedLifestyleSettings) < variantCount {
			if attempt == 0 {
				req = provider.ChatRequest{
					Model: "agnes-2.0-flash",
					Messages: []provider.ChatMessage{
						{Role: "system", Content: systemPrompt},
						{Role: "user", Content: fmt.Sprintf("Need at least %d lifestyle settings, got %d. Return a corrected JSON.", variantCount, len(analysis.SuggestedLifestyleSettings))},
					},
				}
				continue
			}
			return nil, fmt.Errorf("only %d lifestyle settings, need at least %d", len(analysis.SuggestedLifestyleSettings), variantCount)
		}

		return &analysis, nil
	}

	return nil, fmt.Errorf("analysis failed after repairs")
}

func stripFences(s string) string {
	s = strings.TrimSpace(s)
	re := regexp.MustCompile("^```(?:json)?\\s*|\\s*```$")
	return re.ReplaceAllString(s, "")
}
