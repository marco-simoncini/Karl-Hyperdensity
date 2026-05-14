package leaseslate

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	MilestoneResourceLeaseActionSlateReadiness = "hyperdensity_resource_lease_action_slate_readiness_v1"
	SlateID                                    = "hyperdensity_action_slate_readiness_v1"
	DonorIndexID                               = "hyperdensity_donor_index_v1"
	ReceiverIndexID                            = "hyperdensity_receiver_index_v1"
)

var forbiddenPositiveClaims = []string{
	"guaranteed savings active",
	"universal performance improvement",
	"production autonomous apply",
	"windows total ram hotplug supported",
	"logical vcpu hotplug supported",
	"1000 production workloads proven",
	"dashboard is source of truth",
	"inventory hyperdensity engine",
	"action slate auto-applies",
	"dashboard applies runtime changes",
}

var skipKeys = map[string]bool{
	"forbiddenPhrases": true, "blockerCodes": true, "blockers": true,
	"remediationCodes": true, "remediations": true, "unsupportedFamilies": true,
	"excludedShellKinds": true,
}

// ValidateDonorIndex checks donor index reference invariants.
func ValidateDonorIndex(doc map[string]interface{}) error {
	if doc["indexId"] != DonorIndexID {
		return fmt.Errorf("indexId must be %s", DonorIndexID)
	}
	donors, ok := doc["donors"].([]interface{})
	if !ok {
		return fmt.Errorf("donors array required")
	}
	for _, item := range donors {
		d, ok := item.(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid donor entry")
		}
		if err := validateIndexShellFields(d, "donor"); err != nil {
			return err
		}
		if rr, _ := d["rollbackReadiness"].(string); rr == "" || rr == "unknown" {
			return fmt.Errorf("donor %v rollbackReadiness must not be unknown", d["shellId"])
		}
	}
	return rejectForbiddenPositiveClaims(doc)
}

// ValidateReceiverIndex checks receiver index reference invariants.
func ValidateReceiverIndex(doc map[string]interface{}) error {
	if doc["indexId"] != ReceiverIndexID {
		return fmt.Errorf("indexId must be %s", ReceiverIndexID)
	}
	receivers, ok := doc["receivers"].([]interface{})
	if !ok {
		return fmt.Errorf("receivers array required")
	}
	for _, item := range receivers {
		r, ok := item.(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid receiver entry")
		}
		if err := validateIndexShellFields(r, "receiver"); err != nil {
			return err
		}
	}
	return rejectForbiddenPositiveClaims(doc)
}

func validateIndexShellFields(entry map[string]interface{}, role string) error {
	if id, _ := entry["shellId"].(string); id == "" {
		return fmt.Errorf("%s missing shellId", role)
	}
	cb, ok := entry["claimBoundary"].([]interface{})
	if !ok || len(cb) == 0 {
		return fmt.Errorf("%s missing claimBoundary", role)
	}
	refs, ok := entry["evidenceRefs"].([]interface{})
	if !ok || len(refs) == 0 {
		return fmt.Errorf("%s missing evidenceRefs", role)
	}
	return nil
}

// ValidateResourceLeaseCandidate enforces Sprint 3 lease gates.
func ValidateResourceLeaseCandidate(doc map[string]interface{}) error {
	if v, ok := doc["autoApplyAllowed"].(bool); !ok || v {
		return fmt.Errorf("lease candidate autoApplyAllowed must be false")
	}
	if v, ok := doc["productionMutationAllowed"].(bool); !ok || v {
		return fmt.Errorf("lease candidate productionMutationAllowed must be false")
	}
	cb, ok := doc["claimBoundary"].([]interface{})
	if !ok || len(cb) == 0 {
		return fmt.Errorf("lease candidate missing claimBoundary")
	}
	refs, ok := doc["evidenceRefs"].([]interface{})
	if !ok || len(refs) == 0 {
		return fmt.Errorf("lease candidate missing evidenceRefs")
	}
	return rejectForbiddenPositiveClaims(doc)
}

