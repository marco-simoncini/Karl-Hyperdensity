package windowsfluidvirt

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestValidateFakeRuntimeBoundary(t *testing.T) {
	boundary := DefaultWindowsFluidVirtFakeRuntimeBoundary()
	if err := ValidateFakeRuntimeBoundary(boundary); err != nil {
		t.Fatalf("expected valid fake runtime boundary, got error: %v", err)
	}
}

func TestReplayFixtureFromTemporaryFile(t *testing.T) {
	boundary := DefaultWindowsFluidVirtFakeRuntimeBoundary()
	replay := NewWindowsFluidVirtNodeActuatorReadonlyReplayMinimal()
	dir := t.TempDir()
	path := filepath.Join(dir, "node_actuator_replay.json")
	raw, err := json.Marshal(replay)
	if err != nil {
		t.Fatalf("marshal replay: %v", err)
	}
	if err := os.WriteFile(path, raw, 0o600); err != nil {
		t.Fatalf("write temp replay fixture: %v", err)
	}

	loaded, err := ReplayFixtureFromTemporaryFile(path, boundary)
	if err != nil {
		t.Fatalf("replay fixture from temp file failed: %v", err)
	}
	if loaded.RuntimeMutationEnabled || loaded.ActuatorRuntimeEnabled || loaded.CgroupWriteEnabled {
		t.Fatalf("loaded replay must remain readonly")
	}
}

func TestReplayFixtureRejectsRealCgroupPath(t *testing.T) {
	boundary := DefaultWindowsFluidVirtFakeRuntimeBoundary()
	if _, err := ReplayFixtureFromTemporaryFile("/sys/fs/cgroup/cpu.max", boundary); err == nil {
		t.Fatalf("expected rejection for real cgroup path")
	}
}
