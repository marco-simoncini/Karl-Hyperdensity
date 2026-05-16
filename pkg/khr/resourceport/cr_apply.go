package resourceport

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/crdv1alpha1"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/host"
)

// ApplyCRGate is the result of apply-cr safety checks.
type ApplyCRGate struct {
	Allowed bool
	Reason  string
}

// CleanupResult summarizes sandbox CR deletion.
type CleanupResult struct {
	ClusterContext string `json:"clusterContext,omitempty"`
	Namespace      string `json:"namespace"`
	Deleted        int    `json:"deleted"`
	Selector       string `json:"selector"`
	NoProductionMutation bool `json:"noProductionMutation"`
}

// ValidateApplyCRGate enforces opt-in sandbox apply (never default).
func ValidateApplyCRGate(opts LoopOptions) ApplyCRGate {
	if !opts.ApplyCR {
		return ApplyCRGate{Reason: "apply-cr is false (default)"}
	}
	if !opts.EmitCR {
		return ApplyCRGate{Reason: "emit-cr=true required before apply-cr"}
	}
	if !opts.SandboxConfirm {
		return ApplyCRGate{Reason: "missing --i-understand-this-is-sandbox confirmation"}
	}
	if gate := validateLoopGate(opts); !gate.Allowed {
		return ApplyCRGate{Reason: gate.Reason}
	}
	return ApplyCRGate{Allowed: true}
}

// ApplyCRDocuments kubectl-applies cluster-scoped ResourcePort CRs (sandbox only).
func ApplyCRDocuments(opts LoopOptions, ports []crdv1alpha1.ResourcePort) ([]string, error) {
	if gate := ValidateApplyCRGate(opts); !gate.Allowed {
		return nil, fmt.Errorf("apply-cr blocked: %s", gate.Reason)
	}
	if opts.ClusterContext == "" {
		return nil, fmt.Errorf("cluster context is required for apply-cr")
	}
	var names []string
	for _, rp := range ports {
		paths, err := RenderCRFiles(opts.OutputDir, []crdv1alpha1.ResourcePort{rp})
		if err != nil {
			return names, err
		}
		if len(paths) == 0 {
			continue
		}
		args := []string{"--context", opts.ClusterContext, "apply", "-f", paths[0]}
		out, err := exec.Command("kubectl", args...).CombinedOutput()
		if err != nil {
			return names, FormatKubectlError(err, out)
		}
		names = append(names, rp.Metadata.Name)
	}
	return names, nil
}

// CleanupAppliedCRs deletes ResourcePorts managed by karl-host-runtime in the sandbox namespace label.
func CleanupAppliedCRs(opts LoopOptions) (CleanupResult, error) {
	ns := host.NormalizeNamespace(opts.Namespace)
	res := CleanupResult{
		ClusterContext:       opts.ClusterContext,
		Namespace:            ns,
		NoProductionMutation: true,
		Selector:             LabelManagedBy + "=" + ManagedByValue + "," + LabelSandboxNamespace + "=" + ns,
	}
	if opts.ClusterContext == "" {
		return res, fmt.Errorf("cluster context is required for cleanup-cr")
	}
	if host.ProductionNamespaceBlocked(res.Namespace) {
		return res, fmt.Errorf("production namespace blocked: %s", res.Namespace)
	}
	args := []string{
		"--context", opts.ClusterContext,
		"delete", "resourceports",
		"-l", res.Selector,
		"--ignore-not-found=true",
	}
	out, err := exec.Command("kubectl", args...).CombinedOutput()
	if err != nil {
		return res, FormatKubectlError(err, out)
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, line := range lines {
		if strings.Contains(line, "deleted") {
			res.Deleted++
		}
	}
	if res.Deleted == 0 && strings.Contains(string(out), `"No resources found"`) {
		res.Deleted = 0
	}
	return res, nil
}

// CandidatesToCRs converts loop candidates to cluster-scoped CR documents with metadata.
func CandidatesToCRs(candidates []Candidate, meta CRDocumentMeta) []crdv1alpha1.ResourcePort {
	out := make([]crdv1alpha1.ResourcePort, 0, len(candidates))
	for _, c := range candidates {
		out = append(out, CandidateToCR(c, meta))
	}
	return out
}
