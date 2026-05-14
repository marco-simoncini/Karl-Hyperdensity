package guardedauto

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const MilestoneGuardedAutoPolicyEngine = "hyperdensity_guarded_auto_policy_engine_v1"

var allowedPolicyModes = map[string]bool{
	"recommendation_only": true, "operator_controlled": true,
	"guarded_auto_candidate": true, "guarded_auto_sandbox_ready": true,
	"guarded_auto_nonprod_ready": true, "production_canary_blocked": true,
	"production_auto_blocked": true,
}

var candidateGateStatuses = map[string]string{
	"dryRunStatus": "valid", "rollbackStatus": "ready", "sloGuardStatus": "passed",
	"noRegressionStatus": "certified", "donorHealthStatus": "preserved",
	"killSwitchStatus": "clear", "circuitBreakerStatus": "closed",
	"rateLimitStatus": "available", "cooldownStatus": "expired", "blastRadiusStatus": "available",
}

var forbiddenPositiveClaims = []string{
	"auto-apply executed",
	"guaranteed savings active",
	"universal performance improvement",
	"production autonomous apply",
	"autonomous production mutation",
	"windows total ram hotplug supported",
	"logical vcpu hotplug supported",
	"1000 production workloads proven",
	"dashboard is policy source of truth",
	"dashboard policy source of truth",
	"fluidvirt policy authority",
	"inventory applies runtime",
	"dashboard applies runtime",
}

var skipKeys = map[string]bool{
	"forbiddenPhrases": true, "blockerCodes": true, "blockers": true,
	"remediationCodes": true, "reason": true, "decisionReason": true,
}

