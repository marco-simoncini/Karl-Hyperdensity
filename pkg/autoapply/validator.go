package autoapply

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const MilestoneGuardedAutoApplySandboxNonprod = "hyperdensity_guarded_auto_apply_sandbox_nonprod_v1"

var allowedAutoExecutionScopes = map[string]bool{
	"sandbox": true, "nonprod": true, "sandbox_and_nonprod": true,
}

var forbiddenAutoExecutionScopes = map[string]bool{
	"production": true, "production_canary": true, "production_auto": true,
}

var executableCandidateStates = map[string]bool{
	"sandbox_ready": true, "nonprod_ready": true,
}

var requiredAuditEventTypes = []string{
	"auto_selection", "preflight_recheck", "safety_reserved", "auto_apply_requested",
	"fluidvirt_invoked", "mutation_observed", "post_verify_evaluated",
	"rollback_decision", "rollback_executed", "auto_apply_closed",
}

var forbiddenPositiveClaims = []string{
	"guaranteed savings active",
	"universal performance improvement",
	"production autonomous apply",
	"autonomous production mutation",
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

func requireEvidenceRefs(doc map[string]interface{}) error {
	refs, ok := doc["evidenceRefs"].([]interface{})
	if !ok || len(refs) == 0 {
		return fmt.Errorf("evidenceRefs required")
	}
	return nil
}

func boolOr(v interface{}) bool {
	b, _ := v.(bool)
	return b
}

// ValidateAutoApplySelection checks selection invariants.
func ValidateAutoApplySelection(doc map[string]interface{}) error {
	selected := boolOr(doc["selected"])
	state, _ := doc["candidateState"].(string)
	prodScope := boolOr(doc["productionScope"])
	if selected {
		if !executableCandidateStates[state] {
			return fmt.Errorf("selected=true forbidden for candidateState %s", state)
		}
		if prodScope || !boolOr(doc["sandboxOrNonProd"]) {
			return fmt.Errorf("selected=true requires productionScope=false and sandboxOrNonProd=true")
		}
	}
	if state == "candidate_only" && selected {
		return fmt.Errorf("candidate_only cannot be selected")
	}
	return requireClaimBoundary(doc)
}

// ValidateAutoPreflightRecheck checks preflight invariants.
func ValidateAutoPreflightRecheck(doc map[string]interface{}) error {
	status, _ := doc["preflightStatus"].(string)
	if status == "passed" {
		if !boolOr(doc["killSwitchStillClear"]) || !boolOr(doc["circuitBreakerStillClosed"]) {
			return fmt.Errorf("passed preflight requires kill switch clear and circuit breaker closed")
		}
	}
	return requireClaimBoundary(doc)
}

// ValidateAutoApplyRequest checks request invariants.
func ValidateAutoApplyRequest(doc map[string]interface{}) error {
	if doc["requestedByComponent"] != "Karl-Hyperdensity" {
		return fmt.Errorf("requestedByComponent must be Karl-Hyperdensity")
	}
	if doc["actuator"] != "FluidVirt" {
		return fmt.Errorf("actuator must be FluidVirt")
	}
	if err := requireFalse(doc, "productionScope"); err != nil {
		return err
	}
	return requireFalse(doc, "productionAutonomousApplyAllowed")
}

// ValidateAutoFluidvirtInvocation checks invocation invariants.
func ValidateAutoFluidvirtInvocation(doc map[string]interface{}, preflightPassed, reservationReserved bool, killClear, circuitClosed bool) error {
	if doc["actuator"] != "FluidVirt" {
		return fmt.Errorf("actuator must be FluidVirt")
	}
	if doc["invocationMode"] != "guarded_auto_sandbox_or_nonprod" {
		return fmt.Errorf("invocationMode must be guarded_auto_sandbox_or_nonprod")
	}
	if err := requireFalse(doc, "productionScope"); err != nil {
		return err
	}
	if err := requireFalse(doc, "productionAutonomousApplyAllowed"); err != nil {
		return err
	}
	if err := requireFalse(doc, "rawRuntimeControlsExposed"); err != nil {
		return err
	}
	if boolOr(doc["mutationExecuted"]) {
		if !preflightPassed {
			return fmt.Errorf("mutationExecuted requires successful preflight")
		}
		if !reservationReserved {
			return fmt.Errorf("mutationExecuted requires safety reservation")
		}
		if !killClear || !circuitClosed {
			return fmt.Errorf("mutationExecuted requires kill switch clear and circuit breaker closed")
		}
	}
	return requireClaimBoundary(doc)
}

// ValidateAutoMutationObservation checks mutation observation invariants.
func ValidateAutoMutationObservation(doc map[string]interface{}) error {
	if !boolOr(doc["identityPreserved"]) {
		return fmt.Errorf("identityPreserved must be true")
	}
	for _, key := range []string{"rebootUsed", "recreateUsed", "rolloutUsed", "migrationUsed"} {
		if boolOr(doc[key]) {
			return fmt.Errorf("%s must be false", key)
		}
	}
	return requireClaimBoundary(doc)
}

// ValidateAutoPostVerifyResult checks post-verify invariants.
func ValidateAutoPostVerifyResult(doc map[string]interface{}) error {
	status, _ := doc["verifyStatus"].(string)
	if status == "passed" {
		if !boolOr(doc["runtimeDeltaVerified"]) {
			return fmt.Errorf("passed post-verify requires runtimeDeltaVerified=true")
		}
		if doc["sloGuardStatus"] != "passed" {
			return fmt.Errorf("passed post-verify requires sloGuardStatus=passed")
		}
		if doc["donorHealthStatus"] != "preserved" {
			return fmt.Errorf("passed post-verify requires donorHealthStatus=preserved")
		}
		if doc["noRegressionStatus"] != "certified" {
			return fmt.Errorf("passed post-verify requires noRegressionStatus=certified")
		}
		if !boolOr(doc["rollbackReady"]) || boolOr(doc["rollbackRequired"]) {
			return fmt.Errorf("passed post-verify requires rollbackReady=true and rollbackRequired=false")
		}
		if boolOr(doc["mutationMatchesPlan"]) == false && doc["mutationMatchesPlan"] != nil {
			if m, ok := doc["mutationMatchesPlan"].(bool); ok && !m {
				return fmt.Errorf("passed post-verify requires mutationMatchesPlan=true")
			}
		}
	}
	return requireClaimBoundary(doc)
}

// ValidateAutoRollbackExecution checks rollback execution invariants.
func ValidateAutoRollbackExecution(doc map[string]interface{}) error {
	if boolOr(doc["executed"]) && !boolOr(doc["rollbackPassed"]) && boolOr(doc["healthRestored"]) {
		return fmt.Errorf("failed rollback cannot have healthRestored=true")
	}
	return requireClaimBoundary(doc)
}

// ValidateAutoApplyAuditTrail checks audit trail completeness.
func ValidateAutoApplyAuditTrail(doc map[string]interface{}) error {
	events, ok := doc["events"].([]interface{})
	if !ok || len(events) == 0 {
		return fmt.Errorf("audit trail events required")
	}
	found := map[string]bool{}
	for _, item := range events {
		e, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		et, _ := e["eventType"].(string)
		found[et] = true
	}
	for _, want := range requiredAuditEventTypes {
		if !found[want] {
			return fmt.Errorf("audit trail missing event type %s", want)
		}
	}
	return requireClaimBoundary(doc)
}

// ValidateGuardedAutoApplySandboxSurface validates ConfigMap-ready surface.
func ValidateGuardedAutoApplySandboxSurface(doc map[string]interface{}) error {
	if doc["milestone"] != MilestoneGuardedAutoApplySandboxNonprod {
		return fmt.Errorf("milestone must be %s", MilestoneGuardedAutoApplySandboxNonprod)
	}
	if err := requireTrue(doc, "sandboxAutoApplyAllowed"); err != nil {
		return err
	}
	if err := requireTrue(doc, "nonProdAutoApplyAllowed"); err != nil {
		return err
	}
	if err := requireTrue(doc, "autoApplyExecutionEnabled"); err != nil {
		return err
	}
	if err := requireFalse(doc, "productionAutonomousApplyAllowed"); err != nil {
		return err
	}
	if err := requireFalse(doc, "productionScope"); err != nil {
		return err
	}
	if err := requireFalse(doc, "productionMutationAllowed"); err != nil {
		return err
	}
	if err := requireFalse(doc, "rawRuntimeControlsExposed"); err != nil {
		return err
	}
	scope, _ := doc["autoExecutionScope"].(string)
	if forbiddenAutoExecutionScopes[scope] || !allowedAutoExecutionScopes[scope] {
		return fmt.Errorf("invalid autoExecutionScope: %s", scope)
	}
	inv, ok := doc["safetyInvariants"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("safetyInvariants required")
	}
	for _, key := range []string{"productionAutonomousApplyAllowed", "productionScope", "productionMutationAllowed", "rawRuntimeControlsExposed", "dashboardExecutor", "fluidvirtPolicyAuthority", "fluidvirtMarketController", "inventoryRuntimeApply"} {
		if v, ok := inv[key].(bool); !ok || v {
			return fmt.Errorf("safetyInvariants.%s must be false", key)
		}
	}
	sot, ok := doc["sourceOfTruth"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("sourceOfTruth required")
	}
	if !strings.Contains(strings.ToLower(fmt.Sprint(sot["runtimeActuator"])), "fluidvirt") {
		return fmt.Errorf("runtimeActuator must be FluidVirt")
	}
	if strings.Contains(strings.ToLower(fmt.Sprint(sot["executionProjection"])), "executor") {
		return fmt.Errorf("Dashboard must not be executor")
	}

	var sandboxSuccess, nonprodSuccess, rollbackExecuted, deniedCandidateOnly, deniedProdScope, deniedWindows, deniedKillSwitch, deniedCircuitBreaker, deniedSynthetic bool

	if sels, ok := doc["selections"].([]interface{}); ok {
		for _, item := range sels {
			s, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			aid, _ := s["actionId"].(string)
			if aid == "action-sandbox-ready-005" && boolOr(s["selected"]) {
				sandboxSuccess = true
			}
			if aid == "action-nonprod-ready-010" && boolOr(s["selected"]) {
				nonprodSuccess = true
			}
			if aid == "action-cpu-lease-001" && !boolOr(s["selected"]) {
				deniedCandidateOnly = true
			}
			if aid == "action-prod-auto-012" && !boolOr(s["selected"]) {
				deniedProdScope = true
			}
			if aid == "action-windows-remediation-004" && !boolOr(s["selected"]) {
				deniedWindows = true
			}
			if aid == "action-synthetic-shadow-009" && !boolOr(s["selected"]) {
				deniedSynthetic = true
			}
			if err := ValidateAutoApplySelection(s); err != nil {
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
			aid, _ := p["actionId"].(string)
			if aid == "action-cpu-lease-006" && p["preflightStatus"] == "blocked" {
				deniedKillSwitch = true
			}
			if aid == "action-circuit-open-013" && p["preflightStatus"] == "blocked" {
				deniedCircuitBreaker = true
			}
		}
	}
	if invocs, ok := doc["fluidvirtInvocations"].([]interface{}); ok {
		for _, item := range invocs {
			inv, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			if inv["autoApplyRequestId"] == "auto-request-sandbox-005" && boolOr(inv["mutationExecuted"]) {
				sandboxSuccess = true
			}
			if inv["autoApplyRequestId"] == "auto-request-nonprod-010" && boolOr(inv["mutationExecuted"]) {
				nonprodSuccess = true
			}
			if err := ValidateAutoFluidvirtInvocation(inv, true, true, true, true); err != nil {
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
			if err := ValidateAutoMutationObservation(m); err != nil {
				return fmt.Errorf("mutationObservations: %w", err)
			}
		}
	}
	if pvs, ok := doc["postVerifyResults"].([]interface{}); ok {
		for _, item := range pvs {
			pv, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			if err := ValidateAutoPostVerifyResult(pv); err != nil {
				return fmt.Errorf("postVerifyResults: %w", err)
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
			if err := ValidateAutoRollbackExecution(re); err != nil {
				return fmt.Errorf("rollbackExecutions: %w", err)
			}
		}
	}
	if audits, ok := doc["auditTrail"].([]interface{}); ok {
		for _, item := range audits {
			a, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			if err := ValidateAutoApplyAuditTrail(a); err != nil {
				return fmt.Errorf("auditTrail: %w", err)
			}
		}
	}
	if !sandboxSuccess || !nonprodSuccess || !rollbackExecuted || !deniedCandidateOnly || !deniedProdScope || !deniedWindows || !deniedKillSwitch || !deniedCircuitBreaker || !deniedSynthetic {
		return fmt.Errorf("surface must include sandbox/nonprod success, rollback, and all denial scenarios")
	}
	return rejectForbiddenPositiveClaims(doc)
}

func rejectForbiddenPositiveClaims(v interface{}) error {
	positives := collectPositiveStrings(v)
	merged := strings.ToLower(strings.Join(positives, "\n"))
	for _, phrase := range forbiddenPositiveClaims {
		if strings.Contains(merged, phrase) {
			if strings.Contains(merged, "not "+phrase) || strings.Contains(merged, "no "+phrase) {
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
			if isPositiveKey(k) {
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

func isPositiveKey(k string) bool {
	switch k {
	case "claimBoundary", "claimBoundaries":
		return true
	default:
		return strings.HasSuffix(k, "Phrases") && !strings.HasPrefix(k, "forbidden")
	}
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

// ValidateSprint8Examples validates all Sprint 8 reference examples.
func ValidateSprint8Examples(repoRoot string) error {
	files := map[string]func(map[string]interface{}) error{
		"auto-apply-execution-policy-reference.json": func(d map[string]interface{}) error {
			if err := requireFalse(d, "productionAutonomousApplyAllowed"); err != nil {
				return err
			}
			return requireClaimBoundary(d)
		},
		"auto-apply-selection-reference.json":           ValidateAutoApplySelection,
		"auto-preflight-recheck-reference.json":         ValidateAutoPreflightRecheck,
		"auto-apply-request-reference.json":             ValidateAutoApplyRequest,
		"auto-fluidvirt-invocation-reference.json":      func(d map[string]interface{}) error { return ValidateAutoFluidvirtInvocation(d, true, true, true, true) },
		"auto-mutation-observation-reference.json":      ValidateAutoMutationObservation,
		"auto-post-verify-result-reference.json":        ValidateAutoPostVerifyResult,
		"auto-rollback-decision-reference.json":         func(d map[string]interface{}) error { return requireClaimBoundary(d) },
		"auto-rollback-execution-reference.json":        ValidateAutoRollbackExecution,
		"auto-apply-audit-trail-reference.json":         ValidateAutoApplyAuditTrail,
		"auto-apply-safety-reservation-reference.json":  func(d map[string]interface{}) error { return requireClaimBoundary(d) },
		"guarded-auto-apply-sandbox-nonprod-reference.json": ValidateGuardedAutoApplySandboxSurface,
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
		if name != "guarded-auto-apply-sandbox-nonprod-reference.json" {
			if err := rejectForbiddenPositiveClaims(doc); err != nil {
				return fmt.Errorf("%s: %w", name, err)
			}
		}
	}
	return nil
}

// SchemaFilesRequiredSprint8 returns Sprint 8 schema basenames.
func SchemaFilesRequiredSprint8() []string {
	return []string{
		"guarded-auto-apply-sandbox-nonprod-v1.schema.json",
		"auto-apply-execution-policy-v1.schema.json",
		"auto-apply-selection-v1.schema.json",
		"auto-preflight-recheck-v1.schema.json",
		"auto-apply-request-v1.schema.json",
		"auto-fluidvirt-invocation-v1.schema.json",
		"auto-mutation-observation-v1.schema.json",
		"auto-post-verify-result-v1.schema.json",
		"auto-rollback-decision-v1.schema.json",
		"auto-rollback-execution-v1.schema.json",
		"auto-apply-audit-trail-v1.schema.json",
		"auto-apply-safety-reservation-v1.schema.json",
	}
}
