#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

required_files=(
  "README.md"
  "docs/architecture/hyperdensity-overview.md"
  "docs/architecture/runtime-overlay-model.md"
  "docs/architecture/control-room-ui-concept.md"
  "docs/contracts/linux-shell-compliance-v1.md"
  "docs/contracts/resource-equilibrium-v1.md"
  "docs/contracts/fleet-equilibrium-onboarding-v1.md"
  "docs/contracts/shell-factory-v1.md"
  "docs/contracts/shell-claim-v1.md"
  "docs/contracts/shell-claim-template-profile-pack-v1.md"
  "docs/contracts/release-support-matrix-v1.md"
  "docs/contracts/evidence-bundle-demo-scenario-pack-v1.md"
  "docs/contracts/live-resource-authority-v1.md"
  "docs/contracts/action-slate-v1.md"
  "docs/contracts/guarded-auto-sandbox-v1.md"
  "docs/contracts/auto-rollback-controller-v1.md"
  "docs/contracts/blast-radius-policy-v1.md"
  "docs/contracts/policy-pack-v1.md"
  "docs/contracts/policy-pack-consistency-checker-v1.md"
  "docs/contracts/admission-guard-enforce-simulation-v1.md"
  "docs/contracts/mutate-preview-apply-dry-run-v1.md"
  "docs/runbooks/operator-runbook-v1.md"
  "docs/releases/technical-preview-release-notes-v1.md"
  "docs/releases/technical-preview-readiness-gate-v1.md"
  "docs/releases/technical-preview-release-candidate-gate-v1.md"
  "docs/demos/technical-preview-demo-guide-v1.md"
  "docs/releases/technical-preview-documentation-pack-v1.md"
  "docs/overcommit/resource-equilibrium-and-safe-overcommit.md"
  "docs/migration/dashboard-to-hyperdensity-extraction-plan.md"
  "docs/khr/KHR_LINUX_MVP_DESIGN.md"
  "docs/khr/KHR_LINUX_AGENT_RUNBOOK.md"
  "docs/khr/KHR_LINUX_CGROUP_ENVELOPE_MODEL.md"
  "docs/khr/KHR_SAFETY_AND_DRY_RUN_MODEL.md"
  "docs/khr/KHR_AUDIT_AND_APPLY_GATES.md"
  "docs/khr/KHR_LINUX_READONLY_DISCOVERY.md"
  "docs/khr/KHR_CGROUP_PATH_POLICY.md"
  "docs/khr/KHR_LINUX_READONLY_TELEMETRY.md"
  "docs/khr/KHR_TELEMETRY_EVIDENCE_MODEL.md"
  "docs/khr/KHR_LOCAL_EVIDENCE_BUNDLE.md"
  "docs/khr/KHR_GRANDE_PADRE_EVIDENCE_HANDOFF.md"
  "docs/khr/KHR_EVIDENCE_INTEGRITY_MODEL.md"
  "docs/khr/KHR_EVIDENCE_ARTIFACT_MANIFEST.md"
  "docs/khr/KARL_INFRASTRUCTURE_SCOPE.md"
  "docs/khr/KHR_RELEASE_READINESS_MAP.md"
  "docs/khr/TECHNICAL_PREVIEW_READINESS.md"
  "docs/khr/TECHNICAL_PREVIEW_PACKAGE.md"
  "docs/khr/TECHNICAL_PREVIEW_OPERATOR_RUNBOOK.md"
  "docs/khr/KHR_BOOTSTRAP_CONSUMER_EXPECTATIONS.md"
  "docs/khr/KHR_INSTALLER_PROFILE_EXPECTATIONS.md"
  "docs/contracts/khr/khr-contract-manifest.yaml"
  "docs/khr/KHR_TECHNICAL_PREVIEW_POST_INSTALL_VERIFICATION.md"
  "docs/khr/BETA_READINESS_GAP_ANALYSIS.md"
  "docs/khr/KHR_CONTRACT_FREEZE_PLAN.md"
  "docs/khr/RESOURCELEASE_TP_FREEZE_CANDIDATE.md"
  "docs/khr/RESOURCEPORT_TP_FREEZE_CANDIDATE.md"
  "docs/khr/NATIVE_LIVE_TP_FREEZE_CANDIDATE.md"
  "docs/ingest/KHR_EVIDENCE_INGEST_CONTRACT.md"
  "docs/ingest/EVIDENCEBUNDLE_CRD.md"
  "docs/ingest/EVIDENCEINGESTREQUEST_CRD.md"
  "docs/ingest/GRANDE_PADRE_EVIDENCE_INDEXING_MODEL.md"
  "docs/ingest/EVIDENCE_INGEST_SECURITY_BOUNDARIES.md"
  "docs/grandepadre/GRANDE_PADRE_EVIDENCE_STORE_SKELETON.md"
  "docs/grandepadre/GRANDE_PADRE_EVIDENCE_INDEXES.md"
  "docs/grandepadre/GRANDE_PADRE_BLOCKED_REMEDIABLE_INDEX.md"
  "docs/grandepadre/GRANDE_PADRE_RECOMMENDATION_ENGINE.md"
  "docs/grandepadre/GRANDE_PADRE_ACTION_SLATE_MODEL.md"
  "docs/grandepadre/GRANDE_PADRE_DONOR_RECEIVER_INDEX.md"
  "docs/grandepadre/GRANDE_PADRE_RISK_PRIORITY_MODEL.md"
  "docs/windowsfluidvirt/WINDOWS_FLUIDVIRT_LANE_RECONCILIATION.md"
  "docs/windowsfluidvirt/WINDOWS_FLUIDVIRT_PLANNING_ONLY_SAFETY.md"
  "docs/extraction/HYPERDENSITY_REAL_EXTRACTION_AUDIT.md"
  "docs/extraction/HYPERDENSITY_PACKAGE_TARGET_PLAN.md"
  "docs/extraction/HYPERDENSITY_KHR_DUPLICATION_REPORT.md"
  "docs/roadmap/KHR_HYPERDENSITY_CORRECTED_ROADMAP.md"
  "docs/extraction/HYPERDENSITY_GOLDEN_ANCHOR_M1.md"
  "docs/extraction/HYPERDENSITY_SUMMARY_FIELD_MAPPING_M2.md"
  "docs/extraction/HYPERDENSITY_SUMMARY_FIELD_MAPPING_M3.md"
  "docs/extraction/HYPERDENSITY_SUMMARY_FIELD_MAPPING_M4.md"
  "docs/extraction/HYPERDENSITY_SUMMARY_FIELD_MAPPING_M5.md"
  "docs/extraction/HYPERDENSITY_SUMMARY_FIELD_MAPPING_M6.md"
  "docs/extraction/HYPERDENSITY_SUMMARY_FIELD_MAPPING_M7.md"
  "docs/extraction/HYPERDENSITY_PARITY_MATRIX_M1_M7.md"
  "docs/extraction/HYPERDENSITY_M1_M7_EXTRACTION_READINESS.md"
  "docs/extraction/HYPERDENSITY_CONTRACTKIT_MODULE.md"
  "docs/extraction/HYPERDENSITY_CONTRACTKIT_FIXTURE_MANIFEST.md"
  "docs/extraction/HYPERDENSITY_CONTRACTKIT_RELEASE_TAGGING.md"
  "docs/extraction/HYPERDENSITY_CONTRACTKIT_VERSION_MODEL.md"
  "docs/extraction/HYPERDENSITY_CONTRACTKIT_NO_REPUBLISH_POLICY.md"
  "docs/extraction/HYPERDENSITY_CONTRACTKIT_CONSUMER_CI_HARDENING.md"
  "docs/extraction/HYPERDENSITY_CONTRACTKIT_CONSUMER_AUDIT.md"
  "docs/extraction/HYPERDENSITY_CONTRACTKIT_CONSUMER_POLICY.md"
  "docs/extraction/HYPERDENSITY_CONTRACTKIT_SECOND_CONSUMER_PLAYBOOK.md"
  "docs/extraction/templates/audit_contractkit_module_pin.sh"
  "docs/extraction/templates/CONTRACTKIT_CONSUMER_DECISION_RECORD.md"
  "docs/extraction/HYPERDENSITY_CONTRACTKIT_CLAIMPOLICY.md"
  "docs/extraction/HYPERDENSITY_CLAIMPOLICY_SURFACE_MAPPING.md"
  "docs/extraction/HYPERDENSITY_CLAIMPOLICY_SURFACE_TRACEABILITY.md"
  "docs/extraction/HYPERDENSITY_CLAIMPOLICY_TRACEABILITY_TOKEN_GUARD.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_EXTRACTION_BOUNDARY.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_EXTRACTION_PHASES.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_DEPENDENCY_GUARDS.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_PURE_PACKAGE_SKELETON.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_PURE_CANDIDATE_AUDIT.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_EXECUTION_TYPES_AUDIT.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_EXECUTION_TYPES_CONTRACT.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_EXECUTION_TYPES_DRIFT_GUARD.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_HELPERS_AUDIT.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_HELPERS_DEFERRED.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_PRIMITIVES_CONTRACT.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_PREREQUISITES.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_BOUNDARY.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_INTERFACE_PROPOSAL.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_REAUDIT_CRITERIA.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_STUB_READINESS.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_PURE_CANDIDATES_CONTRACT.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_HARDENING.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_SHADOW_TESTS.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_WIRING_PROPOSAL.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_PATH_WIRING.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_PILOT_OBSERVATION_WIRING.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_PILOT_OBSERVATION_HARDENING.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_LIVE_OBSERVATION_PROPOSAL.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_LIVE_OBSERVATION_STAGED_WIRING.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_LIVE_OBSERVATION_SHADOW_HARDENING.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_LIVE_OBSERVATION_FLIP_CRITERIA.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_LIVE_OBSERVATION_FLIP.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_LIVE_OBSERVATION_SEMANTIC_PROTOTYPE.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_LIVE_OBSERVATION_BRANCH_SWAP_CRITERIA.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_LIVE_OBSERVATION_BRANCH_SWAP.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_OBSERVATION_REAUDIT.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_BROAD_OBSERVATION_DECISION.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_PROPOSAL.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_CRITERIA.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_SHADOW_MATRIX.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_STAGED_WRAPPERS.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_WRAPPER_HARDENING.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_WIRING_READINESS.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_CALLSITE_WIRING.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_POST_WIRING_HARDENING.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_FLIP_CRITERIA.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_FLIP_RISKS.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_CANDIDATE_RUNTIME_READINESS.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_BRANCH_LOGIC.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_CANDIDATE_RUNTIME_FLIP.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_ACTIVATION.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_POST_ACTIVATION_HARDENING.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_MIGRATION_BOUNDARY.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_OBSERVATION_AUDIT.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_OBSERVATION_CRITERIA.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_OBSERVATION_SHADOW_MATRIX.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_LOCAL_HELPER_CLASSIFICATION.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_OBSERVATION_STAGED_WRAPPERS.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_OBSERVATION_STAGED_WRAPPER_CRITERIA.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_LOCAL_HELPER_SHADOW_MATRIX.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_LOCAL_HELPER_CRITERIA.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_FULL_HELPER_STAGED_WRAPPERS.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_CALLSITE_WIRING_READINESS.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_CALLSITE_WIRING.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_CANDIDATE_RUNTIME_STAGING.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_CANDIDATE_RUNTIME_REVALIDATION.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_ACTIVATION_READINESS.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_ACTIVATION_RISKS.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_ACTIVATION.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_POST_ACTIVATION_HARDENING.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_MIGRATION_BOUNDARY.md"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_REMAINING_SURFACE_DECISION.md"
  "docs/extraction/HYPERDENSITY_KHR_ROADMAP_TRANSITION_NOTE.md"
  "docs/extraction/HYPERDENSITY_KHR_ARCHITECTURE_MEMORY.md"
  "docs/extraction/HYPERDENSITY_KHR_STORAGE_SEMANTICS.md"
  "docs/extraction/HYPERDENSITY_KHR_NETWORK_SDN_SEMANTICS.md"
  "docs/extraction/HYPERDENSITY_KHR_RESOURCELEASE_DIRECTION.md"
  "docs/extraction/HYPERDENSITY_KHR_ROADMAP_PHASES.md"
  "docs/extraction/HYPERDENSITY_KHR_RESOURCELEASE_MINIMAL_CONTRACT.md"
  "docs/extraction/HYPERDENSITY_KHR_RESOURCELEASE_STORAGE_CONTRACT.md"
  "docs/extraction/HYPERDENSITY_KHR_RESOURCELEASE_NETWORK_CONTRACT.md"
  "docs/extraction/HYPERDENSITY_KHR_RESOURCELEASE_PROVIDER_CONTRACT.md"
  "docs/extraction/HYPERDENSITY_KHR_RESOURCELEASE_EXAMPLES.md"
  "docs/extraction/HYPERDENSITY_KHR_INVENTORY_SCAN_METHOD.md"
  "docs/extraction/HYPERDENSITY_KHR_KUBEVIRT_CAPABILITY_INVENTORY.md"
  "docs/extraction/HYPERDENSITY_KHR_OVN_SDN_CAPABILITY_INVENTORY.md"
  "docs/extraction/HYPERDENSITY_KHR_RESOURCELEASE_INVENTORY_MAPPING.md"
  "docs/extraction/HYPERDENSITY_KHR_RESOURCELEASE_JSON_SCHEMA.md"
  "docs/contracts/khr/resourcelease.schema.json"
  "docs/contracts/khr/resourcelease.schema.manifest.json"
  "docs/contracts/khr/examples/resourcelease-windows-daas-khr-native.json"
  "docs/contracts/khr/examples/resourcelease-public-cloud-kubevirt-fallback.json"
  "docs/contracts/khr/examples/resourcelease-baremetal-native.json"
  "docs/extraction/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_OBSERVATION_NEXT_SURFACE_DECISION.md"
  "pkg/hyperdensity/contractkit/go.mod"
  "pkg/hyperdensity/contractkit/testdata/dashboard/hyperdensity_parity_manifest_m1_m7.json"
  "testdata/dashboard/parent_fabric_summary_redacted.golden.json"
)

