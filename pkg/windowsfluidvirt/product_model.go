package windowsfluidvirt

import (
	"fmt"
)

type WindowsHyperdensityTargetKind string

const (
	TargetKindStandaloneWindowsVM WindowsHyperdensityTargetKind = "standalone-windows-vm"
	TargetKindPoolChildWindowsVM  WindowsHyperdensityTargetKind = "pool-child-windows-vm"
)

type WindowsHyperdensityRuntimeMode string

const (
	RuntimeModePrearmedFluidEnvelopeV2 WindowsHyperdensityRuntimeMode = "prearmed-fluid-envelope-v2"
)

type WindowsCPUEnvelope struct {
	Mechanism        string `json:"mechanism"`
	FloorCPUMax      string `json:"floorCpuMax"`
	CeilingCPUMax    string `json:"ceilingCpuMax"`
	CurrentCPUMax    string `json:"currentCpuMax"`
	ActuatorRequired bool   `json:"actuatorRequired"`
}

type WindowsMemoryEnvelope struct {
	Mechanism    string `json:"mechanism"`
	FloorBytes   int64  `json:"floorBytes"`
	CeilingBytes int64  `json:"ceilingBytes"`
	CurrentBytes int64  `json:"currentBytes"`
	QMPRequired  bool   `json:"qmpRequired"`
}

type WindowsGuestRequirements struct {
	FluidShellRequired bool `json:"fluidShellRequired"`
	GuestAckRequired   bool `json:"guestAckRequired"`
}

type WindowsRuntimeGuarantees struct {
	NoReboot              bool `json:"noReboot"`
	NoRecreate            bool `json:"noRecreate"`
	NoRollout             bool `json:"noRollout"`
	NoLiveMigration       bool `json:"noLiveMigration"`
	SameQEMU              bool `json:"sameQemu"`
	SameBoot              bool `json:"sameBoot"`
	RollbackRequired      bool `json:"rollbackRequired"`
	ReturnToFloorRequired bool `json:"returnToFloorRequired"`
}

type WindowsHyperdensityTarget struct {
	TargetID                 string                             `json:"targetId"`
	VMRef                    string                             `json:"vmRef"`
	Namespace                string                             `json:"namespace"`
	TargetKind               WindowsHyperdensityTargetKind      `json:"targetKind"`
	RuntimeMode              WindowsHyperdensityRuntimeMode     `json:"runtimeMode"`
	HyperdensityReady        bool                               `json:"hyperdensityReady"`
	CompliancePhase          WindowsHyperdensityCompliancePhase `json:"compliancePhase"`
	NodeName                 string                             `json:"nodeName"`
	VirtLauncherPodRef       string                             `json:"virtLauncherPodRef"`
	PodUID                   string                             `json:"podUid"`
	QemuPID                  string                             `json:"qemuPid"`
	QemuStartTime            string                             `json:"qemuStartTime"`
	MachineGuidHash          string                             `json:"machineGuidHash"`
	LastBootTime             string                             `json:"lastBootTime"`
	CPU                      WindowsCPUEnvelope                 `json:"cpu"`
	Memory                   WindowsMemoryEnvelope              `json:"memory"`
	Guest                    WindowsGuestRequirements           `json:"guest"`
	Guarantees               WindowsRuntimeGuarantees           `json:"guarantees"`
	Blockers                 []string                           `json:"blockers"`
	EvidenceRefs             []string                           `json:"evidenceRefs"`
	PoolScalingRequested     bool                               `json:"poolScalingRequested,omitempty"`
	LogicalCPUScalingClaimed bool                               `json:"logicalCpuScalingClaimed,omitempty"`
	VCPUHotplugRequested     bool                               `json:"vcpuHotplugRequested,omitempty"`
}

