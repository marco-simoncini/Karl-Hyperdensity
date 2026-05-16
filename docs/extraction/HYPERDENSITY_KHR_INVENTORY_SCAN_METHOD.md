# Hyperdensity / KHR — inventory scan method (Sprint 90)

## Summary

Documents the **read-only** inventory scan performed for Sprint 90. No repository modifications outside Karl-Hyperdensity and Karl-Dashboard. No runtime changes.

---

## Scan method

1. **Workspace root**: `/home/m.simoncini/GitHub`
2. **Tooling**: `rg` (ripgrep) with glob filters `*.{yaml,yml,go,md}`; `glob` for path discovery
3. **Scope**: KubeVirt/VM/storage terms and OVN/SDN/network terms
4. **Output**: Findings recorded in inventory docs; parity goldens in Dashboard

---

## Repositories read (not modified)

| Repository | Role in scan |
|------------|----------------|
| **Karl-Hyperdensity** | Primary — provider examples, KubeVirt API sketches, docs |
| **Karl-Dashboard** | Primary — Parent Fabric live discovery, VM runtime, artifacts |
| **FluidVirt** | Read-only — upstream KubeVirt fork; multus network wiring |
| **audit-work/repos/Karl-technology/KARL** | Read-only — hosts `vm/karl_instances/5.winsrv_vm.yaml` |
| **kubevirt-upstream-v1.8.1** | Read-only — upstream reference (not exhaustively indexed) |
| **Karl-OS-ISO**, **Karl-Installer**, **Karl-Inventory**, **Karl-Warden** | Read-only — no Sprint 90 file changes; spot-checked |

---

## Search terms (KubeVirt)

`VirtualMachine`, `VirtualMachineInstance`, `VirtualMachinePool`, `VMPool`, `kubevirt.io`, `DataVolume`, `PVC`, `persistentVolumeClaim`, `ephemeral`, `containerDisk`, `cloudInitNoCloud`, `claimName`, `karl-os-nfs`, `winsrv`, `goldenImage`, `snapshot`, `pool.kubevirt.io`

---

## Search terms (OVN/SDN)

`OVN`, `ovn`, `SDN`, `sdn`, `NetworkAttachmentDefinition`, `NAD`, `logical switch`, `logical router`, `logical port`, `ACL`, `NAT`, `SNAT`, `DNAT`, `DHCP`, `FloatingIP`, `tenant network`, `gateway`, `multus`, `cni`, `kube-ovn`, `KubeOVN`

---

## Limits of local clone

- **Absence in Karl-Hyperdensity / Karl-Dashboard does not remove user-provided requirements** (Sprint 88–89).
- OVN **controller/schema source** may live outside indexed repos (operational cluster + Dashboard runtime only).
- Artifact YAML under `Karl-Dashboard/artifacts/` reflects **live cluster captures**, not product CRDs.

---

## Rules

- No edits to read-only repos
- No CRD apply, no controller, no API/runtime change
- Inventory anchors future schema/contracts only

---

## Related

- `HYPERDENSITY_KHR_KUBEVIRT_CAPABILITY_INVENTORY.md`
- `HYPERDENSITY_KHR_OVN_SDN_CAPABILITY_INVENTORY.md`
- `HYPERDENSITY_KHR_RESOURCELEASE_INVENTORY_MAPPING.md`
