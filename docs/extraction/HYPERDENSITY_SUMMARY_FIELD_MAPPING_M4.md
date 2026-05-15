# Hyperdensity summary field mapping — M4 canonical golden helpers (Sprint 20)

## Purpose

Sprint 20 adds **pure helpers** in `pkg/hyperdensity/contracts` for stable JSON of `ParentFabricSummary`:

- `CanonicalSummaryJSON` — canonical bytes for tests and extraction aids  
- `WriteCanonicalSummary` — optional write to disk (e.g. env-guarded golden refresh)  
- `CompareSummaryGolden` — compare an actual summary to golden JSON by **re-canonicalizing** both sides (whitespace-insensitive)

## Boundaries

- **Test-only / extraction aid:** these helpers exist so Dashboard (or other consumers) can pin contract shape in tests without touching cluster or HTTP. They are **not** a runtime API surface for the console or Parent Fabric handlers.
- **Does not replace Dashboard summary:** Dashboard `?view=summary` remains the native transport; the contract is a **mapped, claim-safe** projection for Hyperdensity alignment tests.
- **No `supportsApply` → `applyAllowed`:** Dashboard `executionEngine.supportsApply` must **not** set contract `applyAllowed` to `true` on the M1/M3/M4 anchors; mapping rules stay in `mapping.go` and mapper tests.

## JSON stability

- Output uses `encoding/json` **Encoder** with `SetEscapeHTML(false)` and **two-space** indent (`SetIndent("", "  ")`).
- Struct field order follows `ParentFabricSummary` in `summary.go` (default `encoding/json` ordering).
- `CanonicalSummaryJSON` returns bytes **without** a trailing newline (the encoder’s trailing newline is stripped); `WriteCanonicalSummary` appends a single `\n` for POSIX text files.

## Related

- M3 mapper: `HYPERDENSITY_SUMMARY_FIELD_MAPPING_M3.md`  
- Dashboard golden generator test: `HYPERDENSITY_GOLDEN_ANCHOR_M4.md` (Karl-Dashboard)
