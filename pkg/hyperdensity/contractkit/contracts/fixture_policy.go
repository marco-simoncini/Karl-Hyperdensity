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

// MissingOptionalGeneratedAtDefault is the stable contract generatedAt when the Dashboard
// fixture has generatedAt null, absent, or empty (M7 missing-optional edge; test-only mapper).
const MissingOptionalGeneratedAtDefault = "redacted-generatedAt-unavailable"

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

// ValidateMissingOptionalFieldsEdge validates the third fixture edge: optional Dashboard
// fields absent (null generatedAt, absent decisionEngine counts, optional providerMode).
// kubeVirtLegacy.present must be true explicitly in the fixture — never inferred from absence.
func ValidateMissingOptionalFieldsEdge(summary ParentFabricSummary, meta RedactedLiveSummaryMetadata) error {
	if summary.ExecutionEngine.ApplyAllowed {
		return fmt.Errorf("missing optional edge: applyAllowed must be false")
	}
	if summary.WindowsLane.Enabled {
		return fmt.Errorf("missing optional edge: windows lane must be disabled")
	}
	if meta.DashboardGeneratedAtUnavailable {
		if summary.GeneratedAt != MissingOptionalGeneratedAtDefault {
			return fmt.Errorf("missing optional edge: generatedAt must default to %q when Dashboard generatedAt is unavailable, got %q",
				MissingOptionalGeneratedAtDefault, summary.GeneratedAt)
		}
	}
	if summary.ParentPool.DonorCount < 0 || summary.ParentPool.ReceiverCount < 0 {
		return fmt.Errorf("missing optional edge: donor/receiver counts must be non-negative")
	}
	if meta.DashboardCountsAbsent {
		if summary.ParentPool.DonorCount != 0 || summary.ParentPool.ReceiverCount != 0 {
			return fmt.Errorf("missing optional edge: absent decisionEngine must map to donor/receiver 0, got %d/%d",
				summary.ParentPool.DonorCount, summary.ParentPool.ReceiverCount)
		}
	}
	if !summary.KubeVirtLegacy.Present {
		return fmt.Errorf("missing optional edge: kubeVirtLegacy.present must be true (M1-M7 anchor; explicit in fixture, not inferred when block absent)")
	}
	if err := ValidateSummary(summary); err != nil {
		return fmt.Errorf("missing optional edge: %w", err)
	}
	if err := ValidateNoForbiddenClaims(summary); err != nil {
		return fmt.Errorf("missing optional edge: %w", err)
	}
	if err := ValidateApplySemantics(summary, meta.DashboardSupportsApply); err != nil {
		return fmt.Errorf("missing optional edge: %w", err)
	}
	return nil
}