for required in "${required_files[@]}"; do
  if [[ ! -f "$required" ]]; then
    echo "[validate] ERROR: missing required file: $required" >&2
    exit 1
  fi
done

if [[ -x "${ROOT_DIR}/scripts/validate_parentfabric_pure_deps.sh" ]]; then
  "${ROOT_DIR}/scripts/validate_parentfabric_pure_deps.sh"
fi

go test ./...
(
  cd pkg/hyperdensity/contractkit
  go test ./...
)
python3 scripts/validate_json.py
python3 scripts/validate_khr_examples.py

if [[ -x "${ROOT_DIR}/scripts/validate_resourcelease_schema.sh" ]]; then
  "${ROOT_DIR}/scripts/validate_resourcelease_schema.sh"
fi

if [[ -x "${ROOT_DIR}/scripts/validate_resourceport_schema.sh" ]]; then
  "${ROOT_DIR}/scripts/validate_resourceport_schema.sh"
fi

if [[ -x "${ROOT_DIR}/scripts/validate_windows_live_scale_contract.sh" ]]; then
  "${ROOT_DIR}/scripts/validate_windows_live_scale_contract.sh"
fi

if [[ -x "${ROOT_DIR}/scripts/validate_lane_discovery.sh" ]]; then
  "${ROOT_DIR}/scripts/validate_lane_discovery.sh"
