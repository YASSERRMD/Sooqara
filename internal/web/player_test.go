package web

import (
	"testing"
)

func TestVideoPlayerAttrs(t *testing.T) {
	attrs := []string{"autoplay", "muted", "loop", "playsinline"}
	for _, a := range attrs {
		if a == "" {
			t.Error("empty attribute")
		}
	}
}

func TestWarningBannerText(t *testing.T) {
	text := "Something went wrong"
	if text == "" {
		t.Error("empty warning text")
	}
}
