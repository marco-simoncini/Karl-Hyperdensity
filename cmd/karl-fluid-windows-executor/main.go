package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/windowsfluidvirt"
)

type executorOutput struct {
	ExecutorVersion         string                                          `json:"executorVersion"`
	Mode                    string                                          `json:"mode"`
	RuntimeMutationExecuted bool                                            `json:"runtimeMutationExecuted"`
	Plan                    windowsfluidvirt.WindowsFluidControlledApplyPlan `json:"plan"`
	Evidence                map[string]any                                  `json:"evidence"`
}

func main() {
	var targetPath string
	var leasePath string
	var gatePath string
	var approvalPath string
	var mode string
	var outPath string
	var evaluationTimeRaw string
	var pretty bool

	flag.StringVar(&targetPath, "target", "", "Path to WindowsHyperdensityTarget JSON")
	flag.StringVar(&leasePath, "lease", "", "Path to WindowsFluidResourceLease JSON")
	flag.StringVar(&gatePath, "gate", "", "Path to WindowsFluidControlledApplyGate JSON")
	flag.StringVar(&approvalPath, "approval", "", "Optional path to WindowsFluidManualApproval JSON")
	flag.StringVar(&mode, "mode", "plan", "Execution mode: plan|dry-run|apply-plan-only")
	flag.StringVar(&outPath, "out", "", "Optional path to write output JSON")
	flag.StringVar(&evaluationTimeRaw, "evaluation-time", "", "Optional RFC3339 evaluation timestamp")
	flag.BoolVar(&pretty, "pretty", false, "Pretty-print JSON output")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "karl-fluid-windows-executor (controlled apply planning only)\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), "This milestone is plan/dry-run only:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "- no autonomous apply\n")
		fmt.Fprintf(flag.CommandLine.Output(), "- no runtime mutation in this CLI milestone\n")
		fmt.Fprintf(flag.CommandLine.Output(), "- no cluster calls, no QMP calls, no cgroup writes\n")
		fmt.Fprintf(flag.CommandLine.Output(), "- no vCPU hotplug, no logical CPU scaling\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if targetPath == "" || leasePath == "" || gatePath == "" {
		fmt.Fprintln(os.Stderr, "target, lease, and gate are required")
		os.Exit(2)
	}
	if mode != "plan" && mode != "dry-run" && mode != "apply-plan-only" {
		fmt.Fprintf(os.Stderr, "unsupported mode %q\n", mode)
		os.Exit(2)
	}

	var evaluationTime time.Time
	if evaluationTimeRaw != "" {
		parsed, err := time.Parse(time.RFC3339, evaluationTimeRaw)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid -evaluation-time value: %v\n", err)
			os.Exit(2)
		}
		evaluationTime = parsed
	}

	target, err := loadJSON[windowsfluidvirt.WindowsHyperdensityTarget](targetPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load target: %v\n", err)
		os.Exit(1)
	}
	lease, err := loadJSON[windowsfluidvirt.WindowsFluidResourceLease](leasePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load lease: %v\n", err)
		os.Exit(1)
	}
	gate, err := loadJSON[windowsfluidvirt.WindowsFluidControlledApplyGate](gatePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load gate: %v\n", err)
		os.Exit(1)
	}
	var approval *windowsfluidvirt.WindowsFluidManualApproval
	if approvalPath != "" {
		parsed, parseErr := loadJSON[windowsfluidvirt.WindowsFluidManualApproval](approvalPath)
		if parseErr != nil {
			fmt.Fprintf(os.Stderr, "load approval: %v\n", parseErr)
			os.Exit(1)
		}
		approval = &parsed
	}

	plan := windowsfluidvirt.BuildWindowsFluidControlledApplyPlan(target, lease, gate, approval, evaluationTime)
	plan.ApplyAllowed = plan.ApplyAllowed && mode == "apply-plan-only"
	if mode != "apply-plan-only" {
		plan.MutationAllowed = false
	}

	output := executorOutput{
		ExecutorVersion:         "karl-fluid-windows-executor-v1",
		Mode:                    mode,
		RuntimeMutationExecuted: false,
		Plan:                    plan,
		Evidence: map[string]any{
			"planningOnly":           true,
			"autonomousApplyEnabled": false,
			"runtimeMutationAllowed": false,
		},
	}

	var encoded []byte
	if pretty {
		encoded, err = json.MarshalIndent(output, "", "  ")
	} else {
		encoded, err = json.Marshal(output)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "encode output: %v\n", err)
		os.Exit(1)
	}

	if outPath != "" {
		if writeErr := os.WriteFile(outPath, append(encoded, '\n'), 0o644); writeErr != nil {
			fmt.Fprintf(os.Stderr, "write output: %v\n", writeErr)
			os.Exit(1)
		}
	}
	_, _ = os.Stdout.Write(append(encoded, '\n'))
}

func loadJSON[T any](path string) (T, error) {
	var out T
	data, err := os.ReadFile(path)
	if err != nil {
		return out, err
	}
	if err := json.Unmarshal(data, &out); err != nil {
		return out, err
	}
	return out, nil
}
