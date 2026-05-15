package contracts

import (
	"strings"
	"testing"
)

func TestAllowedDashboardSummaryFixtureFields_stable(t *testing.T) {
	got := AllowedDashboardSummaryFixtureFields()
	if len(got) == 0 {
		t.Fatal("expected non-empty allowlist")
	}
	seen := make(map[string]struct{}, len(got))
	for _, f := range got {
		if strings.TrimSpace(f) == "" {
			t.Fatalf("empty allowlist entry")
		}
		if _, ok := seen[f]; ok {
			t.Fatalf("duplicate allowlist entry: %q", f)
		}
		seen[f] = struct{}{}
	}
	if len(seen) != len(got) {
		t.Fatal("internal: slice length mismatch")
	}
}

func TestValidateContractClaimSafe_ok(t *testing.T) {
	s := ParentFabricSummary{
		APIVersion:  SummaryAPIVersion,
		GeneratedAt: "2026-05-15T12:00:00Z",
		Source:      "live-capture-redacted-x",
		ExecutionEngine: BuildClaimSafeExecutionEngine(false, "dry_run_only", false),
		WindowsLane:    WindowsLaneSummary{Enabled: false},
		KubeVirtLegacy: KubeVirtLegacySummary{Present: true, ProviderMode: "kubevirt.legacy.v1"},
		Hyperdensity: HyperdensityPosture{
			RecommendationOnly: true,
			OperatorControlled: true,
		},
	}
	meta := RedactedLiveSummaryMetadata{
		DashboardSupportsApply:   false,
		ExecutionSummaryCategory: "dry_run_only",
	}
	if err := ValidateContractClaimSafe(s, meta); err != nil {
		t.Fatal(err)
	}
}

func TestValidateContractClaimSafe_applyAllowed(t *testing.T) {
	s := ParentFabricSummary{
		APIVersion:  SummaryAPIVersion,
		GeneratedAt: "2026-05-15T12:00:00Z",
		Source:      "live-capture-redacted-x",
		ExecutionEngine: ExecutionEngineSummary{
			Mode:            "operator_controlled",
			ApplyAllowed:    true,
			DryRunSupported: true,
		},
		WindowsLane:    WindowsLaneSummary{Enabled: false},
		KubeVirtLegacy: KubeVirtLegacySummary{Present: true},
		Hyperdensity: HyperdensityPosture{
			RecommendationOnly: true,
			OperatorControlled: true,
		},
	}
	if err := ValidateContractClaimSafe(s, RedactedLiveSummaryMetadata{}); err == nil {
		t.Fatal("expected error when applyAllowed true")
	}
}

func TestValidateSupportsApplyFalseEdge_ok(t *testing.T) {
	s := ParentFabricSummary{
		APIVersion:  SummaryAPIVersion,
		GeneratedAt: "2026-05-15T12:00:00Z",
		Source:      "live-capture-redacted-x",
		ExecutionEngine: BuildClaimSafeExecutionEngine(false, "dry_run_only", false),
		WindowsLane:    WindowsLaneSummary{Enabled: false},
		KubeVirtLegacy: KubeVirtLegacySummary{Present: true},
		Hyperdensity: HyperdensityPosture{
			RecommendationOnly: true,
			OperatorControlled: true,
		},
	}
	meta := RedactedLiveSummaryMetadata{
		DashboardSupportsApply:   false,
		ExecutionSummaryCategory: "dry_run_only",
	}
	if err := ValidateSupportsApplyFalseEdge(s, meta); err != nil {
		t.Fatal(err)
	}
}

func TestValidateSupportsApplyFalseEdge_wrongCategoryWhenDryRunSupported(t *testing.T) {
	s := ParentFabricSummary{
		APIVersion:  SummaryAPIVersion,
		GeneratedAt: "2026-05-15T12:00:00Z",
		Source:      "live-capture-redacted-x",
		ExecutionEngine: ExecutionEngineSummary{
			Mode:            "operator_controlled",
			ApplyAllowed:    false,
			DryRunSupported: true,
		},
		WindowsLane:    WindowsLaneSummary{Enabled: false},
		KubeVirtLegacy: KubeVirtLegacySummary{Present: true},
		Hyperdensity: HyperdensityPosture{
			RecommendationOnly: true,
			OperatorControlled: true,
		},
	}
	meta := RedactedLiveSummaryMetadata{
		DashboardSupportsApply:   false,
		ExecutionSummaryCategory: "live",
	}
	if err := ValidateSupportsApplyFalseEdge(s, meta); err == nil {
		t.Fatal("expected error when dryRunSupported but category not dry_run_only")
	}
}

func TestValidateSupportsApplyFalseEdge_rejectsSupportsApplyTrueMeta(t *testing.T) {
	s := ParentFabricSummary{
		APIVersion:  SummaryAPIVersion,
		GeneratedAt: "2026-05-15T12:00:00Z",
		Source:      "live-capture-redacted-x",
		ExecutionEngine: BuildClaimSafeExecutionEngine(true, "dry_run_only", false),
		WindowsLane:    WindowsLaneSummary{Enabled: false},
		KubeVirtLegacy: KubeVirtLegacySummary{Present: true},
		Hyperdensity: HyperdensityPosture{
			RecommendationOnly: true,
			OperatorControlled: true,
		},
	}
	meta := RedactedLiveSummaryMetadata{DashboardSupportsApply: true}
	if err := ValidateSupportsApplyFalseEdge(s, meta); err == nil {
		t.Fatal("expected error when meta.DashboardSupportsApply true")
	}
}
