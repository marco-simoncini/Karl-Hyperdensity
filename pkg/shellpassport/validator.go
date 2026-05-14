package shellpassport

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	MilestoneShellPassportFactory = "hyperdensity_shell_passport_factory_v1"
	RegistryID                    = "hyperdensity_shell_registry_v1"
)

var canonicalShellKinds = []string{
	"container_linux",
	"container_linux_replica",
	"vm_linux",
	"vm_windows",
	"vm_windows_pool_replica",
	"daas_pool_replica",
}

var excludedShellKinds = []string{
	"container_windows",
	"container_windows_replica",
	"vm_linux_pool_replica",
}

var poolExcludedNameTokens = []string{"master", "template", "controller", "root"}

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

// ValidateShellPassportFactory checks factory reference invariants.
func ValidateShellPassportFactory(doc map[string]interface{}) error {
	if doc["factoryId"] != MilestoneShellPassportFactory {
		return fmt.Errorf("factoryId must be %s", MilestoneShellPassportFactory)
	}
	if v, ok := doc["productionMutationAllowed"].(bool); !ok || v {
		return fmt.Errorf("productionMutationAllowed must be false")
	}
	if v, ok := doc["autoApplyAllowed"].(bool); !ok || v {
		return fmt.Errorf("autoApplyAllowed must be false")
	}
	kinds, ok := doc["canonicalShellKinds"].([]interface{})
	if !ok || len(kinds) < len(canonicalShellKinds) {
		return fmt.Errorf("canonicalShellKinds incomplete")
	}
	for _, want := range canonicalShellKinds {
		if !sliceContainsString(kinds, want) {
			return fmt.Errorf("missing canonical shell kind %s", want)
		}
	}
	return rejectForbiddenPositiveClaims(doc)
}

// ValidateShellRegistrySurface validates ConfigMap-ready registry surface.
func ValidateShellRegistrySurface(doc map[string]interface{}) error {
	if doc["registryId"] != RegistryID {
		return fmt.Errorf("registryId must be %s", RegistryID)
	}
	shells, ok := doc["shells"].([]interface{})
	if !ok || len(shells) == 0 {
		return fmt.Errorf("shells required")
	}
	var donor, receiver, blocked, remediable, protected int
	kindsSeen := map[string]bool{}
	for _, item := range shells {
		sh, ok := item.(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid shell entry")
		}
		if err := validateShellEntry(sh); err != nil {
			return err
		}
		kind, _ := sh["shellKind"].(string)
		kindsSeen[kind] = true
		if isEligible(sh, "donor") {
			donor++
		}
		if isEligible(sh, "receiver") {
			receiver++
		}
		if b, _ := sh["blocked"].(bool); b {
			blocked++
		}
		if r, _ := sh["remediable"].(bool); r {
			remediable++
		}
		if p, _ := sh["protected"].(bool); p {
			protected++
		}
	}
	for _, k := range canonicalShellKinds {
		if !kindsSeen[k] {
			return fmt.Errorf("registry missing canonical shell kind example: %s", k)
		}
	}
	if donor < 1 || receiver < 1 || blocked < 1 || remediable < 1 || protected < 1 {
		return fmt.Errorf("registry must include donor, receiver, blocked, remediable, and protected shells")
	}
	return rejectForbiddenPositiveClaims(doc)
}

