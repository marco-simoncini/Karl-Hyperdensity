package telemetry

import (
	"os"
	"strings"
	"time"
)

// Evidence wraps observability metadata for Hyperdensity / Grande Padre consumers.
type Evidence struct {
	ObservedAt     string   `json:"observedAt"`
	Source         string   `json:"source"`
	Confidence     string   `json:"confidence"`
	Warnings       []string `json:"warnings"`
	BlockedReasons []string `json:"blockedReasons"`
}

const sourceCgroupV2 = "cgroup-v2"

// BuildEvidence merges policy warnings/blocks with reader output and sets confidence.
func BuildEvidence(warnings, blocked []string, m MetricsBundle) Evidence {
	return Evidence{
		ObservedAt:     evidenceNow().Format(time.RFC3339),
		Source:         sourceCgroupV2,
		Confidence:     confidenceFor(m),
		Warnings:       dedupeStrings(warnings),
		BlockedReasons: dedupeStrings(blocked),
	}
}

func evidenceNow() time.Time {
	s := os.Getenv("KHR_TEST_TELEMETRY_NOW")
	if s == "" {
		return time.Now().UTC()
	}
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		t2, err2 := time.Parse(time.RFC3339, s)
		if err2 != nil {
			return time.Now().UTC()
		}
		return t2.UTC()
	}
	return t.UTC()
}

func confidenceFor(m MetricsBundle) string {
	hasCPU := len(m.CPUStat) > 0
	hasMem := m.MemoryCurrent != ""
	switch {
	case hasCPU && hasMem:
		return "high"
	case hasCPU || hasMem:
		return "medium"
	default:
		return "low"
	}
}

func dedupeStrings(in []string) []string {
	seen := make(map[string]struct{})
	var out []string
	for _, s := range in {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	if len(out) == 0 {
		return []string{}
	}
	return out
}
