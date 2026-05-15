// Package safety enforces hard rules: no real apply in Sprint 6 builds.
package safety

const (
	// FlagAllowUnsafeApply is the CLI flag name (audit-only in Sprint 6).
	FlagAllowUnsafeApply = "allow-unsafe-apply"
)

// MutationsForbidden is true whenever host mutation paths must remain disabled.
// Sprint 6: always true — writes are never enabled from this binary regardless of CLI flags.
func MutationsForbidden(allowUnsafeApply bool) bool {
	_ = allowUnsafeApply
	return true
}

// UnsafeApplyRequested reports whether the operator passed the unsafe apply flag (audit signal only).
func UnsafeApplyRequested(allowUnsafeApply bool) bool {
	return allowUnsafeApply
}
