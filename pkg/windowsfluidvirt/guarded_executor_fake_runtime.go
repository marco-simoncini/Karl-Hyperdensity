package windowsfluidvirt

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func LoadGuardedExecutorFakeRuntimeReplayFixtureFromTemporaryFile(path string) (WindowsFluidVirtGuardedExecutorFakeRuntimeReplay, error) {
	clean := filepath.Clean(path)
	if strings.HasPrefix(clean, "/sys/fs/cgroup") {
		return WindowsFluidVirtGuardedExecutorFakeRuntimeReplay{}, fmt.Errorf("real cgroup path forbidden")
	}
	if !strings.Contains(clean, os.TempDir()) {
		return WindowsFluidVirtGuardedExecutorFakeRuntimeReplay{}, fmt.Errorf("fixture must be loaded from temporary path")
	}
	raw, err := os.ReadFile(clean)
	if err != nil {
		return WindowsFluidVirtGuardedExecutorFakeRuntimeReplay{}, err
	}
	if hasForbiddenMaterial(string(raw)) {
		return WindowsFluidVirtGuardedExecutorFakeRuntimeReplay{}, fmt.Errorf("fixture includes forbidden raw runtime/secret material")
	}
	var replay WindowsFluidVirtGuardedExecutorFakeRuntimeReplay
	if err := json.Unmarshal(raw, &replay); err != nil {
		return WindowsFluidVirtGuardedExecutorFakeRuntimeReplay{}, err
	}
	if err := ValidateWindowsFluidVirtGuardedExecutorFakeRuntimeReplay(replay); err != nil {
		return WindowsFluidVirtGuardedExecutorFakeRuntimeReplay{}, err
	}
	return replay, nil
}
