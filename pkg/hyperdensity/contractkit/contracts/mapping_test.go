package contracts

import "testing"

func TestMapSupportsApplyNeverEnablesContractApply(t *testing.T) {
	if MapSupportsApplyToContractApplyAllowed(true) {
		t.Fatal("supportsApply true must not yield applyAllowed true")
	}
	if MapSupportsApplyToContractApplyAllowed(false) {
		t.Fatal("supportsApply false must yield applyAllowed false")
	}
}

func TestValidateApplySemanticsRejectsSupportsApplyLeak(t *testing.T) {
	s := ParentFabricSummary{
		ExecutionEngine: ExecutionEngineSummary{ApplyAllowed: true},
		Hyperdensity:    HyperdensityPosture{OperatorControlled: true},
	}
	if err := ValidateApplySemantics(s, true); err == nil {
		t.Fatal("expected error when applyAllowed true with supportsApply true")
	}
}

func TestBuildClaimSafeExecutionEngineDryRunOnly(t *testing.T) {
	eng := BuildClaimSafeExecutionEngine(true, "dry_run_only", false)
	if eng.ApplyAllowed {
		t.Fatal("applyAllowed must be false")
	}
	if !eng.DryRunSupported {
		t.Fatal("dryRunSupported expected true for dry_run_only category")
	}
}

func TestNormalizeExecutionModeDefault(t *testing.T) {
	if NormalizeExecutionMode("") != "operator_controlled" {
		t.Fatal("expected default operator_controlled")
	}
}
