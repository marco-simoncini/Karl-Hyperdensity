package resourcelease

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/cgroup"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/crdv1alpha1"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/host"
)

const (
	ApplyStateBlocked  = "blocked"
	ApplyStateApplied  = "applied"
	ApplyStatePending  = "pending"
	RollbackStateIdle  = "idle"
	RollbackStateDone  = "restored"
	VerificationStatePass = "pass"
	VerificationStateFail = "fail"
)

// GuardedApplySandboxOptions configures opt-in sandbox ResourceLease apply.
type GuardedApplySandboxOptions struct {
	DryRunAgainstPortOptions
	ApplyResourceLease bool
	SandboxConfirm     bool
}

// VerificationOutcome is post-apply cgroup observation.
type VerificationOutcome struct {
	State                  string          `json:"state"`
	ObservedCPUMax         string          `json:"observedCpuMax,omitempty"`
	ExpectedCPUMax         string          `json:"expectedCpuMax,omitempty"`
	ObservedMemoryMax      string          `json:"observedMemoryMax,omitempty"`
	ExpectedMemoryMax      string          `json:"expectedMemoryMax,omitempty"`
	ObservedMemoryHigh     string          `json:"observedMemoryHigh,omitempty"`
	ExpectedMemoryHigh     string          `json:"expectedMemoryHigh,omitempty"`
	NoRestart              bool            `json:"noRestart"`
	NoRollout              bool            `json:"noRollout"`
	NoRecreate             bool            `json:"noRecreate"`
	NoProductionMutation   bool            `json:"noProductionMutation"`
	LiveScalePolicy        LiveScalePolicy `json:"liveScalePolicy,omitempty"`
}

// GuardedApplySandboxResult is CLI output for resourcelease-guarded-apply.
type GuardedApplySandboxResult struct {
	Mode                   string                  `json:"mode"`
	Applied                bool                    `json:"applied"`
	Blocked                bool                    `json:"blocked"`
	ApplyState             string                  `json:"applyState"`
	Reason                 string                  `json:"reason,omitempty"`
	BlockedReason          string                  `json:"blockedReason,omitempty"`
	DryRun                 DryRunAgainstPortResult   `json:"dryRun"`
	Baseline               Baseline                `json:"baseline"`
	BaselineRef            string                  `json:"baselineRef,omitempty"`
	ApplyEvidenceRef       string                  `json:"applyEvidenceRef,omitempty"`
	Verification           VerificationOutcome     `json:"verification"`
	RollbackState          string                  `json:"rollbackState"`
	SafetyGates            []string                `json:"safetyGates,omitempty"`
	CgroupPath             string                  `json:"cgroupPath,omitempty"`
	NoProductionMutation   bool                    `json:"noProductionMutation"`
	Namespace              string                  `json:"namespace,omitempty"`
	ClusterContext         string                  `json:"clusterContext,omitempty"`
}

// RollbackSandboxOptions restores a captured cgroup baseline.
type RollbackSandboxOptions struct {
	Config         *host.Config
	BaselineID     string
	SandboxDir     string
	AllowPathPrefix string
}

// RollbackSandboxResult is CLI output for resourcelease-rollback.
type RollbackSandboxResult struct {
	Mode                 string           `json:"mode"`
	RolledBack           bool             `json:"rolledBack"`
	Blocked              bool             `json:"blocked"`
	Reason               string           `json:"reason,omitempty"`
	Baseline             Baseline         `json:"baseline"`
	Rollback             RollbackResult   `json:"rollback"`
	Verification         VerificationOutcome `json:"verification"`
	RollbackState        string           `json:"rollbackState"`
	NoProductionMutation bool             `json:"noProductionMutation"`
}

