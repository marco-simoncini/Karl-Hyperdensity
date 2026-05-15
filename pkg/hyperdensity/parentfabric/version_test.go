package parentfabric

import "testing"

func TestParentFabricVersionConstants(t *testing.T) {
	if ParentFabricPackageVersion == "" {
		t.Fatal("ParentFabricPackageVersion must be non-empty")
	}
	if ParentFabricRuntimeOwnership != "dashboard-runtime-owner" {
		t.Fatalf("unexpected ParentFabricRuntimeOwnership: %q", ParentFabricRuntimeOwnership)
	}
	if ParentFabricExtractionMode != "pure-core-skeleton" {
		t.Fatalf("unexpected ParentFabricExtractionMode: %q", ParentFabricExtractionMode)
	}
}
