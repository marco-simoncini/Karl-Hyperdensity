# Hyperdensity Parent Fabric — resource exchange local helper classification (Sprint 79)

## Summary

Sprint 78 listed-helper inventory covered **`hyperdensityObservedPodCPURequest`** only. Sprint 79 classifies **local observation-like helpers** used by resource_exchange but **excluded** from the Sprint 79 shadow matrix.

---

## hyperdensityObservedPodContainerReady

| Attribute | Detail |
|-----------|--------|
| **Where** | `resource_exchange_stage_apply.go` (definition + calls); `resource_exchange_stage_apply_chain.go` (calls) |
| **Call sites (resource_exchange)** | **12** |
| **Definition** | **Local** to `hyperdensity_parent_fabric_resource_exchange_stage_apply.go` |
| **Observation class** | Pod container **ready** status from `status.containerStatuses` |
| **In Sprint 79 shadow matrix** | **No** |
| **Future matrix** | **Yes** — if resource_exchange wiring proceeds (Sprint 80+) |
| **Safe for candidate helper** | **Yes** — with dedicated candidate delegating to local legacy |
| **Risk if ignored** | Exchange stage-apply gates on ready state; wiring CPU-only would miss readiness regressions |

---

## hyperdensityObservedPodContainerRestartCount

| Attribute | Detail |
|-----------|--------|
| **Where** | `resource_exchange_stage_apply.go`, `resource_exchange_stage_apply_chain.go` |
| **Call sites (resource_exchange)** | **12** |
| **Definition** | `hyperdensity_parent_fabric_workload_helpers.go` (shared legacy) |
| **Observation class** | Container restart count from pod status |
| **In Sprint 79 shadow matrix** | **No** |
| **Future matrix** | **Yes** — with CPU + ready helpers before wired flip |
| **Safe for candidate helper** | **Yes** — delegate to workload_helpers legacy |
| **Risk if ignored** | Restart drift during exchange apply/rollback checks |

---

## Decision (Sprint 79–80)

| Decision | Value |
|----------|-------|
| Include in Sprint 79 shadow matrix | **No** |
| Include in Sprint 80 staged wrappers | **No** |
| Include in Sprint 81 shadow matrix | **Yes** (candidates only) |
| Include in Sprint 81 staged wrappers | **No** |
| Classify as local-observation-like | **Yes** |
| Copy to Hyperdensity | **No** |
| Use for broad observation | **No** |
| Defer to Sprint 80+ if wiring continues | **Yes** |

---

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_OBSERVATION_SHADOW_MATRIX.md`
- Dashboard `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_LOCAL_HELPER_CLASSIFICATION_M76.md`


---

## Sprint 83 (call-site wiring)

Sprint 83 wires all 32 production `resource_exchange_*` observation call sites to full-helper staged wrappers (8 CPU + 12 ready + 12 restart). `ResourceExchangeObservationWiredV1` and `ResourceExchangeObservationCandidateRuntimeUsedV1` remain **false**; effective runtime path is **legacy**. Direct candidate calls in production remain forbidden. See `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_CALLSITE_WIRING.md` (Hyperdensity) and `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_RESOURCE_EXCHANGE_CALLSITE_WIRING_M83.md` (Dashboard).
