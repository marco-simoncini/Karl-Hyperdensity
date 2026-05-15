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
	var evaluationTimeArg string
	flag.StringVar(&fixturePath, "fixture", "", "Path to dry-run replay fixture JSON")
	flag.StringVar(&bundlePath, "bundle", "", "Path to runtime evidence bundle JSON")
	flag.StringVar(&evaluationTimeArg, "evaluation-time", "", "Optional RFC3339 evaluation timestamp")
	flag.Parse()

	if fixturePath == "" && bundlePath == "" {
		fmt.Fprintln(os.Stderr, "provide either -fixture or -bundle")
		os.Exit(2)
	}

	var options windowsfluidvirt.DryRunEvaluationOptions
	if evaluationTimeArg != "" {
		parsed, err := time.Parse(time.RFC3339, evaluationTimeArg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid -evaluation-time value: %v\n", err)
			os.Exit(2)
		}
		options.EvaluationTime = parsed
	}

	var result windowsfluidvirt.DryRunEvaluationResult
	if fixturePath != "" {
		fixture, err := windowsfluidvirt.LoadDryRunReplayFixture(fixturePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed replay evaluation: %v\n", err)
			os.Exit(1)
		}
		result = windowsfluidvirt.EvaluateWindowsFluidRuntimeDryRunWithOptions(fixture.Bundle, options)
	} else {
		data, err := os.ReadFile(bundlePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed reading bundle: %v\n", err)
			os.Exit(1)
		}
		var bundle windowsfluidvirt.WindowsFluidRuntimeEvidenceBundle
		if err := json.Unmarshal(data, &bundle); err != nil {
			fmt.Fprintf(os.Stderr, "failed decoding bundle: %v\n", err)
			os.Exit(1)
		}
		result = windowsfluidvirt.EvaluateWindowsFluidRuntimeDryRunWithOptions(bundle, options)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(result); err != nil {
		fmt.Fprintf(os.Stderr, "failed encoding result: %v\n", err)
		os.Exit(1)
	}
}
