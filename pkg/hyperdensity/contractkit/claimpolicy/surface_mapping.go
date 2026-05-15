// Surface mapping (Sprint 37–39): controlled documentation of how claim-policy catalog IDs
// relate to Karl-Dashboard Parent Fabric surfaces, traced files, and required content tokens.
// Test and doc contract only; no runtime wiring, no Dashboard filesystem reads in Hyperdensity.
package claimpolicy

import (
	"fmt"
	"sort"
	"strings"
)

// ParentFabricSurface names a Parent Fabric / Dashboard builder family (documentation token).
type ParentFabricSurface string

const (
	SurfaceExecutionEngine            ParentFabricSurface = "execution_engine"
	SurfaceWindowsLane                ParentFabricSurface = "windows_lane"
	SurfaceKubeVirtLegacyProvider     ParentFabricSurface = "kubevirt_legacy_provider"
	SurfacePolicyPack                 ParentFabricSurface = "policy_pack"
	SurfaceReleaseSupportMatrix       ParentFabricSurface = "release_support_matrix"
	SurfaceLiveResourceAuthority      ParentFabricSurface = "live_resource_authority"
	SurfaceRuntimeImportFreeze        ParentFabricSurface = "runtime_import_freeze"
	SurfaceHyperdensityRecommendation ParentFabricSurface = "hyperdensity_recommendation"
)

// SurfaceClaimMapping ties a catalog ClaimID to a conceptual Parent Fabric field/value,
// traced Dashboard files, and required substring tokens (Sprint 39). RuntimeImportAllowed is always false.
type SurfaceClaimMapping struct {
	Surface                   string
	ClaimID                   string
	Field                     string
	ExpectedValue             string
	RuntimeImportAllowed      bool
	Notes                     string
	DashboardFiles            []string
	DashboardRequiredTokens   []string
}

var surfaceMappingRows []SurfaceClaimMapping

func normalizeDashboardFiles(files []string) []string {
	seen := make(map[string]struct{}, len(files))
	var out []string
	for _, f := range files {
		f = strings.TrimSpace(f)
		if f == "" {
			continue
		}
		if _, dup := seen[f]; dup {
			continue
		}
		seen[f] = struct{}{}
		out = append(out, f)
	}
	sort.Strings(out)
	return out
}

func normalizeRequiredTokens(tokens []string) []string {
	seen := make(map[string]struct{}, len(tokens))
	var out []string
	for _, t := range tokens {
		t = strings.TrimSpace(t)
		if t == "" {
			continue
		}
		if _, dup := seen[t]; dup {
			continue
		}
		seen[t] = struct{}{}
		out = append(out, t)
	}
	sort.Strings(out)
	return out
}