func requireFalse(doc map[string]interface{}, key string) error {
	if v, ok := doc[key].(bool); !ok || v {
		return fmt.Errorf("%s must be false", key)
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

// ValidateGuardedAutoPolicy checks policy invariants.
func ValidateGuardedAutoPolicy(doc map[string]interface{}) error {
	if err := requireFalse(doc, "autoApplyExecutionEnabled"); err != nil {
		return err
	}
	if err := requireFalse(doc, "productionAutonomousApplyAllowed"); err != nil {
		return err
	}
	if err := requireFalse(doc, "productionScope"); err != nil {
		return err
	}
	if boolOr(doc["allowWindowsEvidenceGatedAuto"]) || boolOr(doc["allowSyntheticProductionProof"]) {
		return fmt.Errorf("allowWindowsEvidenceGatedAuto and allowSyntheticProductionProof must be false")
	}
	return requireClaimBoundary(doc)
}

// ValidateGuardedAutoCandidate checks candidate classification invariants.
func ValidateGuardedAutoCandidate(doc map[string]interface{}) error {
	if err := requireFalse(doc, "autoApplyExecutionEnabled"); err != nil {
		return err
	}
	if err := requireFalse(doc, "productionAutonomousApplyAllowed"); err != nil {
		return err
	}
	if err := requireFalse(doc, "productionScope"); err != nil {
		return err
	}
	state, _ := doc["candidateState"].(string)
	if state == "candidate_only" || state == "sandbox_ready" || state == "nonprod_ready" {
		if !boolOr(doc["rollbackReady"]) {
			return fmt.Errorf("candidate requires rollbackReady=true")
		}
		if !boolOr(doc["sloGuardPassed"]) {
			return fmt.Errorf("candidate requires sloGuardPassed=true")
		}
		if !boolOr(doc["noRegressionCertified"]) {
			return fmt.Errorf("candidate requires noRegressionCertified=true")
		}
		if !boolOr(doc["donorHealthPreserved"]) {
			return fmt.Errorf("candidate requires donorHealthPreserved=true")
		}
		if !boolOr(doc["killSwitchClear"]) {
			return fmt.Errorf("candidate requires killSwitchClear=true")
		}
		if !boolOr(doc["circuitBreakerClosed"]) {
			return fmt.Errorf("candidate requires circuitBreakerClosed=true")
		}
		if !boolOr(doc["rateLimitAvailable"]) {
			return fmt.Errorf("candidate requires rateLimitAvailable=true")
		}
		if !boolOr(doc["cooldownExpired"]) {
			return fmt.Errorf("candidate requires cooldownExpired=true")
		}
		if doc["ledgerRecordRef"] == nil && doc["performanceProofRef"] == nil {
			return fmt.Errorf("candidate requires ledgerRecordRef or performanceProofRef")
		}
		if err := requireEvidenceRefs(doc); err != nil {
			return err
		}
	}
	return requireClaimBoundary(doc)
}

// ValidateEligibilityDecision checks eligibility decision invariants.
func ValidateEligibilityDecision(doc map[string]interface{}) error {
	if err := requireFalse(doc, "autoApplyExecutionEnabled"); err != nil {
		return err
	}
	if err := requireFalse(doc, "productionAutonomousApplyAllowed"); err != nil {
		return err
	}
	if err := requireFalse(doc, "productionScope"); err != nil {
		return err
	}
	if boolOr(doc["eligibleAsProductionCanary"]) {
		return fmt.Errorf("eligibleAsProductionCanary must be false in Sprint 7")
	}
	decision, _ := doc["decision"].(string)
	switch decision {
	case "guarded_auto_candidate":
		if !boolOr(doc["eligibleAsGuardedAutoCandidate"]) {
			return fmt.Errorf("guarded_auto_candidate decision requires eligibleAsGuardedAutoCandidate=true")
		}
		for key, want := range candidateGateStatuses {
			got, _ := doc[key].(string)
			if got != want {
				if key == "receiverHealthStatus" && (got == "neutral_accepted" || got == "preserved") {
					continue
				}
				if key == "ledgerStatus" && (got == "record_present" || got == "not_required") {
					continue
				}
				return fmt.Errorf("guarded_auto_candidate requires %s=%s, got %s", key, want, got)
			}
		}
		if ws, _ := doc["windowsEvidenceGateStatus"].(string); ws == "gated" {
			return fmt.Errorf("Windows evidence-gated action cannot be guarded_auto_candidate")
		}
		if syn, _ := doc["syntheticProductionStatus"].(string); syn == "synthetic_shadow" {
			return fmt.Errorf("synthetic proof cannot be guarded_auto_candidate")
		}
	case "denied_windows_evidence_gated":
		if boolOr(doc["eligibleAsGuardedAutoCandidate"]) {
			return fmt.Errorf("denied_windows_evidence_gated cannot be eligible candidate")
		}
	case "denied_synthetic_production":
		if boolOr(doc["eligibleAsGuardedAutoCandidate"]) {
			return fmt.Errorf("denied_synthetic_production cannot be eligible candidate")
		}
	case "denied_regression", "denied_kill_switch_active", "denied_blast_radius_exceeded":
		if boolOr(doc["eligibleAsGuardedAutoCandidate"]) {
			return fmt.Errorf("%s cannot be eligible candidate", decision)
		}
	case "manual_only":
		if boolOr(doc["eligibleAsGuardedAutoCandidate"]) {
			return fmt.Errorf("manual_only cannot be eligible candidate")
		}
	}
	return requireClaimBoundary(doc)
}

// ValidateGuardedAutoPolicyEngineSurface validates ConfigMap-ready policy engine surface.
func ValidateGuardedAutoPolicyEngineSurface(doc map[string]interface{}) error {
	if doc["milestone"] != MilestoneGuardedAutoPolicyEngine {
		return fmt.Errorf("milestone must be %s", MilestoneGuardedAutoPolicyEngine)
	}
	if err := requireFalse(doc, "autoApplyExecutionEnabled"); err != nil {
		return err
	}
	if err := requireFalse(doc, "productionAutonomousApplyAllowed"); err != nil {
		return err
	}
	if err := requireFalse(doc, "productionScope"); err != nil {
		return err
	}
	mode, _ := doc["policyMode"].(string)
	if mode == "production_auto_with_policy" || !allowedPolicyModes[mode] {
		return fmt.Errorf("invalid or forbidden policyMode: %s", mode)
	}
	inv, ok := doc["safetyInvariants"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("safetyInvariants required")
	}
	for _, key := range []string{"autoApplyExecutionEnabled", "productionAutonomousApplyAllowed", "productionScope", "dashboardPolicySourceOfTruth", "fluidvirtPolicyAuthority", "inventoryRuntimeApply"} {
		if v, ok := inv[key].(bool); !ok || v {
			return fmt.Errorf("safetyInvariants.%s must be false", key)
		}
	}
	sot, ok := doc["sourceOfTruth"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("sourceOfTruth required")
	}
	if !strings.Contains(strings.ToLower(fmt.Sprint(sot["runtimeConstraintEvidence"])), "fluidvirt") {
		return fmt.Errorf("runtime constraint evidence must be FluidVirt")
	}
	if strings.Contains(strings.ToLower(fmt.Sprint(sot["policyProjection"])), "source of truth") {
		return fmt.Errorf("Dashboard must not be policy source of truth")
	}

	var cpuCandidate, sandboxReady, deniedRegression, deniedWindows, deniedKillSwitch, deniedBlastRadius, manualOnly, deniedSynthetic bool
	if cands, ok := doc["candidates"].([]interface{}); ok {
		for _, item := range cands {
			c, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			aid, _ := c["actionId"].(string)
			state, _ := c["candidateState"].(string)
			if aid == "action-cpu-lease-001" && state == "candidate_only" {
				cpuCandidate = true
			}
			if state == "sandbox_ready" {
				sandboxReady = true
			}
			if err := ValidateGuardedAutoCandidate(c); err != nil {
				return fmt.Errorf("candidates: %w", err)
			}
		}
	}
	if decisions, ok := doc["eligibilityDecisions"].([]interface{}); ok {
		for _, item := range decisions {
			d, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			dec, _ := d["decision"].(string)
			switch dec {
			case "guarded_auto_candidate":
				if aid, _ := d["actionId"].(string); aid == "action-cpu-lease-001" {
					cpuCandidate = true
				}
			case "guarded_auto_sandbox_ready":
				sandboxReady = true
			case "denied_regression":
				deniedRegression = true
			case "denied_windows_evidence_gated":
				deniedWindows = true
			case "denied_kill_switch_active":
				deniedKillSwitch = true
			case "denied_blast_radius_exceeded":
				deniedBlastRadius = true
			case "manual_only":
				manualOnly = true
			case "denied_synthetic_production":
				deniedSynthetic = true
			}
			if err := ValidateEligibilityDecision(d); err != nil {
				return fmt.Errorf("eligibilityDecisions: %w", err)
			}
		}
	}
	if policies, ok := doc["policies"].([]interface{}); ok {
		for _, item := range policies {
			p, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			if err := ValidateGuardedAutoPolicy(p); err != nil {
				return fmt.Errorf("policies: %w", err)
			}
		}
	}
	if !cpuCandidate || !sandboxReady || !deniedRegression || !deniedWindows || !deniedKillSwitch || !deniedBlastRadius || !manualOnly || !deniedSynthetic {
		return fmt.Errorf("surface must include CPU candidate, sandbox-ready, denied regression, Windows gated, kill switch, blast radius, manual-only, and synthetic denied scenarios")
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
	case "claimBoundary", "claimBoundaries", "allowedPhrases", "conditionalPhrases":
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

// ValidateSprint7Examples validates all Sprint 7 reference examples.
func ValidateSprint7Examples(repoRoot string) error {
	files := map[string]func(map[string]interface{}) error{
		"guarded-auto-policy-reference.json":            ValidateGuardedAutoPolicy,
		"guarded-auto-eligibility-decision-reference.json": ValidateEligibilityDecision,
		"guarded-auto-candidate-reference.json":           ValidateGuardedAutoCandidate,
		"blast-radius-budget-reference.json":              func(d map[string]interface{}) error { return requireClaimBoundary(d) },
		"kill-switch-state-reference.json":                func(d map[string]interface{}) error { return requireClaimBoundary(d) },
		"circuit-breaker-state-reference.json":            func(d map[string]interface{}) error { return requireClaimBoundary(d) },
		"auto-apply-rate-limit-reference.json":            func(d map[string]interface{}) error { return requireClaimBoundary(d) },
		"auto-apply-cooldown-reference.json":              func(d map[string]interface{}) error { return requireClaimBoundary(d) },
		"policy-scope-reference.json":                     func(d map[string]interface{}) error { return requireClaimBoundary(d) },
		"policy-denial-reason-reference.json":             func(d map[string]interface{}) error { return requireClaimBoundary(d) },
		"guarded-auto-audit-requirement-reference.json":   func(d map[string]interface{}) error { return requireClaimBoundary(d) },
		"guarded-auto-policy-engine-reference.json":       ValidateGuardedAutoPolicyEngineSurface,
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
		if name != "guarded-auto-policy-engine-reference.json" {
			if err := rejectForbiddenPositiveClaims(doc); err != nil {
				return fmt.Errorf("%s: %w", name, err)
			}
		}
	}
	return nil
}

// SchemaFilesRequiredSprint7 returns Sprint 7 schema basenames.
func SchemaFilesRequiredSprint7() []string {
	return []string{
		"guarded-auto-policy-engine-v1.schema.json",
		"guarded-auto-policy-v1.schema.json",
		"guarded-auto-eligibility-decision-v1.schema.json",
		"guarded-auto-candidate-v1.schema.json",
		"blast-radius-budget-v1.schema.json",
		"kill-switch-state-v1.schema.json",
		"circuit-breaker-state-v1.schema.json",
		"auto-apply-rate-limit-v1.schema.json",
		"auto-apply-cooldown-v1.schema.json",
		"policy-scope-v1.schema.json",
		"policy-denial-reason-v1.schema.json",
		"guarded-auto-audit-requirement-v1.schema.json",
	}
}
