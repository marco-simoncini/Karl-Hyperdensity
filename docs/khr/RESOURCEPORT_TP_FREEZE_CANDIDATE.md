# ResourcePort TP Freeze Candidate (KHR-AF)

| Field | Value |
|-------|-------|
| **Sprint** | KHR-AF |
| **Contract** | `runtime.karl.io/v1alpha1` / `ResourcePort` |
| **Verdict** | **APPROVED — TP freeze candidate** (capability truth + observation binding) |
| **Cluster apply on ISO** | **Disabled** — preview/disabled |

**NOT production ready.** Native-live capability truth is consumed by native-live lane evidence — see `NATIVE_LIVE_TP_FREEZE_CANDIDATE.md`.

---

## Freeze candidate decision

| Layer | TP freeze candidate? | Notes |
|-------|----------------------|-------|
| CRD (`resourceports.runtime.karl.io`) | **Yes** | Additive-only |
| `RESOURCEPORT_CONTRACT.md` | **Yes** | Capability + hotplug vocabulary |
| Reporting loop (JSON default) | **Yes** | Read-only observation |
| CR preview / cluster apply | **No** | Sandbox-only; ISO disabled |
| Inventory as fact source | **Beta** | Candidacy documented |

Validation: `./scripts/khr_native_live_freeze_check.sh` (shared native-live bundle checks)

---

## Current fields (frozen)

| Field | Freeze |
|-------|--------|
| `spec.provider` | **Yes** — `khr.native`, `kubevirt.compatibility`, … |
| `spec.shellRef` / `spec.cellRef` | **Yes** |
| `spec.capabilities[]` | **Yes** — token vocabulary |
| `spec.hotplug` | **Yes** — cpu/memory/disk/network booleans |
| `spec.constraints` | **Yes** — open schema shape |
| `spec.ports` | **Yes** — legacy matrix for ResourceLease dry-run |
| `status.observedAt` | **Yes** |
| `status.conditions[]` | **Yes** |

---

## Native-live binding semantics

| Signal | ResourcePort meaning |
|--------|---------------------|
| Provider `khr.native` | Native cgroup path |
| Hotplug cpu/memory `true` | Live-scale capable in sandbox |
| Label `karl.io/sandbox-namespace` | Sandbox namespace binding |
| CR name prefix `khr-` | Cluster-scoped preview CRs |

Evidence: certification runs include `cr-preview/resourceport-*-native-live-*-port.json`.

---

## Experimental / unsupported

| Item | Status |
|------|--------|
| Cluster apply from ISO | **Unsupported** |
| Autonomous ResourcePort reconcile | **Unsupported** |
| Production namespace CR write | **Unsupported** |
| Windows hotplug on ResourcePort | **Experimental** observation |

---

## Related

- `RESOURCEPORT_CONTRACT.md`, `RESOURCEPORT_REPORTING_LOOP.md`, `RESOURCEPORT_CR_PREVIEW.md`
- `NATIVE_LIVE_TP_FREEZE_CANDIDATE.md`
