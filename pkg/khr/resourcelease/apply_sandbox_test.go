package resourcelease

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/crdv1alpha1"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/host"
)

func testLeasePort(t *testing.T) (*crdv1alpha1.ResourceLease, *crdv1alpha1.ResourcePort) {
	t.Helper()
	root := filepath.Join("..", "..", "..", "examples", "khr")
	leaseRaw, _ := os.ReadFile(filepath.Join(root, "resourcelease-linux-envelope-dry-run.json"))
	portRaw, _ := os.ReadFile(filepath.Join(root, "resourceport-linux-envelope-for-dryrun.json"))
	var lease crdv1alpha1.ResourceLease
	var port crdv1alpha1.ResourcePort
	if err := json.Unmarshal(leaseRaw, &lease); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(portRaw, &port); err != nil {
		t.Fatal(err)
	}
	return &lease, &port
}

func testHostCfg(enabled bool) *host.Config {
	cfg := &host.Config{}
	cfg.Spec.HostID = "host-test"
	cfg.Spec.LinuxOnly = true
	cfg.Spec.SandboxMode = true
	cfg.Spec.SandboxApplyEnabled = enabled
	cfg.Spec.AllowedNamespaces = []string{"karl-sandbox"}
	cfg.Spec.AllowedLabels = map[string]string{"karl.io/khr-sandbox": "true"}
	return cfg
}

func TestGuardedApplyBlockedByDefault(t *testing.T) {
	lease, port := testLeasePort(t)
	res, err := GuardedApply(testHostCfg(false), lease, port, &CellContext{DonorPlatform: "linux", ReceiverPlatform: "linux"}, "karl-sandbox", map[string]string{"karl.io/khr-sandbox": "true"}, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if res.Applied || !res.Blocked {
		t.Fatalf("res=%+v", res)
	}
}

func TestGuardedApplySandboxWhenEnabled(t *testing.T) {
	lease, port := testLeasePort(t)
	dir := t.TempDir()
	res, err := GuardedApply(testHostCfg(true), lease, port, &CellContext{DonorPlatform: "linux", ReceiverPlatform: "linux"}, "karl-sandbox", map[string]string{"karl.io/khr-sandbox": "true"}, dir)
	if err != nil {
		t.Fatal(err)
	}
	if !res.Applied || res.Blocked {
		t.Fatalf("res=%+v", res)
	}
}

func TestDryRunLease(t *testing.T) {
	lease, port := testLeasePort(t)
	dr := DryRun(lease, port, &CellContext{DonorPlatform: "linux", ReceiverPlatform: "linux"})
	if !dr.Allowed {
		t.Fatalf("dr=%+v", dr)
	}
}
