package windowsfluidvirt

type BlockerSeverity string

const (
	SeverityCritical BlockerSeverity = "critical"
	SeverityHigh     BlockerSeverity = "high"
	SeverityMedium   BlockerSeverity = "medium"
)

type BlockerCategory string

const (
	CategoryQMP         BlockerCategory = "qmp"
	CategoryGuest       BlockerCategory = "guest"
	CategoryIdentity    BlockerCategory = "identity"
	CategoryMigration   BlockerCategory = "migration"
	CategorySafety      BlockerCategory = "safety"
	CategoryEnvironment BlockerCategory = "environment"
	CategoryIntegration BlockerCategory = "integration"
	CategoryExecution   BlockerCategory = "execution_safety"
)

type ResultingPhase string

const (
	PhaseBlocked     ResultingPhase = "BLOCKED"
	PhaseQuarantined ResultingPhase = "QUARANTINED"
)

type BlockerDefinition struct {
	ID                  string
	Severity            BlockerSeverity
	Category            BlockerCategory
	Message             string
	Remediable          bool
	ResultingPhase      ResultingPhase
	EvidenceRequirement string
}

const (
	BlockerQMPSocketUnavailable               = "qmp_socket_unavailable"
	BlockerGuestAgentUnavailable              = "guest_agent_unavailable"
	BlockerKarlAgentFluidModuleMissing        = "karl_agent_fluid_module_missing"
	BlockerPendingRebootDetected              = "pending_reboot_detected"
	BlockerQemuPIDChanged                     = "qemu_pid_changed"
	BlockerLastBootChanged                    = "last_boot_changed"
	BlockerMachineGUIDChanged                 = "machine_guid_changed"
	BlockerLiveMigrationRequired              = "live_migration_required"
	BlockerVMIRecreateRequired                = "vmi_recreate_required"
	BlockerVirtLauncherPodChanged             = "virt_launcher_pod_changed"
	BlockerNodeChanged                        = "node_changed"
	BlockerMemoryDriverUnverified             = "memory_driver_unverified"
	BlockerMemoryReturnNotSafe                = "memory_return_not_safe"
	BlockerCPUTopologyNotConfirmed            = "cpu_topology_not_confirmed"
	BlockerGuestMemoryNotConfirmed            = "guest_memory_not_confirmed"
	BlockerRollbackNotReady                   = "rollback_not_ready"
	BlockerReturnToFloorNotReady              = "return_to_floor_not_ready"
	BlockerQMPAckMissing                      = "qmp_ack_missing"
	BlockerGuestAckMissing                    = "guest_ack_missing"
	BlockerHotplugErrorDetected               = "hotplug_error_detected"
	BlockerCriticalWindowsEventDetected       = "critical_windows_event_detected"
	BlockerDashboard443TouchRisk              = "dashboard_443_touch_risk"
	BlockerCandidate8888Unavailable           = "candidate_8888_unavailable"
	BlockerWindowsAgentRepoNotPresentInTarget = "windows_agent_repo_not_present_in_target_repos"
	BlockerFutureApplyExecutorDisabled        = "future_apply_executor_disabled"
)

