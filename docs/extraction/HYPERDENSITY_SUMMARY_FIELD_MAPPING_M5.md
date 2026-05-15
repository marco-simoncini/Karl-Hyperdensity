# Hyperdensity summary field mapping — M5 redacted live capture validation (Sprint 21)

## Purpose

Sprint 21 adds **pure** helpers for validating a **redacted** live (or live-like) `?view=summary` projection after it is mapped to `ParentFabricSummary`:

- `RedactedLiveSummaryMetadata` — carries `DashboardSupportsApply` from the Dashboard JSON **before** mapping, so `ValidateApplySemantics` can be checked against the same capture.
- `ValidateRedactedLiveSummaryFixture` — enforces source markers (`redacted` or `live-capture-redacted`), Windows lane disabled, and `ValidateSummary` / `ValidateNoForbiddenClaims` / `ValidateApplySemantics`.

## Boundaries

- **Test-only / extraction pipeline:** no runtime API, no cluster access, no HTTP clients in Karl-Hyperdensity.
- **No secrets:** validators only inspect the contract shape and markers; callers must redact captures before committing fixtures (see Karl-Dashboard `HYPERDENSITY_GOLDEN_ANCHOR_M5.md`).
- **Does not replace Dashboard summary:** native `GET .../parent-fabric?view=summary` response format stays owned by Dashboard; this package validates the **mapped** contract for alignment tests.

## Related

- M4 canonical JSON: `HYPERDENSITY_SUMMARY_FIELD_MAPPING_M4.md`
- Dashboard live-redacted fixture test: `docs/hyperdensity/HYPERDENSITY_GOLDEN_ANCHOR_M5.md`
