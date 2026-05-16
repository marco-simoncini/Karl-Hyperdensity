# Hyperdensity Parent Fabric — resource exchange local helper criteria (Sprint 81+)

Criteria for staged wrappers and wiring of ready/restart helpers after Sprint 81 shadow matrix.

---

## 1. Staged wrappers (future sprint)

| Requirement | Detail |
|-------------|--------|
| Dedicated sprint | Ready and restart wrappers separate from CPU-only shortcut |
| AND gate | Same as CPU: `WiredV1 && CandidateRuntimeUsedV1` |
| No production use | Until explicit wiring sprint |
| wrapper == candidate == legacy | Per helper, 8+ cases each |

---

## 2. Wiring readiness

Before replacing legacy calls in `resource_exchange_*`:

| Decision | Options |
|----------|---------|
| Coverage scope | **CPU-only** vs **CPU + ready + restart** |
| `ContainerReady` locality | Defined in `resource_exchange_stage_apply.go` — wiring stays Dashboard-local |
| `ContainerRestartCount` | Lives in `workload_helpers.go` — no Hyperdensity copy without re-audit |
| Call-site count unchanged until flip | Audit must show 8 CPU + 12 ready + 12 restart legacy until wiring sprint |

---

## 3. Policy gates

- No `ObservationWiredV1` flip from local helper work.
- No rollback / VM / admission_guard bundling.
- No direct candidate calls in production.
- No API/payload drift.

---

## 4. Rollback strategy

| Step | Action |
|------|--------|
| 1 | `ResourceExchangeObservationWiredV1 = false` |
| 2 | Restore legacy ready/restart/CPU calls if wired |
| 3 | Re-run parity + audit |

---

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_LOCAL_HELPER_SHADOW_MATRIX.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_OBSERVATION_STAGED_WRAPPER_CRITERIA.md`


---

## Sprint 83 (call-site wiring)

Sprint 83 wires all 32 production `resource_exchange_*` observation call sites to full-helper staged wrappers (8 CPU + 12 ready + 12 restart). `ResourceExchangeObservationWiredV1` and `ResourceExchangeObservationCandidateRuntimeUsedV1` remain **false**; effective runtime path is **legacy**. Direct candidate calls in production remain forbidden. See `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_CALLSITE_WIRING.md` (Hyperdensity) and `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_CALLSITE_WIRING_M83.md` (Dashboard).
