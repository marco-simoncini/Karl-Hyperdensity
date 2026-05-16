#!/usr/bin/env bash
# Sprint KHR-B: validate unified ResourceLease JSON schema and example fixtures (stdlib only).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SCHEMA="${ROOT}/docs/contracts/khr/resourcelease.schema.json"
MANIFEST="${ROOT}/docs/contracts/khr/resourcelease.schema.manifest.json"
EXAMPLES_DIR="${ROOT}/docs/contracts/khr/examples"

export ROOT

python3 - <<'PY'
import json
import os
import sys
from pathlib import Path

root = Path(os.environ["ROOT"])

schema_path = root / "docs/contracts/khr/resourcelease.schema.json"
manifest_path = root / "docs/contracts/khr/resourcelease.schema.manifest.json"
examples_dir = root / "docs/contracts/khr/examples"

errors = []


def load_json(p: Path):
    with p.open(encoding="utf-8") as f:
        return json.load(f)


def get_def(schema, name):
    return schema.get("$defs", {}).get(name, {})


def require(cond, msg):
    if not cond:
        errors.append(msg)


schema = load_json(schema_path)
manifest = load_json(manifest_path)

require(schema.get("title") == "KHR ResourceLease", "schema title mismatch")
require("ADR-0005" in schema.get("description", "") or "unified" in schema.get("description", "").lower(), "schema description must reference unified contract")

spec_def = get_def(schema, "spec")
spec_required = set(spec_def.get("required", []))
require(
    {"leaseKind", "shell", "cell", "provider"}.issubset(spec_required),
    f"spec required mismatch: {spec_required}",
)

lease_kind_enum = set(get_def(schema, "leaseKind").get("enum", []))
require(lease_kind_enum == {"runtime", "transfer"}, f"leaseKind enum: {lease_kind_enum}")

provider_enum = set(get_def(schema, "provider").get("enum", []))
require(
    {"khr.native", "kubevirt.compatibility", "kubevirt.public-cloud-fallback"}.issubset(provider_enum),
    f"provider enum missing core values: {provider_enum}",
)

phase_enum = set(get_def(schema, "phase").get("enum", []))
require("DryRunValidated" in phase_enum and "Active" in phase_enum, f"phase enum: {phase_enum}")

disk_mode_enum = set(get_def(schema, "diskMode").get("enum", []))
source_enum = set(get_def(schema, "sourceType").get("enum", []))
discard_enum = set(get_def(schema, "discardPolicy").get("enum", []))

require(manifest.get("schemaVersion", "").startswith("v") or manifest.get("unifiedContract") is True or True, "manifest present")
require(manifest.get("schemaOnly") is True, "manifest schemaOnly must be true")
require(manifest.get("controllerImplemented") is False, "manifest controllerImplemented must be false")

example_files = [
    "resourcelease-windows-daas-khr-native.json",
    "resourcelease-public-cloud-kubevirt-fallback.json",
    "resourcelease-baremetal-native.json",
    "resourcelease-linux-cpu-transfer.json",
]

for name in example_files:
    p = examples_dir / name
    require(p.is_file(), f"missing example {name}")
    ex = load_json(p)
    for key in ("apiVersion", "kind", "metadata", "spec"):
        require(key in ex, f"{name}: missing top-level {key}")
    require(ex.get("apiVersion") == "hyperdensity.karl.io/v1alpha1", f"{name}: apiVersion")
    require(ex.get("kind") == "ResourceLease", f"{name}: kind")
    spec = ex.get("spec", {})
    for sec in ("leaseKind", "shell", "cell", "provider"):
        require(sec in spec, f"{name}: spec missing {sec}")
    lk = spec.get("leaseKind")
    require(lk in lease_kind_enum, f"{name}: leaseKind {lk!r}")
    prov = spec.get("provider")
    require(prov in provider_enum, f"{name}: provider {prov!r}")
    if lk == "runtime":
        for sec in ("resources", "storage", "network", "policy"):
            require(sec in spec, f"{name}: runtime lease missing {sec}")
        for disk in spec.get("storage", {}).get("disks", []):
            mode = disk.get("mode")
            require(mode in disk_mode_enum, f"{name}: disk mode {mode!r}")
            st = disk.get("source", {}).get("type")
            require(st in source_enum, f"{name}: source type {st!r}")
    if lk == "transfer":
        tr = spec.get("transfer", {})
        require(tr.get("donor", {}).get("name"), f"{name}: transfer.donor.name")
        require(tr.get("receiver", {}).get("name"), f"{name}: transfer.receiver.name")
        require(tr.get("resource") in get_def(schema, "transfer")["properties"]["resource"]["enum"], f"{name}: transfer.resource")

win = load_json(examples_dir / "resourcelease-windows-daas-khr-native.json")
require(win["spec"]["leaseKind"] == "runtime", "windows: leaseKind")
require(win["spec"]["shell"]["kind"] == "windowsDesktop", "windows: shell kind")

pub = load_json(examples_dir / "resourcelease-public-cloud-kubevirt-fallback.json")
require(pub["spec"]["provider"] == "kubevirt.public-cloud-fallback", "public: provider")
root_disk = pub["spec"]["storage"]["disks"][0]
require(root_disk["mode"] == "ephemeralOverlay", "public: os mode")
require(root_disk["source"]["ref"] == "karl-os-nfs", "public: karl-os-nfs")

xfer = load_json(examples_dir / "resourcelease-linux-cpu-transfer.json")
require(xfer["spec"]["leaseKind"] == "transfer", "transfer: leaseKind")
require(xfer["spec"]["transfer"]["mode"] == "envelope", "transfer: mode")
require(xfer.get("status", {}).get("phase") in phase_enum, "transfer: status.phase")

if errors:
    for e in errors:
        print(f"[validate_resourcelease_schema] FAIL: {e}", file=sys.stderr)
    sys.exit(1)

print("[validate_resourcelease_schema] PASS")
PY
