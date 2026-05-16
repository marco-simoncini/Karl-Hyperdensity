# Hyperdensity Parent Fabric — workload helpers copy deferred (Sprint 48)

## Decision

**`copy-deferred`** for `hyperdensity_parent_fabric_workload_helpers.go`.

No `helpers.go`, golden, or copy-contract in Karl-Hyperdensity for Sprint 48.

## Blockers

| Blocker | Detail |
|---------|--------|
| **Cross-package coupling** | 35+ functions call `server` symbols defined outside this file (execution actions, nested map helpers, pilot types, deployment state builder). |
| **KubeVirt / K8s paths** | Six functions build REST paths containing `kubevirt.io` or Kubernetes API prefixes as string literals. |
| **Runtime observation** | Functions read untyped `map[string]interface{}` workload/pod shapes and produce `HyperdensityPilotObservedState` — runtime-bound, not DTO-only. |
| **Apply semantics** | Execution mode / ready-reason / mechanism helpers are part of the apply orchestration surface. |

## What would be required before copy

1. ~~Extract stdlib-only nested-map + quantity primitives~~ — **Sprint 49 done** (`parentfabric/primitives`); not wired to Dashboard yet.
2. Define **adapter** package in Dashboard for API paths and live observation (explicit sprint).
3. Re-audit with `go list -deps` on a narrowed function allowlist.

## Placeholder

`pkg/hyperdensity/parentfabric/workload/doc.go` — package reserved; **no** implementation.

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_HELPERS_AUDIT.md`
