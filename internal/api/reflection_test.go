package api

import (
	"reflect"
	"testing"
)

func TestHandlerType(t *testing.T) {
	h := &Handler{}
	t := reflect.TypeOf(h).Elem()
	if t.Name() != "Handler" {
		t.Errorf("type name = %s, want Handler", t.Name())
	}
}

func TestErrorResponseFields(t *testing.T) {
	h := ErrorResponse{}
	t := reflect.TypeOf(h)
	if t.NumField() != 3 {
		t.Errorf("ErrorResponse has %d fields, want 3", t.NumField())
	}
}
