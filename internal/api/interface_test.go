package api

import (
	"testing"
)

func TestHandlerImplementsHTTPHandler(t *testing.T) {
	var _ http.Handler = (*Handler)(nil)
}
