# Hyperdensity / KHR — architecture memory (canonical, Sprint 88)

## Summary

**Sprint 88** canonizes the **KARL Engine / KHR** roadmap as verifiable architecture memory. No runtime changes, no adapter wiring, no flag flips. This document is the **canonical** reference for product direction after Parent Fabric apply/resource_exchange boundary closure (Sprint 87).

---

## 1. Scope

| In scope | Out of scope (Sprint 88) |
|----------|---------------------------|
| Platform roles, Shell/Cell model | Runtime code changes |
| KubeVirt / public-cloud fallback rules | Broad observation |
| Storage & network semantics (references) | rollback/VM/admission wiring |
| Hyperdensity governance role | KHR agent, cluster, ISO |
| ResourceLease direction (reference) | New operational providers |

---

## 2. Product direction

**KARL Engine / KHR** is the **native host runtime** future. Users interact with **Shell** experiences (VM-like, desktop, app, container, session). Workloads materialize as **Cell** units on hosts. **Hyperdensity** governs resource intelligence. **KubeVirt** and **OVN/SDN (Dashboard)** remain **compatibility providers** — not the product model.

---

## 3. Platform roles

| Component | Role |
|-----------|------|
| **KARL Engine / KHR** | Native host runtime |
| **Shell** | User-facing workload experience |
| **Cell** | Executable unit materialized on host |
| **Hyperdensity** | Governance / resource intelligence (Resource Market, ResourceLease, donor/receiver, rollback, capacity prediction) |
| **KubeVirt** | Compatibility provider / fallback provider |
| **OVN/SDN Dashboard** | Compatibility source-of-truth for network semantics |

---

## 4. Shell / Cell model

- **Shell**: what the user sees — VM-like desktop, app session, container shell, etc.
- **Cell**: what runs on the host — cgroup envelope, disks, network attachments, provider binding.
- Product contracts (claims, leases, evidence) attach to **Shell**; execution binds to **Cell**.

---

## 5. KHR provider model

Providers are explicit and selectable per lease:

- `khr.native` — direct KHR host runtime where available
- `kubevirt.compatibility` — KubeVirt VM/VMPool mapping
- `kubevirt.public-cloud-fallback` — when KHR cannot run on public cloud
- `ovn.compatibility` — Dashboard OVN/SDN mapping
- `kubernetes.cni`, `cloud.vpc`, `baremetal.bridge`, `baremetal.vlan` — future/backends

---

## 6. KubeVirt compatibility rule

KubeVirt semantics **must not be lost**:

- Ephemeral PVC patterns
- VM / VMPool compatibility
- Provider fallback when KHR unavailable

KubeVirt is a **backend**, not the definition of Shell/Cell/ResourceLease.

---

## 7. Public cloud fallback rule

Where **KHR cannot run directly** (public cloud constraints), KARL **continues** to generate VM/VMPool via **KubeVirt provider**. This is an **explicit fallback**, not a conceptual downgrade. The product model remains **Shell / Cell / ResourceLease**.

---

## 8. Hyperdensity role

Hyperdensity governs:

- Resource Market and **ResourceLease**
- Donor/receiver equilibrium
- Rollback and evidence chains
- Capacity prediction and dry-run decision loops

Hyperdensity does **not** replace KHR execution; it **governs** resource intelligence over providers.

---

## 9. Storage semantics

See `HYPERDENSITY_KHR_STORAGE_SEMANTICS.md`. KARL-native primitives: **EphemeralDisk**, modes (ephemeralOverlay, ephemeralClone, scratch, readonly, persistent), sources (pvc, image, snapshot, volume, goldenImage), discardPolicy, **promote-to-image**.

---

## 10. Network semantics

See `HYPERDENSITY_KHR_NETWORK_SDN_SEMANTICS.md`. KARL-native: **KARLNetwork**, **CellNetwork**, **ShellNetwork**, **NetworkAttachment**, **NetworkLease**, **NetworkPolicy**, plus OVN mapping preservation.

---

## 11. ResourceLease direction

See `HYPERDENSITY_KHR_RESOURCELEASE_DIRECTION.md`. Future lease dimensions: cpu, memory, storage, network, policy, provider, rollback, evidence, expiration, promotion.

---

## 12. Non-goals

- Broad `ObservationWiredV1` / `ProductionWiredV1`
- Indefinite micro-sprint adapter without architecture consolidation
- Removing KubeVirt or OVN compatibility paths
- Windows runtime enablement in observation sprints
- Dashboard `pkg/hyperdensity/parentfabric` runtime import

---

## 13. Next roadmap phases

See `HYPERDENSITY_KHR_ROADMAP_PHASES.md`. Phase 1 (Sprint 88) = architecture memory. Phase 2+ = ResourceLease contract, inventory facts, Hyperdensity dry-run, KubeVirt operational loop, KHR native seed, Shell/Cell providers, Windows DaaS profile.

---

## 14. Risks

- Treating KubeVirt fallback as permanent product model (must stay labeled fallback).
- Losing ephemeral PVC or OVN semantics during KARL-native promotion.
- Adapter drift without contract anchors for storage/network.

---

## 15. Related documents

- `HYPERDENSITY_KHR_STORAGE_SEMANTICS.md`
- `HYPERDENSITY_KHR_NETWORK_SDN_SEMANTICS.md`
- `HYPERDENSITY_KHR_RESOURCELEASE_DIRECTION.md`
- `HYPERDENSITY_KHR_ROADMAP_PHASES.md`
- `HYPERDENSITY_KHR_ROADMAP_TRANSITION_NOTE.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_MIGRATION_BOUNDARY.md`


---

## Sprint 89 (ResourceLease minimal contract)

Sprint 89 adds ResourceLease minimal contract sketch (storage/network/provider/examples). No CRD, no controller, no runtime. See HYPERDENSITY_KHR_RESOURCELEASE_MINIMAL_CONTRACT.md and related Sprint 89 docs.


---

## Sprint 90 (inventory facts)

Sprint 90 adds read-only KubeVirt and OVN/SDN capability inventory mapped to ResourceLease contract. No CRD, no controller, no runtime. See HYPERDENSITY_KHR_KUBEVIRT_CAPABILITY_INVENTORY.md and related Sprint 90 docs.
