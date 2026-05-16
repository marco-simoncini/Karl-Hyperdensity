#!/usr/bin/env bash
set -euo pipefail
# shellcheck source=scripts/khr_runtime_sandbox_lib.sh
source "$(dirname "${BASH_SOURCE[0]}")/khr_runtime_sandbox_lib.sh"

NS="${KHR_RUNTIME_SANDBOX_NS}"
CTX="${KHR_RUNTIME_CLUSTER_CONTEXT}"
MANIFEST_DIR="$(khr_runtime_sandbox_manifest_dir)"
OUT="$(mktemp)"

khr_runtime_assert_cluster_context
khr_runtime_block_production_namespace "${NS}"
khr_runtime_assert_namespace_arg "${NS}"

khr_runtime_log "applying sandbox manifests (namespace, configmap, test target)"
kubectl --context "${CTX}" apply -f "${MANIFEST_DIR}/namespace.yaml"
kubectl --context "${CTX}" apply -f "${MANIFEST_DIR}/configmap-karl-host-runtime.yaml"
kubectl --context "${CTX}" apply -f "${MANIFEST_DIR}/test-target-linux.yaml"
kubectl --context "${CTX}" -n "${NS}" rollout status deployment/khr-runtime-linux-target --timeout=120s

khr_runtime_assert_sandbox_namespace_labels
khr_runtime_assert_workload_labels

{
  echo "{"
  echo "  \"phase\": \"preflight\","
  echo "  \"status\": \"PASS\","
  echo "  \"clusterContext\": \"${CTX}\","
  echo "  \"namespace\": \"${NS}\","
  echo "  \"requiredLabel\": \"${KHR_RUNTIME_SANDBOX_LABEL}\","
  echo "  \"sandboxApplyEnabledDefault\": false,"
  echo "  \"blockedSurfaces\": [\"kubevirt\", \"libvirt\", \"qmp\", \"windows\", \"autonomous-apply\"],"
  echo "  \"productionNamespacesBlocked\": $(printf '%s\n' "${KHR_RUNTIME_PRODUCTION_NS_BLOCKLIST[@]}" | jq -R . | jq -s .),"
  echo "  \"at\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\""
  echo "}"
} > "${OUT}"

khr_runtime_artifact_json "preflight.json" "${OUT}"
khr_runtime_artifact_text "preflight.log" "preflight PASS namespace=${NS} context=${CTX}"
echo "[khr_runtime_sandbox_preflight] PASS"
