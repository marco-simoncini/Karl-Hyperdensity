package host

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/crdv1alpha1"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/flightrecorder"
)

const (
	ApplyStateIdle             = "idle"
	ApplyStateDryRunAllowed    = "dry-run-allowed"
	ApplyStateApplied          = "applied"
	ApplyStateRollbackRestored = "rollback-restored"
	DefaultStaleHeartbeat      = 2 * time.Minute
)

// HeartbeatOptions configures the host-heartbeat loop (read-only).
type HeartbeatOptions struct {
	Config           *Config
	NodeName         string
	Namespace        string
	ClusterContext   string
	RequiredContext  string
	SandboxDir       string
	Iterations       int
	Interval         time.Duration
	StaleThreshold   time.Duration
	OutputPath       string
	PriorHeartbeatAt string
}

// HeartbeatIteration is one heartbeat cycle output.
type HeartbeatIteration struct {
	Index            int                `json:"index"`
	ObservedAt       string             `json:"observedAt"`
	Host             crdv1alpha1.Host   `json:"host"`
	RuntimeSession   RuntimeSession     `json:"runtimeSession"`
	Stale            bool               `json:"stale"`
	StaleReason      string             `json:"staleReason,omitempty"`
	NoMutation       bool               `json:"noMutation"`
	LastApplyState   string             `json:"lastApplyState"`
}

// HeartbeatResult is CLI output for host-heartbeat mode.
type HeartbeatResult struct {
	Mode                 string               `json:"mode"`
	Blocked              bool                 `json:"blocked"`
	Reason               string               `json:"reason,omitempty"`
	Namespace            string               `json:"namespace,omitempty"`
	ClusterContext       string               `json:"clusterContext,omitempty"`
	NoMutation           bool                 `json:"noMutation"`
	NoProductionMutation bool                 `json:"noProductionMutation"`
	StaleDetected        bool                 `json:"staleDetected"`
	Iterations           []HeartbeatIteration `json:"iterations,omitempty"`
	FlightRecorder       []flightrecorder.Event `json:"flightRecorder,omitempty"`
}

// RunHostHeartbeat executes periodic Host status JSON emission (no cluster mutation).
func RunHostHeartbeat(opts HeartbeatOptions) (HeartbeatResult, error) {
	sess := InitRuntimeSession(opts.Config)
	flightrecorder.InitContext(flightrecorder.SessionContext{
		RuntimeSessionID:      sess.RuntimeSessionID,
		HostRuntimeInstanceID: sess.HostRuntimeInstanceID,
		CorrelationID:         sess.CorrelationID,
	})
	res := HeartbeatResult{
		Mode:                 "host-heartbeat",
		Namespace:            NormalizeNamespace(opts.Namespace),
		ClusterContext:       opts.ClusterContext,
		NoMutation:           true,
		NoProductionMutation: true,
	}
	if opts.Config == nil {
		res.Blocked = true
		res.Reason = "config is nil"
		return res, nil
	}
	if !opts.Config.Spec.SandboxMode || !opts.Config.Spec.LinuxOnly {
		res.Blocked = true
		res.Reason = "sandboxMode and linuxOnly required"
		return res, nil
	}
	if ProductionNamespaceBlocked(res.Namespace) {
		res.Blocked = true
		res.Reason = "production namespace blocked: " + res.Namespace
		return res, nil
	}
	if opts.RequiredContext != "" && opts.ClusterContext != "" && opts.ClusterContext != opts.RequiredContext {
		res.Blocked = true
		res.Reason = fmt.Sprintf("cluster context %q != required %q", opts.ClusterContext, opts.RequiredContext)
		return res, nil
	}
	if opts.Iterations <= 0 {
		opts.Iterations = 1
	}
	if opts.StaleThreshold <= 0 {
		opts.StaleThreshold = DefaultStaleHeartbeat
	}

	flightrecorder.Record("host-heartbeat", "loop start", sess.RuntimeSessionID)

	for i := 1; i <= opts.Iterations; i++ {
		now := time.Now().UTC()
		ports := discoverActiveResourcePorts(opts)
		leases := discoverActiveResourceLeases(opts.SandboxDir)
		applyState := discoverLastApplyState(opts.SandboxDir)
		host := BuildHostHeartbeatStatus(opts.Config, opts.NodeName, ports, leases, applyState, now)
		iter := HeartbeatIteration{
			Index:          i,
			ObservedAt:     now.Format(time.RFC3339),
			Host:           host,
			RuntimeSession: CurrentRuntimeSession(),
			NoMutation:     true,
			LastApplyState: applyState,
		}
		if opts.PriorHeartbeatAt != "" && i == 1 {
			if DetectStaleHeartbeat(opts.PriorHeartbeatAt, opts.StaleThreshold, now) {
				iter.Stale = true
				iter.StaleReason = "prior heartbeat older than threshold"
				res.StaleDetected = true
			}
		}
		res.Iterations = append(res.Iterations, iter)
		flightrecorder.Record("host-heartbeat", "tick", fmt.Sprintf("index=%d applyState=%s ports=%d", i, applyState, len(ports)))

		if opts.OutputPath != "" {
			if err := writeHostStatusFile(opts.OutputPath, host); err != nil {
				return res, err
			}
		}
		if i < opts.Iterations && opts.Interval > 0 {
			time.Sleep(opts.Interval)
		}
	}

	flightrecorder.Record("host-heartbeat", "loop complete", "")
	res.FlightRecorder = flightrecorder.Snapshot()
	return res, nil
}

