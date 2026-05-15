# contractkit/blockers — Dashboard runtime import (Sprint 31–34 / M14)

## Boundary

`pkg/hyperdensity/contractkit/blockers` exposes **stable string constants** for gate/blocker IDs. Sprint 31 introduced the first Dashboard **production** import; Sprint 32 extends it to VM readonly observation; Sprint 33 completes catalog `no_production_mutation` on policy pack, consistency checker, release matrix, and live resource authority limitation surfaces; **Sprint 34** adds bash audits + **M17** runtime import freeze.

| Allowed in runtime | Still test-only |
|--------------------|-----------------|
| `contractkit/blockers` ID constants | `contractkit/contracts` (DTOs, manifest, golden) |

Constants are byte-identical to prior string literals — **no JSON or API change**.

## Surfaces (Dashboard)

- **Sprint 31:** VM lane readiness / evidence refresh / runtime evidence collector; guarded auto-execution ledger `no_production_mutation`.
- **Sprint 32:** VM readonly observation (intake, package, dry-run pass, submission gate, real submission policy remediation, runtime live probe remediation, operator submission / grant / approval preflights); acknowledgement and checklist strings that match catalog `no_windows_lane`.
- **Sprint 33:** Policy pack `RuleID`, policy-pack consistency `requiredSafetyGates`, release support matrix `LimitationID`, live resource authority `LimitationID` for catalog `no_production_mutation`.

Import alias used in Dashboard:

```go
hpblockers "github.com/marco-simoncini/Karl-Hyperdensity/pkg/hyperdensity/contractkit/blockers"
```

## Guards (Dashboard tests + scripts)

- No `contractkit/contracts` in production `.go`.
- `TestHyperdensityRuntimeContractkitImportWhitelist` + `contractkit_runtime_import_allowlist.txt`: only listed production files may import `blockers`; only exact import path `.../pkg/hyperdensity/contractkit/blockers` is permitted under `contractkit/`.
- Parity runner: `audit_contractkit_runtime_imports.sh`, `audit_hyperdensity_blocker_literals.sh` (see M17).

## Not in scope

- Execution/apply paths, handler flow, Parent Fabric assembly order.
- Replacing `dry_run_only` where it denotes **execution category** or **SupportLevel** (same string, different layer).
- Replacing preflight check names like `windows_lane_disabled` (not catalog `windows_disabled`).
- New Kubernetes clients, cluster IO, or npm.

## Version

Consumers pin tagged module `v0.1.5-khr-m1-m16` (Sprint 38 adds `DashboardFiles` traceability on claimpolicy surface mappings; runtime remains `blockers` only per M17).

## Sprint 34 — freeze (Dashboard M17)

Runtime production import from `contractkit` is **closed** at `contractkit/blockers` only. Automated audits and an explicit importer allowlist live under `Karl-Dashboard/kubernetes-console/scripts/hyperdensity/`. See **`Karl-Dashboard/docs/hyperdensity/HYPERDENSITY_CONTRACTKIT_RUNTIME_IMPORT_FREEZE_M17.md`**.

## Next slices

- Optional: other catalog-aligned IDs in policy matrix rows where strings duplicate catalog verbatim (beyond `no_production_mutation`).
