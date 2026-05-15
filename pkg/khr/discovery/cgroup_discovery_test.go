package discovery

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/cgroup"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/crdv1alpha1"
)

func mustMkdir(t *testing.T, p string) {
	t.Helper()
	if err := os.MkdirAll(p, 0o755); err != nil {
		t.Fatal(err)
	}
}

func TestRun_FoundHeuristic(t *testing.T) {
	t.Setenv("KHR_TEST_CGROUP_VERSION", "v2")
	root := t.TempDir()
	target := filepath.Join(root, "karl.slice", "karl-shell-dev-linux-systemd-001.scope")
	mustMkdir(t, target)
	cell := &crdv1alpha1.Cell{}
	cell.Spec.ShellRef.Name = "dev-linux-systemd-001"
	cell.Metadata.Name = "my-cell"
	out := Run("agent-1", root, "", cell)
	if out.SelectedPath != target {
		t.Fatalf("selected %q want %q blocked=%v warn=%v", out.SelectedPath, target, out.BlockedReasons, out.Warnings)
	}
	if !out.MutationsForbidden {
		t.Fatal("mutations must be forbidden")
	}
}

func TestRun_NotFoundNonFatal(t *testing.T) {
	t.Setenv("KHR_TEST_CGROUP_VERSION", "v2")
	root := t.TempDir()
	// Intentionally no karl.slice — no heuristic match.
	cell := &crdv1alpha1.Cell{}
	cell.Spec.ShellRef.Name = "missing-shell-xyz"
	cell.Metadata.Name = "missing-cell"
	out := Run("agent-1", root, "", cell)
	if out.SelectedPath != "" {
		t.Fatalf("expected empty selected got %q", out.SelectedPath)
	}
	if len(out.BlockedReasons) == 0 {
		t.Fatal("expected aggregate blocked reason")
	}
}

func TestRun_ProviderHandleRebased(t *testing.T) {
	t.Setenv("KHR_TEST_CGROUP_VERSION", "v2")
	root := t.TempDir()
	hostRel := filepath.Join("karl.slice", "unit.scope")
	mustMkdir(t, filepath.Join(root, hostRel))
	cell := &crdv1alpha1.Cell{}
	cell.Spec.ShellRef.Name = "x"
	cell.Metadata.Name = "y"
	h := map[string]string{"cgroupPath": filepath.Join(cgroup.UnifiedCgroupMount, hostRel)}
	raw, _ := json.Marshal(h)
	cell.Spec.ProviderHandle = raw
	out := Run("agent-1", root, "", cell)
	want := filepath.Join(root, hostRel)
	if out.SelectedPath != want {
		t.Fatalf("selected %q want %q cand=%v blocked=%v", out.SelectedPath, want, out.CandidatePaths, out.BlockedReasons)
	}
}

func TestRun_InvalidProviderHandleWarning(t *testing.T) {
	t.Setenv("KHR_TEST_CGROUP_VERSION", "v2")
	root := t.TempDir()
	mustMkdir(t, filepath.Join(root, "karl.slice", "karl-shell-z.scope"))
	cell := &crdv1alpha1.Cell{}
	cell.Spec.ShellRef.Name = "z"
	cell.Spec.ProviderHandle = json.RawMessage(`{`)
	out := Run("agent-1", root, "", cell)
	if len(out.Warnings) == 0 || out.SelectedPath == "" {
		t.Fatalf("want parse warning and heuristic success, got warn=%v sel=%q", out.Warnings, out.SelectedPath)
	}
}

func TestRun_ExplicitBlockedByPrefixHeuristicWins(t *testing.T) {
	t.Setenv("KHR_TEST_CGROUP_VERSION", "v2")
	root := t.TempDir()
	mustMkdir(t, filepath.Join(root, "outside", "bad"))
	mustMkdir(t, filepath.Join(root, "karl.slice", "karl-shell-dev-linux-systemd-001.scope"))
	cell := &crdv1alpha1.Cell{}
	cell.Spec.ShellRef.Name = "dev-linux-systemd-001"
	cell.Metadata.Name = "demo-cell"
	h := map[string]string{"cgroupPath": filepath.Join(root, "outside", "bad")}
	raw, _ := json.Marshal(h)
	cell.Spec.ProviderHandle = raw
	pfx := filepath.Join(root, "karl.slice")
	out := Run("agent-1", root, pfx, cell)
	want := filepath.Join(root, "karl.slice", "karl-shell-dev-linux-systemd-001.scope")
	if out.SelectedPath != want {
		t.Fatalf("selected %q want %q blocked=%v", out.SelectedPath, want, out.BlockedReasons)
	}
}
