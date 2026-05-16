package executiontypes

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
)

func TestDefaultSourceManifest_matchesGolden(t *testing.T) {
	got, err := json.Marshal(DefaultSourceManifest())
	if err != nil {
		t.Fatal(err)
	}
	want, err := os.ReadFile(filepath.Join("testdata", "execution_types_source_manifest.golden.json"))
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

func TestExecutionSummaryJSONTags_reflectMatch(t *testing.T) {
	got := ExecutionSummaryJSONTags()
	want := DefaultSourceManifest().ExecutionSummaryJSONTags
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ExecutionSummaryJSONTags=%v manifest=%v", got, want)
	}
	if len(got) != 21 {
		t.Fatalf("tag count=%d want 21", len(got))
	}
}

func TestEngineSpineJSONTags_reflectMatch(t *testing.T) {
	got := EngineSpineJSONTags()
	want := DefaultSourceManifest().EngineSpineJSONTags
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("EngineSpineJSONTags=%v manifest=%v", got, want)
	}
}

func TestJSONTags_noEmptyOrDuplicate(t *testing.T) {
	for _, label := range []struct {
		name string
		tags []string
	}{
		{"executionSummary", ExecutionSummaryJSONTags()},
		{"engineSpine", EngineSpineJSONTags()},
	} {
		if err := validateNoEmptyOrDuplicateTags(label.name, label.tags); err != nil {
			t.Fatal(err)
		}
	}
}

func TestContractDocument_alignsWithSourceManifest(t *testing.T) {
	doc := DefaultContractDocument()
	m := DefaultSourceManifest()
	if doc.SourceFile != m.SourceFile {
		t.Fatalf("sourceFile doc=%q manifest=%q", doc.SourceFile, m.SourceFile)
	}
	if doc.SourceStats.TypeDefinitionCount != m.SourceTypeDefinitionCount {
		t.Fatalf("typeDefinitionCount doc=%d manifest=%d", doc.SourceStats.TypeDefinitionCount, m.SourceTypeDefinitionCount)
	}
	if !reflect.DeepEqual(doc.SourceStats.SourceImports, m.SourceImportSet) {
		t.Fatalf("sourceImports doc=%v manifest=%v", doc.SourceStats.SourceImports, m.SourceImportSet)
	}
	copied := append([]string(nil), doc.CopiedTypes...)
	sort.Strings(copied)
	if !reflect.DeepEqual(copied, m.CopiedTypes) {
		t.Fatalf("copiedTypes doc=%v manifest=%v", copied, m.CopiedTypes)
	}
}
