#!/usr/bin/env bash
# KHR-I: generate Host status JSON on karl-metal-01@ovh (read-only; no CR apply).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
CTX="${KHR_RUNTIME_CLUSTER_CONTEXT:-karl-metal-01@ovh}"
NODE="${KHR_HOST_NODE_NAME:-karl-metal-01}"
CFG="${ROOT}/examples/khr/runtime-sandbox/karl-host-runtime-config-default.yaml"
EVIDENCE="${ROOT}/docs/evidence/khr-host-registration"
RUN_ID="${KHR_HOST_RUN_ID:-$(date -u +%Y%m%dT%H%M%SZ)}"
RUN_DIR="${EVIDENCE}/${RUN_ID}"

current="$(kubectl config current-context 2>/dev/null || true)"
if [[ "${current}" != "${CTX}" ]]; then
  echo "BLOCKED: context=${current:-<none>} required=${CTX}" >&2
  exit 2
fi

mkdir -p "${RUN_DIR}"
BIN="${ROOT}/bin/karl-host-runtime"
mkdir -p "${ROOT}/bin"
(cd "${ROOT}" && go build -o "${BIN}" ./cmd/karl-host-runtime)

"${BIN}" -mode=host-status -config="${CFG}" -node-name="${NODE}" \
  -namespace=khr-runtime-sandbox -port-name=khr-runtime-sandbox-port \
  > "${RUN_DIR}/host-status.json"

jq -e '.kind == "Host" and .status.safetyMode == "sandbox"' "${RUN_DIR}/host-status.json" >/dev/null

jq -n \
  --arg sprint "KHR-I" \
  --arg status "PASS" \
  --arg ctx "${CTX}" \
  --arg node "${NODE}" \
  --arg runId "${RUN_ID}" \
  --arg at "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
  '{
    sprint: $sprint,
    status: $status,
    clusterContext: $ctx,
    nodeName: $node,
    noProductionMutation: true,
    noController: true,
    hostCRApplied: false,
    runId: $runId,
    at: $at
  }' > "${RUN_DIR}/summary.json"

cp "${RUN_DIR}/summary.json" "${EVIDENCE}/summary.json"
echo "[khr_host_registration_evidence] PASS -> ${EVIDENCE}/summary.json"
