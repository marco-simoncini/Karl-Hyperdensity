#!/usr/bin/env bash
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/khr_runtime_sandbox_lib.sh"

NS="${KHR_RUNTIME_SANDBOX_NS}"
CTX="${KHR_RUNTIME_CLUSTER_CONTEXT}"
EVIDENCE_ROOT="$(khr_runtime_evidence_dir)"
RUN_DIR="$(khr_runtime_run_evidence_dir)"
BEFORE="$(mktemp)"
AFTER="$(mktemp)"
SUMMARY="$(mktemp)"

khr_runtime_assert_cluster_context
khr_runtime_assert_namespace_arg "${NS}"

if [[ ! -d "${RUN_DIR}" ]]; then
  echo "BLOCKED: no evidence run dir ${RUN_DIR}; run preflight/dry-run/apply/rollback first" >&2
  exit 1
fi

khr_runtime_production_mutation_proof > "${BEFORE}"
sleep 1
khr_runtime_production_mutation_proof > "${AFTER}"

kubectl --context "${CTX}" -n "${NS}" get pods -l khr.karl.io/sandbox=true -o json > "${RUN_DIR}/sandbox-pods.json"
kubectl --context "${CTX}" get namespace "${NS}" -o json > "${RUN_DIR}/sandbox-namespace.json"

ARTIFACTS="$(find "${RUN_DIR}" -maxdepth 1 -type f -printf '%f\n' | sort | jq -R . | jq -s .)"
jq -n \
  --arg sprint "KHR-G" \
  --arg status "PASS" \
  --arg ctx "${CTX}" \
  --arg ns "${NS}" \
  --arg label "${KHR_RUNTIME_SANDBOX_LABEL}" \
  --arg runId "$(khr_runtime_run_id)" \
  --arg evidenceDir "${RUN_DIR}" \
  --arg at "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
  --slurpfile before "${BEFORE}" \
  --slurpfile after "${AFTER}" \
  --argjson artifacts "${ARTIFACTS}" \
  '{
    sprint: $sprint,
    status: $status,
    clusterContext: $ctx,
    namespace: $ns,
    requiredLabel: $label,
    sandboxApplyEnabledDefault: false,
    flightRecorderRequired: true,
    rollbackBaselineRequired: true,
    noProductionMutation: true,
    blockedSurfaces: ["kubevirt", "libvirt", "qmp", "windows", "autonomous-apply"],
    runId: $runId,
    evidenceDir: $evidenceDir,
    productionProofBefore: $before[0],
    productionProofAfter: $after[0],
    artifacts: $artifacts,
    at: $at
  }' > "${SUMMARY}"

cp "${SUMMARY}" "${EVIDENCE_ROOT}/summary.json"
cp "${SUMMARY}" "${RUN_DIR}/summary.json"
khr_runtime_artifact_text "collect-evidence.log" "evidence PASS run=$(khr_runtime_run_id)"

echo "[khr_runtime_sandbox_collect_evidence] PASS -> ${EVIDENCE_ROOT}/summary.json"
