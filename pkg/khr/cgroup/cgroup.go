// Package cgroup provides cgroup v2 detection and envelope planning without writes (Sprint 5+).
package cgroup

import (
	"os"
	"path/filepath"
	"strings"
)

// envTestCgroupVersion forces DetectVersion for golden tests (non-production).
const envTestCgroupVersion = "KHR_TEST_CGROUP_VERSION"

// Version is cgroup hierarchy flavour.
type Version string

const (
	V2      Version = "v2"
	V1      Version = "v1"
	Unknown Version = "unknown"
)

// DetectVersion returns cgroup v2 if unified hierarchy is mounted (best-effort).
// When KHR_TEST_CGROUP_VERSION is set to v1|v2|unknown, that value wins (tests only).
func DetectVersion() Version {
	if ev := strings.TrimSpace(os.Getenv(envTestCgroupVersion)); ev != "" {
		switch Version(strings.ToLower(ev)) {
		case V2, V1, Unknown:
			return Version(strings.ToLower(ev))
		default:
			return Unknown
		}
	}
	st, err := os.Stat(filepath.Join(UnifiedCgroupMount, "cgroup.controllers"))
	if err == nil && !st.IsDir() {
		return V2
	}
	if st, err := os.Stat(UnifiedCgroupMount); err == nil && st.IsDir() {
		return V1
	}
	return Unknown
}

// EnvelopePlan describes intended cgroup writes without performing them.
type EnvelopePlan struct {
	CgroupVersion  Version  `json:"cgroupVersion"`
	CPUMaxDelta    string   `json:"cpuMaxDelta,omitempty"`
	MemoryMaxDelta string   `json:"memoryMaxDelta,omitempty"`
	WouldWrite     bool     `json:"wouldWrite"`
	WritePaths     []string `json:"writePaths,omitempty"`
}

// PlanEnvelope computes a no-op plan for CPU/memory envelope adjustments.
// Sprint 6: allowWrite is ignored; real cgroup writes are always disabled in this agent build.
func PlanEnvelope(_ bool, cpuMaxDelta, memoryMaxDelta string) EnvelopePlan {
	return EnvelopePlan{
		CgroupVersion:  DetectVersion(),
		CPUMaxDelta:    cpuMaxDelta,
		MemoryMaxDelta: memoryMaxDelta,
		WouldWrite:     false,
		WritePaths:     nil,
	}
}
