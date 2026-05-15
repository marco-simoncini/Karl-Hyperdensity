# claimpolicy — surface mapping (Sprint 37–39)

## Purpose

Sprint 37 adds **`SurfaceClaimMapping`** rows (see `pkg/hyperdensity/contractkit/claimpolicy/surface_mapping.go`) that document how each **claim-policy catalog ID** aligns with **conceptual Parent Fabric / Karl-Dashboard builder surfaces** (`ParentFabricSurface`).

**Sprint 38** adds **`DashboardFiles`** on each row: relative paths under `kubernetes-console/` to real `hyperdensity_parent_fabric_*` sources or `scripts/hyperdensity/` audit scripts — see **`HYPERDENSITY_CLAIMPOLICY_SURFACE_TRACEABILITY.md`**.

**Sprint 39** adds **`DashboardRequiredTokens`**: stable substring expectations per row; validated in Hyperdensity without reading Dashboard files; mechanical content checks run in Dashboard tests — see **`HYPERDENSITY_CLAIMPOLICY_TRACEABILITY_TOKEN_GUARD.md`**.

This is a **contract / test / documentation** artifact only:

- **No** runtime enforcement in `claimpolicy`.
- **No** Karl-Dashboard production import of `claimpolicy` (M17 unchanged).
- **No** API payload, JSON ordering, or Parent Fabric behavior change.

## Surfaces (minimal set)

| Surface constant | String value | Typical Parent Fabric alignment |
|------------------|--------------|----------------------------------|
| `SurfaceExecutionEngine` | `execution_engine` | Apply / exchange / governance posture builders. |
| `SurfaceWindowsLane` | `windows_lane` | VM readonly observation, Windows lane planning-only. |
| `SurfaceKubeVirtLegacyProvider` | `kubevirt_legacy_provider` | Legacy KubeVirt provider markers vs replacement narrative. |
| `SurfacePolicyPack` | `policy_pack` | `RuleID` rows aligned with catalog tokens. |
| `SurfaceReleaseSupportMatrix` | `release_support_matrix` | `LimitationID` rows. |
| `SurfaceLiveResourceAuthority` | `live_resource_authority` | Live authority limitation rows (e.g. `no_production_mutation`). |
| `SurfaceRuntimeImportFreeze` | `runtime_import_freeze` | M17: `contractkit/blockers` only in runtime `pkg/server`. |
| `SurfaceHyperdensityRecommendation` | `hyperdensity_recommendation` | Recommendation-only surfaces without apply authority. |

## Sprint 37–39 rules

- **`RuntimeImportAllowed`:** always **`false`** for every mapping row (claimpolicy remains test-only on Dashboard).
- **Every catalog `ClaimID`:** must appear in at least one mapping (`ValidateSurfaceMappings`).
- **`DashboardFiles`:** non-empty unless **`Notes`** documents **`future-only`**; paths are **relative** (`pkg/server/...` or `scripts/hyperdensity/...`); validated by `ValidateDashboardFileTraceability`.
- **`DashboardRequiredTokens`:** non-empty whenever **`DashboardFiles`** is non-empty; sorted, unique per row; no path-like tokens; validated by `ValidateDashboardRequiredTokens` (included from `ValidateSurfaceMappings` / `ValidateDashboardFileTraceability`).
- **KubeVirt:** `kubevirt_legacy_provider` (compatibility marker) and `no_generic_kubevirt_replacement` (forbidden narrative) use **distinct** `Field` / semantics.
- **Windows:** `no_windows_hyperdensity_apply` maps to **Windows lane** with apply **disabled**; `windows_lane_disabled` maps to preflight check name vocabulary.
- **Schema / manifest:** `ContractKitVersion` and `FixtureManifestVersion` stay on Sprint 26 / M9 anchors — module semver only bumps (Sprint 39 **`v0.1.8-khr-m1-m18`**, Sprint 40 **`v0.1.9-khr-m1-m19`**).

## API

| Function | Role |
|----------|------|
| `SurfaceMappings()` | All rows, stable sort order. |
| `MappingsForClaim(id)` | Filter by `ClaimID`. |
| `ValidateSurfaceMappings()` | Invariants for tests and Dashboard parity (paths + tokens, Sprint 39). |
| `DashboardFilesForClaim(id)` | Sprint 38: sorted unique traced paths for a claim. |
| `RequiredTokensForClaim(id)` | Sprint 39: sorted unique union of required substring tokens for a claim. |
| `ValidateDashboardRequiredTokens()` | Sprint 39: token invariants (no Dashboard filesystem reads). |
| `ValidateDashboardFileTraceability()` | Sprint 38–39: path + token trace invariants (delegates to `ValidateSurfaceMappings`). |

## Validation

```bash
( cd pkg/hyperdensity/contractkit && go test ./claimpolicy -count=1 )
./scripts/validate.sh
```

## Related

- `HYPERDENSITY_CLAIMPOLICY_SURFACE_TRACEABILITY.md`
- `HYPERDENSITY_CLAIMPOLICY_TRACEABILITY_TOKEN_GUARD.md`
- `HYPERDENSITY_CONTRACTKIT_CLAIMPOLICY.md`
- `Karl-Dashboard/docs/hyperdensity/HYPERDENSITY_CLAIMPOLICY_SURFACE_MAPPING_M20.md`
- `Karl-Dashboard/docs/hyperdensity/HYPERDENSITY_CLAIMPOLICY_TRACEABILITY_M21.md`
