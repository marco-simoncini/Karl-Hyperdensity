package evidence

import (
	"sort"
	"strings"
)

// BuildBlockedRemediableIndex aggregates blocked / not-ready rows by cellRef.
func BuildBlockedRemediableIndex(indices []EvidenceIndex) []BlockedRemediableIndex {
	type key struct {
		ns, name string
	}
	groups := map[key][]EvidenceIndex{}
	for _, idx := range indices {
		if !isBlockedLike(idx) {
			continue
		}
		if idx.CellRef == nil {
			k := key{"_", "_"}
			groups[k] = append(groups[k], idx)
			continue
		}
		k := key{strings.TrimSpace(idx.CellRef.Namespace), strings.TrimSpace(idx.CellRef.Name)}
		groups[k] = append(groups[k], idx)
	}
	var keys []key
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].ns == keys[j].ns {
			return keys[i].name < keys[j].name
		}
		return keys[i].ns < keys[j].ns
	})
	out := make([]BlockedRemediableIndex, 0, len(keys))
	for _, k := range keys {
		rows := groups[k]
		sort.Slice(rows, func(i, j int) bool { return rows[i].IndexedAt < rows[j].IndexedAt })
		last := rows[len(rows)-1]
		var reasons []string
		seen := map[string]bool{}
		for _, r := range rows {
			for _, b := range r.BlockedReasons {
				if !seen[b] {
					seen[b] = true
					reasons = append(reasons, b)
				}
			}
		}
		var cell *CellRefLite
		if last.CellRef != nil {
			cell = &CellRefLite{
				APIVersion: last.CellRef.APIVersion,
				Kind:       last.CellRef.Kind,
				Namespace:  last.CellRef.Namespace,
				Name:       last.CellRef.Name,
			}
		}
		out = append(out, BlockedRemediableIndex{
			CellRef:          cell,
			BlockedReasons:   reasons,
			RemediationHints: remediationHintsFor(reasons),
			LastArtifactID:   last.ArtifactID,
			Confidence:       last.Confidence,
			TrustTier:        last.TrustTier,
		})
	}
	return out
}

func isBlockedLike(idx EvidenceIndex) bool {
	return !idx.ReadyForGrandePadre || len(idx.BlockedReasons) > 0 || idx.TrustTier == TrustIntegrityFailed
}

func remediationHintsFor(blocked []string) []string {
	var out []string
	seen := map[string]bool{}
	for _, b := range blocked {
		lb := strings.ToLower(b)
		var hint string
		switch {
		case strings.Contains(lb, "cgroup"):
			hint = "Re-run discovery with a valid cgroup root and path policy, then collect evidence again."
		case strings.Contains(lb, "lease") || strings.Contains(lb, "resource"):
			hint = "Provide ResourceLease / ResourcePort inputs or accept dry-run skip before expecting readiness."
		case strings.Contains(lb, "digest") || strings.Contains(lb, "integrity"):
			hint = "Regenerate manifest and digest from canonical bundle JSON; digest match is integrity, not authorization."
		default:
			hint = "Review blockedReasons on the cell and remediate before promoting beyond dry-run indexing."
		}
		if !seen[hint] {
			seen[hint] = true
			out = append(out, hint)
		}
	}
	return out
}
