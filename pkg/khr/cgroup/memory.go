package cgroup

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const minSandboxMemoryBytes int64 = 32 * 1024 * 1024 // 32Mi floor

// ParseMemoryValue parses cgroup memory.max or memory.high (max or bytes).
func ParseMemoryValue(line string) (bytes int64, unlimited bool, err error) {
	line = strings.TrimSpace(line)
	if line == "" || strings.EqualFold(line, "max") {
		return 0, true, nil
	}
	v, err := strconv.ParseInt(line, 10, 64)
	if err != nil {
		return 0, false, err
	}
	return v, false, nil
}

// FormatMemoryValue formats bytes for memory.max / memory.high.
func FormatMemoryValue(bytes int64) string {
	if bytes <= 0 {
		return "max"
	}
	return strconv.FormatInt(bytes, 10)
}

// ReadMemoryMax reads memory.max from dir under allowPathPrefix.
func ReadMemoryMax(dir, allowPathPrefix string) (string, error) {
	return readMemoryFile(dir, allowPathPrefix, "memory.max")
}

// ReadMemoryHigh reads memory.high from dir under allowPathPrefix.
func ReadMemoryHigh(dir, allowPathPrefix string) (string, error) {
	return readMemoryFile(dir, allowPathPrefix, "memory.high")
}

// WriteMemoryMax writes memory.max under allowPathPrefix.
func WriteMemoryMax(dir, allowPathPrefix, value string) error {
	return writeMemoryFile(dir, allowPathPrefix, "memory.max", value)
}

// WriteMemoryHigh writes memory.high under allowPathPrefix.
func WriteMemoryHigh(dir, allowPathPrefix, value string) error {
	return writeMemoryFile(dir, allowPathPrefix, "memory.high", value)
}

// ComputeMemoryTarget applies scaleUp/scaleDown/envelope to current limit.
func ComputeMemoryTarget(currentBytes int64, currentUnlimited bool, deltaBytes int64, mode string) (int64, error) {
	base := currentBytes
	if currentUnlimited {
		base = 256 * 1024 * 1024 // default when current is max
	}
	switch strings.ToLower(mode) {
	case "scaleup":
		return base + deltaBytes, nil
	case "scaledown":
		next := base - deltaBytes
		if next < minSandboxMemoryBytes {
			next = minSandboxMemoryBytes
		}
		return next, nil
	case "envelope":
		if deltaBytes <= 0 {
			return 0, fmt.Errorf("envelope memory requires positive bytes")
		}
		return deltaBytes, nil
	default:
		return 0, fmt.Errorf("unsupported memory mode %q", mode)
	}
}

func readMemoryFile(dir, allowPathPrefix, name string) (string, error) {
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
	data, err := os.ReadFile(filepath.Join(resolved, name))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

func writeMemoryFile(dir, allowPathPrefix, name, value string) error {
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
	target := filepath.Join(resolved, name)
	return os.WriteFile(target, []byte(strings.TrimSpace(value)+"\n"), 0o644)
}
