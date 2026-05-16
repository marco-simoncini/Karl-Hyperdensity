// Command khr-control-graph exports the unified KHR control graph (KHR-X).
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/actionapproval"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/certregistry"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/controlgraph"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/host"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/lanediscovery"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/resourcefuture"
)

func main() {
	configPath := flag.String("config", "", "KarlHostRuntimeConfig path")
	clusterContext := flag.String("cluster-context", "", "kube context")
	registryPath := flag.String("registry", "", "certification registry JSON")
	simulationPath := flag.String("simulation", "", "resourcefuture simulation JSON")
	approvalPath := flag.String("approvals", "", "action approval bundle JSON")
	sprint := flag.String("sprint", "KHR-X", "sprint label")
	out := flag.String("out", "", "output graph JSON")
	discoveryOnly := flag.String("discovery", "", "use existing lane-discovery JSON instead of cluster")
	flag.Parse()

	now := time.Now().UTC()
	var disc lanediscovery.Result
	var err error

	if *discoveryOnly != "" {
		data, err := os.ReadFile(*discoveryOnly)
		if err != nil {
			fatal(err)
		}
		if err := json.Unmarshal(data, &disc); err != nil {
			fatal(err)
		}
	} else {
		if *configPath == "" {
			fmt.Fprintln(os.Stderr, "-config or -discovery required")
			os.Exit(2)
		}
		cfg, err := host.LoadConfig(*configPath)
		if err != nil {
			fatal(err)
		}
		ctx := *clusterContext
		if ctx == "" {
			ctx = "karl-metal-01@ovh"
		}
		disc, err = lanediscovery.Run(lanediscovery.Options{
			Config: cfg, ClusterContext: ctx, RequiredContext: ctx,
		})
		if err != nil {
			fatal(err)
		}
	}

	in := controlgraph.BuildInput{
		Sprint: *sprint, ClusterContext: disc.ClusterContext,
		Discovery: disc, Now: now,
	}
	if t, err := time.Parse(time.RFC3339, disc.ObservedAt); err == nil {
		in.ObservedAt = t
	}

	if *registryPath != "" {
		reg, err := certregistry.LoadJSON(*registryPath)
		if err != nil {
			fatal(err)
		}
		in.Registry = &reg
	}
	if *simulationPath != "" {
		data, err := os.ReadFile(*simulationPath)
		if err != nil {
			fatal(err)
		}
		var sim resourcefuture.SimulationResult
		if err := json.Unmarshal(data, &sim); err != nil {
			fatal(err)
		}
		in.Simulation = &sim
	}
	if *approvalPath != "" {
		b, err := actionapproval.LoadBundleJSON(*approvalPath)
		if err != nil {
			fatal(err)
		}
		in.Approvals = b.Approvals
	}

	g := controlgraph.Build(in)
	data, err := json.MarshalIndent(g, "", "  ")
	if err != nil {
		fatal(err)
	}
	if *out == "" {
		fmt.Println(string(data))
		return
	}
	if err := os.WriteFile(*out, append(data, '\n'), 0o644); err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
