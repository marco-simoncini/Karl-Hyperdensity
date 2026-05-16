#!/usr/bin/env bash
# Sprint 91: validate KHR ResourceLease JSON schema and example fixtures (stdlib only).
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


def enum_values(schema, def_name, prop_chain):
    """Walk $ref chain to find enum at prop_chain e.g. ('spec','provider') or disk mode."""
    cur = schema
    for part in prop_chain:
        if "$ref" in cur:
            ref = cur["$ref"]
            if ref.startswith("#/$defs/"):
                cur = get_def(schema, ref.split("/")[-1])
        props = cur.get("properties", {})
        if part not in props:
            return []
        cur = props[part]
    if "$ref" in cur:
        ref = cur["$ref"]
        if ref.startswith("#/$defs/"):
            cur = get_def(schema, ref.split("/")[-1])
    return cur.get("enum", [])


def require(cond, msg):
    if not cond:
        errors.append(msg)


schema = load_json(schema_path)
manifest = load_json(manifest_path)

require(schema.get("title") == "KHR ResourceLease", "schema title mismatch")
require("apiVersion" in schema.get("properties", {}), "schema missing apiVersion")
require("kind" in schema.get("properties", {}), "schema missing kind")
require("metadata" in schema.get("properties", {}), "schema missing metadata")
require("spec" in schema.get("properties", {}), "schema missing spec")
require("status" in schema.get("properties", {}), "schema missing optional status")

spec_def = get_def(schema, "spec")
spec_required = set(spec_def.get("required", []))
require(
    {"shell", "cell", "provider", "resources", "storage", "network", "policy"}.issubset(spec_required),
    f"spec required mismatch: {spec_required}",
)

provider_enum = set(get_def(schema, "provider").get("enum", []))
require(
    {"khr.native", "kubevirt.compatibility", "kubevirt.public-cloud-fallback"}.issubset(provider_enum),
    f"provider enum missing core values: {provider_enum}",
)

disk_mode_enum = set(get_def(schema, "diskMode").get("enum", []))
require(
    {"ephemeralOverlay", "ephemeralClone", "scratch", "readonly", "persistent"}.issubset(disk_mode_enum),
    f"diskMode enum mismatch: {disk_mode_enum}",
)

source_enum = set(get_def(schema, "sourceType").get("enum", []))
require(
    {"pvc", "image", "snapshot", "volume", "goldenImage"}.issubset(source_enum),
    f"sourceType enum mismatch: {source_enum}",
)

discard_enum = set(get_def(schema, "discardPolicy").get("enum", []))
require(
    {"deleteOnStop", "keepOnFailure", "promoteOnRequest"}.issubset(discard_enum),
    f"discardPolicy enum mismatch: {discard_enum}",
)

require(manifest.get("schemaOnly") is True, "manifest schemaOnly must be true")
require(manifest.get("crdApplied") is False, "manifest crdApplied must be false")

example_files = [
    "resourcelease-windows-daas-khr-native.json",
    "resourcelease-public-cloud-kubevirt-fallback.json",
    "resourcelease-baremetal-native.json",
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
    for sec in ("shell", "cell", "provider", "resources", "storage", "network", "policy"):
        require(sec in spec, f"{name}: spec missing {sec}")
    prov = spec.get("provider")
    require(prov in provider_enum, f"{name}: provider {prov!r} not in enum")
    for disk in spec.get("storage", {}).get("disks", []):
        mode = disk.get("mode")
        require(mode in disk_mode_enum, f"{name}: disk mode {mode!r}")
        st = disk.get("source", {}).get("type")
        require(st in source_enum, f"{name}: source type {st!r}")
        dp = disk.get("discardPolicy")
        if dp is not None:
            require(dp in discard_enum, f"{name}: discardPolicy {dp!r}")

# Windows DaaS checks
win = load_json(examples_dir / "resourcelease-windows-daas-khr-native.json")
require(win["spec"]["shell"]["kind"] == "windowsDesktop", "windows: shell kind")
require(win["spec"]["provider"] == "khr.native", "windows: provider")
require("rdp-GW" in win["spec"]["network"]["exposure"]["ingress"], "windows: ingress")
require(win["spec"]["policy"].get("tenantIsolation") == "strict", "windows: tenantIsolation")

# Public cloud fallback
pub = load_json(examples_dir / "resourcelease-public-cloud-kubevirt-fallback.json")
require(pub["spec"]["provider"] == "kubevirt.public-cloud-fallback", "public: provider")
root_disk = pub["spec"]["storage"]["disks"][0]
require(root_disk["mode"] == "ephemeralOverlay", "public: os mode")
require(root_disk["source"]["type"] == "pvc", "public: source type")
require(root_disk["source"]["ref"] == "karl-os-nfs", "public: karl-os-nfs")
require(pub["spec"]["network"]["providerNetwork"]["provider"] == "kubevirt.legacy.ovn", "public: ovn provider")

# Baremetal
bm = load_json(examples_dir / "resourcelease-baremetal-native.json")
require(bm["spec"]["provider"] == "khr.native", "baremetal: provider")
bm_net = bm["spec"]["network"]["providerNetwork"]["provider"]
require(bm_net in ("baremetal.bridge", "baremetal.vlan"), f"baremetal: network provider {bm_net}")

if errors:
    for e in errors:
        print(f"[validate_resourcelease_schema] FAIL: {e}", file=sys.stderr)
    sys.exit(1)

print("[validate_resourcelease_schema] PASS")
PY
