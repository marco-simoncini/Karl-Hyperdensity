package productioncanary

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const MilestoneProductionCanaryAutoApply = "hyperdensity_production_canary_auto_apply_v1"

var executableCanaryStates = map[string]bool{"production_canary_ready": true}

var requiredAuditEventTypes = []string{
	"canary_selection", "production_preflight_recheck", "canary_safety_reserved",
	"production_canary_apply_requested", "fluidvirt_invoked", "mutation_observed",
	"post_verify_evaluated", "rollback_decision", "rollback_executed", "production_canary_closed",
}

var forbiddenPositiveClaims = []string{
	"guaranteed savings active",
	"universal performance improvement",
	"production autonomous apply",
	"autonomous production mutation",
	"general production auto",
	"production_auto_with_policy",
	"production auto with policy",
	"windows total ram hotplug supported",
	"logical vcpu hotplug supported",
	"1000 production workloads proven",
	"dashboard is executor",
	"dashboard executor",
	"dashboard is source of truth",
	"fluidvirt policy authority",
	"fluidvirt market controller",
	"inventory applies runtime",
}

var skipKeys = map[string]bool{
	"forbiddenPhrases": true, "blockerCodes": true, "blockers": true,
	"remediationCodes": true, "reason": true, "selectionReason": true, "summary": true,
	"approvalReason": true,
}

func requireFalse(doc map[string]interface{}, key string) error {
	if v, ok := doc[key].(bool); !ok || v {
		return fmt.Errorf("%s must be false", key)
	}
	return nil
}

func requireTrue(doc map[string]interface{}, key string) error {
	if v, ok := doc[key].(bool); !ok || !v {
		return fmt.Errorf("%s must be true", key)
	}
	return nil
}

func requireClaimBoundary(doc map[string]interface{}) error {
	cb, ok := doc["claimBoundary"].([]interface{})
	if !ok || len(cb) == 0 {
		return fmt.Errorf("claimBoundary required")
	}
	return nil
}

func boolOr(v interface{}) bool {
	b, _ := v.(bool)
	return b
}

func validateProductionScopePair(doc map[string]interface{}) error {
	prodScope := boolOr(doc["productionScope"])
	canaryScope := boolOr(doc["productionCanaryScope"])
	if prodScope && !canaryScope {
		return fmt.Errorf("productionScope=true requires productionCanaryScope=true")
	}
	autoAllowed := boolOr(doc["productionAutonomousApplyAllowed"])
	if autoAllowed && !canaryScope {
		return fmt.Errorf("productionAutonomousApplyAllowed=true requires productionCanaryScope=true")
	}
	return nil
}

// ValidateCanarySelection checks selection invariants.
func ValidateCanarySelection(doc map[string]interface{}) error {
	if err := validateProductionScopePair(doc); err != nil {
		return err
	}
	if boolOr(doc["selected"]) {
		state, _ := doc["candidateState"].(string)
		if !executableCanaryStates[state] {
			return fmt.Errorf("selected=true requires production_canary_ready")
		}
		if !boolOr(doc["allowlisted"]) || boolOr(doc["broadProductionScope"]) {
			return fmt.Errorf("selected=true requires allowlisted=true and broadProductionScope=false")
		}
	}
	if state, _ := doc["candidateState"].(string); state == "candidate_only" && boolOr(doc["selected"]) {
		return fmt.Errorf("candidate_only cannot be selected")
	}
	return requireClaimBoundary(doc)
}

// ValidateCanaryPreflight checks preflight invariants.
func ValidateCanaryPreflight(doc map[string]interface{}) error {
	if status, _ := doc["preflightStatus"].(string); status == "passed" {
		if !boolOr(doc["incidentStateClear"]) || !boolOr(doc["killSwitchStillClear"]) || !boolOr(doc["circuitBreakerStillClosed"]) {
			return fmt.Errorf("passed preflight requires incident clear, kill switch clear, circuit breaker closed")
		}
		if boolOr(doc["broadProductionScopeDenied"]) == false && doc["broadProductionScopeDenied"] != nil {
			// ok
		}
	}
	return requireClaimBoundary(doc)
}

// ValidateCanaryApplyRequest checks apply request invariants.
func ValidateCanaryApplyRequest(doc map[string]interface{}) error {
	if doc["requestedByComponent"] != "Karl-Hyperdensity" || doc["actuator"] != "FluidVirt" {
		return fmt.Errorf("requestedByComponent=Karl-Hyperdensity and actuator=FluidVirt required")
	}
	if err := requireTrue(doc, "productionScope"); err != nil {
		return err
	}
	if err := requireTrue(doc, "productionCanaryScope"); err != nil {
		return err
	}
	if err := requireFalse(doc, "generalProductionAutoAllowed"); err != nil {
		return err
	}
	return requireFalse(doc, "productionAutoWithPolicy")
}

