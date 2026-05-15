package contracts

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCanonicalSummaryJSON_roundTrip(t *testing.T) {
	s := ParentFabricSummary{
		APIVersion:  SummaryAPIVersion,
		GeneratedAt: "2026-05-15T12:00:00Z",
		Source:      "test",
		ParentPool: ParentPoolSummary{
			UsageSummary:  "u",
			DonorCount:    1,
			ReceiverCount: 2,
		},
		ExecutionEngine: ExecutionEngineSummary{
			Mode:            "operator_controlled",
			AutonomousMode:  false,
			ApplyAllowed:    false,
			DryRunSupported: true,
		},
		WindowsLane: WindowsLaneSummary{
			Enabled:  false,
			Reason:   "r",
			Blockers: []string{"no_windows_lane"},
		},
		KubeVirtLegacy: KubeVirtLegacySummary{
			Present:      true,
			ProviderMode: "kubevirt.legacy.v1",
		},
		Hyperdensity: HyperdensityPosture{
			RecommendationOnly: true,
			OperatorControlled: true,
		},
	}
	b, err := CanonicalSummaryJSON(s)
	if err != nil {
		t.Fatal(err)
	}
	if len(b) > 0 && b[len(b)-1] == '\n' {
		t.Fatal("canonical json must not end with newline")
	}
	parsed, err := ParseParentFabricSummary(b)
	if err != nil {
		t.Fatal(err)
	}
	b2, err := CanonicalSummaryJSON(parsed)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != string(b2) {
		t.Fatalf("second canonicalization changed output:\n%s\nvs\n%s", b, b2)
	}
}

func TestCompareSummaryGolden_match(t *testing.T) {
	s := ParentFabricSummary{
		APIVersion:  SummaryAPIVersion,
		GeneratedAt: "2026-05-15T12:00:00Z",
		Source:      "src",
		ParentPool: ParentPoolSummary{
			DonorCount:    3,
			ReceiverCount: 2,
		},
		ExecutionEngine: ExecutionEngineSummary{
			Mode:            "operator_controlled",
			DryRunSupported: true,
		},
		WindowsLane: WindowsLaneSummary{
			Enabled: false,
		},
		KubeVirtLegacy: KubeVirtLegacySummary{Present: true},
		Hyperdensity: HyperdensityPosture{
			RecommendationOnly: true,
			OperatorControlled: true,
		},
	}
	canon, err := CanonicalSummaryJSON(s)
	if err != nil {
		t.Fatal(err)
	}
	// Golden with extra outer whitespace should still match after parse + canonicalize.
	golden := append([]byte("\n\t "), append(canon, []byte(" \n")...)...)
	if err := CompareSummaryGolden(s, golden); err != nil {
		t.Fatal(err)
	}
}

func TestCompareSummaryGolden_mismatch(t *testing.T) {
	a := ParentFabricSummary{
		APIVersion:  SummaryAPIVersion,
		GeneratedAt: "2026-05-15T12:00:00Z",
		Source:      "a",
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
	b := a
	b.Source = "b"
	golden, err := CanonicalSummaryJSON(b)
	if err != nil {
		t.Fatal(err)
	}
	if err := CompareSummaryGolden(a, golden); err == nil {
		t.Fatal("expected mismatch")
	}
}

func TestWriteCanonicalSummary(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.json")
	s := ParentFabricSummary{
		APIVersion:  SummaryAPIVersion,
		GeneratedAt: "2026-05-15T12:00:00Z",
		Source:      "w",
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
	if err := WriteCanonicalSummary(path, s); err != nil {
		t.Fatal(err)
	}
	if err := CompareSummaryGolden(s, mustRead(t, path)); err != nil {
		t.Fatal(err)
	}
}

func mustRead(t *testing.T, path string) []byte {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return b
}