fi

if [[ -x "${ROOT_DIR}/scripts/validate_resourcefuture_simulation.sh" ]]; then
  "${ROOT_DIR}/scripts/validate_resourcefuture_simulation.sh"
fi

if [[ -x "${ROOT_DIR}/scripts/validate_shell_cell_schema.sh" ]]; then
  "${ROOT_DIR}/scripts/validate_shell_cell_schema.sh"
fi

if [[ -x "${ROOT_DIR}/scripts/validate_shelllease_gatewayroute_schema.sh" ]]; then
  "${ROOT_DIR}/scripts/validate_shelllease_gatewayroute_schema.sh"
fi

if [[ -x "${ROOT_DIR}/scripts/validate_karl_host_runtime.sh" ]]; then
  "${ROOT_DIR}/scripts/validate_karl_host_runtime.sh"
fi

if [[ -x "${ROOT_DIR}/scripts/validate_khr_runtime_sandbox.sh" ]]; then
  "${ROOT_DIR}/scripts/validate_khr_runtime_sandbox.sh"
fi

if [[ -x "${ROOT_DIR}/scripts/validate_host_schema.sh" ]]; then
  "${ROOT_DIR}/scripts/validate_host_schema.sh"
fi

if [[ -x "${ROOT_DIR}/scripts/validate_resourceport_loop.sh" ]]; then
  "${ROOT_DIR}/scripts/validate_resourceport_loop.sh"
