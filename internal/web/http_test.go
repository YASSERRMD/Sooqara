package web

import (
	"net/http"
	"testing"
)

func TestServeFiles(t *testing.T) {
	mux := http.NewServeMux()
	ServeFiles(mux)
}

func TestHTTPServeMux(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})
}
