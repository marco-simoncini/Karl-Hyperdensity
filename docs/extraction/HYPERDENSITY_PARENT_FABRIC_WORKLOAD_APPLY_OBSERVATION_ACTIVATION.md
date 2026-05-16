# Hyperdensity Parent Fabric — apply observation activation (Sprint 75)

## Summary

**Sprint 75** executes the **activation flip**: `hyperdensityWorkloadAdapterApplyObservationWiredV1 = true`. **`hyperdensityWorkloadAdapterApplyObservationCandidateRuntimeUsedV1`** was already **`true`** (Sprint 74). With the AND gate, **apply wrappers now take the candidate branch**. Behavior remains **legacy-equivalent** because candidate helpers ≡ legacy on the 8×4 matrix.

**Not** broad observation — `ObservationWiredV1` and `ProductionWiredV1` remain **`false`**.

---

## 1. Scope

| Item | Sprint 75 |
|------|-----------|
| `ApplyObservationWiredV1 = true` | **Yes** (Dashboard `adapter_v1.go`) |
| `ApplyObservationCandidateRuntimeUsedV1` | **true** (unchanged from Sprint 74) |
| Candidate branch in apply wrappers | **Active** |
| Runtime output change | **None** (candidate ≡ legacy) |
| Broad `ObservationWiredV1` | **false** |

---

## 2. Non-goals

- Broad observation flip.
- `ProductionWiredV1 = true`.
- Changing `apply.go` call sites, wrapper branch logic, or candidate helper implementations.
- resource_exchange, rollback, VM runtime, admission_guard.
- Dashboard `parentfabric` import.
- API/JSON ordering changes.

---

## 3. Preconditions from Sprint 70–74

| Sprint | Deliverable |
|--------|-------------|
| **70** | Four wrapper call sites in `apply.go` |
| **71** | Post-wiring 8×4 PASS |
| **72–73** | Flip criteria + candidate readiness + branch logic |
| **74** | `CandidateRuntimeUsedV1 = true` (staging; legacy path while Wired false) |

---

## 4. Flag change

**Only** in `hyperdensity_parent_fabric_workload_adapter_v1.go`:

```text
hyperdensityWorkloadAdapterApplyObservationWiredV1: false → true
```

Unchanged:

- `hyperdensityWorkloadAdapterApplyObservationCandidateRuntimeUsedV1 = true` (`candidate_v1.go`)
- `hyperdensityWorkloadAdapterObservationWiredV1 = false`
- `hyperdensityWorkloadAdapterProductionWiredV1 = false`

---

## 5. Candidate branch activation

```go
if ApplyObservationWiredV1 && ApplyObservationCandidateRuntimeUsedV1 {
    return candidate(...)  // now taken
}
return legacy(...)
```

Both flags **true** → **candidate** path active in all four apply observation wrappers.

---

## 6. Behavior equivalence statement

| Assertion | Value |
|-----------|-------|
| Effective wrapper implementation | **candidate** |
| wrapper ≡ candidate ≡ legacy (8×4) | **PASS** |
| Apply plan target fields | **unchanged** vs pre-activation |
| `apply.go` direct candidate calls | **0** |
| Broad observation | **disabled** |

---

## 7. Rollback

1. Set `ApplyObservationWiredV1` to **`false`** in `adapter_v1.go`.
2. Re-run parity + apply audit.
3. Wrappers return to legacy branch (CandidateUsed may remain true).
4. No `apply.go` changes required.

---

## 8. Required parity/audit coverage

- `TestHyperdensityParentFabricWorkloadApplyObservationActivation`
- `TestHyperdensityParentFabricWorkloadApplyObservationBranchSwapGuard`
- `audit_workload_apply_observation.sh` Sprint 75 guards
- Historical Sprint 65–74 goldens remain snapshots where `applyObservationWired: false`

---

## 9. Risks

- Assuming activation enables broad observation (it does **not**).
- Candidate ≡ legacy invariant broken by future candidate helper edits.
- Bundled `ObservationWiredV1=true` in same PR.

---

## 10. Next hardening step

Post-activation monitoring sprint: re-run 8×4 under load fixtures; document any drift if candidate helpers diverge from legacy; keep broad observation blocked until dedicated policy sprint.

---

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_CANDIDATE_RUNTIME_FLIP.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_BRANCH_LOGIC.md`
