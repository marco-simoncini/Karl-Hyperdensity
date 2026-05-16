package resourcelease

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/crdv1alpha1"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/host"
)

const (
	DryRunDecisionAllowed = "allowed"
	DryRunDecisionBlocked = "blocked"
	AnnotationResourcePortRef = "khr.karl.io/resource-port-ref"

	SandboxMaxMilliCPU      int64 = 500
	SandboxMaxMemoryBytes   int64 = 512 * 1024 * 1024
)

// TransferAmount is the requested transfer quantity on a lease.
type TransferAmount struct {
	MilliCPU  int64  `json:"milliCpu,omitempty"`
	Bytes     int64  `json:"bytes,omitempty"`
	Direction string `json:"direction,omitempty"` // scaleUp | scaleDown (memory)
}

// DryRunAgainstPortOptions configures sandbox ResourceLease dry-run against cluster ResourcePorts.
type DryRunAgainstPortOptions struct {
	Config          *host.Config
	Lease           *crdv1alpha1.ResourceLease
	Namespace       string
	Labels          map[string]string
	ClusterContext  string
	RequiredContext string
	ResourcePortRef string
	SandboxDir      string
	BaselineID      string
	Ports           []crdv1alpha1.ResourcePort
}

// DryRunAgainstPortResult is CLI JSON for resourcelease-dryrun mode.
type DryRunAgainstPortResult struct {
	Mode                   string       `json:"mode"`
	Allowed                bool         `json:"allowed"`
	Blocked                bool         `json:"blocked"`
	DryRunDecision         string       `json:"dryRunDecision"`
	Reason                 string       `json:"reason,omitempty"`
	BlockedReason          string       `json:"blockedReason,omitempty"`
	Baseline               Baseline     `json:"baseline"`
	RollbackPlan           []string     `json:"rollbackPlan"`
	VerificationPlan       []string     `json:"verificationPlan"`
	RollbackPlanRef        string       `json:"rollbackPlanRef,omitempty"`
	RollbackPlanStatus     string       `json:"rollbackPlanStatus,omitempty"`
	VerificationPlanRef    string       `json:"verificationPlanRef,omitempty"`
	VerificationPlanStatus string       `json:"verificationPlanStatus,omitempty"`
	SourceResourcePortRef  string       `json:"sourceResourcePortRef,omitempty"`
	LeaseKind              string       `json:"leaseKind,omitempty"`
	Resource               string         `json:"resource,omitempty"`
	RequestedAmount        *TransferAmount `json:"requestedAmount,omitempty"`
	ClusterContext         string       `json:"clusterContext,omitempty"`
	Namespace              string       `json:"namespace,omitempty"`
	NoMutation             bool         `json:"noMutation"`
	NoApply                bool         `json:"noApply"`
	MatchedResourcePort    string       `json:"matchedResourcePort,omitempty"`
	ResourcePortCheck      string       `json:"resourcePortCheck,omitempty"`
}

// ListClusterResourcePorts returns cluster-scoped ResourcePorts for the sandbox namespace label.
func ListClusterResourcePorts(clusterContext, sandboxNamespace string) ([]crdv1alpha1.ResourcePort, error) {
	selector := "karl.io/sandbox-namespace=" + sandboxNamespace
	args := []string{
		"--context", clusterContext,
		"get", "resourceports",
		"-l", selector,
		"-o", "json",
	}
	out, err := exec.Command("kubectl", args...).Output()
	if err != nil {
		return nil, fmt.Errorf("kubectl get resourceports: %w", err)
	}
	var list struct {
		Items []crdv1alpha1.ResourcePort `json:"items"`
	}
	if err := json.Unmarshal(out, &list); err != nil {
		return nil, err
	}
	return list.Items, nil
}

