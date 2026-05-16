# Hyperdensity Parent Fabric — apply observation candidate-runtime readiness (Sprint 73)

## Summary

**Sprint 73** was **readiness/proposal-only**. **Sprint 74** executed staging flip (`CandidateRuntimeUsedV1=true`; Wired **false**; runtime legacy). for a future flip of `hyperdensityWorkloadAdapterApplyObservationCandidateRuntimeUsedV1 = true`. **No flag flip, no wiring, no runtime changes** in Sprint 73.

**Sprint 72** documented flip criteria and risks. This sprint certifies readiness and clarifies candidate-runtime semantics.

---

## 1. Scope

| Item | Sprint 73 |
|------|-----------|
| Candidate-runtime readiness doc | **Yes** |
| Branch logic clarification doc | **Yes** |
| Dashboard readiness golden + test | **Yes** |
| `ApplyObservationCandidateRuntimeUsedV1 = true` | **No** |
| `ApplyObservationWiredV1 = true` | **No** |
| `apply.go` / wrapper / candidate code changes | **No** |

---

## 2. Non-goals

- Flipping any adapter flag.
- Changing wrapper branch logic or candidate helper implementations.
- Broad observation (`ObservationWiredV1`) or `ProductionWiredV1` flip.
- resource_exchange, rollback, VM runtime, admission_guard.
- Dashboard `pkg/hyperdensity/parentfabric` import.
- API response, JSON ordering, or apply payload shape changes.
- Changing `workload_helpers.go` verdict (remains **`copy-deferred`**).

---

## 3. Current post-Sprint 72 state

| Item | Value |
|------|-------|
| `apply.go` wrapper call sites | **4** |
| `apply.go` legacy observation call sites | **0** |
| `ApplyObservationCandidateRuntimeUsedV1` | **`false`** |
| `ApplyObservationWiredV1` | **`false`** |
| `ObservationWiredV1` | **`false`** |
| `ProductionWiredV1` | **`false`** |
| Runtime wrapper path | **legacy** (both apply flags false) |
| Post-wiring 8×4 | **PASS** (Sprint 71) |
| Flip criteria | **Documented** (Sprint 72) |

---

## 4. Candidate-runtime semantics

`ApplyObservationCandidateRuntimeUsedV1` is a **subflag** indicating that apply observation wrappers may route to **candidate helpers** when combined with `ApplyObservationWiredV1`.

| Flag value | Meaning |
|------------|---------|
| **`false`** (current) | Wrappers never select candidate at runtime; legacy path always taken. |
| **`true`** (future) | Wrapper **may** select candidate **only when** `ApplyObservationWiredV1` is also **`true`**. |

**Critical:** Setting **only** `ApplyObservationCandidateRuntimeUsedV1=true` with `ApplyObservationWiredV1=false` produces **no runtime behavior change** under current branch logic (AND gate). A dedicated flip sprint may set the flag for certification/audit while behavior stays legacy-equivalent until a later `ApplyObservationWiredV1` flip.

Candidate helpers exist for shadow/hardening (`wrapper ≡ legacy ≡ candidate`); they are **not** invoked from `apply.go` directly.

---

## 5. Wrapper branch logic to verify

See **`HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_BRANCH_LOGIC.md`**.

Dashboard test `TestHyperdensityParentFabricWorkloadApplyObservationCandidateRuntimeReadiness` verifies:

- All four wrappers use the same AND condition.
- `apply.go` does not call candidate helpers.
- Flags remain **false**.

---

## 6. Required future flip changes

Dedicated **candidate-runtime flip sprint** (not Sprint 73):

| # | Change |
|---|--------|
| 1 | Set `hyperdensityWorkloadAdapterApplyObservationCandidateRuntimeUsedV1 = true` in `apply_observation_candidate_v1.go` |
| 2 | Update readiness/flip goldens |
| 3 | Re-run 8×4 matrix (wrapper ≡ legacy ≡ candidate) |
| 4 | Extend `audit_workload_apply_observation.sh` flip guards |
| 5 | Confirm runtime still legacy-equivalent if `ApplyObservationWiredV1` remains **false** |

**Separate sprint** for `ApplyObservationWiredV1=true` when candidate path should become live.

---

## 7. Required rollback plan

| Step | Action |
|------|--------|
| 1 | Set `ApplyObservationCandidateRuntimeUsedV1` to **`false`** |
| 2 | Re-run parity + apply audit |
| 3 | Restore goldens to pre-flip values |
| 4 | If `apply.go` was changed: restore four legacy calls (Sprint 70 rollback) |
| 5 | Do **not** touch resource_exchange, rollback, VM |

---

## 8. Required parity/audit coverage

- `TestHyperdensityParentFabricWorkloadApplyObservationCandidateRuntimeReadiness` (Sprint 73)
- `TestHyperdensityParentFabricWorkloadApplyObservationPostWiringHardening`
- `TestHyperdensityParentFabricWorkloadApplyObservationBranchSwapGuard`
- `audit_workload_apply_observation.sh` Sprint 65–73
- `test_hyperdensity_parity.sh`

---

## 9. Risks

- Assuming `CandidateRuntimeUsedV1=true` alone changes apply behavior (it does **not** until `ApplyObservationWiredV1=true`).
- Flipping both flags in one sprint without re-hardening.
- Confusing candidate-runtime flip with broad `ObservationWiredV1`.
- Modifying branch logic to OR-gate without dedicated review.

---

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_BRANCH_LOGIC.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_FLIP_CRITERIA.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_FLIP_RISKS.md`
