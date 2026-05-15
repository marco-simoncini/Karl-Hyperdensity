package contracts

// Module semver (go.mod / git tag) — distinct from logical contract schema version.
const ContractKitModuleVersion = "v0.1.0-khr-m1-m9"

// Git tag on Karl-Hyperdensity parent repo (nested module prefix required for go get).
const ContractKitGitTag = "pkg/hyperdensity/contractkit/v0.1.0-khr-m1-m9"

// ContractKitSchemaVersion is the logical contract/manifest schema epoch (ContractKitVersion).
const ContractKitSchemaVersion = ContractKitVersion

// FixtureManifestVersion is the JSON manifest envelope version (manifestVersion field).
const FixtureManifestVersion = "hyperdensity.parity.manifest/v1"

// ReleaseMetadata documents the three version layers for test-only parity tooling.
type ReleaseMetadata struct {
	ModuleVersion          string `json:"moduleVersion"`
	GitTag                 string `json:"gitTag"`
	SchemaVersion          string `json:"schemaVersion"`
	FixtureManifestVersion string `json:"fixtureManifestVersion"`
}

// ReleaseInfo returns the locked release metadata for M1–M9 parity anchors.
func ReleaseInfo() ReleaseMetadata {
	return ReleaseMetadata{
		ModuleVersion:          ContractKitModuleVersion,
		GitTag:                 ContractKitGitTag,
		SchemaVersion:          ContractKitSchemaVersion,
		FixtureManifestVersion: FixtureManifestVersion,
	}
}