func EvaluateWindowsHyperdensityTarget(target WindowsHyperdensityTarget) WindowsHyperdensityTarget {
	blockers := make([]string, 0, 10)
	if target.RuntimeMode != RuntimeModePrearmedFluidEnvelopeV2 {
		blockers = append(blockers, BlockerPoolScalingAsMechanism)
	}
	if target.TargetKind != TargetKindStandaloneWindowsVM && target.TargetKind != TargetKindPoolChildWindowsVM {
		blockers = append(blockers, BlockerActuatorTargetAmbiguous)
	}
	if target.PoolScalingRequested {
		blockers = append(blockers, BlockerPoolScalingAsMechanism)
	}
	if target.LogicalCPUScalingClaimed {
		blockers = append(blockers, BlockerLeaseRequestsVCPUHotplug)
	}
	if target.VCPUHotplugRequested {
		blockers = append(blockers, BlockerLeaseRequestsVCPUHotplug)
	}
	if target.CPU.Mechanism != "cgroup-v2-cpu-max" || !target.CPU.ActuatorRequired {
		blockers = append(blockers, BlockerNodeFluidActuatorUnavailable)
	}
	if target.Memory.Mechanism != "qmp-balloon" || !target.Memory.QMPRequired {
		blockers = append(blockers, BlockerRAMBalloonUnavailable)
	}
	if target.Guest.FluidShellRequired && target.MachineGuidHash == "" {
		blockers = append(blockers, BlockerKarlAgentFluidModuleMissing)
	}
	if !target.Guarantees.ReturnToFloorRequired {
		blockers = append(blockers, BlockerReturnToFloorNotReady)
	}
	if !target.Guarantees.RollbackRequired {
		blockers = append(blockers, BlockerRollbackNotReady)
	}
	if target.CompliancePhase != ComplianceHyperdensityReadyWindowsShell {
		blockers = append(blockers, BlockerCPUTopologyNotConfirmed)
	}
	target.Blockers = dedupe(blockers)
	target.HyperdensityReady = target.CompliancePhase == ComplianceHyperdensityReadyWindowsShell && len(target.Blockers) == 0
	return target
}

type WindowsFluidLeaseKind string

const (
	LeaseKindCPUEntitlement   WindowsFluidLeaseKind = "cpu-entitlement"
	LeaseKindRAMBalloon       WindowsFluidLeaseKind = "ram-balloon"
	LeaseKindCombinedEnvelope WindowsFluidLeaseKind = "combined-envelope"
)

type WindowsFluidLeaseStatus string

const (
	LeaseStatusPrepared         WindowsFluidLeaseStatus = "prepared"
	LeaseStatusDryRunAccepted   WindowsFluidLeaseStatus = "dryRunAccepted"
	LeaseStatusApplying         WindowsFluidLeaseStatus = "applying"
	LeaseStatusActive           WindowsFluidLeaseStatus = "active"
	LeaseStatusReturningToFloor WindowsFluidLeaseStatus = "returningToFloor"
	LeaseStatusRolledBack       WindowsFluidLeaseStatus = "rolledBack"
	LeaseStatusBlocked          WindowsFluidLeaseStatus = "blocked"
	LeaseStatusQuarantined      WindowsFluidLeaseStatus = "quarantined"
)

type WindowsFluidLeaseRequest struct {
	CPUMax      string `json:"cpuMax,omitempty"`
	MemoryBytes int64  `json:"memoryBytes,omitempty"`
}

type WindowsFluidResourceLease struct {
	LeaseID                string                   `json:"leaseId"`
	TargetRef              string                   `json:"targetRef"`
	LeaseKind              WindowsFluidLeaseKind    `json:"leaseKind"`
	Requested              WindowsFluidLeaseRequest `json:"requested"`
	Previous               WindowsFluidLeaseRequest `json:"previous"`
	RollbackTarget         WindowsFluidLeaseRequest `json:"rollbackTarget"`
	ReturnToFloorTarget    WindowsFluidLeaseRequest `json:"returnToFloorTarget"`
	TTLSeconds             int64                    `json:"ttlSeconds"`
	Reason                 string                   `json:"reason"`
	Risk                   string                   `json:"risk"`
	PolicySnapshot         map[string]any           `json:"policySnapshot"`
	ActionSlateRef         string                   `json:"actionSlateRef"`
	ActuatorRequestRef     string                   `json:"actuatorRequestRef"`
	QMPRequestRef          string                   `json:"qmpRequestRef"`
	GuestEvidenceRef       string                   `json:"guestEvidenceRef"`
	AuditBundleRef         string                   `json:"auditBundleRef"`
	Status                 WindowsFluidLeaseStatus  `json:"status"`
	Blockers               []string                 `json:"blockers"`
	EvidenceRefs           []string                 `json:"evidenceRefs"`
	RequestedMechanism     string                   `json:"requestedMechanism,omitempty"`
	LogicalCPUScalingClaim bool                     `json:"logicalCpuScalingClaim,omitempty"`
	RequestsVCPUHotplug    bool                     `json:"requestsVcpuHotplug,omitempty"`
	RequestsVMSpecPatch    bool                     `json:"requestsVmSpecPatch,omitempty"`
}