// ValidateActionSlateEntry enforces per-action readiness invariants.
func ValidateActionSlateEntry(doc map[string]interface{}) error {
	if v, ok := doc["autoApplyAllowed"].(bool); !ok || v {
		return fmt.Errorf("action entry autoApplyAllowed must be false")
	}
	if v, ok := doc["productionMutationAllowed"].(bool); !ok || v {
		return fmt.Errorf("action entry productionMutationAllowed must be false")
	}
	dry, _ := doc["dryRunStatus"].(string)
	if dry == "" {
		return fmt.Errorf("action entry missing dryRunStatus")
	}
	rb, _ := doc["rollbackStatus"].(string)
	if rb == "" {
		return fmt.Errorf("action entry missing rollbackStatus")
	}
	cb, ok := doc["claimBoundary"].([]interface{})
	if !ok || len(cb) == 0 {
		return fmt.Errorf("action entry missing claimBoundary")
	}
	refs, ok := doc["evidenceRefs"].([]interface{})
	if !ok || len(refs) == 0 {
		return fmt.Errorf("action entry missing evidenceRefs")
	}
	ready, _ := doc["readyForApply"].(string)
	blockers, _ := doc["blockers"].([]interface{})
	if ready == "operator_controlled_ready" || ready == "true" {
		if len(blockers) > 0 {
			return fmt.Errorf("action with blockers cannot be operator_controlled_ready")
		}
		if ready == "true" {
			return fmt.Errorf("readyForApply must be operator_controlled_ready not bare true")
		}
	}
	if ready == "operator_controlled_ready" && (dry == "blocked" || strings.EqualFold(rb, "unknown")) {
		return fmt.Errorf("blocked/unknown rollback cannot be operator_controlled_ready")
	}
	return rejectForbiddenPositiveClaims(doc)
}

// ValidateActionSlateReadiness validates ConfigMap-ready action slate surface.
func ValidateActionSlateReadiness(doc map[string]interface{}) error {
	if doc["milestone"] != MilestoneResourceLeaseActionSlateReadiness {
		return fmt.Errorf("milestone must be %s", MilestoneResourceLeaseActionSlateReadiness)
	}
	if v, ok := doc["noFullNxNPairing"].(bool); !ok || !v {
		return fmt.Errorf("noFullNxNPairing must be true")
	}
	evaluated := intFrom(doc["evaluatedPairCount"])
	maxPairs := intFrom(doc["maxEvaluatedPairs"])
	if maxPairs > 0 && evaluated > maxPairs {
		return fmt.Errorf("evaluatedPairCount %d exceeds maxEvaluatedPairs %d", evaluated, maxPairs)
	}
	full := intFrom(doc["fullPairSpace"])
	if evaluated > 0 && full > 0 && evaluated >= full && doc["noFullNxNPairing"] == true {
		// allow equal only when fullPairSpace equals evaluated (degenerate); require avoided count or topK
		avoided := intFrom(doc["avoidedPairCount"])
		if avoided == 0 && full > evaluated {
			return fmt.Errorf("noFullNxNPairing requires avoided pairs when fullPairSpace > evaluatedPairCount")
		}
	}
	if full > evaluated && intFrom(doc["avoidedPairCount"]) < 1 {
		return fmt.Errorf("avoidedPairCount required when not evaluating full pair space")
	}
	inv, ok := doc["safetyInvariants"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("safetyInvariants required")
	}
	for _, key := range []string{"autoApplyAllowed", "productionMutationAllowed", "guaranteedSavingsClaimed",
		"universalPerformanceImprovementClaimed", "windowsTotalRamHotplugClaimed", "logicalVcpuHotplugClaimed",
		"dashboardAppliesRuntimeChanges"} {
		if v, ok := inv[key].(bool); !ok || v {
			return fmt.Errorf("safetyInvariants.%s must be false", key)
		}
	}
	sot, ok := doc["sourceOfTruth"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("sourceOfTruth required")
	}
	if !strings.Contains(strings.ToLower(fmt.Sprint(sot["dryRunReadinessEvidence"])), "fluidvirt") {
		return fmt.Errorf("dry-run readiness evidence source must be FluidVirt")
	}
	if strings.Contains(strings.ToLower(fmt.Sprint(sot["actionSlateProjection"])), "source of truth") {
		return fmt.Errorf("Dashboard must not be source of truth")
	}
	entries, ok := doc["actionEntries"].([]interface{})
	if !ok || len(entries) == 0 {
		return fmt.Errorf("actionEntries required")
	}
	var readyOp, blocked, remediation int
	var hasCPU, hasMem, hasWindowsRemediation bool
	for _, item := range entries {
		ae, ok := item.(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid action entry")
		}
		if err := ValidateActionSlateEntry(ae); err != nil {
			return err
		}
		ready, _ := ae["readyForApply"].(string)
		switch ready {
		case "operator_controlled_ready":
			readyOp++
		case "blocked":
			blocked++
		case "remediation_only":
			remediation++
		}
		res, _ := ae["resource"].(string)
		if res == "cpu" {
			hasCPU = true
		}
		if res == "memory_envelope" {
			hasMem = true
		}
		if ready == "remediation_only" && strings.Contains(fmt.Sprint(ae["receiverShellId"]), "windows") {
			hasWindowsRemediation = true
		}
	}
	if readyOp < 1 || blocked < 1 || remediation < 1 {
		return fmt.Errorf("surface must include operator-ready, blocked, and remediation-only actions")
	}
	if !hasCPU || !hasMem {
		return fmt.Errorf("surface must include cpu and memory_envelope lease actions")
	}
	if !hasWindowsRemediation {
		return fmt.Errorf("surface must include Windows remediation-only action")
	}
	if intFrom(doc["protectedExcludedCount"]) < 1 {
		return fmt.Errorf("protectedExcludedCount must be >= 1")
	}
	if candidates, ok := doc["leaseCandidates"].([]interface{}); ok {
		for _, item := range candidates {
			c, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			if err := ValidateResourceLeaseCandidate(c); err != nil {
				return err
			}
		}
	}
	return rejectForbiddenPositiveClaims(doc)
}

