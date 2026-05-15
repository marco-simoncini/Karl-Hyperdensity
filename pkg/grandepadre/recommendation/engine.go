package recommendation

import (
	"strings"
	"time"

	gpevidence "github.com/marco-simoncini/Karl-Hyperdensity/pkg/grandepadre/evidence"
)

// EngineOptions configures the local recommendation skeleton.
type EngineOptions struct {
	Tenant        string
	DryRunOnly    bool
	UnsignedLabel gpevidence.UnsignedDigestTrustPolicy
	GeneratedAt   time.Time
}

// Now returns the clock used for generatedAt (tests inject via EngineOptions.GeneratedAt).
func (o EngineOptions) now() time.Time {
	if !o.GeneratedAt.IsZero() {
		return o.GeneratedAt.UTC()
	}
	return gpevidence.NowFunc().UTC()
}

// FilterIndicesByTenant keeps rows whose cell namespace matches tenant when tenant is non-empty.
func FilterIndicesByTenant(indices []gpevidence.EvidenceIndex, tenant string) []gpevidence.EvidenceIndex {
	t := strings.TrimSpace(tenant)
	if t == "" {
		return indices
	}
	var out []gpevidence.EvidenceIndex
	for _, idx := range indices {
		if idx.CellRef != nil && strings.TrimSpace(idx.CellRef.Namespace) == t {
			out = append(out, idx)
		}
	}
	return out
}

func inferResource(idx gpevidence.EvidenceIndex) Resource {
	for _, b := range idx.BlockedReasons {
		lb := strings.ToLower(b)
		if strings.Contains(lb, "memory") || strings.Contains(lb, "mem") {
			return ResourceMemory
		}
		if strings.Contains(lb, "cpu") {
			return ResourceCPU
		}
	}
	for _, w := range idx.Warnings {
		lw := strings.ToLower(w)
		if strings.Contains(lw, "memory") {
			return ResourceMemory
		}
		if strings.Contains(lw, "cpu") {
			return ResourceCPU
		}
	}
	return ResourceUnknown
}

func envelopeMode(idx gpevidence.EvidenceIndex) EnvelopeMode {
	for _, s := range append(append([]string{}, idx.BlockedReasons...), idx.Warnings...) {
		ls := strings.ToLower(s)
		if strings.Contains(ls, "cgroup") || strings.Contains(ls, "envelope") {
			return ModeEnvelope
		}
	}
	if idx.ReadyForGrandePadre {
		return ModeEnvelope
	}
	return ModeUnknown
}

