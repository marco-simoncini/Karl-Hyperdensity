package agent

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/audit"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	return filepath.Join("..", "..", "..")
}

func TestRunDryRunCLIUnsafeFlagAuditOnly(t *testing.T) {
	t.Setenv("KHR_TEST_CGROUP_VERSION", "v2")
	root := repoRoot(t)
	lease, err := os.ReadFile(filepath.Join(root, "examples", "khr", "resourcelease-linux-envelope-full.json"))
	if err != nil {
		t.Fatal(err)
	}
	port, err := os.ReadFile(filepath.Join(root, "examples", "khr", "resourceport-linux-envelope-full.json"))
	if err != nil {
		t.Fatal(err)
	}
	out, err := RunDryRunCLI(lease, port, nil, true, "-1ms", "")
	if err != nil {
		t.Fatal(err)
	}
	if !out.MutationsForbidden || !out.FutureApplyGateRequired || !out.UnsafeApplyFlagPresent {
		t.Fatalf("flags: %+v", out)
	}
	if len(out.Audit) != 1 || out.Audit[0].Code != audit.UnsafeApplyFlagNonOperational().Code {
		t.Fatalf("audit: %+v", out.Audit)
	}
	if out.CgroupEnvelopePlan.WouldWrite {
		t.Fatal("cgroup plan must not advertise writes in Sprint 6")
	}
}

func TestRunDryRunCLIJSONStableShape(t *testing.T) {
	t.Setenv("KHR_TEST_CGROUP_VERSION", "v2")
	root := repoRoot(t)
	lease, _ := os.ReadFile(filepath.Join(root, "examples", "khr", "golden", "inputs", "lease-blocked-mode.json"))
	port, _ := os.ReadFile(filepath.Join(root, "examples", "khr", "resourceport-linux-envelope-full.json"))
	out, err := RunDryRunCLI(lease, port, nil, false, "", "")
	if err != nil {
		t.Fatal(err)
	}
	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	if _, ok := m["audit"]; ok {
		t.Fatal("did not expect audit when unsafe flag absent")
	}
}
