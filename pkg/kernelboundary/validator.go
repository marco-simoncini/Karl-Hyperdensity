package kernelboundary

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	MilestoneProductionKernelBoundary = "hyperdensity_production_kernel_boundary_v1"
	ClaimPolicyV2ID                   = "hyperdensity_claim_policy_v2"
	RuntimeActuatorName               = "FluidVirt"
)

var forbiddenPositiveClaims = []string{
	"guaranteed savings active",
	"universal performance improvement",
	"production autonomous apply",
	"windows total ram hotplug supported",
	"logical vcpu hotplug supported",
	"1000 production workloads proven",
	"dashboard is source of truth",
	"inventory hyperdensity engine",
}

// ValidateClaimPolicyV2 enforces conservative Sprint 1 claim boundaries.
func ValidateClaimPolicyV2(doc map[string]interface{}) error {
	if doc["claimPolicyId"] != ClaimPolicyV2ID {
		return fmt.Errorf("claimPolicyId must be %s", ClaimPolicyV2ID)
	}
	mustFalse := []string{
		"guaranteedSavingsAllowed",
		"universalPerformanceImprovementAllowed",
		"logicalVcpuHotplugClaimAllowed",
		"windowsTotalRamHotplugClaimAllowed",
		"ramAboveOriginalClaimAllowed",
		"productionAutonomousApplyAllowed",
		"syntheticFleetProductionClaimAllowed",
	}
	for _, k := range mustFalse {
		if v, ok := doc[k].(bool); !ok || v {
			return fmt.Errorf("%s must be false", k)
		}
	}
	return rejectForbiddenPhrasesInDoc(doc)
}

// ValidateProductionKernelBoundary checks repository responsibility contract.
func ValidateProductionKernelBoundary(doc map[string]interface{}) error {
	if doc["boundaryId"] != MilestoneProductionKernelBoundary {
		return fmt.Errorf("boundaryId must be %s", MilestoneProductionKernelBoundary)
	}
	inv, ok := doc["safetyInvariants"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("safetyInvariants required")
	}
	for _, k := range []string{"productionAutonomousApplyAllowed", "guaranteedSavingsAllowed", "universalPerformanceImprovementAllowed", "dashboardMutationSourceOfTruth", "inventoryHyperdensityEngine"} {
		if v, ok := inv[k].(bool); !ok || v {
			return fmt.Errorf("safetyInvariants.%s must be false", k)
		}
	}
	components, ok := doc["components"].([]interface{})
	if !ok {
		return fmt.Errorf("components required")
	}
	var fluidVirtActuator, dashboardProjection bool
	for _, c := range components {
		cm, ok := c.(map[string]interface{})
		if !ok {
			continue
		}
		id, _ := cm["componentId"].(string)
		ra, _ := cm["runtimeActuator"].(bool)
		po, _ := cm["projectionOnly"].(bool)
		if id == "FluidVirt" && ra {
			fluidVirtActuator = true
		}
		if id == "Karl-Dashboard" && po && !ra {
			dashboardProjection = true
		}
		if id == "Karl-Inventory" && ra {
			return fmt.Errorf("Inventory must not be runtime actuator")
		}
	}
	if !fluidVirtActuator {
		return fmt.Errorf("FluidVirt must be declared runtime actuator")
	}
	if !dashboardProjection {
		return fmt.Errorf("Dashboard must be projection only")
	}
	return rejectForbiddenPhrasesInDoc(doc)
}

// ValidateRuntimeMutationResult ensures actuator is FluidVirt and disruptive actions false.
func ValidateRuntimeMutationResult(doc map[string]interface{}) error {
	if doc["actuator"] != RuntimeActuatorName {
		return fmt.Errorf("actuator must be FluidVirt")
	}
	for _, k := range []string{"rebootUsed", "recreateUsed", "rolloutUsed", "migrationUsed"} {
		if v, ok := doc[k].(bool); !ok || v {
			return fmt.Errorf("%s must be false in Sprint 1 reference", k)
		}
	}
	return rejectForbiddenPhrasesInDoc(doc)
}

