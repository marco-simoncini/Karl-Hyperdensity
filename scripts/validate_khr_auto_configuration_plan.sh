#!/usr/bin/env bash
# KHR-CK: validate KARL 2.0 auto-configuration plan docs and plan-only fixture (no runtime).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
export ROOT

echo "[validate_khr_auto_configuration_plan] Doc presence..."
python3 - <<'PY'
import json
import os
import sys
from pathlib import Path

root = Path(os.environ["ROOT"])
errors = []

cross = {
    "hyperdensity": root / "docs/khr/KHR_AUTO_CONFIGURATION_PLAN.md",
    "installer": Path(os.environ.get("KARL_INSTALLER_ROOT", "/home/m.simoncini/GitHub/Karl-Installer"))
    / "docs/khr/INSTALLER_KARL2_BAREMETAL_KHR_NATIVE_PROFILE.md",
    "iso": Path(os.environ.get("KARL_ISO_ROOT", "/home/m.simoncini/GitHub/Karl-OS-ISO"))
    / "docs/khr/ISO_KARL2_AUTO_CONFIGURATION_BOUNDARY.md",
    "dashboard": Path(os.environ.get("KARL_DASHBOARD_ROOT", "/home/m.simoncini/GitHub/Karl-Dashboard"))
    / "docs/khr/DASHBOARD_KARL2_AUTO_CONFIGURATION.md",
    "inventory": Path(os.environ.get("KARL_INVENTORY_ROOT", "/home/m.simoncini/GitHub/Karl-Inventory"))
    / "docs/khr/INVENTORY_KARL2_AUTO_INGEST.md",
    "rdpgw": Path(os.environ.get("RDPGW_ROOT", "/home/m.simoncini/rdp-GW"))
    / "docs/khr/RDPGW_KARL2_REFERENCE_CONFIGURATION.md",
}

for name, path in cross.items():
    if not path.is_file():
        errors.append(f"missing {name}: {path}")

plan = cross["hyperdensity"]
if plan.is_file():
    text = plan.read_text(encoding="utf-8").lower()
    for phrase in (
        "crd foundation",
        "host-runtime preview",
        "resourceport loop",
        "dry-run",
        "guarded apply",
        "governance",
        "plan-only",
        "no global",
        "hyperdensity",
        "shell-workload-list",
    ):
        if phrase not in text:
            errors.append(f"KHR_AUTO_CONFIGURATION_PLAN missing phrase: {phrase}")

fixture = root / "examples/khr/karl2-baremetal-auto-configuration-plan.json"
if not fixture.is_file():
    errors.append("missing karl2-baremetal-auto-configuration-plan.json")
else:
    doc = json.loads(fixture.read_text(encoding="utf-8"))
    for key, val in [
        ("planOnly", True),
        ("globalAutoEnable", False),
        ("globalDefaultsChanged", False),
        ("runtimeMutation", False),
        ("rolloutInSprint", False),
        ("firstAutoConfiguredModule", "hyperdensity"),
        ("dashboardProviderProfile", "khr-native"),
    ]:
        if doc.get(key) != val:
            errors.append(f"fixture {key}={doc.get(key)!r} expected {val!r}")
    order = doc.get("bootstrapOrder") or []
    expected = [
        "crd-foundation",
        "host-runtime-preview",
        "resourceport-loop",
        "resourcelease-dryrun",
        "guarded-apply-preflight",
        "guarded-apply",
        "guarded-apply-policy",
        "governance",
    ]
    if order != expected:
        errors.append(f"fixture bootstrapOrder mismatch: {order}")

if errors:
    for e in errors:
        print(f"[validate_khr_auto_configuration_plan] FAIL: {e}", file=sys.stderr)
    raise SystemExit(1)
print("[validate_khr_auto_configuration_plan] PASS")
PY

if [[ -x "${ROOT}/scripts/guard_khr_docs_scope.sh" ]]; then
  echo "[validate_khr_auto_configuration_plan] Doc scope guard..."
  "${ROOT}/scripts/guard_khr_docs_scope.sh"
fi

echo "[validate_khr_auto_configuration_plan] OK"
