package windowsfluidvirt

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNodeActuatorMVPDryRunAccepted(t *testing.T) {
	dir := t.TempDir()
	cgroupDir := filepath.Join(dir, "cgroup")
	if err := os.MkdirAll(cgroupDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	target := filepath.Join(cgroupDir, "cpu.max")
	if err := os.WriteFile(target, []byte("300000 100000\n"), 0o644); err != nil {
		t.Fatalf("seed cpu.max: %v", err)
	}
	request, allowlist := validActuatorInput(cgroupDir)
	evidence, err := EvaluateNodeFluidActuatorMVP(NodeFluidActuatorMVPInput{
		Mode:           NodeActuatorModeDryRun,
		Request:        request,
		Allowlist:      allowlist,
		EvaluationTime: time.Date(2026, 5, 7, 22, 5, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if !evidence.Allowed {
		t.Fatalf("expected dry run to be allowed, blockers=%v", evidence.Blockers)
	}
	if evidence.AfterCPUMax != "300000 100000" {
		t.Fatalf("dry-run must not mutate cpu.max, got %s", evidence.AfterCPUMax)
	}
}

func TestNodeActuatorMVPApplyWritesCPUMax(t *testing.T) {
	dir := t.TempDir()
	cgroupDir := filepath.Join(dir, "cgroup")
	if err := os.MkdirAll(cgroupDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	target := filepath.Join(cgroupDir, "cpu.max")
	if err := os.WriteFile(target, []byte("300000 100000\n"), 0o644); err != nil {
		t.Fatalf("seed cpu.max: %v", err)
	}
	request, allowlist := validActuatorInput(cgroupDir)
	evidence, err := EvaluateNodeFluidActuatorMVP(NodeFluidActuatorMVPInput{
		Mode:           NodeActuatorModeApply,
		Request:        request,
		Allowlist:      allowlist,
		EvaluationTime: time.Date(2026, 5, 7, 22, 5, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if !evidence.Allowed {
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

func TestNodeActuatorMVPBlocksStaleRequest(t *testing.T) {
	dir := t.TempDir()
	cgroupDir := filepath.Join(dir, "cgroup")
	if err := os.MkdirAll(cgroupDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	target := filepath.Join(cgroupDir, "cpu.max")
	if err := os.WriteFile(target, []byte("300000 100000\n"), 0o644); err != nil {
		t.Fatalf("seed cpu.max: %v", err)
	}
	request, allowlist := validActuatorInput(cgroupDir)
	request.CreatedAt = "2026-05-07T21:00:00Z"
	request.TTLSeconds = 60
	evidence, err := EvaluateNodeFluidActuatorMVP(NodeFluidActuatorMVPInput{
		Mode:           NodeActuatorModeDryRun,
		Request:        request,
		Allowlist:      allowlist,
		EvaluationTime: time.Date(2026, 5, 7, 22, 5, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if evidence.Allowed {
		t.Fatal("stale request must be blocked")
	}
	assertHas(t, evidence.Blockers, "stale_request")
}

func validActuatorInput(cgroupDir string) (NodeFluidActuatorMVPRequest, NodeFluidActuatorMVPAllowlist) {
	request := NodeFluidActuatorMVPRequest{
		RequestID:       "req-1",
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
		PreviousCPUMax:  "300000 100000",
		RequestedCPUMax: "600000 100000",
		RollbackCPUMax:  "300000 100000",
		TTLSeconds:      600,
		Reason:          "test",
		PolicyVersion:   "windows-fluid-node-actuator-v1",
		CreatedAt:       "2026-05-07T22:00:00Z",
	}
	allowlist := NodeFluidActuatorMVPAllowlist{
		NodeName:               request.NodeName,
		Namespace:              request.Namespace,
		VMName:                 request.VMName,
		PodUID:                 request.PodUID,
		QemuPID:                request.QemuPID,
		QemuStartTime:          request.QemuStartTime,
		CgroupPath:             request.CgroupPath,
		AllowedControllers:     []string{"cpu.max"},
		MinCPUMax:              "100000 100000",
		MaxCPUMax:              "800000 100000",
		AllowParentCgroupWrite: false,
		AllowArbitraryWrite:    false,
	}
	return request, allowlist
}
