package windowsfluidvirt

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestComplianceReplayModelBuildsAndValidates(t *testing.T) {
	replay := NewWindowsFluidVirtComplianceReplayMinimal()
	if !replay.ComplianceReplayAvailable || !replay.ComplianceReplayExecuted || !replay.ComplianceReplayDeterministic {
		t.Fatalf("compliance replay must be available, executed, deterministic")
	}
	if !replay.AuditHashChainAvailable || !replay.AuditHashChainVerified {
		t.Fatalf("audit hash chain must be available and verified")
	}
	if replay.RuntimeMutationEnabled || replay.ActuatorRuntimeEnabled || replay.CgroupWriteEnabled {
		t.Fatalf("runtime mutation surfaces must be disabled")
	}
	if replay.QMPCommandExecutionAllowed || replay.QGACommandExecutionAllowed {
		t.Fatalf("qmp/qga execution must stay disabled")
	}
	if replay.ReplayValidationResult.ControlledApplyReady || replay.ReplayValidationResult.RuntimeMVPReady {
		t.Fatalf("controlled apply/runtime mvp readiness must remain false")
	}
	if err := ValidateWindowsFluidVirtComplianceReplay(replay); err != nil {
		t.Fatalf("validate replay: %v", err)
	}
}

func TestComplianceReplayRequiredChecksAndForbiddenActionsPresent(t *testing.T) {
	replay := NewWindowsFluidVirtComplianceReplayMinimal()
	requiredChecks := []string{
		"product_model_claim_boundary_valid",
		"actuator_contract_boundary_valid",
		"readonly_replay_boundary_valid",
		"fake_runtime_boundary_valid",
		"cgroup_write_disabled",
		"qmp_execution_disabled",
		"qga_execution_disabled",
		"autonomous_apply_disabled",
		"production_mutation_disabled",
		"windows_ga_claim_disabled",
		"windows_production_ready_claim_disabled",
		"windows_execution_ready_default_disabled",
		"vcpu_hotplug_claim_disabled",
		"logical_cpu_scaling_claim_disabled",
		"pool_scaling_claim_disabled",
		"raw_runtime_controls_disabled",
		"blocker_taxonomy_present",
		"audit_hash_chain_verified",
		"no_secret_material_present",
	}
	for _, id := range requiredChecks {
		if !containsComplianceCheck(replay.ComplianceChecks, id) {
			t.Fatalf("missing compliance check: %s", id)
		}
	}

	requiredForbidden := []string{
		"execute_node_actuator_runtime",
		"write_cgroup_cpu_max",
		"mutate_qmp_balloon",
		"execute_qmp_command",
		"execute_qga_command",
		"touch_real_cgroup",
		"touch_real_qmp",
		"touch_real_qga",
		"enable_runtime_actuator",
		"enable_controlled_apply",
		"enable_autonomous_apply",
		"enable_production_auto",
		"mark_windows_execution_ready",
		"claim_windows_ga",
		"claim_windows_production_ready",
		"expose_raw_runtime_controls",
		"touch_dashboard",
		"touch_inventory",
	}
	for _, action := range requiredForbidden {
		if !containsString(replay.ForbiddenComplianceReplayActions, action) {
			t.Fatalf("missing forbidden action: %s", action)
		}
	}
}

func TestLoadComplianceReplayFixtureFromTemporaryFile(t *testing.T) {
	replay := NewWindowsFluidVirtComplianceReplayMinimal()
	dir := t.TempDir()
	path := filepath.Join(dir, "compliance_replay.json")
	raw, err := json.Marshal(replay)
	if err != nil {
		t.Fatalf("marshal replay: %v", err)
	}
	if err := os.WriteFile(path, raw, 0o600); err != nil {
		t.Fatalf("write temp replay fixture: %v", err)
	}
	loaded, err := LoadComplianceReplayFixtureFromTemporaryFile(path)
	if err != nil {
		t.Fatalf("load replay fixture from temp: %v", err)
	}
	if loaded.ReplayValidationResult.RealRuntimeTouched || loaded.ReplayValidationResult.RealCgroupTouched {
		t.Fatalf("loaded replay must remain runtime-untouched")
	}
}

func containsComplianceCheck(checks []WindowsFluidVirtComplianceCheck, id string) bool {
	for _, check := range checks {
		if check.CheckID == id {
			return true
		}
	}
	return false
}

func containsString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
