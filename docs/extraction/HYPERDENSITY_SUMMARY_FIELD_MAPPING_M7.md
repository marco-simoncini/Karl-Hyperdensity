# Hyperdensity summary field mapping — M7 missing optional fields edge (Sprint 23)

## Purpose

Third **test-only** fixture edge: redacted live-like Dashboard JSON with **optional fields absent or null**, mapped to a claim-safe `ParentFabricSummary`.

## Covered optional absences

| Dashboard field | Fixture behavior | Contract / mapper default |
|-----------------|------------------|---------------------------|
| `generatedAt` | `null`, absent, or `""` | `generatedAt`: **`redacted-generatedAt-unavailable`** (`MissingOptionalGeneratedAtDefault`) |
| `decisionEngine` | absent or counts missing | `parentPool.donorCount` / `receiverCount`: **`0`** |
| `kubeVirtLegacyRedacted.providerMode` | absent or empty | `providerMode` omitted or empty (no claim beyond legacy presence) |
| `kubeVirtLegacyRedacted.present` | **must be explicit `true`** | **Not inferred** when the block is absent — see below |

## KubeVirt legacy (M1–M7)

- **`kubeVirtLegacy.present: true` is required** on the contract for all M1–M7 anchors (`ValidateSummary`).
- Missing-optional fixtures **must set `kubeVirtLegacyRedacted.present: true` explicitly** in JSON.
- The test-only mapper **does not** default `present=true` when the block is missing (avoids implying live legacy posture without evidence).

## Helper

- **`ValidateMissingOptionalFieldsEdge`** — uses `RedactedLiveSummaryMetadata` (`DashboardGeneratedAtUnavailable`, `DashboardCountsAbsent`, `DashboardSupportsApply`, `ExecutionSummaryCategory`) plus standard summary validators.

## Boundaries

- **Test-only:** no runtime API, no cluster, no Dashboard imports in Karl-Hyperdensity.
- **Claim-safe:** `applyAllowed` remains **false**; Windows lane **disabled**.

## Related

- M6 dual edge + allowlist: `HYPERDENSITY_SUMMARY_FIELD_MAPPING_M6.md`
- Dashboard M7 anchor: `docs/hyperdensity/HYPERDENSITY_GOLDEN_ANCHOR_M7.md` (Karl-Dashboard)
