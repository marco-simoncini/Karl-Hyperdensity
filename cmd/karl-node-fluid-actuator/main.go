package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/windowsfluidvirt"
)

func main() {
	var mode string
	var requestPath string
	var allowlistPath string
	var evidenceOut string
	var dryRun bool
	var killSwitchPath string
	var evaluationTimeRaw string

	flag.StringVar(&mode, "mode", "dry-run", "Execution mode: dry-run|apply|rollback|return-to-floor")
	flag.StringVar(&requestPath, "request", "", "Path to actuator request JSON.")
	flag.StringVar(&allowlistPath, "allowlist", "", "Path to allowlist JSON.")
	flag.StringVar(&evidenceOut, "evidence-out", "", "Path to write evidence JSON.")
	flag.BoolVar(&dryRun, "dry-run", false, "Force dry-run (no write).")
	flag.StringVar(&killSwitchPath, "kill-switch", "", "Path to kill-switch marker file.")
	flag.StringVar(&evaluationTimeRaw, "evaluation-time", "", "Deterministic evaluation time in RFC3339.")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "karl-node-fluid-actuator (MVP, lab-controlled)\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Safety scope:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "- validates request/allowlist identity and cgroup path\n")
		fmt.Fprintf(flag.CommandLine.Output(), "- supports cpu.max only\n")
		fmt.Fprintf(flag.CommandLine.Output(), "- no parent cgroup writes\n")
		fmt.Fprintf(flag.CommandLine.Output(), "- no arbitrary file writes\n")
		fmt.Fprintf(flag.CommandLine.Output(), "- supports dry-run/apply/rollback/return-to-floor\n")
		fmt.Fprintf(flag.CommandLine.Output(), "- no Kubernetes or QMP calls\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if requestPath == "" || allowlistPath == "" {
		fatalf("both -request and -allowlist are required")
	}
	evaluationTime, err := parseTime(evaluationTimeRaw)
	if err != nil {
		fatalf("invalid -evaluation-time: %v", err)
	}
	request, err := windowsfluidvirt.LoadNodeFluidActuatorMVPRequest(requestPath)
	if err != nil {
		fatalf("load request: %v", err)
	}
	allowlist, err := windowsfluidvirt.LoadNodeFluidActuatorMVPAllowlist(allowlistPath)
	if err != nil {
		fatalf("load allowlist: %v", err)
	}
	evidence, err := windowsfluidvirt.EvaluateNodeFluidActuatorMVP(windowsfluidvirt.NodeFluidActuatorMVPInput{
		Action:          windowsfluidvirt.KARLNodeFluidActuatorAction(mode),
		Request:         request,
		Allowlist:       allowlist,
		KillSwitchPath:  killSwitchPath,
		EvaluationTime:  evaluationTime,
		EvidenceOutPath: evidenceOut,
		DryRun:          dryRun,
	})
	if err != nil {
		fatalf("evaluate actuator request: %v", err)
	}
	if evidenceOut != "" {
		if err := writeJSON(evidenceOut, evidence); err != nil {
			fatalf("write evidence: %v", err)
		}
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(evidence); err != nil {
		fatalf("encode evidence: %v", err)
	}
	if evidence.Decision == windowsfluidvirt.NodeActuatorDecisionRejected || evidence.Decision == windowsfluidvirt.NodeActuatorDecisionBlocked {
		os.Exit(2)
	}
}

func writeJSON(path string, value any) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(value)
}

func parseTime(raw string) (time.Time, error) {
	if raw == "" {
		return time.Now().UTC(), nil
	}
	t, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return time.Time{}, err
	}
	return t.UTC(), nil
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
