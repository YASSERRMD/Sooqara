package observability

// PricingModel defines token-to-cost rates for AI models.
type PricingModel struct {
	InputPerM    float64 // cost per million input tokens
	OutputPerM   float64 // cost per million output tokens
	ImagePerCall float64 // flat cost per image generation call
	VideoPerCall float64 // flat cost per video generation call
}

// DefaultPricing returns the standard Agnes AI pricing model.
func DefaultPricing() PricingModel {
	return PricingModel{
		InputPerM:    0.003,
		OutputPerM:   0.006,
		ImagePerCall: 0.04,
		VideoPerCall: 0.25,
	}
}

// CalculateChatCost computes the USD cost for a chat completion.
func (pm PricingModel) CalculateChatCost(inputTokens, outputTokens int64) float64 {
	inputCost := float64(inputTokens) / 1_000_000 * pm.InputPerM
	outputCost := float64(outputTokens) / 1_000_000 * pm.OutputPerM
	return roundToCents(inputCost + outputCost)
}

// CalculateImageCost returns the flat cost for an image generation call.
func (pm PricingModel) CalculateImageCost(count int) float64 {
	return roundToCents(pm.ImagePerCall * float64(count))
}

// CalculateVideoCost returns the flat cost for a video generation call.
func (pm PricingModel) CalculateVideoCost(count int) float64 {
	return roundToCents(pm.VideoPerCall * float64(count))
}

// roundToCents rounds a float to two decimal places.
func roundToCents(v float64) float64 {
	return float64(int(v*100+0.5)) / 100
}
