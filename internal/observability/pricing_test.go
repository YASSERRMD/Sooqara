package observability

import "testing"

func TestRoundToCents(t *testing.T) {
	tests := []struct {
		input  float64
		expect float64
	}{
		{0.0004, 0.0},
		{0.0005, 0.0},
		{0.001, 0.0},
		{0.004, 0.0},
		{0.005, 0.01},
		{0.01, 0.01},
		{1.234, 1.23},
		{1.235, 1.24},
		{0.04, 0.04},
		{0.25, 0.25},
	}
	for _, tt := range tests {
		got := roundToCents(tt.input)
		if got != tt.expect {
			t.Errorf("roundToCents(%f) = %f, want %f", tt.input, got, tt.expect)
		}
	}
}
