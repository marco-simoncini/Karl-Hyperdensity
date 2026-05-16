# Hyperdensity Parent Fabric — remaining observation surface decision (Sprint 87)

## Summary

Classifies observation surfaces **after** apply track completion (Sprint 77) and resource_exchange track completion (Sprint 87). No wiring authorized except audit/proposal unless a dedicated sprint chain exists.

---

## 1. rollback

| Field | Value |
|-------|-------|
| **status** | legacy |
| **risk** | safety-critical |
| **next allowed action** | audit / proposal only |
| **wiring** | **forbidden** until dedicated shadow matrix and rollback policy |

---

## 2. VM runtime

| Field | Value |
|-------|-------|
| **status** | legacy |
| **risk** | high |
| **next allowed action** | audit / proposal only |
| **wiring** | **forbidden**; no Windows runtime claim in observation sprints |

---

## 3. admission_guard

| Field | Value |
|-------|-------|
| **status** | legacy |
| **risk** | policy-critical |
| **next allowed action** | classification / audit only |
| **wiring** | **forbidden** without dedicated policy sprint |

---

## 4. usage.go / other-review

| Field | Value |
|-------|-------|
| **status** | review-needed |
| **risk** | medium |
| **next allowed action** | classification sprint |
| **wiring** | per-file decision after audit |

---

## 5. broad observation

| Field | Value |
|-------|-------|
| **status** | explicitly disabled |
| `ObservationWiredV1` | **remains false** |
| `ProductionWiredV1` | **remains false** |
| **automatic flip** | **not authorized** after apply or resource_exchange completion |

---

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_BROAD_OBSERVATION_DECISION.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_MIGRATION_BOUNDARY.md`
- `HYPERDENSITY_KHR_ROADMAP_TRANSITION_NOTE.md`


---

## Sprint 88 (KHR architecture memory)

Sprint 88 canonizes KHR/KARL architecture memory, storage semantics, and network/OVN semantics. No runtime/adapter changes. KubeVirt remains compatibility provider and public-cloud fallback. See HYPERDENSITY_KHR_ARCHITECTURE_MEMORY.md and related Sprint 88 docs.
