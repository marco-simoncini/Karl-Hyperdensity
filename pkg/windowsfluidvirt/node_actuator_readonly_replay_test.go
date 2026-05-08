package windowsfluidvirt

import "testing"

func TestReadonlyReplayDefaults(t *testing.T) {
	replay := NewWindowsFluidVirtNodeActuatorReadonlyReplayMinimal()

	if replay.ReplayID != "windows_fluidvirt_node_actuator_mvp_readonly_replay_v1" {
		t.Fatalf("unexpected replay id: %s", replay.ReplayID)
	}
	if replay.ReplayMode != "readonly_replay_only" {
		t.Fatalf("unexpected replay mode: %s", replay.ReplayMode)
	}
	if !replay.NodeActuatorMvpReplayAvailable {
		t.Fatalf("readonly replay availability must be true")
	}
	if replay.NodeActuatorMvpPortedAsRuntime || replay.ActuatorRuntimeEnabled || replay.RuntimeMutationEnabled {
		t.Fatalf("runtime must remain disabled")
	}
	if replay.CgroupWriteEnabled || replay.CPUMaxMutationEnabled {
		t.Fatalf("cgroup write/cpu.max mutation must remain disabled")
	}
	if replay.QMPCommandExecutionAllowed || replay.QGACommandExecutionAllowed {
		t.Fatalf("qmp/qga execution must remain disabled")
	}
	if replay.RawRuntimeControlsExposed {
		t.Fatalf("raw runtime controls must remain disabled")
	}
	if replay.AutonomousApplyAllowed || replay.ProductionMutationAllowed {
		t.Fatalf("autonomous/production mutation must remain disabled")
	}
	if replay.WindowsGaClaimAllowed || replay.WindowsProductionReadyClaimAllowed || replay.WindowsExecutionReadyByDefault {
		t.Fatalf("windows ga/production/execution ready claims must remain disabled")
	}
	if replay.VCPUHotplugClaimAllowed || replay.LogicalCPUScalingClaimAllowed || replay.PoolScalingClaimAllowed {
		t.Fatalf("vcpu/logical/pool claims must remain disabled")
	}
	if replay.LiveMigrationClaimAllowed || replay.RebootRecreateRolloutMechanismAllowed {
		t.Fatalf("migration/reboot/recreate/rollout claims must remain disabled")
	}
	if replay.ReplayValidationResult.ValidationState != "passed" || replay.ReplayValidationResult.RuntimeMVPReady {
		t.Fatalf("validation must pass while runtime MVP remains not ready")
	}
}

func TestReadonlyReplayContainsExpectedModelOnlyData(t *testing.T) {
	replay := NewWindowsFluidVirtNodeActuatorReadonlyReplayMinimal()

	if len(replay.ReplayedCPUEntitlementTransitions) < 6 {
		t.Fatalf("expected replayed cpu transitions")
	}
	for _, tr := range replay.ReplayedCPUEntitlementTransitions {
		if tr.TransitionState != TransitionReplayedModelOnly || tr.AppliedToRealCgroup || !tr.FakeRuntimeOnly {
			t.Fatalf("transition %s must stay model-only and fake-runtime only", tr.TransitionID)
		}
	}
	if len(replay.ReplayedSafetyChecks) == 0 {
		t.Fatalf("expected replayed safety checks")
	}
	if len(replay.ReplayedBlockers) == 0 {
		t.Fatalf("expected replayed blockers")
	}
	if len(replay.ReplayedAuditEvents) == 0 {
		t.Fatalf("expected replayed audit events")
	}
	for _, event := range replay.ReplayedAuditEvents {
		if event.ContainsSecretMaterial {
			t.Fatalf("audit event %s must not contain secret material", event.EventID)
		}
	}
	if len(replay.ForbiddenReplayActions) == 0 {
		t.Fatalf("expected forbidden replay actions")
	}
}
