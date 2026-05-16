package resourcefuture

import (
	"fmt"
	"time"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/crdv1alpha1"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/host"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/lanediscovery"
)

// Options configures resourcefuture-simulate.
type Options struct {
	Config          *host.Config
	ClusterContext  string
	RequiredContext string
	NodeName        string
}

// Run executes read-only ResourceFuture simulation (KHR-R).
func Run(opts Options) (SimulationResult, error) {
	res := SimulationResult{
		Mode:           ModeSimulate,
		ClusterContext: opts.ClusterContext,
		Safety:         DefaultSafetyPolicy(),
		Summary:        map[string]any{},
	}
	if gate := validateGate(opts); !gate.Allowed {
		res.Blocked = true
		res.Reason = gate.Reason
		return res, nil
	}

	ld, err := lanediscovery.Run(lanediscovery.Options{
		Config:          opts.Config,
		ClusterContext:  opts.ClusterContext,
		RequiredContext: opts.RequiredContext,
	})
	if err != nil {
		return res, err
	}
	if ld.Blocked {
		res.Blocked = true
		res.Reason = "lane discovery blocked: " + ld.Reason
		return res, nil
	}

	nodeName := opts.NodeName
	ports := []crdv1alpha1.ObjectRef{}
	for _, p := range ld.DiscoveredResourcePorts {
		if p.ClusterObserved {
			parts := splitRef(p.Ref)
			if len(parts) >= 2 {
				ports = append(ports, crdv1alpha1.ObjectRef{Name: parts[1], Namespace: parts[0]})
			}
		}
	}
	hostStatus := host.BuildHostStatus(opts.Config, nodeName, ports)
	now := time.Now().UTC().Format(time.RFC3339)
	res.ObservedAt = now
	res.ClusterContext = ld.ClusterContext

	posture := PostureSummary{
		Sandbox:               opts.Config.Spec.SandboxMode,
		Preview:               !opts.Config.Spec.SandboxApplyEnabled,
		ProductionUnsupported: true,
		LaneCount:             len(ld.LaneCapabilities),
		BlockedStateCount:     len(ld.BlockedStates),
		Blockers:              []string{"KHR-R simulation only", "no autonomous apply"},
	}

	res.Input = SimulateInput{
		HostStatus:     hostStatus,
		LaneDiscovery:  ld,
		ResourcePorts:  ld.DiscoveredResourcePorts,
		ResourceLeases: synthesizeLeaseRefs(ld),
		PostureSummary: posture,
	}

	plans, sat, blocked, compat, live, restart, bundle := buildForecasts(
		ld.DiscoveredCells, ld.DiscoveredResourcePorts,
	)
	res.CandidateScalePlans = plans
	res.SaturationForecast = sat
	res.BlockedConstraints = blocked
	res.CompatibilityFallbackPrediction = compat
	res.LiveInPlaceEligibility = live
	res.RestartRequiredPrediction = restart
	res.Forecasts = bundle

	res.Summary["candidatePlans"] = len(plans)
	res.Summary["saturationEntries"] = len(sat)
	res.Summary["blockedConstraints"] = len(blocked)
	res.Summary["liveInPlaceEligibleCount"] = countEligible(live)
	res.Summary["nativeLiveEligibleCount"] = countEligibleLane(live, lanediscovery.LaneNativeLive)
	res.Summary["restartRequiredCount"] = countRestartRequired(restart)
	res.Summary["compatibilityFallbackCount"] = len(compat)
	return res, nil
}

func validateGate(opts Options) struct {
	Allowed bool
	Reason  string
} {
	cfg := opts.Config
	if cfg == nil {
		return struct{ Allowed bool; Reason string }{Reason: "config is nil"}
	}
	if !cfg.Spec.SandboxMode {
		return struct{ Allowed bool; Reason string }{Reason: "sandboxMode required"}
	}
	if !cfg.Spec.ResourceFutureSimulationEnabled {
		return struct{ Allowed bool; Reason string }{Reason: "resourceFutureSimulationEnabled must be true"}
	}
	if opts.RequiredContext != "" && opts.ClusterContext != "" && opts.ClusterContext != opts.RequiredContext {
		return struct{ Allowed bool; Reason string }{
			Reason: fmt.Sprintf("cluster context %q != required %q", opts.ClusterContext, opts.RequiredContext),
		}
	}
	return struct{ Allowed bool; Reason string }{Allowed: true}
}

func synthesizeLeaseRefs(ld lanediscovery.Result) []SimulatedLeaseRef {
	out := make([]SimulatedLeaseRef, 0, len(ld.DiscoveredCells))
	for _, c := range ld.DiscoveredCells {
		out = append(out, SimulatedLeaseRef{
			Ref: fmt.Sprintf("%s/ResourceLease/%s-scale-sim", c.Namespace, c.Name),
			Resource: "cpu", Mode: "envelope", DryRunOnly: true,
		})
	}
	return out
}

func splitRef(ref string) []string {
	parts := make([]string, 0, 3)
	cur := ""
	for _, r := range ref {
		if r == '/' {
			if cur != "" {
				parts = append(parts, cur)
				cur = ""
			}
			continue
		}
		cur += string(r)
	}
	if cur != "" {
		parts = append(parts, cur)
	}
	return parts
}

func countEligible(items []LiveInPlaceEligibility) int {
	n := 0
	for _, i := range items {
		if i.Eligible {
			n++
		}
	}
	return n
}

func countEligibleLane(items []LiveInPlaceEligibility, lane string) int {
	n := 0
	for _, i := range items {
		if i.Eligible && i.Lane == lane {
			n++
		}
	}
	return n
}

func countRestartRequired(items []RestartRequiredPrediction) int {
	n := 0
	for _, i := range items {
		if i.Required {
			n++
		}
	}
	return n
}
