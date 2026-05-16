package resourcelease

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/cgroup"
)

// Baseline captures pre-apply sandbox state for rollback.
type Baseline struct {
	ID            string            `json:"id"`
	SandboxDir    string            `json:"sandboxDir"`
	MarkerPath    string            `json:"markerPath"`
	MarkerBytes   []byte            `json:"-"`
	CgroupCPUPath string            `json:"cgroupCpuPath,omitempty"`
	CPUMaxBefore  string            `json:"cpuMaxBefore,omitempty"`
	CPUMaxApplied string            `json:"cpuMaxApplied,omitempty"`
	Extra         map[string]string `json:"extra,omitempty"`
}

// CaptureBaseline records sandbox marker bytes before apply (or after for tests).
func CaptureBaseline(id, sandboxDir string) (Baseline, error) {
	b := Baseline{
		ID:         id,
		SandboxDir: sandboxDir,
		MarkerPath: filepath.Join(sandboxDir, "apply-marker.txt"),
		Extra:      map[string]string{"source": "karl-host-runtime"},
	}
	if raw, err := os.ReadFile(b.MarkerPath); err == nil {
		b.MarkerBytes = append([]byte(nil), raw...)
	}
	return b, nil
}

// CaptureCgroupBaseline records cpu.max before sandbox apply.
func CaptureCgroupBaseline(id, sandboxDir, cgroupPath, allowPrefix string, milliCPU int64) (Baseline, error) {
	b := Baseline{
		ID:            id,
		SandboxDir:    sandboxDir,
		MarkerPath:    filepath.Join(sandboxDir, "apply-marker.txt"),
		CgroupCPUPath: cgroupPath,
		Extra: map[string]string{
			"source":          "karl-host-runtime",
			"allowPathPrefix": allowPrefix,
		},
	}
	if err := os.MkdirAll(sandboxDir, 0o755); err != nil {
		return b, err
	}
	before, err := readCPUMaxOrMax(cgroupPath, allowPrefix)
	if err != nil {
		return b, err
	}
	b.CPUMaxBefore = before
	return b, SaveBaseline(b)
}

func readCPUMaxOrMax(cgroupPath, allowPrefix string) (string, error) {
	val, err := cgroup.ReadCPUMax(cgroupPath, allowPrefix)
	if err != nil {
		if os.IsNotExist(err) {
			return "max", nil
		}
		return "", err
	}
	if val == "" {
		return "max", nil
	}
	return val, nil
}

// SaveBaseline persists baseline JSON under sandboxDir.
func SaveBaseline(b Baseline) error {
	if b.SandboxDir == "" {
		return os.ErrInvalid
	}
	raw, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(b.SandboxDir, "baseline-"+b.ID+".json"), raw, 0o644)
}

// LoadBaseline loads baseline JSON from sandboxDir.
func LoadBaseline(id, sandboxDir string) (Baseline, error) {
	raw, err := os.ReadFile(filepath.Join(sandboxDir, "baseline-"+id+".json"))
	if err != nil {
		return Baseline{}, err
	}
	var b Baseline
	if err := json.Unmarshal(raw, &b); err != nil {
		return Baseline{}, err
	}
	return b, nil
}

// RollbackResult describes rollback outcome.
type RollbackResult struct {
	RolledBack bool     `json:"rolledBack"`
	Blocked    bool     `json:"blocked"`
	Reason     string   `json:"reason,omitempty"`
	Actions    []string `json:"actions,omitempty"`
}

// RollbackBaseline restores sandbox marker from baseline or removes apply marker.
func RollbackBaseline(b Baseline) RollbackResult {
	if b.SandboxDir == "" {
		return RollbackResult{Blocked: true, Reason: "empty sandboxDir"}
	}
	actions := []string{"restore sandbox apply marker"}
	if len(b.MarkerBytes) == 0 {
		_ = os.Remove(b.MarkerPath)
		actions = append(actions, "removed apply-marker.txt")
		return RollbackResult{RolledBack: true, Actions: actions}
	}
	if err := os.MkdirAll(b.SandboxDir, 0o755); err != nil {
		return RollbackResult{Blocked: true, Reason: err.Error()}
	}
	if err := os.WriteFile(b.MarkerPath, b.MarkerBytes, 0o644); err != nil {
		return RollbackResult{Blocked: true, Reason: err.Error()}
	}
	return RollbackResult{RolledBack: true, Actions: actions}
}
