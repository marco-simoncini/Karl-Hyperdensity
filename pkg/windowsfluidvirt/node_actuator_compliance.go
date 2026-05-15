package windowsfluidvirt

import "time"

type WindowsHyperdensityCompliancePhase string

const (
	ComplianceDiscoveredWindowsVM           WindowsHyperdensityCompliancePhase = "DISCOVERED_WINDOWS_VM"
	ComplianceAssessedWindowsVM             WindowsHyperdensityCompliancePhase = "ASSESSED_WINDOWS_VM"
	ComplianceBlockedWithRemediation        WindowsHyperdensityCompliancePhase = "BLOCKED_WITH_REMEDIATION"
	ComplianceFluidEnvelopeCandidate        WindowsHyperdensityCompliancePhase = "FLUID_ENVELOPE_CANDIDATE"
	ComplianceHyperdensityReadyWindowsShell WindowsHyperdensityCompliancePhase = "HYPERDENSITY_READY_WINDOWS_SHELL"
)

type WindowsHyperdensityRisk string

const (
	RiskLow    WindowsHyperdensityRisk = "low"
	RiskMedium WindowsHyperdensityRisk = "medium"
	RiskHigh   WindowsHyperdensityRisk = "high"
)

type WindowsHyperdensityRemediationID string

const (
	RemediationInstallOrEnableFluidShell    WindowsHyperdensityRemediationID = "install_or_enable_fluidShell"
	RemediationEnableQMPAccess              WindowsHyperdensityRemediationID = "enable_qmp_access"
	RemediationVerifyQMPSocketMapping       WindowsHyperdensityRemediationID = "verify_qmp_socket_mapping"
	RemediationEnableRAMBalloonDriver       WindowsHyperdensityRemediationID = "enable_ram_balloon_driver"
	RemediationVerifyQMPBalloon             WindowsHyperdensityRemediationID = "verify_qmp_balloon"
	RemediationDeployNodeFluidActuator      WindowsHyperdensityRemediationID = "deploy_node_fluid_actuator"
	RemediationAuthorizeCPUCgroupTarget     WindowsHyperdensityRemediationID = "authorize_cpu_cgroup_target"
	RemediationApplyReadyAnnotations        WindowsHyperdensityRemediationID = "apply_hyperdensity_ready_annotations"
	RemediationDisableLiveMigrationForFluid WindowsHyperdensityRemediationID = "disable_livemigration_for_fluid_actions"
	RemediationConfigurePrearmedEnvelope    WindowsHyperdensityRemediationID = "configure_prearmed_envelope"
	RemediationVerifyReturnToFloor          WindowsHyperdensityRemediationID = "verify_return_to_floor"
	RemediationVerifyRollback               WindowsHyperdensityRemediationID = "verify_rollback"
	RemediationQuarantineIdentityMismatch   WindowsHyperdensityRemediationID = "quarantine_identity_mismatch"
	RemediationBlockPoolScalingAsMechanism  WindowsHyperdensityRemediationID = "block_pool_scaling_as_mechanism"
)

type WindowsRemediationDefinition struct {
	ID               WindowsHyperdensityRemediationID `json:"id"`
	Automatable      bool                             `json:"automatable"`
	Risk             WindowsHyperdensityRisk          `json:"risk"`
	RequiredEvidence []string                         `json:"requiredEvidence"`
	TargetComponent  string                           `json:"targetComponent"`
	Rollback         string                           `json:"rollback"`
	BlockerResolved  string                           `json:"blockerResolved"`
}

