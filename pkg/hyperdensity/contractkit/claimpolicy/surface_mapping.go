// Surface mapping (Sprint 37–38): controlled documentation of how claim-policy catalog IDs
// relate to Karl-Dashboard Parent Fabric conceptual surfaces and traced Dashboard files.
// Test and doc contract only; no runtime wiring and no Dashboard imports from this package.
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

// SurfaceClaimMapping ties a catalog ClaimID to a conceptual Parent Fabric field/value and
// Karl-Dashboard file paths (relative to kubernetes-console/). Sprint 37–38: RuntimeImportAllowed is always false.
type SurfaceClaimMapping struct {
	Surface              string
	ClaimID              string
	Field                string
	ExpectedValue        string
	RuntimeImportAllowed bool
	Notes                string
	DashboardFiles       []string
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
		},
		{
			Surface:              string(SurfaceKubeVirtLegacyProvider),
			ClaimID:              string(ClaimKubeVirtLegacyProvider),
			Field:                "provider_marker",
			ExpectedValue:        "kubevirt_legacy_provider",
			RuntimeImportAllowed: false,
			Notes:                "Legacy KubeVirt / VM compatibility markers in live inventory and matrix surfaces.",
			DashboardFiles: normalizeDashboardFiles([]string{
				"pkg/server/hyperdensity_parent_fabric_live.go",
				"pkg/server/hyperdensity_parent_fabric_release_support_matrix_v1.go",
			}),
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
		},
		{
			Surface:              string(SurfaceLiveResourceAuthority),
			ClaimID:              string(ClaimNoProductionMutation),
			Field:                "LimitationID",
			ExpectedValue:        "no_production_mutation",
			RuntimeImportAllowed: false,
			Notes:                "Live resource authority limitation rows reference the same production-mutation gate id.",
			DashboardFiles: normalizeDashboardFiles([]string{
				"pkg/server/hyperdensity_parent_fabric_live_resource_authority_v1.go",
			}),
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
		},
	}
	for i := range rows {
		if len(rows[i].DashboardFiles) == 0 {
			panic("claimpolicy: mapping must include DashboardFiles or document future-only in Notes (Sprint 38)")
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

// SurfaceMappings returns all claim ↔ surface rows in stable order.
func SurfaceMappings() []SurfaceClaimMapping {
	out := make([]SurfaceClaimMapping, len(surfaceMappingRows))
	for i := range surfaceMappingRows {
		m := surfaceMappingRows[i]
		m.DashboardFiles = append([]string(nil), m.DashboardFiles...)
		out[i] = m
	}
	return out
}

// MappingsForClaim returns mappings for a single ClaimID, in stable global order.
func MappingsForClaim(id string) []SurfaceClaimMapping {
	var out []SurfaceClaimMapping
	for _, m := range surfaceMappingRows {
		if m.ClaimID == id {
			mc := m
			mc.DashboardFiles = append([]string(nil), m.DashboardFiles...)
			out = append(out, mc)
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

// ValidateSurfaceMappings checks invariants for Sprint 37 mapping tables.
func ValidateSurfaceMappings() error {
	seen := make(map[string]struct{}, len(surfaceMappingRows))
	for _, m := range surfaceMappingRows {
		if !Known(m.ClaimID) {
			return fmt.Errorf("surface mapping references unknown claim id %q", m.ClaimID)
		}
		if m.RuntimeImportAllowed {
			return fmt.Errorf("surface mapping for claim %q must have RuntimeImportAllowed=false in Sprint 37", m.ClaimID)
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

// ValidateDashboardFileTraceability checks Sprint 38 file path invariants for all mappings.
func ValidateDashboardFileTraceability() error {
	if err := ValidateSurfaceMappings(); err != nil {
		return err
	}
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
