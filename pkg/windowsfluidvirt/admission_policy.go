package windowsfluidvirt

type AdmissionPhase string

const (
	AdmissionAdmittedForFutureApply AdmissionPhase = "ADMITTED_FOR_FUTURE_APPLY"
	AdmissionDenied                 AdmissionPhase = "DENIED"
	AdmissionBlocked                AdmissionPhase = "BLOCKED"
	AdmissionQuarantined            AdmissionPhase = "QUARANTINED"
	AdmissionNeedsMoreEvidence      AdmissionPhase = "NEEDS_MORE_EVIDENCE"
)

type EvidenceLevel string

const (
	EvidenceLevelInsufficient          EvidenceLevel = "insufficient"
	EvidenceLevelPartial               EvidenceLevel = "partial"
	EvidenceLevelDryRunReady           EvidenceLevel = "dryrun-ready"
	EvidenceLevelFutureApplyAdmissible EvidenceLevel = "future-apply-admissible"
)

type RequestedAdmissionAction string

const (
	RequestedActionCertifyShell       RequestedAdmissionAction = "certify-shell"
	RequestedActionPrepareCPULease    RequestedAdmissionAction = "prepare-cpu-lease"
	RequestedActionPrepareMemoryLease RequestedAdmissionAction = "prepare-memory-lease"
	RequestedActionEvidenceRefresh    RequestedAdmissionAction = "evidence-refresh"
	RequestedActionReturnToFloorCheck RequestedAdmissionAction = "return-to-floor-check"
	RequestedActionQuarantine         RequestedAdmissionAction = "quarantine"
)

type BlockerPriority string

const (
	BlockerPriorityP0Quarantine  BlockerPriority = "P0_QUARANTINE"
	BlockerPriorityP1HardBlock   BlockerPriority = "P1_HARD_BLOCK"
	BlockerPriorityP2Capability  BlockerPriority = "P2_CAPABILITY_BLOCK"
	BlockerPriorityP3Environment BlockerPriority = "P3_ENVIRONMENT_BLOCK"
	BlockerPriorityUnknown       BlockerPriority = "UNKNOWN"
)

type BlastRadiusPolicy struct {
	Scope  string `json:"scope"`
	MaxVMs int64  `json:"maxVms"`
}

type RollbackPolicy struct {
	RequireRollbackReady bool `json:"requireRollbackReady"`
}

type ReturnToFloorPolicy struct {
	RequireReturnToFloorReady bool `json:"requireReturnToFloorReady"`
	RequireMemorySafety       bool `json:"requireMemorySafety"`
}

type TTLPolicy struct {
	MaxLeaseTtlSeconds int64 `json:"maxLeaseTtlSeconds"`
}

type WindowsFluidPolicyPack struct {
	PolicyVersion                  string              `json:"policyVersion"`
	AllowedRuntimeMode             string              `json:"allowedRuntimeMode"`
	RequireCertifiedFluidShell     bool                `json:"requireCertifiedFluidShell"`
	RequireNoLiveMigration         bool                `json:"requireNoLiveMigration"`
	RequireNoReboot                bool                `json:"requireNoReboot"`
	RequireNoRecreate              bool                `json:"requireNoRecreate"`
	RequireSameNode                bool                `json:"requireSameNode"`
	RequireSameVirtLauncherPod     bool                `json:"requireSameVirtLauncherPod"`
	RequireSameQemuProcess         bool                `json:"requireSameQemuProcess"`
	RequireSameLastBoot            bool                `json:"requireSameLastBoot"`
	RequireSameMachineIdentity     bool                `json:"requireSameMachineIdentity"`
	RequireQmpAck                  bool                `json:"requireQmpAck"`
	RequireGuestAck                bool                `json:"requireGuestAck"`
	RequireRollbackReady           bool                `json:"requireRollbackReady"`
	RequireReturnToFloorReady      bool                `json:"requireReturnToFloorReady"`
	AllowPoolReplicaModel          bool                `json:"allowPoolReplicaModel"`
	AllowGenericWindowsVm          bool                `json:"allowGenericWindowsVm"`
	AllowMutationInThisPhase       bool                `json:"allowMutationInThisPhase"`
	MaxEvidenceAgeSeconds          int64               `json:"maxEvidenceAgeSeconds"`
	MaxLeaseTtlSeconds             int64               `json:"maxLeaseTtlSeconds"`
	MinEvidenceScoreForFutureApply int64               `json:"minEvidenceScoreForFutureApply"`
	BlastRadiusPolicy              BlastRadiusPolicy   `json:"blastRadiusPolicy"`
	BlockerPriorityPolicy          map[string]string   `json:"blockerPriorityPolicy"`
	RollbackPolicy                 RollbackPolicy      `json:"rollbackPolicy"`
	ReturnToFloorPolicy            ReturnToFloorPolicy `json:"returnToFloorPolicy"`
	TTLPolicy                      TTLPolicy           `json:"ttlPolicy"`
}

