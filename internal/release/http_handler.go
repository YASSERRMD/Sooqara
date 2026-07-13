package release

import (
	"encoding/json"
	"net/http"
	"runtime"

	"github.com/yasserrmd/sooqara/internal/version"
)

// HandleBuildInfo registers an HTTP handler that returns build info as JSON.
func HandleBuildInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	info := map[string]string{
		"version":    version.Version,
		"commit":     version.Commit,
		"built":      version.BuildTime,
		"go_version": runtime.Version(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}
