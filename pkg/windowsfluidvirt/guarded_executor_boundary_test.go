package windowsfluidvirt

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestGuardedExecutorBoundaryDefaults(t *testing.T) {
	boundary := NewWindowsFluidVirtGuardedExecutorBoundaryMinimal()
	if !boundary.ExecutorBoundaryDefined {
		t.Fatalf("executor boundary must be defined")
	}
	if boundary.ExecutorEnabled || boundary.ExecutorRuntimeAvailable || boundary.ExecutorExecuted {
		t.Fatalf("executor must remain disabled and unavailable")
	}
	if boundary.ControlledApplyEnabled || boundary.ControlledApplyReady {
		t.Fatalf("controlled apply must remain disabled/not ready")
	}
	if boundary.RuntimeMutationEnabled || boundary.ActuatorRuntimeEnabled || boundary.CgroupWriteEnabled {
		t.Fatalf("runtime mutation and cgroup write must remain disabled")
	}
	if boundary.QMPCommandExecutionAllowed || boundary.QGACommandExecutionAllowed {
		t.Fatalf("qmp/qga execution must remain disabled")
	}
	if boundary.AutonomousApplyAllowed || boundary.ProductionMutationAllowed {
		t.Fatalf("autonomous and production mutation must remain disabled")
	}
	if boundary.WindowsGaClaimAllowed || boundary.WindowsProductionReadyClaimAllowed || boundary.WindowsExecutionReadyByDefault {
		t.Fatalf("windows ga/production/execution-ready claims must remain disabled")
	}
	if boundary.ExecutorReadinessScorecard.ScorecardState != "executor_boundary_defined_not_ready" {
		t.Fatalf("unexpected scorecard state: %s", boundary.ExecutorReadinessScorecard.ScorecardState)
	}
	if err := ValidateWindowsFluidVirtGuardedExecutorBoundary(boundary); err != nil {
		t.Fatalf("validate boundary: %v", err)
	}
}

func TestGuardedExecutorInputOutputBoundariesForbidRawControls(t *testing.T) {
	boundary := NewWindowsFluidVirtGuardedExecutorBoundaryMinimal()
	forbiddenInput := []string{"raw_cgroup_path", "raw_qmp_command", "raw_qga_command", "raw_secret", "token"}
	for _, item := range forbiddenInput {
		if !containsString(boundary.ExecutorInputBoundary.ForbiddenInputTypes, item) {
			t.Fatalf("missing forbidden input type: %s", item)
		}
	}
	forbiddenOutput := []string{"runtime_apply_result", "cgroup_write_result", "qmp_command_result", "qga_command_result", "raw_secret_output"}
	for _, item := range forbiddenOutput {
		if !containsString(boundary.ExecutorOutputBoundary.ForbiddenOutputTypes, item) {
			t.Fatalf("missing forbidden output type: %s", item)
		}
	}
	if boundary.ExecutorInputBoundary.RealRuntimeCommandInputAllowed || boundary.ExecutorInputBoundary.RawCgroupPathInputAllowed || boundary.ExecutorInputBoundary.RawQMPCommandInputAllowed || boundary.ExecutorInputBoundary.RawQGACommandInputAllowed {
		t.Fatalf("input boundary must forbid raw runtime command inputs")
	}
	if boundary.ExecutorOutputBoundary.RuntimeMutationOutputAllowed || boundary.ExecutorOutputBoundary.RawRuntimeControlOutputAllowed || boundary.ExecutorOutputBoundary.SecretOutputAllowed {
		t.Fatalf("output boundary must forbid mutation and secret outputs")
	}
}

func TestGuardedExecutorBoundaryListsPresent(t *testing.T) {
	boundary := NewWindowsFluidVirtGuardedExecutorBoundaryMinimal()
	if len(boundary.RequiredBeforeExecutorEnablement) == 0 {
		t.Fatalf("requiredBeforeExecutorEnablement must be populated")
	}
	if len(boundary.ExecutorBlockingReasons) == 0 {
		t.Fatalf("executorBlockingReasons must be populated")
	}
	requiredForbidden := []string{
		"enable_executor",
		"execute_windows_executor",
		"enable_controlled_apply",
		"execute_controlled_apply",
		"execute_node_actuator_runtime",
		"write_cgroup_cpu_max",
		"execute_qmp_command",
		"execute_qga_command",
		"enable_autonomous_apply",
		"enable_production_auto",
		"mark_windows_execution_ready",
		"claim_windows_ga",
		"claim_windows_production_ready",
		"expose_raw_runtime_controls",
	}
	for _, action := range requiredForbidden {
		if !containsString(boundary.ForbiddenExecutorBoundaryActions, action) {
			t.Fatalf("missing forbidden boundary action: %s", action)
		}
	}
}

func TestLoadGuardedExecutorBoundaryFixtureFromTemporaryFile(t *testing.T) {
	boundary := NewWindowsFluidVirtGuardedExecutorBoundaryMinimal()
	path := filepath.Join(t.TempDir(), "guarded_executor_boundary.json")
	raw, err := json.Marshal(boundary)
	if err != nil {
		t.Fatalf("marshal boundary: %v", err)
	}
	if err := os.WriteFile(path, raw, 0o600); err != nil {
		t.Fatalf("write boundary fixture: %v", err)
	}
	loaded, err := LoadGuardedExecutorBoundaryFixtureFromTemporaryFile(path)
	if err != nil {
		t.Fatalf("load boundary fixture: %v", err)
	}
	if loaded.ExecutorEnabled || loaded.ExecutorRuntimeAvailable {
		t.Fatalf("loaded boundary must keep executor disabled")
	}
}