// GuardedApplyAgainstResourcePorts runs dry-run then optional CPU cgroup apply in sandbox.
func GuardedApplyAgainstResourcePorts(opts GuardedApplySandboxOptions) (GuardedApplySandboxResult, error) {
	res := GuardedApplySandboxResult{
		Mode:                 "resourcelease-guarded-apply",
		Namespace:            host.NormalizeNamespace(opts.Namespace),
		ClusterContext:       opts.ClusterContext,
		NoProductionMutation: true,
		ApplyState:           ApplyStatePending,
		RollbackState:        RollbackStateIdle,
		Verification: VerificationOutcome{
			NoRestart:            true,
			NoProductionMutation: true,
		},
	}
	policy := DefaultLiveScalePolicy()
	res.Verification.LiveScalePolicy = policy
	res.Verification.NoRollout = policy.NoRollout
	res.Verification.NoRecreate = policy.NoRecreate
	res.SafetyGates = []string{
		"linux-only",
		"sandbox-namespace-allowlist",
		"label-allowlist",
		"cluster-context-guard",
		"dry-run-required",
		"cpu-and-ram-live-scale",
		"sandbox-cpu-cap-500m",
		"sandbox-memory-delta-cap",
		"no-restart",
		"no-rollout",
		"no-recreate",
		"no-production-mutation",
	}

	if gate := validateGuardedApplyGate(opts); !gate.Allowed {
		return finishGuardedBlocked(res, gate.Reason), nil
	}
	if opts.Lease != nil {
		if err := ValidateLiveScaleLease(opts.Lease); err != nil {
			return finishGuardedBlocked(res, err.Error()), nil
		}
	}

	dr, err := DryRunAgainstResourcePorts(opts.DryRunAgainstPortOptions)
	res.DryRun = dr
	if err != nil {
		return res, err
	}
	if dr.Blocked || !dr.Allowed || dr.DryRunDecision != DryRunDecisionAllowed {
		return finishGuardedBlocked(res, "dry-run not allowed: "+dr.Reason), nil
	}
	if dr.RollbackPlanRef == "" {
		return finishGuardedBlocked(res, "rollbackPlanRef required for guarded apply"), nil
	}
	if dr.VerificationPlanRef == "" {
		return finishGuardedBlocked(res, "verificationPlanRef required for guarded apply"), nil
	}
	if len(dr.RollbackPlan) == 0 || len(dr.VerificationPlan) == 0 {
		return finishGuardedBlocked(res, "rollbackPlan and verificationPlan must be present"), nil
	}

	cgPath, prefix, err := sandboxCgroupPathAndPrefix(opts.Config, opts.Lease, opts.SandboxDir)
	if err != nil {
		return finishGuardedBlocked(res, err.Error()), nil
	}
	res.CgroupPath = cgPath

	switch dr.Resource {
	case "cpu":
		return applySandboxCPU(opts, res, dr, cgPath, prefix)
	case "memory":
		return applySandboxMemory(opts, res, dr, cgPath, prefix)
	default:
		return finishGuardedBlocked(res, fmt.Sprintf("unsupported resource %q for guarded apply", dr.Resource)), nil
	}
}

func applySandboxCPU(opts GuardedApplySandboxOptions, res GuardedApplySandboxResult, dr DryRunAgainstPortResult, cgPath, prefix string) (GuardedApplySandboxResult, error) {
	if dr.RequestedAmount == nil || dr.RequestedAmount.MilliCPU <= 0 {
		return finishGuardedBlocked(res, "requested cpu amount required"), nil
	}
	if over, reason := amountOverSandboxLimitWithConfig(opts.Config, "cpu", dr.RequestedAmount); over {
		return finishGuardedBlocked(res, reason), nil
	}
	bl, err := CaptureCgroupBaseline(opts.BaselineID, opts.SandboxDir, cgPath, prefix, dr.RequestedAmount.MilliCPU)
	if err != nil {
		return res, err
	}
	res.Baseline = bl
	res.BaselineRef = baselineRef(opts.Namespace, opts.BaselineID)

	expected := cgroup.MilliCPUToMaxLine(dr.RequestedAmount.MilliCPU)
	if err := cgroup.WriteCPUMax(cgPath, prefix, expected); err != nil {
		return res, err
	}
	bl.CPUMaxApplied = expected
	_ = SaveBaseline(bl)

	_ = writeApplyEvidence(filepath.Join(opts.SandboxDir, "apply-evidence-"+opts.BaselineID+".json"), bl, expected)
	res.ApplyEvidenceRef = applyEvidenceRef(opts.Namespace, opts.BaselineID)

	observed, err := cgroup.ReadCPUMax(cgPath, prefix)
	if err != nil {
		res.Verification.State = VerificationStateFail
		return finishGuardedBlocked(res, "verification read failed: "+err.Error()), nil
	}
	res.Verification.ObservedCPUMax = observed
	res.Verification.ExpectedCPUMax = expected
	if strings.TrimSpace(observed) != strings.TrimSpace(expected) {
		res.Verification.State = VerificationStateFail
		return finishGuardedBlocked(res, fmt.Sprintf("cpu.max mismatch: got %q want %q", observed, expected)), nil
	}
	res.Verification.State = VerificationStatePass
	res.Applied = true
	res.Blocked = false
	res.ApplyState = ApplyStateApplied
	res.Reason = "sandbox cpu.max applied under " + cgPath
	return res, nil
}

