package windowsfluidvirt

import "sort"

type WindowsFluidVirtBlocker string

const (
	BlockerMissingGuestWitness            WindowsFluidVirtBlocker = "missing_guest_witness"
	BlockerMissingSameBootProof           WindowsFluidVirtBlocker = "missing_same_boot_proof"
	BlockerMissingReturnToFloorPlan       WindowsFluidVirtBlocker = "missing_return_to_floor_plan"
	BlockerMissingRollbackPlan            WindowsFluidVirtBlocker = "missing_rollback_plan"
	BlockerMissingAuditChain              WindowsFluidVirtBlocker = "missing_audit_chain"
	BlockerNoManualApproval               WindowsFluidVirtBlocker = "no_manual_approval"
	BlockerRawRuntimeControlForbidden     WindowsFluidVirtBlocker = "raw_runtime_control_forbidden"
	BlockerAutonomousApplyForbidden       WindowsFluidVirtBlocker = "autonomous_apply_forbidden"
	BlockerProductionReadyClaimForbidden  WindowsFluidVirtBlocker = "production_ready_claim_forbidden"
	BlockerVCPUHotplugClaimForbidden      WindowsFluidVirtBlocker = "vcpu_hotplug_claim_forbidden"
	BlockerLogicalCPUScalingClaimForbidden WindowsFluidVirtBlocker = "logical_cpu_scaling_claim_forbidden"
	BlockerPoolScalingClaimForbidden      WindowsFluidVirtBlocker = "pool_scaling_claim_forbidden"
)

type WindowsFluidVirtBlockerDefinition struct {
	ID          WindowsFluidVirtBlocker `json:"id"`
	Severity    string                  `json:"severity"`
	Description string                  `json:"description"`
}

var windowsFluidVirtBlockerCatalog = map[WindowsFluidVirtBlocker]WindowsFluidVirtBlockerDefinition{
	BlockerMissingGuestWitness: {
		ID:          BlockerMissingGuestWitness,
		Severity:    "high",
		Description: "Guest witness evidence is missing from the Windows FluidVirt product model.",
	},
	BlockerMissingSameBootProof: {
		ID:          BlockerMissingSameBootProof,
		Severity:    "high",
		Description: "Same-boot continuity proof is missing.",
	},
	BlockerMissingReturnToFloorPlan: {
		ID:          BlockerMissingReturnToFloorPlan,
		Severity:    "critical",
		Description: "Return-to-floor plan is required before any apply progression.",
	},
	BlockerMissingRollbackPlan: {
		ID:          BlockerMissingRollbackPlan,
		Severity:    "critical",
		Description: "Rollback plan is required before any apply progression.",
	},
	BlockerMissingAuditChain: {
		ID:          BlockerMissingAuditChain,
		Severity:    "high",
		Description: "Audit chain reference is missing.",
	},
	BlockerNoManualApproval: {
		ID:          BlockerNoManualApproval,
		Severity:    "high",
		Description: "Manual approval is required for gated preview progression.",
	},
	BlockerRawRuntimeControlForbidden: {
		ID:          BlockerRawRuntimeControlForbidden,
		Severity:    "critical",
		Description: "Raw runtime controls are forbidden in product surface.",
	},
	BlockerAutonomousApplyForbidden: {
		ID:          BlockerAutonomousApplyForbidden,
		Severity:    "critical",
		Description: "Autonomous apply is forbidden in this milestone.",
	},
	BlockerProductionReadyClaimForbidden: {
		ID:          BlockerProductionReadyClaimForbidden,
		Severity:    "critical",
		Description: "Production-ready claim is forbidden.",
	},
	BlockerVCPUHotplugClaimForbidden: {
		ID:          BlockerVCPUHotplugClaimForbidden,
		Severity:    "critical",
		Description: "vCPU hotplug claim is forbidden.",
	},
	BlockerLogicalCPUScalingClaimForbidden: {
		ID:          BlockerLogicalCPUScalingClaimForbidden,
		Severity:    "critical",
		Description: "Logical CPU scaling claim is forbidden.",
	},
	BlockerPoolScalingClaimForbidden: {
		ID:          BlockerPoolScalingClaimForbidden,
		Severity:    "critical",
		Description: "Pool scaling claim is forbidden.",
	},
}

func WindowsFluidVirtBlockerCatalog() []WindowsFluidVirtBlockerDefinition {
	keys := make([]string, 0, len(windowsFluidVirtBlockerCatalog))
	for blocker := range windowsFluidVirtBlockerCatalog {
		keys = append(keys, string(blocker))
	}
	sort.Strings(keys)
	result := make([]WindowsFluidVirtBlockerDefinition, 0, len(keys))
	for _, key := range keys {
		result = append(result, windowsFluidVirtBlockerCatalog[WindowsFluidVirtBlocker(key)])
	}
	return result
}

func IsKnownWindowsFluidVirtBlocker(blocker WindowsFluidVirtBlocker) bool {
	_, ok := windowsFluidVirtBlockerCatalog[blocker]
	return ok
}