func init() {
	rows := []SurfaceClaimMapping{
		{
			Surface:              string(SurfaceExecutionEngine),
			ClaimID:              string(ClaimDryRunOnly),
			Field:                "execution_category",
			ExpectedValue:        "dry_run_only",
			RuntimeImportAllowed: false,
			Notes:                "Resource exchange / apply-stage builders when the token denotes execution category (not SupportLevel rows).",
			DashboardFiles: normalizeDashboardFiles([]string{
				"pkg/server/hyperdensity_parent_fabric_release_support_matrix_v1.go",
				"pkg/server/hyperdensity_parent_fabric_resource_exchange_stage_apply.go",
			}),
			DashboardRequiredTokens: normalizeRequiredTokens([]string{"dry_run", "dry_run_only"}),
		},
		{
			Surface:              string(SurfaceExecutionEngine),
			ClaimID:              string(ClaimNoAutonomousApply),
			Field:                "apply_authority",
			ExpectedValue:        "no_autonomous_apply",
			RuntimeImportAllowed: false,
			Notes:                "No unattended or autonomous apply on Hyperdensity Parent Fabric builders.",
			DashboardFiles: normalizeDashboardFiles([]string{
				"pkg/server/hyperdensity_parent_fabric_resource_exchange_stage_apply.go",
				"pkg/server/hyperdensity_parent_fabric_release_support_matrix_v1.go",
				"pkg/server/hyperdensity_parent_fabric_live_resource_authority_v1.go",
			}),
			DashboardRequiredTokens: normalizeRequiredTokens([]string{"autonomous", "no_autonomous_apply"}),
		},
		{
			Surface:              string(SurfaceExecutionEngine),
			ClaimID:              string(ClaimOperatorControlled),
			Field:                "governance_posture",
			ExpectedValue:        "operator_controlled",
			RuntimeImportAllowed: false,
			Notes:                "Submission / grant / approval flows that require explicit operator control.",
			DashboardFiles: normalizeDashboardFiles([]string{
				"pkg/server/hyperdensity_parent_fabric_vm_readonly_observation_operator_submission_v1.go",
				"pkg/server/hyperdensity_parent_fabric_vm_readonly_observation_operator_grant_v1.go",
			}),
			DashboardRequiredTokens: normalizeRequiredTokens([]string{"operator", "Operator"}),
		},
		{
			Surface:              string(SurfaceHyperdensityRecommendation),
			ClaimID:              string(ClaimRecommendationOnly),
			Field:                "surface_mode",
			ExpectedValue:        "recommendation_only",
			RuntimeImportAllowed: false,
			Notes:                "Recommendation and value surfaces without implying live apply or mutation authority.",
			DashboardFiles: normalizeDashboardFiles([]string{
				"pkg/server/hyperdensity_parent_fabric_policy_pack_v1.go",
				"pkg/server/hyperdensity_parent_fabric_release_support_matrix_v1.go",
			}),
			DashboardRequiredTokens: normalizeRequiredTokens([]string{"recommendation", "recommendation_only"}),
		},
		{
			Surface:              string(SurfaceKubeVirtLegacyProvider),
			ClaimID:              string(ClaimKubeVirtLegacyProvider),
			Field:                "provider_marker",
			ExpectedValue:        "kubevirt_legacy_provider",
			RuntimeImportAllowed: false,
			Notes:                "Legacy KubeVirt / VM compatibility markers in live inventory and matrix surfaces; live.go uses kubevirt API paths and centos-legacy probe metadata.",
			DashboardFiles: normalizeDashboardFiles([]string{
				"pkg/server/hyperdensity_parent_fabric_live.go",
				"pkg/server/hyperdensity_parent_fabric_release_support_matrix_v1.go",
			}),
			DashboardRequiredTokens: normalizeRequiredTokens([]string{"kubevirt", "legacy"}),
		},
		{
			Surface:              string(SurfaceKubeVirtLegacyProvider),
			ClaimID:              string(ClaimNoGenericKubeVirtReplacement),
			Field:                "replacement_narrative",
			ExpectedValue:        "forbid_generic_non_kubevirt",
			RuntimeImportAllowed: false,
			Notes:                "Distinct from provider_marker: forbids implying a generic non-KubeVirt replacement for workloads.",
			DashboardFiles: normalizeDashboardFiles([]string{
				"pkg/server/hyperdensity_parent_fabric_live_resource_authority_v1.go",
				"pkg/server/hyperdensity_parent_fabric_vm_readonly_observation_operator_submission_v1.go",
			}),
			DashboardRequiredTokens: normalizeRequiredTokens([]string{"kubevirt", "generic"}),
		},
		{
			Surface:              string(SurfacePolicyPack),
			ClaimID:              string(ClaimNoProductionMutation),
			Field:                "RuleID",
			ExpectedValue:        "no_production_mutation",
			RuntimeImportAllowed: false,
			Notes:                "Policy pack rules that carry the M1-aligned blocker id as RuleID.",
			DashboardFiles: normalizeDashboardFiles([]string{
				"pkg/server/hyperdensity_parent_fabric_policy_pack_v1.go",
				"pkg/server/hyperdensity_parent_fabric_policy_pack_consistency_v1.go",
			}),
			DashboardRequiredTokens: normalizeRequiredTokens([]string{"NoProductionMutation", "no_production_mutation"}),
		},
		{
			Surface:              string(SurfaceReleaseSupportMatrix),
			ClaimID:              string(ClaimNoProductionMutation),
			Field:                "LimitationID",
			ExpectedValue:        "no_production_mutation",
			RuntimeImportAllowed: false,
			Notes:                "Release support matrix limitation rows using the same catalog token.",
			DashboardFiles: normalizeDashboardFiles([]string{
				"pkg/server/hyperdensity_parent_fabric_release_support_matrix_v1.go",
			}),
			DashboardRequiredTokens: normalizeRequiredTokens([]string{"NoProductionMutation", "no_production_mutation"}),
		},
		{
			Surface:              string(SurfaceLiveResourceAuthority),
			ClaimID:              string(ClaimNoProductionMutation),
			Field:                "LimitationID",
			ExpectedValue:        "no_production_mutation",
			RuntimeImportAllowed: false,
			Notes:                "Live resource authority limitation rows use hpblockers.IDNoProductionMutation; the snake-case id is not spelled as a string literal in this file.",
			DashboardFiles: normalizeDashboardFiles([]string{
				"pkg/server/hyperdensity_parent_fabric_live_resource_authority_v1.go",
			}),
			DashboardRequiredTokens: normalizeRequiredTokens([]string{"IDNoProductionMutation", "keep_evidence_namespace_scope"}),
		},
		{
			Surface:              string(SurfaceRuntimeImportFreeze),
			ClaimID:              string(ClaimNoRuntimeContractsImport),
			Field:                "dashboard_pkg_server_non_test_import",
			ExpectedValue:        "contractkit_blockers_only",
			RuntimeImportAllowed: false,
			Notes:                "M17 freeze: only contractkit/blockers in runtime; contracts and claimpolicy test-only.",
			DashboardFiles: normalizeDashboardFiles([]string{
				"scripts/hyperdensity/audit_contractkit_runtime_imports.sh",
				"scripts/hyperdensity/contractkit_runtime_import_allowlist.txt",
			}),
			DashboardRequiredTokens: normalizeRequiredTokens([]string{"contractkit", "blockers"}),
		},
		{
			Surface:              string(SurfaceWindowsLane),
			ClaimID:              string(ClaimNoWindowsHyperdensityApply),
			Field:                "windows_hyperdensity_apply",
			ExpectedValue:        "disabled",
			RuntimeImportAllowed: false,
			Notes:                "Windows Hyperdensity apply remains off; planning-only safety posture.",
			DashboardFiles: normalizeDashboardFiles([]string{
				"pkg/server/hyperdensity_parent_fabric_vm_lane_readiness_v1.go",
				"pkg/server/hyperdensity_parent_fabric_vm_runtime_evidence_collector_v1.go",
			}),
			DashboardRequiredTokens: normalizeRequiredTokens([]string{"windows", "disabled"}),
		},
		{
			Surface:              string(SurfaceWindowsLane),
			ClaimID:              string(ClaimWindowsLaneDisabled),
			Field:                "preflight_check_name",
			ExpectedValue:        "windows_lane_disabled",
			RuntimeImportAllowed: false,
			Notes:                "VM readonly observation preflights use this check name; distinct from catalog windows_disabled.",
			DashboardFiles: normalizeDashboardFiles([]string{
				"pkg/server/hyperdensity_parent_fabric_vm_readonly_observation_real_submission_policy_v1.go",
				"pkg/server/hyperdensity_parent_fabric_vm_readonly_observation_operator_submission_v1.go",
				"pkg/server/hyperdensity_parent_fabric_vm_readonly_observation_operator_grant_v1.go",
				"pkg/server/hyperdensity_parent_fabric_vm_runtime_evidence_collector_v1.go",
			}),
			DashboardRequiredTokens: normalizeRequiredTokens([]string{"windows_lane_disabled"}),
		},
	}
	for i := range rows {
		if len(rows[i].DashboardFiles) == 0 {
			panic("claimpolicy: mapping must include DashboardFiles or document future-only in Notes (Sprint 38)")
		}
		if len(rows[i].DashboardRequiredTokens) == 0 {
			panic("claimpolicy: mapping with DashboardFiles must include DashboardRequiredTokens (Sprint 39)")
		}
	}
	sort.Slice(rows, func(i, j int) bool {
		a, b := rows[i], rows[j]
		if a.ClaimID != b.ClaimID {
			return a.ClaimID < b.ClaimID
		}
		if a.Surface != b.Surface {
			return a.Surface < b.Surface
		}
		if a.Field != b.Field {
			return a.Field < b.Field
		}
		return a.ExpectedValue < b.ExpectedValue
	})
	surfaceMappingRows = rows
}

