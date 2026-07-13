package api

import (
	"errors"
	"fmt"
	"testing"
)

func TestErrorWrapping(t *testing.T) {
	base := errors.New("base error")
	wrapped := fmt.Errorf("wrapped: %w", base)
	if !errors.Is(wrapped, base) {
		t.Error("errors.Is should return true for wrapped error")
	}
}

func TestErrorNew(t *testing.T) {
	err := errors.New("test error")
	if err == nil {
		t.Fatal("expected non-nil error")
	}
}