func validateShellEntry(sh map[string]interface{}) error {
	if id, _ := sh["shellId"].(string); id == "" {
		return fmt.Errorf("shell missing shellId")
	}
	cb, ok := sh["claimBoundary"].([]interface{})
	if !ok || len(cb) == 0 {
		return fmt.Errorf("shell %v missing claimBoundary", sh["shellId"])
	}
	if sp, _ := sh["supportProfile"].(string); sp == "" {
		return fmt.Errorf("shell %v missing supportProfile", sh["shellId"])
	}
	kind, _ := sh["shellKind"].(string)
	for _, ex := range excludedShellKinds {
		if kind == ex {
			return fmt.Errorf("excluded shell kind enrolled without override: %s", kind)
		}
	}
	targetRef, _ := sh["targetRef"].(string)
	low := strings.ToLower(targetRef)
	for _, tok := range poolExcludedNameTokens {
		if strings.Contains(low, tok) && (kind == "vm_windows_pool_replica" || kind == "daas_pool_replica") {
			if strings.Contains(low, "master-win11") || strings.Contains(low, "template") || strings.Contains(low, "controller") {
				return fmt.Errorf("pool excluded member must not be eligible: %s", targetRef)
			}
		}
	}
	blocked, _ := sh["blocked"].(bool)
	if blocked {
		if isEligible(sh, "donor") || isEligible(sh, "receiver") {
			return fmt.Errorf("blocked shell cannot be donor/receiver eligible: %v", sh["shellId"])
		}
	}
	mem, _ := sh["memoryEnvelope"].(map[string]interface{})
	if mem != nil {
		raw, _ := json.Marshal(mem)
		lowMem := strings.ToLower(string(raw))
		if strings.Contains(lowMem, "total_ram_hotplug") || strings.Contains(lowMem, "windows_total_ram_hotplug") {
			return fmt.Errorf("windows total ram hotplug overclaim on shell %v", sh["shellId"])
		}
	}
	cpu, _ := sh["cpuEntitlement"].(map[string]interface{})
	if cpu != nil {
		raw, _ := json.Marshal(cpu)
		lowCPU := strings.ToLower(string(raw))
		if strings.Contains(lowCPU, "logical_vcpu_hotplug") || strings.Contains(lowCPU, "logical vcpu hotplug") {
			return fmt.Errorf("logical vcpu hotplug overclaim on shell %v", sh["shellId"])
		}
	}
	return rejectForbiddenPositiveClaims(sh)
}

func isEligible(sh map[string]interface{}, role string) bool {
	var field string
	if role == "donor" {
		field = "donorEligibility"
	} else {
		field = "receiverEligibility"
	}
	v, _ := sh[field].(string)
	return strings.EqualFold(v, "eligible")
}

func ValidateShellEnrollmentResult(doc map[string]interface{}) error {
	if v, ok := doc["productionMutationAllowed"].(bool); !ok || v {
		return fmt.Errorf("productionMutationAllowed must be false")
	}
	if v, ok := doc["autoApplyAllowed"].(bool); !ok || v {
		return fmt.Errorf("autoApplyAllowed must be false")
	}
	blocked, _ := doc["blocked"].(bool)
	if blocked {
		if d, _ := doc["eligibleAsDonor"].(bool); d {
			return fmt.Errorf("blocked shell cannot be eligibleAsDonor")
		}
		if r, _ := doc["eligibleAsReceiver"].(bool); r {
			return fmt.Errorf("blocked shell cannot be eligibleAsReceiver")
		}
	}
	return rejectForbiddenPositiveClaims(doc)
}

func ValidateShellCapabilityEvidence(doc map[string]interface{}) error {
	if doc["source"] != "FluidVirt" || doc["actuator"] != "FluidVirt" {
		return fmt.Errorf("capability evidence must come from FluidVirt actuator")
	}
	for _, k := range []string{"rebootRequired", "recreateRequired", "migrationRequired"} {
		if v, ok := doc[k].(bool); !ok || v {
			return fmt.Errorf("%s must be false for in-place movement certification", k)
		}
	}
	unsupported, _ := doc["unsupportedFamilies"].([]interface{})
	for _, item := range unsupported {
		s, _ := item.(string)
		if s == "windows_total_ram_hotplug" || s == "logical_vcpu_hotplug" {
			continue
		}
	}
	supported, _ := doc["runtimeMutationSupportedFamilies"].([]interface{})
	for _, item := range supported {
		s, _ := item.(string)
		if s == "logical_vcpu_hotplug" || s == "windows_total_ram_hotplug" {
			return fmt.Errorf("unsupported family must not appear in supported list: %s", s)
		}
	}
	return rejectForbiddenPositiveClaims(doc)
}