func cloneMapping(m SurfaceClaimMapping) SurfaceClaimMapping {
	m.DashboardFiles = append([]string(nil), m.DashboardFiles...)
	m.DashboardRequiredTokens = append([]string(nil), m.DashboardRequiredTokens...)
	return m
}

// SurfaceMappings returns all claim ↔ surface rows in stable order.
func SurfaceMappings() []SurfaceClaimMapping {
	out := make([]SurfaceClaimMapping, len(surfaceMappingRows))
	for i := range surfaceMappingRows {
		out[i] = cloneMapping(surfaceMappingRows[i])
	}
	return out
}

// MappingsForClaim returns mappings for a single ClaimID, in stable global order.
func MappingsForClaim(id string) []SurfaceClaimMapping {
	var out []SurfaceClaimMapping
	for _, m := range surfaceMappingRows {
		if m.ClaimID == id {
			out = append(out, cloneMapping(m))
		}
	}
	return out
}

// DashboardFilesForClaim returns the sorted, de-duplicated union of DashboardFiles for a claim id.
func DashboardFilesForClaim(id string) []string {
	seen := make(map[string]struct{})
	var out []string
	for _, m := range surfaceMappingRows {
		if m.ClaimID != id {
			continue
		}
		for _, f := range m.DashboardFiles {
			if _, ok := seen[f]; ok {
				continue
			}
			seen[f] = struct{}{}
			out = append(out, f)
		}
	}
	sort.Strings(out)
	return out
}