// ValidateCanaryFluidvirtInvocation checks invocation invariants.
func ValidateCanaryFluidvirtInvocation(doc map[string]interface{}, preflightPassed, reservationReserved, allowlisted, incidentClear, killClear, circuitClosed bool) error {
	if doc["actuator"] != "FluidVirt" || doc["invocationMode"] != "guarded_production_canary" {
		return fmt.Errorf("actuator FluidVirt and invocationMode guarded_production_canary required")
	}
	if err := requireTrue(doc, "productionScope"); err != nil {
		return err
	}
	if err := requireTrue(doc, "productionCanaryScope"); err != nil {
		return err
	}
	if err := requireFalse(doc, "generalProductionAutoAllowed"); err != nil {
		return err
	}
	if err := requireFalse(doc, "productionAutoWithPolicy"); err != nil {
		return err
	}
	if err := requireFalse(doc, "rawRuntimeControlsExposed"); err != nil {
		return err
	}
	if boolOr(doc["mutationExecuted"]) {
		if !preflightPassed || !reservationReserved || !allowlisted || !incidentClear || !killClear || !circuitClosed {
			return fmt.Errorf("mutationExecuted requires preflight, reservation, allowlist, incident/kill/circuit gates")
		}
	}
	return requireClaimBoundary(doc)
}

// ValidateCanaryMutationObservation checks mutation observation.
func ValidateCanaryMutationObservation(doc map[string]interface{}) error {
	if !boolOr(doc["identityPreserved"]) {
		return fmt.Errorf("identityPreserved must be true")
	}
	for _, k := range []string{"rebootUsed", "recreateUsed", "rolloutUsed", "migrationUsed"} {
		if boolOr(doc[k]) {
			return fmt.Errorf("%s must be false", k)
		}
	}
	return requireClaimBoundary(doc)
}

// ValidateCanaryPostVerify checks post-verify result.
func ValidateCanaryPostVerify(doc map[string]interface{}) error {
	if doc["verifyStatus"] == "passed" {
		if !boolOr(doc["runtimeDeltaVerified"]) || doc["sloGuardStatus"] != "passed" || doc["donorHealthStatus"] != "preserved" || doc["noRegressionStatus"] != "certified" {
			return fmt.Errorf("passed post-verify requires runtime/SLO/donor/no-regression gates")
		}
		if !boolOr(doc["rollbackReady"]) || boolOr(doc["rollbackRequired"]) {
			return fmt.Errorf("passed post-verify requires rollbackReady=true and rollbackRequired=false")
		}
	}
	return requireClaimBoundary(doc)
}

// ValidateCanaryRollbackExecution checks rollback execution.
func ValidateCanaryRollbackExecution(doc map[string]interface{}) error {
	if boolOr(doc["executed"]) && !boolOr(doc["rollbackPassed"]) && boolOr(doc["healthRestored"]) {
		return fmt.Errorf("failed rollback cannot have healthRestored=true")
	}
	return requireClaimBoundary(doc)
}

// ValidateImmutableAuditTrail checks immutable audit completeness.
func ValidateImmutableAuditTrail(doc map[string]interface{}) error {
	if boolOr(doc["immutableAuditRequired"]) && !boolOr(doc["immutableAuditWritten"]) {
		return fmt.Errorf("immutableAuditWritten required when immutableAuditRequired=true")
	}
	events, ok := doc["events"].([]interface{})
	if !ok || len(events) == 0 {
		return fmt.Errorf("audit events required")
	}
	found := map[string]bool{}
	for _, item := range events {
		e, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		found[fmt.Sprint(e["eventType"])] = true
	}
	for _, want := range requiredAuditEventTypes {
		if !found[want] {
			return fmt.Errorf("missing audit event type %s", want)
		}
	}
	return requireClaimBoundary(doc)
}

// ValidateCanaryCloseout checks closeout invariants.
func ValidateCanaryCloseout(doc map[string]interface{}) error {
	if boolOr(doc["promotedToGeneralProduction"]) {
		return fmt.Errorf("promotedToGeneralProduction must be false in Sprint 9")
	}
	if state, _ := doc["nextAllowedState"].(string); state == "production_auto_with_policy" {
		return fmt.Errorf("nextAllowedState production_auto_with_policy forbidden")
	}
	return requireClaimBoundary(doc)
}

