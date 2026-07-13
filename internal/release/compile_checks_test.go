package release

import "net/http"

// compile-time verification
var (
	_ http.HandlerFunc = HandleBuildInfo
)