// DryRunAgainstResourcePorts evaluates a lease against sandbox ResourcePort CRs (read-only).
func DryRunAgainstResourcePorts(opts DryRunAgainstPortOptions) (DryRunAgainstPortResult, error) {
	res := DryRunAgainstPortResult{
		Mode:           "resourcelease-dryrun",
		Namespace:      host.NormalizeNamespace(opts.Namespace),
		ClusterContext: opts.ClusterContext,
		NoMutation:     true,
		NoApply:        true,
		RollbackPlan:   []string{"revert cgroup cpu.max/memory.max to prior values captured pre-apply"},
		VerificationPlan: []string{
			"read back cgroup limits",
			"compare cgroup.events pressure",
			"confirm process RSS/CPU throttle metrics",
		},
		RollbackPlanStatus:     "planned",
		VerificationPlanStatus: "planned",
	}
	if opts.BaselineID == "" {
		opts.BaselineID = "sandbox-default"
	}
	if opts.SandboxDir == "" {
		opts.SandboxDir = "/tmp/khr-resourcelease-dryrun"
	}
	bl, _ := CaptureBaseline(opts.BaselineID, opts.SandboxDir)
	res.Baseline = bl

	if gate := validateDryRunAgainstPortGate(opts); !gate.Allowed {
		return finishBlocked(res, gate.Reason), nil
	}
	lease := opts.Lease
	if lease == nil {
		return finishBlocked(res, "lease input is nil"), nil
	}

	leaseKind := strings.TrimSpace(lease.Spec.LeaseKind)
	if leaseKind == "" {
		leaseKind = "transfer"
	}
	res.LeaseKind = leaseKind
	if leaseKind != "transfer" && leaseKind != "runtime" {
		return finishBlocked(res, fmt.Sprintf("leaseKind %q not supported in sandbox dry-run", leaseKind)), nil
	}

	if err := validateLeaseLabels(opts, lease); err != nil {
		return finishBlocked(res, err.Error()), nil
	}

	ports := opts.Ports
	if len(ports) == 0 && opts.ClusterContext != "" {
		var err error
		ports, err = ListClusterResourcePorts(opts.ClusterContext, res.Namespace)
		if err != nil {
			return res, err
		}
	}

	port, portRef, err := matchResourcePort(ports, lease, opts.ResourcePortRef)
	if err != nil {
		return finishBlocked(res, err.Error()), nil
	}
	res.SourceResourcePortRef = portRef
	res.MatchedResourcePort = port.Metadata.Name

	_, _, resource, _, ok := lease.Spec.EffectiveTransfer()
	if !ok {
		return finishBlocked(res, "donor/receiver kind and name are required"), nil
	}
	res.Resource = resource

	amount, amountRaw, err := effectiveTransferAmount(lease)
	if err != nil {
		return finishBlocked(res, err.Error()), nil
	}
	res.RequestedAmount = amount
	if over, reason := amountOverSandboxLimitWithConfig(opts.Config, resource, amount); over {
		return finishBlocked(res, reason), nil
	}
	_ = amountRaw

	res.RollbackPlanRef, res.VerificationPlanRef = governancePlanRefs(lease)

	ctx := &CellContext{DonorPlatform: "linux", ReceiverPlatform: "linux"}
	dr := DryRun(lease, port, ctx)
	res.ResourcePortCheck = dr.ResourcePortCheck
	res.RollbackPlan = dr.RollbackPlan
	res.VerificationPlan = dr.VerificationPlan
	if dr.Blocked || !dr.Allowed {
		return finishBlocked(res, dr.Reason), nil
	}

	res.Allowed = true
	res.Blocked = false
	res.DryRunDecision = DryRunDecisionAllowed
	res.Reason = dr.Reason
	return res, nil
}

type dryRunGate struct {
	Allowed bool
	Reason  string
}

func validateDryRunAgainstPortGate(opts DryRunAgainstPortOptions) dryRunGate {
	cfg := opts.Config
	if cfg == nil {
		return dryRunGate{Reason: "config is nil"}
	}
	if !cfg.Spec.SandboxMode || !cfg.Spec.LinuxOnly {
		return dryRunGate{Reason: "sandboxMode and linuxOnly required"}
	}
	ns := host.NormalizeNamespace(opts.Namespace)
	if host.ProductionNamespaceBlocked(ns) {
		return dryRunGate{Reason: "production namespace blocked: " + ns}
	}
	if !host.NamespaceAllowed(cfg, ns) {
		return dryRunGate{Reason: "namespace not in allowedNamespaces allowlist"}
	}
	if !host.LabelsAllowlistMatch(cfg, opts.Labels) {
		return dryRunGate{Reason: "label allowlist mismatch"}
	}
	if opts.RequiredContext != "" && opts.ClusterContext != "" && opts.ClusterContext != opts.RequiredContext {
		return dryRunGate{Reason: fmt.Sprintf("cluster context %q != required %q", opts.ClusterContext, opts.RequiredContext)}
	}
	return dryRunGate{Allowed: true}
}

func validateLeaseLabels(opts DryRunAgainstPortOptions, lease *crdv1alpha1.ResourceLease) error {
	if len(opts.Config.Spec.AllowedLabels) == 0 {
		return nil
	}
	for k, want := range opts.Config.Spec.AllowedLabels {
		got := ""
		if lease.Metadata.Labels != nil {
			got = lease.Metadata.Labels[k]
		}
		if got != want {
			return fmt.Errorf("lease label allowlist mismatch: %s=%q want %q", k, got, want)
		}
	}
	return nil
}

