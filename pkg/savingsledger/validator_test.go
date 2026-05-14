package savingsledger_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/savingsledger"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

func TestSprint5SchemaFilesExist(t *testing.T) {
	root := repoRoot(t)
	for _, name := range savingsledger.SchemaFilesRequiredSprint5() {
		if _, err := os.Stat(filepath.Join(root, "schemas", name)); err != nil {
			t.Fatalf("missing schema %s: %v", name, err)
		}
	}
}

func TestSprint5ExamplesValidate(t *testing.T) {
	if err := savingsledger.ValidateSprint5Examples(repoRoot(t)); err != nil {
		t.Fatalf("ValidateSprint5Examples: %v", err)
	}
}

func TestGuaranteedSavingsClaimedRejected(t *testing.T) {
	doc := map[string]interface{}{
		"milestone": savingsledger.MilestoneRealizedSavingsLedger,
		"guaranteedSavingsAllowed": false, "guaranteedSavingsClaimed": true,
		"estimatedValueCountedAsGuaranteed": false, "syntheticValueCountedAsProduction": false,
		"safetyInvariants": map[string]interface{}{
			"guaranteedSavingsAllowed": false, "guaranteedSavingsClaimed": true,
			"estimatedValueCountedAsGuaranteed": false, "syntheticValueCountedAsProduction": false,
			"dashboardAccountingSourceOfTruth": false, "fluidvirtAccountingAuthority": false,
		},
		"sourceOfTruth": map[string]interface{}{"ledgerContracts": "Karl-Hyperdensity", "movementMeasurements": "FluidVirt"},
		"records": []interface{}{},
	}
	if err := savingsledger.ValidateRealizedSavingsLedgerSurface(doc); err == nil {
		t.Fatal("expected guaranteedSavingsClaimed rejection")
	}
}

func TestSyntheticCountedAsProductionRejected(t *testing.T) {
	doc := map[string]interface{}{
		"milestone": savingsledger.MilestoneRealizedSavingsLedger,
		"guaranteedSavingsAllowed": false, "guaranteedSavingsClaimed": false,
		"estimatedValueCountedAsGuaranteed": false, "syntheticValueCountedAsProduction": true,
		"safetyInvariants": map[string]interface{}{
			"guaranteedSavingsAllowed": false, "guaranteedSavingsClaimed": false,
			"estimatedValueCountedAsGuaranteed": false, "syntheticValueCountedAsProduction": true,
			"dashboardAccountingSourceOfTruth": false, "fluidvirtAccountingAuthority": false,
		},
		"sourceOfTruth": map[string]interface{}{"ledgerContracts": "Karl-Hyperdensity", "movementMeasurements": "FluidVirt"},
		"records": []interface{}{},
	}
	if err := savingsledger.ValidateRealizedSavingsLedgerSurface(doc); err == nil {
		t.Fatal("expected syntheticValueCountedAsProduction rejection")
	}
}

func TestGuaranteeEligibleWithoutSLORejected(t *testing.T) {
	doc := map[string]interface{}{
		"claimClassification": "eligible_for_future_guarantee",
		"guaranteeEligibleForFuture": true,
		"sloPreserved": false, "donorHealthPreserved": true, "rollbackReady": true,
		"durationSeconds": 300, "unitPrice": 0.12, "netRealizedValue": 0.01,
		"synthetic": false,
		"claimBoundary": []interface{}{"x"}, "evidenceRefs": []interface{}{"y"},
	}
	if err := savingsledger.ValidateMovementAccountingRecord(doc); err == nil {
		t.Fatal("expected sloPreserved rejection")
	}
}

func TestNegativeNetValueFutureGuaranteeRejected(t *testing.T) {
	doc := map[string]interface{}{
		"claimClassification": "eligible_for_future_guarantee",
		"guaranteeEligibleForFuture": true,
		"sloPreserved": true, "donorHealthPreserved": true, "rollbackReady": true,
		"durationSeconds": 300, "unitPrice": 0.12, "netRealizedValue": -0.01,
		"synthetic": false,
		"claimBoundary": []interface{}{"x"}, "evidenceRefs": []interface{}{"y"},
	}
	if err := savingsledger.ValidateMovementAccountingRecord(doc); err == nil {
		t.Fatal("expected negative net value rejection")
	}
}

func TestGuaranteedSavingsActiveClassificationRejected(t *testing.T) {
	doc := map[string]interface{}{
		"claimClassification": "guaranteed_savings_active",
		"claimBoundary": []interface{}{"x"}, "evidenceRefs": []interface{}{"y"},
	}
	if err := savingsledger.ValidateMovementAccountingRecord(doc); err == nil {
		t.Fatal("expected guaranteed_savings_active rejection")
	}
}

func TestEstimatedOnlyMustBeExcluded(t *testing.T) {
	doc := map[string]interface{}{
		"claimClassification": "eligible_for_future_guarantee",
		"estimatedOnly": true,
		"claimBoundary": []interface{}{"x"},
	}
	if err := savingsledger.ValidateSavingsClaimClassification(doc); err == nil {
		t.Fatal("expected estimatedOnly rejection")
	}
}

func TestPlaceholderPriceNotFutureEligible(t *testing.T) {
	rec := map[string]interface{}{
		"claimClassification": "eligible_for_future_guarantee",
		"guaranteeEligibleForFuture": true,
		"sloPreserved": true, "donorHealthPreserved": true, "rollbackReady": true,
		"durationSeconds": 0, "unitPrice": 0.0, "netRealizedValue": 0.01,
		"synthetic": false,
		"claimBoundary": []interface{}{"x"}, "evidenceRefs": []interface{}{"y"},
	}
	if err := savingsledger.ValidateMovementAccountingRecord(rec); err == nil {
		t.Fatal("expected missing duration/price rejection")
	}
}
