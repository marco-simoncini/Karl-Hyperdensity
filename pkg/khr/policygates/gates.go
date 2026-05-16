// Package policygates evaluates read-only certification policy gates (KHR-V).
package policygates

import (
	"strings"
	"time"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/certregistry"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/lanediscovery"
)

const (
	EligibilityEligible = "eligible"
	EligibilityBlocked  = "blocked"
)

// Gates defines required certification predicates for lane use (read-only).
type Gates struct {
	NoRestart                 bool `json:"noRestart"`
	NoRollout                 bool `json:"noRollout"`
	NoRecreate                bool `json:"noRecreate"`
	NoInterruption            bool `json:"noInterruption"`
	ShellContinuityRequired   bool `json:"shellContinuityRequired"`
	RollbackRequired          bool `json:"rollbackRequired"`
	EvidenceFreshnessRequired bool `json:"evidenceFreshnessRequired"`
}

// DefaultNativeLiveGates is the KHR-V gate set for native-live simulation.
func DefaultNativeLiveGates() Gates {
	return Gates{
		NoRestart: true, NoRollout: true, NoRecreate: true, NoInterruption: true,
		ShellContinuityRequired: true, RollbackRequired: true,
		EvidenceFreshnessRequired: true,
	}
}

// EligibilityOutcome is the gate evaluation result for one lane.
type EligibilityOutcome struct {
	EligibilityState string   `json:"eligibilityState"`
	Eligible         bool     `json:"eligible"`
	Blocked          bool     `json:"blocked"`
	BlockedReason    string   `json:"blockedReason,omitempty"`
	StaleEvidence    bool     `json:"staleEvidence"`
	UncertifiedLane  bool     `json:"uncertifiedLane"`
	FailedGates      []string `json:"failedGates,omitempty"`
}

// Evaluate applies gates to a lane using the certification registry.
func Evaluate(lane string, reg *certregistry.Registry, gates Gates, now time.Time) EligibilityOutcome {
	out := EligibilityOutcome{
		EligibilityState: EligibilityBlocked,
		Blocked:          true,
	}
	if reg == nil {
		out.UncertifiedLane = true
		out.BlockedReason = "uncertified lane: certification registry not loaded"
		return out
	}
	entry := reg.FindByLane(lane)
	if entry == nil {
		out.UncertifiedLane = true
		out.BlockedReason = "uncertified lane: no registry entry for " + lane
		return out
	}
	if entry.CertificationState != certregistry.CertStateCertified {
		out.BlockedReason = "certification state " + entry.CertificationState
		out.FailedGates = append(out.FailedGates, "certificationState")
		return out
	}
	if gates.EvidenceFreshnessRequired && !certregistry.IsFresh(*entry, now) {
		out.StaleEvidence = true
		out.BlockedReason = "stale certification evidence"
		out.FailedGates = append(out.FailedGates, "evidenceFreshnessRequired")
		return out
	}
	failed := checkAttestation(entry.Attestation, gates)
	if len(failed) > 0 {
		out.BlockedReason = "policy gate failed: " + strings.Join(failed, ", ")
		out.FailedGates = failed
		return out
	}
	out.EligibilityState = EligibilityEligible
	out.Eligible = true
	out.Blocked = false
	return out
}

func checkAttestation(a certregistry.LaneAttestation, gates Gates) []string {
	var failed []string
	if gates.NoRestart && !a.NoRestart {
		failed = append(failed, "noRestart")
	}
	if gates.NoRollout && !a.NoRollout {
		failed = append(failed, "noRollout")
	}
	if gates.NoRecreate && !a.NoRecreate {
		failed = append(failed, "noRecreate")
	}
	if gates.NoInterruption && !a.NoInterruption {
		failed = append(failed, "noInterruption")
	}
	if gates.ShellContinuityRequired && !a.ShellContinuityPreserved {
		failed = append(failed, "shellContinuityRequired")
	}
	if gates.RollbackRequired && !a.RollbackObserved {
		failed = append(failed, "rollbackRequired")
	}
	return failed
}

// RequiresRegistry reports whether lane type participates in certification gating.
func RequiresRegistry(lane string) bool {
	return lane == lanediscovery.LaneNativeLive
}
