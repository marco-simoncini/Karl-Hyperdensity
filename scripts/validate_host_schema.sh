#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
export ROOT
python3 - <<'PY'
import json
import os
from pathlib import Path

root = Path(os.environ["ROOT"])
path = root / "docs/contracts/khr/examples/host-karl-metal-01.json"
doc = json.loads(path.read_text(encoding="utf-8"))
assert doc["kind"] == "Host"
assert doc["apiVersion"] == "runtime.karl.io/v1alpha1"
assert doc["spec"]["hostId"]
assert doc["spec"]["nodeName"] == "karl-metal-01"
assert doc["spec"]["runtimeMode"] == "sandbox"
assert doc["status"]["safetyMode"] == "sandbox"
assert doc["status"]["runtimeVersion"]
print("[validate_host_schema] PASS")
PY
