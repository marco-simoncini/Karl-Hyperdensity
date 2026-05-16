package resourcelease

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/cgroup"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/crdv1alpha1"
)

func guardedApplyOpts(lease *crdv1alpha1.ResourceLease, dir string, apply, confirm bool) GuardedApplySandboxOptions {
	cfg := sandboxDryRunCfg()
	cfg.Spec.AllowPathPrefixes = []string{filepath.Join(dir, "karl.slice")}
	return GuardedApplySandboxOptions{
		DryRunAgainstPortOptions: DryRunAgainstPortOptions{
			Config:     cfg,
			Lease:      lease,
			Namespace:  "khr-runtime-sandbox",
			Labels:     map[string]string{"khr.karl.io/sandbox": "true"},
			Ports:      []crdv1alpha1.ResourcePort{sandboxPort()},
			SandboxDir: dir,
			BaselineID: "test-baseline",
		},
		ApplyResourceLease: apply,
		SandboxConfirm:     confirm,
	}
}

func TestGuardedApplyAgainstPortsBlockedByDefault(t *testing.T) {
	lease := loadSandboxLease(t, "resourcelease-dryrun-allowed.json")
	dir := t.TempDir()
	res, err := GuardedApplyAgainstResourcePorts(guardedApplyOpts(lease, dir, false, false))
	if err != nil {
		t.Fatal(err)
	}
	if !res.Blocked || res.Applied {
		t.Fatalf("res=%+v", res)
	}
}

func TestGuardedApplyRequiresConfirmation(t *testing.T) {
	lease := loadSandboxLease(t, "resourcelease-dryrun-allowed.json")
	dir := t.TempDir()
	res, err := GuardedApplyAgainstResourcePorts(guardedApplyOpts(lease, dir, true, false))
	if err != nil {
		t.Fatal(err)
	}
	if !res.Blocked {
		t.Fatalf("res=%+v", res)
	}
}

func TestGuardedApplyRequiresDryRunAllowed(t *testing.T) {
	lease := loadSandboxLease(t, "resourcelease-dryrun-blocked-missing-port.json")
	dir := t.TempDir()
	res, err := GuardedApplyAgainstResourcePorts(guardedApplyOpts(lease, dir, true, true))
	if err != nil {
		t.Fatal(err)
	}
	if !res.Blocked {
		t.Fatalf("res=%+v", res)
	}
}

func TestGuardedApplyBlockedNoRollbackPlanRef(t *testing.T) {
	lease := loadSandboxLease(t, "resourcelease-guarded-apply-blocked-no-rollback.json")
	dir := t.TempDir()
	res, err := GuardedApplyAgainstResourcePorts(guardedApplyOpts(lease, dir, true, true))
	if err != nil {
		t.Fatal(err)
	}
	if !res.Blocked {
		t.Fatalf("res=%+v", res)
	}
}

func TestGuardedApplyBlockedProductionNamespace(t *testing.T) {
	lease := loadSandboxLease(t, "resourcelease-dryrun-allowed.json")
	dir := t.TempDir()
	opts := guardedApplyOpts(lease, dir, true, true)
	opts.Namespace = "karl-system"
	res, err := GuardedApplyAgainstResourcePorts(opts)
	if err != nil {
		t.Fatal(err)
	}
	if !res.Blocked {
		t.Fatalf("res=%+v", res)
	}
}

func TestGuardedApplySavesBaselineAndAppliesCPU(t *testing.T) {
	lease := loadSandboxLease(t, "resourcelease-dryrun-allowed.json")
	dir := t.TempDir()
	prefix := filepath.Join(dir, "karl.slice")
	slice := filepath.Join(prefix, "resourcelease", "khr-runtime-sandbox-lease-allowed")
	if err := os.MkdirAll(slice, 0o755); err != nil {
		t.Fatal(err)
	}
	_ = cgroup.WriteCPUMax(slice, prefix, "max")

	res, err := GuardedApplyAgainstResourcePorts(guardedApplyOpts(lease, dir, true, true))
	if err != nil {
		t.Fatal(err)
	}
	if !res.Applied || res.Blocked {
		t.Fatalf("res=%+v", res)
	}
	if res.Baseline.CgroupCPUPath == "" || res.Baseline.CPUMaxBefore == "" {
		t.Fatalf("baseline=%+v", res.Baseline)
	}
	if res.Verification.State != VerificationStatePass {
		t.Fatalf("verification=%+v", res.Verification)
	}
}

func memoryApplyOpts(lease *crdv1alpha1.ResourceLease, dir, baselineID string) GuardedApplySandboxOptions {
	opts := guardedApplyOpts(lease, dir, true, true)
	opts.BaselineID = baselineID
	return opts
}