// RequiredTokensForClaim returns the sorted, de-duplicated union of DashboardRequiredTokens for a claim id.
func RequiredTokensForClaim(id string) []string {
	seen := make(map[string]struct{})
	var out []string
	for _, m := range surfaceMappingRows {
		if m.ClaimID != id {
			continue
		}
		for _, tok := range m.DashboardRequiredTokens {
			if _, ok := seen[tok]; ok {
				continue
			}
			seen[tok] = struct{}{}
			out = append(out, tok)
		}
	}
	sort.Strings(out)
	return out
}

func validateDashboardPath(p string) error {
	if strings.TrimSpace(p) == "" {
		return fmt.Errorf("dashboard path must be non-empty")
	}
	if strings.HasPrefix(p, "/") || strings.HasPrefix(p, "\\") {
		return fmt.Errorf("dashboard path must be relative, got %q", p)
	}
	if strings.Contains(p, "..") {
		return fmt.Errorf("dashboard path must not contain '..', got %q", p)
	}
	if !(strings.HasPrefix(p, "pkg/server/") || strings.HasPrefix(p, "scripts/hyperdensity/")) {
		return fmt.Errorf("dashboard path %q must start with pkg/server/ or scripts/hyperdensity/", p)
	}
	return nil
}

func validateTokenShape(tok string) error {
	if strings.TrimSpace(tok) == "" {
		return fmt.Errorf("required token must be non-empty")
	}
	if strings.ContainsAny(tok, `/\`) || strings.Contains(tok, "..") {
		return fmt.Errorf("required token %q must not look like a path fragment", tok)
	}
	return nil
}

// validateSurfaceMappingCore checks claim IDs, duplicate mapping keys, and catalog coverage.
func validateSurfaceMappingCore() error {
	seen := make(map[string]struct{}, len(surfaceMappingRows))
	for _, m := range surfaceMappingRows {
		if !Known(m.ClaimID) {
			return fmt.Errorf("surface mapping references unknown claim id %q", m.ClaimID)
		}
		if m.RuntimeImportAllowed {
			return fmt.Errorf("surface mapping for claim %q must have RuntimeImportAllowed=false", m.ClaimID)
		}
		key := m.ClaimID + "\x00" + m.Surface + "\x00" + m.Field + "\x00" + m.ExpectedValue
		if _, dup := seen[key]; dup {
			return fmt.Errorf("duplicate surface mapping key for claim %q surface %q field %q value %q", m.ClaimID, m.Surface, m.Field, m.ExpectedValue)
		}
		seen[key] = struct{}{}
	}
	for _, rule := range Catalog() {
		if len(MappingsForClaim(rule.ID)) == 0 {
			return fmt.Errorf("catalog claim %q has no surface mapping", rule.ID)
		}
	}
	return nil
}

func validateDashboardFilePaths() error {
	for _, m := range surfaceMappingRows {
		if len(m.DashboardFiles) == 0 {
			if strings.Contains(strings.ToLower(m.Notes), "future-only") {
				continue
			}
			return fmt.Errorf("mapping claim=%q surface=%q field=%q has empty DashboardFiles (document future-only in Notes or add paths)", m.ClaimID, m.Surface, m.Field)
		}
		rowSeen := make(map[string]struct{}, len(m.DashboardFiles))
		for _, p := range m.DashboardFiles {
			if err := validateDashboardPath(p); err != nil {
				return fmt.Errorf("claim %q surface %q: %w", m.ClaimID, m.Surface, err)
			}
			if _, dup := rowSeen[p]; dup {
				return fmt.Errorf("duplicate dashboard path %q within mapping claim=%q surface=%q field=%q", p, m.ClaimID, m.Surface, m.Field)
			}
			rowSeen[p] = struct{}{}
		}
	}
	for _, rule := range Catalog() {
		if len(DashboardFilesForClaim(rule.ID)) == 0 {
			return fmt.Errorf("catalog claim %q has no traced dashboard files", rule.ID)
		}
	}
	return nil
}

// validateDashboardRequiredTokensBody validates token rows (caller supplies core + path checks as needed).
func validateDashboardRequiredTokensBody() error {
	for _, m := range surfaceMappingRows {
		if len(m.DashboardFiles) == 0 {
			continue
		}
		if len(m.DashboardRequiredTokens) == 0 {
			return fmt.Errorf("mapping claim=%q surface=%q field=%q has DashboardFiles but empty DashboardRequiredTokens", m.ClaimID, m.Surface, m.Field)
		}
		rowTok := make(map[string]struct{}, len(m.DashboardRequiredTokens))
		for _, tok := range m.DashboardRequiredTokens {
			if err := validateTokenShape(tok); err != nil {
				return fmt.Errorf("claim %q surface %q: %w", m.ClaimID, m.Surface, err)
			}
			if _, dup := rowTok[tok]; dup {
				return fmt.Errorf("duplicate required token %q within mapping claim=%q surface=%q field=%q", tok, m.ClaimID, m.Surface, m.Field)
			}
			rowTok[tok] = struct{}{}
		}
	}
	for _, rule := range Catalog() {
		if len(RequiredTokensForClaim(rule.ID)) == 0 {
			return fmt.Errorf("catalog claim %q has no required dashboard tokens", rule.ID)
		}
	}
	return nil
}

// ValidateDashboardRequiredTokens checks Sprint 39 token invariants (no filesystem I/O).
func ValidateDashboardRequiredTokens() error {
	if err := validateSurfaceMappingCore(); err != nil {
		return err
	}
	return validateDashboardRequiredTokensBody()
}

// ValidateSurfaceMappings checks full mapping contract (Sprint 39: file paths + required tokens).
func ValidateSurfaceMappings() error {
	if err := validateSurfaceMappingCore(); err != nil {
		return err
	}
	if err := validateDashboardFilePaths(); err != nil {
		return err
	}
	return validateDashboardRequiredTokensBody()
}

// ValidateDashboardFileTraceability checks Sprint 38–39 path + token contract (no filesystem reads here).
func ValidateDashboardFileTraceability() error {
	return ValidateSurfaceMappings()
}
