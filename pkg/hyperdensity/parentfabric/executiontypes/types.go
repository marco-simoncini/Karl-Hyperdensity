package executiontypes

import (
	"encoding/json"
	"sort"
)

const (
	ExecutionTypesPackageVersion   = "v0.0.0-sprint46"
	ExecutionTypesSourceFile         = "pkg/server/hyperdensity_parent_fabric_execution_types.go"
	ExecutionTypesExtractionMode     = "copy-contract-no-runtime-wiring"
	ExecutionTypesSourceTypeCount      = 152
	ExecutionEngineTopLevelFieldCount  = 150
)

// HyperdensityCPUQuantity and HyperdensityMemoryQuantity are defined in Dashboard
// hyperdensity_parent_fabric.go; duplicated here so HyperdensityExecutionSummary is self-contained.
type HyperdensityCPUQuantity struct {
	Quantity   string `json:"quantity"`
	Millicores int64  `json:"millicores"`
}

type HyperdensityMemoryQuantity struct {
	Quantity string `json:"quantity"`
	Bytes    int64  `json:"bytes"`
}

// HyperdensityExecutionSummary is copied from Dashboard execution_types (pure DTO).
type HyperdensityExecutionSummary struct {
	DryRunOnlyObjects          int                        `json:"dryRunOnlyObjects"`
	ApplyReadyObjects          int                        `json:"applyReadyObjects"`
	BlockedByPolicyObjects     int                        `json:"blockedByPolicyObjects"`
	BlockedByConfidenceObjects int                        `json:"blockedByConfidenceObjects"`
	BlockedByGovernanceObjects int                        `json:"blockedByGovernanceObjects"`
	BlockedByRuntimeObjects    int                        `json:"blockedByRuntimeObjects"`
	AppliedObjects             int                        `json:"appliedObjects"`
	AppliedWithFreezeObjects   int                        `json:"appliedWithFreezeObjects"`
	RolledBackObjects          int                        `json:"rolledBackObjects"`
	StagedForRetryObjects      int                        `json:"stagedForRetryObjects"`
	PlannedReclaimCPU          HyperdensityCPUQuantity    `json:"plannedReclaimCpu"`
	PlannedReclaimMemory       HyperdensityMemoryQuantity `json:"plannedReclaimMemory"`
	PlannedBurstCPU            HyperdensityCPUQuantity    `json:"plannedBurstCpu"`
	PlannedBurstMemory         HyperdensityMemoryQuantity `json:"plannedBurstMemory"`
	RemainingBlockedCPU        HyperdensityCPUQuantity    `json:"remainingBlockedCpu"`
	RemainingBlockedMemory     HyperdensityMemoryQuantity `json:"remainingBlockedMemory"`
	WaitingApprovalCPU         HyperdensityCPUQuantity    `json:"waitingApprovalCpu"`
	WaitingApprovalMemory      HyperdensityMemoryQuantity `json:"waitingApprovalMemory"`
	DryRunOnlyCPU              HyperdensityCPUQuantity    `json:"dryRunOnlyCpu"`
	DryRunOnlyMemory           HyperdensityMemoryQuantity `json:"dryRunOnlyMemory"`
	Summary                    string                     `json:"summary"`
}

// HyperdensityExecutionEngineSpine is the contract subset of Dashboard HyperdensityExecutionEngine
// (Summary + apply posture fields only). Nested surface types are deferred.
type HyperdensityExecutionEngineSpine struct {
	Summary           HyperdensityExecutionSummary `json:"summary"`
	SupportsApply     bool                         `json:"supportsApply"`
	SupportedSurfaces []string                     `json:"supportedSurfaces,omitempty"`
	ApplyNotes        []string                     `json:"applyNotes,omitempty"`
}

// ExecutionEngineSpineJSONKeys are the copied top-level json keys (stable contract ordering).
var ExecutionEngineSpineJSONKeys = []string{
	"summary",
	"supportsApply",
	"supportedSurfaces",
	"applyNotes",
}

// CopiedTypeNames lists types present in this package from the Dashboard source file.
var CopiedTypeNames = []string{
	"HyperdensityCPUQuantity",
	"HyperdensityMemoryQuantity",
	"HyperdensityExecutionSummary",
	"HyperdensityExecutionEngineSpine",
}

// ContractDocument is the stable Sprint 46 contract snapshot (golden source).
type ContractDocument struct {
	ContractVersion   string                         `json:"contractVersion"`
	SourceFile        string                         `json:"sourceFile"`
	ExtractionMode    string                         `json:"extractionMode"`
	CopyScope         string                         `json:"copyScope"`
	SourceStats       ContractSourceStats            `json:"sourceStats"`
	CopiedTypes       []string                       `json:"copiedTypes"`
	DeferredNote      string                         `json:"deferredNote"`
	EngineSpineKeys   []string                       `json:"engineSpineKeys"`
	EngineSpine       HyperdensityExecutionEngineSpine `json:"engineSpine"`
	SummaryZeroValue  HyperdensityExecutionSummary   `json:"summaryZeroValue"`
}

type ContractSourceStats struct {
	LineCount                         int      `json:"lineCount"`
	TypeDefinitionCount               int      `json:"typeDefinitionCount"`
	SourceImports                     []string `json:"sourceImports"`
	ExecutionEngineTopLevelFieldCount int      `json:"executionEngineTopLevelFieldCount"`
}

// DefaultContractDocument returns the canonical zero-value contract document.
func DefaultContractDocument() ContractDocument {
	return ContractDocument{
		ContractVersion: ExecutionTypesPackageVersion,
		SourceFile:      ExecutionTypesSourceFile,
		ExtractionMode:  ExecutionTypesExtractionMode,
		CopyScope:       "partial-spine-and-summary",
		SourceStats: ContractSourceStats{
			LineCount:                         4571,
			TypeDefinitionCount:               ExecutionTypesSourceTypeCount,
			SourceImports:                     []string{"time"},
			ExecutionEngineTopLevelFieldCount: ExecutionEngineTopLevelFieldCount,
		},
		CopiedTypes:      append([]string(nil), CopiedTypeNames...),
		DeferredNote:     "Remaining nested HyperdensityExecutionEngine surface types (~146 top-level fields) stay in Dashboard until per-surface extraction sprints.",
		EngineSpineKeys:  append([]string(nil), ExecutionEngineSpineJSONKeys...),
		EngineSpine:      HyperdensityExecutionEngineSpine{},
		SummaryZeroValue: HyperdensityExecutionSummary{},
	}
}

// CanonicalContractJSON returns deterministic JSON for golden comparison (sorted keys at top level via struct order).
func CanonicalContractJSON() ([]byte, error) {
	doc := DefaultContractDocument()
	keys := append([]string(nil), doc.EngineSpineKeys...)
	sort.Strings(keys)
	doc.EngineSpineKeys = keys
	copied := append([]string(nil), doc.CopiedTypes...)
	sort.Strings(copied)
	doc.CopiedTypes = copied
	return json.Marshal(doc)
}