func DefaultWindowsHyperdensityRemediationTaxonomy() map[WindowsHyperdensityRemediationID]WindowsRemediationDefinition {
	return map[WindowsHyperdensityRemediationID]WindowsRemediationDefinition{
		RemediationInstallOrEnableFluidShell: {
			ID:               RemediationInstallOrEnableFluidShell,
			Automatable:      true,
			Risk:             RiskMedium,
			RequiredEvidence: []string{"guest_module_inventory", "guest_ack_probe"},
			TargetComponent:  "Karl-Inventory/windows-agent",
			Rollback:         "disable_fluid_shell_module_and_restore_previous_agent_profile",
			BlockerResolved:  BlockerKarlAgentFluidModuleMissing,
		},
		RemediationEnableQMPAccess: {
			ID:               RemediationEnableQMPAccess,
			Automatable:      true,
			Risk:             RiskMedium,
			RequiredEvidence: []string{"qmp_socket_probe", "qmp_capability_probe"},
			TargetComponent:  "Karl-Hyperdensity/qmp-runtime",
			Rollback:         "restore_previous_qmp_access_policy",
			BlockerResolved:  BlockerQMPSocketUnavailable,
		},
		RemediationVerifyQMPSocketMapping: {
			ID:               RemediationVerifyQMPSocketMapping,
			Automatable:      false,
			Risk:             RiskHigh,
			RequiredEvidence: []string{"vm_to_pod_mapping", "qmp_socket_path_mapping"},
			TargetComponent:  "KubeVirt/vmi-runtime",
			Rollback:         "quarantine_shell_and_rebuild_runtime_mapping",
			BlockerResolved:  BlockerActuatorTargetAmbiguous,
		},
		RemediationEnableRAMBalloonDriver: {
			ID:               RemediationEnableRAMBalloonDriver,
			Automatable:      true,
			Risk:             RiskMedium,
			RequiredEvidence: []string{"guest_driver_inventory", "guest_memory_adapter_probe"},
			TargetComponent:  "Karl-Inventory/windows-guest",
			Rollback:         "restore_previous_guest_driver_profile",
			BlockerResolved:  BlockerRAMBalloonUnavailable,
		},
		RemediationVerifyQMPBalloon: {
			ID:               RemediationVerifyQMPBalloon,
			Automatable:      true,
			Risk:             RiskLow,
			RequiredEvidence: []string{"qmp_query_balloon", "libvirt_balloon_sample"},
			TargetComponent:  "Karl-Hyperdensity/qmp-runtime",
			Rollback:         "keep_memory_lease_disabled_until_balloon_verification_passes",
			BlockerResolved:  BlockerRAMBalloonUnavailable,
		},
		RemediationDeployNodeFluidActuator: {
			ID:               RemediationDeployNodeFluidActuator,
			Automatable:      true,
			Risk:             RiskHigh,
			RequiredEvidence: []string{"node_actuator_health", "node_actuator_attestation"},
			TargetComponent:  "Karl-OS-ISO/node-runtime",
			Rollback:         "disable_actuator_endpoint_and_revert_to_no-mutation-mode",
			BlockerResolved:  BlockerNodeFluidActuatorUnavailable,
		},
		RemediationAuthorizeCPUCgroupTarget: {
			ID:               RemediationAuthorizeCPUCgroupTarget,
			Automatable:      false,
			Risk:             RiskHigh,
			RequiredEvidence: []string{"allowlist_entry", "resolved_cgroup_path", "qemu_pid_start_time"},
			TargetComponent:  "Karl-Hyperdensity/node-actuator",
			Rollback:         "remove_allowlist_entry_and_return_cpu_to_floor",
			BlockerResolved:  BlockerCgroupPathMismatch,
		},
		RemediationApplyReadyAnnotations: {
			ID:               RemediationApplyReadyAnnotations,
			Automatable:      true,
			Risk:             RiskLow,
			RequiredEvidence: []string{"policy_annotations_snapshot"},
			TargetComponent:  "Karl-Hyperdensity/compliance-engine",
			Rollback:         "restore_previous_annotation_set",
			BlockerResolved:  BlockerCPUTopologyNotConfirmed,
		},
		RemediationDisableLiveMigrationForFluid: {
			ID:               RemediationDisableLiveMigrationForFluid,
			Automatable:      false,
			Risk:             RiskHigh,
			RequiredEvidence: []string{"migration_policy_snapshot", "vmim_zero_check"},
			TargetComponent:  "KubeVirt/cluster-policy",
			Rollback:         "restore_migration_policy_after_fluid_operation_window",
			BlockerResolved:  BlockerLiveMigrationRequired,
		},
		RemediationConfigurePrearmedEnvelope: {
			ID:               RemediationConfigurePrearmedEnvelope,
			Automatable:      false,
			Risk:             RiskHigh,
			RequiredEvidence: []string{"cpu_floor_ceiling", "ram_floor_ceiling", "envelope_annotations"},
			TargetComponent:  "Karl-Hyperdensity/windows-shell",
			Rollback:         "revert_to_previous_shell_profile",
			BlockerResolved:  BlockerCPUTopologyNotConfirmed,
		},
		RemediationVerifyReturnToFloor: {
			ID:               RemediationVerifyReturnToFloor,
			Automatable:      true,
			Risk:             RiskMedium,
			RequiredEvidence: []string{"return_to_floor_result", "guest_ack_result"},
			TargetComponent:  "Karl-Hyperdensity/runtime-evidence",
			Rollback:         "block_new_leases_until_return_to_floor_proof_is_rebuilt",
			BlockerResolved:  BlockerReturnToFloorNotReady,
		},
		RemediationVerifyRollback: {
			ID:               RemediationVerifyRollback,
			Automatable:      true,
			Risk:             RiskMedium,
			RequiredEvidence: []string{"rollback_result", "rollback_readiness_probe"},
			TargetComponent:  "Karl-Hyperdensity/runtime-evidence",
			Rollback:         "disable_lease_transition_until_rollback_readiness_is_restored",
			BlockerResolved:  BlockerRollbackNotReady,
		},
		RemediationQuarantineIdentityMismatch: {
			ID:               RemediationQuarantineIdentityMismatch,
			Automatable:      true,
			Risk:             RiskHigh,
			RequiredEvidence: []string{"qemu_pid_mismatch", "pod_uid_mismatch", "vm_uid_mismatch"},
			TargetComponent:  "Karl-Hyperdensity/identity-guard",
			Rollback:         "clear_quarantine_only_after_fresh_continuity_evidence",
			BlockerResolved:  BlockerQemuPIDChanged,
		},
		RemediationBlockPoolScalingAsMechanism: {
			ID:               RemediationBlockPoolScalingAsMechanism,
			Automatable:      true,
			Risk:             RiskLow,
			RequiredEvidence: []string{"pool_context", "requested_scaling_mechanism"},
			TargetComponent:  "Karl-Hyperdensity/compliance-engine",
			Rollback:         "force_single_shell_lease_model_and_disable_pool_mechanism",
			BlockerResolved:  BlockerPoolScalingAsMechanism,
		},
	}
}

