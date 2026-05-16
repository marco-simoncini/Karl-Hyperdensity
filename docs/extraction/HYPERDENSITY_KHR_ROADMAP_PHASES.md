# Hyperdensity / KHR — roadmap phases (post Sprint 87)

## Summary

Updated phase model after resource_exchange boundary closure (Sprint 87) and architecture memory canonization (Sprint 88).

---

## Phase 0 — Parent Fabric safety foundation (complete)

- Apply observation track complete (Sprint 75–77).
- Resource_exchange observation track complete (Sprint 78–87).
- `ObservationWiredV1` / `ProductionWiredV1` remain **false** (broad observation not authorized).
- rollback, VM runtime, admission_guard remain **legacy** separate surfaces.

---

## Phase 1 — KHR architecture memory (Sprint 88, current)

- `HYPERDENSITY_KHR_ARCHITECTURE_MEMORY.md` (canonical)
- Storage semantics (`HYPERDENSITY_KHR_STORAGE_SEMANTICS.md`)
- Network/SDN semantics (`HYPERDENSITY_KHR_NETWORK_SDN_SEMANTICS.md`)
- ResourceLease direction
- Parity goldens in Dashboard — **no runtime change**

---

## Phase 2 — ResourceLease minimal contract

- CRD/schema sketch for ResourceLease dimensions
- Storage + network sections from Phase 1
- Provider enum including KubeVirt public-cloud fallback

---

## Phase 3 — Inventory facts / capacity facts

- Host/cell inventory
- OVN/SDN inventory task (if not in clone)
- KubeVirt capability inventory

---

## Phase 4 — Hyperdensity dry-run decision loop

- Donor/receiver dry-run
- Capacity prediction integration
- No production auto-apply without gates

---

## Phase 5 — KubeVirt-backed operational loop

- VM/VMPool operations via compatibility provider
- Ephemeral PVC preservation tests
- Public cloud fallback paths validated

---

## Phase 6 — KHR native host runtime seed

- Read-only discovery, cgroup envelope (existing KHR docs)
- No broad observation flip

---

## Phase 7 — Shell/Cell provider model

- Provider selection in lease
- Shell experience mapping

---

## Phase 8 — Windows DaaS KHR-native storage/network profile

- goldenImage + ephemeralOverlay OS disk
- persistent profile disk
- scratch deleteOnStop
- tenant-isolated network, rdp-GW ingress

---

## Related

- `HYPERDENSITY_KHR_ROADMAP_TRANSITION_NOTE.md`
- `HYPERDENSITY_PARENT_FABRIC_EXTRACTION_PHASES.md`
- `docs/roadmap/KHR_HYPERDENSITY_CORRECTED_ROADMAP.md`


---

## Sprint 89 (ResourceLease minimal contract)

Sprint 89 adds ResourceLease minimal contract sketch (storage/network/provider/examples). No CRD, no controller, no runtime. See HYPERDENSITY_KHR_RESOURCELEASE_MINIMAL_CONTRACT.md and related Sprint 89 docs.
