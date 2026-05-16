package host

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/flightrecorder"
)

func TestDetectStaleHeartbeat(t *testing.T) {
	old := time.Now().UTC().Add(-5 * time.Minute).Format(time.RFC3339)
	if !DetectStaleHeartbeat(old, time.Minute, time.Now().UTC()) {
		t.Fatal("expected stale")
	}
	fresh := time.Now().UTC().Format(time.RFC3339)
	if DetectStaleHeartbeat(fresh, time.Minute, time.Now().UTC()) {
		t.Fatal("expected fresh")
	}
}

func TestRunHostHeartbeatNoMutation(t *testing.T) {
	ResetRuntimeSession()
	flightrecorder.Reset()
	res, err := RunHostHeartbeat(HeartbeatOptions{
		Config:     testHostConfig(),
		NodeName:   "karl-metal-01",
		Namespace:  "khr-runtime-sandbox",
		Iterations: 2,
		Interval:   time.Millisecond,
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.Blocked || !res.NoMutation || !res.NoProductionMutation {
		t.Fatalf("res=%+v", res)
	}
	if len(res.Iterations) != 2 {
		t.Fatalf("iterations=%d", len(res.Iterations))
	}
	sid := res.Iterations[0].RuntimeSession.RuntimeSessionID
	if sid == "" || sid != res.Iterations[1].RuntimeSession.RuntimeSessionID {
		t.Fatal("runtime session id must be stable across ticks")
	}
}

func TestRunHostHeartbeatStaleSimulation(t *testing.T) {
	ResetRuntimeSession()
	res, err := RunHostHeartbeat(HeartbeatOptions{
		Config:           testHostConfig(),
		Namespace:        "khr-runtime-sandbox",
		PriorHeartbeatAt: time.Now().UTC().Add(-10 * time.Minute).Format(time.RFC3339),
		StaleThreshold:   time.Minute,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !res.StaleDetected || !res.Iterations[0].Stale {
		t.Fatalf("res=%+v", res)
	}
}

func TestDiscoverLastApplyStateFromSandbox(t *testing.T) {
	dir := t.TempDir()
	state := discoverLastApplyState(dir)
	if state != ApplyStateIdle {
		t.Fatalf("state=%q", state)
	}
	if err := osWriteBaseline(dir); err != nil {
		t.Fatal(err)
	}
	state = discoverLastApplyState(dir)
	if state != ApplyStateApplied {
		t.Fatalf("state=%q", state)
	}
}

func osWriteBaseline(dir string) error {
	path := filepath.Join(dir, "baseline-test.json")
	return os.WriteFile(path, []byte(`{"cpuMaxApplied":"25000 100000"}`), 0o644)
}
