#!/usr/bin/env bash
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/khr_runtime_sandbox_lib.sh"

NS="${KHR_RUNTIME_SANDBOX_NS}"
MANIFEST_DIR="$(khr_runtime_sandbox_manifest_dir)"
BIN="$(khr_runtime_build_binary)"
CFG="${MANIFEST_DIR}/karl-host-runtime-config-apply.yaml"
SANDBOX_DIR="$(mktemp -d)"
BASELINE_OUT="$(mktemp)"
ROLLBACK_OUT="$(mktemp)"

khr_runtime_assert_cluster_context
khr_runtime_assert_namespace_arg "${NS}"
khr_runtime_assert_sandbox_namespace_labels

# Ensure marker exists from prior apply step or create minimal state
echo "pre-rollback-marker" > "${SANDBOX_DIR}/apply-marker.txt"
cp "${SANDBOX_DIR}/apply-marker.txt" "${SANDBOX_DIR}/baseline-copy.txt"

"${BIN}" -mode=rollback -config="${CFG}" -sandbox-dir="${SANDBOX_DIR}" -baseline-id=sandbox-khr-g > "${ROLLBACK_OUT}"
if ! jq -e '.rolledBack == true' "${ROLLBACK_OUT}" >/dev/null 2>&1; then
  echo "BLOCKED: rollback failed" >&2
  cat "${ROLLBACK_OUT}" >&2
  exit 1
fi

{
  echo "{"
  echo "  \"phase\": \"rollback\","
  echo "  \"status\": \"PASS\","
  echo "  \"namespace\": \"${NS}\","
  echo "  \"baselineRequired\": true,"
  echo "  \"rollbackResult\": $(cat "${ROLLBACK_OUT}"),"
  echo "  \"at\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\""
  echo "}"
} > "${BASELINE_OUT}"

khr_runtime_artifact_json "rollback-baseline.json" "${BASELINE_OUT}"
khr_runtime_artifact_json "rollback-result.json" "${ROLLBACK_OUT}"
khr_runtime_artifact_text "rollback.log" "rollback baseline PASS"

echo "[khr_runtime_sandbox_rollback] PASS"
