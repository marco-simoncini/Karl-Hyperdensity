# Hyperdensity Parent Fabric — workload helper prerequisites (Sprint 49)

## Purpose

Unblock a **future** re-audit of `hyperdensity_parent_fabric_workload_helpers.go` **without** copying it in Sprint 49.

## Delivered (Sprint 49)

| Item | Status |
|------|--------|
| `parentfabric/primitives` stdlib package | **Done** |
| Golden `primitives_contract.golden.json` | **Done** |
| Dashboard primitive loci audit (M34) | **Done** |
| `workload_helpers` copy | **Still deferred** |

## Still required before workload copy

1. **Dashboard adapter implementation** for path + observation (Sprint 50 documents boundary only).
2. **Observed-state** builders remain runtime-bound — classified in Dashboard M35 fixture.
3. **Re-audit criteria** met (see **`HYPERDENSITY_PARENT_FABRIC_WORKLOAD_REAUDIT_CRITERIA.md`**) — then narrow copy to pure allowlist (3 functions).

## Rules (unchanged)

- Dashboard **runtime owner**
- **No** production import of `pkg/hyperdensity/parentfabric`
- **No** API / JSON ordering / apply changes

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_HELPERS_DEFERRED.md`
- `HYPERDENSITY_PARENT_FABRIC_PRIMITIVES_CONTRACT.md`