func PrepareWindowsFluidResourceLease(target WindowsHyperdensityTarget, lease WindowsFluidResourceLease) WindowsFluidResourceLease {
	lease.Blockers = nil
	lease.Status = LeaseStatusPrepared
	blockers := make([]string, 0, 10)
	if !target.HyperdensityReady {
		blockers = append(blockers, BlockerCPUTopologyNotConfirmed)
	}
	if lease.TargetRef == "" || lease.LeaseID == "" {
		blockers = append(blockers, BlockerActuatorTargetAmbiguous)
	}
	if lease.TTLSeconds <= 0 {
		blockers = append(blockers, BlockerStaleActuatorRequest)
	}
	if lease.RollbackTarget.CPUMax == "" && lease.RollbackTarget.MemoryBytes == 0 {
		blockers = append(blockers, BlockerRollbackNotReady)
	}
	if lease.ReturnToFloorTarget.CPUMax == "" && lease.ReturnToFloorTarget.MemoryBytes == 0 {
		blockers = append(blockers, BlockerReturnToFloorNotReady)
	}
	if lease.RequestsVCPUHotplug {
		blockers = append(blockers, BlockerLeaseRequestsVCPUHotplug)
	}
	if lease.RequestsVMSpecPatch {
		blockers = append(blockers, BlockerLeaseRequestsVMSpecPatch)
	}
	if lease.LogicalCPUScalingClaim {
		blockers = append(blockers, BlockerLeaseRequestsVCPUHotplug)
	}
	if lease.RequestedMechanism == "pool-scaling" {
		blockers = append(blockers, BlockerPoolScalingAsMechanism)
	}
	if (lease.LeaseKind == LeaseKindCPUEntitlement || lease.LeaseKind == LeaseKindCombinedEnvelope) && !target.CPU.ActuatorRequired {
		blockers = append(blockers, BlockerNodeFluidActuatorUnavailable)
	}
	if (lease.LeaseKind == LeaseKindRAMBalloon || lease.LeaseKind == LeaseKindCombinedEnvelope) && !target.Memory.QMPRequired {
		blockers = append(blockers, BlockerRAMBalloonUnavailable)
	}
	lease.Blockers = dedupe(blockers)
	if len(lease.Blockers) > 0 {
		lease.Status = LeaseStatusBlocked
	}
	return lease
}

func EvaluateWindowsFluidLeasePreconditions(target WindowsHyperdensityTarget, lease WindowsFluidResourceLease) []string {
	blockers := make([]string, 0, 8)
	if !target.HyperdensityReady {
		blockers = append(blockers, BlockerCPUTopologyNotConfirmed)
	}
	if target.PoolScalingRequested {
		blockers = append(blockers, BlockerPoolScalingAsMechanism)
	}
	if target.CPU.ActuatorRequired && target.CPU.CurrentCPUMax == "" {
		blockers = append(blockers, BlockerNodeFluidActuatorUnavailable)
	}
	if target.Memory.QMPRequired && target.Memory.CurrentBytes == 0 {
		blockers = append(blockers, BlockerRAMBalloonUnavailable)
	}
	if target.Guest.GuestAckRequired && lease.GuestEvidenceRef == "" {
		blockers = append(blockers, BlockerGuestAckMissing)
	}
	if lease.RollbackTarget.CPUMax == "" && lease.RollbackTarget.MemoryBytes == 0 {
		blockers = append(blockers, BlockerRollbackNotReady)
	}
	if lease.ReturnToFloorTarget.CPUMax == "" && lease.ReturnToFloorTarget.MemoryBytes == 0 {
		blockers = append(blockers, BlockerReturnToFloorNotReady)
	}
	return dedupe(blockers)
}

