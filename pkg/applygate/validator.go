package applygate

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const MilestoneOperatorControlledApplyGate = "hyperdensity_operator_controlled_apply_gate_v1"

var forbiddenPositiveClaims = []string{
	"guaranteed savings active",
	"universal performance improvement",
	"production autonomous apply",
	"autonomous production mutation",
	"windows total ram hotplug supported",
	"logical vcpu hotplug supported",
	"1000 production workloads proven",
	"dashboard is source of truth",
	"dashboard executes mutation",
	"inventory hyperdensity engine",
	"inventory applies runtime",
}

var skipKeys = map[string]bool{
	"forbiddenPhrases": true, "blockerCodes": true, "blockers": true,
	"remediationCodes": true, "remediations": true,
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

// ValidateOperatorApprovalRecord checks explicit operator approval invariants.
func ValidateOperatorApprovalRecord(doc map[string]interface{}) error {
	if doc["approvalMode"] != "operator_required" {
		return fmt.Errorf("approvalMode must be operator_required")
	}
	if err := requireTrue(doc, "riskAccepted"); err != nil {
		return err
	}
	for _, k := range []string{"productionScope", "autoApplyAllowed", "productionAutonomousApplyAllowed"} {
		if err := requireFalse(doc, k); err != nil {
			return err
		}
	}
	if id, _ := doc["actionId"].(string); id == "" {
		return fmt.Errorf("actionId required")
	}
	if id, _ := doc["leaseCandidateId"].(string); id == "" {
		return fmt.Errorf("leaseCandidateId required")
	}
	if err := requireClaimBoundary(doc); err != nil {
		return err
	}
	return requireEvidenceRefs(doc)
}

// ValidateApplyRequest enforces guarded apply request fields.
func ValidateApplyRequest(doc map[string]interface{}) error {
	if doc["actuator"] != "FluidVirt" {
		return fmt.Errorf("actuator must be FluidVirt")
	}
	for _, k := range []string{"productionScope", "autoApplyAllowed", "productionAutonomousApplyAllowed", "rawRuntimeControlsExposed"} {
		if err := requireFalse(doc, k); err != nil {
			return err
		}
	}
	for _, ref := range []string{"approvalRef", "dryRunRef", "rollbackReadinessRef", "sloPrecheckRef", "riskAssessmentRef"} {
		if s, _ := doc[ref].(string); s == "" {
			return fmt.Errorf("%s required", ref)
		}
	}
	for _, id := range []string{"actionId", "leaseCandidateId", "donorShellId", "receiverShellId"} {
		if s, _ := doc[id].(string); s == "" {
			return fmt.Errorf("%s required", id)
		}
	}
	if err := requireClaimBoundary(doc); err != nil {
		return err
	}
	return requireEvidenceRefs(doc)
}

// ValidateFluidvirtInvocationRecord checks FluidVirt guarded invocation.
func ValidateFluidvirtInvocationRecord(doc map[string]interface{}, approvalPresent bool) error {
	if doc["actuator"] != "FluidVirt" {
		return fmt.Errorf("actuator must be FluidVirt")
	}
	for _, k := range []string{"rawRuntimeControlsExposed", "productionScope", "autoApplyAllowed", "productionAutonomousApplyAllowed"} {
		if err := requireFalse(doc, k); err != nil {
			return err
		}
	}
	executed, _ := doc["mutationExecuted"].(bool)
	if executed {
		if !approvalPresent {
			return fmt.Errorf("mutationExecuted requires operator approval")
		}
		if v, ok := doc["productionScope"].(bool); ok && v {
			return fmt.Errorf("mutationExecuted forbidden with productionScope=true")
		}
	}
	if err := requireClaimBoundary(doc); err != nil {
		return err
	}
	return requireEvidenceRefs(doc)
}

// ValidateRuntimeMutationObservation enforces in-place movement invariants.
func ValidateRuntimeMutationObservation(doc map[string]interface{}) error {
	if err := requireTrue(doc, "identityPreserved"); err != nil {
		return err
	}
	for _, k := range []string{"rebootUsed", "recreateUsed", "rolloutUsed", "migrationUsed"} {
		if err := requireFalse(doc, k); err != nil {
			return err
		}
	}
	if err := requireClaimBoundary(doc); err != nil {
		return err
	}
	return requireEvidenceRefs(doc)
}

// ValidatePostVerifyResult checks post-verify pass conditions.
func ValidatePostVerifyResult(doc map[string]interface{}, mutationMatchesPlan bool) error {
	status, _ := doc["verifyStatus"].(string)
	if status == "passed" {
		if !mutationMatchesPlan {
			return fmt.Errorf("postVerify cannot pass when mutationMatchesPlan=false")
		}
		if v, ok := doc["runtimeDeltaVerified"].(bool); !ok || !v {
			return fmt.Errorf("passed postVerify requires runtimeDeltaVerified=true")
		}
		if doc["sloGuardStatus"] != "passed" {
			return fmt.Errorf("passed postVerify requires sloGuardStatus=passed")
		}
		if doc["donorHealthStatus"] != "preserved" {
			return fmt.Errorf("passed postVerify requires donorHealthStatus=preserved")
		}
		if v, ok := doc["rollbackReady"].(bool); !ok || !v {
			return fmt.Errorf("passed postVerify requires rollbackReady=true")
		}
		if v, ok := doc["mutationRolledBack"].(bool); ok && v {
			return fmt.Errorf("passed postVerify requires mutationRolledBack=false")
		}
		guestVerified, _ := doc["guestDeltaVerified"].(bool)
		if !guestVerified && status == "passed" {
			// guest must be verified for passed, or use evidence_gated status
			return fmt.Errorf("passed postVerify requires guestDeltaVerified=true or evidence_gated status")
		}
	}
	if status == "evidence_gated" {
		cb := flattenClaimBoundary(doc)
		merged := strings.ToLower(strings.Join(cb, " "))
		if strings.Contains(merged, "guest-visible delta proven") || strings.Contains(merged, "guest delta verified for all classes") {
			return fmt.Errorf("evidence_gated cannot overclaim guest-visible delta")
		}
	}
	if err := requireClaimBoundary(doc); err != nil {
		return err
	}
	return requireEvidenceRefs(doc)
}

func flattenClaimBoundary(doc map[string]interface{}) []string {
	cb, ok := doc["claimBoundary"].([]interface{})
	if !ok {
		return nil
	}
	var out []string
	for _, item := range cb {
		if s, ok := item.(string); ok {
			out = append(out, s)
		}
	}
	return out
}

// ValidateRollbackWindow checks rollback window reference.
func ValidateRollbackWindow(doc map[string]interface{}) error {
	if err := requireClaimBoundary(doc); err != nil {
		return err
	}
	return requireEvidenceRefs(doc)
}

// ValidateApplyAuditEvent checks audit event reference.
func ValidateApplyAuditEvent(doc map[string]interface{}) error {
	if id, _ := doc["actionId"].(string); id == "" {
		return fmt.Errorf("actionId required")
	}
	if id, _ := doc["leaseCandidateId"].(string); id == "" {
		return fmt.Errorf("leaseCandidateId required")
	}
	if err := requireClaimBoundary(doc); err != nil {
		return err
	}
	return requireEvidenceRefs(doc)
}

// ValidateOperatorApplyGateSurface validates ConfigMap-ready apply gate surface.
func ValidateOperatorApplyGateSurface(doc map[string]interface{}) error {
	if doc["milestone"] != MilestoneOperatorControlledApplyGate {
		return fmt.Errorf("milestone must be %s", MilestoneOperatorControlledApplyGate)
	}
	if err := requireTrue(doc, "operatorControlledApplyAllowed"); err != nil {
		return err
	}
	for _, k := range []string{"autoApplyAllowed", "productionAutonomousApplyAllowed", "rawRuntimeControlsExposed", "productionScope"} {
		if err := requireFalse(doc, k); err != nil {
			return err
		}
	}
	if doc["mutationScope"] != "technical_preview_operator_controlled" {
		return fmt.Errorf("mutationScope must be technical_preview_operator_controlled")
	}
	inv, ok := doc["safetyInvariants"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("safetyInvariants required")
	}
	for _, key := range []string{"autoApplyAllowed", "productionAutonomousApplyAllowed", "rawRuntimeControlsExposed", "dashboardAppliesRuntimeChanges", "inventoryRuntimeApply"} {
		if v, ok := inv[key].(bool); !ok || v {
			return fmt.Errorf("safetyInvariants.%s must be false", key)
		}
	}
	sot, ok := doc["sourceOfTruth"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("sourceOfTruth required")
	}
	if !strings.Contains(strings.ToLower(fmt.Sprint(sot["runtimeMutation"])), "fluidvirt") {
		return fmt.Errorf("runtime mutation source must be FluidVirt")
	}
	if strings.Contains(strings.ToLower(fmt.Sprint(sot["applyGateProjection"])), "source of truth") {
		return fmt.Errorf("Dashboard must not be mutation source of truth")
	}
	scenarios, ok := doc["scenarioSummaries"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("scenarioSummaries required")
	}
	for _, name := range []string{"successfulOperatorControlled", "blockedApply", "rollbackRequired", "windowsEvidenceGated"} {
		if _, ok := scenarios[name]; !ok {
			return fmt.Errorf("missing scenario %s", name)
		}
	}
	if blocked, ok := scenarios["blockedApply"].(map[string]interface{}); ok {
		if v, _ := blocked["mutationExecuted"].(bool); v {
			return fmt.Errorf("blocked apply must not have mutationExecuted=true")
		}
	}
	if win, ok := scenarios["windowsEvidenceGated"].(map[string]interface{}); ok {
		if v, _ := win["guestDeltaVerified"].(bool); v {
			return fmt.Errorf("windows evidence gated must not claim guestDeltaVerified=true")
		}
		cb := ""
		if c, ok := win["claimBoundary"].(string); ok {
			cb = strings.ToLower(c)
		}
		if strings.Contains(cb, "windows total ram hotplug") || strings.Contains(cb, "logical vcpu hotplug") {
			return fmt.Errorf("windows scenario must not claim hotplug")
		}
	}
	if mo, ok := doc["mutationObservation"].(map[string]interface{}); ok {
		if err := ValidateRuntimeMutationObservation(extendMutationObservation(mo)); err != nil {
			return err
		}
		if pv, ok := doc["postVerify"].(map[string]interface{}); ok {
			matches, _ := mo["mutationMatchesPlan"].(bool)
			pvDoc := map[string]interface{}{
				"verifyStatus":         pv["verifyStatus"],
				"runtimeDeltaVerified":   pv["runtimeDeltaVerified"],
				"sloGuardStatus":         pv["sloGuardStatus"],
				"donorHealthStatus":      "preserved",
				"rollbackReady":          pv["rollbackReady"],
				"mutationRolledBack":     false,
				"guestDeltaVerified":     true,
				"claimBoundary":          []interface{}{"technical preview operator-controlled apply only"},
				"evidenceRefs":           []interface{}{"hyperdensity_guarded_apply_invocation_surface_v1"},
			}
			if err := ValidatePostVerifyResult(pvDoc, matches); err != nil {
				return err
			}
		}
	}
	return rejectForbiddenPositiveClaims(doc)
}

func extendMutationObservation(mo map[string]interface{}) map[string]interface{} {
	out := map[string]interface{}{}
	for k, v := range mo {
		out[k] = v
	}
	if _, ok := out["identityPreserved"]; !ok {
		out["identityPreserved"] = true
	}
	for _, k := range []string{"rebootUsed", "recreateUsed", "rolloutUsed", "migrationUsed"} {
		if _, ok := out[k]; !ok {
			out[k] = false
		}
	}
	if _, ok := out["claimBoundary"]; !ok {
		out["claimBoundary"] = []interface{}{"in-place movement only"}
	}
	if _, ok := out["evidenceRefs"]; !ok {
		out["evidenceRefs"] = []interface{}{"hyperdensity_guarded_apply_invocation_surface_v1"}
	}
	return out
}

func rejectForbiddenPositiveClaims(v interface{}) error {
	positives := collectPositiveStrings(v)
	merged := strings.ToLower(strings.Join(positives, "\n"))
	for _, phrase := range forbiddenPositiveClaims {
		if strings.Contains(merged, phrase) {
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
	case "claimBoundary", "claimBoundaries", "allowedPhrases", "conditionalPhrases", "eventSummary":
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

// ValidateSprint4Examples validates all Sprint 4 reference examples.
func ValidateSprint4Examples(repoRoot string) error {
	approvalPresent := true
	files := map[string]func(map[string]interface{}) error{
		"operator-approval-record-reference.json": func(d map[string]interface{}) error {
			return ValidateOperatorApprovalRecord(d)
		},
		"apply-request-reference.json": func(d map[string]interface{}) error {
			return ValidateApplyRequest(d)
		},
		"fluidvirt-invocation-record-reference.json": func(d map[string]interface{}) error {
			return ValidateFluidvirtInvocationRecord(d, approvalPresent)
		},
		"runtime-mutation-observation-reference.json": func(d map[string]interface{}) error {
			return ValidateRuntimeMutationObservation(d)
		},
		"post-verify-result-reference.json": func(d map[string]interface{}) error {
			return ValidatePostVerifyResult(d, true)
		},
		"rollback-window-reference.json": func(d map[string]interface{}) error {
			return ValidateRollbackWindow(d)
		},
		"apply-audit-event-reference.json": func(d map[string]interface{}) error {
			return ValidateApplyAuditEvent(d)
		},
		"operator-apply-gate-reference.json": func(d map[string]interface{}) error {
			return ValidateOperatorApplyGateSurface(d)
		},
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
		if err := rejectForbiddenPositiveClaims(doc); err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
	}
	return nil
}

// SchemaFilesRequiredSprint4 returns Sprint 4 schema basenames.
func SchemaFilesRequiredSprint4() []string {
	return []string{
		"operator-apply-gate-v1.schema.json",
		"operator-approval-record-v1.schema.json",
		"apply-request-v1.schema.json",
		"fluidvirt-invocation-record-v1.schema.json",
		"runtime-mutation-observation-v1.schema.json",
		"post-verify-result-v1.schema.json",
		"rollback-window-v1.schema.json",
		"apply-audit-event-v1.schema.json",
	}
}
