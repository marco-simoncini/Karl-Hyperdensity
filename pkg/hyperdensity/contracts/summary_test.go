package contracts

import (
	"os"
	"path/filepath"
	"testing"
)

func moduleRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("go.mod not found")
		}
		dir = parent
	}
}

func TestGoldenParentFabricSummaryParse(t *testing.T) {
	root := moduleRoot(t)
	path := filepath.Join(root, "testdata", "dashboard", "parent_fabric_summary_redacted.golden.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	s, err := ParseParentFabricSummary(data)
	if err != nil {
		t.Fatal(err)
	}
	if s.APIVersion != SummaryAPIVersion {
		t.Fatalf("apiVersion=%q", s.APIVersion)
	}
	if err := ValidateSummary(s); err != nil {
		t.Fatalf("ValidateSummary: %v", err)
	}
	if err := ValidateNoForbiddenClaims(s); err != nil {
		t.Fatalf("ValidateNoForbiddenClaims: %v", err)
	}
}

func TestGoldenWindowsLaneDisabled(t *testing.T) {
	root := moduleRoot(t)
	path := filepath.Join(root, "testdata", "dashboard", "parent_fabric_summary_redacted.golden.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	s, err := ParseParentFabricSummary(data)
	if err != nil {
		t.Fatal(err)
	}
	if s.WindowsLane.Enabled {
		t.Fatal("windows lane must be disabled in golden")
	}
	foundNoWindows := false
	for _, b := range s.WindowsLane.Blockers {
		if b == "no_windows_lane" {
			foundNoWindows = true
		}
	}
	if !foundNoWindows {
		t.Fatal("golden must include no_windows_lane blocker")
	}
}

func TestValidateNoForbiddenClaimsRejectsWindowsEnabled(t *testing.T) {
	s := ParentFabricSummary{
		APIVersion:  SummaryAPIVersion,
		GeneratedAt: "2026-05-15T12:00:00Z",
		Source:      "test",
		ExecutionEngine: ExecutionEngineSummary{
			DryRunSupported: true,
			ApplyAllowed:    false,
		},
		WindowsLane: WindowsLaneSummary{Enabled: true},
		KubeVirtLegacy: KubeVirtLegacySummary{
			Present: true,
		},
		Hyperdensity: HyperdensityPosture{RecommendationOnly: true, OperatorControlled: true},
	}
	if err := ValidateNoForbiddenClaims(s); err == nil {
		t.Fatal("expected error when windows enabled")
	}
}

func TestValidateSummaryRequiresDryRun(t *testing.T) {
	s := ParentFabricSummary{
		APIVersion:  SummaryAPIVersion,
		GeneratedAt: "2026-05-15T12:00:00Z",
		Source:      "test",
		ExecutionEngine: ExecutionEngineSummary{
			DryRunSupported: false,
		},
		KubeVirtLegacy: KubeVirtLegacySummary{Present: true},
		Hyperdensity:   HyperdensityPosture{RecommendationOnly: true},
	}
	if err := ValidateSummary(s); err == nil {
		t.Fatal("expected error when dryRunSupported false")
	}
}
