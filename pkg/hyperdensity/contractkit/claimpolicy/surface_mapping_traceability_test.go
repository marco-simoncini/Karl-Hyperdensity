package claimpolicy

import (
	"path/filepath"
	"sort"
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

func TestValidateDashboardRequiredTokens(t *testing.T) {
	if err := ValidateDashboardRequiredTokens(); err != nil {
		t.Fatal(err)
	}
}

func TestEveryMappingWithFilesHasRequiredTokens(t *testing.T) {
	for _, m := range SurfaceMappings() {
		if len(m.DashboardFiles) == 0 {
			continue
		}
		if len(m.DashboardRequiredTokens) == 0 {
			t.Fatalf("mapping claim=%q surface=%q field=%q has DashboardFiles but no DashboardRequiredTokens", m.ClaimID, m.Surface, m.Field)
		}
	}
}

func TestNoDuplicateTokensPerMappingRow(t *testing.T) {
	for _, m := range SurfaceMappings() {
		seen := make(map[string]struct{}, len(m.DashboardRequiredTokens))
		for _, tok := range m.DashboardRequiredTokens {
			if _, dup := seen[tok]; dup {
				t.Fatalf("duplicate token %q in mapping claim=%q surface=%q field=%q", tok, m.ClaimID, m.Surface, m.Field)
			}
			seen[tok] = struct{}{}
		}
	}
}

func TestRequiredTokensNoPathLikeFragments(t *testing.T) {
	for _, m := range SurfaceMappings() {
		for _, tok := range m.DashboardRequiredTokens {
			if strings.ContainsAny(tok, `/\`) || strings.Contains(tok, "..") {
				t.Fatalf("token %q must not look like a path fragment (mapping claim=%q surface=%q)", tok, m.ClaimID, m.Surface)
			}
		}
	}
}

func TestRequiredTokensForClaimSortedUnique(t *testing.T) {
	for _, rule := range Catalog() {
		toks := RequiredTokensForClaim(rule.ID)
		if len(toks) == 0 {
			t.Fatalf("claim %q must expose required tokens", rule.ID)
		}
		if !sort.StringsAreSorted(toks) {
			t.Fatalf("RequiredTokensForClaim(%q) not sorted: %#v", rule.ID, toks)
		}
		prev := ""
		for _, tok := range toks {
			if tok == prev {
				t.Fatalf("duplicate in RequiredTokensForClaim union for %q", rule.ID)
			}
			prev = tok
		}
	}
}

func TestNoRuntimeContractsImportTokenSet(t *testing.T) {
	var acc []string
	for _, m := range SurfaceMappings() {
		if m.ClaimID != string(ClaimNoRuntimeContractsImport) {
			continue
		}
		acc = append(acc, m.DashboardRequiredTokens...)
	}
	joined := strings.Join(acc, "\n")
	if !strings.Contains(joined, "contractkit") || !strings.Contains(joined, "blockers") {
		t.Fatalf("no_runtime_contracts_import must require contractkit and blockers substrings, got %#v", acc)
	}
}

func TestWindowsLaneDisabledRequiredToken(t *testing.T) {
	for _, m := range SurfaceMappings() {
		if m.ClaimID != string(ClaimWindowsLaneDisabled) {
			continue
		}
		joined := strings.Join(m.DashboardRequiredTokens, "\n")
		if !strings.Contains(joined, "windows_lane_disabled") {
			t.Fatalf("windows_lane_disabled claim mapping must list windows_lane_disabled token, got %#v", m.DashboardRequiredTokens)
		}
	}
}

func TestNoProductionMutationMappingsUseBlockerIDToken(t *testing.T) {
	for _, m := range SurfaceMappings() {
		if m.ClaimID != string(ClaimNoProductionMutation) {
			continue
		}
		joined := strings.Join(m.DashboardRequiredTokens, " ")
		if !strings.Contains(joined, "IDNoProductionMutation") || !strings.Contains(joined, "hpblockers") {
			t.Fatalf("no_production_mutation mapping must trace hpblockers + IDNoProductionMutation, got %#v (surface=%q field=%q)", m.DashboardRequiredTokens, m.Surface, m.Field)
		}
	}
}
