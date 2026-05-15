// Package audit records non-mutating safety and operator intent signals (Sprint 6).
package audit

// Level is a coarse severity for structured audit lines.
type Level string

const (
	LevelInfo    Level = "info"
	LevelWarning Level = "warning"
)

// Record is one audit line suitable for JSON embedding in CLI output.
type Record struct {
	Level   Level  `json:"level"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

// UnsafeApplyFlagNonOperational warns that the unsafe apply flag is acknowledged but does not enable writes.
func UnsafeApplyFlagNonOperational() Record {
	return Record{
		Level:   LevelWarning,
		Code:    "KHR_AUDIT_UNSAFE_APPLY_NON_OPERATIONAL",
		Message: "--allow-unsafe-apply is acknowledged but non-operational in this build",
		Detail:  "Real cgroup/systemd writes remain disabled; a future apply gate, policy bundle attestation, and operator approval are required before any mutation path is enabled",
	}
}
