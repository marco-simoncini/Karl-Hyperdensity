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
		fmt.Fprintf(flag.CommandLine.Output(), "karl-fluid-executor (hard-disabled)\n")
		fmt.Fprintf(flag.CommandLine.Output(), "No runtime mutation. No QMP commands are sent. No CPU/RAM apply is performed.\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Usage:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  karl-fluid-executor -fixture <path> [-evaluation-time <RFC3339>]\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  karl-fluid-executor -governance <path> -revalidation <path> -attestation <path> [-killswitch <path>] [-evaluation-time <RFC3339>]\n\n")
		flag.PrintDefaults()
	}

	var fixturePath string
	var governancePath string
	var revalidationPath string
	var attestationPath string
	var killSwitchPath string
	var evaluationTimeArg string

	flag.StringVar(&fixturePath, "fixture", "", "Path to executor replay fixture JSON")
	flag.StringVar(&governancePath, "governance", "", "Path to governance contract JSON")
	flag.StringVar(&revalidationPath, "revalidation", "", "Path to pre-apply revalidation JSON")
	flag.StringVar(&attestationPath, "attestation", "", "Path to policy attestation JSON")
	flag.StringVar(&killSwitchPath, "killswitch", "", "Optional path to kill switch JSON")
	flag.StringVar(&evaluationTimeArg, "evaluation-time", "", "Optional RFC3339 evaluation timestamp")
	flag.Parse()

	if fixturePath == "" && (governancePath == "" || revalidationPath == "" || attestationPath == "") {
		fmt.Fprintln(os.Stderr, "provide -fixture or (-governance -revalidation -attestation)")
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

	var result windowsfluidvirt.FutureApplyExecutorEvaluationResult
	if fixturePath != "" {
		fixture, err := windowsfluidvirt.LoadExecutorReplayFixture(fixturePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed loading executor fixture: %v\n", err)
			os.Exit(1)
		}
		result = windowsfluidvirt.EvaluateWindowsFluidFutureApplyExecutor(windowsfluidvirt.FutureApplyExecutorEvaluationInput{
			GovernanceContract: fixture.GovernanceContract,
			Revalidation:       fixture.Revalidation,
			Attestation:        fixture.Attestation,
			KillSwitch:         fixture.KillSwitch,
			EvaluationTime:     evalTime,
		})
	} else {
		governance, err := loadJSON[windowsfluidvirt.WindowsFluidApplyGovernanceContract](governancePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed loading governance: %v\n", err)
			os.Exit(1)
		}
		revalidation, err := loadJSON[windowsfluidvirt.WindowsFluidPreApplyRevalidationContract](revalidationPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed loading revalidation: %v\n", err)
			os.Exit(1)
		}
		attestation, err := loadJSON[windowsfluidvirt.WindowsFluidPolicyAttestation](attestationPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed loading attestation: %v\n", err)
			os.Exit(1)
		}

		var killSwitch *windowsfluidvirt.WindowsFluidKillSwitch
		if killSwitchPath != "" {
			parsed, err := loadJSON[windowsfluidvirt.WindowsFluidKillSwitch](killSwitchPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed loading killswitch: %v\n", err)
				os.Exit(1)
			}
			killSwitch = &parsed
		}

		result = windowsfluidvirt.EvaluateWindowsFluidFutureApplyExecutor(windowsfluidvirt.FutureApplyExecutorEvaluationInput{
			GovernanceContract: governance,
			Revalidation:       revalidation,
			Attestation:        attestation,
			KillSwitch:         killSwitch,
			EvaluationTime:     evalTime,
		})
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(result); err != nil {
		fmt.Fprintf(os.Stderr, "failed encoding output: %v\n", err)
		os.Exit(1)
	}
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
