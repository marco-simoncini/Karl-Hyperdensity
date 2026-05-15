# claimpolicy — Dashboard file traceability (Sprint 38)

## Purpose

Sprint 38 extends each **`SurfaceClaimMapping`** with **`DashboardFiles`**: relative paths under Karl-Dashboard **`kubernetes-console/`** pointing at real **`hyperdensity_parent_fabric_*`** builders or **`scripts/hyperdensity/`** audit artifacts. Helpers:

- `DashboardFilesForClaim(id)` — sorted, de-duplicated union of paths for a catalog claim.
- `ValidateDashboardFileTraceability()` — validates path shape, non-empty trace lists, and mapping invariants (includes `ValidateSurfaceMappings`).

## Rules

| Rule | Detail |
|------|--------|
| Path shape | Must start with `pkg/server/` or `scripts/hyperdensity/`; **no** absolute paths; **no** `..`. |
| Duplicates | **No** duplicate path within a single mapping row; union per claim must match unique path count. |
| Runtime | **`RuntimeImportAllowed`** remains **`false`** for every row. |
| Empty trace | Forbidden unless **`Notes`** documents **`future-only`** (Sprint 38 ships concrete paths for all current rows). |
| Enforcement | **None** — contract / test / documentation only; **no** API or Parent Fabric behavior change. |
| Schema | **`ContractKitVersion`** and manifest envelope **unchanged**; nested module semver bumps only. |

## Relationship to Sprint 37

Sprint 37 introduced surface tokens and mapping rows; Sprint 38 **anchors** those rows to on-disk Dashboard sources for reviewer traceability. Karl-Dashboard **runtime** still **must not** import `claimpolicy` (M17).

## Validation

```bash
( cd pkg/hyperdensity/contractkit && go test ./claimpolicy -count=1 )
./scripts/validate.sh
```

## Related

- `HYPERDENSITY_CLAIMPOLICY_SURFACE_MAPPING.md`
- `HYPERDENSITY_CONTRACTKIT_CLAIMPOLICY.md`
- `Karl-Dashboard/docs/hyperdensity/HYPERDENSITY_CLAIMPOLICY_TRACEABILITY_M21.md`
