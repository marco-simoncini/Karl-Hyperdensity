# Hyperdensity / KHR — KubeVirt capability inventory (Sprint 90)

## Summary

Read-only inventory of KubeVirt-related capabilities in local clones, mapped toward ResourceLease storage/provider (Sprint 89).

---

## 1. Scope

Inventory only — no runtime, no CRD apply, no controller.

---

## 2. Repositories scanned

Karl-Hyperdensity, Karl-Dashboard, FluidVirt (read-only), audit-work/KARL (read-only), kubevirt-upstream (reference).

---

## 3. Search terms

See `HYPERDENSITY_KHR_INVENTORY_SCAN_METHOD.md`.

---

## 4. KubeVirt objects found

### Karl-Hyperdensity

| Path | Kind / concept |
|------|----------------|
| `examples/providers/kubevirt/shell-windows-desktop-kubevirt-legacy.yaml` | Shell + ShellClass, `kubevirt.legacy.v1` provider |
| `examples/providers/kubevirt/shell-linux-vm-kubevirt-legacy.yaml` | Linux VM Shell example |
| `examples/providers/kubevirt/cell-kubevirt-vm-handle.yaml` | Cell handle |
| `examples/providers/kubevirt/resourcelease-kubevirt-guarded-example.yaml` | ResourceLease guarded example |
| `api/providers/kubevirt/kubevirt-legacy-provider.yaml` | Provider contract sketch |
| `docs/providers/KUBEVIRT_LEGACY_PROVIDER_CONTRACT.md` | Provider documentation |
| `docs/adr/ADR-0002-kubevirt-as-legacy-provider.md` | ADR: KubeVirt as legacy provider |

### Karl-Dashboard (runtime — observation only, not modified)

| Path | Capability |
|------|------------|
| `pkg/server/hyperdensity_parent_fabric_live.go` | Live discovery: `/apis/kubevirt.io/v1/virtualmachines`, `virtualmachineinstances`, `virtualmachineinstancemigrations` |
| `pkg/server/hyperdensity_parent_fabric_workload_adapter_v1.go` | `VirtualMachinePath`, `VirtualMachineInstancePath` API path helpers |
| `pkg/server/hyperdensity_parent_fabric_vm_linux_*.go` | VM Linux CPU/memory runtime (legacy surface) |
| `pkg/server/hyperdensity_parent_fabric_windows_*.go` | Windows pool / VMPool governance references |
| `pkg/server/hyperdensity_parent_fabric_unified_cross_os_runtime_engine_v1.go` | `VirtualMachinePool` object IDs |
| `artifacts/**/win11_pool*.yaml`, `vms_*.yaml` | Captured VMPool / VM cluster state |

### FluidVirt (read-only)

| Path | Capability |
|------|------------|
| `pkg/virt-controller/services/template.go` | Multus `NetworkToResource` for VMI templates |

---

## 5. Ephemeral PVC findings

### User-provided compatibility reference

| Field | Value |
|-------|-------|
| **path** | `vm/karl_instances/5.winsrv_vm.yaml` |
| **status in Karl-Hyperdensity / Karl-Dashboard** | **not found in local clone** |
| **status in audit-work mirror** | **found** at `audit-work/repos/Karl-technology/KARL/vm/karl_instances/5.winsrv_vm.yaml` |

**Semantic** (verified in audit-work copy):

```yaml
volumes:
- name: os
  ephemeral:
    persistentVolumeClaim:
      claimName: karl-os-nfs
```

Disk name: `os`. Kind: `VirtualMachinePool` (`pool.kubevirt.io/v1alpha1`).

**ResourceLease mapping**: `storage.disks[]` with `role=os`, `mode=ephemeralOverlay`, `source.type=pvc`, `sourceRef=karl-os-nfs`.

**Validation rule**: absence in primary clone **does not** fail Sprint 90 — requirement preserved.

### Other ephemeral / PVC evidence (Dashboard artifacts)

- `artifacts/the-father-vm-linux-cpu-live-bidirectional-breakthrough/**` — `ephemeral` volume + `claimName` on launcher pods
- Multiple VM YAML captures with `persistentVolumeClaim` / `claimName` (RWO roots)

Hyperdensity **examples** do not embed raw KubeVirt ephemeral PVC YAML; mapping is **documental** per Sprint 89 storage contract.

---

## 6. VM/VMPool findings

| Finding | Location |
|---------|----------|
| VirtualMachinePool | audit-work `5.winsrv_vm.yaml`; Dashboard artifacts `pool.kubevirt.io/v1beta1` |
| VirtualMachine / VMI discovery | `hyperdensity_parent_fabric_live.go` |
| VMPool governance IDs | `hyperdensity_parent_fabric_windows_rollback_and_pool_governance_v1.go` |
| Pool member naming | `hyperdensity_auto_scope_policy.go` (e.g. `win11-pool-0`) |

**ResourceLease provider mapping**: `kubevirt.compatibility` for VM/VMPool; `kubevirt.public-cloud-fallback` when KHR unavailable.

---

## 7. Windows/DaaS findings

| Item | Evidence |
|------|----------|
| Windows Shell on KubeVirt | `shell-windows-desktop-kubevirt-legacy.yaml` |
| Windows pool runtime | Dashboard windows_* + unified_cross_os_runtime_engine |
| OS disk ephemeral PVC pattern | user reference `5.winsrv_vm.yaml` (audit-work) |
| goldenImage semantics | Sprint 88–89 docs + Shell examples (not raw goldenImage CR in Hyperdensity examples) |

---

## 8. ResourceLease mapping

| KubeVirt / source | ResourceLease |
|-------------------|---------------|
| Ephemeral PVC volume | `storage.disks[].mode=ephemeralOverlay`, `source.type=pvc` |
| PVC claimName | `sourceRef` / `source.ref` |
| containerDisk / image | `source.type=image` |
| snapshot | `source.type=snapshot` |
| golden image lineage | `source.type=goldenImage` |
| VM / VMPool | `provider=kubevirt.compatibility` or `kubevirt.public-cloud-fallback` |

---

## 9. Gaps

- `vm/karl_instances/5.winsrv_vm.yaml` **not** in Karl-Hyperdensity or Karl-Dashboard tree (only audit-work mirror).
- No applied ResourceLease CRD reconciling KubeVirt.
- Hyperdensity examples use KARL `Shell`/`Cell` CRs, not raw KubeVirt VM YAML in-repo.

---

## 10. Risks

- Treating audit-work path as guaranteed in all clones.
- Losing ephemeral PVC semantic when promoting to KARL-native EphemeralDisk.

---

## 11. Recommended next sprint

**Sprint 91** — JSON Schema file for ResourceLease under `docs/contracts/` or contractkit fixture manifest entry; still no cluster apply. Optional: copy user reference into `examples/compatibility/` as non-applied sample.

---

## Related

- `HYPERDENSITY_KHR_RESOURCELEASE_INVENTORY_MAPPING.md`
- `HYPERDENSITY_KHR_RESOURCELEASE_STORAGE_CONTRACT.md`
