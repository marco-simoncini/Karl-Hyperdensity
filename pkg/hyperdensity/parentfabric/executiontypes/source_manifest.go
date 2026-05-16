package executiontypes

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// SourceManifest is the static anti-drift contract for the Sprint 46 copy (Sprint 47 guard).
// It does not read Dashboard sources at runtime.
type SourceManifest struct {
	SourceFile                 string   `json:"sourceFile"`
	SourceImportSet            []string `json:"sourceImportSet"`
	SourceTypeDefinitionCount  int      `json:"sourceTypeDefinitionCount"`
	ExecutionSummaryFieldCount int      `json:"executionSummaryFieldCount"`
	ExecutionSummaryJSONTags   []string `json:"executionSummaryJsonTags"`
	EngineSpineJSONTags        []string `json:"engineSpineJsonTags"`
	CopiedTypes                []string `json:"copiedTypes"`
	DeferredReason             string   `json:"deferredReason"`
}

var (
	expectedSourceImportSet           = []string{"time"}
	expectedExecutionSummaryJSONTags = []string{
		"dryRunOnlyObjects",
		"applyReadyObjects",
		"blockedByPolicyObjects",
		"blockedByConfidenceObjects",
		"blockedByGovernanceObjects",
		"blockedByRuntimeObjects",
		"appliedObjects",
		"appliedWithFreezeObjects",
		"rolledBackObjects",
		"stagedForRetryObjects",
		"plannedReclaimCpu",
		"plannedReclaimMemory",
		"plannedBurstCpu",
		"plannedBurstMemory",
		"remainingBlockedCpu",
		"remainingBlockedMemory",
		"waitingApprovalCpu",
		"waitingApprovalMemory",
		"dryRunOnlyCpu",
		"dryRunOnlyMemory",
		"summary",
	}
	expectedEngineSpineJSONTags = []string{
		"summary",
		"supportsApply",
		"supportedSurfaces",
		"applyNotes",
	}
)

// DefaultSourceManifest returns the canonical manifest for the copied contract slice.
func DefaultSourceManifest() SourceManifest {
	summaryTags := ExecutionSummaryJSONTags()
	engineTags := EngineSpineJSONTags()
	copied := append([]string(nil), CopiedTypeNames...)
	sort.Strings(copied)
	return SourceManifest{
		SourceFile:                 ExecutionTypesSourceFile,
		SourceImportSet:            append([]string(nil), expectedSourceImportSet...),
		SourceTypeDefinitionCount:  ExecutionTypesSourceTypeCount,
		ExecutionSummaryFieldCount: len(summaryTags),
		ExecutionSummaryJSONTags:   summaryTags,
		EngineSpineJSONTags:        engineTags,
		CopiedTypes:                copied,
		DeferredReason:             defaultDeferredReason(),
	}
}

func defaultDeferredReason() string {
	return "Remaining nested HyperdensityExecutionEngine surface types (~146 top-level fields) stay in Dashboard until per-surface extraction sprints."
}

// ExecutionSummaryJSONTags returns json tag names from HyperdensityExecutionSummary (struct field order).
func ExecutionSummaryJSONTags() []string {
	return structJSONTags(reflect.TypeOf(HyperdensityExecutionSummary{}))
}

// EngineSpineJSONTags returns json tag names from HyperdensityExecutionEngineSpine (struct field order).
func EngineSpineJSONTags() []string {
	return structJSONTags(reflect.TypeOf(HyperdensityExecutionEngineSpine{}))
}

func structJSONTags(typ reflect.Type) []string {
	tags := make([]string, 0, typ.NumField())
	for i := 0; i < typ.NumField(); i++ {
		raw := typ.Field(i).Tag.Get("json")
		if raw == "" || raw == "-" {
			continue
		}
		name := raw
		if idx := strings.Index(name, ","); idx >= 0 {
			name = name[:idx]
		}
		if name == "" {
			continue
		}
		tags = append(tags, name)
	}
	return tags
}

// ValidateSourceManifest checks manifest invariants against the copied structs.
func ValidateSourceManifest() error {
	m := DefaultSourceManifest()
	if m.SourceFile != ExecutionTypesSourceFile {
		return fmt.Errorf("sourceFile: got %q want %q", m.SourceFile, ExecutionTypesSourceFile)
	}
	if !reflect.DeepEqual(m.SourceImportSet, expectedSourceImportSet) {
		return fmt.Errorf("sourceImportSet: got %v want %v", m.SourceImportSet, expectedSourceImportSet)
	}
	if m.SourceTypeDefinitionCount != ExecutionTypesSourceTypeCount {
		return fmt.Errorf("sourceTypeDefinitionCount: got %d want %d", m.SourceTypeDefinitionCount, ExecutionTypesSourceTypeCount)
	}
	if m.ExecutionSummaryFieldCount != 21 {
		return fmt.Errorf("executionSummaryFieldCount: got %d want 21", m.ExecutionSummaryFieldCount)
	}
	if err := validateJSONTags("executionSummaryJsonTags", m.ExecutionSummaryJSONTags, expectedExecutionSummaryJSONTags); err != nil {
		return err
	}
	if err := validateJSONTags("engineSpineJsonTags", m.EngineSpineJSONTags, expectedEngineSpineJSONTags); err != nil {
		return err
	}
	copied := append([]string(nil), CopiedTypeNames...)
	sort.Strings(copied)
	if !reflect.DeepEqual(m.CopiedTypes, copied) {
		return fmt.Errorf("copiedTypes: got %v want %v", m.CopiedTypes, copied)
	}
	if m.DeferredReason == "" {
		return errors.New("deferredReason must be non-empty")
	}
	return nil
}

func validateJSONTags(label string, got, want []string) error {
	if err := validateNoEmptyOrDuplicateTags(label, got); err != nil {
		return err
	}
	if !reflect.DeepEqual(got, want) {
		return fmt.Errorf("%s: got %v want %v", label, got, want)
	}
	return nil
}

func validateNoEmptyOrDuplicateTags(label string, tags []string) error {
	seen := make(map[string]struct{}, len(tags))
	for _, tag := range tags {
		if tag == "" {
			return fmt.Errorf("%s: empty json tag", label)
		}
		if _, ok := seen[tag]; ok {
			return fmt.Errorf("%s: duplicate json tag %q", label, tag)
		}
		seen[tag] = struct{}{}
	}
	return nil
}
