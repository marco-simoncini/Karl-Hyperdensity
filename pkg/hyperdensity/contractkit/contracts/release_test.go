package contracts

import "testing"

func TestReleaseInfo_lockedVersions(t *testing.T) {
	info := ReleaseInfo()
	if info.ModuleVersion != ContractKitModuleVersion {
		t.Fatalf("moduleVersion=%q", info.ModuleVersion)
	}
	if info.GitTag != ContractKitGitTag {
		t.Fatalf("gitTag=%q", info.GitTag)
	}
	if info.SchemaVersion != ContractKitVersion {
		t.Fatalf("schemaVersion=%q", info.SchemaVersion)
	}
	if info.SchemaVersion != ContractKitSchemaVersion {
		t.Fatalf("ContractKitSchemaVersion mismatch")
	}
	if info.FixtureManifestVersion != FixtureManifestVersion {
		t.Fatalf("fixtureManifestVersion=%q", info.FixtureManifestVersion)
	}
}

func TestValidateFixtureManifest_rejectsWrongManifestVersion(t *testing.T) {
	m := FixtureManifest{
		ManifestVersion:    "hyperdensity.parity.manifest/v0",
		ContractKitVersion: ContractKitVersion,
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
		t.Fatal("expected manifestVersion mismatch error")
	}
}
