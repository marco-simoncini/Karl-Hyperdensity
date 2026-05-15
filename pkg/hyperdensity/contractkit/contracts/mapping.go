package contracts

import (
	"fmt"
	"strings"
)

// MapSupportsApplyToContractApplyAllowed maps Dashboard executionEngine.supportsApply
// to Hyperdensity contract applyAllowed. Technical capability must not become a
// product apply claim: this always returns false for M1/M3 claim-safe anchors.
func MapSupportsApplyToContractApplyAllowed(supportsApply bool) bool {
	_ = supportsApply
	return false
}

// IsClaimSafeApplyAllowed reports whether applyAllowed is consistent with operator control.
func IsClaimSafeApplyAllowed(applyAllowed, operatorControlled bool) bool {
	if !applyAllowed {
		return true
	}
	return operatorControlled
}

// InferDryRunSupported derives dryRunSupported from Dashboard execution summary signals.
func InferDryRunSupported(summaryCategory string, supportsApply bool) bool {
	if strings.EqualFold(strings.TrimSpace(summaryCategory), "dry_run_only") {
		return true
	}
	return supportsApply
}

// NormalizeExecutionMode returns a stable lowercase execution mode label.
func NormalizeExecutionMode(mode string) string {
	m := strings.TrimSpace(strings.ToLower(mode))
	if m == "" {
		return "operator_controlled"
	}
	return m
}

// BuildClaimSafeExecutionEngine builds contract executionEngine fields from Dashboard summary inputs.
func BuildClaimSafeExecutionEngine(supportsApply bool, summaryCategory string, autonomousMode bool) ExecutionEngineSummary {
	return ExecutionEngineSummary{
		Mode:            NormalizeExecutionMode("operator_controlled"),
		ApplyAllowed:    MapSupportsApplyToContractApplyAllowed(supportsApply),
		DryRunSupported: InferDryRunSupported(summaryCategory, supportsApply),
		AutonomousMode:  autonomousMode,
	}
}

// ValidateApplySemantics ensures Dashboard supportsApply does not leak into contract applyAllowed.
func ValidateApplySemantics(contract ParentFabricSummary, dashboardSupportsApply bool) error {
	if dashboardSupportsApply && contract.ExecutionEngine.ApplyAllowed {
		return fmt.Errorf("dashboard supportsApply must not map to contract applyAllowed true")
	}
	if !IsClaimSafeApplyAllowed(contract.ExecutionEngine.ApplyAllowed, contract.Hyperdensity.OperatorControlled) {
		return fmt.Errorf("contract applyAllowed requires operatorControlled posture")
	}
	return nil
}