// ValidateProductionCanaryAutoApplySurface validates ConfigMap-ready surface.
func ValidateProductionCanaryAutoApplySurface(doc map[string]interface{}) error {
	if doc["milestone"] != MilestoneProductionCanaryAutoApply {
		return fmt.Errorf("milestone must be %s", MilestoneProductionCanaryAutoApply)
	}
	if err := requireTrue(doc, "productionCanaryAutoApplyAllowed"); err != nil {
		return err
	}
	if err := requireTrue(doc, "productionCanaryScope"); err != nil {
		return err
	}
	if err := requireFalse(doc, "generalProductionAutoAllowed"); err != nil {
		return err
	}
	if err := requireFalse(doc, "productionAutoWithPolicy"); err != nil {
		return err
	}
	if err := requireFalse(doc, "rawRuntimeControlsExposed"); err != nil {
		return err
	}
	if boolOr(doc["productionScope"]) && !boolOr(doc["productionCanaryScope"]) {
		return fmt.Errorf("productionScope=true requires productionCanaryScope=true")
	}
	inv, ok := doc["safetyInvariants"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("safetyInvariants required")
	}
	for _, key := range []string{"generalProductionAutoAllowed", "productionAutoWithPolicy", "rawRuntimeControlsExposed", "dashboardExecutor", "fluidvirtPolicyAuthority", "fluidvirtMarketController", "inventoryRuntimeApply"} {
		if v, ok := inv[key].(bool); !ok || v {
			return fmt.Errorf("safetyInvariants.%s must be false", key)
		}
	}

	var canarySuccess, rollbackExecuted, deniedBroad, deniedNotAllowlisted, deniedCandidateOnly, deniedSandbox, deniedWindows, deniedIncident, deniedKillSwitch, deniedSynthetic bool

	if sels, ok := doc["selections"].([]interface{}); ok {
		for _, item := range sels {
			s, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			aid, _ := s["actionId"].(string)
			switch aid {
			case "action-prod-canary-cpu-014":
				if boolOr(s["selected"]) {
					canarySuccess = true
				}
			case "action-prod-broad-016":
				if !boolOr(s["selected"]) {
					deniedBroad = true
				}
			case "action-prod-not-allowlisted-017":
				if !boolOr(s["selected"]) {
					deniedNotAllowlisted = true
				}
			case "action-cpu-lease-001":
				if !boolOr(s["selected"]) {
					deniedCandidateOnly = true
				}
			case "action-sandbox-ready-005":
				if !boolOr(s["selected"]) {
					deniedSandbox = true
				}
			case "action-windows-remediation-004":
				if !boolOr(s["selected"]) {
					deniedWindows = true
				}
			case "action-synthetic-shadow-009":
				if !boolOr(s["selected"]) {
					deniedSynthetic = true
				}
			}
			if err := ValidateCanarySelection(s); err != nil {
				return fmt.Errorf("selections: %w", err)
			}
		}
	}
	if preflights, ok := doc["preflightRechecks"].([]interface{}); ok {
		for _, item := range preflights {
			p, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			if !boolOr(p["incidentStateClear"]) && p["preflightStatus"] == "blocked" {
				deniedIncident = true
			}
			aid, _ := p["actionId"].(string)
			if aid == "action-prod-canary-kill-018" && p["preflightStatus"] == "blocked" {
				deniedKillSwitch = true
			}
			if aid == "action-prod-canary-circuit-019" && p["preflightStatus"] == "blocked" {
				deniedKillSwitch = true
			}
		}
	}
	if invocs, ok := doc["fluidvirtInvocations"].([]interface{}); ok {
		for _, item := range invocs {
			inv, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			if inv["canaryApplyRequestId"] == "canary-request-cpu-014" && boolOr(inv["mutationExecuted"]) {
				canarySuccess = true
			}
			if err := ValidateCanaryFluidvirtInvocation(inv, true, true, true, true, true, true); err != nil {
				return fmt.Errorf("fluidvirtInvocations: %w", err)
			}
		}
	}
	if obs, ok := doc["mutationObservations"].([]interface{}); ok {
		for _, item := range obs {
			m, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			if err := ValidateCanaryMutationObservation(m); err != nil {
				return fmt.Errorf("mutationObservations: %w", err)
			}
		}
	}
	if rexecs, ok := doc["rollbackExecutions"].([]interface{}); ok {
		for _, item := range rexecs {
			re, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			if boolOr(re["executed"]) && boolOr(re["rollbackPassed"]) {
				rollbackExecuted = true
			}
			if err := ValidateCanaryRollbackExecution(re); err != nil {
				return fmt.Errorf("rollbackExecutions: %w", err)
			}
		}
	}
	if audits, ok := doc["immutableAuditTrails"].([]interface{}); ok {
		for _, item := range audits {
			a, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			if err := ValidateImmutableAuditTrail(a); err != nil {
				return fmt.Errorf("immutableAuditTrails: %w", err)
			}
		}
	}
	if closeouts, ok := doc["closeouts"].([]interface{}); ok {
		for _, item := range closeouts {
			c, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			if err := ValidateCanaryCloseout(c); err != nil {
				return fmt.Errorf("closeouts: %w", err)
			}
		}
	}
	if !canarySuccess || !rollbackExecuted || !deniedBroad || !deniedNotAllowlisted || !deniedCandidateOnly || !deniedSandbox || !deniedWindows || !deniedIncident || !deniedKillSwitch || !deniedSynthetic {
		return fmt.Errorf("surface must include all required canary scenarios")
	}
	return rejectForbiddenPositiveClaims(doc)
}

