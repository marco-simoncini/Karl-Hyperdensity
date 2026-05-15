package contracts

import "testing"

func TestValidateRedactedLiveSummaryFixture_ok(t *testing.T) {
	s := ParentFabricSummary{
		APIVersion:  SummaryAPIVersion,
		GeneratedAt: "2026-05-15T12:00:00Z",
		Source:      "karl-hyperdensity-live-capture-redacted-m5",
		ParentPool: ParentPoolSummary{
			DonorCount:    2,
			ReceiverCount: 1,
		},
		ExecutionEngine: BuildClaimSafeExecutionEngine(true, "dry_run_only", false),
		WindowsLane: WindowsLaneSummary{
			Enabled: false,
		},
		KubeVirtLegacy: KubeVirtLegacySummary{Present: true, ProviderMode: "kubevirt.legacy.v1"},
		Hyperdensity: HyperdensityPosture{
			RecommendationOnly: true,
			OperatorControlled: true,
		},
	}
	meta := RedactedLiveSummaryMetadata{DashboardSupportsApply: true, ExecutionSummaryCategory: "dry_run_only"}
	if err := ValidateRedactedLiveSummaryFixture(s, meta); err != nil {
		t.Fatal(err)
	}
}

func TestValidateRedactedLiveSummaryFixture_sourceMarker(t *testing.T) {
	s := ParentFabricSummary{
		APIVersion:  SummaryAPIVersion,
		GeneratedAt: "2026-05-15T12:00:00Z",
		Source:      "prod-live-no-redaction-marker",
		ExecutionEngine: ExecutionEngineSummary{
			Mode:            "operator_controlled",
			DryRunSupported: true,
		},
		KubeVirtLegacy: KubeVirtLegacySummary{Present: true},
		Hyperdensity: HyperdensityPosture{
			RecommendationOnly: true,
			OperatorControlled: true,
		},
	}
	if err := ValidateRedactedLiveSummaryFixture(s, RedactedLiveSummaryMetadata{}); err == nil {
		t.Fatal("expected error for missing redaction marker in source")
	}
}

func TestValidateRedactedLiveSummaryFixture_redactedSubstring(t *testing.T) {
	s := ParentFabricSummary{
		APIVersion:  SummaryAPIVersion,
		GeneratedAt: "2026-05-15T12:00:00Z",
		Source:      "fixture-redacted-manual",
		ExecutionEngine: ExecutionEngineSummary{
			Mode:            "operator_controlled",
			DryRunSupported: true,
		},
		WindowsLane:    WindowsLaneSummary{Enabled: false},
		KubeVirtLegacy: KubeVirtLegacySummary{Present: true},
		Hyperdensity: HyperdensityPosture{
			RecommendationOnly: true,
			OperatorControlled: true,
		},
	}
	if err := ValidateRedactedLiveSummaryFixture(s, RedactedLiveSummaryMetadata{}); err != nil {
		t.Fatal(err)
	}
}

func TestValidateRedactedLiveSummaryFixture_windowsEnabled(t *testing.T) {
	s := ParentFabricSummary{
		APIVersion:  SummaryAPIVersion,
		GeneratedAt: "2026-05-15T12:00:00Z",
		Source:      "live-capture-redacted-x",
		ExecutionEngine: ExecutionEngineSummary{
			Mode:            "operator_controlled",
			DryRunSupported: true,
		},
		WindowsLane:    WindowsLaneSummary{Enabled: true},
		KubeVirtLegacy: KubeVirtLegacySummary{Present: true},
		Hyperdensity: HyperdensityPosture{
			RecommendationOnly: true,
			OperatorControlled: true,
		},
	}
	if err := ValidateRedactedLiveSummaryFixture(s, RedactedLiveSummaryMetadata{}); err == nil {
		t.Fatal("expected error when windows enabled")
	}
}
