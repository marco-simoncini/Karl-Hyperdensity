package telemetry

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/cgroup"
)

// MetricsBundle is cgroup v2 scalar evidence collected read-only.
type MetricsBundle struct {
	CPUStat       map[string]int64    `json:"cpuStat,omitempty"`
	MemoryCurrent string              `json:"memoryCurrent,omitempty"`
	MemoryMax     string              `json:"memoryMax,omitempty"`
	MemoryEvents  map[string]int64    `json:"memoryEvents,omitempty"`
	IOStat        []map[string]string `json:"ioStat,omitempty"`
}

// ReadCgroupV2Metrics reads best-effort cgroup v2 files under resolvedDir (already validated).
func ReadCgroupV2Metrics(resolvedDir string) (m MetricsBundle, warnings []string, blocked []string) {
	m.CPUStat = map[string]int64{}
	m.MemoryEvents = map[string]int64{}

	read := func(name string) []byte {
		b, w, bl := cgroup.ReadFileInResolvedDir(resolvedDir, name)
		warnings = append(warnings, w...)
		blocked = append(blocked, bl...)
		return b
	}

	if b := read("cpu.stat"); len(b) > 0 {
		parsed, w := parseKVInt64Lines(string(b))
		warnings = append(warnings, w...)
		for k, v := range parsed {
			m.CPUStat[k] = v
		}
	}

	if b := read("memory.current"); len(b) > 0 {
		m.MemoryCurrent = strings.TrimSpace(string(b))
	}

	if b := read("memory.max"); len(b) > 0 {
		m.MemoryMax = strings.TrimSpace(string(b))
	}

	if b := read("memory.events"); len(b) > 0 {
		parsed, w := parseKVInt64Lines(string(b))
		warnings = append(warnings, w...)
		for k, v := range parsed {
			m.MemoryEvents[k] = v
		}
	}

	if b := read("io.stat"); len(b) > 0 {
		m.IOStat = parseIOStat(string(b))
	}

	if len(m.CPUStat) == 0 && m.MemoryCurrent == "" {
		blocked = append(blocked, "telemetry unusable: both cpu.stat and memory.current are missing or unreadable under cgroup path")
	}
	return m, warnings, blocked
}

func parseKVInt64Lines(s string) (map[string]int64, []string) {
	out := make(map[string]int64)
	var warns []string
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		key := fields[0]
		valStr := fields[1]
		n, err := strconv.ParseInt(valStr, 10, 64)
		if err != nil {
			warns = append(warns, fmt.Sprintf("skip non-int field %q in cgroup kv line", key))
			continue
		}
		out[key] = n
	}
	return out, warns
}

func parseIOStat(s string) []map[string]string {
	var blocks [][]string
	cur := []string{}
	for _, line := range strings.Split(s, "\n") {
		if strings.TrimSpace(line) == "" {
			if len(cur) > 0 {
				blocks = append(blocks, cur)
				cur = []string{}
			}
			continue
		}
		cur = append(cur, strings.TrimSpace(line))
	}
	if len(cur) > 0 {
		blocks = append(blocks, cur)
	}
	var out []map[string]string
	for _, lines := range blocks {
		m := make(map[string]string)
		for _, line := range lines {
			fields := strings.Fields(line)
			if len(fields) < 2 {
				continue
			}
			m[fields[0]] = strings.Join(fields[1:], " ")
		}
		if len(m) > 0 {
			out = append(out, m)
		}
	}
	return out
}
