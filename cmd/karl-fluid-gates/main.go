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
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "karl-fluid-gates (non-executable)\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Executor remains disabled. No runtime mutation. No QMP commands are sent. No CPU/RAM apply is performed.\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Usage:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  karl-fluid-gates -fixture <path> [-mode gate|gateset] [-evaluation-time <RFC3339>]\n")
		flag.PrintDefaults()
	}

	var fixturePath string
	var mode string
	var evaluationTimeArg string

	flag.StringVar(&fixturePath, "fixture", "", "Path to unlock gate fixture JSON")
	flag.StringVar(&mode, "mode", "auto", "Evaluation mode: auto, gate, or gateset")
	flag.StringVar(&evaluationTimeArg, "evaluation-time", "", "Optional RFC3339 evaluation timestamp")
	flag.Parse()

	if fixturePath == "" {
		fmt.Fprintln(os.Stderr, "provide -fixture")
		os.Exit(2)
	}

	var evalTime time.Time
	if evaluationTimeArg != "" {
		parsed, err := time.Parse(time.RFC3339, evaluationTimeArg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid -evaluation-time value: %v\n", err)
			os.Exit(2)
		}
		evalTime = parsed
	}

	fixture, err := windowsfluidvirt.LoadUnlockGateReplayFixture(fixturePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed loading fixture: %v\n", err)
		os.Exit(1)
	}

	selectedMode := mode
	if selectedMode == "auto" {
		if fixture.GateID == "" {
			selectedMode = "gateset"
		} else {
			selectedMode = "gate"
		}
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	switch selectedMode {
	case "gate":
		result := windowsfluidvirt.EvaluateWindowsFluidUnlockGate(windowsfluidvirt.UnlockGateEvaluationInput{
			GateID:             fixture.GateID,
			EvidenceBundle:     fixture.EvidenceBundle,
			GovernanceContract: fixture.GovernanceContract,
			ExecutorOutput:     fixture.ExecutorOutput,
			Attestation:        fixture.Attestation,
			ParityEvidence:     fixture.ParityEvidence,
			EvaluationTime:     evalTime,
		})
		if err := enc.Encode(result); err != nil {
			fmt.Fprintf(os.Stderr, "failed encoding output: %v\n", err)
			os.Exit(1)
		}
	case "gateset":
		result := windowsfluidvirt.EvaluateWindowsFluidUnlockGateSet(windowsfluidvirt.UnlockGateSetEvaluationInput{
			EvidenceBundle:     fixture.EvidenceBundle,
			GovernanceContract: fixture.GovernanceContract,
			ExecutorOutput:     fixture.ExecutorOutput,
			Attestation:        fixture.Attestation,
			ParityEvidence:     fixture.ParityEvidence,
			EvaluationTime:     evalTime,
		})
		if err := enc.Encode(result); err != nil {
			fmt.Fprintf(os.Stderr, "failed encoding output: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "unsupported -mode %q (use auto|gate|gateset)\n", selectedMode)
		os.Exit(2)
	}
}
