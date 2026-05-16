# Hyperdensity Parent Fabric — apply observation candidate-runtime flip (Sprint 74)

## Summary

**Sprint 74** executed a **staging flip only**. **Sprint 75** activated `ApplyObservationWiredV1=true` (candidate branch active).: `hyperdensityWorkloadAdapterApplyObservationCandidateRuntimeUsedV1 = true` in Dashboard. **`hyperdensityWorkloadAdapterApplyObservationWiredV1` remains `false`**. With the AND gate certified in Sprint 73, **runtime behavior stays legacy-equivalent** — the candidate branch is **not** active.

---

## 1. Scope

| Item | Sprint 74 |
|------|-----------|
| `ApplyObservationCandidateRuntimeUsedV1 = true` | **Yes** (Dashboard) |
| `ApplyObservationWiredV1 = true` | **No** |
| Runtime behavior change | **No** (AND gate) |
| `apply.go` / wrapper branch logic changes | **No** |
| Hyperdensity Go code | **No** |

---

## 2. Non-goals

- Activation flip (`ApplyObservationWiredV1 = true`) — future sprint.
- Broad `ObservationWiredV1` or `ProductionWiredV1` flip.
- Changing wrapper AND branch logic or candidate helper semantics.
- resource_exchange, rollback, VM runtime, admission_guard.
- Dashboard `parentfabric` import.
- API/JSON ordering / apply payload changes.

---

## 3. Preconditions from Sprint 72–73

| Sprint | Deliverable |
|--------|-------------|
| **72** | Flip criteria + risks |
| **73** | Candidate-runtime readiness + branch logic (AND gate × 4) |
| **71** | Post-wiring 8×4 PASS |
| **70** | Four wrapper call sites in `apply.go` |

---

## 4. Flag change

**Only** in `hyperdensity_parent_fabric_workload_adapter_apply_observation_candidate_v1.go`:

```text
hyperdensityWorkloadAdapterApplyObservationCandidateRuntimeUsedV1: false → true
```

Unchanged:

- `hyperdensityWorkloadAdapterApplyObservationCandidateV1 = true`
- `hyperdensityWorkloadAdapterApplyObservationShadowReadyV1 = true`
- `hyperdensityWorkloadAdapterApplyObservationWiredV1 = false` (in `adapter_v1.go`)

---

## 5. AND gate behavior

```go
if ApplyObservationWiredV1 && ApplyObservationCandidateRuntimeUsedV1 {
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
| wrapper ≡ legacy ≡ candidate (8×4) | **PASS** |
| `apply.go` direct candidate calls | **0** |
| API / apply payload | **unchanged** |

---

## 7. Rollback

Set `ApplyObservationCandidateRuntimeUsedV1` back to **`false`** in `apply_observation_candidate_v1.go`. Re-run parity + apply audit. No `apply.go` changes required for rollback.

---

## 8. Required parity/audit coverage

- `TestHyperdensityParentFabricWorkloadApplyObservationCandidateRuntimeFlip`
- `TestHyperdensityParentFabricWorkloadApplyObservationBranchSwapGuard`
- `TestHyperdensityParentFabricWorkloadApplyObservationPostWiringHardening`
- `audit_workload_apply_observation.sh` Sprint 74 guards
- Historical Sprint 66–73 goldens remain snapshots (candidate flag **false** at capture time)

---

## 9. Risks

- Operators assume flag true implies candidate runtime active (it does **not** until Wired true).
- Accidental `ApplyObservationWiredV1=true` in same PR (bundled activation).
- Skipping 8×4 re-run after a later activation flip.

---

## 10. Next activation criteria

Dedicated sprint may set `ApplyObservationWiredV1 = true` when approved:

- Sprint 74 flip green;
- 8×4 matrix re-run post-activation;
- Flip criteria Sprint 72 §7 satisfied;
- Rollback documented;
- `ObservationWiredV1` and `ProductionWiredV1` remain **false**.

---

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_BRANCH_LOGIC.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_CANDIDATE_RUNTIME_READINESS.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_FLIP_CRITERIA.md`
