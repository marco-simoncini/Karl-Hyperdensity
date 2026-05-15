// Command khr-linux-agent is a dry-run-first Linux MVP skeleton (Sprint 5+; discovery Sprint 7).
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/agent"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/crdv1alpha1"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/resourcelease"
)

func main() {
	mode := flag.String("mode", "", "one of: validate-config, dry-run, print-capabilities, discover-cgroups, read-telemetry, collect-evidence")
	configPath := flag.String("config", "", "path to agent YAML/JSON config")
	cellInputPath := flag.String("cell-input", "", "optional path to Cell JSON (discover-cgroups, read-telemetry); required for collect-evidence")
	cgroupRoot := flag.String("cgroup-root", "", "optional cgroup scan root (default /sys/fs/cgroup) for discover-cgroups")
	allowPathPrefix := flag.String("allow-path-prefix", "", "optional path prefix policy (discover-cgroups, read-telemetry, collect-evidence)")
	telemetryCgroupPath := flag.String("cgroup-path", "", "resolved cgroup directory to sample (read-telemetry)")
	telemetryOutputPath := flag.String("telemetry-output", "", "optional path to write the same JSON as stdout (read-telemetry)")
	evidenceOutputPath := flag.String("evidence-output", "", "optional path to write the same JSON as stdout (collect-evidence)")
	evidenceManifestOut := flag.String("evidence-manifest-output", "", "optional path to write artifact manifest JSON (collect-evidence)")
	evidenceDigestOut := flag.String("evidence-digest-output", "", "optional path to write bundle SHA256 hex line (collect-evidence)")
	signingMode := flag.String("signing-mode", "none", "none|local-dev: local integrity signing for collect-evidence (local-dev is not production security)")
	signingKeyFile := flag.String("signing-key-file", "", "Ed25519 private key PEM when -signing-mode=local-dev (collect-evidence)")
	artifactID := flag.String("artifact-id", "", "optional artifact id recorded in evidence manifest (collect-evidence)")
	leasePath := flag.String("lease-input", "", "path to ResourceLease JSON (dry-run)")
	portPath := flag.String("resource-port-input", "", "path to ResourcePort JSON (dry-run)")
	cellCtxPath := flag.String("cell-context", "", "optional path to CellContext JSON (dry-run)")
	allowUnsafe := flag.Bool("allow-unsafe-apply", false, "non-operational in Sprint 6: emits audit only; never enables writes")
	cpuDelta := flag.String("cpu-delta", "", "optional cpu.max delta string for envelope dry-run plan (simulation only)")
	memDelta := flag.String("memory-delta", "", "optional memory.max delta string for envelope dry-run plan (simulation only)")
	flag.Parse()

	out := map[string]interface{}{
		"tool":    "khr-linux-agent",
		"version": agent.AgentVersion,
		"mode":    *mode,
	}

	if *mode == "" {
		out["error"] = "missing -mode"
		emit(out, 2)
	}

	switch *mode {
	case "validate-config":
		if *configPath == "" {
			out["error"] = "missing -config"
			emit(out, 2)
		}
		cfg, err := agent.LoadConfig(*configPath)
		if err != nil {
			out["error"] = err.Error()
			emit(out, 2)
		}
		errs := agent.ValidateConfig(cfg)
		out["valid"] = len(errs) == 0
		out["validationErrors"] = errs
		emit(out, boolExit(len(errs) > 0))

	case "print-capabilities":
		if *configPath == "" {
			out["error"] = "missing -config"
			emit(out, 2)
		}
		cfg, err := agent.LoadConfig(*configPath)
		if err != nil {
			out["error"] = err.Error()
			emit(out, 2)
		}
		if errs := agent.ValidateConfig(cfg); len(errs) > 0 {
			out["validationErrors"] = errs
			out["error"] = "invalid config"
			emit(out, 2)
		}
		b, err := agent.PrintCapabilitiesJSON(cfg, *allowUnsafe)
		if err != nil {
			out["error"] = err.Error()
			emit(out, 2)
		}
		os.Stdout.Write(b)
		os.Stdout.Write([]byte("\n"))
		os.Exit(0)

	case "dry-run":
		if *configPath == "" || *leasePath == "" || *portPath == "" {
			out["error"] = "dry-run requires -config, -lease-input, and -resource-port-input"
			emit(out, 2)
		}
		cfg, err := agent.LoadConfig(*configPath)
		if err != nil {
			out["error"] = err.Error()
			emit(out, 2)
		}
		if errs := agent.ValidateConfig(cfg); len(errs) > 0 {
			out["validationErrors"] = errs
			out["error"] = "invalid config"
			emit(out, 2)
		}
		leaseRaw, err := os.ReadFile(*leasePath)
		if err != nil {
			out["error"] = err.Error()
			emit(out, 2)
		}
		portRaw, err := os.ReadFile(*portPath)
		if err != nil {
			out["error"] = err.Error()
			emit(out, 2)
		}
		var ctx *resourcelease.CellContext
		if *cellCtxPath != "" {
			raw, err := os.ReadFile(*cellCtxPath)
			if err != nil {
				out["error"] = err.Error()
				emit(out, 2)
			}
			ctx = &resourcelease.CellContext{}
			if err := json.Unmarshal(raw, ctx); err != nil {
				out["error"] = err.Error()
				emit(out, 2)
			}
		}
		cliOut, err := agent.RunDryRunCLI(leaseRaw, portRaw, ctx, *allowUnsafe, *cpuDelta, *memDelta)
		if err != nil {
			out["error"] = err.Error()
			emit(out, 2)
		}
		b, err := json.MarshalIndent(cliOut, "", "  ")
		if err != nil {
			out["error"] = err.Error()
			emit(out, 2)
		}
		os.Stdout.Write(b)
		os.Stdout.Write([]byte("\n"))
		os.Exit(0)

	case "discover-cgroups":
		if *configPath == "" {
			out["error"] = "discover-cgroups requires -config"
			emit(out, 2)
		}
		cfg, err := agent.LoadConfig(*configPath)
		if err != nil {
			out["error"] = err.Error()
			emit(out, 2)
		}
		if errs := agent.ValidateConfig(cfg); len(errs) > 0 {
			out["validationErrors"] = errs
			out["error"] = "invalid config"
			emit(out, 2)
		}
		var cell *crdv1alpha1.Cell
		if *cellInputPath != "" {
			raw, err := os.ReadFile(*cellInputPath)
			if err != nil {
				out["error"] = err.Error()
				emit(out, 2)
			}
			cell = &crdv1alpha1.Cell{}
			if err := json.Unmarshal(raw, cell); err != nil {
				out["error"] = err.Error()
				emit(out, 2)
			}
		}
		cliOut := agent.RunDiscoverCgroupsCLI(cfg, cell, *cgroupRoot, *allowPathPrefix)
		b, err := json.MarshalIndent(cliOut, "", "  ")
		if err != nil {
			out["error"] = err.Error()
			emit(out, 2)
		}
		os.Stdout.Write(b)
		os.Stdout.Write([]byte("\n"))
		os.Exit(0)

	case "read-telemetry":
		if *configPath == "" || *telemetryCgroupPath == "" {
			out["error"] = "read-telemetry requires -config and -cgroup-path"
			emit(out, 2)
		}
		cfg, err := agent.LoadConfig(*configPath)
		if err != nil {
			out["error"] = err.Error()
			emit(out, 2)
		}
		if errs := agent.ValidateConfig(cfg); len(errs) > 0 {
			out["validationErrors"] = errs
			out["error"] = "invalid config"
			emit(out, 2)
		}
		var cell *crdv1alpha1.Cell
		if *cellInputPath != "" {
			raw, err := os.ReadFile(*cellInputPath)
			if err != nil {
				out["error"] = err.Error()
				emit(out, 2)
			}
			cell = &crdv1alpha1.Cell{}
			if err := json.Unmarshal(raw, cell); err != nil {
				out["error"] = err.Error()
				emit(out, 2)
			}
		}
		cliOut, err := agent.RunReadTelemetryCLI(cfg, *telemetryCgroupPath, *allowPathPrefix, cell, *telemetryOutputPath)
		if err != nil {
			out["error"] = err.Error()
			emit(out, 2)
		}
		b, err := json.MarshalIndent(cliOut, "", "  ")
		if err != nil {
			out["error"] = err.Error()
			emit(out, 2)
		}
		os.Stdout.Write(b)
		os.Stdout.Write([]byte("\n"))
		os.Exit(0)

	case "collect-evidence":
		if *configPath == "" || *cellInputPath == "" {
			out["error"] = "collect-evidence requires -config and -cell-input"
			emit(out, 2)
		}
		cfg, err := agent.LoadConfig(*configPath)
		if err != nil {
			out["error"] = err.Error()
			emit(out, 2)
		}
		if errs := agent.ValidateConfig(cfg); len(errs) > 0 {
			out["validationErrors"] = errs
			out["error"] = "invalid config"
			emit(out, 2)
		}
		cellRaw, err := os.ReadFile(*cellInputPath)
		if err != nil {
			out["error"] = err.Error()
			emit(out, 2)
		}
		cell := &crdv1alpha1.Cell{}
		if err := json.Unmarshal(cellRaw, cell); err != nil {
			out["error"] = err.Error()
			emit(out, 2)
		}
		var leaseRaw, portRaw []byte
		if *leasePath != "" {
			leaseRaw, err = os.ReadFile(*leasePath)
			if err != nil {
				out["error"] = err.Error()
				emit(out, 2)
			}
		}
		if *portPath != "" {
			portRaw, err = os.ReadFile(*portPath)
			if err != nil {
				out["error"] = err.Error()
				emit(out, 2)
			}
		}
		var ctx *resourcelease.CellContext
		if *cellCtxPath != "" {
			raw, err := os.ReadFile(*cellCtxPath)
			if err != nil {
				out["error"] = err.Error()
				emit(out, 2)
			}
			ctx = &resourcelease.CellContext{}
			if err := json.Unmarshal(raw, ctx); err != nil {
				out["error"] = err.Error()
				emit(out, 2)
			}
		}
		var integ *agent.CollectEvidenceIntegrityOpts
		man := strings.TrimSpace(*evidenceManifestOut)
		dig := strings.TrimSpace(*evidenceDigestOut)
		sign := strings.TrimSpace(*signingMode)
		if man != "" || dig != "" || (sign != "" && sign != "none") {
			integ = &agent.CollectEvidenceIntegrityOpts{
				ManifestOutputPath: *evidenceManifestOut,
				DigestOutputPath:   *evidenceDigestOut,
				SigningMode:        *signingMode,
				SigningKeyFile:     *signingKeyFile,
				ArtifactID:         *artifactID,
			}
		}
		cliOut, err := agent.RunCollectEvidenceCLI(cfg, cell, *cgroupRoot, *allowPathPrefix, leaseRaw, portRaw, ctx, *evidenceOutputPath, *allowUnsafe, *cpuDelta, *memDelta, integ)
		if err != nil {
			out["error"] = err.Error()
			emit(out, 2)
		}
		b, err := json.MarshalIndent(cliOut, "", "  ")
		if err != nil {
			out["error"] = err.Error()
			emit(out, 2)
		}
		os.Stdout.Write(b)
		os.Stdout.Write([]byte("\n"))
		os.Exit(0)

	default:
		out["error"] = fmt.Sprintf("unknown mode %q", *mode)
		emit(out, 2)
	}
}

func emit(v map[string]interface{}, code int) {
	b, _ := json.MarshalIndent(v, "", "  ")
	os.Stdout.Write(b)
	os.Stdout.Write([]byte("\n"))
	os.Exit(code)
}

func boolExit(hasErr bool) int {
	if hasErr {
		return 2
	}
	return 0
}
