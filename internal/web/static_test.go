package web

import (
	"strings"
	"testing"
)

func TestStaticAssetPaths(t *testing.T) {
	paths := []string{"/static/css/style.css", "/static/js/app.js", "/static/images/logo.png"}
	for _, p := range paths {
		if !strings.HasPrefix(p, "/static/") {
			t.Errorf("path %q doesn't start with /static/", p)
		}
	}
}
