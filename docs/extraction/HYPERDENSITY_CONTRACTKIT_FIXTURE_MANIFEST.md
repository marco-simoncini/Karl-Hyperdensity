# Hyperdensity contractkit fixture manifest (Sprint 26)

## Purpose

Shared **JSON manifest** listing M1–M8 Hyperdensity parity cases: milestone, Dashboard-shaped fixture path, contract golden path, and claim-safe metadata. **Test-only** — no runtime or handler import.

## Locations

| Consumer | Path |
|----------|------|
| contractkit (canonical schema + example) | `pkg/hyperdensity/contractkit/testdata/dashboard/hyperdensity_parity_manifest_m1_m7.json` |
| Karl-Dashboard (consumer paths) | `kubernetes-console/pkg/server/testdata/hyperdensity_parity_manifest_m1_m7.json` |

Case **IDs** and **metadata** align; **file paths** differ per repo layout.

## Versioning (three layers — see `HYPERDENSITY_CONTRACTKIT_VERSION_MODEL.md`)

| Field | Matches |
|-------|---------|
| `manifestVersion` | `FixtureManifestVersion` (`hyperdensity.parity.manifest/v1`) |
| `contractKitVersion` | `ContractKitVersion` / schema (`v0.0.0-sprint26`) |
| `go.mod` module pin | `ContractKitModuleVersion` (`v0.1.0-khr-m1-m9`) — **not** in JSON manifest |

## API (contractkit/contracts)

- `ParseFixtureManifest([]byte)`
- `ValidateFixtureManifest(FixtureManifest)`
- `Version()` / `ContractKitVersion`
- `ExpectedM1M8CaseIDs()` — canonical M1–M8 case ID list (Sprint 27)
- `CaseIDs(m)` / `ManifestCaseIDSet(m)` — drift guards

Validation rules: unique case IDs; non-empty paths; `claimSafe: true`; `windowsEnabled: false`; `kubeVirtLegacyRequired: true`; version match.

Contractkit tests assert the **example manifest** contains **exactly** the `ExpectedM1M8CaseIDs()` set.

## Dashboard usage

`TestHyperdensityParityManifest` reads the consumer manifest, validates via contractkit, checks all referenced files exist under `pkg/server/testdata/`, and asserts the case ID set **matches exactly** `ExpectedM1M8CaseIDs()` (no silent drift).

## Release tagging

Optional git tag on parent repo: `contractkit/v0.1.0-khr-m1-m9` — see `HYPERDENSITY_CONTRACTKIT_RELEASE_TAGGING.md`. Pseudo-versions remain valid until tags are applied.

## Boundaries

- No HTTP, cluster, or npm.
- Dashboard **must not** import manifest helpers in non-test code.
- Hyperdensity root module re-exports types via `pkg/hyperdensity/contracts/alias.go` for in-repo use only.

## Related

- `HYPERDENSITY_CONTRACTKIT_MODULE.md`
- Dashboard `HYPERDENSITY_PARITY_MANIFEST_M9.md`