var CanonicalBlockers = map[string]BlockerDefinition{
	BlockerQMPSocketUnavailable: {
		ID:                  BlockerQMPSocketUnavailable,
		Severity:            SeverityCritical,
		Category:            CategoryQMP,
		Message:             "QMP socket is unavailable for runtime continuity checks.",
		Remediable:          true,
		ResultingPhase:      PhaseBlocked,
		EvidenceRequirement: "qmp socket probe evidence",
	},
	BlockerGuestAgentUnavailable: {
		ID:                  BlockerGuestAgentUnavailable,
		Severity:            SeverityCritical,
		Category:            CategoryGuest,
		Message:             "Guest agent is unavailable for Windows fluid preflight.",
		Remediable:          true,
		ResultingPhase:      PhaseBlocked,
		EvidenceRequirement: "guest agent health evidence",
	},
	BlockerKarlAgentFluidModuleMissing: {
		ID:                  BlockerKarlAgentFluidModuleMissing,
		Severity:            SeverityCritical,
		Category:            CategoryIntegration,
		Message:             "Windows KARL Agent module fluidShell is missing.",
		Remediable:          true,
		ResultingPhase:      PhaseBlocked,
		EvidenceRequirement: "agent module registry evidence",
	},
	BlockerPendingRebootDetected: {
		ID:                  BlockerPendingRebootDetected,
		Severity:            SeverityCritical,
		Category:            CategoryGuest,
		Message:             "Pending reboot detected in Windows guest.",
		Remediable:          true,
		ResultingPhase:      PhaseBlocked,
		EvidenceRequirement: "pending reboot probe evidence",
	},
	BlockerQemuPIDChanged: {
		ID:                  BlockerQemuPIDChanged,
		Severity:            SeverityCritical,
		Category:            CategoryIdentity,
		Message:             "QEMU PID changed; runtime continuity is broken.",
		Remediable:          false,
		ResultingPhase:      PhaseQuarantined,
		EvidenceRequirement: "before/after qemu pid evidence",
	},
	BlockerLastBootChanged: {
		ID:                  BlockerLastBootChanged,
		Severity:            SeverityCritical,
		Category:            CategoryIdentity,
		Message:             "Windows last boot changed; no-reboot proof failed.",
		Remediable:          false,
		ResultingPhase:      PhaseQuarantined,
		EvidenceRequirement: "before/after last boot evidence",
	},
	BlockerMachineGUIDChanged: {
		ID:                  BlockerMachineGUIDChanged,
		Severity:            SeverityCritical,
		Category:            CategoryIdentity,
		Message:             "Machine GUID changed; machine identity continuity failed.",
		Remediable:          false,
		ResultingPhase:      PhaseQuarantined,
		EvidenceRequirement: "before/after machine guid evidence",
	},
	BlockerLiveMigrationRequired: {
		ID:                  BlockerLiveMigrationRequired,
		Severity:            SeverityCritical,
		Category:            CategoryMigration,
		Message:             "Live migration is required or active; in-place mode cannot proceed.",
		Remediable:          true,
		ResultingPhase:      PhaseBlocked,
		EvidenceRequirement: "vmim/livemigration state evidence",
	},
	BlockerVMIRecreateRequired: {
		ID:                  BlockerVMIRecreateRequired,
		Severity:            SeverityCritical,
		Category:            CategoryMigration,
		Message:             "VMI recreate is required; same-runtime guarantee failed.",
		Remediable:          true,
		ResultingPhase:      PhaseBlocked,
		EvidenceRequirement: "vmi condition evidence",
	},
	BlockerVirtLauncherPodChanged: {
		ID:                  BlockerVirtLauncherPodChanged,
		Severity:            SeverityCritical,
		Category:            CategoryIdentity,
		Message:             "virt-launcher pod identity changed.",
		Remediable:          false,
		ResultingPhase:      PhaseQuarantined,
		EvidenceRequirement: "before/after virt-launcher pod evidence",
	},
	BlockerNodeChanged: {
		ID:                  BlockerNodeChanged,
		Severity:            SeverityCritical,
		Category:            CategoryIdentity,
		Message:             "Node identity changed; same-node guarantee failed.",
		Remediable:          false,
		ResultingPhase:      PhaseQuarantined,
		EvidenceRequirement: "before/after node evidence",
	},
	BlockerMemoryDriverUnverified: {
		ID:                  BlockerMemoryDriverUnverified,
		Severity:            SeverityHigh,
		Category:            CategorySafety,
		Message:             "Guest memory driver is not verified.",
		Remediable:          true,
		ResultingPhase:      PhaseBlocked,
		EvidenceRequirement: "guest driver truth evidence",
	},
	BlockerMemoryReturnNotSafe: {
		ID:                  BlockerMemoryReturnNotSafe,
		Severity:            SeverityCritical,
		Category:            CategorySafety,
		Message:             "Return-to-floor for memory is not safe.",
		Remediable:          true,
		ResultingPhase:      PhaseBlocked,
		EvidenceRequirement: "return-to-floor safety evidence",
	},
	BlockerCPUTopologyNotConfirmed: {
		ID:                  BlockerCPUTopologyNotConfirmed,
		Severity:            SeverityHigh,
		Category:            CategoryGuest,
		Message:             "CPU topology is not confirmed.",
		Remediable:          true,
		ResultingPhase:      PhaseBlocked,
		EvidenceRequirement: "guest cpu topology evidence",
	},
	BlockerGuestMemoryNotConfirmed: {
		ID:                  BlockerGuestMemoryNotConfirmed,
		Severity:            SeverityHigh,
		Category:            CategoryGuest,
		Message:             "Guest memory truth is not confirmed.",
		Remediable:          true,
		ResultingPhase:      PhaseBlocked,
		EvidenceRequirement: "guest memory evidence",
	},
	BlockerRollbackNotReady: {
		ID:                  BlockerRollbackNotReady,
		Severity:            SeverityCritical,
		Category:            CategorySafety,
		Message:             "Rollback readiness is not proven.",
		Remediable:          true,
		ResultingPhase:      PhaseBlocked,
		EvidenceRequirement: "rollback readiness evidence",
	},
	BlockerReturnToFloorNotReady: {
		ID:                  BlockerReturnToFloorNotReady,
		Severity:            SeverityCritical,
		Category:            CategorySafety,
		Message:             "Return-to-floor readiness is not proven.",
		Remediable:          true,
		ResultingPhase:      PhaseBlocked,
		EvidenceRequirement: "return-to-floor readiness evidence",
	},
	BlockerQMPAckMissing: {
		ID:                  BlockerQMPAckMissing,
		Severity:            SeverityCritical,
		Category:            CategoryQMP,
		Message:             "QMP acknowledgment is missing.",
		Remediable:          true,
		ResultingPhase:      PhaseBlocked,
		EvidenceRequirement: "qmp ack evidence",
	},
	BlockerGuestAckMissing: {
		ID:                  BlockerGuestAckMissing,
		Severity:            SeverityCritical,
		Category:            CategoryGuest,
		Message:             "Guest acknowledgment is missing.",
		Remediable:          true,
		ResultingPhase:      PhaseBlocked,
		EvidenceRequirement: "guest ack evidence",
	},
	BlockerHotplugErrorDetected: {
		ID:                  BlockerHotplugErrorDetected,
		Severity:            SeverityHigh,
		Category:            CategoryQMP,
		Message:             "Hotplug error detected in runtime evidence.",
		Remediable:          true,
		ResultingPhase:      PhaseBlocked,
		EvidenceRequirement: "runtime error evidence",
	},
	BlockerCriticalWindowsEventDetected: {
		ID:                  BlockerCriticalWindowsEventDetected,
		Severity:            SeverityHigh,
		Category:            CategoryGuest,
		Message:             "Critical Windows event detected.",
		Remediable:          true,
		ResultingPhase:      PhaseBlocked,
		EvidenceRequirement: "guest critical events evidence",
	},
	BlockerDashboard443TouchRisk: {
		ID:                  BlockerDashboard443TouchRisk,
		Severity:            SeverityCritical,
		Category:            CategoryEnvironment,
		Message:             "Operation risks touching dashboard on port 443.",
		Remediable:          true,
		ResultingPhase:      PhaseBlocked,
		EvidenceRequirement: "deployment plan isolation evidence",
	},
	BlockerCandidate8888Unavailable: {
		ID:                  BlockerCandidate8888Unavailable,
		Severity:            SeverityMedium,
		Category:            CategoryEnvironment,
		Message:             "Candidate dashboard on port 8888 is unavailable.",
		Remediable:          true,
		ResultingPhase:      PhaseBlocked,
		EvidenceRequirement: "candidate availability check",
	},
	BlockerWindowsAgentRepoNotPresentInTarget: {
		ID:                  BlockerWindowsAgentRepoNotPresentInTarget,
		Severity:            SeverityCritical,
		Category:            CategoryIntegration,
		Message:             "Windows agent repository is not present in target repositories.",
		Remediable:          true,
		ResultingPhase:      PhaseBlocked,
		EvidenceRequirement: "repository discovery evidence",
	},
	BlockerFutureApplyExecutorDisabled: {
		ID:                  BlockerFutureApplyExecutorDisabled,
		Severity:            SeverityCritical,
		Category:            CategoryExecution,
		Message:             "Future apply executor is hard-disabled in this phase.",
		Remediable:          false,
		ResultingPhase:      PhaseBlocked,
		EvidenceRequirement: "governance contract + disabled executor evidence",
	},
}

func LookupBlocker(id string) (BlockerDefinition, bool) {
	blocker, ok := CanonicalBlockers[id]
	return blocker, ok
}
