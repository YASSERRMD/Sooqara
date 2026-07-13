package api

import (
	"testing"
)

func TestStatusCodeConstants(t *testing.T) {
	if http.StatusOK != 200 {
		t.Error("OK should be 200")
	}
	if http.StatusAccepted != 202 {
		t.Error("Accepted should be 202")
	}
	if http.StatusBadRequest != 400 {
		t.Error("Bad Request should be 400")
	}
	if http.StatusNotFound != 404 {
		t.Error("Not Found should be 404")
	}
}
