#!/usr/bin/env bash
# KHR-P: validate Windows live scale lane contract schemas and examples.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
export ROOT

python3 - <<'PY'
import json
import os
from pathlib import Path

root = Path(os.environ["ROOT"])
errors = []


def load(p):
    with p.open(encoding="utf-8") as f:
        return json.load(f)


def require(cond, msg):
    if not cond:
        errors.append(msg)


port_schema = load(root / "docs/contracts/khr/windows-live-scale-resourceport.schema.json")
blocked_schema = load(root / "docs/contracts/khr/windows-live-scale-blocked-state.schema.json")
port_ex = load(root / "examples/khr/windows/resourceport-windows-session.json")
cpu_ex = load(root / "examples/khr/windows/resourcelease-windows-cpu-dryrun-observed.json")
mem_ex = load(root / "examples/khr/windows/resourcelease-windows-memory-dryrun-blocked.json")

wl = port_ex["spec"]["windowsLiveScale"]
require(wl["requiresRestart"] is False, "windows.host-runtime requiresRestart target false")
require(port_ex["spec"]["providerBinding"] == "windows.host-runtime", "port providerBinding")
require(wl["observationOnly"] is True, "port observationOnly")

obs = cpu_ex["status"]["dryRunObservation"]
require(obs["allowed"] is True and obs["noApply"] is True, "cpu dryrun observed")

mem = mem_ex["status"]["dryRunObservation"]
require(mem["blocked"] is True, "memory dryrun blocked")
require(mem["blockedState"] == "requiresRestart", "memory blockedState")
require(mem["liveScaleTarget"] == "compatibility-fallback", "memory fallback target")

enum_states = set(blocked_schema["properties"]["blockedState"]["enum"])
require(enum_states == {"requiresRestart", "requiresReboot", "requiresSessionDrain", "providerUnsupported"}, "blocked enum")

providers = set(port_schema["$defs"]["providerBinding"]["enum"])
require(providers == {"windows.host-runtime", "kubevirt.compatibility"}, "provider enum")

if errors:
    for e in errors:
        print(f"[validate_windows_live_scale_contract] FAIL: {e}")
    raise SystemExit(1)
print("[validate_windows_live_scale_contract] PASS")
PY

go test ./pkg/khr/windowslane/... -count=1
