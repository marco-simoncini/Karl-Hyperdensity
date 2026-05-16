package resourcelease

import (
	"os"
	"path/filepath"
)

// Baseline captures pre-apply sandbox state for rollback.
type Baseline struct {
	ID          string            `json:"id"`
	SandboxDir  string            `json:"sandboxDir"`
	MarkerPath  string            `json:"markerPath"`
	MarkerBytes []byte            `json:"-"`
	Extra       map[string]string `json:"extra,omitempty"`
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
