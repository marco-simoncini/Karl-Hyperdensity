// Package blockers exports Hyperdensity gate/blocker IDs aligned with Karl-Dashboard
// parent-fabric runtime vocabulary (M1 golden anchor). Read-only catalog — no cluster IO.
package blockers

import "sort"

// Severity levels for catalog entries (stable strings for tests and docs).
const (
	SeverityCritical = "critical"
	SeverityHigh     = "high"
	SeverityMedium   = "medium"
	SeverityInfo     = "info"
)

// Gate / blocker IDs used across Dashboard parent-fabric collectors and remediation copy.
const (
	IDNoWindowsLane                  = "no_windows_lane"
	IDNoProductionMutation           = "no_production_mutation"
	IDKeepWindowsLaneDisabled        = "keep_windows_lane_disabled"
	IDWindowsDisabled                = "windows_disabled"
	IDDryRunOnly                     = "dry_run_only"
	IDRuntimeApplyDisabled           = "runtime_apply_disabled"
	IDUnsupportedBroadVMExecution    = "unsupported_broad_vm_execution"
	IDUnsupportedBroadMemoryExec     = "unsupported_broad_memory_execution"
	IDUnsupportedMultiContainerWiden = "unsupported_multi_container_widening"
	IDUnsupportedBroadAutomation     = "unsupported_broad_automation"
)

// Blocker is a catalog entry for a known gate or blocker ID.
type Blocker struct {
	ID       string
	Severity string
	Message  string
}

var catalog = map[string]Blocker{
	IDNoWindowsLane: {
		ID:       IDNoWindowsLane,
		Severity: SeverityCritical,
		Message:  "Windows lane must remain disabled and out-of-scope for live execution.",
	},
	IDNoProductionMutation: {
		ID:       IDNoProductionMutation,
		Severity: SeverityCritical,
		Message:  "Production mutation is not allowed on the active Hyperdensity surface.",
	},
	IDKeepWindowsLaneDisabled: {
		ID:       IDKeepWindowsLaneDisabled,
		Severity: SeverityHigh,
		Message:  "Remediation action: keep Windows lane disabled.",
	},
	IDWindowsDisabled: {
		ID:       IDWindowsDisabled,
		Severity: SeverityInfo,
		Message:  "Observed Windows lane posture: disabled.",
	},
	IDDryRunOnly: {
		ID:       IDDryRunOnly,
		Severity: SeverityHigh,
		Message:  "Execution category is dry-run only; no runtime apply.",
	},
	IDRuntimeApplyDisabled: {
		ID:       IDRuntimeApplyDisabled,
		Severity: SeverityCritical,
		Message:  "Runtime apply path is disabled until operator-controlled gates pass.",
	},
	IDUnsupportedBroadVMExecution: {
		ID:       IDUnsupportedBroadVMExecution,
		Severity: SeverityHigh,
		Message:  "Broad VM execution beyond narrow cohort is unsupported.",
	},
	IDUnsupportedBroadMemoryExec: {
		ID:       IDUnsupportedBroadMemoryExec,
		Severity: SeverityHigh,
		Message:  "Broad memory execution beyond narrow cohort is unsupported.",
	},
	IDUnsupportedMultiContainerWiden: {
		ID:       IDUnsupportedMultiContainerWiden,
		Severity: SeverityHigh,
		Message:  "Multi-container widening is unsupported on the current surface.",
	},
	IDUnsupportedBroadAutomation: {
		ID:       IDUnsupportedBroadAutomation,
		Severity: SeverityHigh,
		Message:  "Broad autonomous automation is unsupported; operator control required.",
	},
}

// Known reports whether id is in the M1 blocker catalog.
func Known(id string) bool {
	_, ok := catalog[id]
	return ok
}

// Severity returns the catalog severity for id, or empty string if unknown.
func Severity(id string) (string, bool) {
	b, ok := catalog[id]
	if !ok {
		return "", false
	}
	return b.Severity, true
}

// Catalog returns all blockers sorted by ID.
func Catalog() []Blocker {
	out := make([]Blocker, 0, len(catalog))
	for _, b := range catalog {
		out = append(out, b)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}
