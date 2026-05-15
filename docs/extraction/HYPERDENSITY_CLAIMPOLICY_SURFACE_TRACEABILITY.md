# claimpolicy — Dashboard file traceability (Sprint 38–39)

## Purpose

Sprint 38 extends each **`SurfaceClaimMapping`** with **`DashboardFiles`**: relative paths under Karl-Dashboard **`kubernetes-console/`** pointing at real **`hyperdensity_parent_fabric_*`** builders or **`scripts/hyperdensity/`** audit artifacts. **Sprint 39** adds **`DashboardRequiredTokens`** on each row (substring expectations for Dashboard test-only parity). Helpers:

- `DashboardFilesForClaim(id)` — sorted, de-duplicated union of paths for a catalog claim.
- `RequiredTokensForClaim(id)` — Sprint 39: sorted, de-duplicated union of required tokens for a claim.
- `ValidateDashboardRequiredTokens()` — Sprint 39: token invariants (still no filesystem reads in Hyperdensity).
- `ValidateDashboardFileTraceability()` — validates path shape, non-empty trace lists, mapping invariants, and required tokens (via `ValidateSurfaceMappings`).

## Rules

| Rule | Detail |
|------|--------|
| Path shape | Must start with `pkg/server/` or `scripts/hyperdensity/`; **no** absolute paths; **no** `..`. |
| Duplicates | **No** duplicate path within a single mapping row; union per claim must match unique path count. |
| Runtime | **`RuntimeImportAllowed`** remains **`false`** for every row. |
| Empty trace | Forbidden unless **`Notes`** documents **`future-only`** (Sprint 38 ships concrete paths for all current rows). |
| Required tokens | Sprint 39: every row with **`DashboardFiles`** lists non-empty **`DashboardRequiredTokens`** (sorted, unique, non–path-like). |
| Enforcement | **None** — contract / test / documentation only; **no** API or Parent Fabric behavior change. |
| Schema | **`ContractKitVersion`** and manifest envelope **unchanged**; nested module semver bumps to **`v0.1.7-khr-m1-m18`** for Sprint 39. |

## Relationship to Sprint 37

Sprint 37 introduced surface tokens and mapping rows; Sprint 38 **anchors** those rows to on-disk Dashboard sources for reviewer traceability. Karl-Dashboard **runtime** still **must not** import `claimpolicy` (M17).

## Validation

```bash
( cd pkg/hyperdensity/contractkit && go test ./claimpolicy -count=1 )
./scripts/validate.sh
```

## Related

- `HYPERDENSITY_CLAIMPOLICY_SURFACE_MAPPING.md`
- `HYPERDENSITY_CLAIMPOLICY_TRACEABILITY_TOKEN_GUARD.md`
- `HYPERDENSITY_CONTRACTKIT_CLAIMPOLICY.md`
- `Karl-Dashboard/docs/hyperdensity/HYPERDENSITY_CLAIMPOLICY_TRACEABILITY_M21.md`
