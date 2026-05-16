# Hyperdensity Parent Fabric — apply observation shadow matrix (Sprint 66)

## Summary

**Sprint 66** designs and validates the **shadow matrix** for the four `apply.go` legacy observation call sites. **Sprint 67** adds **staged wrappers** (see `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_STAGED_WRAPPERS.md`). **No runtime wiring.** **`apply.go` remains legacy.**

---

## 1. Scope

| Item | Sprint 66 |
|------|-----------|
| Shadow matrix doc (this file) | **Yes** |
| Dashboard candidate helpers + shadow tests | **Yes** |
| `apply.go` call-site replacement | **No** |
| `ApplyObservationWiredV1 = true` | **No** |
| Hyperdensity Go adapter code | **No** |

---

## 2. Non-goals

- Wiring `apply.go` to wrappers or candidates.
- `ApplyObservationCandidateRuntimeUsedV1 = true`.
- Branch swap in apply path.
- resource_exchange, rollback, VM runtime, admission_guard observation.
- Broad `ObservationWiredV1` flip.
- Dashboard `parentfabric` import.

---

## 3. Apply call-site matrix

| # | Legacy helper | `apply.go` lines (Sprint 65) | Shadow case |
|---|---------------|------------------------------|-------------|
| 1 | `hyperdensityObservedPodMemoryRequest` | ~271 | Normal / missing / multi / no resources |
| 2 | `hyperdensityObservedPodMemoryLimit` | ~272 | Same matrix |
| 3 | `hyperdensityObservedPodCPURequest` | ~278 | Same matrix |
| 4 | `hyperdensityObservedPodCPULimit` | ~279 | Same matrix |

**Five shadow cases** × **four helpers** = matrix coverage (Dashboard test).

---

## 4. Shadow fixture requirements

- Golden: `hyperdensity_parent_fabric_workload_apply_observation_shadow_matrix.golden.json`
- `applyObservationCandidatePresent: true`
- `applyObservationCandidateRuntimeUsed: false`
- `branchSwapAllowed: false`
- `callSitesRemainLegacy: true`

---

## 5. Future wrapper/candidate shape

```text
hyperdensityWorkloadApplyCandidateObservedPod*V1(pod, containerName)
  → future apply observation wrapper true branch (dedicated sprint)
```

Sprint 66 candidate may delegate to legacy internally; semantic equivalence proven by shadow tests.

---

## 6. Required assertions

Per shadow case and helper:

- `candidate == legacy`
- `apply.go` contains legacy helper names
- `apply.go` does **not** contain candidate helper names
- `ApplyObservationWiredV1 == false`
- `ApplyObservationCandidateRuntimeUsedV1 == false`
- `ObservationWiredV1 == false`
- `ProductionWiredV1 == false`

---

## 7. Rollback

Sprint 66 is test-only — rollback: remove candidate file and shadow tests; flags unchanged.

---

## 8. Risks

Premature `ApplyObservationCandidateRuntimeUsedV1 = true` before dedicated wiring sprint could change apply plan targets.

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_CRITERIA.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_APPLY_OBSERVATION_PROPOSAL.md`