type WindowsVMIdentityEvidence struct {
	VMRef              string `json:"vmRef"`
	Namespace          string `json:"namespace"`
	VMUID              string `json:"vmUid"`
	VMIUID             string `json:"vmiUid"`
	VirtLauncherPodRef string `json:"virtLauncherPodRef"`
	PodUID             string `json:"podUid"`
	NodeName           string `json:"nodeName"`
	QemuPID            string `json:"qemuPid"`
	QemuStartTime      string `json:"qemuStartTime"`
}

type WindowsGuestFluidShellEvidence struct {
	Available   bool `json:"available"`
	GuestAck    bool `json:"guestAck"`
	SameBoot    bool `json:"sameBoot"`
	SameMachine bool `json:"sameMachine"`
}

type WindowsRAMBalloonEvidence struct {
	Available          bool `json:"available"`
	BalloonDriverReady bool `json:"balloonDriverReady"`
	CanReturnToFloor   bool `json:"canReturnToFloor"`
}

type WindowsCPUActuatorCapabilityEvidence struct {
	Available              bool   `json:"available"`
	NodeLocal              bool   `json:"nodeLocal"`
	Allowlisted            bool   `json:"allowlisted"`
	SupportsCPUmax         bool   `json:"supportsCpuMax"`
	CgroupPathMatchesShell bool   `json:"cgroupPathMatchesShell"`
	ActuatorAck            bool   `json:"actuatorAck"`
	ResolvedCgroupPath     string `json:"resolvedCgroupPath"`
}

type WindowsPoolContext struct {
	IsPoolChild            bool `json:"isPoolChild"`
	RequestedAsMechanism   bool `json:"requestedAsMechanism"`
	TreatedAsProvisionOnly bool `json:"treatedAsProvisionOnly"`
}

type EvaluateWindowsHyperdensityReadyComplianceInput struct {
	Identity           WindowsVMIdentityEvidence            `json:"identity"`
	QMPAvailable       bool                                 `json:"qmpAvailable"`
	GuestFluidShell    WindowsGuestFluidShellEvidence       `json:"guestFluidShell"`
	RAMBalloon         WindowsRAMBalloonEvidence            `json:"ramBalloon"`
	CPUActuator        WindowsCPUActuatorCapabilityEvidence `json:"cpuActuator"`
	PolicyAnnotations  map[string]string                    `json:"policyAnnotations"`
	PoolContext        WindowsPoolContext                   `json:"poolContext"`
	RemediationOptions []WindowsHyperdensityRemediationID   `json:"remediationOptions"`
}

