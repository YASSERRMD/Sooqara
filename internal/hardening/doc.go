// Package hardening provides input validation and security utilities for Sooqara.
//
// Functions:
//   - SanitizeTitle: strips control characters from user-supplied titles.
//   - ValidateTitle: checks title is non-empty and within MaxTitleLength.
//   - ValidateURL: ensures a string is a well-formed http/https URL.
//   - ValidateAPIKey: enforces min/max length for API keys.
//   - ValidateCopyLength / ValidateDescriptionLength: enforce output size caps.
//   - ValidateTags: checks tag count and per-tag length.
//   - SafePath: prevents directory traversal in file operations.
//   - ValidateHeaderName: allows only a whitelist of HTTP header names.
//
// All validators return descriptive errors suitable for API responses.
package hardening
