package primitives

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestNormalizeCPUQuantity(t *testing.T) {
	cases := []struct {
		in     string
		mc     int64
		ok     bool
	}{
		{"100m", 100, true},
		{"1", 1000, true},
		{"2", 2000, true},
		{"bad", 0, false},
	}
	for _, c := range cases {
		_, mc, ok := NormalizeCPUQuantity(c.in)
		if ok != c.ok || (ok && mc != c.mc) {
			t.Fatalf("NormalizeCPUQuantity(%q) ok=%v mc=%d want ok=%v mc=%d", c.in, ok, mc, c.ok, c.mc)
		}
	}
}

func TestNormalizeMemoryQuantity(t *testing.T) {
	cases := []struct {
		in    string
		bytes int64
		ok    bool
	}{
		{"128Mi", 128 * 1024 * 1024, true},
		{"1Gi", 1024 * 1024 * 1024, true},
		{"512Ki", 512 * 1024, true},
		{"1000", 1000, true},
		{"nope", 0, false},
	}
	for _, c := range cases {
		_, b, ok := NormalizeMemoryQuantity(c.in)
		if ok != c.ok || (ok && b != c.bytes) {
			t.Fatalf("NormalizeMemoryQuantity(%q) ok=%v bytes=%d want ok=%v bytes=%d", c.in, ok, b, c.ok, c.bytes)
		}
	}
}

func TestNestedHelpers_noPanicOnMissing(t *testing.T) {
	obj := map[string]interface{}{"a": map[string]interface{}{"b": "x"}}
	if _, ok := StringAt(obj, "a", "missing"); ok {
		t.Fatal("expected missing path false")
	}
	if _, ok := MapAt(nil, "a"); ok {
		t.Fatal("nil map")
	}
}

func TestCanonicalContractJSON_matchesGolden(t *testing.T) {
	got, err := CanonicalContractJSON()
	if err != nil {
		t.Fatal(err)
	}
	want, err := os.ReadFile(filepath.Join("testdata", "primitives_contract.golden.json"))
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
			t.Fatalf("unexpected key %q", k)
		}
		if string(gv) != string(wv) {
			t.Fatalf("field %q mismatch:\ngot  %s\nwant %s", k, gv, wv)
		}
	}
}

func TestPrimitivesPackageVersion(t *testing.T) {
	if PrimitivesPackageVersion == "" {
		t.Fatal("PrimitivesPackageVersion empty")
	}
}
