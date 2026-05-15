// Command khr-linux-agent is a dry-run-first Linux MVP skeleton (Sprint 5).
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/agent"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/resourcelease"
)

func main() {
	mode := flag.String("mode", "", "one of: validate-config, dry-run, print-capabilities")
	configPath := flag.String("config", "", "path to agent YAML/JSON config")
	leasePath := flag.String("lease-input", "", "path to ResourceLease JSON (dry-run)")
	portPath := flag.String("resource-port-input", "", "path to ResourcePort JSON (dry-run)")
	cellCtxPath := flag.String("cell-context", "", "optional path to CellContext JSON (dry-run)")
	allowUnsafe := flag.Bool("allow-unsafe-apply", false, "DANGEROUS: required for any future real cgroup writes (disabled in Sprint 5)")
	cpuDelta := flag.String("cpu-delta", "", "optional cpu.max delta string for envelope dry-run plan")
	memDelta := flag.String("memory-delta", "", "optional memory.max delta string for envelope dry-run plan")
	flag.Parse()

	out := map[string]interface{}{
		"tool":    "khr-linux-agent",
		"version": "0.0.1-sprint5",
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
		var portRaw []byte
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
		leaseOut, err := agent.DryRunLease(leaseRaw, portRaw, ctx)
		if err != nil {
			out["error"] = err.Error()
			emit(out, 2)
		}
		planOut, err := agent.DryRunEnvelopePlan(*allowUnsafe, *cpuDelta, *memDelta)
		if err != nil {
			out["error"] = err.Error()
			emit(out, 2)
		}
		var leaseObj interface{}
		var planObj interface{}
		_ = json.Unmarshal(leaseOut, &leaseObj)
		_ = json.Unmarshal(planOut, &planObj)
		out["resourceLeaseDryRun"] = leaseObj
		out["cgroupEnvelopePlan"] = planObj
		out["applyLocked"] = !*allowUnsafe
		emit(out, 0)

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