// BuildHostHeartbeatStatus extends Host status with runtime contract fields (KHR-N).
func BuildHostHeartbeatStatus(cfg *Config, nodeName string, ports, leases []crdv1alpha1.ObjectRef, applyState string, now time.Time) crdv1alpha1.Host {
	h := BuildHostStatus(cfg, nodeName, ports)
	h.Status.LastHeartbeatTime = now.UTC().Format(time.RFC3339)
	h.Status.ActiveResourcePorts = append([]crdv1alpha1.ObjectRef(nil), ports...)
	h.Status.ActiveResourceLeases = append([]crdv1alpha1.ObjectRef(nil), leases...)
	h.Status.LastApplyState = applyState
	sess := CurrentRuntimeSession()
	h.Status.RuntimeSessionID = sess.RuntimeSessionID
	h.Status.HostRuntimeInstanceID = sess.HostRuntimeInstanceID
	h.Status.CorrelationID = sess.CorrelationID
	return h
}

// DetectStaleHeartbeat reports whether lastHeartbeat is older than threshold from reference time.
func DetectStaleHeartbeat(lastHeartbeatISO string, threshold time.Duration, reference time.Time) bool {
	lastHeartbeatISO = strings.TrimSpace(lastHeartbeatISO)
	if lastHeartbeatISO == "" {
		return true
	}
	t, err := time.Parse(time.RFC3339, lastHeartbeatISO)
	if err != nil {
		t, err = time.Parse(time.RFC3339Nano, lastHeartbeatISO)
		if err != nil {
			return true
		}
	}
	return reference.Sub(t) > threshold
}

func discoverActiveResourcePorts(opts HeartbeatOptions) []crdv1alpha1.ObjectRef {
	if opts.ClusterContext == "" || opts.Namespace == "" {
		return nil
	}
	selector := "karl.io/sandbox-namespace=" + opts.Namespace
	args := []string{
		"--context", opts.ClusterContext,
		"get", "resourceports",
		"-l", selector,
		"-o", "json",
	}
	out, err := exec.Command("kubectl", args...).Output()
	if err != nil {
		return nil
	}
	var list struct {
		Items []crdv1alpha1.ResourcePort `json:"items"`
	}
	if err := json.Unmarshal(out, &list); err != nil {
		return nil
	}
	refs := make([]crdv1alpha1.ObjectRef, 0, len(list.Items))
	for _, p := range list.Items {
		refs = append(refs, crdv1alpha1.ObjectRef{
			Name:      p.Metadata.Name,
			Namespace: opts.Namespace,
		})
	}
	return refs
}

func discoverActiveResourceLeases(sandboxDir string) []crdv1alpha1.ObjectRef {
	if sandboxDir == "" {
		return nil
	}
	matches, _ := filepath.Glob(filepath.Join(sandboxDir, "apply-evidence-*.json"))
	if len(matches) == 0 {
		return nil
	}
	return []crdv1alpha1.ObjectRef{{
		Name:      filepath.Base(matches[len(matches)-1]),
		Namespace: "khr-runtime-sandbox",
	}}
}

func discoverLastApplyState(sandboxDir string) string {
	if sandboxDir == "" {
		return ApplyStateIdle
	}
	if matches, _ := filepath.Glob(filepath.Join(sandboxDir, "baseline-*.json")); len(matches) > 0 {
		raw, err := os.ReadFile(matches[len(matches)-1])
		if err == nil {
			var bl struct {
				CPUMaxApplied string `json:"cpuMaxApplied"`
			}
			if json.Unmarshal(raw, &bl) == nil && strings.TrimSpace(bl.CPUMaxApplied) != "" {
				return ApplyStateApplied
			}
		}
	}
	if matches, _ := filepath.Glob(filepath.Join(sandboxDir, "apply-evidence-*.json")); len(matches) > 0 {
		return ApplyStateApplied
	}
	return ApplyStateIdle
}

func writeHostStatusFile(path string, host crdv1alpha1.Host) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	raw, err := json.MarshalIndent(host, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, raw, 0o644)
}
