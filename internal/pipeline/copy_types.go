package pipeline

// CopySet is the structured output from the copy generation stage.
type CopySet struct {
	Title             string   `json:"title"`
	Bullets           []string `json:"bullets"`
	ShortDescription  string   `json:"short_description"`
	LongDescription   string   `json:"long_description"`
	AltText           string   `json:"alt_text"`
	MetaDescription   string   `json:"meta_description"`
	Keywords          []string `json:"keywords"`
	Tone              string   `json:"tone"`
}

// CopySet constraints.
const (
	MaxTitleLen          = 60
	MaxBulletLen         = 120
	MaxShortDescLen      = 300
	MaxAltTextLen        = 125
	MaxMetaDescLen       = 155
	ExpectedBulletCount  = 5
	MinKeywordCount      = 6
	MaxKeywordCount      = 10
	DefaultTone          = "clear and practical"
)
