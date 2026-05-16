# Hyperdensity Parent Fabric — resource exchange candidate-runtime staging (Sprint 84)

## Summary

**Sprint 84** sets `ResourceExchangeObservationCandidateRuntimeUsedV1 = true` while **`ResourceExchangeObservationWiredV1` remains `false`**. With the AND gate on all three wrappers, **runtime stays legacy** — the candidate branch is **not** active. This is **staging only**, not activation.

---

## 1. Scope

| Item | Sprint 84 |
|------|-----------|
| `ResourceExchangeObservationCandidateRuntimeUsedV1 = true` | **Yes** |
| `ResourceExchangeObservationWiredV1 = true` | **No** |
| Production call-site changes | **No** (Sprint 83 wiring retained) |
| Wrapper / candidate helper logic changes | **No** |
| Runtime behavior change | **No** (AND gate) |

---

## 2. Non-goals

- Activation flip (`ResourceExchangeObservationWiredV1 = true`).
- Broad `ObservationWiredV1` / `ProductionWiredV1`.
- Direct candidate calls in production `resource_exchange_*`.
- rollback, VM runtime, admission_guard, apply track changes.
- Dashboard `parentfabric` import.
- Storage/network primitive implementation.

---

## 3. Preconditions from Sprint 78–83

| Sprint | Deliverable |
|--------|-------------|
| 78–81 | Candidates + shadow matrices |
| 82 | Full-helper staged wrappers |
| 83 | 32 production call-sites wired to wrappers (8/12/12) |

---

## 4. Flag change

**Only** in `hyperdensity_parent_fabric_workload_resource_exchange_observation_candidate_v1.go`:

```text
ResourceExchangeObservationCandidateRuntimeUsedV1: false → true
```

Unchanged:

- `ResourceExchangeObservationCandidateV1 = true`
- `ResourceExchangeObservationShadowReadyV1 = true`
- `ResourceExchangeObservationWiredV1 = false`

---

## 5. AND gate behavior

```go
if ResourceExchangeObservationWiredV1 && ResourceExchangeObservationCandidateRuntimeUsedV1 {
    return candidate(...)
}
return legacy(...)
```

With **Wired=false**, **CandidateUsed=true** → condition **false** → **legacy** path.

---

## 6. Runtime behavior statement

| Assertion | Value |
|-----------|-------|
| Effective wrapper path | **legacy** |
| Candidate branch active | **no** |
| Production uses wrappers | **yes** (not direct candidates) |
| API / payload | **unchanged** |

---

## 7. Shadow re-validation

Full-helper **24-case** matrix (CPU 8 + ready 8 + restart 8): **wrapper ≡ candidate ≡ legacy** after flip.

---

## 8. Rollback

Set `ResourceExchangeObservationCandidateRuntimeUsedV1 = false` in candidate file; re-run parity.

---

## 9. Risks

Accidental `WiredV1=true` without dedicated activation sprint. Confusion that CandidateUsed=true implies active candidate path.

---

## 10. Recommended next sprint

**Activation readiness** or **activation flip** (`ResourceExchangeObservationWiredV1=true`) in a **dedicated** sprint with separate approval — not combined with unrelated surfaces.

---

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_CANDIDATE_RUNTIME_REVALIDATION.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_CALLSITE_WIRING.md`


---

## Sprint 85 (activation readiness)

Sprint 85 is readiness-only for `ResourceExchangeObservationWiredV1=true`. No flag changes. Sprint 86 may execute activation flip if approved. See `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_ACTIVATION_READINESS.md`.
