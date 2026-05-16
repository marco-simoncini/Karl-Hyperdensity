#!/usr/bin/env bash
# KHR-AW: offline test — manifest + config guards (no cluster).
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/khr_tp_live_scope1_lib.sh"

ROOT="$(khr_scope1_root)"
MANIFEST_DIR="$(khr_scope1_manifest_dir)"
CFG="${MANIFEST_DIR}/karl-host-runtime-config-scope1.yaml"

for f in \
  "${MANIFEST_DIR}/namespace-runtime.yaml" \
  "${MANIFEST_DIR}/namespace-rdpgw.yaml" \
  "${MANIFEST_DIR}/configmap-karl-host-runtime-scope1.yaml" \
  "${MANIFEST_DIR}/karl-host-runtime-preview-deployment.yaml" \
  "${CFG}"; do
  [[ -f "${f}" ]] || { echo "missing ${f}" >&2; exit 1; }
done

khr_scope1_assert_config_safe "${CFG}"
rg -q 'resourcePortLoopEnabled: false' "${CFG}"
rg -q 'sandboxApplyEnabled: false' "${CFG}"
rg -q 'khr.karl.io/sandbox' "${MANIFEST_DIR}/karl-host-runtime-preview-deployment.yaml"

echo "[khr_tp_live_scope1_test] PASS"