type EvaluateWindowsHyperdensityReadyComplianceOutput struct {
	CompliancePhase    WindowsHyperdensityCompliancePhase `json:"compliancePhase"`
	Blockers           []string                           `json:"blockers"`
	RemediationActions []string                           `json:"remediationActions"`
	AutomatableActions []string                           `json:"automatableActions"`
	ManualActions      []string                           `json:"manualActions"`
	Risk               WindowsHyperdensityRisk            `json:"risk"`
	EvidenceSummary    map[string]any                     `json:"evidenceSummary"`
}

func EvaluateWindowsHyperdensityReadyCompliance(input EvaluateWindowsHyperdensityReadyComplianceInput) EvaluateWindowsHyperdensityReadyComplianceOutput {
	taxonomy := DefaultWindowsHyperdensityRemediationTaxonomy()
	blockers := make([]string, 0, 8)
	actions := make([]WindowsHyperdensityRemediationID, 0, 8)
	phase := ComplianceDiscoveredWindowsVM

	if input.Identity.VMRef != "" && input.Identity.Namespace != "" {
		phase = ComplianceAssessedWindowsVM
	}
	if input.PoolContext.RequestedAsMechanism {
		blockers = append(blockers, BlockerPoolScalingAsMechanism)
		actions = append(actions, RemediationBlockPoolScalingAsMechanism)
	}
	if !input.GuestFluidShell.Available {
		blockers = append(blockers, BlockerKarlAgentFluidModuleMissing)
		actions = append(actions, RemediationInstallOrEnableFluidShell)
	}
	if !input.QMPAvailable {
		blockers = append(blockers, BlockerQMPSocketUnavailable)
		actions = append(actions, RemediationEnableQMPAccess, RemediationVerifyQMPSocketMapping)
	}
	if !input.RAMBalloon.Available || !input.RAMBalloon.BalloonDriverReady {
		blockers = append(blockers, BlockerRAMBalloonUnavailable)
		actions = append(actions, RemediationEnableRAMBalloonDriver, RemediationVerifyQMPBalloon)
	}
	if !input.CPUActuator.Available || !input.CPUActuator.NodeLocal || !input.CPUActuator.SupportsCPUmax {
		blockers = append(blockers, BlockerNodeFluidActuatorUnavailable)
		actions = append(actions, RemediationDeployNodeFluidActuator)
	}
	if !input.CPUActuator.CgroupPathMatchesShell {
		blockers = append(blockers, BlockerCgroupPathMismatch)
		actions = append(actions, RemediationAuthorizeCPUCgroupTarget, RemediationQuarantineIdentityMismatch)
	}
	if !input.GuestFluidShell.SameBoot || !input.GuestFluidShell.SameMachine {
		blockers = append(blockers, BlockerQemuPIDChanged)
		actions = append(actions, RemediationQuarantineIdentityMismatch)
	}
	if !input.RAMBalloon.CanReturnToFloor {
		blockers = append(blockers, BlockerReturnToFloorNotReady)
		actions = append(actions, RemediationVerifyReturnToFloor)
	}
	if !input.GuestFluidShell.GuestAck {
		blockers = append(blockers, BlockerGuestAckMissing)
		actions = append(actions, RemediationInstallOrEnableFluidShell)
	}

	blockers = dedupe(blockers)
	actions = dedupeRemediations(actions)
	if len(blockers) > 0 {
		phase = ComplianceBlockedWithRemediation
	} else {
		phase = ComplianceFluidEnvelopeCandidate
		readyAnnotations := input.PolicyAnnotations["hyperdensity.karl.io/windows-ready"] == "true"
		if readyAnnotations && input.PoolContext.TreatedAsProvisionOnly && input.CPUActuator.Allowlisted && input.CPUActuator.ActuatorAck {
			phase = ComplianceHyperdensityReadyWindowsShell
		} else {
			actions = append(actions, RemediationApplyReadyAnnotations, RemediationConfigurePrearmedEnvelope)
		}
	}

	automatable := make([]string, 0, len(actions))
	manual := make([]string, 0, len(actions))
	for _, action := range actions {
		def, ok := taxonomy[action]
		if !ok {
			continue
		}
		if def.Automatable {
			automatable = append(automatable, string(action))
		} else {
			manual = append(manual, string(action))
		}
	}

	risk := RiskLow
	if len(blockers) > 0 {
		risk = RiskMedium
	}
	for _, blocker := range blockers {
		def, ok := LookupBlocker(blocker)
		if !ok {
			continue
		}
		if def.Severity == SeverityCritical {
			risk = RiskHigh
			break
		}
	}

	return EvaluateWindowsHyperdensityReadyComplianceOutput{
		CompliancePhase:    phase,
		Blockers:           blockers,
		RemediationActions: remediationsToStrings(actions),
		AutomatableActions: dedupe(automatable),
		ManualActions:      dedupe(manual),
		Risk:               risk,
		EvidenceSummary: map[string]any{
			"vmRef":                    input.Identity.VMRef,
			"namespace":                input.Identity.Namespace,
			"isPoolChild":              input.PoolContext.IsPoolChild,
			"poolProvisioningOnly":     input.PoolContext.TreatedAsProvisionOnly,
			"qmpAvailable":             input.QMPAvailable,
			"guestFluidShellAvailable": input.GuestFluidShell.Available,
			"ramBalloonReady":          input.RAMBalloon.Available && input.RAMBalloon.BalloonDriverReady,
			"cpuActuatorReady":         input.CPUActuator.Available && input.CPUActuator.NodeLocal && input.CPUActuator.SupportsCPUmax,
			"cgroupPath":               input.CPUActuator.ResolvedCgroupPath,
			"guestAck":                 input.GuestFluidShell.GuestAck,
			"sameBoot":                 input.GuestFluidShell.SameBoot,
		},
	}
}