// BuildRecommendations turns evidence rows into dry-run ActionRecommendations.
func BuildRecommendations(indices []gpevidence.EvidenceIndex, donors, receivers []gpevidence.CellRefLite, o EngineOptions) []ActionRecommendation {
	var recs []ActionRecommendation
	n := 0
	for _, idx := range indices {
		prefix := strings.TrimSpace(idx.ArtifactID)
		if prefix == "" {
			prefix = "artifact"
		}

		if idx.TrustTier == gpevidence.TrustIntegrityFailed {
			n++
			recs = append(recs, ActionRecommendation{
				ActionID:      actionID(prefix+"-integrity", n),
				ActionType:    ActionRemediate,
				TargetCellRef: cloneRef(idx.CellRef),
				Resource:      ResourceUnknown,
				Mode:          ModeUnknown,
				Confidence:    idx.Confidence,
				Risk:          RiskBlocked,
				Priority:      PriorityHigh,
				Reasons: []string{
					"Bundle digest or manifest alignment failed (IntegrityFailed).",
					"Do not use this evidence for production promotion until integrity is restored.",
				},
				Prerequisites: []string{"Regenerate manifest/digest from canonical bundle JSON.", "Future KHR apply gate must reject mutations on failed integrity."},
				DryRunOnly:    o.DryRunOnly,
				ApplyAllowed:  false,
			})
			continue
		}

		if idx.TrustTier == gpevidence.TrustDevOnly {
			n++
			recs = append(recs, ActionRecommendation{
				ActionID:      actionID(prefix+"-devonly", n),
				ActionType:    ActionObserve,
				TargetCellRef: cloneRef(idx.CellRef),
				Resource:      inferResource(idx),
				Mode:          envelopeMode(idx),
				Confidence:    idx.Confidence,
				Risk:          RiskMedium,
				Priority:      PriorityMedium,
				Reasons: []string{
					"DevOnly trust tier: local-dev integrity only; not production PKI.",
					"Recommendations are non-production; no apply.",
				},
				Prerequisites: []string{"Do not treat DevOnly-signed bundles as production donors."},
				DryRunOnly:    o.DryRunOnly,
				ApplyAllowed:  false,
			})
			continue
		}

		if strings.ToLower(strings.TrimSpace(idx.Confidence)) == "low" {
			n++
			recs = append(recs, ActionRecommendation{
				ActionID:      actionID(prefix+"-lowconf", n),
				ActionType:    ActionCollectMoreEvidence,
				TargetCellRef: cloneRef(idx.CellRef),
				Resource:      inferResource(idx),
				Mode:          envelopeMode(idx),
				Confidence:    idx.Confidence,
				Risk:          RiskMedium,
				Priority:      PriorityHigh,
				Reasons:       []string{"Confidence is low; collect additional read-only evidence before market-style placement."},
				Prerequisites: []string{"Re-run collect-evidence with stable clocks and complete lease/port inputs if applicable."},
				DryRunOnly:    o.DryRunOnly,
				ApplyAllowed:  false,
			})
		}

		if idx.ReadyForGrandePadre && len(idx.BlockedReasons) == 0 &&
			strings.ToLower(strings.TrimSpace(idx.Confidence)) == "high" &&
			(idx.TrustTier == gpevidence.TrustIntegrityVerified || idx.TrustTier == gpevidence.TrustUnsigned) {
			n++
			var donor *gpevidence.CellRefLite
			for _, d := range donors {
				if idx.CellRef != nil && d.Namespace == idx.CellRef.Namespace && d.Name == idx.CellRef.Name {
					donor = cloneRef(&d)
					break
				}
			}
			recs = append(recs, ActionRecommendation{
				ActionID:      actionID(prefix+"-observe", n),
				ActionType:    ActionObserve,
				TargetCellRef: cloneRef(idx.CellRef),
				DonorCellRef:  donor,
				Resource:      inferResource(idx),
				Mode:          envelopeMode(idx),
				Confidence:    idx.Confidence,
				Risk:          RiskLow,
				Priority:      PriorityMedium,
				Reasons:       []string{"Cell evidence is ready at high confidence; continue read-only observation for predictive market signals."},
				Prerequisites: []string{
					"No apply: await future KHR apply gate and Hyperdensity controller decisions.",
				},
				DryRunOnly:   o.DryRunOnly,
				ApplyAllowed: false,
			})
			n++
			recs = append(recs, ActionRecommendation{
				ActionID:      actionID(prefix+"-prep-lease", n),
				ActionType:    ActionPrepareResourceLease,
				TargetCellRef: cloneRef(idx.CellRef),
				DonorCellRef:  donor,
				Resource:      ResourceCPU,
				Mode:          ModeEnvelope,
				Confidence:    idx.Confidence,
				Risk:          RiskLow,
				Priority:      PriorityLow,
				Reasons:       []string{"Dry-run skeleton: prepare ResourceLease drafts offline before any future apply window."},
				Prerequisites: []string{
					"Use existing khr-linux-agent dry-run mode with lease/port fixtures.",
					"applyAllowed remains false in Sprint 13.",
				},
				DryRunOnly:   o.DryRunOnly,
				ApplyAllowed: false,
			})
			continue
		}

		if isBlockedLike(idx) && idx.TrustTier != gpevidence.TrustIntegrityFailed {
			n++
			var recv *gpevidence.CellRefLite
			for _, r := range receivers {
				if idx.CellRef != nil && r.Namespace == idx.CellRef.Namespace && r.Name == idx.CellRef.Name {
					recv = cloneRef(&r)
					break
				}
			}
			recs = append(recs, ActionRecommendation{
				ActionID:        actionID(prefix+"-remediate", n),
				ActionType:      ActionRemediate,
				TargetCellRef:   cloneRef(idx.CellRef),
				ReceiverCellRef: recv,
				Resource:        inferResource(idx),
				Mode:            envelopeMode(idx),
				Confidence:      idx.Confidence,
				Risk:            RiskHigh,
				Priority:        PriorityHigh,
				Reasons:         append([]string{"Evidence is blocked or not ready for Grande Padre."}, idx.BlockedReasons...),
				Prerequisites:   []string{"Resolve blockedReasons locally; ingest is not apply."},
				DryRunOnly:      o.DryRunOnly,
				ApplyAllowed:    false,
			})
		}
	}
	return recs
}

func cloneRef(c *gpevidence.CellRefLite) *gpevidence.CellRefLite {
	if c == nil {
		return nil
	}
	cp := *c
	return &cp
}

// BuildActionSlate assembles the full slate from an in-memory evidence store.
func BuildActionSlate(s *gpevidence.Store, o EngineOptions) ActionSlate {
	s.DeduplicateBySha256()
	filtered := FilterIndicesByTenant(s.Snapshot(), o.Tenant)
	var blockedRows []gpevidence.EvidenceIndex
	for _, idx := range filtered {
		if isBlockedLike(idx) {
			blockedRows = append(blockedRows, idx)
		}
	}
	rem := gpevidence.BuildBlockedRemediableIndex(filtered)
	donors := DonorCandidates(filtered)
	receivers := ReceiverCandidates(filtered)
	recs := BuildRecommendations(filtered, donors, receivers, o)
	slate := ActionSlate{
		GeneratedAt:        o.now().Format(time.RFC3339),
		Source:             "local-evidence-store",
		Recommendations:    recs,
		Blocked:            nonNilIndexSlice(blockedRows),
		Remediable:         nonNilRemed(rem),
		DonorCandidates:    donors,
		ReceiverCandidates: receivers,
	}
	summarizeSlate(&slate)
	return slate
}

func nonNilIndexSlice(v []gpevidence.EvidenceIndex) []gpevidence.EvidenceIndex {
	if v == nil {
		return []gpevidence.EvidenceIndex{}
	}
	return v
}

func nonNilRemed(v []gpevidence.BlockedRemediableIndex) []gpevidence.BlockedRemediableIndex {
	if v == nil {
		return []gpevidence.BlockedRemediableIndex{}
	}
	return v
}
