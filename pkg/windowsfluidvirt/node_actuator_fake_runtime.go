package windowsfluidvirt

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func DefaultWindowsFluidVirtFakeRuntimeBoundary() WindowsFluidVirtFakeRuntimeBoundary {
	return WindowsFluidVirtFakeRuntimeBoundary{
		FakeRuntimeOnly:        true,
		UsesTemporaryFilesOnly: true,
		TouchesRealCgroup:      false,
		TouchesRealQMP:         false,
		TouchesRealQGA:         false,
		TouchesHostRuntime:     false,
		RequiresNoPrivileges:   true,
		DeterministicReplay:    true,
		SafeForCI:              true,
		ClaimBoundary:          "fake_runtime_boundary_only",
	}
}

func ValidateFakeRuntimeBoundary(boundary WindowsFluidVirtFakeRuntimeBoundary) error {
	if !boundary.FakeRuntimeOnly || !boundary.UsesTemporaryFilesOnly {
		return fmt.Errorf("fake runtime boundary must stay fake-runtime only and temporary-files only")
	}
	if boundary.TouchesRealCgroup || boundary.TouchesRealQMP || boundary.TouchesRealQGA || boundary.TouchesHostRuntime {
		return fmt.Errorf("fake runtime boundary must not touch real runtime surfaces")
	}
	if !boundary.RequiresNoPrivileges || !boundary.DeterministicReplay || !boundary.SafeForCI {
		return fmt.Errorf("fake runtime boundary must require no privileges and stay deterministic/ci-safe")
	}
	return nil
}

func ReplayFixtureFromTemporaryFile(path string, boundary WindowsFluidVirtFakeRuntimeBoundary) (WindowsFluidVirtNodeActuatorReadonlyReplay, error) {
	if err := ValidateFakeRuntimeBoundary(boundary); err != nil {
		return WindowsFluidVirtNodeActuatorReadonlyReplay{}, err
	}
	clean := filepath.Clean(path)
	if strings.HasPrefix(clean, "/sys/fs/cgroup") {
		return WindowsFluidVirtNodeActuatorReadonlyReplay{}, fmt.Errorf("real cgroup path is forbidden for fake runtime replay")
	}
	if !strings.Contains(clean, os.TempDir()) {
		return WindowsFluidVirtNodeActuatorReadonlyReplay{}, fmt.Errorf("fake runtime replay requires temporary file path")
	}
	raw, err := os.ReadFile(clean)
	if err != nil {
		return WindowsFluidVirtNodeActuatorReadonlyReplay{}, err
	}
	var replay WindowsFluidVirtNodeActuatorReadonlyReplay
	if err := json.Unmarshal(raw, &replay); err != nil {
		return WindowsFluidVirtNodeActuatorReadonlyReplay{}, err
	}
	if replay.RuntimeMutationEnabled || replay.ActuatorRuntimeEnabled || replay.CgroupWriteEnabled {
		return WindowsFluidVirtNodeActuatorReadonlyReplay{}, fmt.Errorf("replay payload violates readonly boundary")
	}
	return replay, nil
}
