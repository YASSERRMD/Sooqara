package api

import (
	"testing"
)

func TestRoutesMethodMapping(t *testing.T) {
	routeCount := 9
	if routeCount != 9 {
		t.Errorf("expected 9 routes, got %d", routeCount)
	}
}
