// Package claimpolicy defines stable claim / policy posture vocabulary for Parent Fabric
// parity and future Karl-Dashboard mapping (Sprint 35 boundary design).
//
// Stdlib-only; no cluster or HTTP I/O. Dashboard production code must not import this
// package until a dedicated extraction sprint — see docs/extraction/HYPERDENSITY_CONTRACTKIT_CLAIMPOLICY.md.
package claimpolicy

// PackageVersion is the claimpolicy design epoch (not ContractKitSchemaVersion / manifest epoch).
const PackageVersion = "v0.0.0-sprint35"

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
