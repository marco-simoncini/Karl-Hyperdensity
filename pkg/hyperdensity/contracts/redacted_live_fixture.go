package contracts

import (
	"fmt"
	"strings"
)

// RedactedLiveSummaryMetadata carries pre-map Dashboard fields needed to validate a
// redacted live capture fixture against the mapped ParentFabricSummary. Pure data — no IO.
type RedactedLiveSummaryMetadata struct {
	// DashboardSupportsApply is executionEngine.supportsApply from the redacted Dashboard JSON before mapping.
	DashboardSupportsApply bool
}

// ValidateRedactedLiveSummaryFixture checks that a mapped ParentFabricSummary is consistent
// with a redacted live / live-like extraction anchor: source marker, Windows posture,
// and standard contract validators including apply semantics vs Dashboard supportsApply.
//
// No HTTP, Kubernetes, or Dashboard imports — suitable for test-only pipelines.
func ValidateRedactedLiveSummaryFixture(summary ParentFabricSummary, meta RedactedLiveSummaryMetadata) error {
	src := strings.ToLower(strings.TrimSpace(summary.Source))
	if src == "" {
		return fmt.Errorf("redacted live fixture: source is required")
	}
	if !strings.Contains(src, "redacted") && !strings.Contains(src, "live-capture-redacted") {
		return fmt.Errorf("redacted live fixture: source %q must contain \"redacted\" or \"live-capture-redacted\"", summary.Source)
	}
	if summary.WindowsLane.Enabled {
		return fmt.Errorf("redacted live fixture: windows lane must be disabled")
	}
	if err := ValidateSummary(summary); err != nil {
		return fmt.Errorf("redacted live fixture: %w", err)
	}
	if err := ValidateNoForbiddenClaims(summary); err != nil {
		return fmt.Errorf("redacted live fixture: %w", err)
	}
	if err := ValidateApplySemantics(summary, meta.DashboardSupportsApply); err != nil {
		return fmt.Errorf("redacted live fixture: %w", err)
	}
	return nil
}
