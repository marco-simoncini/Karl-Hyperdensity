package windowsfluidvirt

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestDefaultGateBlocksApplyWithoutApproval(t *testing.T) {
	target := EvaluateWindowsHyperdensityTarget(baseReadyTarget())
	lease := baseCombinedLease()
	lease.PolicySnapshot = map[string]any{"killSwitchState": "allow", "workloadEvidenceRef": "workload://cpu-ram"}
	gate := DefaultWindowsFluidControlledApplyGate()
	plan := BuildWindowsFluidControlledApplyPlan(target, lease, gate, nil, fixedControlledTime())
	if plan.ApplyAllowed {
		t.Fatal("default gate must block apply without approval")
	}
	assertHas(t, plan.Blockers, ControlledApplyBlockerManualApprovalRequired)
}

func TestApprovedFixtureReachesApplyReady(t *testing.T) {
	target := EvaluateWindowsHyperdensityTarget(baseReadyTarget())
	lease := baseCombinedLease()
	lease.PolicySnapshot = map[string]any{"killSwitchState": "allow", "workloadEvidenceRef": "workload://proof"}
	gate := DefaultWindowsFluidControlledApplyGate()
	gate.CPUApplyEnabled = true
	gate.RAMApplyEnabled = true
	gate.NodeActuatorApplyEnabled = true
	gate.QMPBalloonApplyEnabled = true
	approval := &WindowsFluidManualApproval{
		ApprovalID: "approval-master-win11",
		State:      string(ApprovalApproved),
		CreatedAt:  "2026-05-08T04:00:00Z",
		ExpiresAt:  "2026-05-09T04:00:00Z",
	}
	plan := BuildWindowsFluidControlledApplyPlan(target, lease, gate, approval, fixedControlledTime())
	if !plan.ApplyAllowed {
		t.Fatalf("expected apply ready, blockers=%v", plan.Blockers)
	}
	if plan.PlanPhase != ControlledPlanApplyReady {
		t.Fatalf("expected apply_ready, got %s", plan.PlanPhase)
	}
}

func TestAutonomousApplyRejected(t *testing.T) {
	target := EvaluateWindowsHyperdensityTarget(baseReadyTarget())
	lease := baseCombinedLease()
	lease.PolicySnapshot = map[string]any{"killSwitchState": "allow", "workloadEvidenceRef": "workload://proof"}
	gate := DefaultWindowsFluidControlledApplyGate()
	gate.AutonomousApplyEnabled = true
	approval := &WindowsFluidManualApproval{ApprovalID: "a1", State: string(ApprovalApproved), CreatedAt: "2026-05-08T04:00:00Z"}
	plan := BuildWindowsFluidControlledApplyPlan(target, lease, gate, approval, fixedControlledTime())
	if plan.ApplyAllowed {
		t.Fatal("autonomous apply must be rejected")
	}
	assertHas(t, plan.Blockers, ControlledApplyBlockerAutonomousApplyDenied)
}

func TestDryRunRequired(t *testing.T) {
	slate := WindowsFluidActionSlate{Actions: []WindowsFluidAction{}, Blockers: nil}
	blockers := EvaluateWindowsFluidDryRunGate(DefaultWindowsFluidControlledApplyGate(), baseCombinedLease(), slate)
	assertHas(t, blockers, ControlledApplyBlockerDryRunRequired)
}

func TestKillSwitchBlocks(t *testing.T) {
	target := EvaluateWindowsHyperdensityTarget(baseReadyTarget())
	lease := baseCombinedLease()
	lease.PolicySnapshot = map[string]any{"killSwitchState": "blocked", "workloadEvidenceRef": "workload://proof"}
	gate := DefaultWindowsFluidControlledApplyGate()
	gate.CPUApplyEnabled = true
	gate.RAMApplyEnabled = true
	gate.NodeActuatorApplyEnabled = true
	gate.QMPBalloonApplyEnabled = true
	approval := &WindowsFluidManualApproval{ApprovalID: "a1", State: string(ApprovalApproved), CreatedAt: "2026-05-08T04:00:00Z"}
	plan := BuildWindowsFluidControlledApplyPlan(target, lease, gate, approval, fixedControlledTime())
	assertHas(t, plan.Blockers, ControlledApplyBlockerKillSwitchBlocked)
}

func TestMissingAuditBlocks(t *testing.T) {
	target := EvaluateWindowsHyperdensityTarget(baseReadyTarget())
	lease := baseCombinedLease()
	lease.AuditBundleRef = ""
	lease.PolicySnapshot = map[string]any{"killSwitchState": "allow", "workloadEvidenceRef": "workload://proof"}
	gate := enabledApplyGate()
	approval := approvedManualApproval()
	plan := BuildWindowsFluidControlledApplyPlan(target, lease, gate, approval, fixedControlledTime())
	assertHas(t, plan.Blockers, ControlledApplyBlockerAuditBundleMissing)
}

