package canarymovement

import "testing"

func TestReferenceSurfaceValidates(t *testing.T) {
	doc, err := ReferenceSurface()
	if err != nil {
		t.Fatal(err)
	}
	if err := ValidateSurface(doc); err != nil {
		t.Fatal(err)
	}
}

func TestReferenceSurfaceCanarySafety(t *testing.T) {
	doc, _ := ReferenceSurface()
	if !boolv(doc["productionCanaryMovementExecuted"]) {
		t.Fatal("success path must execute canary movement")
	}
	if boolv(doc["generalProductionAutoAllowed"]) || boolv(doc["productionAutoWithPolicy"]) {
		t.Fatal("general auto must remain disabled")
	}
	closeout := doc["canaryCloseout"].(map[string]interface{})
	if boolv(closeout["promotedToGeneralProduction"]) {
		t.Fatal("must not promote to general production")
	}
}

func TestAccountingSeparation(t *testing.T) {
	doc, _ := ReferenceSurface()
	if err := evaluateAccounting(doc); err != nil {
		t.Fatal(err)
	}
	sep := doc["projectedRealizedSeparation"].(map[string]interface{})
	if num(sep["realizedMovedIdleValue"]) >= num(sep["projectedMovedIdleValue"]) {
		t.Fatal("realized should be less than projected opportunity in reference")
	}
}
