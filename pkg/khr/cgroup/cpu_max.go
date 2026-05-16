package cgroup

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const defaultCPUPeriodUS = 100000

// MilliCPUToMaxLine formats cgroup v2 cpu.max for the given millicpu quota.
func MilliCPUToMaxLine(milliCPU int64) string {
	if milliCPU <= 0 {
		return "max"
	}
	quota := milliCPU * defaultCPUPeriodUS / 1000
	return fmt.Sprintf("%d %d", quota, defaultCPUPeriodUS)
}

// ReadCPUMax reads cpu.max from dir under allowPathPrefix (read-only).
func ReadCPUMax(dir, allowPathPrefix string) (string, error) {
	resolved, err := AbsClean(dir)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(allowPathPrefix) != "" {
		ap, err := AbsClean(allowPathPrefix)
		if err != nil {
			return "", err
		}
		if !IsSubpath(ap, resolved) {
			return "", fmt.Errorf("cgroup path %q outside allowPathPrefix %q", resolved, ap)
		}
	}
	data, err := os.ReadFile(filepath.Join(resolved, "cpu.max"))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

// WriteCPUMax writes cpu.max under allowPathPrefix after path validation.
func WriteCPUMax(dir, allowPathPrefix, maxLine string) error {
	resolved, err := AbsClean(dir)
	if err != nil {
		return err
	}
	if strings.TrimSpace(allowPathPrefix) != "" {
		ap, err := AbsClean(allowPathPrefix)
		if err != nil {
			return err
		}
		if !IsSubpath(ap, resolved) {
			return fmt.Errorf("cgroup path %q outside allowPathPrefix %q", resolved, ap)
		}
	}
	if err := os.MkdirAll(resolved, 0o755); err != nil {
		return err
	}
	target := filepath.Join(resolved, "cpu.max")
	return os.WriteFile(target, []byte(strings.TrimSpace(maxLine)+"\n"), 0o644)
}

// ParseMilliCPUFromMaxLine best-effort parses millicpu from a cpu.max line.
func ParseMilliCPUFromMaxLine(line string) (int64, error) {
	line = strings.TrimSpace(line)
	if line == "" || line == "max" {
		return 0, nil
	}
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return 0, fmt.Errorf("invalid cpu.max line %q", line)
	}
	quota, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, err
	}
	period, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil || period <= 0 {
		return 0, fmt.Errorf("invalid cpu.max period")
	}
	return quota * 1000 / period, nil
}
