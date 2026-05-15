package contracts

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseAndValidateFixtureManifest_contractkitExample(t *testing.T) {
	root := moduleRoot(t)
	path := filepath.Join(root, "testdata", "dashboard", "hyperdensity_parity_manifest_m1_m7.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	m, err := ParseFixtureManifest(data)
	if err != nil {
		t.Fatal(err)
	}
	if err := ValidateFixtureManifest(m); err != nil {
		t.Fatal(err)
	}
	if m.ContractKitVersion != ContractKitVersion {
		t.Fatalf("version mismatch: %q", m.ContractKitVersion)
	}
	set := ManifestCaseIDSet(m)
	expected := ExpectedM1M8CaseIDs()
	if len(set) != len(expected) {
		t.Fatalf("example manifest case count: got %d want %d", len(set), len(expected))
	}
	for _, id := range expected {
		if _, ok := set[id]; !ok {
			t.Fatalf("example manifest missing id %q", id)
		}
	}
}

func TestValidateFixtureManifest_rejectsVersionMismatch(t *testing.T) {
	m := FixtureManifest{
		ManifestVersion:    "1",
		ContractKitVersion: "v0.0.0-old",
		Cases: []FixtureCase{{
			ID:                     "x",
			Milestone:              "M1",
			DashboardFixture:       "a.json",
			ContractGolden:         "b.json",
			ClaimSafe:              true,
			KubeVirtLegacyRequired: true,
		}},
	}
	if err := ValidateFixtureManifest(m); err == nil {
		t.Fatal("expected version mismatch error")
	}
}

func TestValidateFixtureManifest_rejectsWindowsEnabled(t *testing.T) {
	m := FixtureManifest{
		ManifestVersion:    "1",
		ContractKitVersion: ContractKitVersion,
		Cases: []FixtureCase{{
			ID:                     "x",
			Milestone:              "M1",
			DashboardFixture:       "a.json",
			ContractGolden:         "b.json",
			WindowsEnabled:         true,
			ClaimSafe:              true,
			KubeVirtLegacyRequired: true,
		}},
	}
	if err := ValidateFixtureManifest(m); err == nil {
		t.Fatal("expected windows enabled error")
	}
}

func TestVersion_returnsContractKitVersion(t *testing.T) {
	if Version() != ContractKitVersion {
		t.Fatalf("Version()=%q", Version())
	}
}