type WindowsFluidActionType string

const (
	ActionComplianceReplay    WindowsFluidActionType = "complianceReplay"
	ActionBuildActuatorReq    WindowsFluidActionType = "buildActuatorRequest"
	ActionActuatorDryRun      WindowsFluidActionType = "actuatorDryRun"
	ActionCPUEntitlementApply WindowsFluidActionType = "cpuEntitlementApply"
	ActionCPUReturnToFloor    WindowsFluidActionType = "cpuReturnToFloor"
	ActionQMPBalloonApply     WindowsFluidActionType = "qmpBalloonApply"
	ActionRAMReturnToFloor    WindowsFluidActionType = "ramReturnToFloor"
	ActionGuestVerify         WindowsFluidActionType = "guestVerify"
	ActionFinalRestore        WindowsFluidActionType = "finalRestore"
	ActionAuditBundleAppend   WindowsFluidActionType = "auditBundleAppend"
)

type WindowsFluidAction struct {
	ActionID            string                 `json:"actionId"`
	ActionType          WindowsFluidActionType `json:"actionType"`
	MutationAllowed     bool                   `json:"mutationAllowed"`
	Preconditions       []string               `json:"preconditions"`
	RequiredEvidence    []string               `json:"requiredEvidence"`
	ExpectedOutput      string                 `json:"expectedOutput"`
	RollbackAction      string                 `json:"rollbackAction,omitempty"`
	ReturnToFloorAction string                 `json:"returnToFloorAction,omitempty"`
	Blockers            []string               `json:"blockers,omitempty"`
}

