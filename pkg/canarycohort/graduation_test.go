package canarycohort

import "testing"

func TestGraduationEvidenceReadiness(t *testing.T) {
	doc, _ := ReferenceSurface()
	if err := validateGraduationEvidence(doc); err != nil {
		t.Fatal(err)
	}
	grad := doc["guardedPolicyGraduationEvidence"].(map[string]interface{})
	if boolv(grad["productionAutoWithPolicyEnabled"]) {
		t.Fatal("productionAutoWithPolicyEnabled must be false in graduation evidence")
	}
}

func TestGraduationRejectsAutoPolicyEnabled(t *testing.T) {
	doc := validSurface()
	doc["productionAutoWithPolicyEnabled"] = true
	if err := validateGraduationEvidence(doc); err == nil {
		t.Fatal("expected rejection when productionAutoWithPolicyEnabled=true")
	}
}
