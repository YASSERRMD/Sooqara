package export

import (
	"strings"
	"testing"
)

func TestSlugifyEmpty(t *testing.T) {
	got := slugify("")
	if got != "" {
		t.Errorf("slugify(\"\") = %q, want empty", got)
	}
}

func TestSlugifySingleWord(t *testing.T) {
	got := slugify("Hello")
	if got != "hello" {
		t.Errorf("slugify(\"Hello\") = %q, want hello", got)
	}
}

func TestSlugifyMultipleWords(t *testing.T) {
	got := slugify("Hello World Foo Bar")
	if got != "hello-world-foo-bar" {
		t.Errorf("slugify = %q, want hello-world-foo-bar", got)
	}
}

func TestSlugifySpecialChars(t *testing.T) {
	got := slugify("Hello, World!")
	if got != "hello,-world!" {
		t.Errorf("slugify = %q, want hello,-world!", got)
	}
}
