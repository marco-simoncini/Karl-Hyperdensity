// Package safety enforces Sprint 5 hard rules: no real apply without explicit unsafe flag.
package safety

const (
	// FlagAllowUnsafeApply is the CLI flag name (documentation only in Sprint 5).
	FlagAllowUnsafeApply = "allow-unsafe-apply"
)

// MutationsForbidden returns true unless allowUnsafeApply is set.
func MutationsForbidden(allowUnsafeApply bool) bool {
	return !allowUnsafeApply
}
