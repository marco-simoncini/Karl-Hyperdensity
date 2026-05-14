package realmarkettick

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

func TestReferenceSurfaceSafetyFlags(t *testing.T) {
	doc, _ := ReferenceSurface()
	for _, key := range []string{"productionMovementExecuted", "generalProductionAutoAllowed", "productionAutoWithPolicy"} {
		if boolv(doc[key]) {
			t.Fatalf("%s must be false", key)
		}
	}
}
