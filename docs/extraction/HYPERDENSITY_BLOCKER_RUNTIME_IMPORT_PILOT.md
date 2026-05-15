# contractkit/blockers — Dashboard runtime import (Sprint 31–32 / M14)

## Boundary

`pkg/hyperdensity/contractkit/blockers` exposes **stable string constants** for gate/blocker IDs. Sprint 31 introduced the first Dashboard **production** import; Sprint 32 extends it to VM readonly observation surfaces.

| Allowed in runtime | Still test-only |
|--------------------|-----------------|
| `contractkit/blockers` ID constants | `contractkit/contracts` (DTOs, manifest, golden) |

Constants are byte-identical to prior string literals — **no JSON or API change**.

## Surfaces (Dashboard)

- **Sprint 31:** VM lane readiness / evidence refresh / runtime evidence collector; guarded auto-execution ledger `no_production_mutation`.
- **Sprint 32:** VM readonly observation (intake, package, dry-run pass, submission gate, real submission policy remediation, runtime live probe remediation, operator submission / grant / approval preflights); acknowledgement and checklist strings that match catalog `no_windows_lane`.

Import alias used in Dashboard:

```go
hpblockers "github.com/marco-simoncini/Karl-Hyperdensity/pkg/hyperdensity/contractkit/blockers"
```

## Guards (Dashboard tests)

- No `contractkit/contracts` in production `.go`.
- `TestHyperdensityRuntimeContractkitImportWhitelist`: only exact import path `.../pkg/hyperdensity/contractkit/blockers` is permitted under `contractkit/`.

## Not in scope

- Execution/apply paths, handler flow, Parent Fabric assembly order.
- Replacing `dry_run_only` where it denotes **execution category** or **SupportLevel** (same string, different layer).
- Replacing preflight check names like `windows_lane_disabled` (not catalog `windows_disabled`).
- New Kubernetes clients, cluster IO, or npm.

## Version

Consumers stay on tagged module `v0.1.1-khr-m1-m12` unless a future sprint bumps contractkit for catalog changes.

## Next slices

- Policy/release matrix `LimitationID` / `RuleID` rows for `no_production_mutation`.
- Optional: broaden whitelist to document other allowed third-party imports (currently only enforces Hyperdensity `contractkit/*` subtree).