func matchResourcePort(ports []crdv1alpha1.ResourcePort, lease *crdv1alpha1.ResourceLease, explicitRef string) (*crdv1alpha1.ResourcePort, string, error) {
	ref := strings.TrimSpace(explicitRef)
	if ref == "" && lease.Metadata.Annotations != nil {
		ref = strings.TrimSpace(lease.Metadata.Annotations[AnnotationResourcePortRef])
	}
	if ref != "" {
		name := resourcePortNameFromRef(ref)
		for i := range ports {
			if ports[i].Metadata.Name == name {
				return &ports[i], clusterResourcePortRef(ports[i].Metadata.Name), nil
			}
		}
		return nil, "", fmt.Errorf("ResourcePort %q not found in sandbox cluster ports", ref)
	}

	shellRef, cellRef := leaseTargetRefs(lease)
	for i := range ports {
		p := &ports[i]
		if shellRef != "" && p.Spec.ShellRef == shellRef {
			return p, clusterResourcePortRef(p.Metadata.Name), nil
		}
		if cellRef != "" && p.Spec.CellRef == cellRef {
			return p, clusterResourcePortRef(p.Metadata.Name), nil
		}
	}
	return nil, "", fmt.Errorf("no ResourcePort CR matches lease shellRef/cellRef/resourcePortRef")
}

func leaseTargetRefs(lease *crdv1alpha1.ResourceLease) (shellRef, cellRef string) {
	if lease.Spec.Shell.Ref != "" {
		shellRef = lease.Spec.Shell.Ref
	}
	if lease.Spec.Cell.Ref != "" {
		cellRef = lease.Spec.Cell.Ref
	}
	donor, receiver, _, _, ok := lease.Spec.EffectiveTransfer()
	if !ok {
		return shellRef, cellRef
	}
	if shellRef == "" && donor.Kind == "Shell" {
		shellRef = fmt.Sprintf("%s/Shell/%s", donor.Namespace, donor.Name)
	}
	if cellRef == "" && donor.Kind == "Cell" {
		cellRef = fmt.Sprintf("%s/Cell/%s", donor.Namespace, donor.Name)
	}
	if cellRef == "" && receiver.Kind == "Cell" {
		cellRef = fmt.Sprintf("%s/Cell/%s", receiver.Namespace, receiver.Name)
	}
	return shellRef, cellRef
}

func clusterResourcePortRef(name string) string {
	return "cluster/ResourcePort/" + name
}

func resourcePortNameFromRef(ref string) string {
	ref = strings.TrimPrefix(ref, "cluster/ResourcePort/")
	parts := strings.Split(ref, "/")
	return parts[len(parts)-1]
}

func effectiveTransferAmount(lease *crdv1alpha1.ResourceLease) (*TransferAmount, json.RawMessage, error) {
	var raw json.RawMessage
	if lease.Spec.Transfer != nil {
		raw = lease.Spec.Transfer.Amount
	} else {
		raw = lease.Spec.Amount
	}
	if len(raw) == 0 {
		return &TransferAmount{}, raw, nil
	}
	var amt TransferAmount
	if err := json.Unmarshal(raw, &amt); err != nil {
		return nil, raw, fmt.Errorf("invalid transfer amount: %w", err)
	}
	return &amt, raw, nil
}

func amountOverSandboxLimit(resource string, amt *TransferAmount) (bool, string) {
	return amountOverSandboxLimitWithConfig(nil, resource, amt)
}

func amountOverSandboxLimitWithConfig(cfg *host.Config, resource string, amt *TransferAmount) (bool, string) {
	if amt == nil {
		return false, ""
	}
	switch resource {
	case "cpu":
		if amt.MilliCPU > SandboxMaxMilliCPU {
			return true, fmt.Sprintf("requested milliCpu %d exceeds sandbox limit %d", amt.MilliCPU, SandboxMaxMilliCPU)
		}
	case "memory":
		limit := SandboxMaxMemoryDelta(cfg)
		if amt.Bytes > limit {
			return true, fmt.Sprintf("requested memory delta %d exceeds sandbox limit %d", amt.Bytes, limit)
		}
		if amt.Bytes <= 0 {
			return true, "memory scale requires positive bytes delta"
		}
	}
	return false, ""
}

func governancePlanRefs(lease *crdv1alpha1.ResourceLease) (rollbackRef, verifyRef string) {
	if lease.Spec.Governance != nil && lease.Spec.Governance.RollbackPlanRef != nil {
		rp := lease.Spec.Governance.RollbackPlanRef
		if rp.Namespace != "" && rp.Name != "" {
			rollbackRef = fmt.Sprintf("%s/RollbackPlan/%s", rp.Namespace, rp.Name)
		}
	}
	if rollbackRef == "" && lease.Spec.RollbackPlanRef != nil {
		rp := lease.Spec.RollbackPlanRef
		if rp.Namespace != "" && rp.Name != "" {
			rollbackRef = fmt.Sprintf("%s/RollbackPlan/%s", rp.Namespace, rp.Name)
		}
	}
	verifyRef = "sandbox/VerificationPlan/resourcelease-dryrun"
	return rollbackRef, verifyRef
}

func finishBlocked(res DryRunAgainstPortResult, reason string) DryRunAgainstPortResult {
	res.Allowed = false
	res.Blocked = true
	res.DryRunDecision = DryRunDecisionBlocked
	res.Reason = reason
	res.BlockedReason = reason
	return res
}