type WindowsFluidAdmissionDecision struct {
	DecisionID                 string                   `json:"decisionId"`
	ShellRef                   string                   `json:"shellRef"`
	RequestedAction            RequestedAdmissionAction `json:"requestedAction"`
	AdmissionPhase             AdmissionPhase           `json:"admissionPhase"`
	MutationAllowed            bool                     `json:"mutationAllowed"`
	ApplyAllowed               bool                     `json:"applyAllowed"`
	RuntimeMode                string                   `json:"runtimeMode"`
	EvidenceScore              int64                    `json:"evidenceScore"`
	EvidenceLevel              EvidenceLevel            `json:"evidenceLevel"`
	PolicyVersion              string                   `json:"policyVersion"`
	Blockers                   []string                 `json:"blockers"`
	DenialReasons              []string                 `json:"denialReasons"`
	RequiredAdditionalEvidence []string                 `json:"requiredAdditionalEvidence"`
	BlastRadius                BlastRadiusPolicy        `json:"blastRadius"`
	RollbackPolicy             RollbackPolicy           `json:"rollbackPolicy"`
	ReturnToFloorPolicy        ReturnToFloorPolicy      `json:"returnToFloorPolicy"`
	TTLPolicy                  TTLPolicy                `json:"ttlPolicy"`
	AuditRefs                  []string                 `json:"auditRefs"`
	CreatedAt                  string                   `json:"createdAt"`
}

func DefaultWindowsFluidPolicyPack() WindowsFluidPolicyPack {
	return WindowsFluidPolicyPack{
		PolicyVersion:                  "windows-fluid-admission-policy-v1",
		AllowedRuntimeMode:             "in-place-qmp",
		RequireCertifiedFluidShell:     true,
		RequireNoLiveMigration:         true,
		RequireNoReboot:                true,
		RequireNoRecreate:              true,
		RequireSameNode:                true,
		RequireSameVirtLauncherPod:     true,
		RequireSameQemuProcess:         true,
		RequireSameLastBoot:            true,
		RequireSameMachineIdentity:     true,
		RequireQmpAck:                  true,
		RequireGuestAck:                true,
		RequireRollbackReady:           true,
		RequireReturnToFloorReady:      true,
		AllowPoolReplicaModel:          false,
		AllowGenericWindowsVm:          false,
		AllowMutationInThisPhase:       false,
		MaxEvidenceAgeSeconds:          900,
		MaxLeaseTtlSeconds:             900,
		MinEvidenceScoreForFutureApply: 90,
		BlastRadiusPolicy: BlastRadiusPolicy{
			Scope:  "single-vm",
			MaxVMs: 1,
		},
		BlockerPriorityPolicy: DefaultBlockerPriorityPolicyMap(),
		RollbackPolicy: RollbackPolicy{
			RequireRollbackReady: true,
		},
		ReturnToFloorPolicy: ReturnToFloorPolicy{
			RequireReturnToFloorReady: true,
			RequireMemorySafety:       true,
		},
		TTLPolicy: TTLPolicy{
			MaxLeaseTtlSeconds: 900,
		},
	}
}

