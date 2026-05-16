package resourcelease

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/crdv1alpha1"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/host"
)

func sandboxDryRunCfg() *host.Config {
	cfg := testHostCfg(false)
	cfg.Spec.ResourcePortLoopEnabled = true
	return cfg
}

func sandboxPort() crdv1alpha1.ResourcePort {
	return crdv1alpha1.ResourcePort{
		Metadata: crdv1alpha1.ObjectMeta{
			Name: "khr-runtime-sandbox-demo-port",
			Labels: map[string]string{
				"khr.karl.io/sandbox": "true",
			},
		},
		Spec: crdv1alpha1.ResourcePortSpec{
			ShellRef: "khr-runtime-sandbox/Shell/demo",
			CellRef:  "khr-runtime-sandbox/Cell/demo",
			Ports: crdv1alpha1.ResourcePortsMatrix{
				CPU:    crdv1alpha1.ResourceModes{Modes: []string{"envelope", "static"}},
				Memory: crdv1alpha1.ResourceModes{Modes: []string{"envelope", "scaleUp", "scaleDown", "static"}},
			},
		},
	}
}

func loadSandboxLease(t *testing.T, name string) *crdv1alpha1.ResourceLease {
	t.Helper()
	root := filepath.Join("..", "..", "..", "examples", "khr", "runtime-sandbox")
	raw, err := os.ReadFile(filepath.Join(root, name))
	if err != nil {
		t.Fatal(err)
	}
	var lease crdv1alpha1.ResourceLease
	if err := json.Unmarshal(raw, &lease); err != nil {
		t.Fatal(err)
	}
	return &lease
}

func baseDryRunOpts(lease *crdv1alpha1.ResourceLease, sandboxDir string) DryRunAgainstPortOptions {
	return DryRunAgainstPortOptions{
		Config:     sandboxDryRunCfg(),
		Lease:      lease,
		Namespace:  "khr-runtime-sandbox",
		Labels:     map[string]string{"khr.karl.io/sandbox": "true"},
		Ports:      []crdv1alpha1.ResourcePort{sandboxPort()},
		SandboxDir: sandboxDir,
	}
}

func TestDryRunAgainstResourcePortsAllowed(t *testing.T) {
	lease := loadSandboxLease(t, "resourcelease-dryrun-allowed.json")
	opts := baseDryRunOpts(lease, t.TempDir())
	res, err := DryRunAgainstResourcePorts(opts)
	if err != nil {
		t.Fatal(err)
	}
	if !res.Allowed || res.Blocked || res.DryRunDecision != DryRunDecisionAllowed {
		t.Fatalf("res=%+v", res)
	}
	if res.SourceResourcePortRef == "" || len(res.RollbackPlan) == 0 || len(res.VerificationPlan) == 0 {
		t.Fatalf("plans missing: %+v", res)
	}
}

func TestDryRunAgainstResourcePortsBlockedNamespace(t *testing.T) {
	lease := loadSandboxLease(t, "resourcelease-dryrun-allowed.json")
	opts := baseDryRunOpts(lease, t.TempDir())
	opts.Namespace = "karl-system"
	res, err := DryRunAgainstResourcePorts(opts)
	if err != nil {
		t.Fatal(err)
	}
	if !res.Blocked {
		t.Fatalf("res=%+v", res)
	}
}

func TestDryRunAgainstResourcePortsBlockedMissingLabel(t *testing.T) {
	lease := loadSandboxLease(t, "resourcelease-dryrun-blocked-label.json")
	opts := baseDryRunOpts(lease, t.TempDir())
	res, err := DryRunAgainstResourcePorts(opts)
	if err != nil {
		t.Fatal(err)
	}
	if !res.Blocked {
		t.Fatalf("res=%+v", res)
	}
}

func TestDryRunAgainstResourcePortsBlockedOverLimit(t *testing.T) {
	lease := loadSandboxLease(t, "resourcelease-dryrun-blocked-over-limit.json")
	opts := baseDryRunOpts(lease, t.TempDir())
	res, err := DryRunAgainstResourcePorts(opts)
	if err != nil {
		t.Fatal(err)
	}
	if !res.Blocked || res.DryRunDecision != DryRunDecisionBlocked {
		t.Fatalf("res=%+v", res)
	}
}

func TestDryRunAgainstResourcePortsMemoryScaleUp(t *testing.T) {
	lease := loadSandboxLease(t, "resourcelease-memory-scale-up.json")
	opts := baseDryRunOpts(lease, t.TempDir())
	res, err := DryRunAgainstResourcePorts(opts)
	if err != nil {
		t.Fatal(err)
	}
	if !res.Allowed || res.Resource != "memory" {
		t.Fatalf("res=%+v", res)
	}
}

func TestDryRunAgainstResourcePortsMemoryBlockedOverLimit(t *testing.T) {
	lease := loadSandboxLease(t, "resourcelease-memory-blocked-over-limit.json")
	opts := baseDryRunOpts(lease, t.TempDir())
	res, err := DryRunAgainstResourcePorts(opts)
	if err != nil {
		t.Fatal(err)
	}
	if !res.Blocked {
		t.Fatalf("res=%+v", res)
	}
}

func TestDryRunAgainstResourcePortsBlockedRestartAnnotation(t *testing.T) {
	lease := loadSandboxLease(t, "resourcelease-memory-blocked-restart.json")
	opts := baseDryRunOpts(lease, t.TempDir())
	res, err := DryRunAgainstResourcePorts(opts)
	if err != nil {
		t.Fatal(err)
	}
	if !res.Blocked {
		t.Fatalf("res=%+v", res)
	}
}

func TestDryRunAgainstResourcePortsBlockedMissingPort(t *testing.T) {
	lease := loadSandboxLease(t, "resourcelease-dryrun-blocked-missing-port.json")
	opts := baseDryRunOpts(lease, t.TempDir())
	res, err := DryRunAgainstResourcePorts(opts)
	if err != nil {
		t.Fatal(err)
	}
	if !res.Blocked {
		t.Fatalf("res=%+v", res)
	}
}