fi

if [[ -x "${ROOT_DIR}/scripts/validate_crds.sh" ]]; then
  "${ROOT_DIR}/scripts/validate_crds.sh"
fi

if [[ -x "${ROOT_DIR}/scripts/guard_khr_docs_scope.sh" ]]; then
  "${ROOT_DIR}/scripts/guard_khr_docs_scope.sh"
fi

if [[ -x "${ROOT_DIR}/scripts/khr_tp_package_check.sh" ]]; then
  "${ROOT_DIR}/scripts/khr_tp_package_check.sh"
fi

if [[ -x "${ROOT_DIR}/scripts/khr_tp_operator_bundle.sh" ]]; then
  "${ROOT_DIR}/scripts/khr_tp_operator_bundle.sh"
fi

if [[ -x "${ROOT_DIR}/scripts/khr_contract_inventory.sh" ]]; then
  "${ROOT_DIR}/scripts/khr_contract_inventory.sh"
fi

if [[ -x "${ROOT_DIR}/scripts/khr_contract_manifest_check.sh" ]]; then
  "${ROOT_DIR}/scripts/khr_contract_manifest_check.sh"
fi

if [[ -x "${ROOT_DIR}/scripts/khr_tp_post_install_bundle_check.sh" ]]; then
  "${ROOT_DIR}/scripts/khr_tp_post_install_bundle_check.sh"
