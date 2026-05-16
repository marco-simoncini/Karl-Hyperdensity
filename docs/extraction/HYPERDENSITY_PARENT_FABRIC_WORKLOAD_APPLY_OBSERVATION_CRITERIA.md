# Hyperdensity Parent Fabric — apply observation criteria (Sprint 65)

## Purpose

Minimum criteria before creating or activating:

```text
hyperdensityWorkloadAdapterApplyObservationWiredV1 = true
```

---

## Mandatory criteria

| # | Criterion |
|---|-----------|
| 1 | **Dedicated sprint** — not bundled with live/pilot/broad flips |
| 2 | **`ObservationWiredV1` remains `false`** — apply uses subflag only |
| 3 | **`ProductionWiredV1` remains `false`** |
| 4 | **Apply-only golden** — documents every call-site in `apply.go` |
| 5 | **Shadow fixture** per apply observation helper invocation |
| 6 | **Rollback documented** — flag false + legacy helpers restored |
| 7 | **resource_exchange excluded** — no stage-apply observation in same sprint |
| 8 | **VM runtime excluded** |
| 9 | **admission_guard excluded** |
| 10 | **`parentfabric` import forbidden** in Dashboard production runtime |
| 11 | **`apply.go` changes** only in wiring sprint; proposal sprint touches audit/docs only |
| 12 | **Parity + apply audit script PASS** |

---

## Pre-flip checklist

- [x] Sprint 65 proposal + criteria merged (this document).
- [x] Sprint 66 shadow matrix implemented and green (candidate runtime **not** used).
- [x] Sprint 67 staged wrappers present (`apply.go` still legacy).
- [x] Sprint 68 wrapper hardening (8×4) green.
- [x] Sprint 69 wiring readiness certified (`readyForApplyGoCallSiteWiring`).
- [x] Sprint 70 `apply.go` call-site wiring (4 wrappers; flags **false**).
- [x] Sprint 71 post-wiring hardening green.
- [x] Sprint 72 flip criteria + risks documented (no flip).
- [x] Sprint 73 candidate-runtime readiness + branch logic (no flip).
- [x] Sprint 74 `ApplyObservationCandidateRuntimeUsedV1 = true` (staging flip; Wired **false**).
- [x] Sprint 75 `ApplyObservationWiredV1 = true` (activation; candidate branch active; broad **false**).
- [ ] `audit_workload_apply_observation.sh` reports zero forbidden patterns.
- [ ] Golden `applyObservationWired: false` updated only when flip sprint approved.
- [ ] Hyperdensity docs updated; no Dashboard `parentfabric` import.

---

## Explicit exclusions

The following must **not** be enabled via `ApplyObservationWiredV1`:

- Broad observation flip.
- Resource-exchange observation wiring.
- Rollback observed-state wiring.
- VM runtime observation wiring.

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_PROPOSAL.md`