// ValidateResourceLease ensures no autonomous production apply.
func ValidateResourceLease(doc map[string]interface{}) error {
	if v, ok := doc["autoApplyAllowed"].(bool); !ok || v {
		return fmt.Errorf("autoApplyAllowed must be false")
	}
	if v, ok := doc["productionMutationAllowed"].(bool); !ok || v {
		return fmt.Errorf("productionMutationAllowed must be false")
	}
	return rejectForbiddenPhrasesInDoc(doc)
}

func rejectForbiddenPhrasesInDoc(doc map[string]interface{}) error {
	return rejectForbiddenInValue(doc)
}

var skipForbiddenScanKeys = map[string]bool{
	"forbiddenPhrases":           true,
	"forbiddenResponsibilities":  true,
	"forbiddenActions":           true,
	"blockers":                   true,
	"remediations":               true,
}

func rejectForbiddenInValue(v interface{}) error {
	switch t := v.(type) {
	case map[string]interface{}:
		for k, child := range t {
			if skipForbiddenScanKeys[k] {
				continue
			}
			if isPositiveClaimKey(k) {
				if err := scanPositiveClaimValue(child); err != nil {
					return err
				}
				continue
			}
			if err := rejectForbiddenInValue(child); err != nil {
				return err
			}
		}
	case []interface{}:
		for _, item := range t {
			if err := rejectForbiddenInValue(item); err != nil {
				return err
			}
		}
	}
	return nil
}

func isPositiveClaimKey(k string) bool {
	switch k {
	case "allowedPhrases", "conditionalPhrases", "claimBoundary", "claimBoundaries", "allowedResponsibilities", "limitation", "limitations":
		return true
	default:
		return strings.HasSuffix(k, "Phrases") && !strings.HasPrefix(k, "forbidden")
	}
}

func scanPositiveClaimValue(v interface{}) error {
	switch t := v.(type) {
	case string:
		lower := strings.ToLower(t)
		for _, phrase := range forbiddenPositiveClaims {
			if strings.Contains(lower, phrase) {
				return fmt.Errorf("forbidden positive claim phrase: %q", phrase)
			}
		}
	case []interface{}:
		for _, item := range t {
			if err := scanPositiveClaimValue(item); err != nil {
				return err
			}
		}
	}
	return nil
}

// ValidateExampleFile loads and validates a named reference example.
func ValidateExampleFile(repoRoot, fileName string) error {
	path := filepath.Join(repoRoot, "examples", fileName)
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var doc map[string]interface{}
	if err := json.Unmarshal(b, &doc); err != nil {
		return fmt.Errorf("%s: %w", fileName, err)
	}
	switch fileName {
	case "hyperdensity-claim-policy-v2-reference.json":
		return ValidateClaimPolicyV2(doc)
	case "production-kernel-boundary-reference.json":
		return ValidateProductionKernelBoundary(doc)
	case "runtime-mutation-result-reference.json":
		return ValidateRuntimeMutationResult(doc)
	case "resource-lease-reference.json":
		return ValidateResourceLease(doc)
	case "shell-passport-reference.json":
		return rejectForbiddenPhrasesInDoc(doc)
	default:
		return nil
	}
}

// ValidateSprint1Examples validates all Sprint 1 reference examples.
func ValidateSprint1Examples(repoRoot string) error {
	files := []string{
		"hyperdensity-claim-policy-v2-reference.json",
		"production-kernel-boundary-reference.json",
		"runtime-mutation-result-reference.json",
		"resource-lease-reference.json",
		"shell-passport-reference.json",
	}
	for _, f := range files {
		if err := ValidateExampleFile(repoRoot, f); err != nil {
			return fmt.Errorf("example %s: %w", f, err)
		}
	}
	return nil
}

// SchemaFilesRequired returns Sprint 1 schema basenames that must exist.
func SchemaFilesRequired() []string {
	return []string{
		"shell-passport-v1.schema.json",
		"runtime-mutation-result-v1.schema.json",
		"resource-lease-v1.schema.json",
		"hyperdensity-claim-policy-v2.schema.json",
		"production-kernel-boundary-v1.schema.json",
	}
}