func DefaultBlockerPriorityPolicyMap() map[string]string {
	return map[string]string{
		BlockerQemuPIDChanged:               string(BlockerPriorityP0Quarantine),
		BlockerLastBootChanged:              string(BlockerPriorityP0Quarantine),
		BlockerMachineGUIDChanged:           string(BlockerPriorityP0Quarantine),
		BlockerNodeChanged:                  string(BlockerPriorityP0Quarantine),
		BlockerVirtLauncherPodChanged:       string(BlockerPriorityP0Quarantine),
		BlockerHotplugErrorDetected:         string(BlockerPriorityP0Quarantine),
		BlockerCriticalWindowsEventDetected: string(BlockerPriorityP0Quarantine),

		BlockerQMPSocketUnavailable:         string(BlockerPriorityP1HardBlock),
		BlockerGuestAgentUnavailable:        string(BlockerPriorityP1HardBlock),
		BlockerKarlAgentFluidModuleMissing:  string(BlockerPriorityP1HardBlock),
		BlockerPendingRebootDetected:        string(BlockerPriorityP1HardBlock),
		BlockerLiveMigrationRequired:        string(BlockerPriorityP1HardBlock),
		BlockerVMIRecreateRequired:          string(BlockerPriorityP1HardBlock),
		BlockerQMPAckMissing:                string(BlockerPriorityP1HardBlock),
		BlockerGuestAckMissing:              string(BlockerPriorityP1HardBlock),
		BlockerRollbackNotReady:             string(BlockerPriorityP1HardBlock),
		BlockerReturnToFloorNotReady:        string(BlockerPriorityP1HardBlock),
		BlockerMemoryReturnNotSafe:          string(BlockerPriorityP1HardBlock),
		BlockerFutureApplyExecutorDisabled:  string(BlockerPriorityP1HardBlock),
		BlockerNodeFluidActuatorUnavailable: string(BlockerPriorityP1HardBlock),
		BlockerPoolScalingAsMechanism:       string(BlockerPriorityP1HardBlock),
		BlockerStaleActuatorRequest:         string(BlockerPriorityP1HardBlock),
		BlockerLeaseRequestsVCPUHotplug:     string(BlockerPriorityP1HardBlock),
		BlockerLeaseRequestsVMSpecPatch:     string(BlockerPriorityP1HardBlock),

		BlockerMemoryDriverUnverified:            string(BlockerPriorityP2Capability),
		BlockerCPUTopologyNotConfirmed:           string(BlockerPriorityP2Capability),
		BlockerGuestMemoryNotConfirmed:           string(BlockerPriorityP2Capability),
		BlockerRAMBalloonUnavailable:             string(BlockerPriorityP2Capability),
		BlockerActuatorCPUEntitlementOutOfBounds: string(BlockerPriorityP2Capability),

		BlockerDashboard443TouchRisk:     string(BlockerPriorityP3Environment),
		BlockerCandidate8888Unavailable:  string(BlockerPriorityP3Environment),
		BlockerActuatorTargetAmbiguous:   string(BlockerPriorityP3Environment),
		BlockerActuatorPathEscape:        string(BlockerPriorityP3Environment),
		BlockerActuatorParentCgroupWrite: string(BlockerPriorityP3Environment),
		BlockerActuatorArbitraryWrite:    string(BlockerPriorityP3Environment),
		BlockerActuatorReplayDetected:    string(BlockerPriorityP3Environment),
		BlockerCgroupPathMismatch:        string(BlockerPriorityP0Quarantine),
	}
}

func BlockerPriorityForID(id string, policy WindowsFluidPolicyPack) BlockerPriority {
	raw, ok := policy.BlockerPriorityPolicy[id]
	if !ok {
		return BlockerPriorityUnknown
	}
	switch BlockerPriority(raw) {
	case BlockerPriorityP0Quarantine, BlockerPriorityP1HardBlock, BlockerPriorityP2Capability, BlockerPriorityP3Environment:
		return BlockerPriority(raw)
	default:
		return BlockerPriorityUnknown
	}
}
