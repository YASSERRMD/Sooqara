package pipeline

import (
	"strings"
	"testing"
)

func TestBannedPhrasesAreLowercase(t *testing.T) {
	for _, bp := range bannedPhrases {
		if bp != strings.ToLower(bp) {
			t.Errorf("banned phrase %q should be lowercase for case-insensitive matching", bp)
		}
	}
}
