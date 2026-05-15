package contracts

import (
	"fmt"
	"strings"
)

// AllowedDashboardSummaryFixtureFields lists top-level and nested keys permitted in a
// committed redacted Dashboard ?view=summary-shaped JSON fixture (test-only allowlist).
// Order is stable for docs and tests.
func AllowedDashboardSummaryFixtureFields() []string {
	return []string{
		"view",
		"generatedAt",
		"bootstrapSource",
		"decisionEngine.eligibleYielderCount",
		"decisionEngine.eligibleReceiverCount",
		"executionEngine.supportsApply",
		"executionEngine.summary.category",
		"windowsLaneRedacted.enabled",
		"windowsLaneRedacted.reason",
		"windowsLaneRedacted.blockers",
		"kubeVirtLegacyRedacted.present",
		"kubeVirtLegacyRedacted.providerMode",
	}
}

// ValidateContractClaimSafe enforces claim-safe posture on a mapped ParentFabricSummary
// together with Dashboard pre-map metadata: Windows off, no applyAllowed, and standard
// summary / forbidden-claim / apply-semantics validators.
func ValidateContractClaimSafe(summary ParentFabricSummary, meta RedactedLiveSummaryMetadata) error {
	if summary.WindowsLane.Enabled {
		return fmt.Errorf("claim-safe: windows lane must be disabled")
	}
	if summary.ExecutionEngine.ApplyAllowed {
		return fmt.Errorf("claim-safe: applyAllowed must be false")
	}
	if err := ValidateSummary(summary); err != nil {
		return fmt.Errorf("claim-safe: %w", err)
	}
	if err := ValidateNoForbiddenClaims(summary); err != nil {
		return fmt.Errorf("claim-safe: %w", err)
	}
	if err := ValidateApplySemantics(summary, meta.DashboardSupportsApply); err != nil {
		return fmt.Errorf("claim-safe: %w", err)
	}
	return nil
}

// ValidateSupportsApplyFalseEdge validates the supportsApply=false + dry_run_only edge:
// Dashboard did not advertise apply; contract must not claim applyAllowed; if the
// contract reports dryRunSupported, the Dashboard execution summary category must have
// been dry_run_only (pre-map), matching InferDryRunSupported semantics.
func ValidateSupportsApplyFalseEdge(summary ParentFabricSummary, meta RedactedLiveSummaryMetadata) error {
	if meta.DashboardSupportsApply {
		return fmt.Errorf("supportsApply false edge: expected DashboardSupportsApply false")
	}
	if summary.ExecutionEngine.ApplyAllowed {
		return fmt.Errorf("supportsApply false edge: applyAllowed must be false")
	}
	cat := strings.TrimSpace(meta.ExecutionSummaryCategory)
	if summary.ExecutionEngine.DryRunSupported {
		if !strings.EqualFold(cat, "dry_run_only") {
			return fmt.Errorf("supportsApply false edge: dryRunSupported true requires execution summary category dry_run_only, got %q", meta.ExecutionSummaryCategory)
		}
	}
	if summary.WindowsLane.Enabled {
		return fmt.Errorf("supportsApply false edge: windows lane must be disabled")
	}
	return nil
}
