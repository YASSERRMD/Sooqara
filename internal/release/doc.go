// Package release provides end-to-end release tooling for Sooqara.
//
// Functions:
//   - GenerateBuildFlags: produces -ldflags for go build
//   - WriteManifest / ReadManifest: JSON serialization of PackageInfo
//   - PrepareReleaseDir / CopyArtifact / ListArtifacts: staging helpers
//   - PackTarball: creates .tar.gz archives
//   - CheckEnvironment: validates build prerequisites
//   - GenerateReleaseNote: renders GitHub release notes
//   - RunRelease: orchestrates the full release workflow
//   - SHA256Checksum: verifies artifact integrity
//
// Usage:
//
//	release.RunRelease(".", "v1.0.0", commitHash)
package release