func TestMissingRollbackBlocks(t *testing.T) {
	target := EvaluateWindowsHyperdensityTarget(baseReadyTarget())
	lease := baseCombinedLease()
	lease.RollbackTarget = WindowsFluidLeaseRequest{}
	lease.PolicySnapshot = map[string]any{"killSwitchState": "allow", "workloadEvidenceRef": "workload://proof"}
	plan := BuildWindowsFluidControlledApplyPlan(target, lease, enabledApplyGate(), approvedManualApproval(), fixedControlledTime())
	assertHas(t, plan.Blockers, BlockerRollbackNotReady)
}

func TestMissingReturnToFloorBlocks(t *testing.T) {
	target := EvaluateWindowsHyperdensityTarget(baseReadyTarget())
	lease := baseCombinedLease()
	lease.ReturnToFloorTarget = WindowsFluidLeaseRequest{}
	lease.PolicySnapshot = map[string]any{"killSwitchState": "allow", "workloadEvidenceRef": "workload://proof"}
	plan := BuildWindowsFluidControlledApplyPlan(target, lease, enabledApplyGate(), approvedManualApproval(), fixedControlledTime())
	assertHas(t, plan.Blockers, BlockerReturnToFloorNotReady)
}

func TestMissingWorkloadVerificationBlocks(t *testing.T) {
	target := EvaluateWindowsHyperdensityTarget(baseReadyTarget())
	lease := baseCombinedLease()
	lease.PolicySnapshot = map[string]any{"killSwitchState": "allow"}
	plan := BuildWindowsFluidControlledApplyPlan(target, lease, enabledApplyGate(), approvedManualApproval(), fixedControlledTime())
	assertHas(t, plan.Blockers, ControlledApplyBlockerWorkloadVerifyMissing)
}

func TestMasterWin11PlanBuildsCPUAndRAMAndAuditActions(t *testing.T) {
	target := EvaluateWindowsHyperdensityTarget(baseReadyTarget())
	lease := baseCombinedLease()
	lease.PolicySnapshot = map[string]any{"killSwitchState": "allow", "workloadEvidenceRef": "workload://proof"}
	plan := BuildWindowsFluidControlledApplyPlan(target, lease, enabledApplyGate(), approvedManualApproval(), fixedControlledTime())
	if plan.ActuatorRequest["ref"] == "" {
		t.Fatal("expected actuator request in plan")
	}
	if plan.QMPBalloonRequest["ref"] == "" {
		t.Fatal("expected qmp request in plan")
	}
	if plan.GuestVerifyPlan["required"] != true {
		t.Fatal("guest verification must be required")
	}
	if plan.AuditBundlePlan["required"] != true {
		t.Fatal("audit plan must be required")
	}
}

func TestPoolChildControlledPlanAllowedAsIndividualTarget(t *testing.T) {
	target := baseReadyTarget()
	target.TargetKind = TargetKindPoolChildWindowsVM
	target = EvaluateWindowsHyperdensityTarget(target)
	lease := baseCombinedLease()
	lease.PolicySnapshot = map[string]any{"killSwitchState": "allow", "workloadEvidenceRef": "workload://proof"}
	plan := BuildWindowsFluidControlledApplyPlan(target, lease, enabledApplyGate(), approvedManualApproval(), fixedControlledTime())
	if plan.PlanPhase != ControlledPlanApplyReady {
		t.Fatalf("expected apply_ready for pool child individual target, got %s blockers=%v", plan.PlanPhase, plan.Blockers)
	}
}

func TestPoolScalingBlocked(t *testing.T) {
	target := baseReadyTarget()
	target.PoolScalingRequested = true
	target = EvaluateWindowsHyperdensityTarget(target)
	lease := baseCombinedLease()
	lease.PolicySnapshot = map[string]any{"killSwitchState": "allow", "workloadEvidenceRef": "workload://proof"}
	plan := BuildWindowsFluidControlledApplyPlan(target, lease, enabledApplyGate(), approvedManualApproval(), fixedControlledTime())
	assertHas(t, plan.Blockers, BlockerPoolScalingAsMechanism)
}

func TestVCPUHotplugRejected(t *testing.T) {
	target := EvaluateWindowsHyperdensityTarget(baseReadyTarget())
	lease := baseCombinedLease()
	lease.RequestsVCPUHotplug = true
	lease.PolicySnapshot = map[string]any{"killSwitchState": "allow", "workloadEvidenceRef": "workload://proof"}
	plan := BuildWindowsFluidControlledApplyPlan(target, lease, enabledApplyGate(), approvedManualApproval(), fixedControlledTime())
	assertHas(t, plan.Blockers, BlockerLeaseRequestsVCPUHotplug)
}

