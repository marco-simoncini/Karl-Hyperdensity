#!/usr/bin/env bash
# KHR-AL: aggregate read-only TP post-install verification bundle.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${ROOT}"

FAIL=0
RUN_ID="$(date -u +%Y%m%dT%H%M%SZ)"
OUT_DIR="${ROOT}/docs/evidence/khr-tp-post-install-bundle/${RUN_ID}"
mkdir -p "${OUT_DIR}"

record() {
  local name="$1" status="$2"
  echo "${name}:${status}" >> "${OUT_DIR}/checks.txt"
  [[ "${status}" == "PASS" ]] || FAIL=1
  echo "[khr_tp_post_install_bundle] ${name} ${status}"
}

echo "[khr_tp_post_install_bundle] runId=${RUN_ID}"

if [[ -x "${ROOT}/scripts/khr_contract_manifest_check.sh" ]]; then
  "${ROOT}/scripts/khr_contract_manifest_check.sh" >/dev/null 2>&1 && record "khr_contract_manifest_check" "PASS" || record "khr_contract_manifest_check" "FAIL"
else
  record "khr_contract_manifest_check" "SKIP"
fi

if [[ -x "${ROOT}/scripts/khr_tp_operator_bundle.sh" ]]; then
  if "${ROOT}/scripts/khr_tp_operator_bundle.sh" >/tmp/khr_tp_ob.log 2>&1; then
    record "khr_tp_operator_bundle" "PASS"
  else
    record "khr_tp_operator_bundle" "FAIL"
    tail -5 /tmp/khr_tp_ob.log >&2 || true
  fi
else
  record "khr_tp_operator_bundle" "SKIP"
fi

ISO="../Karl-OS-ISO"
if [[ -x "${ISO}/scripts/khr_post_install_verify.sh" ]]; then
  (cd "${ISO}" && ./scripts/khr_post_install_verify.sh) >/dev/null 2>&1 && record "khr_post_install_verify" "PASS" || record "khr_post_install_verify" "FAIL"
else
  record "khr_post_install_verify" "SKIP"
fi

INST="../Karl-Installer"
KARL2_EVIDENCE="${INST}/docs/evidence/khr-installer-crd-foundation/summary.json"
HYBRID_EVIDENCE="${INST}/docs/evidence/khr-hybrid-transition/summary.json"

if [[ -f "${KARL2_EVIDENCE}" ]]; then
  rg -q '"contractSetId".*khr-tp-contract-v1' "${KARL2_EVIDENCE}" && record "installer_karl2_evidence" "PASS" || record "installer_karl2_evidence" "FAIL"
  rg -q '"crdDiffEmpty".*true' "${KARL2_EVIDENCE}" || echo "[khr_tp_post_install_bundle] INFO: karl2 crdDiffEmpty not true (apply may be skipped)"
else
  record "installer_karl2_evidence" "INFO"
  echo "[khr_tp_post_install_bundle] INFO: optional karl2 evidence missing"
fi

if [[ -f "${HYBRID_EVIDENCE}" ]]; then
  rg -q '"kubevirtCompatibility".*true' "${HYBRID_EVIDENCE}" && record "installer_hybrid_evidence" "PASS" || record "installer_hybrid_evidence" "FAIL"
  rg -q '"hostRuntimeEnabled".*false' "${HYBRID_EVIDENCE}" && record "hybrid_host_runtime_disabled" "PASS" || record "hybrid_host_runtime_disabled" "FAIL"
else
  record "installer_hybrid_evidence" "INFO"
  echo "[khr_tp_post_install_bundle] INFO: optional hybrid evidence missing"
fi

ISO_POST="${ISO}/docs/evidence/khr-post-install-verify/summary.json"
[[ -f "${ISO_POST}" ]] && record "iso_post_install_summary" "PASS" || record "iso_post_install_summary" "INFO"

{
  echo "{"
  echo "  \"phase\": \"khr-tp-post-install-bundle\","
  echo "  \"sprint\": \"KHR-AL\","
  echo "  \"runId\": \"${RUN_ID}\","
  echo "  \"contractSetId\": \"khr-tp-contract-v1\","
  echo "  \"status\": \"$([[ "${FAIL}" -eq 0 ]] && echo PASS || echo FAIL)\","
  echo "  \"noAutonomousOrchestration\": true,"
  echo "  \"productionReady\": false,"
  echo "  \"at\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\""
  echo "}"
} > "${OUT_DIR}/summary.json"
cp "${OUT_DIR}/summary.json" "${ROOT}/docs/evidence/khr-tp-post-install-bundle/summary.json" 2>/dev/null || true

[[ "${FAIL}" -eq 0 ]] || exit 1
echo "[khr_tp_post_install_bundle] PASS"
