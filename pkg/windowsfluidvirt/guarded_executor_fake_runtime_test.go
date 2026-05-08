package windowsfluidvirt

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadGuardedExecutorFakeRuntimeReplayFixtureFromTemporaryFile(t *testing.T) {
	replay := NewWindowsFluidVirtGuardedExecutorFakeRuntimeReplayMinimal()
	raw := marshalFixtureJSON(t, replay)
	tempFile := filepath.Join(t.TempDir(), "guarded_executor_fake_runtime_replay.json")
	if err := os.WriteFile(tempFile, raw, 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	loaded, err := LoadGuardedExecutorFakeRuntimeReplayFixtureFromTemporaryFile(tempFile)
	if err != nil {
		t.Fatalf("load replay fixture: %v", err)
	}
	if loaded.ExecutorFakeRuntimeReplayID != replay.ExecutorFakeRuntimeReplayID {
		t.Fatalf("unexpected replay id: %s", loaded.ExecutorFakeRuntimeReplayID)
	}
}

func TestLoadGuardedExecutorReplayRejectsSysFsPath(t *testing.T) {
	if _, err := LoadGuardedExecutorFakeRuntimeReplayFixtureFromTemporaryFile("/sys/fs/cgroup/test.json"); err == nil {
		t.Fatalf("expected /sys/fs/cgroup path rejection")
	}
}

func TestLoadGuardedExecutorReplayRejectsRawQMPQGAAndSecretMaterial(t *testing.T) {
	tempFile := filepath.Join(t.TempDir(), "forbidden-material.json")
	raw := []byte(`{"payload":"raw_qmp_payload raw_qga_payload","token":"abc"}`)
	if err := os.WriteFile(tempFile, raw, 0o600); err != nil {
		t.Fatalf("write forbidden fixture: %v", err)
	}
	if _, err := LoadGuardedExecutorFakeRuntimeReplayFixtureFromTemporaryFile(tempFile); err == nil {
		t.Fatalf("expected forbidden material rejection")
	}
}

func marshalFixtureJSON(t *testing.T, value any) []byte {
	t.Helper()
	raw, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		t.Fatalf("marshal fixture: %v", err)
	}
	return raw
}