fi

if [[ -x "${ROOT_DIR}/scripts/khr_resourcelease_freeze_check.sh" ]]; then
  "${ROOT_DIR}/scripts/khr_resourcelease_freeze_check.sh"
fi

if [[ -x "${ROOT_DIR}/scripts/khr_native_live_freeze_check.sh" ]]; then
  "${ROOT_DIR}/scripts/khr_native_live_freeze_check.sh"
fi

if [[ -x "${ROOT_DIR}/scripts/khr_runtime_observation_federation_check.sh" ]]; then
  "${ROOT_DIR}/scripts/khr_runtime_observation_federation_check.sh"
fi

if [[ -x "${ROOT_DIR}/scripts/khr_tp_live_scope1_test.sh" ]]; then
  "${ROOT_DIR}/scripts/khr_tp_live_scope1_test.sh"
fi

if [[ -x "${ROOT_DIR}/scripts/khr_tp_live_enablement_preflight.sh" ]]; then
  "${ROOT_DIR}/scripts/khr_tp_live_enablement_preflight.sh"
fi

if [[ -x "${ROOT_DIR}/scripts/khr_tp_live_scope2_preflight.sh" ]]; then
  KHR_TP_LIVE_SCOPE2_RUN_ID="${KHR_TP_LIVE_SCOPE2_RUN_ID:-committed-scope2-preflight-khr-az}" \
    "${ROOT_DIR}/scripts/khr_tp_live_scope2_preflight.sh"
fi

if [[ -x "${ROOT_DIR}/scripts/khr_tp_live_scope3_preflight.sh" ]]; then
  KHR_TP_LIVE_SCOPE3_RUN_ID="${KHR_TP_LIVE_SCOPE3_RUN_ID:-committed-scope3-preflight-khr-bb}" \
    "${ROOT_DIR}/scripts/khr_tp_live_scope3_preflight.sh"
fi

if [[ -x "${ROOT_DIR}/scripts/khr_tp_live_scope2_resourceport_loop_run.sh" ]]; then
  KHR_TP_LIVE_SCOPE2_I_UNDERSTAND_MANUAL_LOOP="${KHR_TP_LIVE_SCOPE2_I_UNDERSTAND_MANUAL_LOOP:-true}" \
  KHR_TP_LIVE_SCOPE2_LOOP_RUN_ID="${KHR_TP_LIVE_SCOPE2_LOOP_RUN_ID:-committed-scope2-loop-khr-ba}" \
  KHR_SCOPE2_LOOP_ITERATIONS="${KHR_SCOPE2_LOOP_ITERATIONS:-2}" \
    "${ROOT_DIR}/scripts/khr_tp_live_scope2_resourceport_loop_run.sh"
  KHR_TP_LIVE_SCOPE2_LOOP_RUN_ID="${KHR_TP_LIVE_SCOPE2_LOOP_RUN_ID:-committed-scope2-loop-khr-ba}" \
    "${ROOT_DIR}/scripts/khr_tp_live_scope2_resourceport_loop_verify.sh"
  KHR_TP_LIVE_SCOPE2_LOOP_RUN_ID="${KHR_TP_LIVE_SCOPE2_LOOP_RUN_ID:-committed-scope2-loop-khr-ba}" \
    "${ROOT_DIR}/scripts/khr_tp_live_scope2_resourceport_loop_cleanup.sh"
fi

if [[ -x "${ROOT_DIR}/scripts/khr_tp_live_reference_env_check.sh" ]]; then
  KHR_DASHBOARD_PATH="${KHR_DASHBOARD_PATH:-${ROOT_DIR}/../Karl-Dashboard}" \
    "${ROOT_DIR}/scripts/khr_tp_live_reference_env_check.sh"
fi

schema_count="$(ls -1 schemas/*.json | wc -l | tr -d ' ')"
example_count="$(ls -1 examples/*.json | wc -l | tr -d ' ')"
doc_count="${#required_files[@]}"

echo "[validate] SUCCESS: go tests + JSON validation passed (schemas=${schema_count}, examples=${example_count}, required_docs=${doc_count})"
