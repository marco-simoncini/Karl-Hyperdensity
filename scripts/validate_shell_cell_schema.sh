#!/usr/bin/env bash
# Sprint KHR-D: validate Shell/Cell/ShellClass contract examples.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
export ROOT

python3 - <<'PY'
import json
import os
from pathlib import Path

root = Path(os.environ["ROOT"])
examples = root / "docs/contracts/khr/examples"
errors = []


def load(name):
    with (examples / name).open(encoding="utf-8") as f:
        return json.load(f)


def require(cond, msg):
    if not cond:
        errors.append(msg)


shell = load("shell-linux-dev.json")
cell = load("cell-linux-container.json")
sc = load("shellclass-linux-dev.json")

for doc, kind in ((shell, "Shell"), (cell, "Cell"), (sc, "ShellClass")):
    require(doc.get("kind") == kind, f"{kind} kind")
    require(doc.get("apiVersion") == "runtime.karl.io/v1alpha1", f"{kind} apiVersion")

require(shell["spec"]["providerBinding"], "shell providerBinding")
require(shell["spec"]["runtimeClass"], "shell runtimeClass")
require(shell["spec"]["owner"]["tenant"], "shell owner.tenant")
require(shell["status"]["phase"], "shell status.phase")
require(shell["status"]["conditions"], "shell conditions")

require(cell["spec"]["shellRef"]["name"], "cell shellRef")
require(cell["spec"]["providerBinding"], "cell providerBinding")
require(cell["status"]["observedResourcePorts"], "cell observedResourcePorts")

require(sc["spec"]["defaultProviderBinding"], "shellclass defaultProviderBinding")
require(sc["spec"]["defaultRuntimeClass"], "shellclass defaultRuntimeClass")

if errors:
    for e in errors:
        print(f"[validate_shell_cell_schema] FAIL: {e}")
    raise SystemExit(1)
print("[validate_shell_cell_schema] PASS")
PY
