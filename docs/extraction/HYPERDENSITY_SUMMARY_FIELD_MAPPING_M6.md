# Hyperdensity summary field mapping — M6 fixture policy & dual edge (Sprint 22)

## Covered edges (redacted live-like fixtures → `ParentFabricSummary`)

1. **`executionEngine.supportsApply: true`** with **`executionEngine.summary.category: dry_run_only`**  
   - Contract **`applyAllowed` remains `false`** (claim-safe mapping; see `mapping.go`).

2. **`executionEngine.supportsApply: false`** with **`executionEngine.summary.category: dry_run_only`**  
   - Contract **`applyAllowed` remains `false`**.  
   - **`dryRunSupported`** may be `true` in the contract **only** when the pre-map category is `dry_run_only`, matching `InferDryRunSupported` semantics.

## Helpers (Sprint 22)

- **`AllowedDashboardSummaryFixtureFields()`** — dotted paths allowed in committed Dashboard-shaped JSON fixtures (documentation + tests).
- **`ValidateContractClaimSafe`** — Windows disabled, `applyAllowed` false, plus `ValidateSummary` / `ValidateNoForbiddenClaims` / `ValidateApplySemantics`.
- **`ValidateSupportsApplyFalseEdge`** — asserts the `supportsApply=false` edge against metadata (`RedactedLiveSummaryMetadata`, including **`ExecutionSummaryCategory`** for the dry-run rule).

## Boundaries

- **Test-only:** fixtures and goldens are for extraction alignment and CI; not a runtime API.
- **No Dashboard / HTTP / Kubernetes** in these helpers.

## Related

- M5 redacted live validation: `HYPERDENSITY_SUMMARY_FIELD_MAPPING_M5.md`  
- Dashboard allowlist capture script + M6 anchor: `docs/hyperdensity/HYPERDENSITY_GOLDEN_ANCHOR_M6.md` (Karl-Dashboard)
