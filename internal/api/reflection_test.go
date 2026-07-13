package api

import (
	"reflect"
	"testing"
)

func TestHandlerType(t *testing.T) {
	h := &Handler{}
	typ := reflect.TypeOf(h).Elem()
	if typ.Name() != "Handler" {
		t.Errorf("type name = %s, want Handler", typ.Name())
	}
}

func TestErrorResponseFields(t *testing.T) {
	h := ErrorResponse{}
	typ := reflect.TypeOf(h)
	if typ.NumField() != 3 {
		t.Errorf("ErrorResponse has %d fields, want 3", typ.NumField())
	}
}
