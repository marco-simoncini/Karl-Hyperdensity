package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/certregistry"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/crdv1alpha1"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/policygates"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/flightrecorder"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/host"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/lanediscovery"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/resourcefuture"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/resourcelease"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/resourceport"
)

func main() {
	mode := flag.String("mode", "register-host", "register-host|host-status|host-heartbeat|lane-discovery|resourcefuture-simulate|resourceport-loop|resourceport-cleanup|resourcelease-dryrun|resourcelease-guarded-apply|resourcelease-rollback|report-capabilities|emit-resourceport|dry-run-lease|apply-lease|rollback|flight-recorder")
	nodeName := flag.String("node-name", "", "Kubernetes node name for host-status (default: hostname)")
	clusterContext := flag.String("cluster-context", "", "required cluster context for resourceport-loop discovery")
	emitCR := flag.Bool("emit-cr", false, "write local ResourcePort CR preview files (never kubectl apply by default)")
	applyCR := flag.Bool("apply-cr", false, "kubectl apply ResourcePort CRs (opt-in; requires --emit-cr and sandbox confirmation)")
	sandboxConfirm := flag.Bool("i-understand-this-is-sandbox", false, "explicit confirmation for sandbox CR apply / ResourceLease apply")
	applyResourceLease := flag.Bool("apply-resourcelease", false, "opt-in ResourceLease guarded apply (requires resourcelease-guarded-apply mode)")
	cleanupCR := flag.Bool("cleanup-cr", false, "delete sandbox ResourcePorts managed by karl-host-runtime")
	loopIterations := flag.Int("loop-iterations", 1, "resourceport-loop iteration count")
	loopIntervalMs := flag.Int("loop-interval-ms", 0, "delay between loop iterations")
	loopOutputDir := flag.String("loop-output-dir", "", "directory for CR preview files when --emit-cr=true")
	heartbeatIterations := flag.Int("heartbeat-iterations", 1, "host-heartbeat iteration count")
	heartbeatIntervalMs := flag.Int("heartbeat-interval-ms", 0, "delay between heartbeat iterations")
	heartbeatOutput := flag.String("heartbeat-output", "", "write Host status JSON to path each heartbeat tick")
	priorHeartbeatAt := flag.String("prior-heartbeat-at", "", "RFC3339 timestamp for stale heartbeat simulation")
	configPath := flag.String("config", "", "path to KarlHostRuntimeConfig YAML/JSON")
	leasePath := flag.String("lease-input", "", "ResourceLease JSON for dry-run/apply")
	portPath := flag.String("resource-port-input", "", "ResourcePort JSON for dry-run/apply")
	namespace := flag.String("namespace", "khr-runtime-sandbox", "sandbox namespace for apply gate")
	sandboxDir := flag.String("sandbox-dir", "", "local sandbox directory for guarded apply")
	baselineID := flag.String("baseline-id", "sandbox-default", "rollback baseline id")
	shellRef := flag.String("shell-ref", "khr-runtime-sandbox/Shell/demo", "Shell ref for ResourcePort candidate")
	cellRef := flag.String("cell-ref", "khr-runtime-sandbox/Cell/demo", "Cell ref for ResourcePort candidate")
	portName := flag.String("port-name", "demo-port", "ResourcePort candidate name")
	resourcePortRef := flag.String("resource-port-ref", "", "cluster ResourcePort ref for resourcelease-dryrun (optional)")
	observedPortsPath := flag.String("observed-resourceports-json", "", "ResourcePort list from scope-2 observed-json loop evidence (no cluster CR apply)")
	certRegistryPath := flag.String("cert-registry", "", "KHR-V: certification registry JSON for gated resourcefuture-simulate")
	flag.Parse()

	switch *mode {
	case "flight-recorder":
		emit(flightrecorder.Snapshot())
		return
	case "rollback":
		dir := *sandboxDir
		if dir == "" {
			dir = os.TempDir()
		}
		bl, _ := resourcelease.CaptureBaseline(*baselineID, dir)
		emit(resourcelease.RollbackBaseline(bl))
		return
	}

	cfg, err := loadConfig(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "config: %v\n", err)
		os.Exit(2)
	}
	if *mode == "lane-discovery" {
		if errs := host.ValidateConfigForLaneDiscovery(cfg); len(errs) > 0 {
			fmt.Fprintf(os.Stderr, "config invalid: %v\n", errs)
			os.Exit(2)
		}
	} else if *mode == "resourcefuture-simulate" {
		if errs := host.ValidateConfigForResourceFutureSimulation(cfg); len(errs) > 0 {
			fmt.Fprintf(os.Stderr, "config invalid: %v\n", errs)
			os.Exit(2)
		}
	} else if errs := host.ValidateConfig(cfg); len(errs) > 0 {
		fmt.Fprintf(os.Stderr, "config invalid: %v\n", errs)
		os.Exit(2)
	}

	switch *mode {
	case "lane-discovery":
		ctx := *clusterContext
		if ctx == "" {
			ctx = resourceport.CurrentKubeContext()
		}
		res, err := lanediscovery.Run(lanediscovery.Options{
			Config:          cfg,
			ClusterContext:  ctx,
			RequiredContext: "karl-metal-01@ovh",
		})
		if err != nil {
			fatal(err)
		}
		emit(res)
	case "resourcefuture-simulate":
		ctx := *clusterContext
		if ctx == "" {
			ctx = resourceport.CurrentKubeContext()
		}
		opts := resourcefuture.Options{
			Config:          cfg,
			ClusterContext:  ctx,
			RequiredContext: "karl-metal-01@ovh",
			NodeName:        *nodeName,
		}
		if *certRegistryPath != "" {
			reg, err := certregistry.LoadJSON(*certRegistryPath)
			if err != nil {
				fatal(err)
			}
			opts.Policy = resourcefuture.PolicyContext{
				Registry: &reg,
				Gates:    policygates.DefaultNativeLiveGates(),
				Now:      time.Now().UTC(),
			}
		}
		res, err := resourcefuture.Run(opts)
		if err != nil {
			fatal(err)
		}
		emit(res)
	case "register-host":
		emit(host.RegisterHost(cfg))
	case "host-heartbeat":
		ctx := *clusterContext
		if ctx == "" {
			ctx = resourceport.CurrentKubeContext()
		}
		if *sandboxDir == "" {
			*sandboxDir = filepath.Join(os.TempDir(), "khr-host-heartbeat")
		}
		res, err := host.RunHostHeartbeat(host.HeartbeatOptions{
			Config:           cfg,
			NodeName:         *nodeName,
			Namespace:        *namespace,
			ClusterContext:   ctx,
			RequiredContext:  "karl-metal-01@ovh",
			SandboxDir:       *sandboxDir,
			Iterations:       *heartbeatIterations,
			Interval:         time.Duration(*heartbeatIntervalMs) * time.Millisecond,
			OutputPath:       *heartbeatOutput,
			PriorHeartbeatAt: *priorHeartbeatAt,
		})
		if err != nil {
			fatal(err)
		}
		emit(res)
	case "host-status":
		ports := []crdv1alpha1.ObjectRef{{
			Name:      *portName,
			Namespace: *namespace,
		}}
		flightrecorder.Record("host-status", "host registration status", "")
		emit(host.BuildHostStatus(cfg, *nodeName, ports))
	case "report-capabilities":
		emit(host.ReportCapabilities(cfg))
	case "emit-resourceport":
		emit(resourceport.ReportCandidate(cfg, *shellRef, *cellRef, *namespace, *portName))
	case "resourceport-loop":
		ctx := *clusterContext
		if ctx == "" {
			ctx = resourceport.CurrentKubeContext()
		}
		loopOpts := resourceport.LoopOptions{
			Config:          cfg,
			Namespace:       *namespace,
			Labels:          copyLabels(cfg.Spec.AllowedLabels),
			ClusterContext:  ctx,
			RequiredContext: "karl-metal-01@ovh",
			NodeName:        *nodeName,
			Iterations:      *loopIterations,
			Interval:        time.Duration(*loopIntervalMs) * time.Millisecond,
			EmitCR:          *emitCR,
			ApplyCR:         *applyCR,
			SandboxConfirm:  *sandboxConfirm,
			CleanupCR:       *cleanupCR,
			OutputDir:       *loopOutputDir,
		}
		res, err := resourceport.RunLoop(loopOpts)
		if err != nil {
			fatal(err)
		}
		emit(res)
	case "resourcelease-dryrun":
		if *leasePath == "" {
			fatal(fmt.Errorf("-lease-input is required for resourcelease-dryrun"))
		}
		leaseRaw, err := os.ReadFile(*leasePath)
		if err != nil {
			fatal(err)
		}
		var lease crdv1alpha1.ResourceLease
		if err := json.Unmarshal(leaseRaw, &lease); err != nil {
			fatal(err)
		}
		ctx := *clusterContext
		if ctx == "" {
			ctx = resourceport.CurrentKubeContext()
		}
		var observedPorts []crdv1alpha1.ResourcePort
		if *observedPortsPath != "" {
			observedPorts, err = resourcelease.LoadObservedResourcePortsFromFile(*observedPortsPath)
			if err != nil {
				fatal(err)
			}
		}
		res, err := resourcelease.DryRunAgainstResourcePorts(resourcelease.DryRunAgainstPortOptions{
			Config:          cfg,
			Lease:           &lease,
			Namespace:       *namespace,
			Labels:          copyLabels(cfg.Spec.AllowedLabels),
			ClusterContext:  ctx,
			RequiredContext: "karl-metal-01@ovh",
			ResourcePortRef: *resourcePortRef,
			SandboxDir:      *sandboxDir,
			BaselineID:      *baselineID,
			Ports:           observedPorts,
		})
		if err != nil {
			fatal(err)
		}
		flightrecorder.Record("resourcelease-dryrun", res.DryRunDecision, res.Reason)
		emit(res)
	case "resourcelease-guarded-apply":
		if *leasePath == "" {
			fatal(fmt.Errorf("-lease-input is required for resourcelease-guarded-apply"))
		}
		leaseRaw, err := os.ReadFile(*leasePath)
		if err != nil {
			fatal(err)
		}
		var lease crdv1alpha1.ResourceLease
		if err := json.Unmarshal(leaseRaw, &lease); err != nil {
			fatal(err)
		}
		ctx := *clusterContext
		if ctx == "" {
			ctx = resourceport.CurrentKubeContext()
		}
		var observedPorts []crdv1alpha1.ResourcePort
		if *observedPortsPath != "" {
			observedPorts, err = resourcelease.LoadObservedResourcePortsFromFile(*observedPortsPath)
			if err != nil {
				fatal(err)
			}
		}
		if *sandboxDir == "" {
			*sandboxDir = filepath.Join(os.TempDir(), "khr-resourcelease-guarded-apply")
		}
		host.InitRuntimeSession(cfg)
		host.SetCorrelationID("khr-guarded-apply")
		sess := host.CurrentRuntimeSession()
		flightrecorder.InitContext(flightrecorder.SessionContext{
			RuntimeSessionID:      sess.RuntimeSessionID,
			HostRuntimeInstanceID: sess.HostRuntimeInstanceID,
			CorrelationID:         "khr-guarded-apply",
		})
		res, err := resourcelease.GuardedApplyAgainstResourcePorts(resourcelease.GuardedApplySandboxOptions{
			DryRunAgainstPortOptions: resourcelease.DryRunAgainstPortOptions{
				Config:          cfg,
				Lease:           &lease,
				Namespace:       *namespace,
				Labels:          copyLabels(cfg.Spec.AllowedLabels),
				ClusterContext:  ctx,
				RequiredContext: "karl-metal-01@ovh",
				ResourcePortRef: *resourcePortRef,
				SandboxDir:      *sandboxDir,
				BaselineID:      *baselineID,
				Ports:           observedPorts,
			},
			ApplyResourceLease: *applyResourceLease,
			SandboxConfirm:     *sandboxConfirm,
		})
		if err != nil {
			fatal(err)
		}
		flightrecorder.Record("resourcelease-guarded-apply", res.ApplyState, res.Reason)
		emit(struct {
			resourcelease.GuardedApplySandboxResult
			FlightRecorder []flightrecorder.Event `json:"flightRecorder"`
		}{res, flightrecorder.Snapshot()})
	case "resourcelease-rollback":
		if *sandboxDir == "" {
			*sandboxDir = filepath.Join(os.TempDir(), "khr-resourcelease-guarded-apply")
		}
		prefix := ""
		if len(cfg.Spec.AllowPathPrefixes) > 0 {
			prefix = cfg.Spec.AllowPathPrefixes[0]
		}
		sess := host.CurrentRuntimeSession()
		flightrecorder.InitContext(flightrecorder.SessionContext{
			RuntimeSessionID:      sess.RuntimeSessionID,
			HostRuntimeInstanceID: sess.HostRuntimeInstanceID,
			CorrelationID:         "khr-rollback",
		})
		res, err := resourcelease.RollbackSandbox(resourcelease.RollbackSandboxOptions{
			Config:          cfg,
			BaselineID:      *baselineID,
			SandboxDir:      *sandboxDir,
			AllowPathPrefix: prefix,
		})
		if err != nil {
			fatal(err)
		}
		flightrecorder.Record("resourcelease-rollback", res.RollbackState, res.Reason)
		emit(res)
	case "resourceport-cleanup":
		ctx := *clusterContext
		if ctx == "" {
			ctx = resourceport.CurrentKubeContext()
		}
		clean, err := resourceport.CleanupAppliedCRs(resourceport.LoopOptions{
			Config:         cfg,
			Namespace:      *namespace,
			ClusterContext: ctx,
		})
		if err != nil {
			fatal(err)
		}
		emit(clean)
	case "dry-run-lease":
		lease, port, err := loadLeasePort(*leasePath, *portPath)
		if err != nil {
			fatal(err)
		}
		flightrecorder.Record("dry-run", "lease evaluation", "")
		emit(resourcelease.DryRun(lease, port, &resourcelease.CellContext{DonorPlatform: "linux", ReceiverPlatform: "linux"}))
	case "apply-lease":
		lease, port, err := loadLeasePort(*leasePath, *portPath)
		if err != nil {
			fatal(err)
		}
		labels := copyLabels(cfg.Spec.AllowedLabels)
		res, err := resourcelease.GuardedApply(cfg, lease, port, &resourcelease.CellContext{DonorPlatform: "linux", ReceiverPlatform: "linux"}, *namespace, labels, *sandboxDir)
		if err != nil {
			fatal(err)
		}
		flightrecorder.Record("apply", res.Reason, "")
		emit(res)
	default:
		fmt.Fprintf(os.Stderr, "unknown mode %q\n", *mode)
		os.Exit(2)
	}
}

func copyLabels(in map[string]string) map[string]string {
	if len(in) == 0 {
		return map[string]string{"khr.karl.io/sandbox": "true"}
	}
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func loadConfig(path string) (*host.Config, error) {
	if path == "" {
		return nil, fmt.Errorf("-config is required")
	}
	return host.LoadConfig(path)
}

func loadLeasePort(leasePath, portPath string) (*crdv1alpha1.ResourceLease, *crdv1alpha1.ResourcePort, error) {
	if leasePath == "" || portPath == "" {
		return nil, nil, fmt.Errorf("-lease-input and -resource-port-input are required")
	}
	leaseRaw, err := os.ReadFile(leasePath)
	if err != nil {
		return nil, nil, err
	}
	portRaw, err := os.ReadFile(portPath)
	if err != nil {
		return nil, nil, err
	}
	var lease crdv1alpha1.ResourceLease
	var port crdv1alpha1.ResourcePort
	if err := json.Unmarshal(leaseRaw, &lease); err != nil {
		return nil, nil, err
	}
	if err := json.Unmarshal(portRaw, &port); err != nil {
		return nil, nil, err
	}
	return &lease, &port, nil
}

func emit(v any) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fatal(err)
	}
	fmt.Println(string(b))
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(1)
}
