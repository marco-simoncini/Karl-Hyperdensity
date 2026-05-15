package canarycohort

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

func TestCohortMovementCount(t *testing.T) {
	doc, _ := ReferenceSurface()
	if num(doc["cohortMovementCount"]) < 2 {
		t.Fatal("cohortMovementCount must be >= 2")
	}
}

func TestProductionAutoWithPolicyDisabled(t *testing.T) {
	doc, _ := ReferenceSurface()
	for _, key := range []string{"productionAutoWithPolicy", "productionAutoWithPolicyEnabled", "productionAutoWithPolicyActivated"} {
		if boolv(doc[key]) {
			t.Fatalf("%s must be false", key)
		}
	}
}
