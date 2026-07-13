// Package release provides build-time version injection, changelog formatting,
// and release metadata for Sooqara.
//
// Key functions:
//   - GenerateBuildFlags: produces -ldflags for go build
//   - FormatChangelog: renders ChangelogEntry slices to markdown
//   - Metadata.String: formats release information
//
// Example build command:
//
//	go build -ldflags "$(go run internal/release/build_flags.go)" -o bin/sooqara ./cmd/sooqara
package release
