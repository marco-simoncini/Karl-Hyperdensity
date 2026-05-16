# Hyperdensity Parent Fabric — apply observation staged wrappers (Sprint 67)

## Summary

**Sprint 67** adds **staged wrappers** for apply observation in Dashboard (`apply_observation_wiring_v1.go`). **Sprint 68** hardens wrappers (8×4 matrix). **`apply.go` does not call wrappers yet.** All runtime flags remain **false**.

---

## 1. Scope

| Item | Sprint 67 |
|------|-----------|
| Staged wrapper functions | **Yes** (Dashboard) |
| `apply.go` call-site replacement | **No** |
| `ApplyObservationWiredV1 = true` | **No** |
| `ApplyObservationCandidateRuntimeUsedV1 = true` | **No** |
| Hyperdensity Go code | **No** |

---

## 2. Non-goals

- Wiring `apply.go` to wrappers.
- Branch swap to candidate true branch.
- resource_exchange, rollback, VM, admission_guard.
- Broad `ObservationWiredV1` flip.
- Dashboard `parentfabric` import.

---

## 3. Wrapper shape

```text
hyperdensityWorkloadApplyObservedPod*V1(pod, containerName)
  if ApplyObservationWiredV1 && ApplyObservationCandidateRuntimeUsedV1 → candidate
  else → legacy
```

Four wrappers: memory request/limit, CPU request/limit.

---

## 4. Flag behavior

| Flag | Sprint 67 |
|------|-----------|
| `ApplyObservationWiredV1` | **`false`** → legacy branch |
| `ApplyObservationCandidateRuntimeUsedV1` | **`false`** → legacy branch |
| `ObservationWiredV1` | **`false`** |
| `ProductionWiredV1` | **`false`** |

With both apply flags false, wrappers ≡ legacy ≡ candidate (shadow tests).

---

## 5. Test matrix

5 cases × 4 helpers: wrapper == legacy == candidate.

Guards: wrappers not in `apply.go`; branch-swap guard; audit script.

---

## 6. Future wiring requirements

1. Dedicated sprint to repoint `apply.go` call sites to wrappers.
2. Shadow hardening green with wrappers.
3. Optional flip `ApplyObservationWiredV1` then candidate runtime swap (separate sprints).
4. `ObservationWiredV1` stays false until all surfaces migrated.

---

## 7. Rollback

Remove wiring file reference from tests; delete `apply_observation_wiring_v1.go`. `apply.go` unchanged in Sprint 67.

---

## 8. Risks

Accidental `apply.go` import of wrappers before shadow/flip criteria met.

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_SHADOW_MATRIX.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_CRITERIA.md`
