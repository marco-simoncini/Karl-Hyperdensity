# Hyperdensity contractkit fixture manifest (Sprint 26)

## Purpose

Shared **JSON manifest** listing M1–M8 Hyperdensity parity cases: milestone, Dashboard-shaped fixture path, contract golden path, and claim-safe metadata. **Test-only** — no runtime or handler import.

## Locations

| Consumer | Path |
|----------|------|
| contractkit (canonical schema + example) | `pkg/hyperdensity/contractkit/testdata/dashboard/hyperdensity_parity_manifest_m1_m7.json` |
| Karl-Dashboard (consumer paths) | `kubernetes-console/pkg/server/testdata/hyperdensity_parity_manifest_m1_m7.json` |

Case **IDs** and **metadata** align; **file paths** differ per repo layout.

## Versioning

- Manifest field `contractKitVersion` must match `contracts.ContractKitVersion` (e.g. `v0.0.0-sprint26`).
- `contracts.ContractKitCommitHint` is `"consumer-pinned"` — authoritative revision is the **go.mod pseudo-version** on Karl-Dashboard, not an embedded commit SHA.

## API (contractkit/contracts)

- `ParseFixtureManifest([]byte)`
- `ValidateFixtureManifest(FixtureManifest)`
- `Version()` / `ContractKitVersion`

Validation rules: unique case IDs; non-empty paths; `claimSafe: true`; `windowsEnabled: false`; `kubeVirtLegacyRequired: true`; version match.

## Dashboard usage

`TestHyperdensityParityManifest` reads the consumer manifest, validates via contractkit, checks all referenced files exist under `pkg/server/testdata/`, and asserts expected case IDs.

## Boundaries

- No HTTP, cluster, or npm.
- Dashboard **must not** import manifest helpers in non-test code.
- Hyperdensity root module re-exports types via `pkg/hyperdensity/contracts/alias.go` for in-repo use only.

## Related

- `HYPERDENSITY_CONTRACTKIT_MODULE.md`
- Dashboard `HYPERDENSITY_PARITY_MANIFEST_M9.md`
