package claimpolicy

import (
	"sort"
	"testing"
)

func TestValidateSurfaceMappings(t *testing.T) {
	if err := ValidateSurfaceMappings(); err != nil {
		t.Fatal(err)
	}
}

func TestSurfaceMappingsEveryClaimKnown(t *testing.T) {
	for _, m := range SurfaceMappings() {
		if !Known(m.ClaimID) {
			t.Fatalf("mapping claim id %q must be Known", m.ClaimID)
		}
		if m.RuntimeImportAllowed {
			t.Fatalf("mapping %q/%q must not allow runtime import in Sprint 37–39", m.ClaimID, m.Surface)
		}
	}
}

func TestCriticalClaimsHaveMapping(t *testing.T) {
	critical := []string{
		string(ClaimNoProductionMutation),
		string(ClaimNoAutonomousApply),
		string(ClaimNoWindowsHyperdensityApply),
		string(ClaimNoRuntimeContractsImport),
	}
	for _, id := range critical {
		if len(MappingsForClaim(id)) == 0 {
			t.Fatalf("critical claim %q must have at least one mapping", id)
		}
		if !MustKeepRuntimeDisabled(id) {
			t.Fatalf("test expects %q to be a critical MustKeepRuntimeDisabled claim", id)
		}
	}
}

func TestKubeVirtLegacyDistinctFromReplacement(t *testing.T) {
	var legacyField, replField string
	for _, m := range SurfaceMappings() {
		if m.ClaimID == string(ClaimKubeVirtLegacyProvider) {
			legacyField = m.Field
		}
		if m.ClaimID == string(ClaimNoGenericKubeVirtReplacement) {
			replField = m.Field
		}
	}
	if legacyField == "" || replField == "" {
		t.Fatal("expected mappings for kubevirt legacy and generic replacement")
	}
	if legacyField == replField {
		t.Fatalf("legacy vs replacement mapping must use distinct fields, both were %q", legacyField)
	}
	if string(ClaimKubeVirtLegacyProvider) == string(ClaimNoGenericKubeVirtReplacement) {
		t.Fatal("claim ids must differ")
	}
}

func TestNoWindowsMappedToWindowsLane(t *testing.T) {
	found := false
	for _, m := range SurfaceMappings() {
		if m.ClaimID != string(ClaimNoWindowsHyperdensityApply) {
			continue
		}
		if m.Surface == string(SurfaceWindowsLane) && m.Field == "windows_hyperdensity_apply" && m.ExpectedValue == "disabled" {
			found = true
		}
	}
	if !found {
		t.Fatal("expected no_windows_hyperdensity_apply on windows_lane with apply disabled")
	}
}

func TestNoRuntimeContractsMappedToFreeze(t *testing.T) {
	found := false
	for _, m := range SurfaceMappings() {
		if m.ClaimID == string(ClaimNoRuntimeContractsImport) && m.Surface == string(SurfaceRuntimeImportFreeze) {
			found = true
			if m.ExpectedValue != "contractkit_blockers_only" {
				t.Fatalf("unexpected ExpectedValue: %q", m.ExpectedValue)
			}
		}
	}
	if !found {
		t.Fatal("expected no_runtime_contracts_import mapped to runtime_import_freeze")
	}
}

func TestSurfaceMappingsStableOrder(t *testing.T) {
	got := SurfaceMappings()
	if len(got) < 2 {
		t.Fatal("expected multiple mappings")
	}
	ids := make([]string, len(got))
	for i := range got {
		ids[i] = got[i].ClaimID + "\x00" + got[i].Surface + "\x00" + got[i].Field + "\x00" + got[i].ExpectedValue
	}
	if !sort.StringsAreSorted(ids) {
		t.Fatalf("surface mappings not in stable sorted order: first keys %#v", ids[:min(5, len(ids))])
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func TestDryRunAndRecommendationMappings(t *testing.T) {
	var dry, rec bool
	for _, m := range SurfaceMappings() {
		if m.ClaimID == string(ClaimDryRunOnly) && m.Surface == string(SurfaceExecutionEngine) && m.Field == "execution_category" {
			dry = true
		}
		if m.ClaimID == string(ClaimRecommendationOnly) && m.Surface == string(SurfaceHyperdensityRecommendation) {
			rec = true
		}
	}
	if !dry || !rec {
		t.Fatalf("dry_run_only=%v recommendation_only=%v", dry, rec)
	}
}
