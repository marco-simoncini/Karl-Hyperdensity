# Hyperdensity Parent Fabric — apply observation flip criteria (Sprint 72)

## Summary

**Sprint 72** is **criteria-only**: formal prerequisites for a future flip of `hyperdensityWorkloadAdapterApplyObservationCandidateRuntimeUsedV1` and/or `hyperdensityWorkloadAdapterApplyObservationWiredV1`. **No flag flip, no wiring, no runtime changes** in Sprint 72.

**Sprint 70** wired `apply.go` to four staged wrappers. **Sprint 71** post-wiring hardening PASS. The flip does **not** occur in Sprint 72.

---

## 1. Scope

| Item | Sprint 72 |
|------|-----------|
| Flip criteria + risks documentation | **Yes** |
| Hyperdensity Go adapter code | **No** |
| Dashboard runtime / `apply.go` changes | **No** |
| `ApplyObservationWiredV1 = true` | **No** |
| `ApplyObservationCandidateRuntimeUsedV1 = true` | **No** |
| `ObservationWiredV1 = true` | **No** |
| `ProductionWiredV1 = true` | **No** |

---

## 2. Non-goals

- Enabling apply observation flags in Sprint 72.
- Broad observation (`ObservationWiredV1`) flip.
- resource_exchange, rollback observed-state, VM runtime, admission_guard.
- Dashboard import of `pkg/hyperdensity/parentfabric`.
- Changing API response, JSON ordering, or apply payload shape.
- Moving runtime code from Dashboard to Hyperdensity.
- Changing `workload_helpers.go` verdict (remains **`copy-deferred`**).

---

## 3. Current apply observation state

| Item | Value |
|------|-------|
| `apply.go` wrapper call sites | **4** |
| `apply.go` legacy observation call sites | **0** |
| `ApplyObservationWiredV1` | **`false`** |
| `ApplyObservationCandidateRuntimeUsedV1` | **`false`** |
| `ApplyObservationCandidateV1` | **`true`** (shadow present) |
| `ApplyObservationShadowReadyV1` | **`true`** |
| `ObservationWiredV1` | **`false`** |
| `ProductionWiredV1` | **`false`** |
| Wrapper branch at runtime | **legacy-equivalent** (flags false) |
| Post-wiring 8×4 matrix | **PASS** (Sprint 71) |

Broad observation remains **deliberately disabled**.

---

## 4. Preconditions already satisfied

| # | Precondition | Sprint |
|---|--------------|--------|
| 1 | Apply observation proposal + criteria documented | 65 |
| 2 | Shadow matrix (candidate not runtime-used) | 66 |
| 3 | Staged wrappers in `apply_observation_wiring_v1.go` | 67 |
| 4 | Wrapper hardening 8×4 (pre-wiring) | 68 |
| 5 | Wiring readiness certified | 69 |
| 6 | `apply.go` call-site wiring (4 wrappers) | 70 |
| 7 | Post-wiring hardening: wrapper ≡ legacy ≡ candidate | 71 |
| 8 | `audit_workload_apply_observation.sh` Sprint 65–71 PASS | 71 |
| 9 | No `parentfabric` import in Dashboard runtime | ongoing |
| 10 | `workload_helpers.go` verdict **`copy-deferred`** | ongoing |

---

## 5. Preconditions still required before flip

| # | Precondition | Notes |
|---|--------------|-------|
| 1 | **Dedicated sprint** per flag (or explicit dual approval) | Do not bundle with broad observation or resource_exchange |
| 2 | **Flip criteria doc approved** | This document + `FLIP_RISKS.md` |
| 3 | **Golden update** for flip sprint only | `applyObservationWired` / `applyObservationCandidateRuntimeUsed` |
| 4 | **Parity + apply audit PASS** after flag change | Re-run full apply observation test suite |
| 5 | **8×4 matrix re-run** after each flag change | wrapper ≡ legacy ≡ candidate until candidate runtime enabled |
| 6 | **Rollback drill documented** | See §8 |
| 7 | **No API/payload drift** | Apply plan target fields unchanged in shape |
| 8 | **PR title states flip scope** | e.g. "Sprint N: ApplyObservationCandidateRuntimeUsedV1 flip" |

---

## 6. Candidate-runtime flip criteria

