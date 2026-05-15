# Hyperdensity contractkit — no republish / no repoint policy (Sprint 40)

## Purpose

This document formalizes **release hygiene** for the nested Go module  
`github.com/marco-simoncini/Karl-Hyperdensity/pkg/hyperdensity/contractkit`.

It complements `HYPERDENSITY_CONTRACTKIT_RELEASE_TAGGING.md` and `HYPERDENSITY_CONTRACTKIT_VERSION_MODEL.md`. **No** HTTP API, Parent Fabric runtime behavior, JSON ordering, or manifest schema epoch changes are implied.

## Immutable module versions (Go ecosystem)

Once a **module version string** (for example `v0.1.7-khr-m1-m18`) has been **fetched by the public Go module proxy** (`proxy.golang.org`) or recorded in consumer `go.sum` files, the **zip content for that version is treated as immutable**.

Therefore:

- **Never repoint** a published git tag that names a released module version (same semver string, different commit).
- **Never delete and recreate** a tag for a version that may already have been ingested by the proxy or downstream builds.
- If a release is **wrong or incomplete**, publish a **strictly newer** module semver (new tag + updated `ContractKitModuleVersion` / `ContractKitGitTag`), and update consumer pins.

## Superseded module versions (historical)

| Version | Reason |
|---------|--------|
| `v0.1.5-khr-m1-m16` | First Sprint 38 traceability slice; superseded by **`v0.1.6-khr-m1-m17`**. |
| `v0.1.7-khr-m1-m18` | Short-lived tag; proxy immutability required a new semver; superseded by **`v0.1.8-khr-m1-m18`** (Sprint 39). |

These strings are also listed in code as `contracts.ContractKitSupersededModuleVersions` (test-only guards).

## Current stable (Sprint 40)

- **Before Sprint 40:** consumer pin **`v0.1.8-khr-m1-m18`** was current stable after Sprint 39.
- **After Sprint 40:** **`v0.1.9-khr-m1-m19`** is current stable — adds release-hygiene metadata (`ContractKitCurrentStableModuleVersion`, `IsSupersededModuleVersion`, `CurrentStableReleaseInfo`) **without** changing `ContractKitVersion` or manifest envelope.

## Why not `v0.1.7-khr-m1-m18`

That version string may still resolve to **stale zip content** on the public module mirror even if the git tag was moved or deleted on GitHub. **Do not pin** `v0.1.7-khr-m1-m18` in new work.

## If a release is wrong

1. Fix the code on branch `KHR` (or appropriate branch).  
2. Bump **`ContractKitModuleVersion`** / **`ContractKitGitTag`** to a **new** semver.  
3. Tag and push **`pkg/hyperdensity/contractkit/<new-version>`**.  
4. Update Karl-Dashboard `kubernetes-console/go.mod` and parity docs.  
5. Extend **`ContractKitSupersededModuleVersions`** only when a published semver must be permanently avoided (rare; prefer never publishing broken tags).

## Consumer CI environment (Sprint 41)

CI and private consumers should set **`GOPRIVATE=github.com/marco-simoncini/*`** and **`GONOSUMDB=github.com/marco-simoncini/*`** so nested `contractkit` tags resolve without spurious sumdb/proxy failures. Optionally **`GONOPROXY=github.com/marco-simoncini/*`** for direct VCS fetch. Prefer scoped **`GONOSUMDB`** over global **`GOSUMDB=off`**. Full checklist: **`HYPERDENSITY_CONTRACTKIT_CONSUMER_CI_HARDENING.md`**.

**No** runtime or API surface change is implied.

## Related

- `HYPERDENSITY_CONTRACTKIT_RELEASE_TAGGING.md`
- `HYPERDENSITY_CONTRACTKIT_VERSION_MODEL.md`
- `HYPERDENSITY_CONTRACTKIT_CONSUMER_CI_HARDENING.md`
- `pkg/hyperdensity/contractkit/contracts/release.go`
