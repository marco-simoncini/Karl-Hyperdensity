package releasegate

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const MilestoneDashboardEnterpriseCleanupGA = "hyperdensity_dashboard_enterprise_cleanup_ga_release_gate_v1"

var forbiddenExecutiveRoutePatterns = []string{
	"demo", "mock", "test", "debug", "lab", "reference", "fake", "sample",
}

var forbiddenPositiveClaims = []string{
	"guaranteed savings active",
	"universal performance improvement",
	"general production auto",
	"production auto with policy",
	"production_auto_with_policy",
	"windows total ram hotplug supported",
	"logical vcpu hotplug supported",
	"1000 production workloads proven",
	"dashboard executes runtime changes",
	"dashboard executor",
	"reference payload as production proof",
	"synthetic proof as production proof",
}

var executiveHiddenTracks = map[string]bool{
	"lab": true, "debug": true, "archived": true, "reference_only": true, "synthetic_shadow": true,
}

func boolOr(v interface{}) bool {
	b, _ := v.(bool)
	return b
}

func strOr(v interface{}) string {
	s, _ := v.(string)
	return s
}

func requireFalse(doc map[string]interface{}, key string) error {
	if v, ok := doc[key].(bool); !ok || v {
		return fmt.Errorf("%s must be false", key)
	}
	return nil
}

