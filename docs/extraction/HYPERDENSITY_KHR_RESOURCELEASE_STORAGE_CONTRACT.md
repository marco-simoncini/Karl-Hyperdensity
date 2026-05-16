# Hyperdensity / KHR — ResourceLease storage contract (Sprint 89)

## Summary

Details `spec.storage` for ResourceLease minimal contract. Anchored to `HYPERDENSITY_KHR_STORAGE_SEMANTICS.md`.

---

## storage top-level

```yaml
storage:
  disks: []
  defaultDiscardPolicy: deleteOnStop
  promoteToImage: { enabled, goldenRef, policy }
```

---

## disk fields

| Field | Required | Description |
|-------|----------|-------------|
| `name` | yes | Disk identifier within lease |
| `role` | yes | `os`, `profile`, `scratch`, `data`, … |
| `mode` | yes | See allowed modes |
| `source` | yes | `{ type, ref }` |
| `sourceRef` | optional | PVC name, image ID, snapshot ID |
| `size` | optional | Quantity string |
| `readonly` | optional | bool |
| `discardPolicy` | optional | Overrides default |
| `mountIntent` | optional | boot, data, temp |
| `promotePolicy` | optional | promote-to-image eligibility |

---

## Allowed modes

`ephemeralOverlay`, `ephemeralClone`, `scratch`, `readonly`, `persistent`

---

## Allowed source types

`pvc`, `image`, `snapshot`, `volume`, `goldenImage`

---

## Discard policies

`deleteOnStop`, `keepOnFailure`, `promoteOnRequest`

---

## Windows DaaS profile

| Disk | Configuration |
|------|----------------|
| **os** | role=`os`, mode=`ephemeralOverlay`, source.type=`goldenImage`, discardPolicy=`deleteOnStop` or `keepOnFailure` |
| **profile** | role=`profile`, mode=`persistent`, scope=`user` |
| **scratch** | role=`scratch`, mode=`scratch`, discardPolicy=`deleteOnStop` |

---

## KubeVirt compatibility

| KubeVirt / Dashboard | ResourceLease disk |
|----------------------|------------------|
| Ephemeral PVC volume | mode=`ephemeralOverlay`, source.type=`pvc`, sourceRef=claimName |

**User-provided reference** (may be absent in clone): `vm/karl_instances/5.winsrv_vm.yaml` with `claimName=karl-os-nfs`. Validation must **not** fail if file missing.

---

## Related

- `HYPERDENSITY_KHR_RESOURCELEASE_MINIMAL_CONTRACT.md`
- `HYPERDENSITY_KHR_STORAGE_SEMANTICS.md`


---

## Sprint 90 (inventory facts)

Sprint 90 adds read-only KubeVirt and OVN/SDN capability inventory mapped to ResourceLease contract. No CRD, no controller, no runtime. See HYPERDENSITY_KHR_KUBEVIRT_CAPABILITY_INVENTORY.md and related Sprint 90 docs.


---

## Sprint 91 (ResourceLease JSON Schema)

Sprint 91 adds non-applied JSON Schema and example fixtures under docs/contracts/khr/. No CRD, no controller, no runtime. See HYPERDENSITY_KHR_RESOURCELEASE_JSON_SCHEMA.md.
