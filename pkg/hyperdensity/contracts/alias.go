// Package contracts re-exports contractkit summary DTOs and validators for in-repo callers.
// Prefer importing github.com/marco-simoncini/Karl-Hyperdensity/pkg/hyperdensity/contractkit/contracts externally.
package contracts

import ck "github.com/marco-simoncini/Karl-Hyperdensity/pkg/hyperdensity/contractkit/contracts"

const (
	SummaryAPIVersion                 = ck.SummaryAPIVersion
	MissingOptionalGeneratedAtDefault = ck.MissingOptionalGeneratedAtDefault
	ContractKitVersion              = ck.ContractKitVersion
	ContractKitCommitHint           = ck.ContractKitCommitHint
	ContractKitModuleVersion        = ck.ContractKitModuleVersion
	ContractKitGitTag               = ck.ContractKitGitTag
	ContractKitSchemaVersion        = ck.ContractKitSchemaVersion
	FixtureManifestVersion          = ck.FixtureManifestVersion
)

type (
	ReleaseMetadata          = ck.ReleaseMetadata
	FixtureManifest          = ck.FixtureManifest
	FixtureCase              = ck.FixtureCase
	ParentFabricSummary      = ck.ParentFabricSummary
	ParentPoolSummary        = ck.ParentPoolSummary
	ExecutionEngineSummary   = ck.ExecutionEngineSummary
	WindowsLaneSummary       = ck.WindowsLaneSummary
	KubeVirtLegacySummary    = ck.KubeVirtLegacySummary
	HyperdensityPosture      = ck.HyperdensityPosture
	RedactedLiveSummaryMetadata = ck.RedactedLiveSummaryMetadata
)

var (
	ReleaseInfo                        = ck.ReleaseInfo
	Version                            = ck.Version
	ExpectedM1M8CaseIDs                = ck.ExpectedM1M8CaseIDs
	CaseIDs                            = ck.CaseIDs
	ManifestCaseIDSet                  = ck.ManifestCaseIDSet
	ParseFixtureManifest               = ck.ParseFixtureManifest
	ValidateFixtureManifest            = ck.ValidateFixtureManifest
	ParseParentFabricSummary           = ck.ParseParentFabricSummary
	ValidateSummary                    = ck.ValidateSummary
	ValidateNoForbiddenClaims          = ck.ValidateNoForbiddenClaims
	MapSupportsApplyToContractApplyAllowed = ck.MapSupportsApplyToContractApplyAllowed
	IsClaimSafeApplyAllowed            = ck.IsClaimSafeApplyAllowed
	InferDryRunSupported               = ck.InferDryRunSupported
	NormalizeExecutionMode             = ck.NormalizeExecutionMode
	BuildClaimSafeExecutionEngine      = ck.BuildClaimSafeExecutionEngine
	ValidateApplySemantics             = ck.ValidateApplySemantics
	CanonicalSummaryJSON               = ck.CanonicalSummaryJSON
	WriteCanonicalSummary              = ck.WriteCanonicalSummary
	CompareSummaryGolden               = ck.CompareSummaryGolden
	ValidateRedactedLiveSummaryFixture = ck.ValidateRedactedLiveSummaryFixture
	AllowedDashboardSummaryFixtureFields = ck.AllowedDashboardSummaryFixtureFields
	ValidateContractClaimSafe          = ck.ValidateContractClaimSafe
	ValidateSupportsApplyFalseEdge     = ck.ValidateSupportsApplyFalseEdge
	ValidateMissingOptionalFieldsEdge  = ck.ValidateMissingOptionalFieldsEdge
)
