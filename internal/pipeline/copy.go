package pipeline

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/yasserrmd/sooqara/internal/provider"
	"github.com/yasserrmd/sooqara/internal/store"
)

// Banned phrases that indicate hype copy.
var bannedPhrases = []string{
	"game-changing", "revolutionary", "unlock", "elevate", "seamlessly",
	"cutting-edge", "transform your", "take it to the next level",
	"look no further",
}

// emDashRe matches em dash characters.
var emDashRe = regexp.MustCompile("[—–]")

// GenerateCopy creates e-commerce copy from a ProductAnalysis.
func GenerateCopy(ctx context.Context, p provider.Provider, s *store.Store, job *store.Job, analysis *ProductAnalysis) (*store.Artifact, error) {
	tone := job.Tone
	if tone == "" {
		tone = DefaultTone
	}

	// Build the tool schema from CopySet.
	toolDef := provider.Tool{
		Type: "function",
		Function: provider.FunctionDef{
			Name:        "generate_copy",
			Description: fmt.Sprintf("Generate e-commerce copy for %s in %s tone", analysis.ProductName, tone),
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"title":             map[string]any{"type": "string", "maxLength": MaxTitleLen},
					"bullets":           map[string]any{"type": "array", "items": map[string]any{"type": "string", "maxLength": MaxBulletLen}, "minItems": ExpectedBulletCount, "maxItems": ExpectedBulletCount},
					"short_description": map[string]any{"type": "string", "maxLength": MaxShortDescLen},
					"long_description":  map[string]any{"type": "string"},
					"alt_text":          map[string]any{"type": "string", "maxLength": MaxAltTextLen},
					"meta_description":  map[string]any{"type": "string", "maxLength": MaxMetaDescLen},
					"keywords":          map[string]any{"type": "array", "items": map[string]any{"type": "string"}, "minItems": MinKeywordCount, "maxItems": MaxKeywordCount},
					"tone":              map[string]any{"type": "string"},
				},
				"required": []string{"title", "bullets", "short_description", "long_description", "alt_text", "meta_description", "keywords", "tone"},
			},
		},
	}

	userPrompt := fmt.Sprintf(
		"Generate e-commerce listing copy for: %s (%s). Materials: %v. Colours: %v. Shape: %s. Features: %v.",
		analysis.ProductName, analysis.Category, analysis.Materials, analysis.Colours,
		analysis.ShapeDescription, analysis.DistinguishingFeatures,
	)

	req := provider.ChatRequest{
		Model:    "agnes-2.0-flash",
		Messages: []provider.ChatMessage{{Role: "user", Content: userPrompt}},
		Tools:    []provider.Tool{toolDef},
	}

	copySet, err := runCopyWithRepair(ctx, p, req, tone)
	if err != nil {
		return nil, err
	}

	// Enforce constraints in Go.
	if err := enforceConstraints(copySet, job.Warning); err != nil {
		// Warning is set on the job later
		_ = err
	}

	// Persist
	payload, _ := jsonMarshal(copySet)
	a := store.NewArtifact(job.ID, store.ArtifactCopy, 0)
	a.Payload = &payload
	if err := store.CreateArtifact(s.DB, a); err != nil {
		return nil, fmt.Errorf("persist copy artifact: %w", err)
	}

	return a, nil
}

func runCopyWithRepair(ctx context.Context, p provider.Provider, req provider.ChatRequest, tone string) (*CopySet, error) {
	for attempt := 0; attempt <= 1; attempt++ {
		resp, err := p.Chat(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("chat call (attempt %d): %w", attempt+1, err)
		}
		if len(resp.Choices) == 0 {
			return nil, fmt.Errorf("empty chat response")
		}

		msg := resp.Choices[0].Message
		// Try tool call first
		if len(msg.ToolCalls) > 0 {
			tc := msg.ToolCalls[0]
			var cs CopySet
			if err := json.Unmarshal(tc.Function.Arguments, &cs); err != nil {
				return nil, fmt.Errorf("parse tool call: %w", err)
			}
			cs.Tone = tone
			return &cs, nil
		}

		// Fall back to JSON parsing
		raw := stripFences(msg.Content)
		var cs CopySet
		if err := json.Unmarshal([]byte(raw), &cs); err != nil {
			if attempt == 0 {
				req = provider.ChatRequest{
					Model: "agnes-2.0-flash",
					Messages: []provider.ChatMessage{
						{Role: "user", Content: fmt.Sprintf("Previous output was not valid JSON. Fix it: %s", raw)},
					},
				}
				continue
			}
			return nil, fmt.Errorf("parse copy: %w", err)
		}
		cs.Tone = tone
		return &cs, nil
	}
	return nil, fmt.Errorf("copy generation failed after repairs")
}

func enforceConstraints(cs *CopySet, warning *string) error {
	var warnings []string

	// Strip em dashes
	cs.Title = emDashRe.ReplaceAllString(cs.Title, ",")
	cs.ShortDescription = emDashRe.ReplaceAllString(cs.ShortDescription, ",")
	cs.LongDescription = emDashRe.ReplaceAllString(cs.LongDescription, ",")
	cs.AltText = emDashRe.ReplaceAllString(cs.AltText, ",")
	cs.MetaDescription = emDashRe.ReplaceAllString(cs.MetaDescription, ",")
	for i := range cs.Bullets {
		cs.Bullets[i] = emDashRe.ReplaceAllString(cs.Bullets[i], ",")
	}

	// Check banned phrases
	for _, bp := range bannedPhrases {
		if strings.Contains(strings.ToLower(cs.Title), bp) {
			warnings = append(warnings, fmt.Sprintf("title contains banned phrase: %s", bp))
		}
	}

	// Truncate title
	if len(cs.Title) > MaxTitleLen {
		truncated := truncateAtWord(cs.Title, MaxTitleLen)
		cs.Title = truncated
		warnings = append(warnings, fmt.Sprintf("title truncated from %d to %d chars", len(cs.Title), MaxTitleLen))
	}

	if len(warnings) > 0 {
		w := strings.Join(warnings, "; ")
		*warning = &w
	}

	return nil
}

func truncateAtWord(s string, max int) string {
	if len(s) <= max {
		return s
	}
	sub := s[:max]
	if idx := strings.LastIndex(sub, " "); idx > 0 {
		return sub[:idx]
	}
	return sub
}
