package cgroup

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ValidateCgroupPathForTelemetry checks that cgroupPath resolves to a directory suitable for read-only metric collection.
// When allowPathPrefix is non-empty, the resolved path must lie under that prefix (after EvalSymlinks).
func ValidateCgroupPathForTelemetry(cgroupPath, allowPathPrefix string) (resolved string, warnings []string, blocked []string) {
	cp, err := AbsClean(cgroupPath)
	if err != nil {
		return "", nil, []string{"invalid cgroup path: " + err.Error()}
	}
	real, err := filepath.EvalSymlinks(cp)
	if err != nil {
		warnings = append(warnings, fmt.Sprintf("cannot resolve cgroup path symlinks: %v", err))
		return "", warnings, nil
	}
	realAbs, err := AbsClean(real)
	if err != nil {
		return "", warnings, []string{"resolved cgroup path invalid: " + err.Error()}
	}
	if strings.TrimSpace(allowPathPrefix) != "" {
		ap, err := AbsClean(allowPathPrefix)
		if err != nil {
			return "", warnings, []string{"invalid allowPathPrefix: " + err.Error()}
		}
		if !IsSubpath(ap, realAbs) {
			return "", warnings, []string{fmt.Sprintf("resolved cgroup path %q is outside allowPathPrefix %q", realAbs, ap)}
		}
	}
	st, err := os.Stat(realAbs)
	if err != nil {
		warnings = append(warnings, fmt.Sprintf("cannot stat cgroup path: %v", err))
		return "", warnings, nil
	}
	if !st.IsDir() {
		return "", warnings, []string{fmt.Sprintf("cgroup path %q is not a directory", realAbs)}
	}
	return realAbs, warnings, nil
}

// ReadFileInResolvedDir reads a file relative to dirResolved (already symlink-resolved).
// The joined path must resolve under dirResolved after EvalSymlinks (blocks escape).
func ReadFileInResolvedDir(dirResolved, relName string) (data []byte, warnings []string, blocked []string) {
	if strings.Contains(relName, string(filepath.Separator)) || relName == "" || relName == "." || relName == ".." {
		return nil, nil, []string{fmt.Sprintf("invalid relative metric name %q", relName)}
	}
	p := filepath.Join(dirResolved, relName)
	real, err := filepath.EvalSymlinks(p)
	if err != nil {
		warnings = append(warnings, fmt.Sprintf("cannot resolve %q: %v", relName, err))
		return nil, warnings, nil
	}
	realAbs, err := AbsClean(real)
	if err != nil {
		return nil, warnings, []string{relName + ": invalid resolved path"}
	}
	if !IsSubpath(dirResolved, realAbs) {
		return nil, warnings, []string{fmt.Sprintf("metric file %q resolves outside cgroup directory (symlink escape blocked)", relName)}
	}
	b, err := os.ReadFile(realAbs)
	if err != nil {
		warnings = append(warnings, fmt.Sprintf("cannot read %q: %v", relName, err))
		return nil, warnings, nil
	}
	return b, warnings, nil
}