func requireClaimBoundary(doc map[string]interface{}) error {
	if cbs, ok := doc["claimBoundaries"].([]interface{}); ok && len(cbs) > 0 {
		return nil
	}
	switch cb := doc["claimBoundary"].(type) {
	case string:
		if strings.TrimSpace(cb) == "" {
			return fmt.Errorf("claimBoundary required")
		}
	case []interface{}:
		if len(cb) == 0 {
			return fmt.Errorf("claimBoundary required")
		}
	default:
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

func scanForbiddenClaims(text string) error {
	lower := strings.ToLower(text)
	for _, phrase := range forbiddenPositiveClaims {
		if strings.Contains(lower, phrase) {
			if strings.Contains(lower, "remains disabled") || strings.Contains(lower, "not a universal") ||
				strings.Contains(lower, "not production proof") || strings.Contains(lower, "blocked") ||
				strings.Contains(lower, "forbidden") || strings.Contains(lower, "rejected") {
				continue
			}
			return fmt.Errorf("forbidden claim phrase: %s", phrase)
		}
	}
	return nil
}

// ValidateGAReleaseGate checks GA release gate invariants.
func ValidateGAReleaseGate(doc map[string]interface{}) error {
	if m, _ := doc["milestone"].(string); m != MilestoneDashboardEnterpriseCleanupGA {
		return fmt.Errorf("milestone must be %s", MilestoneDashboardEnterpriseCleanupGA)
	}
	for _, key := range []string{
		"generalProductionAutoAllowed", "productionAutoWithPolicy", "guaranteedSavingsClaimed",
		"universalPerformanceImprovementClaimed", "referencePayloadCountedAsProduction",
		"syntheticProofCountedAsProduction", "dashboardRuntimeControlsExposed", "dashboardExecutor",
		"fluidvirtPolicyAuthority", "inventoryRuntimeExecutor",
	} {
		if err := requireFalse(doc, key); err != nil {
			return err
		}
	}
	decision, _ := doc["releaseDecision"].(string)
	if decision == "ga_pass" {
		if !boolOr(doc["gaReady"]) {
			return fmt.Errorf("releaseDecision=ga_pass requires gaReady=true")
		}
		if checks, ok := doc["checks"].([]interface{}); ok {
			for _, item := range checks {
				ch, ok := item.(map[string]interface{})
				if !ok {
					continue
				}
				if boolOr(ch["requiredForGA"]) {
					status, _ := ch["status"].(string)
					if status != "passed" && status != "not_applicable" {
						return fmt.Errorf("releaseDecision=ga_pass with failed required check %s", ch["checkId"])
					}
				}
			}
		}
	}
	if decision == "ga_pass" && !boolOr(doc["generalProductionAutoAllowed"]) {
		// Sprint 10: general auto disabled means ga_pass is only valid if scope excludes general auto explicitly
		if boolOr(doc["gaReady"]) {
			return fmt.Errorf("releaseDecision=ga_pass incompatible with generalProductionAutoAllowed=false in Sprint 10")
		}
	}
	if decision != "ga_pass" && decision != "ga_blocked" && decision != "canary_only" && decision != "preview_only" && decision != "lab_only" {
		return fmt.Errorf("invalid releaseDecision: %s", decision)
	}
	if cbs, ok := doc["claimBoundaries"].([]interface{}); ok && len(cbs) > 0 {
		cbText := fmt.Sprintf("%v", cbs)
		return scanForbiddenClaims(cbText)
	}
	if err := requireClaimBoundary(doc); err != nil {
		return err
	}
	return scanForbiddenClaims(fmt.Sprintf("%v", doc["claimBoundary"]))
}

// ValidateSurfaceClassification checks per-surface rules.
func ValidateSurfaceClassification(doc map[string]interface{}) error {
	if err := requireClaimBoundary(doc); err != nil {
		return err
	}
	track, _ := doc["releaseTrack"].(string)
	route := strings.ToLower(strOr(doc["route"]))
	if boolOr(doc["referenceOnly"]) && boolOr(doc["canBeUsedForProductionProof"]) {
		return fmt.Errorf("referenceOnly surface cannot be production proof")
	}
	if boolOr(doc["syntheticOrShadow"]) && boolOr(doc["canBeUsedForProductionProof"]) {
		return fmt.Errorf("syntheticOrShadow surface cannot be production proof")
	}
	if executiveHiddenTracks[track] && boolOr(doc["visibleInExecutive"]) {
		return fmt.Errorf("releaseTrack=%s cannot be visibleInExecutive", track)
	}
	for _, pat := range forbiddenExecutiveRoutePatterns {
		if strings.Contains(route, pat) && boolOr(doc["visibleInExecutive"]) {
			return fmt.Errorf("route containing %s cannot be visibleInExecutive", pat)
		}
	}
	if track == "production_canary" {
		cb := strings.ToLower(fmt.Sprintf("%v", doc["claimBoundary"]))
		if strings.Contains(cb, "general production auto") && !strings.Contains(cb, "remains disabled") {
			return fmt.Errorf("production canary cannot be labeled general production auto")
		}
	}
	return nil
}

// ValidateReleaseGateDecision checks decision document.
func ValidateReleaseGateDecision(doc map[string]interface{}) error {
	if comp, _ := doc["decidedByComponent"].(string); comp != "Karl-Hyperdensity" {
		return fmt.Errorf("decidedByComponent must be Karl-Hyperdensity")
	}
	if next, _ := doc["nextAllowedState"].(string); next == "production_auto_with_policy" {
		return fmt.Errorf("nextAllowedState=production_auto_with_policy forbidden in Sprint 10")
	}
	return requireClaimBoundary(doc)
}

// ValidateClaimGate checks claim gate invariants.
func ValidateClaimGate(doc map[string]interface{}) error {
	if err := requireClaimBoundary(doc); err != nil {
		return err
	}
	if err := requireEvidenceRefs(doc); err != nil {
		return err
	}
	if boolOr(doc["forbiddenClaimsPresent"]) {
		return fmt.Errorf("forbiddenClaimsPresent must be false")
	}
	if allowed, ok := doc["allowedClaims"].([]interface{}); ok {
		for _, item := range allowed {
			if s, ok := item.(string); ok {
				if err := scanForbiddenClaims(s); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// ValidateEvidenceGate checks evidence separation.
func ValidateEvidenceGate(doc map[string]interface{}) error {
	if err := requireClaimBoundary(doc); err != nil {
		return err
	}
	if err := requireEvidenceRefs(doc); err != nil {
		return err
	}
	prod := map[string]bool{}
	for _, item := range doc["productionEvidenceRefs"].([]interface{}) {
		prod[strOr(item)] = true
	}
	for _, item := range doc["referenceOnlyRefs"].([]interface{}) {
		if prod[strOr(item)] {
			return fmt.Errorf("referenceOnlyRefs cannot be in productionEvidenceRefs")
		}
	}
	for _, item := range doc["syntheticShadowRefs"].([]interface{}) {
		if prod[strOr(item)] {
			return fmt.Errorf("syntheticShadowRefs cannot be in productionEvidenceRefs")
		}
	}
	return nil
}

// ValidateReferencePayloadGate checks reference payload gate.
func ValidateReferencePayloadGate(doc map[string]interface{}) error {
	if err := requireClaimBoundary(doc); err != nil {
		return err
	}
	return requireFalse(doc, "countedAsProduction")
}

// ValidateSyntheticProofGate checks synthetic proof gate.
func ValidateSyntheticProofGate(doc map[string]interface{}) error {
	if err := requireClaimBoundary(doc); err != nil {
		return err
	}
	return requireFalse(doc, "countedAsProduction")
}

// ValidateCountConsistencyGate checks count consistency gate includes mismatch representation.
func ValidateCountConsistencyGate(doc map[string]interface{}) error {
	if err := requireClaimBoundary(doc); err != nil {
		return err
	}
	if err := requireEvidenceRefs(doc); err != nil {
		return err
	}
	mismatchCount := 0
	if mismatches, ok := doc["mismatches"].([]interface{}); ok {
		mismatchCount = len(mismatches)
	}
	if intFrom(doc["mismatchCount"]) != mismatchCount {
		return fmt.Errorf("mismatchCount must match mismatches length")
	}
	return nil
}

func intFrom(v interface{}) int {
	switch n := v.(type) {
	case float64:
		return int(n)
	case int:
		return n
	default:
		return 0
	}
}

// ValidateEnterpriseSurfaceManifest checks manifest invariants.
func ValidateEnterpriseSurfaceManifest(doc map[string]interface{}) error {
	if !boolOr(doc["enterpriseMode"]) {
		return fmt.Errorf("enterpriseMode must be true")
	}
	if mode, _ := doc["defaultNavigationMode"].(string); mode != "enterprise" {
		return fmt.Errorf("defaultNavigationMode must be enterprise")
	}
	if surfaces, ok := doc["surfaces"].([]interface{}); ok {
		for _, item := range surfaces {
			surf, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			if err := ValidateSurfaceClassification(surf); err != nil {
				return fmt.Errorf("surface %v: %w", surf["surfaceId"], err)
			}
		}
	}
	return requireClaimBoundary(doc)
}

// ValidateReferenceFile validates the main GA release gate reference payload.
func ValidateReferenceFile(path string) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var doc map[string]interface{}
	if err := json.Unmarshal(raw, &doc); err != nil {
		return err
	}
	if err := ValidateGAReleaseGate(doc); err != nil {
		return err
	}
	if cg, ok := doc["claimGate"].(map[string]interface{}); ok {
		if err := ValidateClaimGate(cg); err != nil {
			return err
		}
	}
	if eg, ok := doc["evidenceGate"].(map[string]interface{}); ok {
		if err := ValidateEvidenceGate(eg); err != nil {
			return err
		}
	}
	if rpg, ok := doc["referencePayloadGate"].(map[string]interface{}); ok {
		if err := ValidateReferencePayloadGate(rpg); err != nil {
			return err
		}
	}
	if spg, ok := doc["syntheticProofGate"].(map[string]interface{}); ok {
		if err := ValidateSyntheticProofGate(spg); err != nil {
			return err
		}
	}
	if ccg, ok := doc["countConsistencyGate"].(map[string]interface{}); ok {
		if err := ValidateCountConsistencyGate(ccg); err != nil {
			return err
		}
	}
	if sm, ok := doc["surfaceManifest"].(map[string]interface{}); ok {
		if err := ValidateEnterpriseSurfaceManifest(sm); err != nil {
			return err
		}
	}
	decision, _ := doc["releaseDecision"].(string)
	if decision != "canary_only" && decision != "ga_blocked" {
		if !boolOr(doc["gaReady"]) {
			return fmt.Errorf("releaseDecision must be canary_only or ga_blocked unless all GA gates pass")
		}
	}
	return nil
}

// ValidateExamplesDir validates sprint 10 examples under dir.
func ValidateExamplesDir(dir string) error {
	ref := filepath.Join(dir, "ga-release-gate-reference.json")
	return ValidateReferenceFile(ref)
}