// ValidateActionDryrunReadiness ensures FluidVirt dry-run evidence only.
func ValidateActionDryrunReadiness(doc map[string]interface{}) error {
	if doc["source"] != "FluidVirt" || doc["actuator"] != "FluidVirt" {
		return fmt.Errorf("dry-run readiness must be FluidVirt-sourced")
	}
	if v, ok := doc["mutationExecuted"].(bool); !ok || v {
		return fmt.Errorf("mutationExecuted must be false")
	}
	if dry, _ := doc["dryRunStatus"].(string); dry == "" {
		return fmt.Errorf("missing dryRunStatus")
	}
	return rejectForbiddenPositiveClaims(doc)
}

// ValidateRollbackReadiness checks rollback evidence.
func ValidateRollbackReadiness(doc map[string]interface{}) error {
	if status, _ := doc["rollbackStatus"].(string); status == "" {
		return fmt.Errorf("missing rollbackStatus")
	}
	refs, ok := doc["evidenceRefs"].([]interface{})
	if !ok || len(refs) == 0 {
		return fmt.Errorf("missing evidenceRefs")
	}
	return nil
}

// ValidateSloPrecheck checks SLO precheck placeholder.
func ValidateSloPrecheck(doc map[string]interface{}) error {
	if status, _ := doc["precheckStatus"].(string); status == "" {
		return fmt.Errorf("missing precheckStatus")
	}
	return nil
}

// ValidateRiskAssessment checks risk assessment reference.
func ValidateRiskAssessment(doc map[string]interface{}) error {
	if status, _ := doc["riskStatus"].(string); status == "" {
		return fmt.Errorf("missing riskStatus")
	}
	return nil
}

func intFrom(v interface{}) int {
	switch t := v.(type) {
	case float64:
		return int(t)
	case int:
		return t
	default:
		return 0
	}
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

// ValidateSprint3Examples validates all Sprint 3 reference examples.
func ValidateSprint3Examples(repoRoot string) error {
	files := map[string]func(map[string]interface{}) error{
		"donor-index-reference.json":              ValidateDonorIndex,
		"receiver-index-reference.json":           ValidateReceiverIndex,
		"resource-lease-candidate-reference.json": ValidateResourceLeaseCandidate,
		"action-slate-entry-reference.json":       ValidateActionSlateEntry,
		"action-slate-readiness-reference.json":   ValidateActionSlateReadiness,
		"action-dryrun-readiness-reference.json":  ValidateActionDryrunReadiness,
		"rollback-readiness-reference.json":       ValidateRollbackReadiness,
		"slo-precheck-reference.json":             ValidateSloPrecheck,
		"risk-assessment-reference.json":          ValidateRiskAssessment,
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
	}
	return nil
}

// SchemaFilesRequiredSprint3 returns Sprint 3 schema basenames.
func SchemaFilesRequiredSprint3() []string {
	return []string{
		"donor-index-v1.schema.json",
		"receiver-index-v1.schema.json",
		"resource-lease-candidate-v1.schema.json",
		"action-slate-readiness-v1.schema.json",
		"action-slate-entry-v1.schema.json",
		"action-dryrun-readiness-v1.schema.json",
		"rollback-readiness-v1.schema.json",
		"slo-precheck-v1.schema.json",
		"risk-assessment-v1.schema.json",
	}
}
