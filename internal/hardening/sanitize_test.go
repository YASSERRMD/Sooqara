package hardening

import "testing"

func TestSanitizeTitleBasic(t *testing.T) {
	got := SanitizeTitle("Hello World")
	if got != "Hello World" {
		t.Errorf("expected 'Hello World', got %q", got)
	}
}

func TestSanitizeTitleTrimSpaces(t *testing.T) {
	got := SanitizeTitle("  Hello  ")
	if got != "Hello" {
		t.Errorf("expected 'Hello', got %q", got)
	}
}

func TestSanitizeTitleControlChars(t *testing.T) {
	got := SanitizeTitle("Hello\x00World\x1f")
	if got != "HelloWorld" {
		t.Errorf("expected 'HelloWorld', got %q", got)
	}
}

func TestSanitizeTitleNewlines(t *testing.T) {
	got := SanitizeTitle("Line1\nLine2")
	if got != "Line1\nLine2" {
		t.Errorf("expected 'Line1\\nLine2', got %q", got)
	}
}

func TestSanitizeTitleEmpty(t *testing.T) {
	got := SanitizeTitle("")
	if got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}

func TestSanitizeTitleTabs(t *testing.T) {
	got := SanitizeTitle("\tHello\t")
	if got != "Hello" {
		t.Errorf("expected 'Hello', got %q", got)
	}
}