func BuildWindowsFluidActionSlate(target WindowsHyperdensityTarget, lease WindowsFluidResourceLease) WindowsFluidActionSlate {
	blockers := EvaluateWindowsFluidLeasePreconditions(target, lease)
	// ActionType still names future gated phases; MutationAllowed stays false on every
	// step so the default slate is planning/readiness only (no execution claims).
	actions := []WindowsFluidAction{
		{ActionID: "a1", ActionType: ActionComplianceReplay, MutationAllowed: false, Preconditions: []string{"target-observed"}, RequiredEvidence: []string{"compliance-input"}, ExpectedOutput: "compliance-ready"},
		{ActionID: "a2", ActionType: ActionBuildActuatorReq, MutationAllowed: false, Preconditions: []string{"compliance-ready"}, RequiredEvidence: []string{"target-identity", "cpu-bounds"}, ExpectedOutput: "actuator-request"},
		{ActionID: "a3", ActionType: ActionActuatorDryRun, MutationAllowed: false, Preconditions: []string{"actuator-request-built"}, RequiredEvidence: []string{"allowlist", "kill-switch"}, ExpectedOutput: "dry-run-accepted"},
		{ActionID: "a4", ActionType: ActionCPUEntitlementApply, MutationAllowed: false, Preconditions: []string{"dry-run-accepted"}, RequiredEvidence: []string{"actuator-request"}, ExpectedOutput: "cpu-entitlement-applied", RollbackAction: "a5", ReturnToFloorAction: "a5"},
		{ActionID: "a5", ActionType: ActionCPUReturnToFloor, MutationAllowed: false, Preconditions: []string{"cpu-entitlement-applied"}, RequiredEvidence: []string{"cpu-before-after"}, ExpectedOutput: "cpu-floor-restored"},
		{ActionID: "a6", ActionType: ActionQMPBalloonApply, MutationAllowed: false, Preconditions: []string{"qmp-available"}, RequiredEvidence: []string{"qmp-request"}, ExpectedOutput: "ram-balloon-applied", RollbackAction: "a7", ReturnToFloorAction: "a7"},
		{ActionID: "a7", ActionType: ActionRAMReturnToFloor, MutationAllowed: false, Preconditions: []string{"ram-balloon-applied"}, RequiredEvidence: []string{"qmp-balloon-before-after"}, ExpectedOutput: "ram-floor-restored"},
		{ActionID: "a8", ActionType: ActionGuestVerify, MutationAllowed: false, Preconditions: []string{"cpu-ram-updated"}, RequiredEvidence: []string{"guest-ack"}, ExpectedOutput: "guest-verified"},
		{ActionID: "a9", ActionType: ActionFinalRestore, MutationAllowed: false, Preconditions: []string{"verification-complete"}, RequiredEvidence: []string{"restore-targets"}, ExpectedOutput: "final-restore-complete"},
		{ActionID: "a10", ActionType: ActionAuditBundleAppend, MutationAllowed: false, Preconditions: []string{"final-restore-complete"}, RequiredEvidence: []string{"bundle-index", "attestation"}, ExpectedOutput: "audit-chain-appended"},
	}
	return WindowsFluidActionSlate{
		ActionID:                "windows-fluid-slate-" + shortHash(target.TargetID+"|"+lease.LeaseID),
		ShellRef:                target.VMRef,
		ActionType:              ActionEvidenceRefresh,
		DryRunState:             DryRunReady,
		RuntimeMode:             string(RuntimeModePrearmedFluidEnvelopeV2),
		MutationAllowed:         false,
		ApplyAllowed:            false,
		RollbackReady:           len(EvaluateWindowsFluidRollbackReadiness(lease)) == 0,
		ReturnToFloorReady:      len(EvaluateWindowsFluidReturnToFloorReadiness(lease)) == 0,
		QMPReady:                target.Memory.QMPRequired,
		GuestAckReady:           lease.GuestEvidenceRef != "",
		NoRebootProof:           target.Guarantees.NoReboot,
		SameQemuProof:           target.Guarantees.SameQEMU,
		SameNodeProof:           target.NodeName != "",
		SamePodProof:            target.PodUID != "",
		TargetRef:               target.TargetID,
		LeaseRef:                lease.LeaseID,
		Actions:                 actions,
		RuntimeMutationExecuted: false,
		Blockers:                blockers,
		EvidenceRefs:            dedupe(append([]string{}, lease.EvidenceRefs...)),
		ValidForSeconds:         lease.TTLSeconds,
	}
}

func EvaluateWindowsFluidRollbackReadiness(lease WindowsFluidResourceLease) []string {
	blockers := make([]string, 0, 2)
	if lease.RollbackTarget.CPUMax == "" && lease.RollbackTarget.MemoryBytes == 0 {
		blockers = append(blockers, BlockerRollbackNotReady)
	}
	return dedupe(blockers)
}

func EvaluateWindowsFluidReturnToFloorReadiness(lease WindowsFluidResourceLease) []string {
	blockers := make([]string, 0, 2)
	if lease.ReturnToFloorTarget.CPUMax == "" && lease.ReturnToFloorTarget.MemoryBytes == 0 {
		blockers = append(blockers, BlockerReturnToFloorNotReady)
	}
	return dedupe(blockers)
}

func EvaluateWindowsFluidAuditReadiness(lease WindowsFluidResourceLease, slate WindowsFluidActionSlate) []string {
	blockers := make([]string, 0, 3)
	if lease.AuditBundleRef == "" {
		blockers = append(blockers, BlockerActuatorReplayDetected)
	}
	if lease.GuestEvidenceRef == "" {
		blockers = append(blockers, BlockerGuestAckMissing)
	}
	if slate.RuntimeMutationExecuted {
		blockers = append(blockers, BlockerActuatorArbitraryWrite)
	}
	return dedupe(blockers)
}

func RequireNoRuntimeMutation(slate WindowsFluidActionSlate) error {
	if slate.RuntimeMutationExecuted {
		return fmt.Errorf("runtime mutation detected in evaluator path")
	}
	return nil
}
