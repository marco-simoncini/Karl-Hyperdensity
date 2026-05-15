// Package discovery implements read-only cgroup path discovery for KHR Linux (Sprint 7).
package discovery

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/cgroup"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/crdv1alpha1"
)

const discoveryModeReadOnly = "read-only"

type providerHandle struct {
	CgroupPath string `json:"cgroupPath"`
}

// CgroupDiscoveryOutput is the JSON shape for `discover-cgroups` mode (Sprint 7).
type CgroupDiscoveryOutput struct {
	Tool               string   `json:"tool,omitempty"`
	Version            string   `json:"version,omitempty"`
	Mode               string   `json:"mode,omitempty"`
	AgentID            string   `json:"agentId"`
	CgroupVersion      string   `json:"cgroupVersion"`
	DiscoveryMode      string   `json:"discoveryMode"`
	ScannedRoot        string   `json:"scannedRoot"`
	AllowedPathPrefix  string   `json:"allowedPathPrefix"`
	CandidatePaths     []string `json:"candidatePaths"`
	SelectedPath       string   `json:"selectedPath,omitempty"`
	BlockedReasons     []string `json:"blockedReasons"`
	Warnings           []string `json:"warnings"`
	MutationsForbidden bool     `json:"mutationsForbidden"`
}

// Run performs read-only discovery. Missing paths are non-fatal (blockedReasons + empty selectedPath).
func Run(agentID, scannedRoot, allowPathPrefix string, cell *crdv1alpha1.Cell) *CgroupDiscoveryOutput {
	out := &CgroupDiscoveryOutput{
		AgentID:            agentID,
		CgroupVersion:      string(cgroup.DetectVersion()),
		DiscoveryMode:      discoveryModeReadOnly,
		AllowedPathPrefix:  strings.TrimSpace(allowPathPrefix),
		CandidatePaths:     []string{},
		BlockedReasons:     []string{},
		Warnings:           []string{},
		MutationsForbidden: true,
	}
	root := strings.TrimSpace(scannedRoot)
	if root == "" {
		root = cgroup.DefaultScannedRoot()
	}
	absRoot, err := cgroup.AbsClean(root)
	if err != nil {
		out.BlockedReasons = append(out.BlockedReasons, "invalid scannedRoot: "+err.Error())
		out.ScannedRoot = root
		return out
	}
	out.ScannedRoot = absRoot

	if st, err := os.Lstat(absRoot); err != nil {
		out.Warnings = append(out.Warnings, fmt.Sprintf("scannedRoot not accessible: %v", err))
	} else if !st.IsDir() {
		out.Warnings = append(out.Warnings, "scannedRoot is not a directory")
	}

	if out.AllowedPathPrefix != "" {
		ap, err := cgroup.AbsClean(out.AllowedPathPrefix)
		if err != nil {
			out.BlockedReasons = append(out.BlockedReasons, "invalid allowPathPrefix: "+err.Error())
		} else {
			out.AllowedPathPrefix = ap
			if !cgroup.IsSubpath(absRoot, ap) {
				out.Warnings = append(out.Warnings, fmt.Sprintf("allowPathPrefix %q is not under scannedRoot (policy may reject all candidates)", ap))
			}
		}
	}

	rawCandidates, parseWarn := buildCandidates(absRoot, cell)
	out.Warnings = append(out.Warnings, parseWarn...)

	seen := make(map[string]struct{})
	var ordered []string
	for _, c := range rawCandidates {
		c = strings.TrimSpace(c)
		if c == "" {
			continue
		}
		ac, err := cgroup.AbsClean(c)
		if err != nil {
			out.Warnings = append(out.Warnings, fmt.Sprintf("skip invalid candidate %q: %v", c, err))
			continue
		}
		if _, ok := seen[ac]; ok {
			continue
		}
		seen[ac] = struct{}{}
		ordered = append(ordered, ac)
	}
	out.CandidatePaths = ordered

	for _, cand := range ordered {
		ok, warn, blocked := cgroup.DiscoverableDir(absRoot, cand, out.AllowedPathPrefix)
		out.Warnings = append(out.Warnings, warn...)
		if ok {
			out.SelectedPath = cand
			return out
		}
		if len(blocked) > 0 {
			out.BlockedReasons = append(out.BlockedReasons, fmt.Sprintf("candidate %q: %s", cand, strings.Join(blocked, "; ")))
		}
	}

	if out.SelectedPath == "" {
		out.BlockedReasons = append(out.BlockedReasons, "no discoverable cgroup directory matched heuristics or providerHandle under scannedRoot (read-only discovery only)")
	}
	return out
}

func buildCandidates(scannedRoot string, cell *crdv1alpha1.Cell) (out []string, warnings []string) {
	var shellName, cellName string
	if cell != nil {
		shellName = strings.TrimSpace(cell.Spec.ShellRef.Name)
		cellName = strings.TrimSpace(cell.Metadata.Name)
		explicit, pw := cgroupPathFromCell(cell)
		warnings = append(warnings, pw...)
		if explicit != "" {
			out = append(out, explicit)
			out = append(out, rebaseHostCgroupToScan(explicit, scannedRoot))
		}
	}

	karl := filepath.Join(scannedRoot, "karl.slice")
	if shellName != "" {
		out = append(out, filepath.Join(karl, "karl-shell-"+shellName+".scope"))
		out = append(out, filepath.Join(karl, "karl-shell-"+shellName))
	}
	if cellName != "" {
		out = append(out, filepath.Join(karl, cellName))
	}
	out = append(out, karl)
	return out, warnings
}

func cgroupPathFromCell(cell *crdv1alpha1.Cell) (path string, parseWarnings []string) {
	if cell == nil || len(cell.Spec.ProviderHandle) == 0 {
		return "", nil
	}
	var h providerHandle
	if err := json.Unmarshal(cell.Spec.ProviderHandle, &h); err != nil {
		return "", []string{"providerHandle parse: " + err.Error()}
	}
	return strings.TrimSpace(h.CgroupPath), nil
}

func rebaseHostCgroupToScan(explicit, scannedRoot string) string {
	hostClean := filepath.Clean(cgroup.UnifiedCgroupMount)
	exp := filepath.Clean(explicit)
	if exp == hostClean || strings.HasPrefix(exp, hostClean+string(filepath.Separator)) {
		rel, err := filepath.Rel(hostClean, exp)
		if err != nil || rel == "." {
			return scannedRoot
		}
		return filepath.Join(scannedRoot, rel)
	}
	return explicit
}
