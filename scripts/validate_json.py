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

    print(
        f"[validate_json] OK: parsed {schema_count} schema files and {example_count} example files"
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
