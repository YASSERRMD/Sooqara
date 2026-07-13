package hardening

import "testing"

func TestHeaderNotAllowedErrorImplementsError(t *testing.T) {
	var err error = &HeaderNotAllowedError{name: "test"}
	_ = err // verifies compile-time interface satisfaction
}
