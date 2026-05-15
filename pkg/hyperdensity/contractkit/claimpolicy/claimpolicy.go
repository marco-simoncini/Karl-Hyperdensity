// Package claimpolicy defines stable claim / policy posture vocabulary, a minimal
// claim-policy catalog (Sprint 35–36), Parent Fabric surface mapping (Sprint 37), and
// Dashboard file traceability metadata (Sprint 38–39: required content tokens).
//
// Stdlib-only; no cluster or HTTP I/O. Karl-Dashboard production code must not import
// this package — see docs/extraction/HYPERDENSITY_CONTRACTKIT_CLAIMPOLICY.md.
package claimpolicy

import "sort"

// PackageVersion is the claimpolicy design epoch (not ContractKitSchemaVersion / manifest epoch).
const PackageVersion = "v0.0.0-sprint39"

// PostureKind describes how narrowly a governance surface may assert readiness or mutation posture.
// Tokens align with existing parent-fabric vocabulary without binding JSON field names here.
type PostureKind string

const (
	PostureEvidenceNamespace  PostureKind = "evidence_namespace_only"
	PostureVisibilityOnly     PostureKind = "visibility_only"
	PostureOperatorControlled PostureKind = "operator_controlled_only"
)

var knownPostures = map[PostureKind]struct{}{
	PostureEvidenceNamespace:  {},
	PostureVisibilityOnly:     {},
	PostureOperatorControlled: {},
}

// KnownPosture reports whether k is a defined posture token.
func KnownPosture(k PostureKind) bool {
	_, ok := knownPostures[k]
	return ok
}

// Postures returns all defined posture kinds in stable sort order (for tests and docs).
func Postures() []PostureKind {
	return []PostureKind{
		PostureEvidenceNamespace,
		PostureOperatorControlled,
		PostureVisibilityOnly,
	}
}

// ClaimPolicyID is a stable catalog identifier for a claim-policy rule (string alias).
type ClaimPolicyID string

// Canonical claim-policy IDs (minimal catalog; Sprint 36). Values are stable wire/doc tokens.
const (
	ClaimDryRunOnly                  ClaimPolicyID = "dry_run_only"
	ClaimKubeVirtLegacyProvider      ClaimPolicyID = "kubevirt_legacy_provider"
	ClaimNoAutonomousApply           ClaimPolicyID = "no_autonomous_apply"
	ClaimNoGenericKubeVirtReplacement ClaimPolicyID = "no_generic_kubevirt_replacement"
	ClaimNoProductionMutation        ClaimPolicyID = "no_production_mutation"
	ClaimNoRuntimeContractsImport    ClaimPolicyID = "no_runtime_contracts_import"
	ClaimNoWindowsHyperdensityApply  ClaimPolicyID = "no_windows_hyperdensity_apply"
	ClaimOperatorControlled          ClaimPolicyID = "operator_controlled"
	ClaimRecommendationOnly          ClaimPolicyID = "recommendation_only"
	ClaimWindowsLaneDisabled         ClaimPolicyID = "windows_lane_disabled"
)

// ClaimPolicyRule is one row of the claim-policy catalog (documentation / parity; no runtime enforcement).
type ClaimPolicyRule struct {
	ID             string
	Severity       string
	RuntimeAllowed bool
	Description    string
}

// Severity levels (aligned with contractkit/blockers style; stable strings).
const (
	severityCritical = "critical"
	severityHigh     = "high"
	severityInfo     = "info"
	severityMedium   = "medium"
)

var (
	catalogRules []ClaimPolicyRule
	byID         map[string]ClaimPolicyRule
)

