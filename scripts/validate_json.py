#!/usr/bin/env python3
"""Validate Hyperdensity JSON schemas and examples with stdlib only."""

from __future__ import annotations

import json
import sys
from pathlib import Path


def fail(message: str) -> None:
    print(f"[validate_json] ERROR: {message}", file=sys.stderr)
    raise SystemExit(1)


def load_json(path: Path) -> dict:
    try:
        data = json.loads(path.read_text(encoding="utf-8"))
    except Exception as exc:  # pragma: no cover
        fail(f"{path} is not valid JSON: {exc}")
    if not isinstance(data, dict):
        fail(f"{path} must contain a top-level JSON object")
    return data


def load_text(path: Path) -> str:
    try:
        return path.read_text(encoding="utf-8")
    except Exception as exc:  # pragma: no cover
        fail(f"{path} is not readable text: {exc}")


def main() -> int:
    repo_root = Path(__file__).resolve().parents[1]
    schema_paths = sorted((repo_root / "schemas").glob("*.json"))
    example_paths = sorted((repo_root / "examples").glob("*.json"))

    if not schema_paths:
        fail("no schema files found under schemas/*.json")
    if not example_paths:
        fail("no example files found under examples/*.json")

    schema_count = 0
    for path in schema_paths:
        data = load_json(path)
        schema_count += 1

        if "$schema" not in data:
            fail(f"{path} is missing required '$schema'")
        if "$id" not in data and "title" not in data:
            fail(f"{path} must contain either '$id' or 'title'")
        if "type" not in data:
            fail(f"{path} is missing required 'type'")

    example_count = 0
    example_by_name = {}
    for path in example_paths:
        data = load_json(path)
        example_by_name[path.name] = data
        example_count += 1

    policy_pack_example = example_by_name.get("policy-pack-reference.json")
    if policy_pack_example is None:
        fail("examples/policy-pack-reference.json is missing")
    required_policy_fields = [
        "policyPackId",
        "policyPackVersion",
        "policyPackMode",
        "enforcementMode",
        "autonomousApplyAllowed",
        "supportedShellKinds",
        "supportedProfiles",
        "factoryRequirements",
        "claimValidationRules",
        "admissionGuardRules",
        "mutatePreviewDefaults",
        "exchangeEligibilityRules",
        "stageApplyRules",
        "shellClaimEvidenceCreateRules",
        "safetyGates",
        "warmupPolicy",
    ]
    for field in required_policy_fields:
        if field not in policy_pack_example:
            fail(f"policy-pack-reference.json missing required field '{field}'")
    if policy_pack_example["policyPackId"] != "hyperdensity_policy_pack_v1":
        fail("policy-pack-reference.json policyPackId must be hyperdensity_policy_pack_v1")
    if policy_pack_example["policyPackMode"] != "visibility_only":
        fail("policy-pack-reference.json policyPackMode must be visibility_only")
    if policy_pack_example["enforcementMode"] != "disabled":
        fail("policy-pack-reference.json enforcementMode must be disabled")
    if policy_pack_example["autonomousApplyAllowed"] is not False:
        fail("policy-pack-reference.json autonomousApplyAllowed must be false")
    if not isinstance(policy_pack_example["supportedShellKinds"], list) or not policy_pack_example["supportedShellKinds"]:
        fail("policy-pack-reference.json supportedShellKinds must be a non-empty array")
    if not isinstance(policy_pack_example["supportedProfiles"], list) or not policy_pack_example["supportedProfiles"]:
        fail("policy-pack-reference.json supportedProfiles must be a non-empty array")
    if not isinstance(policy_pack_example["safetyGates"], list) or not policy_pack_example["safetyGates"]:
        fail("policy-pack-reference.json safetyGates must be a non-empty array")

    consistency_example = example_by_name.get("policy-pack-consistency-reference.json")
    if consistency_example is None:
        fail("examples/policy-pack-consistency-reference.json is missing")
    required_consistency_fields = [
        "consistencyId",
        "consistencyVersion",
        "consistencyMode",
        "consistencyState",
        "policyPackId",
        "policyPackVersion",
        "checkedAt",
        "checkedSectionCount",
        "checkedRuleCount",
        "checkedSafetyGateCount",
        "expectedSectionCount",
        "expectedSafetyGateCount",
        "missingSections",
        "missingRules",
        "missingSafetyGates",
        "driftFindings",
        "sourceSurfaceFindings",
        "invariantFindings",
        "consistent",
        "nextConsistencyAction",
    ]
    for field in required_consistency_fields:
        if field not in consistency_example:
            fail(f"policy-pack-consistency-reference.json missing required field '{field}'")
    if consistency_example["consistencyId"] != "hyperdensity_policy_pack_consistency_checker_v1":
        fail("policy-pack-consistency-reference.json consistencyId must be hyperdensity_policy_pack_consistency_checker_v1")
    if consistency_example["consistencyMode"] != "validation_only":
        fail("policy-pack-consistency-reference.json consistencyMode must be validation_only")
    if consistency_example["policyPackId"] != "hyperdensity_policy_pack_v1":
        fail("policy-pack-consistency-reference.json policyPackId must be hyperdensity_policy_pack_v1")
    if consistency_example["consistent"] is not True:
        fail("policy-pack-consistency-reference.json consistent must be true")
    if not isinstance(consistency_example["missingSections"], list):
        fail("policy-pack-consistency-reference.json missingSections must be an array")
    if not isinstance(consistency_example["missingSafetyGates"], list):
        fail("policy-pack-consistency-reference.json missingSafetyGates must be an array")

    enforce_simulation_example = example_by_name.get("admission-guard-enforce-simulation-reference.json")
    if enforce_simulation_example is None:
        fail("examples/admission-guard-enforce-simulation-reference.json is missing")
    required_enforce_simulation_fields = [
        "simulationId",
        "simulationVersion",
        "simulationMode",
        "enforcementMode",
        "admissionGuardMode",
        "mutatePreviewMode",
        "autonomousApplyAllowed",
        "policyPackId",
        "policyConsistencyRequired",
        "evidenceScope",
        "productionMutationAllowed",
        "simulatedObjects",
        "simulatedDecisions",
        "summary",
        "safetyGates",
        "nextSimulationAction",
    ]
    for field in required_enforce_simulation_fields:
        if field not in enforce_simulation_example:
            fail(f"admission-guard-enforce-simulation-reference.json missing required field '{field}'")
    if enforce_simulation_example["simulationId"] != "hyperdensity_admission_guard_enforce_simulation_v1":
        fail("admission-guard-enforce-simulation-reference.json simulationId must be hyperdensity_admission_guard_enforce_simulation_v1")
    if enforce_simulation_example["simulationMode"] != "simulation_only":
        fail("admission-guard-enforce-simulation-reference.json simulationMode must be simulation_only")
    if enforce_simulation_example["enforcementMode"] != "disabled":
        fail("admission-guard-enforce-simulation-reference.json enforcementMode must be disabled")
    if enforce_simulation_example["admissionGuardMode"] != "audit_only":
        fail("admission-guard-enforce-simulation-reference.json admissionGuardMode must be audit_only")
    if enforce_simulation_example["mutatePreviewMode"] != "audit_preview_only":
        fail("admission-guard-enforce-simulation-reference.json mutatePreviewMode must be audit_preview_only")
    if enforce_simulation_example["autonomousApplyAllowed"] is not False:
        fail("admission-guard-enforce-simulation-reference.json autonomousApplyAllowed must be false")
    if enforce_simulation_example["productionMutationAllowed"] is not False:
        fail("admission-guard-enforce-simulation-reference.json productionMutationAllowed must be false")

    mutate_preview_apply_dry_run_example = example_by_name.get("mutate-preview-apply-dry-run-reference.json")
    if mutate_preview_apply_dry_run_example is None:
        fail("examples/mutate-preview-apply-dry-run-reference.json is missing")
    required_mutate_preview_apply_dry_run_fields = [
        "dryRunId",
        "dryRunVersion",
        "dryRunMode",
        "mutationAllowed",
        "productionMutationAllowed",
        "enforcementMode",
        "admissionGuardMode",
        "mutatePreviewMode",
        "autonomousApplyAllowed",
        "evidenceScope",
        "policyPackId",
        "policyConsistencyRequired",
        "sourceSurface",
        "dryRunTargets",
        "dryRunResults",
        "safetyGates",
        "summary",
        "nextDryRunAction",
    ]
    for field in required_mutate_preview_apply_dry_run_fields:
        if field not in mutate_preview_apply_dry_run_example:
            fail(f"mutate-preview-apply-dry-run-reference.json missing required field '{field}'")
    if mutate_preview_apply_dry_run_example["dryRunId"] != "hyperdensity_mutate_preview_apply_dry_run_v1":
        fail("mutate-preview-apply-dry-run-reference.json dryRunId must be hyperdensity_mutate_preview_apply_dry_run_v1")
    if mutate_preview_apply_dry_run_example["dryRunVersion"] != "v1":
        fail("mutate-preview-apply-dry-run-reference.json dryRunVersion must be v1")
    if mutate_preview_apply_dry_run_example["dryRunMode"] != "server_side_apply_dry_run_only":
        fail("mutate-preview-apply-dry-run-reference.json dryRunMode must be server_side_apply_dry_run_only")
    if mutate_preview_apply_dry_run_example["mutationAllowed"] is not False:
        fail("mutate-preview-apply-dry-run-reference.json mutationAllowed must be false")
    if mutate_preview_apply_dry_run_example["productionMutationAllowed"] is not False:
        fail("mutate-preview-apply-dry-run-reference.json productionMutationAllowed must be false")
    if mutate_preview_apply_dry_run_example["enforcementMode"] != "disabled":
        fail("mutate-preview-apply-dry-run-reference.json enforcementMode must be disabled")
    if mutate_preview_apply_dry_run_example["admissionGuardMode"] != "audit_only":
        fail("mutate-preview-apply-dry-run-reference.json admissionGuardMode must be audit_only")
    if mutate_preview_apply_dry_run_example["mutatePreviewMode"] != "audit_preview_only":
        fail("mutate-preview-apply-dry-run-reference.json mutatePreviewMode must be audit_preview_only")
    if mutate_preview_apply_dry_run_example["autonomousApplyAllowed"] is not False:
        fail("mutate-preview-apply-dry-run-reference.json autonomousApplyAllowed must be false")
    if mutate_preview_apply_dry_run_example["evidenceScope"] != "evidence_namespace_only":
        fail("mutate-preview-apply-dry-run-reference.json evidenceScope must be evidence_namespace_only")
    if mutate_preview_apply_dry_run_example["policyPackId"] != "hyperdensity_policy_pack_v1":
        fail("mutate-preview-apply-dry-run-reference.json policyPackId must be hyperdensity_policy_pack_v1")
    if mutate_preview_apply_dry_run_example["policyConsistencyRequired"] is not True:
        fail("mutate-preview-apply-dry-run-reference.json policyConsistencyRequired must be true")
    if mutate_preview_apply_dry_run_example["sourceSurface"] != "admission_guard_mutate_preview":
        fail("mutate-preview-apply-dry-run-reference.json sourceSurface must be admission_guard_mutate_preview")
    if not isinstance(mutate_preview_apply_dry_run_example["dryRunTargets"], list) or not mutate_preview_apply_dry_run_example["dryRunTargets"]:
        fail("mutate-preview-apply-dry-run-reference.json dryRunTargets must be a non-empty array")
    if not isinstance(mutate_preview_apply_dry_run_example["dryRunResults"], list) or not mutate_preview_apply_dry_run_example["dryRunResults"]:
        fail("mutate-preview-apply-dry-run-reference.json dryRunResults must be a non-empty array")
    if not isinstance(mutate_preview_apply_dry_run_example["safetyGates"], list) or not mutate_preview_apply_dry_run_example["safetyGates"]:
        fail("mutate-preview-apply-dry-run-reference.json safetyGates must be a non-empty array")

    release_support_matrix_example = example_by_name.get("release-support-matrix-reference.json")
    if release_support_matrix_example is None:
        fail("examples/release-support-matrix-reference.json is missing")
    required_release_matrix_fields = [
        "supportMatrixId",
        "supportMatrixVersion",
        "releaseTrack",
        "matrixMode",
        "supportClaimMode",
        "supportMatrixState",
        "enforcementMode",
        "autonomousApplyAllowed",
        "productionMutationAllowed",
        "evidenceScope",
        "policyPackId",
        "policyConsistencyRequired",
        "profilePackId",
        "supportedShellKinds",
        "supportedProfiles",
        "capabilityMatrix",
        "operationMatrix",
        "surfaceMatrix",
        "supportLevels",
        "safetyGates",
        "proofCatalog",
        "limitations",
        "outOfScope",
        "nextReleaseActions",
        "summary",
    ]
    for field in required_release_matrix_fields:
        if field not in release_support_matrix_example:
            fail(f"release-support-matrix-reference.json missing required field '{field}'")
    if release_support_matrix_example["supportMatrixId"] != "hyperdensity_release_support_matrix_v1":
        fail("release-support-matrix-reference.json supportMatrixId must be hyperdensity_release_support_matrix_v1")
    if release_support_matrix_example["supportMatrixVersion"] != "v1":
        fail("release-support-matrix-reference.json supportMatrixVersion must be v1")
    if release_support_matrix_example["releaseTrack"] != "technical_preview":
        fail("release-support-matrix-reference.json releaseTrack must be technical_preview")
    if release_support_matrix_example["matrixMode"] != "release_boundary_visibility_only":
        fail("release-support-matrix-reference.json matrixMode must be release_boundary_visibility_only")
    if release_support_matrix_example["supportClaimMode"] != "evidence_backed_only":
        fail("release-support-matrix-reference.json supportClaimMode must be evidence_backed_only")
    if release_support_matrix_example["enforcementMode"] != "disabled":
        fail("release-support-matrix-reference.json enforcementMode must be disabled")
    if release_support_matrix_example["autonomousApplyAllowed"] is not False:
        fail("release-support-matrix-reference.json autonomousApplyAllowed must be false")
    if release_support_matrix_example["productionMutationAllowed"] is not False:
        fail("release-support-matrix-reference.json productionMutationAllowed must be false")
    if release_support_matrix_example["evidenceScope"] != "evidence_namespace_only":
        fail("release-support-matrix-reference.json evidenceScope must be evidence_namespace_only")
    if release_support_matrix_example["policyPackId"] != "hyperdensity_policy_pack_v1":
        fail("release-support-matrix-reference.json policyPackId must be hyperdensity_policy_pack_v1")
    if release_support_matrix_example["policyConsistencyRequired"] is not True:
        fail("release-support-matrix-reference.json policyConsistencyRequired must be true")
    if release_support_matrix_example["profilePackId"] != "hyperdensity_shell_claim_templates_profile_pack_v1":
        fail("release-support-matrix-reference.json profilePackId must be hyperdensity_shell_claim_templates_profile_pack_v1")
    if not isinstance(release_support_matrix_example["supportedShellKinds"], list) or not release_support_matrix_example["supportedShellKinds"]:
        fail("release-support-matrix-reference.json supportedShellKinds must be a non-empty array")
    if not isinstance(release_support_matrix_example["supportedProfiles"], list) or len(release_support_matrix_example["supportedProfiles"]) < 7:
        fail("release-support-matrix-reference.json supportedProfiles must contain at least 7 entries")
    if not isinstance(release_support_matrix_example["supportLevels"], list) or len(release_support_matrix_example["supportLevels"]) < 10:
        fail("release-support-matrix-reference.json supportLevels must contain all required support levels")
    if "windows_container" in release_support_matrix_example["supportedShellKinds"] or "windows_vm" in release_support_matrix_example["supportedShellKinds"]:
        fail("release-support-matrix-reference.json supportedShellKinds must not include Windows lanes")

    evidence_bundle_example = example_by_name.get("evidence-bundle-demo-scenario-pack-reference.json")
    if evidence_bundle_example is None:
        fail("examples/evidence-bundle-demo-scenario-pack-reference.json is missing")
    required_evidence_bundle_fields = [
        "evidenceBundleId",
        "evidenceBundleVersion",
        "releaseTrack",
        "bundleMode",
        "demoMode",
        "supportMatrixId",
        "policyPackId",
        "profilePackId",
        "enforcementMode",
        "autonomousApplyAllowed",
        "productionMutationAllowed",
        "evidenceScope",
        "supportedClaims",
        "demoScenarios",
        "evidenceCatalog",
        "artifactIndex",
        "proofChecks",
        "safetyGates",
        "demoRunbook",
        "knownLimitations",
        "outOfScope",
        "nextReleaseActions",
        "summary",
    ]
    for field in required_evidence_bundle_fields:
        if field not in evidence_bundle_example:
            fail(f"evidence-bundle-demo-scenario-pack-reference.json missing required field '{field}'")
    if evidence_bundle_example["evidenceBundleId"] != "hyperdensity_evidence_bundle_demo_scenario_pack_v1":
        fail("evidence-bundle-demo-scenario-pack-reference.json evidenceBundleId must be hyperdensity_evidence_bundle_demo_scenario_pack_v1")
    if evidence_bundle_example["evidenceBundleVersion"] != "v1":
        fail("evidence-bundle-demo-scenario-pack-reference.json evidenceBundleVersion must be v1")
    if evidence_bundle_example["releaseTrack"] != "technical_preview":
        fail("evidence-bundle-demo-scenario-pack-reference.json releaseTrack must be technical_preview")
    if evidence_bundle_example["bundleMode"] != "evidence_bundle_only":
        fail("evidence-bundle-demo-scenario-pack-reference.json bundleMode must be evidence_bundle_only")
    if evidence_bundle_example["demoMode"] != "guided_operator_demo":
        fail("evidence-bundle-demo-scenario-pack-reference.json demoMode must be guided_operator_demo")
    if evidence_bundle_example["supportMatrixId"] != "hyperdensity_release_support_matrix_v1":
        fail("evidence-bundle-demo-scenario-pack-reference.json supportMatrixId must be hyperdensity_release_support_matrix_v1")
    if evidence_bundle_example["policyPackId"] != "hyperdensity_policy_pack_v1":
        fail("evidence-bundle-demo-scenario-pack-reference.json policyPackId must be hyperdensity_policy_pack_v1")
    if evidence_bundle_example["profilePackId"] != "hyperdensity_shell_claim_templates_profile_pack_v1":
        fail("evidence-bundle-demo-scenario-pack-reference.json profilePackId must be hyperdensity_shell_claim_templates_profile_pack_v1")
    if evidence_bundle_example["enforcementMode"] != "disabled":
        fail("evidence-bundle-demo-scenario-pack-reference.json enforcementMode must be disabled")
    if evidence_bundle_example["autonomousApplyAllowed"] is not False:
        fail("evidence-bundle-demo-scenario-pack-reference.json autonomousApplyAllowed must be false")
    if evidence_bundle_example["productionMutationAllowed"] is not False:
        fail("evidence-bundle-demo-scenario-pack-reference.json productionMutationAllowed must be false")
    if evidence_bundle_example["evidenceScope"] != "evidence_namespace_only":
        fail("evidence-bundle-demo-scenario-pack-reference.json evidenceScope must be evidence_namespace_only")
    if not isinstance(evidence_bundle_example["supportedClaims"], list) or not evidence_bundle_example["supportedClaims"]:
        fail("evidence-bundle-demo-scenario-pack-reference.json supportedClaims must be non-empty array")
    if not isinstance(evidence_bundle_example["demoScenarios"], list) or len(evidence_bundle_example["demoScenarios"]) < 5:
        fail("evidence-bundle-demo-scenario-pack-reference.json demoScenarios must contain at least 5 entries")
    if not isinstance(evidence_bundle_example["evidenceCatalog"], list) or not evidence_bundle_example["evidenceCatalog"]:
        fail("evidence-bundle-demo-scenario-pack-reference.json evidenceCatalog must be non-empty array")
    if not isinstance(evidence_bundle_example["proofChecks"], list) or not evidence_bundle_example["proofChecks"]:
        fail("evidence-bundle-demo-scenario-pack-reference.json proofChecks must be non-empty array")
    if not isinstance(evidence_bundle_example["safetyGates"], list) or not evidence_bundle_example["safetyGates"]:
        fail("evidence-bundle-demo-scenario-pack-reference.json safetyGates must be non-empty array")

    live_resource_authority_example = example_by_name.get("live-resource-authority-reference.json")
    if live_resource_authority_example is None:
        fail("examples/live-resource-authority-reference.json is missing")
    required_live_resource_authority_fields = [
        "authorityId",
        "authorityName",
        "authorityVersion",
        "authorityMode",
        "authorityState",
        "releaseTrack",
        "operationMode",
        "autonomousApplyAllowed",
        "enforcementMode",
        "productionMutationAllowed",
        "evidenceScope",
        "supportMatrixId",
        "evidenceBundleId",
        "policyPackId",
        "profilePackId",
        "uniformContract",
        "runtimeDrivers",
        "liveResourceIntents",
        "runtimeLeases",
        "capabilityChecks",
        "preflightChecks",
        "dryRunSemantics",
        "applySemantics",
        "verificationSemantics",
        "rollbackSemantics",
        "reconciliationSemantics",
        "safetyGates",
        "supportBoundaries",
        "limitations",
        "outOfScope",
        "summary",
        "nextAuthorityAction",
        "authorityBlocker",
    ]
    for field in required_live_resource_authority_fields:
        if field not in live_resource_authority_example:
            fail(f"live-resource-authority-reference.json missing required field '{field}'")
    if live_resource_authority_example["authorityId"] != "hyperdensity_live_resource_authority_v1":
        fail("live-resource-authority-reference.json authorityId must be hyperdensity_live_resource_authority_v1")
    if live_resource_authority_example["authorityName"] != "KARL Live Resource Authority":
        fail("live-resource-authority-reference.json authorityName must be KARL Live Resource Authority")
    if live_resource_authority_example["authorityVersion"] != "v1":
        fail("live-resource-authority-reference.json authorityVersion must be v1")
    if live_resource_authority_example["authorityMode"] != "unified_runtime_control_surface":
        fail("live-resource-authority-reference.json authorityMode must be unified_runtime_control_surface")
    if live_resource_authority_example["releaseTrack"] != "technical_preview":
        fail("live-resource-authority-reference.json releaseTrack must be technical_preview")
    if live_resource_authority_example["operationMode"] != "operator_controlled_projection":
        fail("live-resource-authority-reference.json operationMode must be operator_controlled_projection")
    if live_resource_authority_example["autonomousApplyAllowed"] is not False:
        fail("live-resource-authority-reference.json autonomousApplyAllowed must be false")
    if live_resource_authority_example["enforcementMode"] != "disabled":
        fail("live-resource-authority-reference.json enforcementMode must be disabled")
    if live_resource_authority_example["productionMutationAllowed"] is not False:
        fail("live-resource-authority-reference.json productionMutationAllowed must be false")
    if live_resource_authority_example["evidenceScope"] != "evidence_namespace_only":
        fail("live-resource-authority-reference.json evidenceScope must be evidence_namespace_only")
    if live_resource_authority_example["supportMatrixId"] != "hyperdensity_release_support_matrix_v1":
        fail("live-resource-authority-reference.json supportMatrixId must be hyperdensity_release_support_matrix_v1")
    if live_resource_authority_example["evidenceBundleId"] != "hyperdensity_evidence_bundle_demo_scenario_pack_v1":
        fail("live-resource-authority-reference.json evidenceBundleId must be hyperdensity_evidence_bundle_demo_scenario_pack_v1")
    if live_resource_authority_example["policyPackId"] != "hyperdensity_policy_pack_v1":
        fail("live-resource-authority-reference.json policyPackId must be hyperdensity_policy_pack_v1")
    if live_resource_authority_example["profilePackId"] != "hyperdensity_shell_claim_templates_profile_pack_v1":
        fail("live-resource-authority-reference.json profilePackId must be hyperdensity_shell_claim_templates_profile_pack_v1")
    if not isinstance(live_resource_authority_example["runtimeDrivers"], list) or len(live_resource_authority_example["runtimeDrivers"]) < 3:
        fail("live-resource-authority-reference.json runtimeDrivers must include at least 3 entries")
    driver_ids = {entry.get("driverId") for entry in live_resource_authority_example["runtimeDrivers"] if isinstance(entry, dict)}
    for driver_id in ("container_linux_pod_resize_driver", "vm_linux_cpu_libvirt_qga_driver", "vm_linux_memory_virtiomem_qmp_driver"):
        if driver_id not in driver_ids:
            fail(f"live-resource-authority-reference.json runtimeDrivers must include {driver_id}")
    for driver in live_resource_authority_example["runtimeDrivers"]:
        if isinstance(driver, dict) and driver.get("exposedAsRawControl") is True:
            fail("live-resource-authority-reference.json runtimeDrivers must not expose raw controls")
    if not isinstance(live_resource_authority_example["uniformContract"], list) or len(live_resource_authority_example["uniformContract"]) < 10:
        fail("live-resource-authority-reference.json uniformContract must contain all required phases")
    phase_ids = {entry.get("phaseId") for entry in live_resource_authority_example["uniformContract"] if isinstance(entry, dict)}
    for phase_id in (
        "live_resource_intent",
        "capability_check",
        "preflight",
        "dry_run_or_dry_run_like_validation",
        "runtime_lease_or_overlay",
        "apply",
        "verify",
        "audit",
        "rollback",
        "reconcile_or_expire",
    ):
        if phase_id not in phase_ids:
            fail(f"live-resource-authority-reference.json uniformContract missing phase {phase_id}")
    if not isinstance(live_resource_authority_example["liveResourceIntents"], list) or not live_resource_authority_example["liveResourceIntents"]:
        fail("live-resource-authority-reference.json liveResourceIntents must be non-empty")
    intent = live_resource_authority_example["liveResourceIntents"][0]
    if not isinstance(intent, dict) or intent.get("rollbackRequired") is not True or intent.get("verificationRequired") is not True:
        fail("live-resource-authority-reference.json liveResourceIntents must require rollback and verification")
    out_of_scope_ids = {entry.get("itemId") for entry in live_resource_authority_example["outOfScope"] if isinstance(entry, dict)}
    if "windows_lane" not in out_of_scope_ids:
        fail("live-resource-authority-reference.json outOfScope must include windows_lane")
    support_boundaries = live_resource_authority_example["supportBoundaries"]
    if not isinstance(support_boundaries, list) or not support_boundaries:
        fail("live-resource-authority-reference.json supportBoundaries must be non-empty")
    rejected_phrases = set()
    for boundary in support_boundaries:
        if isinstance(boundary, dict):
            for phrase in boundary.get("rejectedWording", []):
                if isinstance(phrase, str):
                    rejected_phrases.add(phrase)
    if "VM RAM live resize is generic KubeVirt template mutation." not in rejected_phrases:
        fail("live-resource-authority-reference.json supportBoundaries must reject generic KubeVirt memory template mutation wording")
    if "Windows is supported." not in rejected_phrases:
        fail("live-resource-authority-reference.json supportBoundaries must reject Windows support wording")

    action_slate_example = example_by_name.get("action-slate-reference.json")
    if action_slate_example is None:
        fail("examples/action-slate-reference.json is missing")
    required_action_slate_fields = [
        "actionSlateId",
        "actionSlateVersion",
        "actionSlateState",
        "actionSlateMode",
        "releaseTrack",
        "generatedAt",
        "validForSeconds",
        "expiresAt",
        "operatorControlled",
        "autonomousApplyAllowed",
        "enforcementMode",
        "productionMutationAllowed",
        "evidenceScope",
        "policyPackId",
        "supportMatrixId",
        "liveResourceAuthorityId",
        "evidenceBundleId",
        "actionCount",
        "readyActionCount",
        "blockedActionCount",
        "expiringActionCount",
        "actions",
        "indexes",
        "safetyGates",
        "blockers",
        "limitations",
        "claimBoundaries",
    ]
    for field in required_action_slate_fields:
        if field not in action_slate_example:
            fail(f"action-slate-reference.json missing required field '{field}'")
    if action_slate_example["actionSlateId"] != "hyperdensity_action_slate_v1":
        fail("action-slate-reference.json actionSlateId must be hyperdensity_action_slate_v1")
    if action_slate_example["actionSlateVersion"] != "v1":
        fail("action-slate-reference.json actionSlateVersion must be v1")
    if action_slate_example["actionSlateMode"] != "prevalidated_recommendation_only":
        fail("action-slate-reference.json actionSlateMode must be prevalidated_recommendation_only")
    if action_slate_example["releaseTrack"] != "technical_preview":
        fail("action-slate-reference.json releaseTrack must be technical_preview")
    if action_slate_example["operatorControlled"] is not True:
        fail("action-slate-reference.json operatorControlled must be true")
    if action_slate_example["autonomousApplyAllowed"] is not False:
        fail("action-slate-reference.json autonomousApplyAllowed must be false")
    if action_slate_example["enforcementMode"] != "disabled":
        fail("action-slate-reference.json enforcementMode must be disabled")
    if action_slate_example["productionMutationAllowed"] is not False:
        fail("action-slate-reference.json productionMutationAllowed must be false")
    if action_slate_example["evidenceScope"] != "evidence_namespace_only":
        fail("action-slate-reference.json evidenceScope must be evidence_namespace_only")
    if not isinstance(action_slate_example["actions"], list) or not action_slate_example["actions"]:
        fail("action-slate-reference.json actions must be non-empty")
    if not isinstance(action_slate_example["safetyGates"], list) or len(action_slate_example["safetyGates"]) < 10:
        fail("action-slate-reference.json safetyGates must include all required safety gates")
    gate_ids = {entry.get("gateId") for entry in action_slate_example["safetyGates"] if isinstance(entry, dict)}
    for gate_id in (
        "autonomous_apply_disabled",
        "enforcement_disabled",
        "production_mutation_disabled",
        "raw_runtime_controls_not_exposed",
        "raw_resource_creation_not_allowed",
        "dry_run_required_or_ready",
        "rollback_required_or_ready",
        "technical_preview_boundary_active",
        "support_claims_evidence_backed_only",
        "windows_out_of_scope",
    ):
        if gate_id not in gate_ids:
            fail(f"action-slate-reference.json safetyGates must include {gate_id}")
    for action in action_slate_example["actions"]:
        if not isinstance(action, dict):
            fail("action-slate-reference.json actions must contain objects")
        if action.get("recommendationMode") != "recommendation_only":
            fail("action-slate-reference.json actions must be recommendation_only")
        if action.get("autonomousApplyAllowed") is not False:
            fail("action-slate-reference.json action autonomousApplyAllowed must be false")
        if action.get("productionMutationAllowed") is not False:
            fail("action-slate-reference.json action productionMutationAllowed must be false")
    indexes = action_slate_example["indexes"]
    if not isinstance(indexes, dict):
        fail("action-slate-reference.json indexes must be an object")
    if indexes.get("pairingStrategy") != "top_k_pairing":
        fail("action-slate-reference.json indexes pairingStrategy must be top_k_pairing")
    if indexes.get("fullFleetScanRequired") is not False:
        fail("action-slate-reference.json indexes fullFleetScanRequired must be false")
    if indexes.get("batchingMode") != "rate_limited_operator_controlled":
        fail("action-slate-reference.json indexes batchingMode must be rate_limited_operator_controlled")

    guarded_auto_example = example_by_name.get("guarded-auto-sandbox-reference.json")
    if guarded_auto_example is None:
        fail("examples/guarded-auto-sandbox-reference.json is missing")
    required_guarded_auto_fields = [
        "guardedAutoSandboxId",
        "guardedAutoSandboxVersion",
        "releaseTrack",
        "sandboxMode",
        "sandboxState",
        "allowedNamespace",
        "productionMutationAllowed",
        "enforcementMode",
        "autonomousProductionApplyAllowed",
        "sandboxAutonomousApplyAllowed",
        "operatorKillSwitchRequired",
        "operatorKillSwitchState",
        "actionSlateId",
        "policyPackId",
        "supportMatrixId",
        "evidenceBundleId",
        "liveResourceAuthorityId",
        "maxActionsPerHour",
        "maxConcurrentActions",
        "maxCpuMove",
        "maxMemoryMove",
        "dryRunRequired",
        "rollbackRequired",
        "auditRequired",
        "verificationRequired",
        "eligibleActionCount",
        "executableActionCount",
        "blockedActionCount",
        "candidateActions",
        "safetyGates",
        "blastRadiusBudget",
        "auditTrail",
        "blockers",
        "limitations",
        "claimBoundaries",
    ]
    for field in required_guarded_auto_fields:
        if field not in guarded_auto_example:
            fail(f"guarded-auto-sandbox-reference.json missing required field '{field}'")
    if guarded_auto_example["guardedAutoSandboxId"] != "hyperdensity_guarded_auto_sandbox_v1":
        fail("guarded-auto-sandbox-reference.json guardedAutoSandboxId must be hyperdensity_guarded_auto_sandbox_v1")
    if guarded_auto_example["guardedAutoSandboxVersion"] != "v1":
        fail("guarded-auto-sandbox-reference.json guardedAutoSandboxVersion must be v1")
    if guarded_auto_example["releaseTrack"] != "technical_preview":
        fail("guarded-auto-sandbox-reference.json releaseTrack must be technical_preview")
    if guarded_auto_example["sandboxMode"] != "guarded_auto_evidence_namespace_only":
        fail("guarded-auto-sandbox-reference.json sandboxMode must be guarded_auto_evidence_namespace_only")
    if guarded_auto_example["allowedNamespace"] != "karl-hyperdensity-evidence":
        fail("guarded-auto-sandbox-reference.json allowedNamespace must be karl-hyperdensity-evidence")
    if guarded_auto_example["productionMutationAllowed"] is not False:
        fail("guarded-auto-sandbox-reference.json productionMutationAllowed must be false")
    if guarded_auto_example["enforcementMode"] != "disabled":
        fail("guarded-auto-sandbox-reference.json enforcementMode must be disabled")
    if guarded_auto_example["autonomousProductionApplyAllowed"] is not False:
        fail("guarded-auto-sandbox-reference.json autonomousProductionApplyAllowed must be false")
    if guarded_auto_example["operatorKillSwitchRequired"] is not True:
        fail("guarded-auto-sandbox-reference.json operatorKillSwitchRequired must be true")
    if guarded_auto_example["dryRunRequired"] is not True:
        fail("guarded-auto-sandbox-reference.json dryRunRequired must be true")
    if guarded_auto_example["rollbackRequired"] is not True:
        fail("guarded-auto-sandbox-reference.json rollbackRequired must be true")
    if guarded_auto_example["auditRequired"] is not True:
        fail("guarded-auto-sandbox-reference.json auditRequired must be true")
    if guarded_auto_example["verificationRequired"] is not True:
        fail("guarded-auto-sandbox-reference.json verificationRequired must be true")
    if not isinstance(guarded_auto_example["candidateActions"], list):
        fail("guarded-auto-sandbox-reference.json candidateActions must be an array")
    for candidate in guarded_auto_example["candidateActions"]:
        if not isinstance(candidate, dict):
            fail("guarded-auto-sandbox-reference.json candidateActions must contain objects")
        if candidate.get("productionMutationAllowed") is not False:
            fail("guarded-auto-sandbox-reference.json candidate productionMutationAllowed must be false")
        if candidate.get("namespace") not in ("karl-hyperdensity-evidence", None):
            fail("guarded-auto-sandbox-reference.json candidate namespace must stay in evidence namespace")
    gate_ids = {entry.get("gateId") for entry in guarded_auto_example["safetyGates"] if isinstance(entry, dict)}
    for gate_id in (
        "evidence_namespace_only",
        "production_mutation_disabled",
        "enforcement_disabled",
        "autonomous_production_apply_disabled",
        "operator_kill_switch_available",
        "action_slate_ready_required",
        "dry_run_ready_required",
        "rollback_ready_required",
        "low_risk_required",
        "support_boundary_required",
        "policy_consistency_required",
        "blast_radius_budget_required",
        "audit_required",
        "verification_required",
        "raw_runtime_controls_not_exposed",
        "raw_resource_creation_not_allowed",
        "windows_out_of_scope",
    ):
        if gate_id not in gate_ids:
            fail(f"guarded-auto-sandbox-reference.json safetyGates must include {gate_id}")
    claim_text = "\n".join(guarded_auto_example["claimBoundaries"]).lower()
    for phrase in (
        "guarded auto is evidence-namespace only.",
        "production autonomous apply is disabled.",
        "enforcement is disabled.",
        "no production mutation.",
        "dry-run is required.",
        "rollback proof is required.",
        "low risk is required.",
        "operator kill switch is required.",
        "no raw runtime controls are exposed.",
        "no raw resource creation.",
        "technical preview boundary active.",
        "not ga.",
        "not ha.",
        "not windows.",
        "not generic vm ram template mutation.",
    ):
        if phrase not in claim_text:
            fail(f"guarded-auto-sandbox-reference.json claimBoundaries missing phrase: {phrase}")

    auto_rollback_example = example_by_name.get("auto-rollback-controller-reference.json")
    if auto_rollback_example is None:
        fail("examples/auto-rollback-controller-reference.json is missing")
    required_auto_rollback_fields = [
        "autoRollbackControllerId",
        "autoRollbackControllerVersion",
        "releaseTrack",
        "controllerMode",
        "controllerState",
        "allowedNamespace",
        "productionRollbackAllowed",
        "productionMutationAllowed",
        "enforcementMode",
        "autonomousProductionApplyAllowed",
        "autonomousSandboxRollbackAllowed",
        "operatorKillSwitchRequired",
        "operatorKillSwitchState",
        "actionSlateId",
        "guardedAutoSandboxId",
        "policyPackId",
        "supportMatrixId",
        "evidenceBundleId",
        "liveResourceAuthorityId",
        "rollbackRequired",
        "verificationRequired",
        "auditRequired",
        "rollbackTriggerCount",
        "readyRollbackPlanCount",
        "blockedRollbackPlanCount",
        "rollbackPlans",
        "rollbackTriggers",
        "safetyGates",
        "blockers",
        "limitations",
        "claimBoundaries",
        "auditTrail",
    ]
    for field in required_auto_rollback_fields:
        if field not in auto_rollback_example:
            fail(f"auto-rollback-controller-reference.json missing required field '{field}'")
    if auto_rollback_example["autoRollbackControllerId"] != "hyperdensity_auto_rollback_controller_v1":
        fail("auto-rollback-controller-reference.json autoRollbackControllerId must be hyperdensity_auto_rollback_controller_v1")
    if auto_rollback_example["autoRollbackControllerVersion"] != "v1":
        fail("auto-rollback-controller-reference.json autoRollbackControllerVersion must be v1")
    if auto_rollback_example["releaseTrack"] != "technical_preview":
        fail("auto-rollback-controller-reference.json releaseTrack must be technical_preview")
    if auto_rollback_example["controllerMode"] != "guarded_sandbox_rollback_readiness":
        fail("auto-rollback-controller-reference.json controllerMode must be guarded_sandbox_rollback_readiness")
    if auto_rollback_example["allowedNamespace"] != "karl-hyperdensity-evidence":
        fail("auto-rollback-controller-reference.json allowedNamespace must be karl-hyperdensity-evidence")
    if auto_rollback_example["productionRollbackAllowed"] is not False:
        fail("auto-rollback-controller-reference.json productionRollbackAllowed must be false")
    if auto_rollback_example["productionMutationAllowed"] is not False:
        fail("auto-rollback-controller-reference.json productionMutationAllowed must be false")
    if auto_rollback_example["enforcementMode"] != "disabled":
        fail("auto-rollback-controller-reference.json enforcementMode must be disabled")
    if auto_rollback_example["autonomousProductionApplyAllowed"] is not False:
        fail("auto-rollback-controller-reference.json autonomousProductionApplyAllowed must be false")
    if auto_rollback_example["operatorKillSwitchRequired"] is not True:
        fail("auto-rollback-controller-reference.json operatorKillSwitchRequired must be true")
    if auto_rollback_example["rollbackRequired"] is not True:
        fail("auto-rollback-controller-reference.json rollbackRequired must be true")
    if auto_rollback_example["verificationRequired"] is not True:
        fail("auto-rollback-controller-reference.json verificationRequired must be true")
    if auto_rollback_example["auditRequired"] is not True:
        fail("auto-rollback-controller-reference.json auditRequired must be true")
    if not isinstance(auto_rollback_example["rollbackPlans"], list):
        fail("auto-rollback-controller-reference.json rollbackPlans must be an array")
    for plan in auto_rollback_example["rollbackPlans"]:
        if not isinstance(plan, dict):
            fail("auto-rollback-controller-reference.json rollbackPlans must contain objects")
        if plan.get("productionRollbackAllowed") is not False:
            fail("auto-rollback-controller-reference.json rollback plan productionRollbackAllowed must be false")
        if plan.get("rollbackEligibility") == "eligible" and plan.get("namespace") != "karl-hyperdensity-evidence":
            fail("auto-rollback-controller-reference.json eligible rollback plans must stay in evidence namespace")
    required_trigger_ids = {
        "verification_failed",
        "runtime_not_converged",
        "cgroup_not_converged",
        "qga_libvirt_qmp_not_converged",
        "warning_event_detected",
        "restart_count_changed",
        "receiver_not_improved",
        "donor_became_pressured",
        "slo_degraded",
        "operator_kill_switch_triggered",
        "action_expired_before_verify",
    }
    trigger_ids = {entry.get("triggerId") for entry in auto_rollback_example["rollbackTriggers"] if isinstance(entry, dict)}
    missing_trigger_ids = required_trigger_ids - trigger_ids
    if missing_trigger_ids:
        fail(f"auto-rollback-controller-reference.json rollbackTriggers missing: {sorted(missing_trigger_ids)}")
    gate_ids = {entry.get("gateId") for entry in auto_rollback_example["safetyGates"] if isinstance(entry, dict)}
    for gate_id in (
        "evidence_namespace_only",
        "production_rollback_disabled",
        "production_mutation_disabled",
        "enforcement_disabled",
        "autonomous_production_apply_disabled",
        "operator_kill_switch_available",
        "rollback_source_required",
        "verification_required",
        "audit_required",
        "low_risk_required",
        "support_boundary_required",
        "policy_consistency_required",
        "blast_radius_budget_required",
        "raw_runtime_controls_not_exposed",
        "raw_resource_creation_not_allowed",
        "windows_out_of_scope",
    ):
        if gate_id not in gate_ids:
            fail(f"auto-rollback-controller-reference.json safetyGates must include {gate_id}")
    auto_claim_text = "\n".join(auto_rollback_example["claimBoundaries"]).lower()
    for phrase in (
        "automatic rollback is evidence-namespace only.",
        "production rollback is disabled.",
        "production autonomous apply is disabled.",
        "enforcement is disabled.",
        "rollback source is required.",
        "verification is required.",
        "operator kill switch is required.",
        "no raw runtime controls are exposed.",
        "no raw resource creation.",
        "technical preview boundary active.",
        "not ga.",
        "not ha.",
        "not windows.",
        "not generic vm ram template mutation.",
    ):
        if phrase not in auto_claim_text:
            fail(f"auto-rollback-controller-reference.json claimBoundaries missing phrase: {phrase}")

    blast_radius_example = example_by_name.get("blast-radius-policy-reference.json")
    if blast_radius_example is None:
        fail("examples/blast-radius-policy-reference.json is missing")
    required_blast_radius_fields = [
        "blastRadiusPolicyId",
        "blastRadiusPolicyVersion",
        "releaseTrack",
        "policyMode",
        "policyState",
        "productionAutonomousApplyAllowed",
        "productionMutationAllowed",
        "enforcementMode",
        "evidenceNamespace",
        "guardedAutoSandboxId",
        "actionSlateId",
        "autoRollbackControllerId",
        "policyPackId",
        "supportMatrixId",
        "evidenceBundleId",
        "liveResourceAuthorityId",
        "operatorKillSwitchRequired",
        "operatorKillSwitchState",
        "globalBudget",
        "namespaceBudgets",
        "resourceBudgets",
        "concurrencyLimits",
        "rateLimits",
        "freezeConditions",
        "escalationRules",
        "stopConditions",
        "safetyGates",
        "blockers",
        "limitations",
        "claimBoundaries",
        "auditTrail",
    ]
    for field in required_blast_radius_fields:
        if field not in blast_radius_example:
            fail(f"blast-radius-policy-reference.json missing required field '{field}'")
    if blast_radius_example["blastRadiusPolicyId"] != "hyperdensity_blast_radius_policy_v1":
        fail("blast-radius-policy-reference.json blastRadiusPolicyId must be hyperdensity_blast_radius_policy_v1")
    if blast_radius_example["blastRadiusPolicyVersion"] != "v1":
        fail("blast-radius-policy-reference.json blastRadiusPolicyVersion must be v1")
    if blast_radius_example["releaseTrack"] != "technical_preview":
        fail("blast-radius-policy-reference.json releaseTrack must be technical_preview")
    if blast_radius_example["policyMode"] != "guarded_auto_safety_budget":
        fail("blast-radius-policy-reference.json policyMode must be guarded_auto_safety_budget")
    if blast_radius_example["productionAutonomousApplyAllowed"] is not False:
        fail("blast-radius-policy-reference.json productionAutonomousApplyAllowed must be false")
    if blast_radius_example["productionMutationAllowed"] is not False:
        fail("blast-radius-policy-reference.json productionMutationAllowed must be false")
    if blast_radius_example["enforcementMode"] != "disabled":
        fail("blast-radius-policy-reference.json enforcementMode must be disabled")
    if blast_radius_example["evidenceNamespace"] != "karl-hyperdensity-evidence":
        fail("blast-radius-policy-reference.json evidenceNamespace must be karl-hyperdensity-evidence")
    if blast_radius_example["operatorKillSwitchRequired"] is not True:
        fail("blast-radius-policy-reference.json operatorKillSwitchRequired must be true")

    global_budget = blast_radius_example["globalBudget"]
    required_global_budget_fields = [
        "maxConcurrentActions",
        "maxActionsPerHour",
        "maxActionsPerNamespacePerHour",
        "maxCpuMovePerAction",
        "maxMemoryMovePerAction",
        "maxCpuMovePerHour",
        "maxMemoryMovePerHour",
        "maxFleetPercentTouchedPerHour",
        "maxNamespacePercentTouchedPerHour",
        "allowedRiskLevels",
        "allowedActionStates",
        "requiredDryRunState",
        "requiredRollbackState",
        "requiredVerificationState",
        "requiredAuditState",
    ]
    if not isinstance(global_budget, dict):
        fail("blast-radius-policy-reference.json globalBudget must be an object")
    for field in required_global_budget_fields:
        if field not in global_budget:
            fail(f"blast-radius-policy-reference.json globalBudget missing required field '{field}'")

    if not isinstance(blast_radius_example["namespaceBudgets"], list) or not blast_radius_example["namespaceBudgets"]:
        fail("blast-radius-policy-reference.json namespaceBudgets must be a non-empty array")
    has_production_budget = False
    for budget in blast_radius_example["namespaceBudgets"]:
        if not isinstance(budget, dict):
            fail("blast-radius-policy-reference.json namespaceBudgets entries must be objects")
        for field in (
            "namespace",
            "namespaceClass",
            "autoMode",
            "maxConcurrentActions",
            "maxActionsPerHour",
            "maxCpuMovePerHour",
            "maxMemoryMovePerHour",
            "allowedRiskLevels",
            "productionMutationAllowed",
            "blocker",
        ):
            if field not in budget:
                fail(f"blast-radius-policy-reference.json namespace budget missing field '{field}'")
        if budget["namespaceClass"] == "production":
            has_production_budget = True
            if budget["productionMutationAllowed"] is not False:
                fail("blast-radius-policy-reference.json production namespace budget must keep productionMutationAllowed=false")
            if budget["autoMode"] not in ("disabled", "observe_only", "recommendation_only"):
                fail("blast-radius-policy-reference.json production namespace autoMode must stay non-autonomous")
    if not has_production_budget:
        fail("blast-radius-policy-reference.json must include production namespace budget")

    required_freeze_ids = {
        "warning_event_detected",
        "rollback_failure",
        "verification_failure",
        "donor_pressure_increased",
        "receiver_not_improved",
        "policy_inconsistency",
        "support_boundary_missing",
        "unknown_risk_detected",
        "kill_switch_triggered",
        "audit_gap_detected",
    }
    required_stop_ids = {
        "max_budget_exceeded",
        "high_risk_action_detected",
        "blocked_action_detected",
        "expired_action_detected",
        "unsupported_shell_detected",
        "windows_lane_detected",
        "raw_runtime_control_requested",
        "production_scope_requested",
    }
    required_escalation_ids = {
        "medium_risk_requires_operator_review",
        "high_risk_blocks_auto",
        "unknown_risk_blocks_auto",
        "production_scope_blocks_auto",
        "rollback_missing_blocks_auto",
        "dry_run_missing_blocks_auto",
    }
    freeze_ids = {entry.get("conditionId") for entry in blast_radius_example["freezeConditions"] if isinstance(entry, dict)}
    stop_ids = {entry.get("conditionId") for entry in blast_radius_example["stopConditions"] if isinstance(entry, dict)}
    escalation_ids = {entry.get("ruleId") for entry in blast_radius_example["escalationRules"] if isinstance(entry, dict)}
    missing_freeze = required_freeze_ids - freeze_ids
    missing_stop = required_stop_ids - stop_ids
    missing_escalation = required_escalation_ids - escalation_ids
    if missing_freeze:
        fail(f"blast-radius-policy-reference.json freezeConditions missing: {sorted(missing_freeze)}")
    if missing_stop:
        fail(f"blast-radius-policy-reference.json stopConditions missing: {sorted(missing_stop)}")
    if missing_escalation:
        fail(f"blast-radius-policy-reference.json escalationRules missing: {sorted(missing_escalation)}")

    blast_claim_text = "\n".join(blast_radius_example["claimBoundaries"]).lower()
    for phrase in (
        "blast radius limits are enforced before any guarded auto action.",
        "production autonomous apply is disabled.",
        "enforcement is disabled.",
        "no production mutation.",
        "dry-run is required.",
        "rollback proof is required.",
        "verification is required.",
        "audit is required.",
        "operator kill switch is required.",
        "high and unknown risk block auto.",
        "technical preview boundary active.",
    ):
        if phrase not in blast_claim_text:
            fail(f"blast-radius-policy-reference.json claimBoundaries missing phrase: {phrase}")

    rc_gate_example = example_by_name.get("technical-preview-release-candidate-gate-reference.json")
    if rc_gate_example is None:
        fail("examples/technical-preview-release-candidate-gate-reference.json is missing")
    required_rc_gate_fields = [
        "releaseCandidateGateId",
        "releaseTrack",
        "gateMode",
        "targetReadinessState",
        "supportMatrixId",
        "evidenceBundleId",
        "liveResourceAuthorityId",
        "documentationPackId",
        "policyPackId",
        "profilePackId",
        "enforcementMode",
        "autonomousApplyAllowed",
        "productionMutationAllowed",
        "windowsSupportState",
        "finalGateChecks",
        "blockerChecks",
        "releaseDecision",
        "releaseBlocker",
        "requiredFollowUps",
        "summary",
    ]
    for field in required_rc_gate_fields:
        if field not in rc_gate_example:
            fail(f"technical-preview-release-candidate-gate-reference.json missing required field '{field}'")
    if rc_gate_example["releaseCandidateGateId"] != "hyperdensity_technical_preview_release_candidate_gate_v1":
        fail("technical-preview-release-candidate-gate-reference.json releaseCandidateGateId mismatch")
    if rc_gate_example["releaseTrack"] != "technical_preview":
        fail("technical-preview-release-candidate-gate-reference.json releaseTrack must be technical_preview")
    if rc_gate_example["gateMode"] != "go_no_go_validation":
        fail("technical-preview-release-candidate-gate-reference.json gateMode must be go_no_go_validation")
    if rc_gate_example["targetReadinessState"] != "ready_for_private_technical_preview":
        fail("technical-preview-release-candidate-gate-reference.json targetReadinessState must be ready_for_private_technical_preview")
    if rc_gate_example["enforcementMode"] != "disabled":
        fail("technical-preview-release-candidate-gate-reference.json enforcementMode must be disabled")
    if rc_gate_example["autonomousApplyAllowed"] is not False:
        fail("technical-preview-release-candidate-gate-reference.json autonomousApplyAllowed must be false")
    if rc_gate_example["productionMutationAllowed"] is not False:
        fail("technical-preview-release-candidate-gate-reference.json productionMutationAllowed must be false")
    if rc_gate_example["windowsSupportState"] != "out_of_scope_frozen":
        fail("technical-preview-release-candidate-gate-reference.json windowsSupportState must be out_of_scope_frozen")
    if not isinstance(rc_gate_example["finalGateChecks"], list) or not rc_gate_example["finalGateChecks"]:
        fail("technical-preview-release-candidate-gate-reference.json finalGateChecks must be non-empty")
    if not isinstance(rc_gate_example["blockerChecks"], list) or not rc_gate_example["blockerChecks"]:
        fail("technical-preview-release-candidate-gate-reference.json blockerChecks must be non-empty")
    if rc_gate_example["releaseDecision"] not in (
        "ready_for_private_technical_preview",
        "ready_for_founder_investor_demo_only",
        "ready_for_internal_evidence_release_only",
        "blocked",
    ):
        fail("technical-preview-release-candidate-gate-reference.json releaseDecision has invalid value")

    doc_paths = [
        repo_root / "docs" / "runbooks" / "operator-runbook-v1.md",
        repo_root / "docs" / "releases" / "technical-preview-release-notes-v1.md",
        repo_root / "docs" / "releases" / "technical-preview-readiness-gate-v1.md",
        repo_root / "docs" / "releases" / "technical-preview-release-candidate-gate-v1.md",
        repo_root / "docs" / "demos" / "technical-preview-demo-guide-v1.md",
        repo_root / "docs" / "releases" / "technical-preview-documentation-pack-v1.md",
    ]
    for path in doc_paths:
        if not path.exists():
            fail(f"required documentation file missing: {path}")

    docs = {path.name: load_text(path) for path in doc_paths}
    docs_lower = {name: body.lower() for name, body in docs.items()}

    for name, body in docs_lower.items():
        if "technical preview" not in body:
            fail(f"{name} must explicitly mention Technical Preview")

    runbook = docs_lower["operator-runbook-v1.md"]
    if "evidence namespace" not in runbook:
        fail("operator-runbook-v1.md must mention evidence namespace safety")
    if "operator-controlled" not in runbook:
        fail("operator-runbook-v1.md must mention operator-controlled safety")
    if "warming_up" not in runbook or "not ready" not in runbook:
        fail("operator-runbook-v1.md must state warming_up is not ready")
    if "partial" not in runbook or "not ready" not in runbook:
        fail("operator-runbook-v1.md must state partial is not ready")
    if "blocked" not in runbook or "not ready" not in runbook:
        fail("operator-runbook-v1.md must state blocked is not ready")

    release_notes = docs_lower["technical-preview-release-notes-v1.md"]
    required_release_phrases = [
        "not ga",
        "no enforcement enabled",
        "no autonomous production movement",
        "no production workload mutation",
        "no windows support",
        "no generic kubevirt template memory mutation claim",
    ]
    for phrase in required_release_phrases:
        if phrase not in release_notes:
            fail(f"technical-preview-release-notes-v1.md missing required non-claim phrase: {phrase}")

    pack_index = docs_lower["technical-preview-documentation-pack-v1.md"]
    required_pack_refs = [
        "documentationpackid=hyperdensity_technical_preview_documentation_pack_v1",
        "documentationpackversion=v1",
        "releasetrack=technical_preview",
        "supportmatrixid=hyperdensity_release_support_matrix_v1",
        "evidencebundleid=hyperdensity_evidence_bundle_demo_scenario_pack_v1",
        "liveresourceauthorityid=hyperdensity_live_resource_authority_v1",
        "policypackid=hyperdensity_policy_pack_v1",
        "profilepackid=hyperdensity_shell_claim_templates_profile_pack_v1",
    ]
    for ref in required_pack_refs:
        if ref not in pack_index:
            fail(f"technical-preview-documentation-pack-v1.md missing required reference: {ref}")

    rc_gate_doc = docs_lower["technical-preview-release-candidate-gate-v1.md"]
    required_rc_doc_refs = [
        "releasecandidategateid=hyperdensity_technical_preview_release_candidate_gate_v1",
        "releasetrack=technical_preview",
        "gatemode=go_no_go_validation",
        "targetreadinessstate=ready_for_private_technical_preview",
        "supportmatrixid=hyperdensity_release_support_matrix_v1",
        "evidencebundleid=hyperdensity_evidence_bundle_demo_scenario_pack_v1",
        "liveresourceauthorityid=hyperdensity_live_resource_authority_v1",
        "documentationpackid=hyperdensity_technical_preview_documentation_pack_v1",
        "policypackid=hyperdensity_policy_pack_v1",
        "profilepackid=hyperdensity_shell_claim_templates_profile_pack_v1",
        "enforcementmode=disabled",
        "autonomousapplyallowed=false",
        "productionmutationallowed=false",
        "windowssupportstate=out_of_scope_frozen",
    ]
    for ref in required_rc_doc_refs:
        if ref not in rc_gate_doc:
            fail(f"technical-preview-release-candidate-gate-v1.md missing required reference: {ref}")

    # --- Sprint 1: production kernel boundary v1 ---
    sprint1_schemas = [
        "shell-passport-v1.schema.json",
        "runtime-mutation-result-v1.schema.json",
        "resource-lease-v1.schema.json",
        "hyperdensity-claim-policy-v2.schema.json",
        "production-kernel-boundary-v1.schema.json",
    ]
    for name in sprint1_schemas:
        if name not in {p.name for p in schema_paths}:
            fail(f"missing Sprint 1 schema: {name}")

    claim_policy_v2 = example_by_name.get("hyperdensity-claim-policy-v2-reference.json")
    if claim_policy_v2 is None:
        fail("examples/hyperdensity-claim-policy-v2-reference.json is missing")
    if claim_policy_v2.get("claimPolicyId") != "hyperdensity_claim_policy_v2":
        fail("claim policy v2 reference claimPolicyId invalid")
    for field in (
        "guaranteedSavingsAllowed",
        "universalPerformanceImprovementAllowed",
        "logicalVcpuHotplugClaimAllowed",
        "windowsTotalRamHotplugClaimAllowed",
        "ramAboveOriginalClaimAllowed",
        "productionAutonomousApplyAllowed",
        "syntheticFleetProductionClaimAllowed",
    ):
        if claim_policy_v2.get(field) is not False:
            fail(f"hyperdensity-claim-policy-v2-reference.json {field} must be false")

    boundary_ref = example_by_name.get("production-kernel-boundary-reference.json")
    if boundary_ref is None:
        fail("examples/production-kernel-boundary-reference.json is missing")
    if boundary_ref.get("boundaryId") != "hyperdensity_production_kernel_boundary_v1":
        fail("production-kernel-boundary-reference.json boundaryId invalid")
    safety = boundary_ref.get("safetyInvariants") or {}
    for field in (
        "productionAutonomousApplyAllowed",
        "guaranteedSavingsAllowed",
        "universalPerformanceImprovementAllowed",
        "dashboardMutationSourceOfTruth",
        "inventoryHyperdensityEngine",
    ):
        if safety.get(field) is not False:
            fail(f"production-kernel-boundary safetyInvariants.{field} must be false")

    mutation_ref = example_by_name.get("runtime-mutation-result-reference.json")
    if mutation_ref is None:
        fail("examples/runtime-mutation-result-reference.json is missing")
    if mutation_ref.get("actuator") != "FluidVirt":
        fail("runtime mutation result actuator must be FluidVirt")

    lease_ref = example_by_name.get("resource-lease-reference.json")
    if lease_ref is None:
        fail("examples/resource-lease-reference.json is missing")
    if lease_ref.get("autoApplyAllowed") is not False:
        fail("resource lease autoApplyAllowed must be false")
    if lease_ref.get("productionMutationAllowed") is not False:
        fail("resource lease productionMutationAllowed must be false")

    forbidden_positive_claims = [
        "guaranteed savings active",
        "universal performance improvement",
        "production autonomous apply",
        "windows total ram hotplug supported",
        "logical vcpu hotplug supported",
        "1000 production workloads proven",
        "dashboard is source of truth",
        "inventory hyperdensity engine",
    ]
    sprint1_examples = [
        "hyperdensity-claim-policy-v2-reference.json",
        "production-kernel-boundary-reference.json",
        "runtime-mutation-result-reference.json",
        "resource-lease-reference.json",
        "shell-passport-reference.json",
    ]
    positive_claim_keys = {
        "allowedPhrases",
        "conditionalPhrases",
        "claimBoundary",
        "claimBoundaries",
        "allowedResponsibilities",
        "limitation",
        "limitations",
    }

    def collect_positive_strings(obj, out: list[str]) -> None:
        if isinstance(obj, dict):
            for k, v in obj.items():
                if k in ("forbiddenPhrases", "forbiddenResponsibilities", "forbiddenActions", "blockerCodes", "blockers", "remediationCodes", "remediations", "unsupportedFamilies", "excludedShellKinds"):
                    continue
                if k in positive_claim_keys or (k.endswith("Phrases") and not k.startswith("forbidden")):
                    if isinstance(v, str):
                        out.append(v)
                    elif isinstance(v, list):
                        for item in v:
                            if isinstance(item, str):
                                out.append(item)
                else:
                    collect_positive_strings(v, out)
        elif isinstance(obj, list):
            for item in obj:
                collect_positive_strings(item, out)

    for name in sprint1_examples:
        ex = example_by_name.get(name)
        if ex is None:
            fail(f"missing Sprint 1 example: {name}")
        positives: list[str] = []
        collect_positive_strings(ex, positives)
        merged = "\n".join(positives).lower()
        for phrase in forbidden_positive_claims:
            if phrase in merged:
                fail(f"{name} contains forbidden positive claim phrase in allowed copy: {phrase}")

    # --- Sprint 2: shell passport factory v1 ---
    sprint2_schemas = [
        "shell-passport-factory-v1.schema.json",
        "shell-registry-v1.schema.json",
        "shell-enrollment-result-v1.schema.json",
        "shell-capability-evidence-v1.schema.json",
    ]
    for name in sprint2_schemas:
        if name not in {p.name for p in schema_paths}:
            fail(f"missing Sprint 2 schema: {name}")

    factory_ref = example_by_name.get("shell-passport-factory-reference.json")
    if factory_ref is None:
        fail("examples/shell-passport-factory-reference.json is missing")
    if factory_ref.get("factoryId") != "hyperdensity_shell_passport_factory_v1":
        fail("shell passport factory factoryId invalid")
    if factory_ref.get("productionMutationAllowed") is not False:
        fail("shell passport factory productionMutationAllowed must be false")
    if factory_ref.get("autoApplyAllowed") is not False:
        fail("shell passport factory autoApplyAllowed must be false")

    registry_surface = example_by_name.get("shell-passport-registry-reference.json")
    if registry_surface is None:
        fail("examples/shell-passport-registry-reference.json is missing")
    if registry_surface.get("registryId") != "hyperdensity_shell_registry_v1":
        fail("shell registry surface registryId invalid")

    enrollment_ref = example_by_name.get("shell-enrollment-result-reference.json")
    if enrollment_ref is None:
        fail("examples/shell-enrollment-result-reference.json is missing")
    if enrollment_ref.get("productionMutationAllowed") is not False:
        fail("enrollment productionMutationAllowed must be false")

    capability_ref = example_by_name.get("shell-capability-evidence-reference.json")
    if capability_ref is None:
        fail("examples/shell-capability-evidence-reference.json is missing")
    if capability_ref.get("source") != "FluidVirt":
        fail("capability evidence source must be FluidVirt")

    sprint2_examples = [
        "shell-passport-factory-reference.json",
        "shell-passport-registry-reference.json",
        "shell-enrollment-result-reference.json",
        "shell-capability-evidence-reference.json",
    ]
    for name in sprint2_examples:
        ex = example_by_name.get(name)
        if ex is None:
            fail(f"missing Sprint 2 example: {name}")
        positives2: list[str] = []
        collect_positive_strings(ex, positives2)
        merged2 = "\n".join(positives2).lower()
        for phrase in forbidden_positive_claims:
            if phrase in merged2:
                fail(f"{name} contains forbidden positive claim in allowed copy: {phrase}")

    # --- Sprint 3: resource lease + action slate readiness v1 ---
    sprint3_schemas = [
        "donor-index-v1.schema.json",
        "receiver-index-v1.schema.json",
        "resource-lease-candidate-v1.schema.json",
        "action-slate-readiness-v1.schema.json",
        "action-slate-entry-v1.schema.json",
        "action-dryrun-readiness-v1.schema.json",
        "rollback-readiness-v1.schema.json",
        "slo-precheck-v1.schema.json",
        "risk-assessment-v1.schema.json",
    ]
    for name in sprint3_schemas:
        if name not in {p.name for p in schema_paths}:
            fail(f"missing Sprint 3 schema: {name}")

    slate_readiness = example_by_name.get("action-slate-readiness-reference.json")
    if slate_readiness is None:
        fail("examples/action-slate-readiness-reference.json is missing")
    if slate_readiness.get("milestone") != "hyperdensity_resource_lease_action_slate_readiness_v1":
        fail("action slate readiness milestone invalid")
    if slate_readiness.get("noFullNxNPairing") is not True:
        fail("action slate readiness noFullNxNPairing must be true")
    inv = slate_readiness.get("safetyInvariants") or {}
    if inv.get("autoApplyAllowed") is not False:
        fail("action slate safetyInvariants.autoApplyAllowed must be false")
    if inv.get("productionMutationAllowed") is not False:
        fail("action slate safetyInvariants.productionMutationAllowed must be false")

    donor_idx = example_by_name.get("donor-index-reference.json")
    if donor_idx is None:
        fail("examples/donor-index-reference.json is missing")
    if donor_idx.get("indexId") != "hyperdensity_donor_index_v1":
        fail("donor index indexId invalid")

    receiver_idx = example_by_name.get("receiver-index-reference.json")
    if receiver_idx is None:
        fail("examples/receiver-index-reference.json is missing")
    if receiver_idx.get("indexId") != "hyperdensity_receiver_index_v1":
        fail("receiver index indexId invalid")

    dryrun_ref = example_by_name.get("action-dryrun-readiness-reference.json")
    if dryrun_ref is None:
        fail("examples/action-dryrun-readiness-reference.json is missing")
    if dryrun_ref.get("source") != "FluidVirt":
        fail("dry-run readiness source must be FluidVirt")
    if dryrun_ref.get("mutationExecuted") is not False:
        fail("dry-run readiness mutationExecuted must be false")

    sprint3_examples = [
        "donor-index-reference.json",
        "receiver-index-reference.json",
        "resource-lease-candidate-reference.json",
        "action-slate-readiness-reference.json",
        "action-slate-entry-reference.json",
        "action-dryrun-readiness-reference.json",
        "rollback-readiness-reference.json",
        "slo-precheck-reference.json",
        "risk-assessment-reference.json",
    ]
    for name in sprint3_examples:
        ex = example_by_name.get(name)
        if ex is None:
            fail(f"missing Sprint 3 example: {name}")
        positives3: list[str] = []
        collect_positive_strings(ex, positives3)
        merged3 = "\n".join(positives3).lower()
        for phrase in forbidden_positive_claims:
            if phrase in merged3:
                fail(f"{name} contains forbidden positive claim in allowed copy: {phrase}")

    # --- Sprint 4: operator-controlled apply gate v1 ---
    sprint4_schemas = [
        "operator-apply-gate-v1.schema.json",
        "operator-approval-record-v1.schema.json",
        "apply-request-v1.schema.json",
        "fluidvirt-invocation-record-v1.schema.json",
        "runtime-mutation-observation-v1.schema.json",
        "post-verify-result-v1.schema.json",
        "rollback-window-v1.schema.json",
        "apply-audit-event-v1.schema.json",
    ]
    for name in sprint4_schemas:
        if name not in {p.name for p in schema_paths}:
            fail(f"missing Sprint 4 schema: {name}")

    apply_gate = example_by_name.get("operator-apply-gate-reference.json")
    if apply_gate is None:
        fail("examples/operator-apply-gate-reference.json is missing")
    if apply_gate.get("milestone") != "hyperdensity_operator_controlled_apply_gate_v1":
        fail("operator apply gate milestone invalid")
    if apply_gate.get("operatorControlledApplyAllowed") is not True:
        fail("operatorControlledApplyAllowed must be true")
    if apply_gate.get("autoApplyAllowed") is not False:
        fail("operator apply gate autoApplyAllowed must be false")
    if apply_gate.get("productionScope") is not False:
        fail("operator apply gate productionScope must be false")

    approval_ref = example_by_name.get("operator-approval-record-reference.json")
    if approval_ref is None:
        fail("examples/operator-approval-record-reference.json is missing")
    if approval_ref.get("approvalMode") != "operator_required":
        fail("approval mode must be operator_required")

    apply_req = example_by_name.get("apply-request-reference.json")
    if apply_req is None:
        fail("examples/apply-request-reference.json is missing")
    if apply_req.get("actuator") != "FluidVirt":
        fail("apply request actuator must be FluidVirt")

    invocation_ref = example_by_name.get("fluidvirt-invocation-record-reference.json")
    if invocation_ref is None:
        fail("examples/fluidvirt-invocation-record-reference.json is missing")
    if invocation_ref.get("rawRuntimeControlsExposed") is not False:
        fail("invocation rawRuntimeControlsExposed must be false")

    sprint4_examples = [
        "operator-apply-gate-reference.json",
        "operator-approval-record-reference.json",
        "apply-request-reference.json",
        "fluidvirt-invocation-record-reference.json",
        "runtime-mutation-observation-reference.json",
        "post-verify-result-reference.json",
        "rollback-window-reference.json",
        "apply-audit-event-reference.json",
    ]
    for name in sprint4_examples:
        ex = example_by_name.get(name)
        if ex is None:
            fail(f"missing Sprint 4 example: {name}")
        positives4: list[str] = []
        collect_positive_strings(ex, positives4)
        merged4 = "\n".join(positives4).lower()
        for phrase in forbidden_positive_claims:
            if phrase in merged4:
                fail(f"{name} contains forbidden positive claim in allowed copy: {phrase}")

    forbidden_approved_phrases = [
        "windows is supported.",
        "production autonomous resource movement is supported",
        "enforcement is enabled",
        "autonomous apply is enabled",
        "production mutation is enabled",
        "vm ram live resize is generic kubevirt template mutation.",
    ]
    all_docs_merged = "\n".join(docs_lower.values())
    for phrase in forbidden_approved_phrases:
        if phrase in all_docs_merged:
            fail(f"documentation contains unsafe approved wording: {phrase}")

    print(
        f"[validate_json] OK: parsed {schema_count} schema files and {example_count} example files"
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
