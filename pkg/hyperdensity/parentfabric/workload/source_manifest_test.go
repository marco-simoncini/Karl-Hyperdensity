package workload

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultSourceManifest_matchesGolden(t *testing.T) {
	got, err := json.Marshal(DefaultSourceManifest())
	if err != nil {
		t.Fatal(err)
	}
	want, err := os.ReadFile(filepath.Join("testdata", "workload_pure_candidates_source_manifest.golden.json"))
	if err != nil {
		t.Fatal(err)
	}
	var gotDoc, wantDoc map[string]json.RawMessage
	if err := json.Unmarshal(got, &gotDoc); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(want, &wantDoc); err != nil {
		t.Fatal(err)
	}
	for k, gv := range gotDoc {
		wv, ok := wantDoc[k]
		if !ok {
			t.Fatalf("unexpected manifest key %q", k)
		}
		if string(gv) != string(wv) {
			t.Fatalf("manifest field %q mismatch:\ngot  %s\nwant %s", k, gv, wv)
		}
	}
}

func TestValidateSourceManifest(t *testing.T) {
	if err := ValidateSourceManifest(); err != nil {
		t.Fatal(err)
	}
}
