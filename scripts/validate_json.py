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

    print(
        f"[validate_json] OK: parsed {schema_count} schema files and {example_count} example files"
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
