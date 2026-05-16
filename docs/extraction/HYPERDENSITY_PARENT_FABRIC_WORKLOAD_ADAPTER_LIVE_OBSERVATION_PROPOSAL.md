# Hyperdensity Parent Fabric — live observation proposal (Sprint 58)

## Summary

**Live observation may be considered in a future sprint.** Sprint 58 **only inventories and proposes** — no live wiring, no new production flags enabled.

---

## Future allowed file

- **`hyperdensity_parent_fabric_live.go`** (read-only observation helpers only)

---

## Forbidden (unchanged)

| Area | Status |
|------|--------|
| `hyperdensity_parent_fabric_apply.go` | **Forbidden** |
| `hyperdensity_parent_fabric_resource_exchange_*` | **Forbidden** |
| Rollback observed-state | **Forbidden** |
| VM runtime observation | **Forbidden** |
| Admission guard | **Forbidden** |

---

## Future flags (not enabled in Sprint 58)

| Constant | Sprint 58 | Future |
|----------|-----------|--------|
| `hyperdensityWorkloadAdapterLiveObservationWiredV1` | **not defined / not true** | `true` when live phase approved |
| `hyperdensityWorkloadAdapterObservationWiredV1` | **`false`** | **`false`** until broad observation phase complete |
| `hyperdensityWorkloadAdapterPilotObservationWiredV1` | **`true`** | unchanged until pilot rollback |

---

## Proposed wrapper (future)

- `hyperdensityWorkloadLiveObservedStateV1` — **must not** appear in production until live sprint
- Inventory: Dashboard `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_LIVE_OBSERVATION_INVENTORY_M45.md`

---

## Sprint sequence

| Sprint | Work |
|--------|------|
| 57 | Pilot-only observation wiring |
| 58 | Hardening + live proposal (no new wiring) |
| 59+ | Live read-only observation (TBD) |

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_PILOT_OBSERVATION_HARDENING.md`
