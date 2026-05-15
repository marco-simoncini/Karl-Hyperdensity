package claimpolicy

import (
	"sort"
	"testing"
)

func TestPackageVersion(t *testing.T) {
	if PackageVersion == "" {
		t.Fatal("PackageVersion must be set")
	}
}

func TestKnownPosture(t *testing.T) {
	for _, p := range Postures() {
		if !KnownPosture(p) {
			t.Fatalf("Postures() returned unknown posture %q", p)
		}
	}
	if KnownPosture(PostureKind("unknown_posture")) {
		t.Fatal("unknown posture must not be known")
	}
}

func TestPosturesStableOrder(t *testing.T) {
	got := Postures()
	if len(got) != 3 {
		t.Fatalf("len(Postures)=%d", len(got))
	}
	if got[0] != PostureEvidenceNamespace || got[1] != PostureOperatorControlled || got[2] != PostureVisibilityOnly {
		t.Fatalf("unexpected order: %#v", got)
	}
}

func TestCatalogNonEmptyUniqueStable(t *testing.T) {
	cat := Catalog()
	if len(cat) == 0 {
		t.Fatal("Catalog() must be non-empty")
	}
	seen := make(map[string]struct{}, len(cat))
	var ids []string
	for i := 1; i < len(cat); i++ {
		if cat[i-1].ID >= cat[i].ID {
			t.Fatalf("catalog not strictly sorted by ID at %d: %q >= %q", i, cat[i-1].ID, cat[i].ID)
		}
	}
	for _, r := range cat {
		if r.ID == "" {
			t.Fatal("empty rule ID")
		}
		if _, dup := seen[r.ID]; dup {
			t.Fatalf("duplicate id %q", r.ID)
		}
		seen[r.ID] = struct{}{}
		ids = append(ids, r.ID)
	}
	sort.Strings(ids)
	for _, id := range ids {
		if !Known(id) {
			t.Fatalf("Known(%q) must be true for catalog id", id)
		}
		if Severity(id) == "" {
			t.Fatalf("Severity(%q) must be non-empty", id)
		}
	}
}

func TestForbiddenProductionClaimIDs(t *testing.T) {
	got := ForbiddenProductionClaimIDs()
	if len(got) == 0 {
		t.Fatal("ForbiddenProductionClaimIDs must be non-empty")
	}
	if !sort.StringsAreSorted(got) {
		t.Fatalf("ForbiddenProductionClaimIDs must be sorted: %#v", got)
	}
	for _, id := range got {
		if !Known(id) {
			t.Fatalf("forbidden id %q must be known", id)
		}
		if RuntimeAllowed(id) {
			t.Fatalf("RuntimeAllowed(%q) must be false for forbidden production claim", id)
		}
	}
	for _, id := range got {
		if id == string(ClaimKubeVirtLegacyProvider) {
			t.Fatal("kubevirt_legacy_provider must not appear in ForbiddenProductionClaimIDs")
		}
	}
}

func TestMustKeepRuntimeDisabledCriticalClaims(t *testing.T) {
	cases := []ClaimPolicyID{
		ClaimNoProductionMutation,
		ClaimNoAutonomousApply,
		ClaimNoWindowsHyperdensityApply,
		ClaimNoRuntimeContractsImport,
	}
	for _, c := range cases {
		if !MustKeepRuntimeDisabled(string(c)) {
			t.Fatalf("MustKeepRuntimeDisabled(%q) want true", c)
		}
	}
	if MustKeepRuntimeDisabled(string(ClaimNoGenericKubeVirtReplacement)) {
		t.Fatal("MustKeepRuntimeDisabled must be false for no_generic_kubevirt_replacement (distinct from legacy marker)")
	}
}

func TestKubeVirtLegacyNotReplacement(t *testing.T) {
	if !Known(string(ClaimKubeVirtLegacyProvider)) {
		t.Fatal("legacy provider must be Known")
	}
	if string(ClaimKubeVirtLegacyProvider) == string(ClaimNoGenericKubeVirtReplacement) {
		t.Fatal("legacy provider id must differ from generic replacement claim")
	}
	if !RuntimeAllowed(string(ClaimKubeVirtLegacyProvider)) {
		t.Fatal("kubevirt legacy marker should be RuntimeAllowed true (compatibility, not replacement narrative)")
	}
	if RuntimeAllowed(string(ClaimNoGenericKubeVirtReplacement)) {
		t.Fatal("generic replacement claim must remain RuntimeAllowed false")
	}
}
