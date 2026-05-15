# claimpolicy — Dashboard required token guard (Sprint 39)

## Purpose

Sprint 39 extends **`SurfaceClaimMapping`** with **`DashboardRequiredTokens`**: stable, sorted, unique **substring** tokens that must appear in Karl-Dashboard sources **for mechanical parity in Dashboard tests only**.

Hyperdensity **`claimpolicy`** does **not** read Dashboard files. It only:

- Declares expected tokens next to **`DashboardFiles`** on each mapping row.
- Validates shape and invariants via **`ValidateDashboardRequiredTokens()`** and **`ValidateSurfaceMappings()`** (no filesystem I/O in this repo).

Real file content checks run in **Karl-Dashboard** `TestHyperdensityClaimpolicyTraceabilityParity` (see `HYPERDENSITY_CLAIMPOLICY_TRACEABILITY_TOKEN_GUARD_M22.md` on Dashboard).

## Rules

| Rule | Detail |
|------|--------|
| Non-empty | Any row with **`DashboardFiles`** must list a non-empty **`DashboardRequiredTokens`** slice. |
| Normalization | Tokens are trimmed, de-duplicated per row, sorted lexicographically; union per claim via **`RequiredTokensForClaim`** is sorted and unique. |
| Substrings only | No regex; tokens are plain substrings searched with `strings.Contains`. |
| No secrets / no env | Tokens are public vocabulary only (e.g. gate ids, lane names). |
| No path-like tokens | Tokens must not contain `/`, `\`, or `..` (conservative guard). |
| Runtime | **`RuntimeImportAllowed`** remains **`false`** for every row. M17 runtime import freeze unchanged. |
| API / payload | **No** change to HTTP APIs, JSON field ordering, or Parent Fabric runtime behavior. |
| Schema | **`ContractKitVersion`** (`v0.0.0-sprint26`) and **`FixtureManifestVersion`** unchanged; module semver bumps to **`v0.1.7-khr-m1-m18`**. |

## Version note

- **`v0.1.5-khr-m1-m16`**: first Sprint 38 traceability tag (superseded).
- **`v0.1.6-khr-m1-m17`**: corrected `windows_lane_disabled` file paths for traceability.
- **`v0.1.7-khr-m1-m18`**: adds **`DashboardRequiredTokens`** + contract-side validation + Dashboard mechanical token parity. The live-resource-authority limitation row traces **`hpblockers.IDNoProductionMutation`** (substring **`IDNoProductionMutation`**) rather than spelling **`no_production_mutation`** as a string literal in that file — see **`Notes`** on that mapping row in `surface_mapping.go`.

## Related

- `HYPERDENSITY_CLAIMPOLICY_SURFACE_TRACEABILITY.md`
- `HYPERDENSITY_CLAIMPOLICY_SURFACE_MAPPING.md`
- `HYPERDENSITY_CONTRACTKIT_CLAIMPOLICY.md`