type WindowsCpuEntitlementLease struct {
	LeaseID             string   `json:"leaseId"`
	ShellRef            string   `json:"shellRef"`
	TargetVM            string   `json:"targetVm"`
	Mode                string   `json:"mode"`
	Mechanism           string   `json:"mechanism"`
	FloorCPUMax         string   `json:"floorCpuMax"`
	CeilingCPUMax       string   `json:"ceilingCpuMax"`
	CurrentCPUMax       string   `json:"currentCpuMax"`
	RequestedCPUMax     string   `json:"requestedCpuMax"`
	TTLSeconds          int64    `json:"ttlSeconds"`
	Reason              string   `json:"reason"`
	Risk                string   `json:"risk"`
	RollbackTarget      string   `json:"rollbackTarget"`
	ReturnToFloorTarget string   `json:"returnToFloorTarget"`
	RequireSameQEMU     bool     `json:"requireSameQemu"`
	RequireSameBoot     bool     `json:"requireSameBoot"`
	RequireSamePod      bool     `json:"requireSamePod"`
	RequireSameNode     bool     `json:"requireSameNode"`
	RequireGuestAck     bool     `json:"requireGuestAck"`
	RequireActuatorAck  bool     `json:"requireActuatorAck"`
	Status              string   `json:"status"`
	Blockers            []string `json:"blockers"`
	EvidenceRefs        []string `json:"evidenceRefs"`
	DisallowedActions   []string `json:"disallowedActions,omitempty"`
}

func ValidateWindowsCpuEntitlementLease(lease WindowsCpuEntitlementLease) []string {
	blockers := make([]string, 0, 8)
	if lease.LeaseID == "" || lease.ShellRef == "" || lease.TargetVM == "" {
		blockers = append(blockers, BlockerActuatorTargetAmbiguous)
	}
	if lease.Mode != "prearmed-fluid-envelope" {
		blockers = append(blockers, BlockerPoolScalingAsMechanism)
	}
	if lease.Mechanism != "cgroup-v2-cpu-max" {
		blockers = append(blockers, BlockerLeaseRequestsVCPUHotplug)
	}
	if lease.TTLSeconds <= 0 {
		blockers = append(blockers, BlockerStaleActuatorRequest)
	}
	if lease.RequestedCPUMax == "" || lease.FloorCPUMax == "" || lease.CeilingCPUMax == "" {
		blockers = append(blockers, BlockerActuatorTargetAmbiguous)
	}
	if lease.RollbackTarget == "" || lease.ReturnToFloorTarget == "" {
		blockers = append(blockers, BlockerRollbackNotReady)
	}
	if !lease.RequireSameQEMU || !lease.RequireSameBoot || !lease.RequireSamePod || !lease.RequireSameNode {
		blockers = append(blockers, BlockerQemuPIDChanged)
	}
	if !lease.RequireGuestAck {
		blockers = append(blockers, BlockerGuestAckMissing)
	}
	if !lease.RequireActuatorAck {
		blockers = append(blockers, BlockerNodeFluidActuatorUnavailable)
	}
	for _, action := range lease.DisallowedActions {
		switch action {
		case "vcpu-hotplug", "vcpu-unplug":
			blockers = append(blockers, BlockerLeaseRequestsVCPUHotplug)
		case "vm-spec-patch":
			blockers = append(blockers, BlockerLeaseRequestsVMSpecPatch)
		case "livemigration":
			blockers = append(blockers, BlockerLiveMigrationRequired)
		}
	}
	return dedupe(blockers)
}

