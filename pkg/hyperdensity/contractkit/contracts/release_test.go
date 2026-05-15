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

func TestCurrentStableReleaseInfo_matchesModuleVersion(t *testing.T) {
	stable := CurrentStableReleaseInfo()
	if stable.ModuleVersion != ContractKitModuleVersion {
		t.Fatalf("CurrentStableReleaseInfo moduleVersion=%q want %q", stable.ModuleVersion, ContractKitModuleVersion)
	}
	if stable.GitTag != ContractKitGitTag {
		t.Fatalf("CurrentStableReleaseInfo gitTag=%q", stable.GitTag)
	}
}

func TestContractKitCurrentStableModuleVersion_alias(t *testing.T) {
	if ContractKitCurrentStableModuleVersion != ContractKitModuleVersion {
		t.Fatalf("ContractKitCurrentStableModuleVersion=%q ContractKitModuleVersion=%q", ContractKitCurrentStableModuleVersion, ContractKitModuleVersion)
	}
}

func TestIsSupersededModuleVersion(t *testing.T) {
	if !IsSupersededModuleVersion("v0.1.5-khr-m1-m16") {
		t.Fatal("expected v0.1.5-khr-m1-m16 superseded")
	}
	if !IsSupersededModuleVersion("v0.1.7-khr-m1-m18") {
		t.Fatal("expected v0.1.7-khr-m1-m18 superseded")
	}
	if IsSupersededModuleVersion(ContractKitModuleVersion) {
		t.Fatalf("current module version must not be superseded: %q", ContractKitModuleVersion)
	}
	if IsSupersededModuleVersion(" " + ContractKitModuleVersion + "  ") {
		t.Fatal("trimmed current version must not match superseded list")
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
