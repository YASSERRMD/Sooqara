package hardening

import (
	"fmt"
	"strings"
)

// SafeHeaders returns a set of HTTP header names that are allowed in user input.
var SafeHeaders = map[string]bool{
	"content-type":     true,
	"accept":           true,
	"x-request-id":     true,
	"x-correlation-id": true,
}

// ValidateHeaderName checks that an HTTP header name is in the allowed set.
func ValidateHeaderName(name string) error {
	if !SafeHeaders[strings.ToLower(strings.TrimSpace(name))] {
		return &HeaderNotAllowedError{name: name}
	}
	return nil
}

// HeaderNotAllowedError is returned when a disallowed header name is encountered.
type HeaderNotAllowedError struct {
	name string
}

func (e *HeaderNotAllowedError) Error() string {
	return fmt.Sprintf("header %q is not allowed", e.name)
}