Minimum criteria before:

```text
hyperdensityWorkloadAdapterApplyObservationCandidateRuntimeUsedV1 = true
```

| # | Criterion |
|---|-----------|
| 1 | **Dedicated sprint** — not bundled with `ApplyObservationWiredV1` unless explicitly approved |
| 2 | **Post-wiring hardening PASS** — Sprint 71 green |
| 3 | **wrapper ≡ legacy ≡ candidate PASS** — 8×4 matrix green immediately before flip |
| 4 | **No resource_exchange / rollback / VM runtime touch** — apply-only |
| 5 | **`audit_workload_apply_observation.sh` PASS** — including flip-sprint guards |
| 6 | **Call-site wiring golden updated** — documents `applyGoCallSitesUseWrappers: true`, flags reflect flip |
| 7 | **Rollback documented** — set flag `false`; wrappers still call legacy until `ApplyObservationWiredV1` also true |
| 8 | **`ObservationWiredV1` remains `false`** | Apply subflag only |
| 9 | **`ProductionWiredV1` remains `false`** | |
| 10 | **Wrapper implementation** routes to candidate when both apply flags true | Verify in `apply_observation_wiring_v1.go` only |

**Recommended:** dedicated candidate-runtime flip sprint (after Sprint 73 readiness), separate from `ApplyObservationWiredV1` flip. **Sprint 73** = readiness only (no flip).

---

## 7. ApplyObservationWired flip criteria

Minimum criteria before:

```text
hyperdensityWorkloadAdapterApplyObservationWiredV1 = true
```

| # | Criterion |
|---|-----------|
| 1 | **Dedicated sprint** — separate from broad observation; separate from candidate flip unless explicitly approved |
| 2 | **Candidate runtime path proven** — `ApplyObservationCandidateRuntimeUsedV1=true` sprint completed **or** declared legacy-equivalent with evidence |
| 3 | **`ObservationWiredV1` remains `false`** | No broad observation |
| 4 | **`ProductionWiredV1` remains `false`** | |
| 5 | **No apply payload drift** — target memory/CPU fields same semantics |
| 6 | **No API/JSON ordering change** | |
| 7 | **8×4 matrix PASS** post-flip | |
| 8 | **Parity runner PASS** | |
| 9 | **Hyperdensity docs updated** | No Dashboard `parentfabric` import |

**Note:** With only `ApplyObservationWiredV1=true` and `ApplyObservationCandidateRuntimeUsedV1=false`, wrappers may still take legacy branch depending on implementation — confirm branch logic before flip.

---

## 8. Required rollback plan

| Step | Action |
|------|--------|
| 1 | Set flipped flag(s) back to **`false`** |
| 2 | If call-site regression: restore four legacy helper calls in `apply.go` (undo Sprint 70) |
| 3 | Re-run `test_hyperdensity_parity.sh` and `audit_workload_apply_observation.sh` |
| 4 | Restore golden snapshots to pre-flip values |
| 5 | Do **not** touch resource_exchange, rollback, VM runtime as part of apply rollback |

---

## 9. Required parity/audit coverage

Future flip sprint(s) must keep green:

- `TestHyperdensityParentFabricWorkloadApplyObservationPostWiringHardening`
- `TestHyperdensityParentFabricWorkloadApplyObservationCallsiteWiring`
- `TestHyperdensityParentFabricWorkloadApplyObservationBranchSwapGuard`
- `audit_workload_apply_observation.sh` (Sprint 65–71+ guards)
- `test_hyperdensity_parity.sh` (go test + runtime import + workload audits)

---

## 10. Risks

- **Flag confusion:** call-site wired ≠ `ApplyObservationWiredV1=true` ≠ broad `ObservationWiredV1=true`.
- **Bundled flip:** changing candidate + wired + broad in one PR.
- **resource_exchange coupling:** stage-apply observation mistaken for pod apply observation.
- **Payload drift:** candidate path returns different strings than legacy for edge cases.
- **Skipped re-hardening:** flipping without 8×4 re-run.

See **`HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_FLIP_RISKS.md`**.

---

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_POST_WIRING_HARDENING.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_FLIP_RISKS.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_CRITERIA.md`
- Dashboard `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_FLIP_CRITERIA_M64.md`
