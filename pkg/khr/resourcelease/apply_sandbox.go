package resourcelease

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/crdv1alpha1"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/host"
)

// ApplyResult is the outcome of a guarded sandbox apply attempt.
type ApplyResult struct {
	Applied        bool     `json:"applied"`
	Blocked        bool     `json:"blocked"`
	Reason         string   `json:"reason,omitempty"`
	SandboxOnly    bool     `json:"sandboxOnly"`
	BaselineID     string   `json:"baselineId,omitempty"`
	WrittenPaths   []string `json:"writtenPaths,omitempty"`
	DryRun         DryRunResult `json:"dryRun"`
}

// GuardedApply runs dry-run first, then optional sandbox-local marker write when allowed.
func GuardedApply(cfg *host.Config, lease *crdv1alpha1.ResourceLease, port *crdv1alpha1.ResourcePort, ctx *CellContext, namespace string, labels map[string]string, sandboxDir string) (ApplyResult, error) {
	dr := DryRun(lease, port, ctx)
	out := ApplyResult{
		SandboxOnly: true,
		DryRun:      dr,
	}
	if !dr.Allowed || dr.Blocked {
		out.Blocked = true
		out.Reason = "dry-run blocked: " + dr.Reason
		return out, nil
	}
	gate := host.SandboxApplyAllowed(cfg, namespace, labels)
	if !gate.Allowed {
		out.Blocked = true
		out.Reason = gate.Reason
		return out, nil
	}
	if sandboxDir == "" {
		out.Blocked = true
		out.Reason = "sandboxDir required for guarded apply"
		return out, nil
	}
	marker := filepath.Join(sandboxDir, "apply-marker.txt")
	if err := os.MkdirAll(sandboxDir, 0o755); err != nil {
		return out, err
	}
	content := fmt.Sprintf("host=%s ns=%s applied-at=sandbox-only\n", cfg.Spec.HostID, namespace)
	if err := os.WriteFile(marker, []byte(content), 0o644); err != nil {
		return out, err
	}
	out.Applied = true
	out.Blocked = false
	out.Reason = "sandbox marker written only; no production cgroup mutation"
	out.BaselineID = "sandbox-" + cfg.Spec.HostID
	out.WrittenPaths = []string{marker}
	return out, nil
}
