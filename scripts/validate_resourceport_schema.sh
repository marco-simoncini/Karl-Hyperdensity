#!/usr/bin/env bash
# Sprint KHR-C: validate ResourcePort observation JSON schema and examples.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
export ROOT

python3 - <<'PY'
import json
import os
from pathlib import Path

root = Path(os.environ["ROOT"])
schema_path = root / "docs/contracts/khr/resourceport.schema.json"
example_path = root / "docs/contracts/khr/examples/resourceport-linux-container.json"
errors = []


def load(p):
    with p.open(encoding="utf-8") as f:
        return json.load(f)


def require(cond, msg):
    if not cond:
        errors.append(msg)


schema = load(schema_path)
ex = load(example_path)

require(schema.get("title") == "KHR ResourcePort observation", "schema title")
spec_req = set(schema["$defs"]["spec"]["required"])
require({"provider", "shellRef", "cellRef", "capabilities", "hotplug"}.issubset(spec_req), f"spec required: {spec_req}")

hotplug_req = set(schema["$defs"]["hotplug"]["required"])
require(hotplug_req == {"cpu", "memory", "disk", "network"}, f"hotplug: {hotplug_req}")

provider_enum = set(schema["$defs"]["provider"]["enum"])
require("kubernetes.cni" in provider_enum, "provider enum")

require(ex["spec"]["provider"] == "kubernetes.cni", "example provider")
require(ex["spec"]["hotplug"]["cpu"] is False, "example hotplug.cpu")
require("cpu.static" in ex["spec"]["capabilities"], "example capabilities")
require(ex["status"]["observedAt"], "example observedAt")
require(ex["status"]["conditions"][0]["type"] == "Observed", "example condition")

if errors:
    for e in errors:
        print(f"[validate_resourceport_schema] FAIL: {e}")
    raise SystemExit(1)
print("[validate_resourceport_schema] PASS")
PY