func init() {
	// Stable lexicographic order by ID (catalog order must remain deterministic).
	catalogRules = []ClaimPolicyRule{
		{
			ID:             string(ClaimDryRunOnly),
			Severity:       severityHigh,
			RuntimeAllowed: true,
			Description:    "Execution and apply surfaces remain dry-run or planning-only; no autonomous apply.",
		},
		{
			ID:             string(ClaimKubeVirtLegacyProvider),
			Severity:       severityInfo,
			RuntimeAllowed: true,
			Description:    "KubeVirt may appear as a legacy / compatibility provider marker; not a generic replacement claim.",
		},
		{
			ID:             string(ClaimNoAutonomousApply),
			Severity:       severityCritical,
			RuntimeAllowed: false,
			Description:    "Autonomous or unattended apply is forbidden on the Hyperdensity surface.",
		},
		{
			ID:             string(ClaimNoGenericKubeVirtReplacement),
			Severity:       severityCritical,
			RuntimeAllowed: false,
			Description:    "No claim that a non-KubeVirt path generically replaces KubeVirt workloads.",
		},
		{
			ID:             string(ClaimNoProductionMutation),
			Severity:       severityCritical,
			RuntimeAllowed: false,
			Description:    "Production mutation is not allowed on the active Hyperdensity surface.",
		},
		{
			ID:             string(ClaimNoRuntimeContractsImport),
			Severity:       severityCritical,
			RuntimeAllowed: false,
			Description:    "contractkit/contracts must not be imported from Dashboard runtime production code (M17 freeze).",
		},
		{
			ID:             string(ClaimNoWindowsHyperdensityApply),
			Severity:       severityCritical,
			RuntimeAllowed: false,
			Description:    "Windows Hyperdensity apply paths remain disabled; planning-only posture.",
		},
		{
			ID:             string(ClaimOperatorControlled),
			Severity:       severityMedium,
			RuntimeAllowed: true,
			Description:    "Material changes require explicit operator control and gates.",
		},
		{
			ID:             string(ClaimRecommendationOnly),
			Severity:       severityMedium,
			RuntimeAllowed: true,
			Description:    "Surface may emit recommendations without implying apply or mutation authority.",
		},
		{
			ID:             string(ClaimWindowsLaneDisabled),
			Severity:       severityHigh,
			RuntimeAllowed: false,
			Description:    "Windows lane remains disabled; aligns with planning-only / no Windows enablement posture.",
		},
	}
	byID = make(map[string]ClaimPolicyRule, len(catalogRules))
	for _, r := range catalogRules {
		if _, dup := byID[r.ID]; dup {
			panic("claimpolicy: duplicate catalog id " + r.ID)
		}
		byID[r.ID] = r
	}
}

// Catalog returns all claim-policy rules in stable order (lexicographic by ID).
func Catalog() []ClaimPolicyRule {
	out := make([]ClaimPolicyRule, len(catalogRules))
	copy(out, catalogRules)
	return out
}

// Known reports whether id is a defined claim-policy catalog identifier.
func Known(id string) bool {
	_, ok := byID[id]
	return ok
}

// Severity returns the catalog severity for id, or empty if unknown.
func Severity(id string) string {
	if r, ok := byID[id]; ok {
		return r.Severity
	}
	return ""
}

// RuntimeAllowed reports whether the catalog marks this claim as compatible with affirmed production posture
// for documentation / parity (not runtime enforcement).
func RuntimeAllowed(id string) bool {
	if r, ok := byID[id]; ok {
		return r.RuntimeAllowed
	}
	return false
}

// MustKeepRuntimeDisabled reports whether this claim implies Hyperdensity runtime must keep
// the named capability disabled (critical freeze anchors; Sprint 36 tests).
func MustKeepRuntimeDisabled(id string) bool {
	switch ClaimPolicyID(id) {
	case ClaimNoProductionMutation, ClaimNoAutonomousApply, ClaimNoWindowsHyperdensityApply, ClaimNoRuntimeContractsImport:
		return true
	default:
		return false
	}
}

// ForbiddenProductionClaimIDs returns claim IDs that label operations or narratives forbidden
// in production Hyperdensity contexts (stable sort order). Excludes kubevirt_legacy_provider,
// which is a compatibility marker, not a replacement-forbidden claim.
func ForbiddenProductionClaimIDs() []string {
	ids := []string{
		string(ClaimNoAutonomousApply),
		string(ClaimNoGenericKubeVirtReplacement),
		string(ClaimNoProductionMutation),
		string(ClaimNoRuntimeContractsImport),
		string(ClaimNoWindowsHyperdensityApply),
	}
	sort.Strings(ids)
	return ids
}
