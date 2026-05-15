package contracts

import "strings"

// Module semver (go.mod / git tag) — distinct from logical contract schema version.
const ContractKitModuleVersion = "v0.1.9-khr-m1-m19"

// Git tag on Karl-Hyperdensity parent repo (nested module prefix required for go get).
const ContractKitGitTag = "pkg/hyperdensity/contractkit/v0.1.9-khr-m1-m19"

// ContractKitCurrentStableModuleVersion is the consumer-facing stable module semver (always equal to ContractKitModuleVersion).
const ContractKitCurrentStableModuleVersion = ContractKitModuleVersion

// ContractKitSupersededModuleVersions lists published module versions that must not be re-pinned or re-tagged; publish a strictly newer semver instead.
var ContractKitSupersededModuleVersions = []string{
	"v0.1.5-khr-m1-m16",
	"v0.1.7-khr-m1-m18",
}

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

// CurrentStableReleaseInfo returns the same metadata as ReleaseInfo (Sprint 40: explicit stable alias for parity/docs).
func CurrentStableReleaseInfo() ReleaseMetadata {
	return ReleaseInfo()
}

// IsSupersededModuleVersion reports whether version is a known superseded nested-module semver (trimmed exact match).
func IsSupersededModuleVersion(version string) bool {
	v := strings.TrimSpace(version)
	for _, s := range ContractKitSupersededModuleVersions {
		if v == s {
			return true
		}
	}
	return false
}
