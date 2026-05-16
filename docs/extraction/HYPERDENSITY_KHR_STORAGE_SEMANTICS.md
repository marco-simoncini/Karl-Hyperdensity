# Hyperdensity / KHR — storage semantics (canonical, Sprint 88)

## Summary

Canonizes **KARL-native storage** primitives and **KubeVirt compatibility mapping**. Documentation-only Sprint 88 — no runtime implementation.

---

## Primitives

| Primitive | Purpose |
|-----------|---------|
| **EphemeralDisk** | Ephemeral storage unit attached to Cell/Shell |
| **DiskProfile** | Reusable disk template (size, mode, source, discard) |
| **DiskAttachment** | Binding of disk profile to Shell/Cell instance |

---

## Modes

| Mode | Semantics |
|------|-----------|
| **ephemeralOverlay** | Copy-on-write overlay atop golden/base image |
| **ephemeralClone** | Fast clone from snapshot/volume |
| **scratch** | Temporary workspace; typically deleteOnStop |
| **readonly** | Read-only attachment |
| **persistent** | Survives stop/restart; user or pool scope |

---

## Sources

`pvc`, `image`, `snapshot`, `volume`, `goldenImage`

---

## Discard policies

| Policy | Behavior |
|--------|----------|
| **deleteOnStop** | Remove ephemeral/scratch on stop |
| **keepOnFailure** | Retain for forensics on failure |
| **promoteOnRequest** | Eligible for golden image promotion workflow |

---

## Required capability

**promote-to-image** — maintenance path for golden image refresh (Windows DaaS OS disk lineage).

---

## KubeVirt compatibility mapping

| KubeVirt / Dashboard | KHR mapping |
|----------------------|-------------|
| Ephemeral PVC on dataVolume/volume | `EphemeralDisk` mode=`ephemeralOverlay`, source=`pvc` |
| PVC claimName backing volume | source=`pvc`, claim reference preserved in compatibility layer |
| Container disk / cloud-init disk | Map via `DiskProfile` + provider-specific adapter |

### User-provided compatibility reference

Pattern from user environment (not required in Hyperdensity clone):

- Path: `vm/karl_instances/5.winsrv_vm.yaml`
- OS disk: volume `ephemeral` backed by `persistentVolumeClaim.claimName=karl-os-nfs`

If this path is absent locally, treat as **user-provided compatibility reference** — validation must not fail.

---

## Windows DaaS target profile

| Disk | Configuration |
|------|----------------|
| **OS disk** | source=`goldenImage`, mode=`ephemeralOverlay` |
| **Profile disk** | mode=`persistent`, scope=`user` (per-user profile) |
| **Scratch disk** | mode=`scratch`, discardPolicy=`deleteOnStop` |
| **Golden maintenance** | optional `promote-to-image` on controlled promotion windows |

---

## Related

- `HYPERDENSITY_KHR_ARCHITECTURE_MEMORY.md`
- `HYPERDENSITY_KHR_RESOURCELEASE_DIRECTION.md`


---

## Sprint 89 (ResourceLease minimal contract)

Sprint 89 adds ResourceLease minimal contract sketch (storage/network/provider/examples). No CRD, no controller, no runtime. See HYPERDENSITY_KHR_RESOURCELEASE_MINIMAL_CONTRACT.md and related Sprint 89 docs.
