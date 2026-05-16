package resourceport

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/flightrecorder"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/host"
)

func loopConfig(enabled bool) *host.Config {
	cfg := &host.Config{}
	cfg.Spec.HostID = "host-loop-test"
	cfg.Spec.LinuxOnly = true
	cfg.Spec.SandboxMode = true
	cfg.Spec.ResourcePortLoopEnabled = enabled
	cfg.Spec.AllowedNamespaces = []string{"khr-runtime-sandbox"}
	cfg.Spec.AllowedLabels = map[string]string{"khr.karl.io/sandbox": "true"}
	return cfg
}

func baseLoopOpts(cfg *host.Config) LoopOptions {
	return LoopOptions{
		Config:    cfg,
		Namespace: "khr-runtime-sandbox",
		Labels:    map[string]string{"khr.karl.io/sandbox": "true"},
		Targets: []SandboxTarget{{
			Namespace: "khr-runtime-sandbox",
			PodName:   "target-1",
			Labels:    map[string]string{"khr.karl.io/sandbox": "true"},
		}},
		Iterations: 1,
	}
}

func TestRunLoopDisabledByDefault(t *testing.T) {
	res, err := RunLoop(baseLoopOpts(loopConfig(false)))
	if err != nil {
		t.Fatal(err)
	}
	if !res.Blocked {
		t.Fatalf("res=%+v", res)
	}
}

func TestRunLoopJSONEmission(t *testing.T) {
	res, err := RunLoop(baseLoopOpts(loopConfig(true)))
	if err != nil {
		t.Fatal(err)
	}
	if res.Blocked || res.EmissionMode != EmissionModeJSONOnly {
		t.Fatalf("res=%+v", res)
	}
	if len(res.Iterations) != 1 || len(res.Iterations[0].ResourcePorts) == 0 {
		t.Fatal("expected resource port emission")
	}
}

func TestRunLoopBlockedNamespace(t *testing.T) {
	opts := baseLoopOpts(loopConfig(true))
	opts.Namespace = "karl-system"
	res, err := RunLoop(opts)
	if err != nil {
		t.Fatal(err)
	}
	if !res.Blocked {
		t.Fatal("production namespace must be blocked")
	}
}

func TestRunLoopLabelAllowlist(t *testing.T) {
	opts := baseLoopOpts(loopConfig(true))
	opts.Labels = map[string]string{"khr.karl.io/sandbox": "false"}
	res, err := RunLoop(opts)
	if err != nil {
		t.Fatal(err)
	}
	if !res.Blocked {
		t.Fatal("label mismatch must block")
	}
}

func TestRunLoopNoProductionMutation(t *testing.T) {
	dir := t.TempDir()
	opts := baseLoopOpts(loopConfig(true))
	opts.EmitCR = false
	opts.OutputDir = dir
	res, err := RunLoop(opts)
	if err != nil {
		t.Fatal(err)
	}
	entries, _ := os.ReadDir(dir)
	if len(entries) > 0 {
		t.Fatal("emit-cr false must not write preview files")
	}
	if !res.NoProductionMutation {
		t.Fatal("expected no production mutation flag")
	}
}

func TestRunLoopFlightRecorderAppend(t *testing.T) {
	flightrecorder.Reset()
	res, err := RunLoop(baseLoopOpts(loopConfig(true)))
	if err != nil {
		t.Fatal(err)
	}
	if len(res.FlightRecorder) < 2 {
		t.Fatalf("events=%d", len(res.FlightRecorder))
	}
}

func TestRunLoopCRPreviewLocalOnly(t *testing.T) {
	dir := t.TempDir()
	opts := baseLoopOpts(loopConfig(true))
	opts.EmitCR = true
	opts.OutputDir = dir
	res, err := RunLoop(opts)
	if err != nil {
		t.Fatal(err)
	}
	if !res.EmitCRApplied || res.EmissionMode != EmissionModeCRPreview {
		t.Fatalf("res=%+v", res)
	}
	if _, err := os.Stat(filepath.Join(dir, "resourceport-khr-runtime-sandbox-target-1-port.json")); err != nil {
		t.Fatal(err)
	}
}

func TestRunLoopClusterContextGuard(t *testing.T) {
	opts := baseLoopOpts(loopConfig(true))
	opts.ClusterContext = "wrong-context"
	opts.RequiredContext = "karl-metal-01@ovh"
	opts.Targets = nil
	res, err := RunLoop(opts)
	if err != nil {
		t.Fatal(err)
	}
	if !res.Blocked {
		t.Fatal("context mismatch must block when discovery required")
	}
}

func TestRunLoopMultipleIterations(t *testing.T) {
	opts := baseLoopOpts(loopConfig(true))
	opts.Iterations = 2
	opts.Interval = time.Millisecond
	res, err := RunLoop(opts)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Iterations) != 2 {
		t.Fatalf("iterations=%d", len(res.Iterations))
	}
}
