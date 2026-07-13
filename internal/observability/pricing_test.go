package observability

import "testing"

func TestCalculateChatCost(t *testing.T) {
	pm := DefaultPricing()
	cost := pm.CalculateChatCost(1000, 500)
	expected := roundToCents(1000/1_000_000*0.003 + 500/1_000_000*0.006)
	if cost != expected {
		t.Errorf("expected %f, got %f", expected, cost)
	}
}

func TestCalculateChatCostZero(t *testing.T) {
	pm := DefaultPricing()
	cost := pm.CalculateChatCost(0, 0)
	if cost != 0 {
		t.Errorf("expected 0, got %f", cost)
	}
}

func TestCalculateImageCost(t *testing.T) {
	pm := DefaultPricing()
	cost := pm.CalculateImageCost(3)
	expected := roundToCents(0.04 * 3)
	if cost != expected {
		t.Errorf("expected %f, got %f", expected, cost)
	}
}

func TestCalculateVideoCost(t *testing.T) {
	pm := DefaultPricing()
	cost := pm.CalculateVideoCost(2)
	expected := roundToCents(0.25 * 2)
	if cost != expected {
		t.Errorf("expected %f, got %f", expected, cost)
	}
}

func TestRoundToCents(t *testing.T) {
	tests := []struct {
		input  float64
		expect float64
	}{
		{0.0004, 0.0},
		{0.0005, 0.0},
		{0.001, 0.0},
		{0.005, 0.01},
		{1.234, 1.23},
		{1.235, 1.24},
	}
	for _, tt := range tests {
		got := roundToCents(tt.input)
		if got != tt.expect {
			t.Errorf("roundToCents(%f) = %f, want %f", tt.input, got, tt.expect)
		}
	}
}

func TestDefaultPricingNonZero(t *testing.T) {
	pm := DefaultPricing()
	if pm.InputPerM <= 0 {
		t.Error("InputPerM should be positive")
	}
	if pm.OutputPerM <= 0 {
		t.Error("OutputPerM should be positive")
	}
	if pm.ImagePerCall <= 0 {
		t.Error("ImagePerCall should be positive")
	}
	if pm.VideoPerCall <= 0 {
		t.Error("VideoPerCall should be positive")
	}
}
