package resourceport

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/cgroup"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/flightrecorder"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/host"
)

const (
	EmissionModeJSONOnly  = "json-only"
	EmissionModeCRPreview = "cr-preview"
	LoopSourceRuntime     = "host-runtime-loop"
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
	EmitCR         bool
	OutputDir      string
	Targets        []SandboxTarget
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
		Source:               LoopSourceRuntime,
		SafetyMode:           "sandbox",
		NoProductionMutation: true,
	}
	if opts.EmitCR {
		res.EmissionMode = EmissionModeCRPreview
	} else {
		res.EmissionMode = EmissionModeJSONOnly
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
			if err := writeCRPreview(opts.OutputDir, iter.ResourcePorts); err != nil {
				return res, err
			}
			res.EmitCRApplied = true
		}

		if i < opts.Iterations && opts.Interval > 0 {
			time.Sleep(opts.Interval)
		}
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

func writeCRPreview(dir string, ports []Candidate) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	for _, p := range ports {
		name := filepath.Join(dir, "resourceport-"+p.Metadata.Name+".json")
		b, err := json.MarshalIndent(p, "", "  ")
		if err != nil {
			return err
		}
		if err := os.WriteFile(name, b, 0o644); err != nil {
			return err
		}
	}
	return nil
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