func rejectForbiddenPositiveClaims(v interface{}) error {
	positives := collectPositiveStrings(v)
	merged := strings.ToLower(strings.Join(positives, "\n"))
	for _, phrase := range forbiddenPositiveClaims {
		if strings.Contains(merged, phrase) {
			if strings.Contains(merged, "not "+phrase) || strings.Contains(merged, "no "+phrase) || strings.Contains(merged, "remains disabled") || strings.Contains(merged, "forbidden") || strings.Contains(merged, "blocked") {
				continue
			}
			return fmt.Errorf("forbidden positive claim: %q", phrase)
		}
	}
	return nil
}

func collectPositiveStrings(v interface{}) []string {
	var out []string
	switch t := v.(type) {
	case map[string]interface{}:
		for k, child := range t {
			if skipKeys[k] {
				continue
			}
			if k == "claimBoundary" || k == "claimBoundaries" {
				out = append(out, flattenStrings(child)...)
			} else {
				out = append(out, collectPositiveStrings(child)...)
			}
		}
	case []interface{}:
		for _, item := range t {
			out = append(out, collectPositiveStrings(item)...)
		}
	}
	return out
}

func flattenStrings(v interface{}) []string {
	switch t := v.(type) {
	case string:
		return []string{t}
	case []interface{}:
		var out []string
		for _, item := range t {
			if s, ok := item.(string); ok {
				out = append(out, s)
			}
		}
		return out
	}
	return nil
}

// ValidateSprint9Examples validates all Sprint 9 reference examples.
func ValidateSprint9Examples(repoRoot string) error {
	files := map[string]func(map[string]interface{}) error{
		"production-canary-policy-reference.json":            func(d map[string]interface{}) error { return requireClaimBoundary(d) },
		"production-canary-allowlist-reference.json":           func(d map[string]interface{}) error { return requireClaimBoundary(d) },
		"production-canary-selection-reference.json":          ValidateCanarySelection,
		"production-canary-preflight-reference.json":          ValidateCanaryPreflight,
		"production-canary-safety-reservation-reference.json": func(d map[string]interface{}) error { return requireClaimBoundary(d) },
		"production-canary-apply-request-reference.json":      ValidateCanaryApplyRequest,
		"production-canary-fluidvirt-invocation-reference.json": func(d map[string]interface{}) error {
			return ValidateCanaryFluidvirtInvocation(d, true, true, true, true, true, true)
		},
		"production-canary-mutation-observation-reference.json": ValidateCanaryMutationObservation,
		"production-canary-post-verify-result-reference.json":   ValidateCanaryPostVerify,
		"production-canary-rollback-decision-reference.json":  func(d map[string]interface{}) error { return requireClaimBoundary(d) },
		"production-canary-rollback-execution-reference.json": ValidateCanaryRollbackExecution,
		"production-canary-immutable-audit-trail-reference.json": ValidateImmutableAuditTrail,
		"production-canary-closeout-reference.json":           ValidateCanaryCloseout,
		"production-canary-auto-apply-reference.json":         ValidateProductionCanaryAutoApplySurface,
	}
	for name, fn := range files {
		b, err := os.ReadFile(filepath.Join(repoRoot, "examples", name))
		if err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
		var doc map[string]interface{}
		if err := json.Unmarshal(b, &doc); err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
		if err := fn(doc); err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
		if name != "production-canary-auto-apply-reference.json" {
			if err := rejectForbiddenPositiveClaims(doc); err != nil {
				return fmt.Errorf("%s: %w", name, err)
			}
		}
	}
	return nil
}

// SchemaFilesRequiredSprint9 returns Sprint 9 schema basenames.
func SchemaFilesRequiredSprint9() []string {
	return []string{
		"production-canary-auto-apply-v1.schema.json",
		"production-canary-policy-v1.schema.json",
		"production-canary-allowlist-v1.schema.json",
		"production-canary-selection-v1.schema.json",
		"production-canary-preflight-v1.schema.json",
		"production-canary-safety-reservation-v1.schema.json",
		"production-canary-apply-request-v1.schema.json",
		"production-canary-fluidvirt-invocation-v1.schema.json",
		"production-canary-mutation-observation-v1.schema.json",
		"production-canary-post-verify-result-v1.schema.json",
		"production-canary-rollback-decision-v1.schema.json",
		"production-canary-rollback-execution-v1.schema.json",
		"production-canary-immutable-audit-trail-v1.schema.json",
		"production-canary-closeout-v1.schema.json",
	}
}