func rejectForbiddenPositiveClaims(v interface{}) error {
	positives := collectPositiveStrings(v)
	merged := strings.ToLower(strings.Join(positives, "\n"))
	for _, phrase := range forbiddenPositiveClaims {
		if strings.Contains(merged, phrase) {
			return fmt.Errorf("forbidden positive claim: %q", phrase)
		}
	}
	return nil
}

var skipKeys = map[string]bool{
	"forbiddenPhrases": true, "blockerCodes": true, "blockers": true,
	"remediationCodes": true, "remediations": true, "unsupportedFamilies": true,
	"excludedShellKinds": true,
}

func collectPositiveStrings(v interface{}) []string {
	var out []string
	switch t := v.(type) {
	case map[string]interface{}:
		for k, child := range t {
			if skipKeys[k] {
				continue
			}
			if isPositiveKey(k) {
				out = append(out, flattenStrings(child)...)
			} else {
				out = append(out, collectPositiveStrings(child)...)
			}
		}
	case []interface{}:
		for _, item := range t {
			out = append(out, collectPositiveStrings(item)...)
		}
	}
	return out
}

func isPositiveKey(k string) bool {
	switch k {
	case "claimBoundary", "claimBoundaries", "allowedPhrases", "conditionalPhrases", "donorEligibility", "receiverEligibility":
		return true
	default:
		return strings.HasSuffix(k, "Phrases") && !strings.HasPrefix(k, "forbidden")
	}
}

func flattenStrings(v interface{}) []string {
	switch t := v.(type) {
	case string:
		return []string{t}
	case []interface{}:
		var out []string
		for _, item := range t {
			if s, ok := item.(string); ok {
				out = append(out, s)
			}
		}
		return out
	}
	return nil
}

func sliceContainsString(arr []interface{}, want string) bool {
	for _, item := range arr {
		if s, ok := item.(string); ok && s == want {
			return true
		}
	}
	return false
}

// ValidateSprint2Examples validates all Sprint 2 reference examples.
func ValidateSprint2Examples(repoRoot string) error {
	files := map[string]func(map[string]interface{}) error{
		"shell-passport-factory-reference.json":      ValidateShellPassportFactory,
		"shell-passport-registry-reference.json":     ValidateShellRegistrySurface,
		"shell-registry-reference.json":              validateShellRegistryMinimal,
		"shell-enrollment-result-reference.json":     ValidateShellEnrollmentResult,
		"shell-capability-evidence-reference.json":   ValidateShellCapabilityEvidence,
	}
	for name, fn := range files {
		b, err := os.ReadFile(filepath.Join(repoRoot, "examples", name))
		if err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
		var doc map[string]interface{}
		if err := json.Unmarshal(b, &doc); err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
		if err := fn(doc); err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
	}
	return nil
}

func validateShellRegistryMinimal(doc map[string]interface{}) error {
	if doc["registryId"] != RegistryID {
		return fmt.Errorf("registryId must be %s", RegistryID)
	}
	shells, ok := doc["shells"].([]interface{})
	if !ok {
		return fmt.Errorf("shells array required")
	}
	for _, item := range shells {
		sh, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		if err := validateShellEntry(sh); err != nil {
			return err
		}
	}
	return nil
}

// SchemaFilesRequiredSprint2 returns Sprint 2 schema basenames.
func SchemaFilesRequiredSprint2() []string {
	return []string{
		"shell-passport-factory-v1.schema.json",
		"shell-registry-v1.schema.json",
		"shell-enrollment-result-v1.schema.json",
		"shell-capability-evidence-v1.schema.json",
	}
}
