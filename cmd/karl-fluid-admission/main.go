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
	var fixturePath string
	var bundlePath string
	var policyPath string
	var requestedAction string
	var evaluationTimeArg string

	flag.StringVar(&fixturePath, "fixture", "", "Path to admission replay fixture JSON")
	flag.StringVar(&bundlePath, "bundle", "", "Path to runtime evidence bundle JSON")
	flag.StringVar(&policyPath, "policy", "", "Optional path to policy pack JSON")
	flag.StringVar(&requestedAction, "requested-action", "", "Optional requested action override")
	flag.StringVar(&evaluationTimeArg, "evaluation-time", "", "Optional RFC3339 evaluation timestamp")
	flag.Parse()

	if fixturePath == "" && bundlePath == "" {
		fmt.Fprintln(os.Stderr, "provide either -fixture or -bundle")
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

	var result windowsfluidvirt.AdmissionEvaluationResult
	if fixturePath != "" {
		fixture, err := windowsfluidvirt.LoadAdmissionReplayFixture(fixturePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed loading fixture: %v\n", err)
			os.Exit(1)
		}
		action := fixture.RequestedAction
		if requestedAction != "" {
			action = windowsfluidvirt.RequestedAdmissionAction(requestedAction)
		}
		result = windowsfluidvirt.EvaluateWindowsFluidAdmission(windowsfluidvirt.AdmissionEvaluationInput{
			Bundle:          fixture.Bundle,
			PolicyPack:      fixture.PolicyPack,
			RequestedAction: action,
			EvaluationTime:  evalTime,
		})
	} else {
		bundleData, err := os.ReadFile(bundlePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed reading bundle: %v\n", err)
			os.Exit(1)
		}
		var bundle windowsfluidvirt.WindowsFluidRuntimeEvidenceBundle
		if err := json.Unmarshal(bundleData, &bundle); err != nil {
			fmt.Fprintf(os.Stderr, "failed decoding bundle: %v\n", err)
			os.Exit(1)
		}

		var policy *windowsfluidvirt.WindowsFluidPolicyPack
		if policyPath != "" {
			policyData, err := os.ReadFile(policyPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed reading policy: %v\n", err)
				os.Exit(1)
			}
			var parsed windowsfluidvirt.WindowsFluidPolicyPack
			if err := json.Unmarshal(policyData, &parsed); err != nil {
				fmt.Fprintf(os.Stderr, "failed decoding policy: %v\n", err)
				os.Exit(1)
			}
			policy = &parsed
		}

		result = windowsfluidvirt.EvaluateWindowsFluidAdmission(windowsfluidvirt.AdmissionEvaluationInput{
			Bundle:          bundle,
			PolicyPack:      policy,
			RequestedAction: windowsfluidvirt.RequestedAdmissionAction(requestedAction),
			EvaluationTime:  evalTime,
		})
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(result); err != nil {
		fmt.Fprintf(os.Stderr, "failed encoding output: %v\n", err)
		os.Exit(1)
	}
}
