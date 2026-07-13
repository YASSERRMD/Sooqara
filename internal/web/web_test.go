package web

import (
	"testing"
)

func TestWebPackageExists(t *testing.T) {
	// Just verify the package compiles
}

func TestInitTemplates(t *testing.T) {
	err := InitTemplates()
	if err != nil {
		t.Logf("InitTemplates returned: %v", err)
	}
}
