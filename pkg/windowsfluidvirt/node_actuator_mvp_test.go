package windowsfluidvirt

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestNodeActuatorMVPDryRunAccepted(t *testing.T) {
	request, allowlist, _ := buildValidNodeActuatorFixture(t)
	request.Action = NodeActuatorActionDryRun
	evidence, err := EvaluateNodeFluidActuatorMVP(NodeFluidActuatorMVPInput{
		Action:         NodeActuatorActionDryRun,
		Request:        request,
		Allowlist:      allowlist,
		EvaluationTime: time.Date(2026, 5, 7, 22, 5, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if evidence.Decision != NodeActuatorDecisionAccepted {
		t.Fatalf("expected dry run to be allowed, blockers=%v", evidence.Blockers)
	}
	if evidence.MutationPerformed {
		t.Fatal("dry-run must never mutate")
	}
}

func TestNodeActuatorMVPApplyWritesCPUMax(t *testing.T) {
	request, allowlist, target := buildValidNodeActuatorFixture(t)
	evidence, err := EvaluateNodeFluidActuatorMVP(NodeFluidActuatorMVPInput{
		Action:         NodeActuatorActionApply,
		Request:        request,
		Allowlist:      allowlist,
		EvaluationTime: time.Date(2026, 5, 7, 22, 5, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if evidence.Decision != NodeActuatorDecisionApplied {
		t.Fatalf("expected apply allowed blockers=%v", evidence.Blockers)
	}
	content, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("read target: %v", err)
	}
	if string(content) != "600000 100000\n" {
		t.Fatalf("unexpected cpu.max content %q", string(content))
	}
}

func TestNodeActuatorMVPRollbackAccepted(t *testing.T) {
	request, allowlist, target := buildValidNodeActuatorFixture(t)
	if err := os.WriteFile(target, []byte("600000 100000\n"), 0o644); err != nil {
		t.Fatalf("seed ceiling: %v", err)
	}
	request.Action = NodeActuatorActionRollback
	request.PreviousCPUMax = "600000 100000"
	request.RollbackCPUMax = "300000 100000"
	evidence, err := EvaluateNodeFluidActuatorMVP(NodeFluidActuatorMVPInput{
		Action:         NodeActuatorActionRollback,
		Request:        request,
		Allowlist:      allowlist,
		EvaluationTime: time.Date(2026, 5, 7, 22, 5, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("rollback evaluate: %v", err)
	}
	if evidence.Decision != NodeActuatorDecisionRolledBack {
		t.Fatalf("expected rolled_back, got %s blockers=%v", evidence.Decision, evidence.Blockers)
	}
	content, _ := os.ReadFile(target)
	if strings.TrimSpace(string(content)) != "300000 100000" {
		t.Fatalf("expected rollback target restored, got %q", string(content))
	}
}

func TestNodeActuatorMVPReturnToFloorAccepted(t *testing.T) {
	request, allowlist, target := buildValidNodeActuatorFixture(t)
	request.Action = NodeActuatorActionReturnToFloor
	request.PreviousCPUMax = "300000 100000"
	request.RollbackCPUMax = "600000 100000"
	if err := os.WriteFile(target, []byte("300000 100000\n"), 0o644); err != nil {
		t.Fatalf("seed floor: %v", err)
	}
	evidence, err := EvaluateNodeFluidActuatorMVP(NodeFluidActuatorMVPInput{
		Action:         NodeActuatorActionReturnToFloor,
		Request:        request,
		Allowlist:      allowlist,
		EvaluationTime: time.Date(2026, 5, 7, 22, 5, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("return-to-floor evaluate: %v", err)
	}
	if evidence.Decision != NodeActuatorDecisionReturnedToFloor {
		t.Fatalf("expected returned_to_floor, got %s blockers=%v", evidence.Decision, evidence.Blockers)
	}
	content, _ := os.ReadFile(target)
	if strings.TrimSpace(string(content)) != "300000 100000" {
		t.Fatalf("expected floor target restored, got %q", string(content))
	}
}

func TestNodeActuatorMVPNegativeMatrix(t *testing.T) {
	t.Run("stale request rejected", func(t *testing.T) {
		request, allowlist, _ := buildValidNodeActuatorFixture(t)
		request.CreatedAt = "2026-05-07T21:00:00Z"
		request.ExpiresAt = "2026-05-07T21:10:00Z"
		evidence, _ := EvaluateNodeFluidActuatorMVP(NodeFluidActuatorMVPInput{
			Action:         NodeActuatorActionApply,
			Request:        request,
			Allowlist:      allowlist,
			EvaluationTime: time.Date(2026, 5, 7, 22, 5, 0, 0, time.UTC),
		})
		assertHas(t, evidence.Blockers, "stale_request")
	})
	t.Run("kill switch blocked", func(t *testing.T) {
		request, allowlist, _ := buildValidNodeActuatorFixture(t)
		kill := filepath.Join(t.TempDir(), "kill.txt")
		if err := os.WriteFile(kill, []byte("blocked\n"), 0o644); err != nil {
			t.Fatalf("write kill: %v", err)
		}
		evidence, _ := EvaluateNodeFluidActuatorMVP(NodeFluidActuatorMVPInput{
			Action:         NodeActuatorActionApply,
			Request:        request,
			Allowlist:      allowlist,
			KillSwitchPath: kill,
			EvaluationTime: time.Date(2026, 5, 7, 22, 5, 0, 0, time.UTC),
		})
		assertHas(t, evidence.Blockers, "kill_switch_blocked")
	})
	t.Run("pod uid mismatch rejected", func(t *testing.T) {
		request, allowlist, _ := buildValidNodeActuatorFixture(t)
		allowlist.PodUID = "other-pod"
		evidence, _ := EvaluateNodeFluidActuatorMVP(NodeFluidActuatorMVPInput{
			Action:    NodeActuatorActionApply,
			Request:   request,
			Allowlist: allowlist,
		})
		assertHas(t, evidence.Blockers, "allowlist_identity_mismatch")
	})
	t.Run("qemu pid mismatch rejected", func(t *testing.T) {
		request, allowlist, _ := buildValidNodeActuatorFixture(t)
		allowlist.QemuPID = "9001"
		evidence, _ := EvaluateNodeFluidActuatorMVP(NodeFluidActuatorMVPInput{
			Action:    NodeActuatorActionApply,
			Request:   request,
			Allowlist: allowlist,
		})
		assertHas(t, evidence.Blockers, "allowlist_identity_mismatch")
	})
	t.Run("qemu start mismatch rejected", func(t *testing.T) {
		request, allowlist, _ := buildValidNodeActuatorFixture(t)
		allowlist.QemuStartTime = "other"
		evidence, _ := EvaluateNodeFluidActuatorMVP(NodeFluidActuatorMVPInput{
			Action:    NodeActuatorActionApply,
			Request:   request,
			Allowlist: allowlist,
		})
		assertHas(t, evidence.Blockers, "allowlist_identity_mismatch")
	})
	t.Run("cgroup path mismatch rejected", func(t *testing.T) {
		request, allowlist, _ := buildValidNodeActuatorFixture(t)
		allowlist.AllowedCgroupPath = filepath.Join(t.TempDir(), "other")
		evidence, _ := EvaluateNodeFluidActuatorMVP(NodeFluidActuatorMVPInput{
			Action:    NodeActuatorActionApply,
			Request:   request,
			Allowlist: allowlist,
		})
		assertHas(t, evidence.Blockers, "allowlist_cgroup_mismatch")
	})
	t.Run("symlink traversal rejected", func(t *testing.T) {
		request, allowlist, target := buildValidNodeActuatorFixture(t)
		if err := os.Remove(target); err != nil {
			t.Fatalf("remove target: %v", err)
		}
		evil := filepath.Join(t.TempDir(), "evil")
		if err := os.WriteFile(evil, []byte("300000 100000\n"), 0o644); err != nil {
			t.Fatalf("write evil: %v", err)
		}
		if err := os.Symlink(evil, target); err != nil {
			t.Fatalf("symlink: %v", err)
		}
		evidence, _ := EvaluateNodeFluidActuatorMVP(NodeFluidActuatorMVPInput{
			Action:    NodeActuatorActionApply,
			Request:   request,
			Allowlist: allowlist,
		})
		assertHas(t, evidence.Blockers, "symlink_or_path_escape_detected")
	})
	t.Run("parent cgroup write rejected", func(t *testing.T) {
		request, allowlist, _ := buildValidNodeActuatorFixture(t)
		allowlist.AllowParentCgroupWrite = true
		evidence, _ := EvaluateNodeFluidActuatorMVP(NodeFluidActuatorMVPInput{
			Action:    NodeActuatorActionApply,
			Request:   request,
			Allowlist: allowlist,
		})
		assertHas(t, evidence.Blockers, "allow_parent_cgroup_write_must_be_false")
	})
	t.Run("controller other than cpu.max rejected", func(t *testing.T) {
		request, allowlist, _ := buildValidNodeActuatorFixture(t)
		request.Controller = "cpu.weight"
		evidence, _ := EvaluateNodeFluidActuatorMVP(NodeFluidActuatorMVPInput{
			Action:    NodeActuatorActionApply,
			Request:   request,
			Allowlist: allowlist,
		})
		assertHas(t, evidence.Blockers, "controller_not_allowed")
	})
	t.Run("requested cpu.max outside bounds rejected", func(t *testing.T) {
		request, allowlist, _ := buildValidNodeActuatorFixture(t)
		request.RequestedCPUMax = "900000 100000"
		evidence, _ := EvaluateNodeFluidActuatorMVP(NodeFluidActuatorMVPInput{
			Action:    NodeActuatorActionApply,
			Request:   request,
			Allowlist: allowlist,
		})
		assertHas(t, evidence.Blockers, "requested_cpu_max_outside_request_bounds")
	})
	t.Run("missing rollback target rejected", func(t *testing.T) {
		request, allowlist, _ := buildValidNodeActuatorFixture(t)
		request.RollbackCPUMax = ""
		evidence, _ := EvaluateNodeFluidActuatorMVP(NodeFluidActuatorMVPInput{
			Action:    NodeActuatorActionApply,
			Request:   request,
			Allowlist: allowlist,
		})
		assertHas(t, evidence.Blockers, "rollback_target_missing")
	})
}

func TestNodeActuatorCLIOnFakeTarget(t *testing.T) {
	request, allowlist, target := buildValidNodeActuatorFixture(t)
	request.Action = NodeActuatorActionDryRun
	reqPath := writeJSONFixture(t, "request.json", request)
	allowPath := writeJSONFixture(t, "allowlist.json", allowlist)
	outPath := filepath.Join(t.TempDir(), "out.json")

	cmd := exec.Command("go", "run", "./cmd/karl-node-fluid-actuator",
		"-mode", "dry-run",
		"-request", reqPath,
		"-allowlist", allowPath,
		"-evaluation-time", "2026-05-07T22:05:00Z",
		"-evidence-out", outPath,
	)
	cmd.Dir = admissionRepoRoot(t)
	raw, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("dry-run cli failed: %v\n%s", err, string(raw))
	}

	request.Action = NodeActuatorActionApply
	reqPath = writeJSONFixture(t, "request-apply.json", request)
	cmd = exec.Command("go", "run", "./cmd/karl-node-fluid-actuator",
		"-mode", "apply",
		"-request", reqPath,
		"-allowlist", allowPath,
		"-evaluation-time", "2026-05-07T22:05:00Z",
		"-evidence-out", outPath,
	)
	cmd.Dir = admissionRepoRoot(t)
	if raw, err = cmd.CombinedOutput(); err != nil {
		t.Fatalf("apply cli failed: %v\n%s", err, string(raw))
	}
	content, _ := os.ReadFile(target)
	if strings.TrimSpace(string(content)) != "600000 100000" {
		t.Fatalf("expected apply mutate target, got %q", string(content))
	}

	request.Action = NodeActuatorActionRollback
	request.PreviousCPUMax = "600000 100000"
	request.RollbackCPUMax = "300000 100000"
	reqPath = writeJSONFixture(t, "request-rollback.json", request)
	cmd = exec.Command("go", "run", "./cmd/karl-node-fluid-actuator",
		"-mode", "rollback",
		"-request", reqPath,
		"-allowlist", allowPath,
		"-evaluation-time", "2026-05-07T22:05:00Z",
		"-evidence-out", outPath,
	)
	cmd.Dir = admissionRepoRoot(t)
	if raw, err = cmd.CombinedOutput(); err != nil {
		t.Fatalf("rollback cli failed: %v\n%s", err, string(raw))
	}
	content, _ = os.ReadFile(target)
	if strings.TrimSpace(string(content)) != "300000 100000" {
		t.Fatalf("expected rollback mutate target, got %q", string(content))
	}

	request.Action = NodeActuatorActionReturnToFloor
	request.PreviousCPUMax = "300000 100000"
	request.RollbackCPUMax = "600000 100000"
	reqPath = writeJSONFixture(t, "request-floor.json", request)
	cmd = exec.Command("go", "run", "./cmd/karl-node-fluid-actuator",
		"-mode", "return-to-floor",
		"-request", reqPath,
		"-allowlist", allowPath,
		"-evaluation-time", "2026-05-07T22:05:00Z",
		"-evidence-out", outPath,
	)
	cmd.Dir = admissionRepoRoot(t)
	if raw, err = cmd.CombinedOutput(); err != nil {
		t.Fatalf("return-to-floor cli failed: %v\n%s", err, string(raw))
	}
	content, _ = os.ReadFile(target)
	if strings.TrimSpace(string(content)) != "300000 100000" {
		t.Fatalf("expected return-to-floor mutate target, got %q", string(content))
	}
}

func buildValidNodeActuatorFixture(t *testing.T) (KARLNodeFluidActuatorRequest, KARLNodeFluidActuatorAllowlist, string) {
	t.Helper()
	cgroupDir := filepath.Join(t.TempDir(), "cgroup")
	if err := os.MkdirAll(cgroupDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	target := filepath.Join(cgroupDir, "cpu.max")
	if err := os.WriteFile(target, []byte("300000 100000\n"), 0o644); err != nil {
		t.Fatalf("seed cpu.max: %v", err)
	}
	now := time.Now().UTC().Truncate(time.Second)
	createdAt := now.Format(time.RFC3339)
	expiresAt := now.Add(10 * time.Minute).Format(time.RFC3339)
	request := KARLNodeFluidActuatorRequest{
		RequestID:       "req-1",
		RequestVersion:  "windows-fluid-node-actuator-mvp-v1",
		Action:          NodeActuatorActionApply,
		Namespace:       "karl",
		VMName:          "master-win11",
		VMUID:           "vm-uid",
		VMIUID:          "vmi-uid",
		PodName:         "virt-launcher-master-win11-kmwgg",
		PodUID:          "pod-uid",
		NodeName:        "karl-lab-metal-01",
		QemuPID:         "96",
		QemuStartTime:   "Thu May 7 18:58:03 2026",
		CgroupPath:      cgroupDir,
		Controller:      "cpu.max",
		PreviousCPUMax:  "300000 100000",
		RequestedCPUMax: "600000 100000",
		RollbackCPUMax:  "300000 100000",
		MinCPUMax:       "300000 100000",
		MaxCPUMax:       "600000 100000",
		TTLSeconds:      600,
		CreatedAt:       createdAt,
		ExpiresAt:       expiresAt,
		Reason:          "test",
		Risk:            "low",
		EvidenceRefs:    []string{"fixture://test"},
		PolicyVersion:   "windows-fluid-node-actuator-v1",
	}
	allowlist := KARLNodeFluidActuatorAllowlist{
		AllowlistID:            "allow-1",
		NodeName:               request.NodeName,
		Namespace:              request.Namespace,
		VMName:                 request.VMName,
		VMUID:                  request.VMUID,
		PodUID:                 request.PodUID,
		QemuPID:                request.QemuPID,
		QemuStartTime:          request.QemuStartTime,
		AllowedCgroupPath:      request.CgroupPath,
		AllowedControllers:     []string{"cpu.max"},
		MinCPUMax:              "300000 100000",
		MaxCPUMax:              "600000 100000",
		AllowParentCgroupWrite: false,
		AllowArbitraryWrite:    false,
		AllowSymlinkTraversal:  false,
		AllowedActions: []KARLNodeFluidActuatorAction{
			NodeActuatorActionDryRun,
			NodeActuatorActionApply,
			NodeActuatorActionRollback,
			NodeActuatorActionReturnToFloor,
		},
		CreatedAt: createdAt,
		ExpiresAt: expiresAt,
	}
	return request, allowlist, target
}

func writeJSONFixture(t *testing.T, filename string, value any) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), filename)
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		t.Fatalf("marshal fixture: %v", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	return path
}
