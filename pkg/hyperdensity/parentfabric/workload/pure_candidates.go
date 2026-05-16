package workload

import "strings"

const (
	WorkloadPackageVersion     = "v0.0.0-sprint52"
	WorkloadHelpersSourceFile  = "pkg/server/hyperdensity_parent_fabric_workload_helpers.go"
)

// AppsWorkloadResource maps a Kubernetes apps workload kind to its API resource segment.
// Matches Dashboard hyperdensityAppsWorkloadResource (stdlib-only copy-contract).
func AppsWorkloadResource(kind string) (resource string, ok bool) {
	switch strings.TrimSpace(kind) {
	case "Deployment":
		return "deployments", true
	case "StatefulSet":
		return "statefulsets", true
	default:
		return "", false
	}
}

// PilotWorkloadTerm returns the pilot terminology for a workload kind.
// Matches Dashboard hyperdensityPilotWorkloadTerm; unknown kinds yield ("workload", false).
func PilotWorkloadTerm(kind string) (term string, ok bool) {
	switch strings.TrimSpace(kind) {
	case "Deployment":
		return "Deployment", true
	case "StatefulSet":
		return "StatefulSet", true
	default:
		return "workload", false
	}
}

// ExecutionSupportsLiveApplyKind reports whether live apply is supported for the kind.
// Matches Dashboard hyperdensityExecutionSupportsLiveApplyKind exactly.
func ExecutionSupportsLiveApplyKind(kind string) bool {
	switch strings.TrimSpace(kind) {
	case "Deployment", "StatefulSet":
		return true
	default:
		return false
	}
}
