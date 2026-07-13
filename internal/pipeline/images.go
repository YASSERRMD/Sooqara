package pipeline

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/yasserrmd/sooqara/internal/provider"
	"github.com/yasserrmd/sooqara/internal/store"
)

// GenerateVariants creates N image variants using seed-locked img2img.
func GenerateVariants(ctx context.Context, p provider.Provider, s *store.Store, job *store.Job, analysis *ProductAnalysis) ([]*store.Artifact, error) {
	if job.Seed == nil {
		s := int64(42)
		job.Seed = &s
	}

	var artifacts []*store.Artifact
	var warnings []string
	successCount := 0

	for i := 0; i < job.VariantCount; i++ {
		seed := *job.Seed + int64(i)
		setting := getSetting(analysis, i)
		prompt := buildPrompt(analysis, setting)

		imgData, err := os.ReadFile(job.SourceImagePath)
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("variant %d: read source: %v", i, err))
			continue
		}

		base64Img := base64.StdEncoding.EncodeToString(imgData)

		req := provider.ImageRequest{
			Model:       "agnes-image-2.1-flash",
			Prompt:      prompt,
			Size:        "1024x1024",
			Seed:        &seed,
			SourceImage: "data:image/jpeg;base64," + base64Img,
		}

		resp, err := p.GenerateImage(ctx, req)
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("variant %d: %v", i, err))
			continue
		}

		if len(resp.Images) == 0 {
			warnings = append(warnings, fmt.Sprintf("variant %d: no images returned", i))
			continue
		}

		// Save via blob store
		imgBytes, err := downloadImage(resp.Images[0].URL)
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("variant %d: download: %v", i, err))
			continue
		}

		path, err := s.Blob.Put(ctx, fmt.Sprintf("variant_%d.jpg", i), bytes.NewReader(imgBytes))
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("variant %d: save: %v", i, err))
			continue
		}

		a := store.NewArtifact(job.ID, store.ArtifactImage, i)
		a.Path = &path
		a.Seed = &seed
		a.Prompt = &prompt
		a.StyleVer = strPtr(StyleV1)
		if err := store.CreateArtifact(s.DB, a); err != nil {
			warnings = append(warnings, fmt.Sprintf("variant %d: persist: %v", i, err))
			continue
		}

		artifacts = append(artifacts, a)
		successCount++
	}

	if successCount == 0 {
		return nil, fmt.Errorf("zero variants succeeded")
	}

	if len(warnings) > 0 {
		w := fmt.Sprintf("partial success: %d/%d variants, warnings: %s", successCount, job.VariantCount, joinStrings(warnings))
		job.Warning = &w
	}

	return artifacts, nil
}

func getSetting(analysis *ProductAnalysis, i int) string {
	settings := analysis.SuggestedLifestyleSettings
	if i < len(settings) {
		return settings[i]
	}
	return "neutral studio background"
}

func buildPrompt(analysis *ProductAnalysis, setting string) string {
	return fmt.Sprintf("%s, %s, %s, %s, %s placed in %s, %s",
		analysis.ProductName,
		joinStrings(analysis.Materials),
		joinStrings(analysis.Colours),
		analysis.ShapeDescription,
		joinStrings(analysis.DistinguishingFeatures),
		setting,
		StyleV1,
	)
}

func joinStrings(ss []string) string {
	if len(ss) == 0 {
		return ""
	}
	result := ss[0]
	for _, s := range ss[1:] {
		result += ", " + s
	}
	return result
}

func downloadImage(url string) ([]byte, error) {
	return []byte{}, nil
}

func strPtr(s string) *string { return &s }
