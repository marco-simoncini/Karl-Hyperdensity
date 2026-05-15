# Parent-fabric summary field mapping — M3 (Sprint 19)

Extends [M2 apply semantics](./HYPERDENSITY_SUMMARY_FIELD_MAPPING_M2.md) with **test-only mapper** helpers used by Karl-Dashboard parity tests.

## Helpers (`pkg/hyperdensity/contracts/mapping.go`)

| Function | Role |
|----------|------|
| `MapSupportsApplyToContractApplyAllowed` | Always returns `false` for claim-safe contract (ignores technical `supportsApply`) |
| `BuildClaimSafeExecutionEngine` | Builds `ExecutionEngineSummary` from Dashboard category + supportsApply |
| `InferDryRunSupported` | `dry_run_only` category → `dryRunSupported: true` |
| `ValidateApplySemantics` | Ensures `supportsApply` did not leak into `applyAllowed` |
| `IsClaimSafeApplyAllowed` | Validates applyAllowed vs operatorControlled |
| `NormalizeExecutionMode` | Stable mode label (default `operator_controlled`) |

## Dashboard test mapper

`hyperdensity_summary_contract_mapper_test.go` maps redacted fixture:

- `decisionEngine.eligibleYielderCount` → `parentPool.donorCount`
- `decisionEngine.eligibleReceiverCount` → `parentPool.receiverCount`
- `executionEngine.supportsApply` → **not** `applyAllowed` (via helpers)
- `executionEngine.summary.category` → `dryRunSupported`
- `windowsLaneRedacted` → `windowsLane`
- `kubeVirtLegacyRedacted` → `kubeVirtLegacy`

No runtime handler imports these helpers in Sprint 19.