func EvaluateWindowsCpuEntitlementLease(lease WindowsCpuEntitlementLease) WindowsCpuEntitlementLease {
	blockers := ValidateWindowsCpuEntitlementLease(lease)
	lease.Blockers = blockers
	if len(blockers) > 0 {
		lease.Status = "rejected"
		return lease
	}
	lease.Status = "accepted_as_lease_plan"
	return lease
}

type KARLNodeFluidActuatorContract struct {
	ActuatorID            string         `json:"actuatorId"`
	NodeName              string         `json:"nodeName"`
	ActuatorVersion       string         `json:"actuatorVersion"`
	RequestID             string         `json:"requestId"`
	ShellRef              string         `json:"shellRef"`
	VMRef                 string         `json:"vmRef"`
	Namespace             string         `json:"namespace"`
	VirtLauncherPodRef    string         `json:"virtLauncherPodRef"`
	PodUID                string         `json:"podUid"`
	QemuPID               string         `json:"qemuPid"`
	QemuStartTime         string         `json:"qemuStartTime"`
	CgroupPath            string         `json:"cgroupPath"`
	AllowedControllers    []string       `json:"allowedControllers"`
	RequestedCPUMax       string         `json:"requestedCpuMax"`
	PreviousCPUMax        string         `json:"previousCpuMax"`
	AppliedCPUMax         string         `json:"appliedCpuMax"`
	RollbackCPUMax        string         `json:"rollbackCpuMax"`
	MutationScope         string         `json:"mutationScope"`
	AllowlistDecision     bool           `json:"allowlistDecision"`
	PolicyDecision        string         `json:"policyDecision"`
	BeforeEvidence        map[string]any `json:"beforeEvidence"`
	AfterEvidence         map[string]any `json:"afterEvidence"`
	RollbackEvidence      map[string]any `json:"rollbackEvidence"`
	ReturnToFloorEvidence map[string]any `json:"returnToFloorEvidence"`
	Blockers              []string       `json:"blockers"`
	CreatedAt             string         `json:"createdAt"`
}

func ValidateKARLNodeFluidActuatorContract(contract KARLNodeFluidActuatorContract) []string {
	blockers := make([]string, 0, 8)
	if contract.ActuatorID == "" || contract.RequestID == "" || contract.NodeName == "" {
		blockers = append(blockers, BlockerActuatorTargetAmbiguous)
	}
	if contract.VMRef == "" || contract.Namespace == "" || contract.VirtLauncherPodRef == "" || contract.PodUID == "" {
		blockers = append(blockers, BlockerActuatorTargetAmbiguous)
	}
	if contract.QemuPID == "" || contract.QemuStartTime == "" {
		blockers = append(blockers, BlockerQemuPIDChanged)
	}
	if contract.CgroupPath == "" {
		blockers = append(blockers, BlockerCgroupPathMismatch)
	}
	if !contract.AllowlistDecision {
		blockers = append(blockers, BlockerActuatorTargetAmbiguous)
	}
	if contract.PolicyDecision != "allowed" && contract.PolicyDecision != "future-signable-allowed" {
		blockers = append(blockers, BlockerNodeFluidActuatorUnavailable)
	}
	if contract.MutationScope != "cpu-entitlement-only" {
		blockers = append(blockers, BlockerActuatorArbitraryWrite)
	}
	allowed := map[string]struct{}{"cpu.max": {}, "cpu.weight": {}}
	for _, controller := range contract.AllowedControllers {
		if _, ok := allowed[controller]; !ok {
			blockers = append(blockers, BlockerActuatorArbitraryWrite)
		}
	}
	if contract.RequestedCPUMax == "" || contract.PreviousCPUMax == "" || contract.AppliedCPUMax == "" || contract.RollbackCPUMax == "" {
		blockers = append(blockers, BlockerActuatorCPUEntitlementOutOfBounds)
	}
	if len(contract.BeforeEvidence) == 0 || len(contract.AfterEvidence) == 0 {
		blockers = append(blockers, BlockerQMPAckMissing)
	}
	if len(contract.RollbackEvidence) == 0 || len(contract.ReturnToFloorEvidence) == 0 {
		blockers = append(blockers, BlockerReturnToFloorNotReady)
	}
	return dedupe(blockers)
}

