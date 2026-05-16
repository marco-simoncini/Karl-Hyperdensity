package executiontypes

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestExecutionTypesConstants(t *testing.T) {
	if ExecutionTypesPackageVersion == "" {
		t.Fatal("ExecutionTypesPackageVersion must be non-empty")
	}
	if ExecutionTypesSourceFile == "" {
		t.Fatal("ExecutionTypesSourceFile must be non-empty")
	}
	if ExecutionTypesExtractionMode != "copy-contract-no-runtime-wiring" {
		t.Fatalf("unexpected ExecutionTypesExtractionMode: %q", ExecutionTypesExtractionMode)
	}
}

func TestExecutionEngineSpineJSONKeys(t *testing.T) {
	want := []string{"summary", "supportsApply", "supportedSurfaces", "applyNotes"}
	if !reflect.DeepEqual(ExecutionEngineSpineJSONKeys, want) {
		t.Fatalf("ExecutionEngineSpineJSONKeys=%v want %v", ExecutionEngineSpineJSONKeys, want)
	}
}

func TestHyperdensityExecutionSummary_fieldCount(t *testing.T) {
	var s HyperdensityExecutionSummary
	if got := reflect.TypeOf(s).NumField(); got != 21 {
		t.Fatalf("HyperdensityExecutionSummary field count=%d want 21", got)
	}
}

func TestCanonicalContractJSON_matchesGolden(t *testing.T) {
	got, err := CanonicalContractJSON()
	if err != nil {
		t.Fatal(err)
	}
	goldenPath := filepath.Join("testdata", "execution_types_contract.golden.json")
	want, err := os.ReadFile(goldenPath)
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
	if len(gotDoc) != len(wantDoc) {
		t.Fatalf("golden key count mismatch: got %d want %d", len(gotDoc), len(wantDoc))
	}
	for k, gv := range gotDoc {
		wv, ok := wantDoc[k]
		if !ok {
			t.Fatalf("unexpected key in generated contract: %q", k)
		}
		if string(gv) != string(wv) {
			t.Fatalf("contract field %q mismatch:\ngot  %s\nwant %s", k, gv, wv)
		}
	}
}

func TestContractDocument_noRuntimeWiringMarkers(t *testing.T) {
	b, err := CanonicalContractJSON()
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	forbidden := []string{
		string([]byte{'k', '8', 's', '.', 'i', 'o', '/'}),
		string([]byte{'k', 'u', 'b', 'e', 'v', 'i', 'r', 't', '.', 'i', 'o', '/'}),
		string([]byte{'n', 'e', 't', '/', 'h', 't', 't', 'p'}),
	}
	for _, f := range forbidden {
		if strings.Contains(s, f) {
			t.Fatalf("contract JSON must not reference forbidden import pattern")
		}
	}
}

func TestHyperdensityCPUQuantity_roundTrip(t *testing.T) {
	q := HyperdensityCPUQuantity{Quantity: "100m", Millicores: 100}
	b, err := json.Marshal(q)
	if err != nil {
		t.Fatal(err)
	}
	var out HyperdensityCPUQuantity
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatal(err)
	}
	if out != q {
		t.Fatalf("round trip: %+v != %+v", out, q)
	}
}