func TestLogicalCPUScalingRejected(t *testing.T) {
	target := baseReadyTarget()
	target.LogicalCPUScalingClaimed = true
	target = EvaluateWindowsHyperdensityTarget(target)
	lease := baseCombinedLease()
	lease.PolicySnapshot = map[string]any{"killSwitchState": "allow", "workloadEvidenceRef": "workload://proof"}
	plan := BuildWindowsFluidControlledApplyPlan(target, lease, enabledApplyGate(), approvedManualApproval(), fixedControlledTime())
	assertHas(t, plan.Blockers, BlockerLeaseRequestsVCPUHotplug)
}

func TestVMSpecPatchRejected(t *testing.T) {
	target := EvaluateWindowsHyperdensityTarget(baseReadyTarget())
	lease := baseCombinedLease()
	lease.RequestsVMSpecPatch = true
	lease.PolicySnapshot = map[string]any{"killSwitchState": "allow", "workloadEvidenceRef": "workload://proof"}
	plan := BuildWindowsFluidControlledApplyPlan(target, lease, enabledApplyGate(), approvedManualApproval(), fixedControlledTime())
	assertHas(t, plan.Blockers, BlockerLeaseRequestsVMSpecPatch)
}

func TestControlledApplyFixtureMatrix(t *testing.T) {
	names := []string{
		"master-win11-controlled-plan.awaiting-approval.json",
		"master-win11-controlled-plan.apply-ready.json",
		"master-win11-controlled-plan.dryrun-blocked.json",
		"master-win11-controlled-plan.kill-switch-blocked.json",
		"master-win11-controlled-plan.missing-audit.blocked.json",
		"master-win11-controlled-plan.missing-return.blocked.json",
		"master-win11-controlled-plan.autonomous-apply.rejected.json",
		"pool-child-controlled-plan.awaiting-approval.json",
		"pool-scaling-controlled-plan.blocked.json",
		"vcpu-hotplug-controlled-plan.rejected.json",
		"logical-cpu-scaling-controlled-plan.rejected.json",
		"vm-spec-patch-controlled-plan.rejected.json",
	}
	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			fixture, err := LoadWindowsFluidControlledApplyFixture(controlledFixturePath(t, name))
			if err != nil {
				t.Fatalf("load fixture: %v", err)
			}
			plan := BuildWindowsFluidControlledApplyPlan(
				fixture.Target,
				fixture.Lease,
				fixture.Gate,
				fixture.Approval,
				fixedControlledTime(),
			)
			if plan.PlanPhase != fixture.ExpectedPlanPhase {
				t.Fatalf("phase mismatch expected=%s got=%s blockers=%v", fixture.ExpectedPlanPhase, plan.PlanPhase, plan.Blockers)
			}
			if plan.ApplyAllowed != fixture.ExpectedApplyAllowed {
				t.Fatalf("applyAllowed mismatch expected=%v got=%v blockers=%v", fixture.ExpectedApplyAllowed, plan.ApplyAllowed, plan.Blockers)
			}
			for _, blocker := range fixture.ExpectedBlockers {
				assertHas(t, plan.Blockers, blocker)
			}
		})
	}
}

func TestControlledExecutorCLIDeterministicOutput(t *testing.T) {
	target := EvaluateWindowsHyperdensityTarget(baseReadyTarget())
	lease := baseCombinedLease()
	lease.PolicySnapshot = map[string]any{"killSwitchState": "allow", "workloadEvidenceRef": "workload://proof"}
	gate := enabledApplyGate()
	approval := approvedManualApproval()

	tmpDir := t.TempDir()
	targetPath := filepath.Join(tmpDir, "target.json")
	leasePath := filepath.Join(tmpDir, "lease.json")
	gatePath := filepath.Join(tmpDir, "gate.json")
	approvalPath := filepath.Join(tmpDir, "approval.json")
	writeJSON(t, targetPath, target)
	writeJSON(t, leasePath, lease)
	writeJSON(t, gatePath, gate)
	writeJSON(t, approvalPath, approval)

	cmd := exec.Command(
		"go", "run", "./cmd/karl-fluid-windows-executor",
		"-target", targetPath,
		"-lease", leasePath,
		"-gate", gatePath,
		"-approval", approvalPath,
		"-mode", "apply-plan-only",
		"-evaluation-time", "2026-05-08T05:00:00Z",
		"-pretty",
	)
	cmd.Dir = governanceRepoRoot(t)
	out1, err := cmd.Output()
	if err != nil {
		t.Fatalf("first cli run failed: %v", err)
	}
	cmd = exec.Command(
		"go", "run", "./cmd/karl-fluid-windows-executor",
		"-target", targetPath,
		"-lease", leasePath,
		"-gate", gatePath,
		"-approval", approvalPath,
		"-mode", "apply-plan-only",
		"-evaluation-time", "2026-05-08T05:00:00Z",
		"-pretty",
	)
	cmd.Dir = governanceRepoRoot(t)
	out2, err := cmd.Output()
	if err != nil {
		t.Fatalf("second cli run failed: %v", err)
	}
	if string(out1) != string(out2) {
		t.Fatal("cli output must be deterministic with fixed evaluation-time")
	}

	var parsed map[string]any
	if err := json.Unmarshal(out1, &parsed); err != nil {
		t.Fatalf("unmarshal cli output: %v", err)
	}
	if parsed["runtimeMutationExecuted"] != false {
		t.Fatal("runtime mutation must stay disabled")
	}
}

