# Hyperdensity Parent Fabric — resource exchange observation activation (Sprint 86)

## Summary

**Sprint 86** flips **`ResourceExchangeObservationWiredV1 = true`** with **`ResourceExchangeObservationCandidateRuntimeUsedV1` already true**. The candidate branch is **active** in resource_exchange wrappers only. Runtime behavior remains **unchanged** because wrapper ≡ candidate ≡ legacy.

---

## 1. Scope

| Item | Sprint 86 |
|------|-----------|
| `ResourceExchangeObservationWiredV1 = true` | **Yes** |
| Candidate branch active (resource_exchange wrappers) | **Yes** |
| `ObservationWiredV1` / `ProductionWiredV1` | **Unchanged (false)** |
| Production call-site changes | **No** |

---

## 2. Non-goals

- Broad observation.
- Direct candidate calls in production.
- rollback, VM, admission_guard changes.
- Dashboard `parentfabric` import.

---

## 3. Preconditions from Sprint 78–85

Sprint 83 call-site wiring, Sprint 84 candidate-runtime staging, Sprint 85 activation readiness.

---

## 4. Flag change

```text
ResourceExchangeObservationWiredV1: false → true
```

`ResourceExchangeObservationCandidateRuntimeUsedV1` remains **true**.

---

## 5. Candidate branch activation

AND gate: `Wired && CandidateUsed` → **true** → wrappers call **candidate** helpers.

---

## 6. Runtime behavior statement

Expected **unchanged** API/payload: candidate helpers delegate to same legacy semantics today.

---

## 7. Production source invariants

Wrapper production **8/12/12**; legacy **0/0/0**; direct candidate **0**.

---

## 8. Rollback

Set `ResourceExchangeObservationWiredV1 = false`; re-run parity; no call-site restore.

---

## 9. Risks

Accidental broad observation flip; semantic drift in candidate helpers.

---

## 10. Recommended next sprint

Resource_exchange **boundary closure** or remaining-surface decision (rollback/VM) — not broad observation.

---

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_POST_ACTIVATION_HARDENING.md`
