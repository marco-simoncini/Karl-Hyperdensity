package contracts

// ContractKitVersion is the semantic version string for this contractkit module release.
// Consumers (e.g. Karl-Dashboard parity tests) should assert manifest contractKitVersion matches.
const ContractKitVersion = "v0.0.0-sprint26"

// ContractKitCommitHint documents that the authoritative module revision is consumer-pinned
// via go.mod pseudo-version (not embedded at runtime). Update docs when tagging releases.
const ContractKitCommitHint = "consumer-pinned"

// Version returns ContractKitVersion for test and manifest guards.
func Version() string {
	return ContractKitVersion
}
