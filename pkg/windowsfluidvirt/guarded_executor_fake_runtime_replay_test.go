package windowsfluidvirt

import "testing"

func TestGuardedExecutorFakeRuntimeReplayDefaults(t *testing.T) {
	replay := NewWindowsFluidVirtGuardedExecutorFakeRuntimeReplayMinimal()
	if !replay.ExecutorFakeRuntimeReplayAvailable || !replay.ExecutorFakeRuntimeReplayExecuted || !replay.ExecutorFakeRuntimeReplayDeterministic {
		t.Fatalf("fake-runtime replay availability/executed/deterministic flags must be true")
	}
	if replay.ExecutorEnabled || replay.ExecutorRuntimeAvailable || replay.ExecutorExecuted {
		t.Fatalf("executor must remain disabled and unavailable")
	}
	if replay.ControlledApplyEnabled || replay.ControlledApplyExecuted || replay.ControlledApplyReady {
		t.Fatalf("controlled apply must remain disabled/not executed/not ready")
	}
	if replay.RuntimeMutationEnabled || replay.CgroupWriteEnabled || replay.CPUMaxMutationEnabled {
		t.Fatalf("runtime mutation/cgroup/cpu mutation must remain disabled")
	}
	if replay.QMPCommandExecutionAllowed || replay.QGACommandExecutionAllowed || replay.RawRuntimeControlsExposed {
		t.Fatalf("qmp/qga/raw controls must remain disabled")
	}
	if replay.AutonomousApplyAllowed || replay.ProductionMutationAllowed {
		t.Fatalf("autonomous/production mutation must remain disabled")
	}
	if replay.WindowsGaClaimAllowed || replay.WindowsProductionReadyClaimAllowed || replay.WindowsExecutionReadyByDefault {
		t.Fatalf("windows ga/production/execution-ready claims must remain disabled")
	}
	if err := ValidateWindowsFluidVirtGuardedExecutorFakeRuntimeReplay(replay); err != nil {
		t.Fatalf("validate replay: %v", err)
	}
	if _, err := ReplayGuardedExecutorFakeRuntime(replay); err != nil {
		t.Fatalf("replay execution should pass in fake-runtime mode: %v", err)
	}
}

func TestGuardedExecutorFakeRuntimeReplayCatalogs(t *testing.T) {
	replay := NewWindowsFluidVirtGuardedExecutorFakeRuntimeReplayMinimal()
	if len(replay.ReplayedExecutorSteps) == 0 || len(replay.ReplayedInputValidations) == 0 || len(replay.ReplayedOutputValidations) == 0 {
		t.Fatalf("replayed steps and validations must be populated")
	}
	if len(replay.ReplayedGateEvaluations) == 0 || len(replay.ReplayedBlockingReasons) == 0 || len(replay.ReplayedAuditEvents) == 0 {
		t.Fatalf("gate evaluations, blocking reasons and audit events must be populated")
	}
	if replay.ExecutorFakeRuntimeReplayValidation.ValidationState != "passed" {
		t.Fatalf("validation result must be passed")
	}
	requiredForbidden := []string{
		"enable_executor",
		"execute_windows_executor_runtime",
		"execute_controlled_apply",
		"enable_controlled_apply",
		"execute_node_actuator_runtime",
		"write_cgroup_cpu_max",
		"execute_qmp_command",
		"execute_qga_command",
		"enable_runtime_actuator",
		"enable_autonomous_apply",
		"enable_production_auto",
		"mark_windows_execution_ready",
		"claim_windows_ga",
		"claim_windows_production_ready",
		"expose_raw_runtime_controls",
	}
	for _, action := range requiredForbidden {
		if !containsString(replay.ForbiddenExecutorFakeRuntimeReplayActions, action) {
			t.Fatalf("missing forbidden replay action: %s", action)
		}
	}
}
