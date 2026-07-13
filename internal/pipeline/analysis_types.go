package pipeline

// ProductAnalysis is the structured output from the vision analysis stage.
type ProductAnalysis struct {
	ProductName              string          `json:"product_name"`
	Category                 string          `json:"category"`
	Materials                []string        `json:"materials"`
	Colours                  []string        `json:"colours"`
	DominantColourHex        *string         `json:"dominant_colour_hex,omitempty"`
	ShapeDescription         string          `json:"shape_description"`
	DistinguishingFeatures   []string        `json:"distinguishing_features"`
	PhotoQuality             PhotoQuality    `json:"photo_quality"`
	SuggestedLifestyleSettings []string      `json:"suggested_lifestyle_settings"`
	Confidence               float64         `json:"confidence"`
}

// PhotoQuality describes the source photo characteristics.
type PhotoQuality struct {
	Lighting   string `json:"lighting"`
	Background string `json:"background"`
	Angle      string `json:"angle"`
}
