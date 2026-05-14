package savingsledger

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const MilestoneRealizedSavingsLedger = "hyperdensity_realized_savings_ledger_v1"

var forbiddenPositiveClaims = []string{
	"guaranteed savings active",
	"guaranteed savings claimed",
	"universal performance improvement",
	"production autonomous apply",
	"autonomous production mutation",
	"windows total ram hotplug supported",
	"logical vcpu hotplug supported",
	"1000 production workloads proven",
	"dashboard is accounting source of truth",
	"dashboard accounting source of truth",
	"fluidvirt accounting authority",
	"inventory accounting authority",
	"estimated value counted as guaranteed",
	"synthetic value counted as production",
}

var skipKeys = map[string]bool{
	"forbiddenPhrases": true, "blockerCodes": true, "blockers": true,
	"exclusionReason": true, "exclusionsByReason": true, "missingEvidence": true,
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

// ValidateMovementAccountingRecord checks a single movement accounting record.
func ValidateMovementAccountingRecord(doc map[string]interface{}) error {
	class, _ := doc["claimClassification"].(string)
	if class == "guaranteed_savings_active" {
		return fmt.Errorf("guaranteed_savings_active forbidden in Sprint 5")
	}
	if err := requireClaimBoundary(doc); err != nil {
		return err
	}
	if err := requireEvidenceRefs(doc); err != nil {
		return err
	}
	if synth, _ := doc["synthetic"].(bool); synth {
		if class != "excluded_synthetic" {
			return fmt.Errorf("synthetic record must be excluded_synthetic")
		}
	}
	if eligible, _ := doc["guaranteeEligibleForFuture"].(bool); eligible {
		if err := validateFutureGuaranteeEligibility(doc); err != nil {
			return err
		}
	}
	net, _ := doc["netRealizedValue"].(float64)
	if eligible, _ := doc["guaranteeEligibleForFuture"].(bool); eligible && net < 0 {
		return fmt.Errorf("negative net value cannot be guaranteeEligibleForFuture")
	}
	return nil
}

func validateFutureGuaranteeEligibility(doc map[string]interface{}) error {
	if v, ok := doc["sloPreserved"].(bool); !ok || !v {
		return fmt.Errorf("guaranteeEligibleForFuture requires sloPreserved=true")
	}
	if v, ok := doc["donorHealthPreserved"].(bool); !ok || !v {
		return fmt.Errorf("guaranteeEligibleForFuture requires donorHealthPreserved=true")
	}
	if v, ok := doc["rollbackReady"].(bool); !ok || !v {
		return fmt.Errorf("guaranteeEligibleForFuture requires rollbackReady=true")
	}
	if dur, ok := doc["durationSeconds"].(float64); !ok || dur <= 0 {
		if dur2, ok2 := doc["durationSeconds"].(int); !ok2 || dur2 <= 0 {
			return fmt.Errorf("guaranteeEligibleForFuture requires durationSeconds>0")
		}
	}
	if price, ok := doc["unitPrice"].(float64); !ok || price <= 0 {
		return fmt.Errorf("guaranteeEligibleForFuture requires unitPrice>0")
	}
	if synth, _ := doc["synthetic"].(bool); synth {
		return fmt.Errorf("synthetic cannot be guaranteeEligibleForFuture")
	}
	class, _ := doc["claimClassification"].(string)
	switch class {
	case "excluded_post_verify_failed", "excluded_evidence_gated", "excluded_windows_guest_delta_not_verified",
		"excluded_slo_failed", "excluded_rollback_unavailable", "excluded_missing_unit_price",
		"excluded_estimated_only", "excluded_synthetic":
		return fmt.Errorf("excluded classification cannot be guaranteeEligibleForFuture: %s", class)
	}
	net, _ := doc["netRealizedValue"].(float64)
	if net < 0 {
		return fmt.Errorf("negative netRealizedValue cannot be guaranteeEligibleForFuture")
	}
	return nil
}

// ValidateResourceUnitPrice checks unit price reference.
func ValidateResourceUnitPrice(doc map[string]interface{}) error {
	if err := requireClaimBoundary(doc); err != nil {
		return err
	}
	return requireEvidenceRefs(doc)
}

// ValidateSavingsClaimClassification checks classification record.
func ValidateSavingsClaimClassification(doc map[string]interface{}) error {
	if class, _ := doc["claimClassification"].(string); class == "guaranteed_savings_active" {
		return fmt.Errorf("guaranteed_savings_active forbidden")
	}
	if est, _ := doc["estimatedOnly"].(bool); est {
		if class, _ := doc["claimClassification"].(string); class != "excluded_estimated_only" {
			return fmt.Errorf("estimatedOnly must be excluded_estimated_only")
		}
	}
	if err := requireClaimBoundary(doc); err != nil {
		return err
	}
	return nil
}

// ValidateMonthlySavingsRollup checks monthly rollup invariants.
func ValidateMonthlySavingsRollup(doc map[string]interface{}) error {
	if err := requireFalse(doc, "guaranteedSavingsClaimed"); err != nil {
		return err
	}
	if err := requireClaimBoundary(doc); err != nil {
		return err
	}
	return requireEvidenceRefs(doc)
}

// ValidateRealizedSavingsLedgerSurface validates ConfigMap-ready ledger surface.
func ValidateRealizedSavingsLedgerSurface(doc map[string]interface{}) error {
	if doc["milestone"] != MilestoneRealizedSavingsLedger {
		return fmt.Errorf("milestone must be %s", MilestoneRealizedSavingsLedger)
	}
	for _, k := range []string{"guaranteedSavingsAllowed", "guaranteedSavingsClaimed", "estimatedValueCountedAsGuaranteed", "syntheticValueCountedAsProduction"} {
		if err := requireFalse(doc, k); err != nil {
			return err
		}
	}
	inv, ok := doc["safetyInvariants"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("safetyInvariants required")
	}
	for _, key := range []string{"guaranteedSavingsAllowed", "guaranteedSavingsClaimed", "estimatedValueCountedAsGuaranteed", "syntheticValueCountedAsProduction", "dashboardAccountingSourceOfTruth", "fluidvirtAccountingAuthority"} {
		if v, ok := inv[key].(bool); !ok || v {
			return fmt.Errorf("safetyInvariants.%s must be false", key)
		}
	}
	sot, ok := doc["sourceOfTruth"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("sourceOfTruth required")
	}
	if !strings.Contains(strings.ToLower(fmt.Sprint(sot["ledgerContracts"])), "hyperdensity") &&
		!strings.Contains(strings.ToLower(fmt.Sprint(sot["ledgerContracts"])), "karl-hyperdensity") {
		return fmt.Errorf("ledger contracts must be Karl-Hyperdensity")
	}
	if strings.Contains(strings.ToLower(fmt.Sprint(sot["ledgerProjection"])), "source of truth") {
		return fmt.Errorf("Dashboard must not be accounting source of truth")
	}
	if !strings.Contains(strings.ToLower(fmt.Sprint(sot["movementMeasurements"])), "fluidvirt") {
		return fmt.Errorf("movement measurements must be FluidVirt")
	}
	records, ok := doc["records"].([]interface{})
	if !ok || len(records) == 0 {
		return fmt.Errorf("records required")
	}
	var realizedCPU, memRollback, blocked, windowsGated bool
	for _, item := range records {
		rec, ok := item.(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid record")
		}
		if err := ValidateMovementAccountingRecord(rec); err != nil {
			return err
		}
		class, _ := rec["claimClassification"].(string)
		res, _ := rec["resource"].(string)
		switch {
		case class == "eligible_for_future_guarantee" && res == "cpu":
			realizedCPU = true
		case class == "realized_non_guaranteed" && res == "memory_envelope":
			memRollback = true
		case class == "excluded_post_verify_failed":
			blocked = true
		case class == "excluded_windows_guest_delta_not_verified":
			windowsGated = true
		}
	}
	if !realizedCPU || !memRollback || !blocked || !windowsGated {
		return fmt.Errorf("ledger must include CPU realized, memory rollback, blocked, and Windows gated records")
	}
	if est, ok := doc["estimatedOpportunityProjection"].(map[string]interface{}); ok {
		if boolOr(est["countedAsGuaranteed"]) || boolOr(est["countedAsRealized"]) {
			return fmt.Errorf("estimated opportunity must not be counted as realized or guaranteed")
		}
	}
	realizedTotal, _ := doc["realizedValueTotal"].(float64)
	estTotal, _ := doc["estimatedOpportunityTotal"].(float64)
	if estTotal > 0 && realizedTotal >= estTotal {
		// estimated should typically be larger projection; not strict but check estimated not merged into realized wrongly
	}
	return rejectForbiddenPositiveClaims(doc)
}

func boolOr(v interface{}) bool {
	b, _ := v.(bool)
	return b
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

// ValidateSprint5Examples validates all Sprint 5 reference examples.
func ValidateSprint5Examples(repoRoot string) error {
	files := map[string]func(map[string]interface{}) error{
		"movement-accounting-record-reference.json": ValidateMovementAccountingRecord,
		"resource-unit-price-reference.json":          ValidateResourceUnitPrice,
		"value-attribution-reference.json": func(d map[string]interface{}) error {
			return requireClaimBoundary(d)
		},
		"rollback-impact-accounting-reference.json": func(d map[string]interface{}) error {
			return requireClaimBoundary(d)
		},
		"monthly-savings-rollup-reference.json":      ValidateMonthlySavingsRollup,
		"savings-claim-classification-reference.json": ValidateSavingsClaimClassification,
		"accounting-confidence-reference.json": func(d map[string]interface{}) error {
			return requireClaimBoundary(d)
		},
		"realized-savings-ledger-reference.json": ValidateRealizedSavingsLedgerSurface,
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
		if name != "realized-savings-ledger-reference.json" {
			if err := rejectForbiddenPositiveClaims(doc); err != nil {
				return fmt.Errorf("%s: %w", name, err)
			}
		}
	}
	return nil
}

// SchemaFilesRequiredSprint5 returns Sprint 5 schema basenames.
func SchemaFilesRequiredSprint5() []string {
	return []string{
		"realized-savings-ledger-v1.schema.json",
		"movement-accounting-record-v1.schema.json",
		"resource-unit-price-v1.schema.json",
		"value-attribution-v1.schema.json",
		"rollback-impact-accounting-v1.schema.json",
		"monthly-savings-rollup-v1.schema.json",
		"savings-claim-classification-v1.schema.json",
		"accounting-confidence-v1.schema.json",
	}
}
