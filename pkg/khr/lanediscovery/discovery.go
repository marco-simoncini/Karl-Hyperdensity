package lanediscovery

import (
	"fmt"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/host"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/resourceport"
)

// Options configures read-only multi-lane discovery.
type Options struct {
	Config          *host.Config
	ClusterContext  string
	RequiredContext string
}

// Run executes lane-discovery against the cluster (read-only).
func Run(opts Options) (Result, error) {
	res := Result{
		Mode:           ModeLaneDiscovery,
		ClusterContext: opts.ClusterContext,
		Safety:         DefaultSafetyPolicy(),
		Summary:        map[string]int{},
	}
	if gate := validateGate(opts); !gate.Allowed {
		res.Blocked = true
		res.Reason = gate.Reason
		return res, nil
	}
	ctx := opts.ClusterContext
	if ctx == "" {
		ctx = resourceport.CurrentKubeContext()
		res.ClusterContext = ctx
	}
	if ctx == "" {
		res.Blocked = true
		res.Reason = "cluster context required for lane-discovery"
		return res, nil
	}
	built, err := buildFromCluster(ctx)
	if err != nil {
		return res, err
	}
	built.Safety = DefaultSafetyPolicy()
	return built, nil
}

type gate struct {
	Allowed bool
	Reason  string
}

func validateGate(opts Options) gate {
	cfg := opts.Config
	if cfg == nil {
		return gate{Reason: "config is nil"}
	}
	if !cfg.Spec.SandboxMode {
		return gate{Reason: "sandboxMode required"}
	}
	if !cfg.Spec.LaneDiscoveryEnabled {
		return gate{Reason: "laneDiscoveryEnabled must be true"}
	}
	if opts.RequiredContext != "" && opts.ClusterContext != "" && opts.ClusterContext != opts.RequiredContext {
		return gate{Reason: fmt.Sprintf("cluster context %q != required %q", opts.ClusterContext, opts.RequiredContext)}
	}
	return gate{Allowed: true}
}
