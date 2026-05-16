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

1. **Dashboard adapter** for KubeVirt/K8s API path builders (explicit sprint).
2. **Observed-state** builders remain runtime-bound — not in primitives.
3. Optional alignment sprint: map Dashboard nested helpers → Hyperdensity `primitives` API (different signatures today).

## Rules (unchanged)

- Dashboard **runtime owner**
- **No** production import of `pkg/hyperdensity/parentfabric`
- **No** API / JSON ordering / apply changes

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_HELPERS_DEFERRED.md`
- `HYPERDENSITY_PARENT_FABRIC_PRIMITIVES_CONTRACT.md`
