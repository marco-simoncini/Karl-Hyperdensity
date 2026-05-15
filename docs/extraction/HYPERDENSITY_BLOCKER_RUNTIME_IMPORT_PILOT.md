# contractkit/blockers — Dashboard runtime import pilot (Sprint 31 / M14)

## Boundary

`pkg/hyperdensity/contractkit/blockers` exposes **stable string constants** for gate/blocker IDs. Sprint 31 is the first Dashboard **production** import of this package.

| Allowed in runtime | Still test-only |
|--------------------|-----------------|
| `contractkit/blockers` ID constants | `contractkit/contracts` (DTOs, manifest, golden) |

Constants are byte-identical to prior string literals — **no JSON or API change**.

## Pilot surfaces (Dashboard)

VM lane readiness / evidence refresh / runtime evidence collector, plus guarded auto-execution ledger safety gate `no_production_mutation`.

Import alias used in Dashboard:

```go
hpblockers "github.com/marco-simoncini/Karl-Hyperdensity/pkg/hyperdensity/contractkit/blockers"
```

## Not in scope

- Execution/apply paths, handler flow, Parent Fabric assembly order.
- Replacing `dry_run_only` where it denotes **execution category** or **SupportLevel** (same string, different layer).
- Replacing preflight check names like `windows_lane_disabled` (not catalog `windows_disabled`).
- New Kubernetes clients, cluster IO, or npm.

## Version

Consumers stay on tagged module `v0.1.1-khr-m1-m12` unless a future sprint bumps contractkit for catalog changes.

## Next slices (documented, not implemented)

- Readonly observation VM surfaces (`no_windows_lane`, `keep_windows_lane_disabled`).
- Policy/release matrix `LimitationID` / `RuleID` rows for `no_production_mutation`.
- Optional: runtime grep/CI guard listing allowed production import paths.
