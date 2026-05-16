#!/usr/bin/env bash
# KHR-X: unified control graph export from real cluster state (read-only).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
# shellcheck source=scripts/khr_runtime_sandbox_lib.sh
source "${ROOT}/scripts/khr_runtime_sandbox_lib.sh"

CTX="${KHR_RUNTIME_CLUSTER_CONTEXT}"
EVIDENCE="${ROOT}/docs/evidence/khr-control-graph"
RUN_ID="${KHR_CONTROL_GRAPH_RUN_ID:-$(date -u +%Y%m%dT%H%M%SZ)}"
RUN_DIR="${EVIDENCE}/${RUN_ID}"
CFG="${ROOT}/examples/khr/runtime-sandbox/karl-host-runtime-config-resourcefuture-simulate.yaml"
REGISTRY="${ROOT}/docs/evidence/khr-certification-registry/registry.json"
BIN_GRAPH="${ROOT}/bin/khr-control-graph"
BIN_HOST="${ROOT}/bin/karl-host-runtime"

khr_runtime_assert_cluster_context
mkdir -p "${RUN_DIR}" "${ROOT}/bin"

(cd "${ROOT}" && go build -o "${BIN_GRAPH}" ./cmd/khr-control-graph)
(cd "${ROOT}" && go build -o "${BIN_HOST}" ./cmd/karl-host-runtime)

APPROVAL_BUNDLE=""
SIM="${RUN_DIR}/simulation-gated.json"
if [[ -f "${REGISTRY}" ]]; then
  "${BIN_HOST}" -mode=resourcefuture-simulate -config="${CFG}" -cluster-context="${CTX}" \
    -cert-registry="${REGISTRY}" > "${SIM}"
fi
LATEST_APPROVAL="$(ls -d "${ROOT}"/docs/evidence/khr-action-approval/20* 2>/dev/null | sort | tail -1 || true)"
if [[ -n "${LATEST_APPROVAL}" && -f "${LATEST_APPROVAL}/pending-bundle.json" ]]; then
  APPROVAL_BUNDLE="${LATEST_APPROVAL}/pending-bundle.json"
fi

ARGS=(-config="${CFG}" -cluster-context="${CTX}" -sprint=KHR-X -out="${RUN_DIR}/control-graph.json")
[[ -f "${REGISTRY}" ]] && ARGS+=(-registry="${REGISTRY}")
[[ -f "${SIM}" ]] && ARGS+=(-simulation="${SIM}")
[[ -n "${APPROVAL_BUNDLE}" ]] && ARGS+=(-approvals="${APPROVAL_BUNDLE}")

khr_runtime_log "export unified control graph"
"${BIN_GRAPH}" "${ARGS[@]}"

jq -e '.graphId == "khr-control-graph-v1"' "${RUN_DIR}/control-graph.json" >/dev/null
jq -e '.noApply == true and .noMutation == true and .noAutonomousOrchestration == true' \
  "${RUN_DIR}/control-graph.json" >/dev/null
jq -e '.nodes | length >= 4' "${RUN_DIR}/control-graph.json" >/dev/null
jq -e '.edges | length >= 2' "${RUN_DIR}/control-graph.json" >/dev/null
jq -e '.correlationId != ""' "${RUN_DIR}/control-graph.json" >/dev/null

echo "no-apply: graph export only" > "${RUN_DIR}/mutation-check.txt"
echo "no-orchestration: read-only discovery" >> "${RUN_DIR}/mutation-check.txt"

jq -n \
  --arg sprint "KHR-X" \
  --arg runId "${RUN_ID}" \
  --arg cluster "${CTX}" \
  --argjson nodes "$(jq '.health.nodeCount' "${RUN_DIR}/control-graph.json")" \
  --argjson orphans "$(jq '.health.orphanCount' "${RUN_DIR}/control-graph.json")" \
  --argjson stale "$(jq '.health.staleCount' "${RUN_DIR}/control-graph.json")" \
  --argjson consistent "$(jq '.health.consistent' "${RUN_DIR}/control-graph.json")" \
  '{sprint:$sprint,runId:$runId,cluster:$cluster,nodeCount:$nodes,orphanCount:$orphans,staleCount:$stale,consistent:$consistent,readOnly:true,noApply:true,noAutonomousOrchestration:true}' \
  > "${RUN_DIR}/summary.json"
cp "${RUN_DIR}/summary.json" "${EVIDENCE}/summary.json"
cp "${RUN_DIR}/control-graph.json" "${EVIDENCE}/control-graph.json"

echo "[khr_control_graph_evidence] PASS ${RUN_DIR}"
