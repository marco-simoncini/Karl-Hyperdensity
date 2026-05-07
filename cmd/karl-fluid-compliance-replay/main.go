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

	flag.StringVar(&inputPath, "input", "", "Path to replay input JSON fixture or direct input payload.")
	flag.StringVar(&evaluationTimeRaw, "evaluation-time", "", "Deterministic replay time in RFC3339 format (UTC recommended).")
	flag.BoolVar(&emitAttestation, "emit-attestation", false, "Emit future-signable attestation envelope (no real signature).")
	flag.BoolVar(&pretty, "pretty", false, "Print JSON output in pretty format.")
	flag.StringVar(&attestationModeRaw, "attestation-mode", string(windowsfluidvirt.AttestationModeFutureSignable), "Attestation signature mode: unsigned-dev or future-signable.")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "karl-fluid-compliance-replay (read-only)\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), "This CLI performs read-only compliance replay and never mutates runtime.\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Safety guarantees:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "- no runtime mutation\n")
		fmt.Fprintf(flag.CommandLine.Output(), "- no CPU/RAM apply\n")
		fmt.Fprintf(flag.CommandLine.Output(), "- no actuator apply\n")
		fmt.Fprintf(flag.CommandLine.Output(), "- no cluster calls\n")
		fmt.Fprintf(flag.CommandLine.Output(), "- future-signable attestation only (no real signature)\n\n")
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

	encoder := json.NewEncoder(os.Stdout)
	if pretty {
		encoder.SetIndent("", "  ")
	}
	if err := encoder.Encode(output); err != nil {
		fatalf("encode output: %v", err)
	}
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
