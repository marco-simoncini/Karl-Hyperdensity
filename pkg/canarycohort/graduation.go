package canarycohort

import "fmt"

func validateGraduationEvidence(doc map[string]interface{}) error {
	if boolv(doc["productionAutoWithPolicy"]) || boolv(doc["productionAutoWithPolicyEnabled"]) || boolv(doc["productionAutoWithPolicyActivated"]) {
		return fmt.Errorf("production_auto_with_policy must remain disabled in Sprint 22")
	}
	grad := doc["guardedPolicyGraduationEvidence"].(map[string]interface{})
	if boolv(grad["productionAutoWithPolicyEnabled"]) {
		return fmt.Errorf("graduation evidence productionAutoWithPolicyEnabled must be false")
	}
	if strv(grad["graduationDecision"]) == "production_auto_with_policy_enabled" {
		return fmt.Errorf("forbidden graduationDecision")
	}
	if boolv(doc["productionAutoWithPolicyReadinessCandidate"]) && !boolv(grad["graduationEvidenceComplete"]) {
		return fmt.Errorf("readiness candidate requires graduationEvidenceComplete")
	}
	if boolv(grad["productionAutoWithPolicyReadinessCandidate"]) && !boolv(grad["graduationEvidenceComplete"]) {
		return fmt.Errorf("graduation evidence inconsistency")
	}
	if boolv(grad["graduationEvidenceComplete"]) {
		for _, key := range []string{
			"failureBudgetWithinLimit", "rollbackEvidenceComplete", "sloEvidenceComplete",
			"accountingEvidenceComplete", "projectedRealizedSeparationComplete",
			"syntheticReferenceSeparationComplete", "dashboardExecutorAbsent", "fluidvirtAuthorityAbsent",
		} {
			if !boolv(grad[key]) {
				return fmt.Errorf("graduationEvidenceComplete with failed check: %s", key)
			}
		}
	}
	return nil
}
