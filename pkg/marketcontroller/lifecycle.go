package marketcontroller

import (
	"fmt"
	"time"
)

const defaultLeaseTTLSeconds = 3600

// RefreshLifecycles expires stale leases, actions, and futures.
func RefreshLifecycles(state *ControllerState, now time.Time) (expiredLeases, expiredFutures, invalidatedActions int) {
	for i, lc := range state.LeaseLifecycles {
		expiresAt, _ := time.Parse(time.RFC3339, strOr(lc["expiresAt"]))
		if !expiresAt.IsZero() && now.After(expiresAt) && !boolOr(lc["expired"]) {
			state.LeaseLifecycles[i]["expired"] = true
			state.LeaseLifecycles[i]["state"] = "expired"
			state.LeaseLifecycles[i]["invalidationReason"] = "lease_ttl_expired"
			expiredLeases++
		}
	}
	for i, lc := range state.FutureLifecycles {
		expiresAt, _ := time.Parse(time.RFC3339, strOr(lc["expiresAt"]))
		if !expiresAt.IsZero() && now.After(expiresAt) && !boolOr(lc["expired"]) {
			state.FutureLifecycles[i]["expired"] = true
			state.FutureLifecycles[i]["state"] = "expired"
			expiredFutures++
		}
	}
	for i, lc := range state.ActionLifecycles {
		expiresAt, _ := time.Parse(time.RFC3339, strOr(lc["expiresAt"]))
		st := strOr(lc["state"])
		if !expiresAt.IsZero() && now.After(expiresAt) && st != "invalidated" && st != "closed" {
			state.ActionLifecycles[i]["state"] = "invalidated"
			state.ActionLifecycles[i]["invalidationReason"] = "action_ttl_expired"
			invalidatedActions++
		}
	}
	return
}

func newLeaseLifecycle(leaseID, actionID, donor, receiver, resource, amount string, now time.Time, ttl int) map[string]interface{} {
	exp := now.Add(time.Duration(ttl) * time.Second)
	return map[string]interface{}{
		"leaseLifecycleId": fmt.Sprintf("lease-lc-%s", leaseID),
		"leaseId":          leaseID,
		"actionId":         actionID,
		"donorShellId":     donor,
		"receiverShellId":  receiver,
		"resource":         resource,
		"amount":           amount,
		"state":            "proposed",
		"createdAt":        now.Format(time.RFC3339),
		"expiresAt":        exp.Format(time.RFC3339),
		"ttlSeconds":       ttl,
		"expired":          false,
		"evidenceRefs":     []interface{}{leaseID},
		"claimBoundary":    "lease lifecycle; not realized until post-verify",
	}
}

func newActionLifecycle(action map[string]interface{}, now time.Time) map[string]interface{} {
	return map[string]interface{}{
		"actionLifecycleId":            fmt.Sprintf("action-lc-%s", strOr(action["actionId"])),
		"actionId":                     action["actionId"],
		"donorShellId":                 action["donorShellId"],
		"receiverShellId":              action["receiverShellId"],
		"resource":                     action["resource"],
		"amount":                       action["amount"],
		"state":                        "generated",
		"executionScopeRecommendation": action["executionScopeRecommendation"],
		"previousState":                "",
		"nextState":                    "queued",
		"createdAt":                    now.Format(time.RFC3339),
		"updatedAt":                    now.Format(time.RFC3339),
		"expiresAt":                    now.Add(24 * time.Hour).Format(time.RFC3339),
		"selectedForExecution":         false,
		"blockers":                     action["blockers"],
		"evidenceRefs":                 []interface{}{action["actionId"]},
		"claimBoundary":                "action lifecycle; projected not realized",
	}
}

func newFutureLifecycle(future map[string]interface{}, now time.Time) map[string]interface{} {
	return map[string]interface{}{
		"futureLifecycleId":   fmt.Sprintf("future-lc-%s", strOr(future["futureId"])),
		"futureId":            future["futureId"],
		"donorShellId":        future["donorShellId"],
		"receiverShellId":     future["receiverShellId"],
		"resource":            future["resource"],
		"amount":              future["amount"],
		"state":               "generated",
		"createdAt":           now.Format(time.RFC3339),
		"refreshedAt":         now.Format(time.RFC3339),
		"expiresAt":           future["expiration"],
		"expired":             false,
		"confidence":          future["confidence"],
		"pressureProbability": future["pressureProbability"],
		"invalidationReasons": future["invalidationReasons"],
		"triggerCondition":    future["triggerCondition"],
		"evidenceRefs":        []interface{}{future["futureId"]},
		"claimBoundary":       "future lifecycle; projected not realized",
	}
}

func advanceActionLifecycle(lc map[string]interface{}, next string, now time.Time) {
	lc["previousState"] = lc["state"]
	lc["state"] = next
	lc["nextState"] = ""
	lc["updatedAt"] = now.Format(time.RFC3339)
}

func refreshFutureLifecycle(lc map[string]interface{}, now time.Time) {
	lc["state"] = "refreshed"
	lc["refreshedAt"] = now.Format(time.RFC3339)
}
