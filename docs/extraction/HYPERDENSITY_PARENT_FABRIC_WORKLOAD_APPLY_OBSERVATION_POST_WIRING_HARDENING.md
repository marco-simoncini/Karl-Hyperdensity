# Hyperdensity Parent Fabric — apply observation post-wiring hardening (Sprint 71)

## Summary

**Sprint 71** hardens the **post-Sprint 70** state: `apply.go` uses four staged wrappers; all apply flags remain **false**; wrappers delegate to legacy (legacy-equivalent behavior). **Hardening-only** — no wiring or flag changes.

---

## 1. Scope

| Item | Sprint 71 |
|------|-----------|
| Post-wiring 8×4 hardening matrix | **Yes** |
| `apply.go` source invariants | **Yes** |
| Flag flips | **No** |

---

## 2. Non-goals

- Changing the four wrapper call sites in `apply.go`.
- `ApplyObservationWiredV1` or `ApplyObservationCandidateRuntimeUsedV1 = true`.
- resource_exchange, rollback, VM, admission_guard.
- Broad `ObservationWiredV1` flip.

---

## 3. Current post-wiring state

| Item | Value |
|------|-------|
| `apply.go` wrapper call sites | **4** |
| `apply.go` legacy observation call sites | **0** |
| `ApplyObservationWiredV1` | **`false`** |
| `ApplyObservationCandidateRuntimeUsedV1` | **`false`** |

---

## 4. Hardening assertions

Per case × helper: **wrapper == legacy == candidate**.

Source: `apply.go` contains wrappers only (no legacy observation calls, no candidate calls).

---

## 5. Wrapper branch behavior

`if ApplyObservationWiredV1 && ApplyObservationCandidateRuntimeUsedV1` → false → **legacy** branch always taken.

---

## 6. No-touch surfaces

resource_exchange, rollback, VM runtime, admission_guard — no apply wrapper/candidate references.

---

## 7. Rollback

Restore four legacy helper calls in `apply.go` (Sprint 70 rollback). Sprint 71 is test/doc only.

---

## 8. Risks

Regression if flags flipped without re-hardening. Confusion that wired call sites imply `ApplyObservationWiredV1=true`.

---

## 9. Future branch-swap path

Separate sprint(s): see **`HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_FLIP_CRITERIA.md`** (Sprint 72). Recommended: candidate runtime flip first, then `ApplyObservationWiredV1=true`.

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_FLIP_CRITERIA.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_FLIP_RISKS.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_CALLSITE_WIRING.md`
