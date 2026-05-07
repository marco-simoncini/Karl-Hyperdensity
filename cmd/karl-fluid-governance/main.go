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
	var admissionPath string
	var bundlePath string
	var policyPath string
	var requestedAction string
	var evaluationTimeArg string

	flag.StringVar(&fixturePath, "fixture", "", "Path to governance replay fixture JSON")
	flag.StringVar(&admissionPath, "admission", "", "Path to admission decision JSON for bundle mode")
	flag.StringVar(&bundlePath, "bundle", "", "Path to runtime evidence bundle JSON")
	flag.StringVar(&policyPath, "policy", "", "Optional path to policy pack JSON")
	flag.StringVar(&requestedAction, "requested-action", "", "Optional requested governance action override")
	flag.StringVar(&evaluationTimeArg, "evaluation-time", "", "Optional RFC3339 evaluation timestamp")
	flag.Parse()

	if fixturePath == "" && (admissionPath == "" || bundlePath == "") {
		fmt.Fprintln(os.Stderr, "provide -fixture or (-admission and -bundle)")
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

	var result windowsfluidvirt.ApplyGovernanceEvaluationResult
	if fixturePath != "" {
		fixture, err := windowsfluidvirt.LoadGovernanceReplayFixture(fixturePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed loading governance fixture: %v\n", err)
			os.Exit(1)
		}
		action := fixture.RequestedAction
		if requestedAction != "" {
			action = windowsfluidvirt.GovernanceRequestedAction(requestedAction)
		}
		result = windowsfluidvirt.EvaluateWindowsFluidApplyGovernance(windowsfluidvirt.ApplyGovernanceEvaluationInput{
			AdmissionDecision: fixture.AdmissionDecision,
			Bundle:            fixture.Bundle,
			PolicyPack:        fixture.PolicyPack,
			RequestedAction:   action,
			EvaluationTime:    evalTime,
		})
	} else {
		admissionData, err := os.ReadFile(admissionPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed reading admission file: %v\n", err)
			os.Exit(1)
		}
		var admission windowsfluidvirt.WindowsFluidAdmissionDecision
		if err := json.Unmarshal(admissionData, &admission); err != nil {
			fmt.Fprintf(os.Stderr, "failed decoding admission file: %v\n", err)
			os.Exit(1)
		}

		bundleData, err := os.ReadFile(bundlePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed reading bundle file: %v\n", err)
			os.Exit(1)
		}
		var bundle windowsfluidvirt.WindowsFluidRuntimeEvidenceBundle
		if err := json.Unmarshal(bundleData, &bundle); err != nil {
			fmt.Fprintf(os.Stderr, "failed decoding bundle file: %v\n", err)
			os.Exit(1)
		}

		var policy *windowsfluidvirt.WindowsFluidPolicyPack
		if policyPath != "" {
			policyData, err := os.ReadFile(policyPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed reading policy file: %v\n", err)
				os.Exit(1)
			}
			var parsed windowsfluidvirt.WindowsFluidPolicyPack
			if err := json.Unmarshal(policyData, &parsed); err != nil {
				fmt.Fprintf(os.Stderr, "failed decoding policy file: %v\n", err)
				os.Exit(1)
			}
			policy = &parsed
		}

		result = windowsfluidvirt.EvaluateWindowsFluidApplyGovernance(windowsfluidvirt.ApplyGovernanceEvaluationInput{
			AdmissionDecision: admission,
			Bundle:            bundle,
			PolicyPack:        policy,
			RequestedAction:   windowsfluidvirt.GovernanceRequestedAction(requestedAction),
			EvaluationTime:    evalTime,
		})
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(result); err != nil {
		fmt.Fprintf(os.Stderr, "failed encoding output: %v\n", err)
		os.Exit(1)
	}
}