type NodeFluidActuatorSafetyModel struct {
	NodeLocalBlastRadiusOnly       bool  `json:"nodeLocalBlastRadiusOnly"`
	TargetAllowlistRequired        bool  `json:"targetAllowlistRequired"`
	ValidateCgroupPath             bool  `json:"validateCgroupPath"`
	ValidatePIDAndStartTime        bool  `json:"validatePidAndStartTime"`
	ValidatePodUID                 bool  `json:"validatePodUid"`
	ValidateVMUID                  bool  `json:"validateVmUid"`
	NoSymlinkTraversal             bool  `json:"noSymlinkTraversal"`
	NoParentCgroupWrites           bool  `json:"noParentCgroupWrites"`
	NoArbitraryFileWrites          bool  `json:"noArbitraryFileWrites"`
	EnforceCPUMaxBounds            bool  `json:"enforceCpuMaxBounds"`
	EnforceTTL                     bool  `json:"enforceTtl"`
	AutomaticReturnToFloor         bool  `json:"automaticReturnToFloor"`
	KillSwitchEnabled              bool  `json:"killSwitchEnabled"`
	AuditLogRequired               bool  `json:"auditLogRequired"`
	DryRunModeSupported            bool  `json:"dryRunModeSupported"`
	PanicSafeRollback              bool  `json:"panicSafeRollback"`
	RejectStaleRequests            bool  `json:"rejectStaleRequests"`
	ReplayProtectionFutureSignable bool  `json:"replayProtectionFutureSignable"`
	MinQuota                       int64 `json:"minQuota"`
	MaxQuota                       int64 `json:"maxQuota"`
	PeriodUsec                     int64 `json:"periodUsec"`
}

func DefaultNodeFluidActuatorSafetyModel() NodeFluidActuatorSafetyModel {
	return NodeFluidActuatorSafetyModel{
		NodeLocalBlastRadiusOnly:       true,
		TargetAllowlistRequired:        true,
		ValidateCgroupPath:             true,
		ValidatePIDAndStartTime:        true,
		ValidatePodUID:                 true,
		ValidateVMUID:                  true,
		NoSymlinkTraversal:             true,
		NoParentCgroupWrites:           true,
		NoArbitraryFileWrites:          true,
		EnforceCPUMaxBounds:            true,
		EnforceTTL:                     true,
		AutomaticReturnToFloor:         true,
		KillSwitchEnabled:              true,
		AuditLogRequired:               true,
		DryRunModeSupported:            true,
		PanicSafeRollback:              true,
		RejectStaleRequests:            true,
		ReplayProtectionFutureSignable: true,
		MinQuota:                       100000,
		MaxQuota:                       800000,
		PeriodUsec:                     100000,
	}
}

type NodeFluidActuatorSafetyInput struct {
	Contract          KARLNodeFluidActuatorContract `json:"contract"`
	ExpectedVMUID     string                        `json:"expectedVmUid"`
	ObservedVMUID     string                        `json:"observedVmUid"`
	ExpectedPodUID    string                        `json:"expectedPodUid"`
	ObservedPodUID    string                        `json:"observedPodUid"`
	ExpectedQemuPID   string                        `json:"expectedQemuPid"`
	ObservedQemuPID   string                        `json:"observedQemuPid"`
	ExpectedQemuStart string                        `json:"expectedQemuStart"`
	ObservedQemuStart string                        `json:"observedQemuStart"`
	ExpectedPrefix    string                        `json:"expectedPrefix"`
	RequestCreatedAt  time.Time                     `json:"requestCreatedAt"`
	EvaluationTime    time.Time                     `json:"evaluationTime"`
	TTLSeconds        int64                         `json:"ttlSeconds"`
	SeenRequestIDs    []string                      `json:"seenRequestIds"`
	WriteTarget       string                        `json:"writeTarget"`
	AllowParentWrite  bool                          `json:"allowParentWrite"`
	ResolvedRealPath  string                        `json:"resolvedRealPath"`
	RequestedQuota    int64                         `json:"requestedQuota"`
}

