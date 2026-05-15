package contracts

// ContractKitVersion is the logical contract schema epoch (manifest contractKitVersion field).
// Distinct from module semver — see ContractKitModuleVersion and ReleaseInfo().
const ContractKitVersion = "v0.0.0-sprint26"

// ContractKitCommitHint documents consumer go.mod pin (semver tag or pseudo-version).
const ContractKitCommitHint = "consumer-pinned; see ContractKitGitTag"

// Version returns ContractKitVersion for test and manifest guards.
func Version() string {
	return ContractKitVersion
}
