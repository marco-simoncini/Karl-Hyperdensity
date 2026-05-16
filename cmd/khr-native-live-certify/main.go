package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/nativelive"
)

func main() {
	sprint := flag.String("sprint", "KHR-T", "sprint label")
	baselinePath := flag.String("baseline", "", "baseline certification JSON")
	outPath := flag.String("out", "", "write certification-summary.json")
	requireMatch := flag.Bool("require-baseline-match", false, "fail when baseline compare does not match")
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: khr-native-live-certify [-baseline file] [-out file] run-metrics.json ...")
		os.Exit(2)
	}
	runs := make([]nativelive.RunMetrics, 0, len(args))
	for _, path := range args {
		b, err := os.ReadFile(path)
		if err != nil {
			fatal(err)
		}
		var m nativelive.RunMetrics
		if err := json.Unmarshal(b, &m); err != nil {
			fatal(err)
		}
		nativelive.NormalizeRunMetrics(&m)
		runs = append(runs, m)
	}
	summary := nativelive.AggregateRuns(*sprint, runs)
	if *baselinePath != "" {
		b, err := os.ReadFile(*baselinePath)
		if err != nil {
			fatal(err)
		}
		var baseline nativelive.BaselineCertification
		if err := json.Unmarshal(b, &baseline); err != nil {
			fatal(err)
		}
		match, diffs := nativelive.CompareBaseline(summary, baseline)
		summary.BaselineMatch = match
		summary.BaselineDiff = diffs
		if *requireMatch && !match {
			fatal(fmt.Errorf("baseline mismatch: %v", diffs))
		}
	}
	if err := nativelive.CheckRegression(summary); err != nil {
		if *outPath != "" {
			_ = writeJSON(*outPath, summary)
		}
		fatal(err)
	}
	if len(runs) > 1 && !nativelive.FingerprintsMatch(summary) {
		summary.RegressionDetected = true
		summary.RegressionReasons = append(summary.RegressionReasons, "run fingerprints differ across repeatable runs")
		summary.Status = nativelive.CertificationFailed
		if *outPath != "" {
			_ = writeJSON(*outPath, summary)
		}
		fatal(fmt.Errorf("non-deterministic run fingerprints"))
	}
	if *outPath == "" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(summary)
		return
	}
	if err := writeJSON(*outPath, summary); err != nil {
		fatal(err)
	}
}

func writeJSON(path string, v any) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	b = append(b, '\n')
	return os.WriteFile(path, b, 0o644)
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