func TestGuardedApplyMemoryScaleUpAndRollback(t *testing.T) {
	lease := loadSandboxLease(t, "resourcelease-memory-scale-up.json")
	dir := t.TempDir()
	prefix := filepath.Join(dir, "karl.slice")
	slice := filepath.Join(prefix, "resourcelease", "khr-runtime-sandbox-lease-memory-up")
	if err := os.MkdirAll(slice, 0o755); err != nil {
		t.Fatal(err)
	}
	_ = cgroup.WriteMemoryMax(slice, prefix, cgroup.FormatMemoryValue(128*1024*1024))
	_ = cgroup.WriteMemoryHigh(slice, prefix, cgroup.FormatMemoryValue(128*1024*1024))

	res, err := GuardedApplyAgainstResourcePorts(memoryApplyOpts(lease, dir, "ram-up"))
	if err != nil {
		t.Fatal(err)
	}
	if !res.Applied || res.Verification.State != VerificationStatePass {
		t.Fatalf("res=%+v", res)
	}
	if !res.Verification.NoRestart || !res.Verification.NoRollout {
		t.Fatalf("policy=%+v", res.Verification)
	}

	rb, err := RollbackSandbox(RollbackSandboxOptions{
		Config:          sandboxDryRunCfg(),
		BaselineID:      "ram-up",
		SandboxDir:      dir,
		AllowPathPrefix: prefix,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !rb.RolledBack || rb.Verification.State != VerificationStatePass {
		t.Fatalf("rb=%+v", rb)
	}
}

func TestGuardedApplyMemoryBlockedOverLimit(t *testing.T) {
	lease := loadSandboxLease(t, "resourcelease-memory-blocked-over-limit.json")
	dir := t.TempDir()
	res, err := GuardedApplyAgainstResourcePorts(memoryApplyOpts(lease, dir, "ram-over"))
	if err != nil {
		t.Fatal(err)
	}
	if !res.Blocked {
		t.Fatalf("res=%+v", res)
	}
}

func TestGuardedApplyMemoryBlockedRestartAnnotation(t *testing.T) {
	lease := loadSandboxLease(t, "resourcelease-memory-blocked-restart.json")
	dir := t.TempDir()
	res, err := GuardedApplyAgainstResourcePorts(memoryApplyOpts(lease, dir, "ram-restart"))
	if err != nil {
		t.Fatal(err)
	}
	if !res.Blocked {
		t.Fatalf("res=%+v", res)
	}
}

func TestGuardedApplyMemoryScaleDown(t *testing.T) {
	lease := loadSandboxLease(t, "resourcelease-memory-scale-down.json")
	dir := t.TempDir()
	prefix := filepath.Join(dir, "karl.slice")
	slice := filepath.Join(prefix, "resourcelease", "khr-runtime-sandbox-lease-memory-down")
	if err := os.MkdirAll(slice, 0o755); err != nil {
		t.Fatal(err)
	}
	start := cgroup.FormatMemoryValue(256 * 1024 * 1024)
	_ = cgroup.WriteMemoryMax(slice, prefix, start)
	_ = cgroup.WriteMemoryHigh(slice, prefix, start)

	res, err := GuardedApplyAgainstResourcePorts(memoryApplyOpts(lease, dir, "ram-down"))
	if err != nil {
		t.Fatal(err)
	}
	if !res.Applied {
		t.Fatalf("res=%+v", res)
	}
}

func TestRollbackSandboxRestoresBaseline(t *testing.T) {
	lease := loadSandboxLease(t, "resourcelease-dryrun-allowed.json")
	dir := t.TempDir()
	prefix := filepath.Join(dir, "karl.slice")
	slice := filepath.Join(prefix, "resourcelease", "khr-runtime-sandbox-lease-allowed")
	if err := os.MkdirAll(slice, 0o755); err != nil {
		t.Fatal(err)
	}
	_ = cgroup.WriteCPUMax(slice, prefix, "max")

	apply, err := GuardedApplyAgainstResourcePorts(guardedApplyOpts(lease, dir, true, true))
	if err != nil {
		t.Fatal(err)
	}
	if !apply.Applied {
		t.Fatalf("apply=%+v", apply)
	}
	rb, err := RollbackSandbox(RollbackSandboxOptions{
		Config:          sandboxDryRunCfg(),
		BaselineID:      "test-baseline",
		SandboxDir:      dir,
		AllowPathPrefix: filepath.Join(dir, "karl.slice"),
	})
	if err != nil {
		t.Fatal(err)
	}
	if !rb.RolledBack || rb.Verification.State != VerificationStatePass {
		t.Fatalf("rb=%+v", rb)
	}
}
