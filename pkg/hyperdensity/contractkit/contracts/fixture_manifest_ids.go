package contracts

// ExpectedM1M8CaseIDs returns the stable M1–M8 parity case IDs for manifest drift guards.
func ExpectedM1M8CaseIDs() []string {
	return []string{
		"m1-blocker-catalog",
		"m2-apply-semantics",
		"m3-mapper",
		"m4-golden-generator",
		"m5-live-redacted-supports-apply-true",
		"m6-live-redacted-supports-apply-false",
		"m7-missing-optional",
		"m8-contractkit-import",
	}
}

// CaseIDs returns case IDs from manifest in declaration order (stable for a given JSON file).
func CaseIDs(m FixtureManifest) []string {
	out := make([]string, 0, len(m.Cases))
	for _, c := range m.Cases {
		out = append(out, c.ID)
	}
	return out
}

// ManifestCaseIDSet returns a set of case IDs for membership checks.
func ManifestCaseIDSet(m FixtureManifest) map[string]struct{} {
	set := make(map[string]struct{}, len(m.Cases))
	for _, c := range m.Cases {
		set[c.ID] = struct{}{}
	}
	return set
}
