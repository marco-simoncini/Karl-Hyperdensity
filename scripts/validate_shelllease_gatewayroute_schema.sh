#!/usr/bin/env bash
# Sprint KHR-E: validate ShellLease + GatewayRoute contract examples.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
export ROOT

python3 - <<'PY'
import json
from pathlib import Path

root = Path(__import__("os").environ["ROOT"])
examples = root / "docs/contracts/khr/examples"
errors = []


def load(name):
    with (examples / name).open(encoding="utf-8") as f:
        return json.load(f)


def require(cond, msg):
    if not cond:
        errors.append(msg)


lease = load("shelllease-demo.json")
rdp = load("gatewayroute-rdp-demo.json")
app = load("gatewayroute-remoteapp-demo.json")

require(lease["kind"] == "ShellLease", "shelllease kind")
require(lease["spec"]["leaseMode"] in ("ephemeral", "persistent", "scheduled"), "leaseMode")
require(lease["spec"]["shellRef"]["name"], "shellRef")
require(lease["status"]["phase"], "shelllease phase")

for gr, proto in ((rdp, "rdp"), (app, "rdp-remoteapp")):
    require(gr["apiVersion"] == "gateway.karl.io/v1alpha1", f"{proto} apiVersion")
    require(gr["spec"]["protocol"] == proto, f"{proto} protocol")
    require(gr["spec"]["shellLeaseRef"]["name"], f"{proto} shellLeaseRef")
    require(gr["status"]["conditions"], f"{proto} conditions")

if errors:
    for e in errors:
        print(f"[validate_shelllease_gatewayroute_schema] FAIL: {e}")
    raise SystemExit(1)
print("[validate_shelllease_gatewayroute_schema] PASS")
PY
