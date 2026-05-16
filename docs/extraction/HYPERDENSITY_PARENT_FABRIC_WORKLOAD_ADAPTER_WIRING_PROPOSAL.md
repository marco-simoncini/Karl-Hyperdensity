# Hyperdensity Parent Fabric — workload adapter wiring proposal (Sprint 55)

## Status

**Sprint 55 was proposal-only.** **Sprint 56** executes **path-only** wiring. **Sprint 57** executes **pilot-only** observed-state wiring. General `ProductionWiredV1` remains **false**. — see **`HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_PATH_WIRING.md`**. **`hyperdensityWorkloadAdapterProductionWiredV1`** remains **`false`** (sub-flag `PathWiredV1` only). **`hyperdensityWorkloadAdapterProductionWiredV1`** remains **`false`**. Full **`workload_helpers.go`** verdict remains **`copy-deferred`**.

**Sprint 54 shadow PASS does not authorize automatic wiring.** A dedicated **wiring sprint** with explicit approval is required.

---

## 1. Scope

| In scope (future wiring sprint) | Owner |
|---------------------------------|-------|
| Replace legacy **path** helper calls with `hyperdensityWorkloadPathAdapterV1` behind a gate | Karl-Dashboard |
| Expand shadow / canary tests per call-site file | Karl-Dashboard |
| Keep Hyperdensity `parentfabric/workload` at **3 pure-candidates only** | Karl-Hyperdensity |

| Out of this proposal’s first phase | Notes |
|------------------------------------|-------|
| Observed-state call-site replacement | Phase 2 after path-only stable |
| Execution / apply mode helpers | **Excluded** — remain legacy |
| Hyperdensity Go adapter code | **No** — Dashboard-only |

---

## 2. Non-goals

- No API response / payload / JSON ordering change in Sprint 55 or implied auto-wiring.
- No Dashboard production import of `pkg/hyperdensity/parentfabric`.
- No full `workload_helpers.go` copy to Hyperdensity.
- No KubeVirt removal; no Windows enablement.
- No change to execution/apply orchestration in wiring phase 1.

---

## 3. Proposed sequence

| Phase | Content | Gate |
|-------|---------|------|
| **3a — Path-only** | Wire `AppsWorkloadPath`, `PodPath`, `PodResizePath`, VM/VMI/GuestOS paths | `hyperdensityWorkloadAdapterPathWiredV1` or flip path sub-flag |
| **3b — Observed-state shadow** | Keep shadow tests; add per-file canary comparing legacy vs adapter on live code paths (test/build tag) | Shadow PASS + inventory reviewed |
| **3c — Observed-state wiring** | Wire `ExtractWorkloadObservedState`, pod snapshot, StatefulSet state | Separate flag; pilot.go first |
| **3d — Apply/execution** | **Last or never** | Do-not-wire classification |

**Sequencing rule:** path-only first; observed-state only after path phase green in CI and manual review.

---

## 4. Feature flag / constant gate recommendation

| Constant | Sprint 55 value | Future wiring sprint |
|----------|-----------------|----------------------|
| `hyperdensityWorkloadAdapterProductionWiredV1` | `false` | `true` only when **all** approved phases complete OR split into sub-flags |
| Recommended addition | — | `hyperdensityWorkloadAdapterPathWiredV1 bool` for phase 3a only |
| Recommended addition | — | `hyperdensityWorkloadAdapterObservationWiredV1 bool` for phase 3c |

**Pattern:** each call site uses:

```go
if hyperdensityWorkloadAdapterPathWiredV1 {
    path, ok := hyperdensityWorkloadPathAdapterV1{}.AppsWorkloadPath(...)
} else {
    path, err := hyperdensityAppsWorkloadAPIPath(...)
}
```

Until wiring sprint, **no** such branches in production.

---

## 5. Rollback

| Step | Action |
|------|--------|
| 1 | Set path/observation wired flags to `false` |
| 2 | Revert call-site replacements (single commit revert per phase) |
| 3 | Keep `hyperdensity_parent_fabric_workload_adapter_v1.go` (adapter shell remains) |
| 4 | Run shadow tests + parity — must PASS on legacy path |

Rollback must not remove adapter file or shadow tests.

---

## 6. Required tests before wiring

| Test | Status (Sprint 55) |
|------|-------------------|
| `TestHyperdensityParentFabricWorkloadAdapterShadow` | **PASS** |
| `TestHyperdensityParentFabricWorkloadAdapterWiringGuard` | **PASS** |
| `audit_workload_adapter_call_sites.sh` | **PASS** |
| Future: path-only golden per wired file | Not started |
| Future: observed-state canary per pilot/live | Not started |

---

## 7. Risks

| Risk | Mitigation |
|------|------------|
| Apply/execution coupling | Path-only first; exclude apply.go phase 1 or wire read-only paths only |
| `time.Now()` in adapter observation methods | Use `extract*At` with explicit `now` at wired call sites |
| Shadow samples ≠ production edge cases | Per-file inventory (M41) + phased rollout |
| Accidental full-file wiring | Wiring guard test + audit script in parity |

---

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_SHADOW_TESTS.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_HARDENING.md`
- Dashboard `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_WIRING_PROPOSAL_M42.md`
- Dashboard `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_ADAPTER_CALLSITE_INVENTORY_M41.md`
