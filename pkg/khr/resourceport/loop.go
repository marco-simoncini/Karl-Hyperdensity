package resourceport

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/cgroup"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/flightrecorder"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/host"
)

const (
	EmissionModeObservedJSON     = "observed-json"
	EmissionModeJSONOnly         = EmissionModeObservedJSON // KHR-K alias
	EmissionModeCRPreview        = "cr-preview"
	EmissionModeCRAppliedSandbox = "cr-applied-sandbox"
	LoopSourceRuntime            = "host-runtime-loop"
)

// SandboxTarget is a read-only observation target in the sandbox namespace.
type SandboxTarget struct {
	Namespace   string            `json:"namespace"`
	PodName     string            `json:"podName"`
	Container   string            `json:"container,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	CgroupPath  string            `json:"cgroupPath,omitempty"`
	HostID      string            `json:"hostId,omitempty"`
}

// LoopOptions configures a sandbox ResourcePort reporting loop.
type LoopOptions struct {
	Config         *host.Config
	Namespace      string
	Labels         map[string]string
	ClusterContext string
	RequiredContext string
	NodeName       string
	Iterations     int
	Interval       time.Duration
	EmitCR          bool
	ApplyCR         bool
	SandboxConfirm  bool
	CleanupCR       bool
	OutputDir       string
	Targets         []SandboxTarget
}

// LoopIteration is one observation cycle.
type LoopIteration struct {
	Index         int         `json:"index"`
	ObservedAt    string      `json:"observedAt"`
	Targets       []SandboxTarget `json:"targets"`
	ResourcePorts []Candidate `json:"resourcePorts"`
}

// LoopResult is the full loop output (JSON-only by default).
type LoopResult struct {
	Mode            string              `json:"mode"`
	Blocked         bool                `json:"blocked"`
	Reason          string              `json:"reason,omitempty"`
	EmissionMode    string              `json:"emissionMode"`
	EmitCR          bool                `json:"emitCR"`
	EmitCRApplied   bool                `json:"emitCRApplied"`
	ApplyCR         bool                `json:"applyCR"`
	ApplyCRBlocked  bool                `json:"applyCRBlocked,omitempty"`
	ApplyCRReason   string              `json:"applyCRReason,omitempty"`
	ApplyCRApplied  bool                `json:"applyCRApplied"`
	AppliedCRNames  []string            `json:"appliedCRNames,omitempty"`
	SandboxConfirm  bool                `json:"sandboxConfirm"`
	CleanupCR       bool                `json:"cleanupCR,omitempty"`
	Cleanup         *CleanupResult      `json:"cleanup,omitempty"`
	ClusterContext  string              `json:"clusterContext,omitempty"`
	Namespace       string              `json:"namespace"`
	SafetyMode      string              `json:"safetyMode"`
	Source          string              `json:"source"`
	NoProductionMutation bool           `json:"noProductionMutation"`
	Iterations      []LoopIteration   `json:"iterations,omitempty"`
	FlightRecorder  []flightrecorder.Event `json:"flightRecorder,omitempty"`
}

// RunLoop executes sandbox ResourcePort reporting (read-only; no ResourceLease apply).
func RunLoop(opts LoopOptions) (LoopResult, error) {
	flightrecorder.Reset()
	res := LoopResult{
		Mode:                 "resourceport-loop",
		Namespace:            opts.Namespace,
		ClusterContext:       opts.ClusterContext,
		EmitCR:               opts.EmitCR,
		ApplyCR:              opts.ApplyCR,
		SandboxConfirm:       opts.SandboxConfirm,
		CleanupCR:            opts.CleanupCR,
		Source:               LoopSourceRuntime,
		SafetyMode:           "sandbox",
		NoProductionMutation: true,
	}
	if opts.EmitCR {
		res.EmissionMode = EmissionModeCRPreview
	} else {
		res.EmissionMode = EmissionModeObservedJSON
	}

	if gate := validateLoopGate(opts); !gate.Allowed {
		res.Blocked = true
		res.Reason = gate.Reason
		flightrecorder.Record("loop-blocked", gate.Reason, opts.Namespace)
		res.FlightRecorder = flightrecorder.Snapshot()
		return res, nil
	}

	if opts.Iterations <= 0 {
		opts.Iterations = 1
	}
	targets := opts.Targets
	if len(targets) == 0 {
		var err error
		targets, err = DiscoverSandboxTargets(opts)
		if err != nil {
			res.Blocked = true
			res.Reason = err.Error()
			flightrecorder.Record("loop-blocked", res.Reason, "")
			res.FlightRecorder = flightrecorder.Snapshot()
			return res, nil
		}
	}

	flightrecorder.Record("loop-start", "resourceport loop", res.EmissionMode)

	for i := 1; i <= opts.Iterations; i++ {
		now := time.Now().UTC()
		iter := LoopIteration{
			Index:      i,
			ObservedAt: now.Format(time.RFC3339),
			Targets:    targets,
		}
		for _, t := range targets {
			portName := t.PodName + "-port"
			shellRef := fmt.Sprintf("%s/Shell/%s", t.Namespace, t.PodName)
			cellRef := fmt.Sprintf("%s/Cell/%s", t.Namespace, t.PodName)
			c := ReportCandidate(opts.Config, shellRef, cellRef, t.Namespace, portName)
			c.Status.ObservedAt = iter.ObservedAt
			c.Metadata.Labels["khr.karl.io/loop-source"] = LoopSourceRuntime
			iter.ResourcePorts = append(iter.ResourcePorts, c)
		}
		res.Iterations = append(res.Iterations, iter)
		flightrecorder.Record("loop-iteration", fmt.Sprintf("iteration %d", i), fmt.Sprintf("ports=%d", len(iter.ResourcePorts)))

		if opts.EmitCR && opts.OutputDir != "" {
			meta := metaFromConfig(opts.Config, iter.ObservedAt, EmissionModeCRPreview, opts.Namespace)
			crs := CandidatesToCRs(iter.ResourcePorts, meta)
			if _, err := RenderCRFiles(opts.OutputDir, crs); err != nil {
				return res, err
			}
			res.EmitCRApplied = true
		}

		if i < opts.Iterations && opts.Interval > 0 {
			time.Sleep(opts.Interval)
		}
	}

	if opts.ApplyCR && len(res.Iterations) > 0 {
		last := res.Iterations[len(res.Iterations)-1]
		if gate := ValidateApplyCRGate(opts); !gate.Allowed {
			res.ApplyCRBlocked = true
			res.ApplyCRReason = gate.Reason
			flightrecorder.Record("apply-cr-blocked", gate.Reason, "")
		} else {
			meta := metaFromConfig(opts.Config, last.ObservedAt, EmissionModeCRAppliedSandbox, opts.Namespace)
			crs := CandidatesToCRs(last.ResourcePorts, meta)
			names, err := ApplyCRDocuments(opts, crs)
			if err != nil {
				return res, err
			}
			res.ApplyCRApplied = true
			res.AppliedCRNames = names
			res.EmissionMode = EmissionModeCRAppliedSandbox
			flightrecorder.Record("apply-cr", "resourceport CR applied", fmt.Sprintf("count=%d", len(names)))
		}
	}

	if opts.CleanupCR {
		clean, err := CleanupAppliedCRs(opts)
		if err != nil {
			return res, err
		}
		res.Cleanup = &clean
		flightrecorder.Record("cleanup-cr", "resourceport CR cleanup", fmt.Sprintf("deleted=%d", clean.Deleted))
	}

	flightrecorder.Record("loop-complete", "resourceport loop finished", "")
	res.FlightRecorder = flightrecorder.Snapshot()
	return res, nil
}

type loopGate struct {
	Allowed bool
	Reason  string
}

func validateLoopGate(opts LoopOptions) loopGate {
	cfg := opts.Config
	if cfg == nil {
		return loopGate{Reason: "config is nil"}
	}
	if !cfg.Spec.ResourcePortLoopEnabled {
		return loopGate{Reason: "resourcePortLoopEnabled is false (default)"}
	}
	if !cfg.Spec.SandboxMode || !cfg.Spec.LinuxOnly {
		return loopGate{Reason: "sandboxMode and linuxOnly required"}
	}
	ns := host.NormalizeNamespace(opts.Namespace)
	if host.ProductionNamespaceBlocked(ns) {
		return loopGate{Reason: "production namespace blocked: " + ns}
	}
	if !host.NamespaceAllowed(cfg, ns) {
		return loopGate{Reason: "namespace not in allowedNamespaces allowlist"}
	}
	if !host.LabelsAllowlistMatch(cfg, opts.Labels) {
		return loopGate{Reason: "label allowlist mismatch"}
	}
	if opts.RequiredContext != "" && opts.ClusterContext != "" && opts.ClusterContext != opts.RequiredContext {
		return loopGate{Reason: fmt.Sprintf("cluster context %q != required %q", opts.ClusterContext, opts.RequiredContext)}
	}
	return loopGate{Allowed: true}
}

// ObserveCgroupPath returns a read-only cgroup path hint for a target (no writes).
func ObserveCgroupPath(cfg *host.Config, t SandboxTarget) string {
	if t.CgroupPath != "" {
		return t.CgroupPath
	}
	root := "/sys/fs/cgroup"
	if cfg != nil && cfg.Spec.CgroupRoot != "" {
		root = cfg.Spec.CgroupRoot
	}
	ver := cgroup.DetectVersion()
	if ver == cgroup.V2 {
		return filepath.Join(root, "kubepods.slice")
	}
	return root
}

// DiscoverSandboxTargets lists sandbox pods when cluster discovery is enabled.
func DiscoverSandboxTargets(opts LoopOptions) ([]SandboxTarget, error) {
	if opts.ClusterContext == "" {
		return []SandboxTarget{{
			Namespace: opts.Namespace,
			PodName:   "local-sandbox-target",
			Labels:    opts.Labels,
			HostID:    opts.Config.Spec.HostID,
		}}, nil
	}
	return discoverTargetsKubectl(opts)
}

func discoverTargetsKubectl(opts LoopOptions) ([]SandboxTarget, error) {
	selector := "khr.karl.io/sandbox=true"
	if len(opts.Labels) > 0 {
		parts := make([]string, 0, len(opts.Labels))
		for k, v := range opts.Labels {
			parts = append(parts, k+"="+v)
		}
		selector = strings.Join(parts, ",")
	}
	raw, err := runKubectlJSON(opts.ClusterContext, opts.Namespace, selector)
	if err != nil {
		return nil, err
	}
	targets := parsePodList(raw, opts)
	if len(targets) == 0 {
		return []SandboxTarget{{
			Namespace: opts.Namespace,
			PodName:   "khr-runtime-linux-target",
			Labels:    opts.Labels,
			HostID:    opts.Config.Spec.HostID,
		}}, nil
	}
	return targets, nil
}
