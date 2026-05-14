package performanceproof

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const MilestoneUniversalSloGuard = "hyperdensity_universal_slo_guard_certified_uplift_v1"

var allowedCertifiedBottlenecks = map[string]bool{
	"cpu_bound": true, "memory_bound": true, "cpu_memory_bound": true,
}

var forbiddenPositiveClaims = []string{
	"universal performance improvement",
	"guaranteed savings active",
	"guaranteed savings claimed",
	"production autonomous apply",
	"autonomous production mutation",
	"windows total ram hotplug supported",
	"logical vcpu hotplug supported",
	"1000 production workloads proven",
	"dashboard is performance source of truth",
	"dashboard performance source of truth",
	"fluidvirt performance claim authority",
	"inventory applies runtime",
}

var skipKeys = map[string]bool{
	"forbiddenPhrases": true, "blockerCodes": true, "blockers": true,
	"exclusionReason": true, "missingEvidence": true,
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

// ValidateSloGuardEvaluation checks SLO guard evaluation invariants.
func ValidateSloGuardEvaluation(doc map[string]interface{}) error {
	status, _ := doc["sloGuardStatus"].(string)
	if status == "passed" {
		if reg, _ := doc["regressionDetected"].(bool); reg {
			return fmt.Errorf("SLO guard cannot pass when regressionDetected=true")
		}
		if rb, _ := doc["rollbackRequired"].(bool); rb {
			return fmt.Errorf("SLO guard cannot pass when rollbackRequired=true")
		}
	}
	if err := requireClaimBoundary(doc); err != nil {
		return err
	}
	return requireEvidenceRefs(doc)
}

// ValidateCertifiedPerformanceUplift checks certified uplift invariants.
func ValidateCertifiedPerformanceUplift(doc map[string]interface{}, bottleneckType string, sloPassed, donorPreserved, noRegression bool) error {
	certified, _ := doc["certified"].(bool)
	if !certified {
		return requireClaimBoundary(doc)
	}
	if !allowedCertifiedBottlenecks[bottleneckType] {
		return fmt.Errorf("certified uplift forbidden for bottleneckType %s", bottleneckType)
	}
	conf, _ := doc["confidence"].(float64)
	if conf < 0.80 {
		return fmt.Errorf("certified uplift requires confidence >= 0.80")
	}
	if !sloPassed {
		return fmt.Errorf("certified uplift requires SLO guard passed")
	}
	if !donorPreserved {
		return fmt.Errorf("certified uplift requires donor health preserved")
	}
	if !noRegression {
		return fmt.Errorf("certified uplift requires noRegressionCertified=true")
	}
	if doc["baselineSampleRef"] == nil && doc["baselineSampleId"] == nil {
		return fmt.Errorf("certified uplift requires baseline sample ref")
	}
	if doc["postMutationSampleRef"] == nil && doc["postMutationSampleId"] == nil {
		return fmt.Errorf("certified uplift requires post-mutation sample ref")
	}
	excluded, ok := doc["excludedFromUniversalClaim"].(bool)
	if !ok || !excluded {
		return fmt.Errorf("certified uplift must have excludedFromUniversalClaim=true")
	}
	if err := requireClaimBoundary(doc); err != nil {
		return err
	}
	if err := requireEvidenceRefs(doc); err != nil {
		return err
	}
	cb := flattenClaimBoundary(doc)
	merged := strings.ToLower(strings.Join(cb, " "))
	if strings.Contains(merged, "universal performance improvement") && !strings.Contains(merged, "not a universal") && !strings.Contains(merged, "not universal") {
		return fmt.Errorf("claimBoundary must not positively claim universal performance improvement")
	}
	return nil
}

// ValidateNoRegressionResult checks no-regression certification.
func ValidateNoRegressionResult(doc map[string]interface{}) error {
	if cert, _ := doc["noRegressionCertified"].(bool); cert {
		if rb, _ := doc["rollbackRequired"].(bool); rb {
			return fmt.Errorf("noRegressionCertified cannot be true when rollbackRequired=true")
		}
		if blocked, _ := doc["regressionBlocked"].(bool); blocked {
			return fmt.Errorf("noRegressionCertified cannot be true when regressionBlocked=true")
		}
	}
	return requireClaimBoundary(doc)
}

// ValidatePerformanceProofClassification checks proof classification.
func ValidatePerformanceProofClassification(doc map[string]interface{}) error {
	status, _ := doc["proofStatus"].(string)
	certified, _ := doc["certifiedUplift"].(bool)
	switch status {
	case "uplift_certified":
		if !certified {
			return fmt.Errorf("uplift_certified requires certifiedUplift=true")
		}
	case "neutral_no_claim", "insufficient_evidence", "not_cpu_ram_bound", "regression_blocked":
		if certified {
			return fmt.Errorf("%s cannot have certifiedUplift=true", status)
		}
	case "no_regression_certified":
		if certified {
			return fmt.Errorf("no_regression_certified must not have certifiedUplift=true")
		}
	}
	return requireClaimBoundary(doc)
}

// ValidateUniversalSloGuardSurface validates ConfigMap-ready SLO guard surface.
func ValidateUniversalSloGuardSurface(doc map[string]interface{}) error {
	if doc["milestone"] != MilestoneUniversalSloGuard {
		return fmt.Errorf("milestone must be %s", MilestoneUniversalSloGuard)
	}
	if err := requireFalse(doc, "universalPerformanceImprovementClaimed"); err != nil {
		return err
	}
	if err := requireFalse(doc, "guaranteedSavingsClaimed"); err != nil {
		return err
	}
	if err := requireFalse(doc, "autoApplyAllowed"); err != nil {
		return err
	}
	if err := requireFalse(doc, "productionAutonomousApplyAllowed"); err != nil {
		return err
	}
	inv, ok := doc["safetyInvariants"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("safetyInvariants required")
	}
	for _, key := range []string{"universalPerformanceImprovementClaimed", "guaranteedSavingsClaimed", "autoApplyAllowed", "productionAutonomousApplyAllowed", "dashboardPerformanceSourceOfTruth", "fluidvirtPerformanceClaimAuthority"} {
		if v, ok := inv[key].(bool); !ok || v {
			return fmt.Errorf("safetyInvariants.%s must be false", key)
		}
	}
	sot, ok := doc["sourceOfTruth"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("sourceOfTruth required")
	}
	if !strings.Contains(strings.ToLower(fmt.Sprint(sot["performanceMeasurements"])), "fluidvirt") {
		return fmt.Errorf("performance measurements must be FluidVirt")
	}
	if strings.Contains(strings.ToLower(fmt.Sprint(sot["sloPerformanceProjection"])), "source of truth") {
		return fmt.Errorf("Dashboard must not be performance source of truth")
	}
	// Scenario coverage
	var cpuUplift, donorNoReg, neutral, regression, windowsGated, rollbackTrigger bool
	if uplifts, ok := doc["certifiedUpliftResults"].([]interface{}); ok {
		for _, item := range uplifts {
			u, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			bt, _ := u["bottleneckType"].(string)
			if bt == "cpu_bound" && boolOr(u["certified"]) {
				cpuUplift = true
			}
			if err := ValidateCertifiedPerformanceUplift(u, bt, true, true, true); err != nil {
				return fmt.Errorf("certifiedUpliftResults: %w", err)
			}
		}
	}
	if nrs, ok := doc["noRegressionResults"].([]interface{}); ok && len(nrs) > 0 {
		donorNoReg = true
	}
	if neutrals, ok := doc["neutralNoClaimResults"].([]interface{}); ok && len(neutrals) > 0 {
		neutral = true
	}
	if regs, ok := doc["regressionBlockedResults"].([]interface{}); ok && len(regs) > 0 {
		regression = true
	}
	if wins, ok := doc["windowsEvidenceGatedResults"].([]interface{}); ok && len(wins) > 0 {
		windowsGated = true
	}
	if triggers, ok := doc["rollbackTriggers"].([]interface{}); ok && len(triggers) > 0 {
		rollbackTrigger = true
	}
	if !cpuUplift || !donorNoReg || !neutral || !regression || !windowsGated || !rollbackTrigger {
		return fmt.Errorf("surface must include CPU uplift, donor no-regression, neutral, regression blocked, Windows gated, and rollback trigger scenarios")
	}
	return rejectForbiddenPositiveClaims(doc)
}

func boolOr(v interface{}) bool {
	b, _ := v.(bool)
	return b
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

func rejectForbiddenPositiveClaims(v interface{}) error {
	positives := collectPositiveStrings(v)
	merged := strings.ToLower(strings.Join(positives, "\n"))
	for _, phrase := range forbiddenPositiveClaims {
		if strings.Contains(merged, phrase) {
			// allow negations
			if strings.Contains(merged, "not "+phrase) || strings.Contains(merged, "not a "+phrase) {
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

// ValidateSprint6Examples validates all Sprint 6 reference examples.
func ValidateSprint6Examples(repoRoot string) error {
	files := map[string]func(map[string]interface{}) error{
		"enrolled-workload-slo-profile-reference.json": func(d map[string]interface{}) error { return requireClaimBoundary(d) },
		"slo-guard-evaluation-reference.json":          ValidateSloGuardEvaluation,
		"donor-health-guard-reference.json":            func(d map[string]interface{}) error { return requireClaimBoundary(d) },
		"receiver-health-guard-reference.json":           func(d map[string]interface{}) error { return requireClaimBoundary(d) },
		"bottleneck-classification-reference.json":       func(d map[string]interface{}) error { return requireClaimBoundary(d) },
		"performance-baseline-sample-reference.json":   func(d map[string]interface{}) error { return requireClaimBoundary(d) },
		"performance-post-mutation-sample-reference.json": func(d map[string]interface{}) error { return requireClaimBoundary(d) },
		"performance-rollback-sample-reference.json":   func(d map[string]interface{}) error { return requireClaimBoundary(d) },
		"no-regression-result-reference.json":          ValidateNoRegressionResult,
		"certified-performance-uplift-reference.json": func(d map[string]interface{}) error {
			return ValidateCertifiedPerformanceUplift(d, "cpu_bound", true, true, true)
		},
		"performance-proof-classification-reference.json": ValidatePerformanceProofClassification,
		"performance-rollback-trigger-reference.json": func(d map[string]interface{}) error { return requireClaimBoundary(d) },
		"universal-slo-guard-reference.json":            ValidateUniversalSloGuardSurface,
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
		if name != "universal-slo-guard-reference.json" {
			if err := rejectForbiddenPositiveClaims(doc); err != nil {
				return fmt.Errorf("%s: %w", name, err)
			}
		}
	}
	return nil
}

// SchemaFilesRequiredSprint6 returns Sprint 6 schema basenames.
func SchemaFilesRequiredSprint6() []string {
	return []string{
		"universal-slo-guard-v1.schema.json",
		"enrolled-workload-slo-profile-v1.schema.json",
		"slo-guard-evaluation-v1.schema.json",
		"donor-health-guard-v1.schema.json",
		"receiver-health-guard-v1.schema.json",
		"bottleneck-classification-v1.schema.json",
		"performance-baseline-sample-v1.schema.json",
		"performance-post-mutation-sample-v1.schema.json",
		"performance-rollback-sample-v1.schema.json",
		"no-regression-result-v1.schema.json",
		"certified-performance-uplift-v1.schema.json",
		"performance-proof-classification-v1.schema.json",
		"performance-rollback-trigger-v1.schema.json",
	}
}
