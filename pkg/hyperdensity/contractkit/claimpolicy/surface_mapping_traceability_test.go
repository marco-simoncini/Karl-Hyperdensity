package claimpolicy

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateDashboardFileTraceability(t *testing.T) {
	if err := ValidateDashboardFileTraceability(); err != nil {
		t.Fatal(err)
	}
}

func TestEveryClaimHasDashboardFiles(t *testing.T) {
	for _, rule := range Catalog() {
		files := DashboardFilesForClaim(rule.ID)
		if len(files) == 0 {
			t.Fatalf("claim %q must have traced dashboard files", rule.ID)
		}
	}
}

func TestDashboardPathsNoAbsOrDotDot(t *testing.T) {
	for _, m := range SurfaceMappings() {
		for _, p := range m.DashboardFiles {
			if filepath.IsAbs(p) {
				t.Fatalf("path must not be absolute: %q", p)
			}
			if strings.Contains(p, "..") {
				t.Fatalf("path must not contain '..': %q", p)
			}
		}
	}
}

func TestNoDuplicatePathsPerClaim(t *testing.T) {
	for _, rule := range Catalog() {
		seen := make(map[string]int)
		for _, m := range SurfaceMappings() {
			if m.ClaimID != rule.ID {
				continue
			}
			for _, p := range m.DashboardFiles {
				seen[p]++
			}
		}
		dedup := DashboardFilesForClaim(rule.ID)
		if len(dedup) != len(seen) {
			t.Fatalf("claim %q: dedup file count %d != unique path count %d", rule.ID, len(dedup), len(seen))
		}
	}
}

func TestRuntimeImportFreezeUsesAuditScripts(t *testing.T) {
	var paths []string
	for _, m := range SurfaceMappings() {
		if m.ClaimID != string(ClaimNoRuntimeContractsImport) {
			continue
		}
		paths = append(paths, m.DashboardFiles...)
	}
	hasAudit := false
	hasAllow := false
	for _, p := range paths {
		if strings.HasSuffix(p, "audit_contractkit_runtime_imports.sh") {
			hasAudit = true
		}
		if strings.HasSuffix(p, "contractkit_runtime_import_allowlist.txt") {
			hasAllow = true
		}
	}
	if !hasAudit || !hasAllow {
		t.Fatalf("no_runtime_contracts_import must trace audit script and allowlist, got %#v", paths)
	}
}

func TestNoProductionMutationTracesPolicyReleaseLive(t *testing.T) {
	var policy, consistency, rel, live bool
	for _, m := range SurfaceMappings() {
		if m.ClaimID != string(ClaimNoProductionMutation) {
			continue
		}
		for _, p := range m.DashboardFiles {
			switch {
			case strings.HasSuffix(p, "hyperdensity_parent_fabric_policy_pack_v1.go"):
				policy = true
			case strings.HasSuffix(p, "hyperdensity_parent_fabric_policy_pack_consistency_v1.go"):
				consistency = true
			case strings.HasSuffix(p, "hyperdensity_parent_fabric_release_support_matrix_v1.go"):
				rel = true
			case strings.HasSuffix(p, "hyperdensity_parent_fabric_live_resource_authority_v1.go"):
				live = true
			}
		}
	}
	if !policy || !consistency || !rel || !live {
		t.Fatalf("no_production_mutation must trace policy_pack, consistency, release matrix, and live authority (policy=%v consistency=%v rel=%v live=%v)", policy, consistency, rel, live)
	}
}

func TestNoWindowsTracesVmLaneFiles(t *testing.T) {
	files := DashboardFilesForClaim(string(ClaimNoWindowsHyperdensityApply))
	var lane, collector bool
	for _, p := range files {
		if strings.HasSuffix(p, "hyperdensity_parent_fabric_vm_lane_readiness_v1.go") {
			lane = true
		}
		if strings.HasSuffix(p, "hyperdensity_parent_fabric_vm_runtime_evidence_collector_v1.go") {
			collector = true
		}
	}
	if !lane || !collector {
		t.Fatalf("no_windows_hyperdensity_apply must trace vm lane + evidence collector, got %#v", files)
	}
}

func TestWindowsLaneDisabledTracesReadonly(t *testing.T) {
	files := DashboardFilesForClaim(string(ClaimWindowsLaneDisabled))
	if len(files) < 2 {
		t.Fatalf("expected multiple readonly traces, got %#v", files)
	}
}
