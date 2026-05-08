package windowsfluidvirt

import "testing"

func TestNodeActuatorContractBoundaryDefaults(t *testing.T) {
	contract := NewWindowsFluidVirtNodeActuatorContractBoundary()

	if contract.ActuatorContractID != "windows_fluidvirt_node_actuator_contract_boundary_v1" {
		t.Fatalf("unexpected contract id: %s", contract.ActuatorContractID)
	}
	if contract.NodeActuatorMvpPorted {
		t.Fatalf("node actuator MVP must not be ported in boundary milestone")
	}
	if contract.ActuatorRuntimeEnabled || contract.RuntimeMutationEnabled {
		t.Fatalf("runtime must stay disabled")
	}
	if contract.CPUBoundary.CgroupWriteEnabled || contract.CPUBoundary.CPUMaxMutationEnabled {
		t.Fatalf("cgroup/cpu.max mutation must stay disabled")
	}
	if contract.RAMBoundary.QMPCommandExecutionAllowed {
		t.Fatalf("qmp command execution must stay disabled")
	}
	if contract.RawRuntimeControlsExposed || contract.RAMBoundary.RawQMPControlExposed {
		t.Fatalf("raw controls must stay disabled")
	}
	if contract.AutonomousApplyAllowed || contract.ProductionMutationAllowed {
		t.Fatalf("autonomous and production mutation must be disabled")
	}
	if contract.WindowsGaClaimAllowed || contract.WindowsProductionReadyClaimAllowed || contract.WindowsExecutionReadyByDefault {
		t.Fatalf("windows ga/production/execution-ready claims must be disabled")
	}
	if contract.VCPUHotplugClaimAllowed || contract.LogicalCPUScalingClaimAllowed || contract.PoolScalingClaimAllowed {
		t.Fatalf("vcpu/logical/pool scaling claims must be disabled")
	}
	if contract.LiveMigrationClaimAllowed || contract.RebootRecreateRolloutMechanismAllowed {
		t.Fatalf("migration and reboot/recreate/rollout claims must be disabled")
	}
	if contract.EnforcementMode != "disabled" {
		t.Fatalf("enforcement mode must be disabled")
	}
}

func TestNodeActuatorContractBoundaryCatalogs(t *testing.T) {
	contract := NewWindowsFluidVirtNodeActuatorContractBoundary()

	requiredGates := []WindowsFluidVirtNodeActuatorGate{
		GateNodeAllowlistDefined,
		GateCgroupV2CPUMaxPathValidated,
		GateHostScopeNodeLocalConfirmed,
		GateManualApprovalRequired,
		GateActiveLeaseRequired,
		GateLeaseTTLRequired,
		GateReturnToFloorPlanRequired,
		GateRollbackPlanRequired,
		GateGuestWitnessRequired,
		GateSameBootProofRequired,
		GateSameQEMUProofRequired,
		GateAuditHashChainRequired,
		GateKillSwitchRequired,
		GateNoRawControlExposure,
		GateNoAutonomousApply,
		GateNoProductionApply,
		GateNoWindowsGAClaim,
	}
	for _, gate := range requiredGates {
		if !containsNodeActuatorGate(contract.RequiredGates, gate) {
			t.Fatalf("missing gate: %s", gate)
		}
	}

	requiredBlockers := []WindowsFluidVirtNodeActuatorBlocker{
		NodeActuatorRuntimeNotPorted,
		NodeActuatorCgroupWritePathNotEnabled,
		NodeActuatorNodeAllowlistMissing,
		NodeActuatorManualApprovalMissing,
		NodeActuatorLeaseTTLMissing,
		NodeActuatorReturnToFloorPlanMissing,
		NodeActuatorRollbackPlanMissing,
		NodeActuatorGuestWitnessMissing,
		NodeActuatorSameBootProofMissing,
		NodeActuatorSameQEMUProofMissing,
		NodeActuatorAuditChainMissing,
		NodeActuatorKillSwitchMissing,
		NodeActuatorProductionReadyForbidden,
		NodeActuatorAutonomousApplyForbidden,
		NodeActuatorRawRuntimeControlForbidden,
	}
	for _, blocker := range requiredBlockers {
		if !containsNodeActuatorBlocker(contract.Blockers, blocker) {
			t.Fatalf("missing blocker: %s", blocker)
		}
	}

	requiredForbiddenActions := []WindowsFluidVirtNodeActuatorForbiddenAction{
		ForbiddenExecuteNodeActuator,
		ForbiddenWriteCgroupCPUMax,
		ForbiddenMutateQMPBalloon,
		ForbiddenExecuteQMPCommand,
		ForbiddenExecuteQGACommand,
		ForbiddenEnableAutonomousApply,
		ForbiddenEnableProductionAUTO,
		ForbiddenMarkExecutionReady,
		ForbiddenClaimWindowsGA,
		ForbiddenClaimWindowsProdReady,
		ForbiddenExposeRawRuntimeControl,
		ForbiddenCreateSecret,
		ForbiddenCommitSecret,
		ForbiddenLogToken,
		ForbiddenPackageOSISO,
		ForbiddenTouchDashboard,
		ForbiddenTouchInventory,
	}
	for _, action := range requiredForbiddenActions {
		if !containsNodeActuatorForbiddenAction(contract.ForbiddenActions, action) {
			t.Fatalf("missing forbidden action: %s", action)
		}
	}
}

func containsNodeActuatorGate(values []WindowsFluidVirtNodeActuatorGate, target WindowsFluidVirtNodeActuatorGate) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func containsNodeActuatorBlocker(values []WindowsFluidVirtNodeActuatorBlocker, target WindowsFluidVirtNodeActuatorBlocker) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func containsNodeActuatorForbiddenAction(values []WindowsFluidVirtNodeActuatorForbiddenAction, target WindowsFluidVirtNodeActuatorForbiddenAction) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
