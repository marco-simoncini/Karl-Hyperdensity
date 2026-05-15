// Surface mapping (Sprint 37): controlled documentation of how claim-policy catalog IDs
// relate to Karl-Dashboard Parent Fabric conceptual surfaces. Test and doc contract only;
// no runtime wiring and no Dashboard imports from this package.
package claimpolicy

import (
	"fmt"
	"sort"
)

// ParentFabricSurface names a Parent Fabric / Dashboard builder family (documentation token).
type ParentFabricSurface string

const (
	SurfaceExecutionEngine            ParentFabricSurface = "execution_engine"
	SurfaceWindowsLane                ParentFabricSurface = "windows_lane"
	SurfaceKubeVirtLegacyProvider     ParentFabricSurface = "kubevirt_legacy_provider"
	SurfacePolicyPack                 ParentFabricSurface = "policy_pack"
	SurfaceReleaseSupportMatrix       ParentFabricSurface = "release_support_matrix"
	SurfaceRuntimeImportFreeze        ParentFabricSurface = "runtime_import_freeze"
	SurfaceHyperdensityRecommendation ParentFabricSurface = "hyperdensity_recommendation"
)

// SurfaceClaimMapping ties a catalog ClaimID to a conceptual Parent Fabric field/value.
// Sprint 37: RuntimeImportAllowed is always false (claimpolicy remains test-only on Dashboard).
type SurfaceClaimMapping struct {
	Surface              string
	ClaimID              string
	Field                string
	ExpectedValue        string
	RuntimeImportAllowed bool
	Notes                string
}

var surfaceMappingRows []SurfaceClaimMapping

func init() {
	rows := []SurfaceClaimMapping{
		{
			Surface:              string(SurfaceExecutionEngine),
			ClaimID:              string(ClaimDryRunOnly),
			Field:                "execution_category",
			ExpectedValue:        "dry_run_only",
			RuntimeImportAllowed: false,
			Notes:                "Resource exchange / apply-stage builders when the token denotes execution category (not SupportLevel rows).",
		},
		{
			Surface:              string(SurfaceExecutionEngine),
			ClaimID:              string(ClaimNoAutonomousApply),
			Field:                "apply_authority",
			ExpectedValue:        "no_autonomous_apply",
			RuntimeImportAllowed: false,
			Notes:                "No unattended or autonomous apply on Hyperdensity Parent Fabric builders.",
		},
		{
			Surface:              string(SurfaceExecutionEngine),
			ClaimID:              string(ClaimOperatorControlled),
			Field:                "governance_posture",
			ExpectedValue:        "operator_controlled",
			RuntimeImportAllowed: false,
			Notes:                "Submission / grant / approval flows that require explicit operator control.",
		},
		{
			Surface:              string(SurfaceHyperdensityRecommendation),
			ClaimID:              string(ClaimRecommendationOnly),
			Field:                "surface_mode",
			ExpectedValue:        "recommendation_only",
			RuntimeImportAllowed: false,
			Notes:                "Recommendation and value surfaces without implying live apply or mutation authority.",
		},
		{
			Surface:              string(SurfaceKubeVirtLegacyProvider),
			ClaimID:              string(ClaimKubeVirtLegacyProvider),
			Field:                "provider_marker",
			ExpectedValue:        "kubevirt_legacy_provider",
			RuntimeImportAllowed: false,
			Notes:                "Legacy KubeVirt runtimeprovider / VM inventory markers; compatibility only.",
		},
		{
			Surface:              string(SurfaceKubeVirtLegacyProvider),
			ClaimID:              string(ClaimNoGenericKubeVirtReplacement),
			Field:                "replacement_narrative",
			ExpectedValue:        "forbid_generic_non_kubevirt",
			RuntimeImportAllowed: false,
			Notes:                "Distinct from provider_marker: forbids implying a generic non-KubeVirt replacement for workloads.",
		},
		{
			Surface:              string(SurfacePolicyPack),
			ClaimID:              string(ClaimNoProductionMutation),
			Field:                "RuleID",
			ExpectedValue:        "no_production_mutation",
			RuntimeImportAllowed: false,
			Notes:                "Policy pack rules that carry the M1-aligned blocker id as RuleID.",
		},
		{
			Surface:              string(SurfaceReleaseSupportMatrix),
			ClaimID:              string(ClaimNoProductionMutation),
			Field:                "LimitationID",
			ExpectedValue:        "no_production_mutation",
			RuntimeImportAllowed: false,
			Notes:                "Release support matrix limitation rows using the same catalog token.",
		},
		{
			Surface:              string(SurfaceRuntimeImportFreeze),
			ClaimID:              string(ClaimNoRuntimeContractsImport),
			Field:                "dashboard_pkg_server_non_test_import",
			ExpectedValue:        "contractkit_blockers_only",
			RuntimeImportAllowed: false,
			Notes:                "M17 freeze: only contractkit/blockers in runtime; contracts and claimpolicy test-only.",
		},
		{
			Surface:              string(SurfaceWindowsLane),
			ClaimID:              string(ClaimNoWindowsHyperdensityApply),
			Field:                "windows_hyperdensity_apply",
			ExpectedValue:        "disabled",
			RuntimeImportAllowed: false,
			Notes:                "Windows Hyperdensity apply remains off; planning-only safety posture.",
		},
		{
			Surface:              string(SurfaceWindowsLane),
			ClaimID:              string(ClaimWindowsLaneDisabled),
			Field:                "preflight_check_name",
			ExpectedValue:        "windows_lane_disabled",
			RuntimeImportAllowed: false,
			Notes:                "VM readonly observation preflights use this check name; distinct from catalog windows_disabled.",
		},
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
	copy(out, surfaceMappingRows)
	return out
}

// MappingsForClaim returns mappings for a single ClaimID, in stable global order.
func MappingsForClaim(id string) []SurfaceClaimMapping {
	var out []SurfaceClaimMapping
	for _, m := range surfaceMappingRows {
		if m.ClaimID == id {
			out = append(out, m)
		}
	}
	return out
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
