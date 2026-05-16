package actionapproval

import (
	"fmt"
	"time"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/certregistry"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/policygates"
)

// CanApprove validates gates before operator approval (no apply).
func CanApprove(a ActionApproval, reg *certregistry.Registry, gates policygates.Gates, now time.Time) error {
	if a.ApprovalState != StatePending {
		return fmt.Errorf("approval state %q is not pending", a.ApprovalState)
	}
	exp, err := time.Parse(time.RFC3339, a.ExpiresAt)
	if err != nil {
		return fmt.Errorf("invalid expiresAt: %w", err)
	}
	if now.After(exp) {
		return fmt.Errorf("approval expired at %s", a.ExpiresAt)
	}
	out := policygates.Evaluate(a.LaneID, reg, gates, now)
	if !out.Eligible {
		if out.StaleEvidence {
			return fmt.Errorf("stale certification blocks approval: %s", out.BlockedReason)
		}
		if len(out.FailedGates) > 0 {
			return fmt.Errorf("policy gate blocks approval: %s", out.BlockedReason)
		}
		return fmt.Errorf("lane not eligible for approval: %s", out.BlockedReason)
	}
	if !a.PolicyGateResult.Eligible {
		return fmt.Errorf("approval policyGateResult not eligible")
	}
	return nil
}

// Approve transitions pending → approved (local evidence only).
func Approve(a ActionApproval, reg *certregistry.Registry, gates policygates.Gates, approvedBy string, now time.Time) (ActionApproval, error) {
	if err := CanApprove(a, reg, gates, now); err != nil {
		return a, err
	}
	a.ApprovalState = StateApproved
	a.ApprovedBy = approvedBy
	a.ApprovedAt = now.UTC().Format(time.RFC3339)
	a.Reason = "operator approved; no apply performed"
	return a, nil
}

// Reject transitions pending → rejected (local evidence only).
func Reject(a ActionApproval, rejectedBy, reason string, now time.Time) (ActionApproval, error) {
	if a.ApprovalState != StatePending {
		return a, fmt.Errorf("approval state %q is not pending", a.ApprovalState)
	}
	a.ApprovalState = StateRejected
	a.ApprovedBy = rejectedBy
	a.ApprovedAt = now.UTC().Format(time.RFC3339)
	if reason == "" {
		reason = "operator rejected"
	}
	a.Reason = reason
	return a, nil
}

// ExpireIfNeeded marks pending approvals past expiresAt as expired.
func ExpireIfNeeded(a ActionApproval, now time.Time) (ActionApproval, bool) {
	if a.ApprovalState != StatePending {
		return a, false
	}
	exp, err := time.Parse(time.RFC3339, a.ExpiresAt)
	if err != nil || !now.After(exp) {
		return a, false
	}
	a.ApprovalState = StateExpired
	a.Reason = "approval TTL expired"
	return a, true
}

// SimulateExpire forces expiry for evidence (sets expiresAt in past and state expired).
func SimulateExpire(a ActionApproval, now time.Time) ActionApproval {
	a.ExpiresAt = now.Add(-time.Hour).UTC().Format(time.RFC3339)
	a.ApprovalState = StateExpired
	a.Reason = "approval TTL expired (simulated)"
	return a
}