func applySandboxMemory(opts GuardedApplySandboxOptions, res GuardedApplySandboxResult, dr DryRunAgainstPortResult, cgPath, prefix string) (GuardedApplySandboxResult, error) {
	if dr.RequestedAmount == nil || dr.RequestedAmount.Bytes <= 0 {
		return finishGuardedBlocked(res, "requested memory bytes delta required"), nil
	}
	if over, reason := amountOverSandboxLimitWithConfig(opts.Config, "memory", dr.RequestedAmount); over {
		return finishGuardedBlocked(res, reason), nil
	}
	_, _, _, mode, ok := opts.Lease.Spec.EffectiveTransfer()
	if !ok {
		return finishGuardedBlocked(res, "lease transfer required"), nil
	}
	bl, err := CaptureMemoryCgroupBaseline(opts.BaselineID, opts.SandboxDir, cgPath, prefix, "memory", mode)
	if err != nil {
		return res, err
	}
	res.Baseline = bl
	res.BaselineRef = baselineRef(opts.Namespace, opts.BaselineID)

	currentBytes, currentUnlimited, err := parseMemoryBaseline(bl.MemoryMaxBefore)
	if err != nil {
		return finishGuardedBlocked(res, err.Error()), nil
	}
	target, err := cgroup.ComputeMemoryTarget(currentBytes, currentUnlimited, dr.RequestedAmount.Bytes, mode)
	if err != nil {
		return finishGuardedBlocked(res, err.Error()), nil
	}
	expected := cgroup.FormatMemoryValue(target)
	if err := cgroup.WriteMemoryHigh(cgPath, prefix, expected); err != nil {
		return res, err
	}
	if err := cgroup.WriteMemoryMax(cgPath, prefix, expected); err != nil {
		return res, err
	}
	bl.MemoryMaxApplied = expected
	bl.MemoryHighApplied = expected
	_ = SaveBaseline(bl)

	_ = writeApplyEvidence(filepath.Join(opts.SandboxDir, "apply-evidence-"+opts.BaselineID+".json"), bl, expected)
	res.ApplyEvidenceRef = applyEvidenceRef(opts.Namespace, opts.BaselineID)

	obsMax, err := cgroup.ReadMemoryMax(cgPath, prefix)
	if err != nil {
		res.Verification.State = VerificationStateFail
		return finishGuardedBlocked(res, "memory.max read failed: "+err.Error()), nil
	}
	obsHigh, err := cgroup.ReadMemoryHigh(cgPath, prefix)
	if err != nil {
		res.Verification.State = VerificationStateFail
		return finishGuardedBlocked(res, "memory.high read failed: "+err.Error()), nil
	}
	res.Verification.ObservedMemoryMax = obsMax
	res.Verification.ObservedMemoryHigh = obsHigh
	res.Verification.ExpectedMemoryMax = expected
	res.Verification.ExpectedMemoryHigh = expected
	if strings.TrimSpace(obsMax) != strings.TrimSpace(expected) || strings.TrimSpace(obsHigh) != strings.TrimSpace(expected) {
		res.Verification.State = VerificationStateFail
		return finishGuardedBlocked(res, fmt.Sprintf("memory mismatch: max=%q high=%q want %q", obsMax, obsHigh, expected)), nil
	}
	res.Verification.State = VerificationStatePass
	res.Applied = true
	res.Blocked = false
	res.ApplyState = ApplyStateApplied
	res.Reason = fmt.Sprintf("sandbox memory.high/memory.max applied (%s) under %s", mode, cgPath)
	return res, nil
}

func parseMemoryBaseline(val string) (bytes int64, unlimited bool, err error) {
	return cgroup.ParseMemoryValue(val)
}

