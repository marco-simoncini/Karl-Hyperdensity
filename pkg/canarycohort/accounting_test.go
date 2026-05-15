package canarycohort

import "testing"

func TestAggregateAccountingRequiresEvidence(t *testing.T) {
	doc := validSurface()
	doc["realizedCompressionAggregation"].(map[string]interface{})["movementEvidencePresent"] = false
	if err := validateAccountingAggregation(doc); err == nil {
		t.Fatal("expected rejection without movement evidence")
	}
}
