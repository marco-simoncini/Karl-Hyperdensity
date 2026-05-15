package cgroup

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// IsSubpath reports whether sub is root or a path strictly inside root (after Clean).
func IsSubpath(root, sub string) bool {
	root = filepath.Clean(root)
	sub = filepath.Clean(sub)
	if root == sub {
		return true
	}
	rel, err := filepath.Rel(root, sub)
	if err != nil {
		return false
	}
	return rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

// AbsClean returns filepath.Abs then filepath.Clean.
func AbsClean(p string) (string, error) {
	if p == "" {
		return "", fmt.Errorf("empty path")
	}
	a, err := filepath.Abs(p)
	if err != nil {
		return "", err
	}
	return filepath.Clean(a), nil
}

// PathUnderOptionalPrefix returns true when prefix is empty/whitespace, or sub is under prefix.
func PathUnderOptionalPrefix(sub, prefix string) bool {
	prefix = strings.TrimSpace(prefix)
	if prefix == "" {
		return true
	}
	pfx, err := AbsClean(prefix)
	if err != nil {
		return false
	}
	s, err := AbsClean(sub)
	if err != nil {
		return false
	}
	return IsSubpath(pfx, s)
}

// DiscoverableDir reports whether path resolves to a directory under scannedRoot, optionally under allowPrefix,
// and that symlink resolution does not escape scannedRoot.
func DiscoverableDir(scannedRoot, path, allowPathPrefix string) (ok bool, warnings []string, blocked []string) {
	root, err := AbsClean(scannedRoot)
	if err != nil {
		return false, nil, []string{"invalid scannedRoot: " + err.Error()}
	}
	p, err := AbsClean(path)
	if err != nil {
		return false, nil, []string{"invalid candidate path: " + err.Error()}
	}
	if !IsSubpath(root, p) {
		return false, nil, []string{fmt.Sprintf("path %q is outside scannedRoot %q", p, root)}
	}
	real, err := filepath.EvalSymlinks(p)
	if err != nil {
		warnings = append(warnings, fmt.Sprintf("cannot resolve symlinks for %q: %v", p, err))
		return false, warnings, nil
	}
	realAbs, err := AbsClean(real)
	if err != nil {
		return false, warnings, []string{"resolved path invalid: " + err.Error()}
	}
	if !IsSubpath(root, realAbs) {
		return false, warnings, []string{fmt.Sprintf("symlink chain escapes scannedRoot (resolved %q)", realAbs)}
	}
	if strings.TrimSpace(allowPathPrefix) != "" {
		ap, err := AbsClean(allowPathPrefix)
		if err != nil {
			return false, warnings, []string{"invalid allowPathPrefix: " + err.Error()}
		}
		if !IsSubpath(ap, realAbs) {
			return false, warnings, []string{fmt.Sprintf("resolved path %q is outside allowPathPrefix %q", realAbs, ap)}
		}
	}
	st, err := os.Stat(realAbs)
	if err != nil {
		warnings = append(warnings, fmt.Sprintf("cannot stat resolved path %q: %v", realAbs, err))
		return false, warnings, nil
	}
	if !st.IsDir() {
		blocked = append(blocked, fmt.Sprintf("resolved path %q is not a directory", realAbs))
		return false, warnings, blocked
	}
	return true, warnings, nil
}
