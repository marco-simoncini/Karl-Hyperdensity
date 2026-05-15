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
	var inputPath string
	var evaluationTimeRaw string
	var emitAttestation bool
	var pretty bool
	var attestationModeRaw string
	var bundleIndexVersion string
	var bundleSubject string
	var previousRunHash string
	var emitBundleIndex bool
	var appendBundleIn string
	var appendBundle bool
	var bundleOut string

	flag.StringVar(&inputPath, "input", "", "Path to replay input JSON fixture or direct input payload.")
	flag.StringVar(&evaluationTimeRaw, "evaluation-time", "", "Deterministic replay time in RFC3339 format (UTC recommended).")
	flag.BoolVar(&emitAttestation, "emit-attestation", false, "Emit future-signable attestation envelope (no real signature).")
	flag.BoolVar(&pretty, "pretty", false, "Print JSON output in pretty format.")
	flag.StringVar(&attestationModeRaw, "attestation-mode", string(windowsfluidvirt.AttestationModeFutureSignable), "Attestation signature mode: unsigned-dev or future-signable.")
	flag.StringVar(&bundleIndexVersion, "bundle-index", "windows-fluid-compliance-replay-bundle-index-v1", "Bundle index version label (used with -emit-bundle-index).")
	flag.StringVar(&bundleSubject, "bundle-subject", "", "Bundle subject reference. Defaults to replay shellRef.")
	flag.StringVar(&previousRunHash, "previous-run-hash", "", "Optional previous run hash for single-run chain linkage.")
	flag.BoolVar(&emitBundleIndex, "emit-bundle-index", false, "Emit single-run replay bundle index (read-only, local hash chain).")
	flag.StringVar(&appendBundleIn, "append-bundle-in", "", "Path to existing bundle index JSON used for append workflow.")
	flag.BoolVar(&appendBundle, "append-bundle", false, "Append run to existing bundle provided by -append-bundle-in.")
	flag.StringVar(&bundleOut, "bundle-out", "", "Optional output path for writing resulting bundle index JSON.")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "karl-fluid-compliance-replay (read-only)\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), "This CLI performs read-only compliance replay and never mutates runtime.\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Safety guarantees:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "- no runtime mutation\n")
		fmt.Fprintf(flag.CommandLine.Output(), "- no CPU/RAM apply\n")
		fmt.Fprintf(flag.CommandLine.Output(), "- no actuator apply\n")
		fmt.Fprintf(flag.CommandLine.Output(), "- no cluster calls\n")
		fmt.Fprintf(flag.CommandLine.Output(), "- future-signable attestation only (no real signature)\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Bundle index notes:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "- emits local deterministic hash chain metadata\n")
		fmt.Fprintf(flag.CommandLine.Output(), "- no keys, no KMS, no real signatures\n")
		fmt.Fprintf(flag.CommandLine.Output(), "- supports single-run emission and append to existing bundle\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Flags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if inputPath == "" {
		fatalf("missing required -input flag")
	}
	evaluationTime, err := parseEvaluationTime(evaluationTimeRaw)
	if err != nil {
		fatalf("invalid -evaluation-time: %v", err)
	}
	input, err := windowsfluidvirt.LoadComplianceReplayInput(inputPath)
	if err != nil {
		fatalf("load replay input: %v", err)
	}
	replay, err := windowsfluidvirt.EvaluateWindowsComplianceReplay(input, inputPath, evaluationTime)
	if err != nil {
		fatalf("evaluate replay: %v", err)
	}
	output := windowsfluidvirt.WindowsComplianceReplayCLIOutput{
		WindowsComplianceReplayOutput: replay,
	}
	if emitAttestation {
		mode, parseErr := windowsfluidvirt.ParseAttestationMode(attestationModeRaw)
		if parseErr != nil {
			fatalf("invalid attestation mode: %v", parseErr)
		}
		attestation, attErr := windowsfluidvirt.BuildWindowsComplianceReplayAttestation(replay, mode, evaluationTime)
		if attErr != nil {
			fatalf("build attestation: %v", attErr)
		}
		output.Attestation = &attestation
	}
	if appendBundle {
		if appendBundleIn == "" {
			fatalf("append mode requires -append-bundle-in")
		}
		if previousRunHash != "" {
			fatalf("-previous-run-hash is not allowed in append mode")
		}
		subject := bundleSubject
		if subject == "" {
			subject = replay.ShellRef
		}
		run, runErr := windowsfluidvirt.BuildWindowsComplianceReplayBundleRun(replay, output.Attestation, "", evaluationTime)
		if runErr != nil {
			fatalf("build bundle run: %v", runErr)
		}
		existingBundle, loadErr := loadBundleIndex(appendBundleIn)
		if loadErr != nil {
			fatalf("load existing bundle index: %v", loadErr)
		}
		if existingBundle.SubjectRef != subject {
			fatalf("bundle subject mismatch: existing=%s requested=%s", existingBundle.SubjectRef, subject)
		}
		appended, appendErr := windowsfluidvirt.AppendWindowsComplianceReplayBundleRun(existingBundle, run, evaluationTime)
		if appendErr != nil {
			fatalf("append bundle run: %v", appendErr)
		}
		output.BundleIndex = &appended
	} else if emitBundleIndex {
		subject := bundleSubject
		if subject == "" {
			subject = replay.ShellRef
		}
		run, runErr := windowsfluidvirt.BuildWindowsComplianceReplayBundleRun(replay, output.Attestation, previousRunHash, evaluationTime)
		if runErr != nil {
			fatalf("build bundle run: %v", runErr)
		}
		bundle, bundleErr := windowsfluidvirt.BuildWindowsComplianceReplayBundleIndex(subject, bundleIndexVersion, []windowsfluidvirt.WindowsComplianceReplayBundleRun{run}, evaluationTime)
		if bundleErr != nil {
			fatalf("build bundle index: %v", bundleErr)
		}
		output.BundleIndex = &bundle
	}

	encoder := json.NewEncoder(os.Stdout)
	if pretty {
		encoder.SetIndent("", "  ")
	}
	if err := encoder.Encode(output); err != nil {
		fatalf("encode output: %v", err)
	}
	if bundleOut != "" {
		if output.BundleIndex == nil {
			fatalf("-bundle-out requires either -emit-bundle-index or -append-bundle")
		}
		if err := writeBundleIndex(bundleOut, *output.BundleIndex); err != nil {
			fatalf("write bundle index: %v", err)
		}
	}
}

func loadBundleIndex(path string) (windowsfluidvirt.WindowsComplianceReplayBundleIndex, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return windowsfluidvirt.WindowsComplianceReplayBundleIndex{}, err
	}
	var index windowsfluidvirt.WindowsComplianceReplayBundleIndex
	if err := json.Unmarshal(data, &index); err != nil {
		return windowsfluidvirt.WindowsComplianceReplayBundleIndex{}, err
	}
	return index, nil
}

func writeBundleIndex(path string, bundle windowsfluidvirt.WindowsComplianceReplayBundleIndex) error {
	data, err := json.MarshalIndent(bundle, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func parseEvaluationTime(raw string) (time.Time, error) {
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
