# Hyperdensity Parent Fabric — apply observation wrapper hardening (Sprint 68)

## Summary

**Sprint 68** expands shadow/hardening coverage for apply observation **staged wrappers** (Sprint 67). **Sprint 69** records wiring readiness. **Sprint 70** wired `apply.go` call sites. **Sprint 71** post-wiring hardening. All apply flags remain **false**.

---

## 1. Scope

| Item | Sprint 68 |
|------|-----------|
| Hardening matrix 8×4 (Dashboard test) | **Yes** |
| Stricter apply observation audit | **Yes** |
| `apply.go` call-site changes | **No** |
| `ApplyObservationWiredV1 = true` | **No** |

---

## 2. Non-goals

- Repointing `apply.go` to wrappers.
- `ApplyObservationCandidateRuntimeUsedV1 = true`.
- resource_exchange, rollback, VM, admission_guard.
- Broad `ObservationWiredV1` flip.
- Hyperdensity Go adapter code or Dashboard `parentfabric` import.

---

## 3. Hardening matrix

| # | Case |
|---|------|
| 1 | Normal container |
| 2 | Missing container |
| 3 | Empty container name |
| 4 | Multi-container main |
| 5 | Multi-container sidecar |
| 6 | Missing resources |
| 7 | Malformed resources map |
| 8 | Nil pod map |

Per case × 4 helpers: **wrapper == legacy == candidate**.

---

## 4. Wrapper branch behavior

With `ApplyObservationWiredV1=false` and `ApplyObservationCandidateRuntimeUsedV1=false`, wrappers always take the **legacy** branch (legacy-equivalent output).

---

## 5. Apply.go no-touch invariant

`apply.go` must retain four legacy helper call sites and must **not** reference wrapper or candidate function names.

---

## 6. Future wiring readiness

Dedicated sprint required to:

1. Repoint `apply.go` call sites to wrappers.
2. Run flip/hardening gates with updated goldens.
3. Keep `ObservationWiredV1=false` until broad policy allows.

---

## 7. Rollback

Remove hardening test/golden; Sprint 67 wrappers unchanged. No runtime rollback needed.

---

## 8. Risks

Insufficient hardening before `apply.go` wiring could miss edge cases (malformed pod maps, nil pod).

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_STAGED_WRAPPERS.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_SHADOW_MATRIX.md`
