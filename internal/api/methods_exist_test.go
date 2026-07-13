package api

import (
	"testing"
)

func TestHandlerHasRoutesMethod(t *testing.T) {
	h := &Handler{}
	_ = h.Routes
}

func TestHandlerHasCreateJobMethod(t *testing.T) {
	h := &Handler{}
	_ = h.CreateJob
}

func TestHandlerHasGetJobMethod(t *testing.T) {
	h := &Handler{}
	_ = h.GetJob
}

func TestHandlerHasCancelJobMethod(t *testing.T) {
	h := &Handler{}
	_ = h.CancelJob
}

func TestHandlerHasRegenerateMethod(t *testing.T) {
	h := &Handler{}
	_ = h.Regenerate
}
