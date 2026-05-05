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
