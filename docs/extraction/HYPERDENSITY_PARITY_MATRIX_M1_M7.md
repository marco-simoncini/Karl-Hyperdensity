# Hyperdensity parity matrix — M1 through M7 (Sprint 24)

Consolidated view of **test-only / extraction-only** work aligning Karl-Dashboard parent-fabric `?view=summary` with Karl-Hyperdensity `ParentFabricSummary` contracts. **Dashboard runtime remains authoritative** for live API responses. **KHR is not an execution arm** — no apply, no Windows enablement, no handler imports.

| Milestone | Anchor | Dashboard artifact | Hyperdensity package/helper | Guarantees | Not covered |
|-----------|--------|--------------------|-----------------------------|------------|-------------|
| **M1** | Blocker catalog + golden summary | `hyperdensity_blocker_catalog_parity_test.go`; `docs/hyperdensity/HYPERDENSITY_GOLDEN_ANCHOR_M1.md` | `pkg/hyperdensity/blockers/`; `pkg/hyperdensity/contracts/summary.go`; `testdata/dashboard/parent_fabric_summary_redacted.golden.json` | 10 stable blocker IDs; `ParentFabricSummary` DTO; `ValidateSummary` / `ValidateNoForbiddenClaims`; Windows **disabled**; `kubeVirtLegacy.present: true`; `recommendationOnly` + `dryRunSupported` on anchor | Live Dashboard JSON byte parity; runtime handler import; production apply |
| **M2** | Apply semantics field mapping | `docs/hyperdensity/HYPERDENSITY_GOLDEN_ANCHOR_M2.md`; compile-only import in blocker parity test | `docs/extraction/HYPERDENSITY_SUMMARY_FIELD_MAPPING_M2.md` | Documents `supportsApply` ≠ `applyAllowed`; compile-time catalog import from Dashboard tests | Mapper tests; live fixtures |
| **M3** | Mapper helpers | `hyperdensity_summary_contract_mapper_test.go`; `testdata/hyperdensity_parent_fabric_summary_dashboard_redacted.json`; `docs/hyperdensity/HYPERDENSITY_GOLDEN_ANCHOR_M3.md` | `pkg/hyperdensity/contracts/mapping.go` (`BuildClaimSafeExecutionEngine`, `ValidateApplySemantics`, `MapSupportsApplyToContractApplyAllowed` → always false) | Test-only mapper; `supportsApply: true` does **not** set `applyAllowed: true`; go.mod tidy + direct Hyperdensity require | Runtime mapper; API shape change |
| **M4** | Canonical golden JSON | `hyperdensity_summary_golden_generator_test.go`; `testdata/hyperdensity_parent_fabric_summary_contract_expected.golden.json`; `docs/hyperdensity/HYPERDENSITY_GOLDEN_ANCHOR_M4.md` | `pkg/hyperdensity/contracts/golden.go` (`CanonicalSummaryJSON`, `CompareSummaryGolden`, `WriteCanonicalSummary`) | Stable two-space JSON; optional `KARL_UPDATE_HYPERDENSITY_GOLDEN=1`; claim-safe contract golden | Dashboard native summary bytes; HTTP capture |
| **M5** | Live redacted fixture validation | `hyperdensity_summary_live_redacted_test.go` (`TestHyperdensitySummaryLiveRedacted`); live redacted fixture + golden; `capture_parent_fabric_summary_redacted.sh` (manual); `docs/hyperdensity/HYPERDENSITY_GOLDEN_ANCHOR_M5.md` | `RedactedLiveSummaryMetadata`; `ValidateRedactedLiveSummaryFixture` | Redacted live-like fixture; source must contain `redacted` / `live-capture-redacted`; cluster-free tests | Automatic live curl in CI; secrets in repo |
| **M6** | Allowlist + supportsApply false edge | Allowlist-only capture script; `TestHyperdensitySummaryLiveRedactedSupportsApplyFalse`; second fixture + golden; `docs/hyperdensity/HYPERDENSITY_GOLDEN_ANCHOR_M6.md` | `AllowedDashboardSummaryFixtureFields`; `ValidateContractClaimSafe`; `ValidateSupportsApplyFalseEdge` | Fixture JSON limited to allowlisted paths; `supportsApply: false` + `dry_run_only` edge; `applyAllowed: false` | Full Dashboard schema; inverse allowlist enforcement in Hyperdensity alone |
| **M7** | Missing optional fields edge | `TestHyperdensitySummaryLiveRedactedMissingOptional`; missing-optional fixture + golden; `mapContractGeneratedAt`; `docs/hyperdensity/HYPERDENSITY_GOLDEN_ANCHOR_M7.md` | `MissingOptionalGeneratedAtDefault`; `ValidateMissingOptionalFieldsEdge`; extended `RedactedLiveSummaryMetadata` | `generatedAt: null` → `redacted-generatedAt-unavailable`; absent `decisionEngine` → counts **0**; **explicit** `kubeVirtLegacyRedacted.present: true` (not inferred) | Inferring KubeVirt from absence; real timestamps when redacted |

## Cross-cutting posture (M1–M7)

- **Test-only / extraction-only:** all milestones are parity tests, fixtures, docs, and pure helpers — no Dashboard handler or Parent Fabric runtime edits.
- **Dashboard runtime authoritative:** `GET /api/hyperdensity/parent-fabric?view=summary` remains owned by Dashboard Go handlers; Hyperdensity contract is a **claim-safe projection** for alignment.
- **Hyperdensity repo state:** first real **`blockers`** + **`contracts`** packages and golden/test helpers exist; suitable for **test dependency import** (already done in Dashboard `pkg/server` tests).
- **KHR is not execution arm:** no apply path, no Windows lane enablement, no new Parent Fabric, no Grande Padre expansion.
- **Windows:** remains **disabled / planning-only** on all anchors; blockers document disabled posture.
- **KubeVirt legacy:** **`present: true` required** on contract for M1–M7 anchors; fixtures must set explicitly — mapper does not assume legacy when block is absent.

## Related

- Extraction readiness: `HYPERDENSITY_M1_M7_EXTRACTION_READINESS.md`
- Dashboard parity runner: `docs/hyperdensity/HYPERDENSITY_PARITY_M1_M7.md` (Karl-Dashboard)