func TestControlledExecutorCLIPhasesFromFixtures(t *testing.T) {
	cases := []struct {
		name          string
		mode          string
		fixture       string
		expectedPhase string
	}{
		{"awaiting approval", "plan", "master-win11-controlled-plan.awaiting-approval.json", "awaiting_approval"},
		{"apply ready", "apply-plan-only", "master-win11-controlled-plan.apply-ready.json", "apply_ready"},
		{"autonomous rejected", "plan", "master-win11-controlled-plan.autonomous-apply.rejected.json", "apply_blocked"},
		{"pool scaling blocked", "plan", "pool-scaling-controlled-plan.blocked.json", "apply_blocked"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fixture, err := LoadWindowsFluidControlledApplyFixture(controlledFixturePath(t, tc.fixture))
			if err != nil {
				t.Fatalf("load fixture: %v", err)
			}
			tmpDir := t.TempDir()
			targetPath := filepath.Join(tmpDir, "target.json")
			leasePath := filepath.Join(tmpDir, "lease.json")
			gatePath := filepath.Join(tmpDir, "gate.json")
			writeJSON(t, targetPath, fixture.Target)
			writeJSON(t, leasePath, fixture.Lease)
			writeJSON(t, gatePath, fixture.Gate)
			args := []string{
				"run", "./cmd/karl-fluid-windows-executor",
				"-target", targetPath,
				"-lease", leasePath,
				"-gate", gatePath,
				"-mode", tc.mode,
				"-evaluation-time", "2026-05-08T05:00:00Z",
			}
			if fixture.Approval != nil {
				approvalPath := filepath.Join(tmpDir, "approval.json")
				writeJSON(t, approvalPath, fixture.Approval)
				args = append(args, "-approval", approvalPath)
			}
			cmd := exec.Command("go", args...)
			cmd.Dir = governanceRepoRoot(t)
			out, err := cmd.Output()
			if err != nil {
				t.Fatalf("cli run failed: %v", err)
			}
			var parsed struct {
				Plan struct {
					PlanPhase string `json:"planPhase"`
				} `json:"plan"`
			}
			if err := json.Unmarshal(out, &parsed); err != nil {
				t.Fatalf("unmarshal cli output: %v", err)
			}
			if parsed.Plan.PlanPhase != tc.expectedPhase {
				t.Fatalf("phase mismatch expected=%s got=%s", tc.expectedPhase, parsed.Plan.PlanPhase)
			}
		})
	}
}

func enabledApplyGate() WindowsFluidControlledApplyGate {
	gate := DefaultWindowsFluidControlledApplyGate()
	gate.CPUApplyEnabled = true
	gate.RAMApplyEnabled = true
	gate.NodeActuatorApplyEnabled = true
	gate.QMPBalloonApplyEnabled = true
	return gate
}

func approvedManualApproval() *WindowsFluidManualApproval {
	return &WindowsFluidManualApproval{
		ApprovalID: "approval-master-win11",
		State:      string(ApprovalApproved),
		CreatedAt:  "2026-05-08T04:00:00Z",
		ExpiresAt:  "2026-05-09T04:00:00Z",
	}
}

func writeJSON(t *testing.T, path string, value any) {
	t.Helper()
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		t.Fatalf("marshal json: %v", err)
	}
	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		t.Fatalf("write json: %v", err)
	}
}

func fixedControlledTime() time.Time {
	return time.Date(2026, 5, 8, 5, 0, 0, 0, time.UTC)
}

func controlledFixturePath(t *testing.T, name string) string {
	t.Helper()
	return filepath.Join(governanceRepoRoot(t), "examples", "windows-fluid-controlled-apply-fixtures", name)
}
