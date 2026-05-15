package windowsfluidvirt

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestControlledApplyPlanBoundaryDefaults(t *testing.T) {
	plan := NewWindowsFluidVirtControlledApplyPlanBoundaryMinimal()
	if !plan.ControlledApplyPlanDefined {
		t.Fatalf("controlled apply plan must be defined")
	}
	if plan.ControlledApplyEnabled || plan.ControlledApplyExecuted || plan.ControlledApplyReady {
		t.Fatalf("controlled apply must remain disabled/not executed/not ready")
	}
	if plan.ExecutorEnabled || plan.RuntimeMutationEnabled || plan.ActuatorRuntimeEnabled {
		t.Fatalf("executor/runtime mutation/actuator runtime must remain disabled")
	}
	if plan.CgroupWriteEnabled || plan.QMPCommandExecutionAllowed || plan.QGACommandExecutionAllowed {
		t.Fatalf("cgroup/qmp/qga execution must remain disabled")
	}
	if plan.AutonomousApplyAllowed || plan.ProductionMutationAllowed {
		t.Fatalf("autonomous/production mutation must remain disabled")
	}
	if plan.WindowsGaClaimAllowed || plan.WindowsProductionReadyClaimAllowed || plan.WindowsExecutionReadyByDefault {
		t.Fatalf("windows ga/production/execution-ready claims must remain disabled")
	}
	if plan.ControlledApplyReadinessScorecard.ScorecardState != "controlled_apply_plan_defined_not_ready" {
		t.Fatalf("unexpected scorecard state: %s", plan.ControlledApplyReadinessScorecard.ScorecardState)
	}
	if err := ValidateWindowsFluidVirtControlledApplyPlanBoundary(plan); err != nil {
		t.Fatalf("validate plan boundary: %v", err)
	}
}

func TestControlledApplyPlanBoundaryHasRequiredLists(t *testing.T) {
	plan := NewWindowsFluidVirtControlledApplyPlanBoundaryMinimal()
	if len(plan.RequiredBeforeControlledApply) == 0 {
		t.Fatalf("requiredBeforeControlledApply must be populated")
	}
	if len(plan.ControlledApplyBlockingReasons) == 0 {
		t.Fatalf("controlledApplyBlockingReasons must be populated")
	}
	if len(plan.ForbiddenControlledApplyPlanActions) == 0 {
		t.Fatalf("forbiddenControlledApplyPlanActions must be populated")
	}
	requiredForbidden := []string{
		"enable_controlled_apply",
		"execute_controlled_apply",
		"enable_executor",
		"execute_windows_executor",
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
	for _, item := range requiredForbidden {
		if !containsString(plan.ForbiddenControlledApplyPlanActions, item) {
			t.Fatalf("missing forbidden action: %s", item)
		}
	}
}

func TestLoadControlledApplyPlanBoundaryFixtureFromTemporaryFile(t *testing.T) {
	plan := NewWindowsFluidVirtControlledApplyPlanBoundaryMinimal()
	dir := t.TempDir()
	path := filepath.Join(dir, "controlled_apply_plan_boundary.json")
	raw, err := json.Marshal(plan)
	if err != nil {
		t.Fatalf("marshal plan: %v", err)
	}
	if err := os.WriteFile(path, raw, 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	loaded, err := LoadControlledApplyPlanBoundaryFixtureFromTemporaryFile(path)
	if err != nil {
		t.Fatalf("load fixture: %v", err)
	}
	if loaded.ControlledApplyEnabled || loaded.ControlledApplyReady {
		t.Fatalf("loaded plan must remain controlled-apply disabled and not-ready")
	}
}
