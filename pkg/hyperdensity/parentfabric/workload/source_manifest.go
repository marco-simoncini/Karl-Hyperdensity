package workload

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
)

// SourceManifest is the static anti-drift contract for the Sprint 52 three-function copy.
type SourceManifest struct {
	SourceFile                    string   `json:"sourceFile"`
	CopyMode                      string   `json:"copyMode"`
	WorkloadHelpersOverallVerdict string   `json:"workloadHelpersOverallVerdict"`
	CopiedFunctions               []string `json:"copiedFunctions"`
	DeferredCategories            []string `json:"deferredCategories"`
	RuntimeImportAllowed          bool     `json:"runtimeImportAllowed"`
}

var (
	expectedCopiedFunctions = []string{
		"hyperdensityAppsWorkloadResource",
		"hyperdensityPilotWorkloadTerm",
		"hyperdensityExecutionSupportsLiveApplyKind",
	}
	expectedDeferredCategories = []string{
		"api_path_builders",
		"observed_state_builders",
		"execution_apply_helpers",
	}
)

// DefaultSourceManifest returns the canonical manifest for the copied pure-candidate slice.
func DefaultSourceManifest() SourceManifest {
	copied := append([]string(nil), expectedCopiedFunctions...)
	sort.Strings(copied)
	deferred := append([]string(nil), expectedDeferredCategories...)
	sort.Strings(deferred)
	return SourceManifest{
		SourceFile:                    WorkloadHelpersSourceFile,
		CopyMode:                      "three-pure-candidates-only",
		WorkloadHelpersOverallVerdict: "copy-deferred",
		CopiedFunctions:               copied,
		DeferredCategories:            deferred,
		RuntimeImportAllowed:          false,
	}
}

// ValidateSourceManifest checks manifest invariants.
func ValidateSourceManifest() error {
	m := DefaultSourceManifest()
	if m.SourceFile != WorkloadHelpersSourceFile {
		return fmt.Errorf("sourceFile: got %q want %q", m.SourceFile, WorkloadHelpersSourceFile)
	}
	if m.CopyMode != "three-pure-candidates-only" {
		return fmt.Errorf("copyMode: got %q", m.CopyMode)
	}
	if m.WorkloadHelpersOverallVerdict != "copy-deferred" {
		return fmt.Errorf("workloadHelpersOverallVerdict: got %q", m.WorkloadHelpersOverallVerdict)
	}
	if m.RuntimeImportAllowed {
		return errors.New("runtimeImportAllowed must be false")
	}
	copied := append([]string(nil), expectedCopiedFunctions...)
	sort.Strings(copied)
	if !reflect.DeepEqual(m.CopiedFunctions, copied) {
		return fmt.Errorf("copiedFunctions: got %v want %v", m.CopiedFunctions, copied)
	}
	deferred := append([]string(nil), expectedDeferredCategories...)
	sort.Strings(deferred)
	if !reflect.DeepEqual(m.DeferredCategories, deferred) {
		return fmt.Errorf("deferredCategories: got %v want %v", m.DeferredCategories, deferred)
	}
	if len(m.CopiedFunctions) != 3 {
		return fmt.Errorf("copiedFunctions count: got %d want 3", len(m.CopiedFunctions))
	}
	return nil
}