type NodeFluidActuatorSafetyResult struct {
	Allowed  bool     `json:"allowed"`
	Blockers []string `json:"blockers"`
}

func EvaluateNodeFluidActuatorSafety(input NodeFluidActuatorSafetyInput, model NodeFluidActuatorSafetyModel) NodeFluidActuatorSafetyResult {
	blockers := ValidateKARLNodeFluidActuatorContract(input.Contract)
	if model.ValidateVMUID && input.ExpectedVMUID != "" && input.ObservedVMUID != input.ExpectedVMUID {
		blockers = append(blockers, BlockerActuatorTargetAmbiguous)
	}
	if model.ValidatePodUID && input.ExpectedPodUID != "" && input.ObservedPodUID != input.ExpectedPodUID {
		blockers = append(blockers, BlockerActuatorTargetAmbiguous)
	}
	if model.ValidatePIDAndStartTime {
		if input.ExpectedQemuPID != "" && input.ObservedQemuPID != input.ExpectedQemuPID {
			blockers = append(blockers, BlockerQemuPIDChanged)
		}
		if input.ExpectedQemuStart != "" && input.ObservedQemuStart != input.ExpectedQemuStart {
			blockers = append(blockers, BlockerQemuPIDChanged)
		}
	}
	if model.ValidateCgroupPath && input.ExpectedPrefix != "" && !hasPrefix(input.Contract.CgroupPath, input.ExpectedPrefix) {
		blockers = append(blockers, BlockerCgroupPathMismatch)
	}
	if model.NoSymlinkTraversal && input.ResolvedRealPath != "" && !hasPrefix(input.ResolvedRealPath, input.ExpectedPrefix) {
		blockers = append(blockers, BlockerActuatorPathEscape)
	}
	if model.NoParentCgroupWrites && !input.AllowParentWrite && input.WriteTarget == "parent" {
		blockers = append(blockers, BlockerActuatorParentCgroupWrite)
	}
	if model.NoArbitraryFileWrites && input.WriteTarget != "" && input.WriteTarget != "cpu.max" && input.WriteTarget != "cpu.weight" {
		blockers = append(blockers, BlockerActuatorArbitraryWrite)
	}
	if model.EnforceCPUMaxBounds {
		if input.RequestedQuota < model.MinQuota || input.RequestedQuota > model.MaxQuota {
			blockers = append(blockers, BlockerActuatorCPUEntitlementOutOfBounds)
		}
	}
	if model.EnforceTTL && input.TTLSeconds > 0 && !input.RequestCreatedAt.IsZero() && !input.EvaluationTime.IsZero() {
		if input.EvaluationTime.Sub(input.RequestCreatedAt) > time.Duration(input.TTLSeconds)*time.Second {
			blockers = append(blockers, BlockerStaleActuatorRequest)
		}
	}
	if model.ReplayProtectionFutureSignable {
		for _, seen := range input.SeenRequestIDs {
			if seen == input.Contract.RequestID {
				blockers = append(blockers, BlockerActuatorReplayDetected)
				break
			}
		}
	}
	blockers = dedupe(blockers)
	return NodeFluidActuatorSafetyResult{Allowed: len(blockers) == 0, Blockers: blockers}
}

func dedupeRemediations(values []WindowsHyperdensityRemediationID) []WindowsHyperdensityRemediationID {
	seen := make(map[WindowsHyperdensityRemediationID]struct{}, len(values))
	result := make([]WindowsHyperdensityRemediationID, 0, len(values))
	for _, value := range values {
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}

func remediationsToStrings(values []WindowsHyperdensityRemediationID) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		result = append(result, string(value))
	}
	return dedupe(result)
}

func hasPrefix(path string, prefix string) bool {
	if prefix == "" {
		return false
	}
	if path == prefix {
		return true
	}
	if len(path) < len(prefix) {
		return false
	}
	return path[:len(prefix)] == prefix
}