// RollbackSandbox restores cgroup baseline captured before guarded apply.
func RollbackSandbox(opts RollbackSandboxOptions) (RollbackSandboxResult, error) {
	out := RollbackSandboxResult{
		Mode:                 "resourcelease-rollback",
		NoProductionMutation: true,
		RollbackState:        RollbackStateIdle,
		Verification: VerificationOutcome{
			NoRestart:            true,
			NoProductionMutation: true,
		},
	}
	if opts.SandboxDir == "" {
		return finishRollbackBlocked(out, "sandboxDir required"), nil
	}
	bl, err := LoadBaseline(opts.BaselineID, opts.SandboxDir)
	if err != nil {
		return finishRollbackBlocked(out, err.Error()), nil
	}
	out.Baseline = bl
	prefix := opts.AllowPathPrefix
	if bl.Extra != nil && bl.Extra["allowPathPrefix"] != "" {
		prefix = bl.Extra["allowPathPrefix"]
	} else if prefix == "" && opts.Config != nil && len(opts.Config.Spec.AllowPathPrefixes) > 0 {
		if _, statErr := os.Stat(opts.Config.Spec.AllowPathPrefixes[0]); statErr == nil && strings.HasPrefix(bl.CgroupCPUPath, opts.Config.Spec.AllowPathPrefixes[0]) {
			prefix = opts.Config.Spec.AllowPathPrefixes[0]
		}
	}
	if prefix == "" {
		prefix = opts.SandboxDir
	}
	if bl.CgroupCPUPath == "" {
		marker := RollbackBaseline(bl)
		out.Rollback = marker
		out.RolledBack = marker.RolledBack
		out.RollbackState = RollbackStateDone
		return out, nil
	}
	policy := DefaultLiveScalePolicy()
	out.Verification.LiveScalePolicy = policy
	out.Verification.NoRollout = policy.NoRollout
	out.Verification.NoRecreate = policy.NoRecreate

	switch bl.Resource {
	case "memory":
		return rollbackMemoryCgroup(out, bl, prefix)
	default:
		return rollbackCPUCgroup(out, bl, prefix)
	}
}

func rollbackCPUCgroup(out RollbackSandboxResult, bl Baseline, prefix string) (RollbackSandboxResult, error) {
	restore := bl.CPUMaxBefore
	if restore == "" {
		restore = "max"
	}
	if err := cgroup.WriteCPUMax(bl.CgroupCPUPath, prefix, restore); err != nil {
		return finishRollbackBlocked(out, err.Error()), nil
	}
	observed, err := cgroup.ReadCPUMax(bl.CgroupCPUPath, prefix)
	if err != nil {
		out.Verification.State = VerificationStateFail
		return finishRollbackBlocked(out, err.Error()), nil
	}
	out.Verification.ObservedCPUMax = observed
	out.Verification.ExpectedCPUMax = restore
	if strings.TrimSpace(observed) != strings.TrimSpace(restore) {
		out.Verification.State = VerificationStateFail
		return finishRollbackBlocked(out, fmt.Sprintf("rollback verify mismatch: %q vs %q", observed, restore)), nil
	}
	out.Verification.State = VerificationStatePass
	out.RolledBack = true
	out.RollbackState = RollbackStateDone
	out.Rollback = RollbackResult{RolledBack: true, Actions: []string{"restored cpu.max from baseline"}}
	return out, nil
}

func rollbackMemoryCgroup(out RollbackSandboxResult, bl Baseline, prefix string) (RollbackSandboxResult, error) {
	restoreMax := bl.MemoryMaxBefore
	if restoreMax == "" {
		restoreMax = "max"
	}
	restoreHigh := bl.MemoryHighBefore
	if restoreHigh == "" {
		restoreHigh = restoreMax
	}
	if err := cgroup.WriteMemoryHigh(bl.CgroupCPUPath, prefix, restoreHigh); err != nil {
		return finishRollbackBlocked(out, err.Error()), nil
	}
	if err := cgroup.WriteMemoryMax(bl.CgroupCPUPath, prefix, restoreMax); err != nil {
		return finishRollbackBlocked(out, err.Error()), nil
	}
	obsMax, err := cgroup.ReadMemoryMax(bl.CgroupCPUPath, prefix)
	if err != nil {
		out.Verification.State = VerificationStateFail
		return finishRollbackBlocked(out, err.Error()), nil
	}
	obsHigh, err := cgroup.ReadMemoryHigh(bl.CgroupCPUPath, prefix)
	if err != nil {
		out.Verification.State = VerificationStateFail
		return finishRollbackBlocked(out, err.Error()), nil
	}
	out.Verification.ObservedMemoryMax = obsMax
	out.Verification.ExpectedMemoryMax = restoreMax
	out.Verification.ObservedMemoryHigh = obsHigh
	out.Verification.ExpectedMemoryHigh = restoreHigh
	if strings.TrimSpace(obsMax) != strings.TrimSpace(restoreMax) {
		out.Verification.State = VerificationStateFail
		return finishRollbackBlocked(out, fmt.Sprintf("rollback memory.max mismatch: %q vs %q", obsMax, restoreMax)), nil
	}
	out.Verification.State = VerificationStatePass
	out.RolledBack = true
	out.RollbackState = RollbackStateDone
	out.Rollback = RollbackResult{RolledBack: true, Actions: []string{"restored memory.high/memory.max from baseline"}}
	return out, nil
}

