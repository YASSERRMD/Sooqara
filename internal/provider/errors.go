package provider

import "errors"

// Typed errors for provider calls.
var (
	ErrRateLimited        = errors.New("provider: rate limited")
	ErrAuth               = errors.New("provider: authentication failed")
	ErrBadRequest         = errors.New("provider: bad request")
	ErrProviderUnavailable = errors.New("provider: service unavailable")
)

// IsRateLimited checks if the error is a rate limit error.
func IsRateLimited(err error) bool {
	return errors.Is(err, ErrRateLimited)
}

// IsAuth checks if the error is an authentication error.
func IsAuth(err error) bool {
	return errors.Is(err, ErrAuth)
}

// IsBadRequest checks if the error is a bad request error.
func IsBadRequest(err error) bool {
	return errors.Is(err, ErrBadRequest)
}

// IsProviderUnavailable checks if the error is a provider unavailable error.
func IsProviderUnavailable(err error) bool {
	return errors.Is(err, ErrProviderUnavailable)
}
