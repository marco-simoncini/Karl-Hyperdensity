# Hyperdensity Parent Fabric — live observation flip criteria

## Summary

**Sprint 61** executed the dedicated flip: **`hyperdensityWorkloadAdapterLiveObservationWiredV1 = true`**. Criteria below were satisfied in Sprint 60; flip applied in Sprint 61.

---

## Minimum criteria before flip

| # | Criterion |
|---|-----------|
| 1 | **Live shadow hardening PASS** — `TestHyperdensityParentFabricWorkloadAdapterLiveObservationShadow` green (7 cases) |
| 2 | **Call-site audit PASS** — `audit_workload_adapter_call_sites.sh` |
| 3 | **Fixture parity** — Pod UID + CPU/memory request/limit shadow DeepEqual vs legacy |
| 4 | **Scope** — No apply / resource_exchange / rollback / VM runtime observation wiring |
| 5 | **Broad flag** — `hyperdensityWorkloadAdapterObservationWiredV1` remains **`false`** |
| 6 | **Rollback** — Documented revert path (flag false + optional direct legacy calls) |
| 7 | **Parity runner** — `test_hyperdensity_parity.sh` green |
| 8 | **Explicit sprint** — Dedicated sprint title/PR; not bundled with unrelated wiring |
| 9 | **Golden update** — `live_observation_shadow.golden.json` sets `liveObservationFlipAllowed: true` only in flip sprint |

---

## Sprint 60–61 status

| Item | Sprint 60 | Sprint 61 |
|------|-----------|-----------|
| Shadow hardening | **Done** | **PASS** (post-flip) |
| `LiveObservationWiredV1` | **`false`** | **`true`** |
| Flip executed | **No** | **Yes** |
| Broad `ObservationWiredV1` | **`false`** | **`false`** |

---

## Forbidden with flip

- `hyperdensityWorkloadAdapterObservationWiredV1 = true` (broad observation)
- `hyperdensityWorkloadAdapterProductionWiredV1 = true`
- Wiring live wrappers outside `live.go`
- Dashboard import of `pkg/hyperdensity/parentfabric`

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_LIVE_OBSERVATION_SHADOW_HARDENING.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_LIVE_OBSERVATION_STAGED_WIRING.md`
