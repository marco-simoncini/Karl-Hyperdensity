# Hyperdensity golden anchor — M1 (Sprint 17)

First **real extraction micro-step**: shared blocker catalog + minimal parent-fabric summary DTO + redacted golden JSON. **No Dashboard handler changes**, **no API response changes**.

## Artifacts

| Path | Role |
|------|------|
| `pkg/hyperdensity/blockers/` | Gate/blocker ID catalog (`Known`, `Severity`, `Catalog`) |
| `pkg/hyperdensity/contracts/` | `ParentFabricSummary` DTO + parse/validate helpers |
| `testdata/dashboard/parent_fabric_summary_redacted.golden.json` | Redacted M1 summary anchor |

## Blocker IDs (M1 minimum)

Aligned with Dashboard parent-fabric copy and collectors:

- `no_windows_lane`, `no_production_mutation`
- `keep_windows_lane_disabled`, `windows_disabled`
- `dry_run_only`, `runtime_apply_disabled`
- `unsupported_broad_vm_execution`, `unsupported_broad_memory_execution`
- `unsupported_multi_container_widening`, `unsupported_broad_automation`

Dashboard references (examples):

- `hyperdensity_parent_fabric_vm_runtime_evidence_collector_v1.go` — gate `no_windows_lane` → `windows_disabled` / `keep_windows_lane_disabled`
- `hyperdensity_parent_fabric_vm_lane_readiness_v1.go` — `no_production_mutation`
- Execution category `dry_run_only` across parent-fabric tests

## Golden summary rules

The golden JSON uses **`hyperdensity.karl.io/parent-fabric-summary/v1`** — a **Hyperdensity contract view**, not a byte-for-byte copy of Dashboard `HyperdensityParentFabricSummaryV1`. Fields are chosen for compatibility checks:

- Windows lane **disabled** with `no_windows_lane` in blockers
- `applyAllowed: false`, `operatorControlled: true`, `recommendationOnly: true`
- `kubeVirtLegacy.present: true`, `dryRunSupported: true`

## Validation

```bash
go test ./pkg/hyperdensity/...
./scripts/validate.sh
```

## Next (M2)

Optional: generate golden from a redacted live Dashboard `?view=summary` capture and diff field mapping table in Karl-Dashboard `docs/hyperdensity/HYPERDENSITY_GOLDEN_ANCHOR_M1.md`.
