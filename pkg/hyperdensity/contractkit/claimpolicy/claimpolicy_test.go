package claimpolicy

import "testing"

func TestPackageVersion(t *testing.T) {
	if PackageVersion == "" {
		t.Fatal("PackageVersion must be set")
	}
}

func TestKnownPosture(t *testing.T) {
	for _, p := range Postures() {
		if !KnownPosture(p) {
			t.Fatalf("Postures() returned unknown posture %q", p)
		}
	}
	if KnownPosture(PostureKind("unknown_posture")) {
		t.Fatal("unknown posture must not be known")
	}
}

func TestPosturesStableOrder(t *testing.T) {
	got := Postures()
	if len(got) != 3 {
		t.Fatalf("len(Postures)=%d", len(got))
	}
	if got[0] != PostureEvidenceNamespace || got[1] != PostureOperatorControlled || got[2] != PostureVisibilityOnly {
		t.Fatalf("unexpected order: %#v", got)
	}
}
