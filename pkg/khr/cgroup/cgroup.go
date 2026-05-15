// Package cgroup provides cgroup v2 detection and envelope planning without writes (Sprint 5).
package cgroup

import (
	"os"
	"path/filepath"
)

const unifiedMount = "/sys/fs/cgroup"

// Version is cgroup hierarchy flavour.
type Version string

const (
	V2      Version = "v2"
	V1      Version = "v1"
	Unknown Version = "unknown"
)

// DetectVersion returns cgroup v2 if unified hierarchy is mounted (best-effort).
func DetectVersion() Version {
	st, err := os.Stat(filepath.Join(unifiedMount, "cgroup.controllers"))
	if err == nil && !st.IsDir() {
		return V2
	}
	// cgroup v1 often has separate controllers; keep stub conservative
	if st, err := os.Stat(unifiedMount); err == nil && st.IsDir() {
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
// allowWrite must be true only when caller passed --allow-unsafe-apply; Sprint 5 tests never enable it.
func PlanEnvelope(allowWrite bool, cpuMaxDelta, memoryMaxDelta string) EnvelopePlan {
	would := allowWrite && (cpuMaxDelta != "" || memoryMaxDelta != "")
	paths := []string{}
	if would {
		paths = append(paths, filepath.Join(unifiedMount, "<slice>", "cpu.max"))
		paths = append(paths, filepath.Join(unifiedMount, "<slice>", "memory.max"))
	}
	return EnvelopePlan{
		CgroupVersion:  DetectVersion(),
		CPUMaxDelta:    cpuMaxDelta,
		MemoryMaxDelta: memoryMaxDelta,
		WouldWrite:     would,
		WritePaths:     paths,
	}
}
