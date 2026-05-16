# Hyperdensity Parent Fabric — apply observation post-activation hardening (Sprint 76)

## Summary

**Sprint 76** is **hardening-only**: certifies the **post-Sprint 75** activation state. Both apply flags are **true**; apply wrappers use the **candidate** branch; **wrapper ≡ candidate ≡ legacy** on the 8×4 matrix. **No flag, wiring, or runtime semantic changes.**

---

## 1. Scope

| Item | Sprint 76 |
|------|-----------|
| Post-activation 8×4 hardening matrix | **Yes** |
| Source invariants (`apply.go`, flags, AND gate) | **Yes** |
| Flag flips | **No** |
| Broad observation | **No** |

---

## 2. Non-goals

- Changing `ApplyObservationWiredV1` or `ApplyObservationCandidateRuntimeUsedV1`.
- Changing wrapper branch logic or candidate helper implementations.
- Modifying `apply.go` call sites.
- `ObservationWiredV1` or `ProductionWiredV1` flip.
- resource_exchange, rollback, VM runtime, admission_guard.
- Dashboard `parentfabric` import.
- Automatic broad observation as next step.

---

## 3. Current post-activation state

| Item | Value |
|------|-------|
| `ApplyObservationWiredV1` | **`true`** |
| `ApplyObservationCandidateRuntimeUsedV1` | **`true`** |
| Candidate branch (apply wrappers) | **active** |
| Runtime path (wrappers) | **candidate** |
| Effective behavior | **legacy-equivalent** (candidate ≡ legacy) |
| `apply.go` wrapper call sites | **4** |
| `apply.go` legacy observation call sites | **0** |
| `ObservationWiredV1` | **`false`** |
| `ProductionWiredV1` | **`false`** |
| `workload_helpers.go` verdict | **`copy-deferred`** |

---

## 4. Hardening matrix

Eight cases × four helpers (memory request/limit, CPU request/limit):

1. Normal container  
2. Missing container  
3. Empty container name  
4. Multi-container main  
5. Multi-container sidecar  
6. Missing resources  
7. Malformed resources map  
8. Nil pod map  

Per case: **wrapper == candidate == legacy**.

---

## 5. Candidate branch invariants

- AND gate present **4×** in `apply_observation_wiring_v1.go`.
- Both flags **true** → wrappers call candidate helpers (not legacy directly).
- `apply.go` calls **wrappers only** (zero direct candidate calls).
- Candidate helpers remain legacy-delegating implementations.

---

## 6. No-touch surfaces

| Surface | Status |
|---------|--------|
| resource_exchange_* | **untouched** |
| rollback observed-state | **untouched** |
| VM runtime observed-state | **untouched** |
| admission_guard_* | **untouched** |
| Broad observation call sites | **not wired** |

---

## 7. Drift risks

- Candidate helper edited to diverge from legacy while flags stay true.
- Accidental `ObservationWiredV1=true` bundled in unrelated PR.
- Wrapper/candidate names appearing in resource_exchange or rollback files.
- Operators assume apply activation implies broad observation enabled.

---

## 8. Rollback

1. Set `ApplyObservationWiredV1` to **`false`** (Sprint 75 rollback) — wrappers return legacy branch; `CandidateRuntimeUsedV1` may stay true.
2. Optionally set `ApplyObservationCandidateRuntimeUsedV1` to **`false`** (Sprint 74 rollback).
3. Re-run parity + apply audit.
4. No `apply.go` changes required for flag-only rollback.

---

## 9. Required parity/audit coverage

- `TestHyperdensityParentFabricWorkloadApplyObservationPostActivationHardening`
- `TestHyperdensityParentFabricWorkloadApplyObservationActivation`
- `TestHyperdensityParentFabricWorkloadApplyObservationBranchSwapGuard`
- `audit_workload_apply_observation.sh` Sprint 76 guards
- Historical goldens Sprint 65–74 (pre-activation snapshots) remain valid as archives

---

## 10. Next migration boundary

The **next boundary is not** automatic broad `ObservationWiredV1`. Acceptable future work:

- Continued apply-only hardening under activation state.
- Explicit broad-observation policy sprint (separate criteria).
- Hyperdensity `workload_helpers` copy boundary review (still **`copy-deferred`**).

Do **not** treat apply activation as permission to wire resource_exchange, rollback, or VM observation via apply flags.

---

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_ACTIVATION.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_BRANCH_LOGIC.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_BROAD_OBSERVATION_DECISION.md`