type applyGate struct {
	Allowed bool
	Reason  string
}

func validateGuardedApplyGate(opts GuardedApplySandboxOptions) applyGate {
	if !opts.ApplyResourceLease {
		return applyGate{Reason: "apply-resourcelease is false (default)"}
	}
	if !opts.SandboxConfirm {
		return applyGate{Reason: "missing --i-understand-this-is-sandbox confirmation"}
	}
	if gate := validateDryRunAgainstPortGate(opts.DryRunAgainstPortOptions); !gate.Allowed {
		return applyGate{Reason: gate.Reason}
	}
	ns := host.NormalizeNamespace(opts.Namespace)
	if ns != "khr-runtime-sandbox" {
		return applyGate{Reason: fmt.Sprintf("namespace must be khr-runtime-sandbox for guarded apply (got %q)", ns)}
	}
	return applyGate{Allowed: true}
}

func sandboxCgroupPathAndPrefix(cfg *host.Config, lease *crdv1alpha1.ResourceLease, sandboxDir string) (path, prefix string, err error) {
	name := "lease-apply"
	if lease != nil && lease.Metadata.Name != "" {
		name = sanitizeLeaseCgroupName(lease.Metadata.Name)
	}
	if cfg != nil && len(cfg.Spec.AllowPathPrefixes) > 0 {
		prefix = cfg.Spec.AllowPathPrefixes[0]
		p := filepath.Join(prefix, "resourcelease", name)
		if _, statErr := os.Stat(prefix); statErr == nil {
			return p, prefix, nil
		}
	}
	if sandboxDir != "" {
		prefix = sandboxDir
		return filepath.Join(sandboxDir, "cgroup", name), prefix, nil
	}
	root := cgroup.UnifiedCgroupMount
	if cfg != nil && cfg.Spec.CgroupRoot != "" {
		root = cfg.Spec.CgroupRoot
	}
	prefix = filepath.Join(root, "karl.slice")
	return filepath.Join(prefix, "khr-runtime-sandbox", "resourcelease", name), prefix, nil
}

func sanitizeLeaseCgroupName(s string) string {
	s = strings.ToLower(s)
	var b strings.Builder
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9', r == '-':
			b.WriteRune(r)
		default:
			b.WriteRune('-')
		}
	}
	out := strings.Trim(b.String(), "-")
	if out == "" {
		return "lease-apply"
	}
	return out
}

func baselineRef(namespace, id string) string {
	return fmt.Sprintf("%s/Baseline/%s", namespace, id)
}

func applyEvidenceRef(namespace, id string) string {
	return fmt.Sprintf("%s/ApplyEvidence/%s", namespace, id)
}

func finishGuardedBlocked(res GuardedApplySandboxResult, reason string) GuardedApplySandboxResult {
	res.Applied = false
	res.Blocked = true
	res.ApplyState = ApplyStateBlocked
	res.Reason = reason
	res.BlockedReason = reason
	return res
}

func finishRollbackBlocked(out RollbackSandboxResult, reason string) RollbackSandboxResult {
	out.Blocked = true
	out.Reason = reason
	return out
}

func writeApplyEvidence(path string, bl Baseline, expected string) error {
	payload := map[string]any{
		"baselineId": bl.ID,
		"cgroupPath": bl.CgroupCPUPath,
		"cpuMaxApplied": expected,
		"at": time.Now().UTC().Format(time.RFC3339),
	}
	raw, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, raw, 0o644)
}
