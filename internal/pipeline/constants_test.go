package pipeline

import (
	"testing"
)

func TestDefaultToneConstant(t *testing.T) {
	if DefaultTone != "clear and practical" {
		t.Errorf("DefaultTone = %q, want clear and practical", DefaultTone)
	}
}

func TestMaxTitleLenConstant(t *testing.T) {
	if MaxTitleLen != 60 {
		t.Errorf("MaxTitleLen = %d, want 60", MaxTitleLen)
	}
}

func TestExpectedBulletCountConstant(t *testing.T) {
	if ExpectedBulletCount != 5 {
		t.Errorf("ExpectedBulletCount = %d, want 5", ExpectedBulletCount)
	}
}

func TestKeywordRange(t *testing.T) {
	if MinKeywordCount != 6 {
		t.Errorf("MinKeywordCount = %d, want 6", MinKeywordCount)
	}
	if MaxKeywordCount != 10 {
		t.Errorf("MaxKeywordCount = %d, want 10", MaxKeywordCount)
	}
}

func TestAllBannedPhrasesNonEmpty(t *testing.T) {
	for i, bp := range bannedPhrases {
		if bp == "" {
			t.Errorf("bannedPhrases[%d] is empty", i)
		}
	}
}
