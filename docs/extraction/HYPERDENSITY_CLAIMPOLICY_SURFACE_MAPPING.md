# claimpolicy — surface mapping (Sprint 37)

## Purpose

Sprint 37 adds **`SurfaceClaimMapping`** rows (see `pkg/hyperdensity/contractkit/claimpolicy/surface_mapping.go`) that document how each **claim-policy catalog ID** aligns with **conceptual Parent Fabric / Karl-Dashboard builder surfaces** (`ParentFabricSurface`).

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
| `SurfaceRuntimeImportFreeze` | `runtime_import_freeze` | M17: `contractkit/blockers` only in runtime `pkg/server`. |
| `SurfaceHyperdensityRecommendation` | `hyperdensity_recommendation` | Recommendation-only surfaces without apply authority. |

## Sprint 37 rules

- **`RuntimeImportAllowed`:** always **`false`** for every mapping row (claimpolicy remains test-only on Dashboard).
- **Every catalog `ClaimID`:** must appear in at least one mapping (`ValidateSurfaceMappings`).
- **KubeVirt:** `kubevirt_legacy_provider` (compatibility marker) and `no_generic_kubevirt_replacement` (forbidden narrative) use **distinct** `Field` / semantics.
- **Windows:** `no_windows_hyperdensity_apply` maps to **Windows lane** with apply **disabled**; `windows_lane_disabled` maps to preflight check name vocabulary.
- **Schema / manifest:** `ContractKitVersion` and `FixtureManifestVersion` stay on Sprint 26 / M9 anchors — module semver only bumps.

## API

| Function | Role |
|----------|------|
| `SurfaceMappings()` | All rows, stable sort order. |
| `MappingsForClaim(id)` | Filter by `ClaimID`. |
| `ValidateSurfaceMappings()` | Invariants for tests and Dashboard parity. |

## Validation

```bash
( cd pkg/hyperdensity/contractkit && go test ./claimpolicy -count=1 )
./scripts/validate.sh
```

## Related

- `HYPERDENSITY_CONTRACTKIT_CLAIMPOLICY.md`
- `Karl-Dashboard/docs/hyperdensity/HYPERDENSITY_CLAIMPOLICY_SURFACE_MAPPING_M20.md`
