# Hyperdensity / KHR — roadmap transition note (Sprint 87)

## Summary

After **resource_exchange observation boundary closure** (Sprint 87), KHR must **not** continue indefinitely with micro-sprint adapter-only work. The next block consolidates **project memory / architecture** and canonizes **storage** and **network** semantics as KARL-native primitives.

---

## Stop: indefinite micro-sprint adapter pattern

Sprint 78–87 completed a full, gated observation chain for **resource_exchange** only. Further adapter sprints without architecture consolidation risk fragmenting the product model.

---

## Canonize: platform roles

| Role | Direction |
|------|-----------|
| **KARL Engine / KHR** | Native runtime |
| **Shell / Cell** | Core workload units |
| **KubeVirt** | Compatibility provider (not the product model) |
| **Hyperdensity** | Governance / resource intelligence |
| **OVN/SDN (Dashboard)** | Compatibility source-of-truth for network semantics |

---

## Canonize: storage primitives

Must appear in architecture / project memory:

- **EphemeralDisk**, **ephemeralOverlay**, **ephemeralClone**
- **scratch**, **readonly**, **persistent**
- **discardPolicy**: `deleteOnStop`, `keepOnFailure`, `promoteOnRequest`
- **promote-to-image**
- **source** kinds: pvc, image, snapshot, volume, goldenImage

**Windows DaaS storage pattern** (preserve in roadmap):

- OS disk: goldenImage + ephemeralOverlay
- Profile disk: persistent
- Scratch disk: deleteOnStop

Existing KubeVirt/Dashboard capabilities must be **preserved** as provider backends.

---

## Canonize: network primitives

- **KARLNetwork**, **CellNetwork**, **ShellNetwork**
- **NetworkAttachment**, **NetworkLease**, **NetworkPolicy**
- **OVN/SDN compatibility mapping** — Dashboard/OVN semantics remain source-of-truth until KARL-native models supersede them

KubeVirt/OVN become **providers/backends**, not the product definition.

---

## Recommended next phase

`khr_architecture_memory_and_storage_network_semantics` — documentation, contracts sketch, and parity anchors without runtime flip until explicitly scheduled.

---

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_MIGRATION_BOUNDARY.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_REMAINING_SURFACE_DECISION.md`
- `docs/roadmap/KHR_HYPERDENSITY_CORRECTED_ROADMAP.md`


---

## Sprint 88 (KHR architecture memory)

Sprint 88 canonizes KHR/KARL architecture memory, storage semantics, and network/OVN semantics. No runtime/adapter changes. KubeVirt remains compatibility provider and public-cloud fallback. See HYPERDENSITY_KHR_ARCHITECTURE_MEMORY.md and related Sprint 88 docs.
