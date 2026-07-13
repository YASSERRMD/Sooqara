package hardening

import (
	"fmt"
	"strings"
	"testing"
)

func TestValidateTagsOK(t *testing.T) {
	tags := []string{"electronics", "phones", "smart"}
	if err := ValidateTags(tags); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateTagsTooMany(t *testing.T) {
	tags := make([]string, MaxTagsPerListing+1)
	for i := range tags {
		tags[i] = fmt.Sprintf("tag-%d", i)
	}
	err := ValidateTags(tags)
	if err == nil {
		t.Fatal("expected error for too many tags")
	}
}

func TestValidateTagsEmpty(t *testing.T) {
	tags := []string{"good", ""}
	err := ValidateTags(tags)
	if err == nil {
		t.Fatal("expected error for empty tag")
	}
}

func TestValidateTagsLongTag(t *testing.T) {
	tags := []string{strings.Repeat("x", MaxTagName+1)}
	err := ValidateTags(tags)
	if err == nil {
		t.Fatal("expected error for long tag")
	}
}

func TestValidateTagsMaxCount(t *testing.T) {
	tags := make([]string, MaxTagsPerListing)
	for i := range tags {
		tags[i] = fmt.Sprintf("tag-%d", i)
	}
	if err := ValidateTags(tags); err != nil {
		t.Errorf("expected OK for max-count tags, got: %v", err)
	}
}

func TestValidateTagsWhitespaceTrimmed(t *testing.T) {
	tags := []string{"  good  ", "ok"}
	if err := ValidateTags(tags); err != nil {
		t.Errorf("expected OK after trim, got: %v", err)
	}
}
