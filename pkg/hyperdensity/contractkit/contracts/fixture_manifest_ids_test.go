package contracts

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExpectedM1M8CaseIDs_stable(t *testing.T) {
	got := ExpectedM1M8CaseIDs()
	if len(got) != 8 {
		t.Fatalf("expected 8 ids, got %d", len(got))
	}
	seen := make(map[string]struct{}, len(got))
	for _, id := range got {
		if id == "" {
			t.Fatal("empty id")
		}
		if _, ok := seen[id]; ok {
			t.Fatalf("duplicate id %q", id)
		}
		seen[id] = struct{}{}
	}
}

func TestCaseIDs_and_ManifestCaseIDSet_exampleManifest(t *testing.T) {
	root := moduleRoot(t)
	path := filepath.Join(root, "testdata", "dashboard", "hyperdensity_parity_manifest_m1_m7.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	m, err := ParseFixtureManifest(data)
	if err != nil {
		t.Fatal(err)
	}
	ids := CaseIDs(m)
	set := ManifestCaseIDSet(m)
	if len(ids) != len(set) {
		t.Fatalf("duplicate ids in manifest: %d cases, %d unique", len(ids), len(set))
	}
	expected := ExpectedM1M8CaseIDs()
	if len(ids) != len(expected) {
		t.Fatalf("case count: got %d want %d", len(ids), len(expected))
	}
	for _, id := range expected {
		if _, ok := set[id]; !ok {
			t.Fatalf("example manifest missing expected id %q", id)
		}
	}
	for id := range set {
		found := false
		for _, e := range expected {
			if e == id {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("example manifest unexpected id %q", id)
		}
	}
}
