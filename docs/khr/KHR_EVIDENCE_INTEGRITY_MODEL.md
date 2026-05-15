# KHR evidence integrity model (Sprint 10)

## Goals

Sprint 10 adds **local integrity metadata** around the `collect-evidence` bundle:

- **Canonical JSON** for hashing (compact `encoding/json` output, `SetEscapeHTML(false)`, no trailing newline).
- **SHA-256** digest over those canonical bytes (`bundleSha256` + `bundleBytes` in the manifest).
- **Artifact manifest** JSON as a sidecar describing the bundle, signing mode, and chain-of-custody hints.
- **Optional `local-dev` signing** (Ed25519 over canonical bytes) — **not** a production trust anchor.

There is **no transport**, **no ingest API**, and **no apply** path: `mutationsForbidden` remains `true` on both the bundle and the manifest.

## CLI flags (`collect-evidence` only)

| Flag | Purpose |
|------|---------|
| `-evidence-manifest-output` | Write artifact manifest JSON (mode `0o600`). |
| `-evidence-digest-output` | Write a single-line lowercase hex SHA-256 + newline. |
| `-signing-mode` | `none` (default) or `local-dev`. |
| `-signing-key-file` | PKCS#8 PEM **Ed25519** private key; **required** when `local-dev`. |
| `-artifact-id` | Optional string recorded as `artifactId` in the manifest. |

`local-dev` additionally **requires** `-evidence-manifest-output` so signature metadata has a defined home.

## Determinism in tests

| Variable | Effect |
|----------|--------|
| `KHR_TEST_INTEGRITY_NOW` | RFC3339 timestamp for `generatedAt` in the manifest. |
| `KHR_TEST_INTEGRITY_CHAIN_STUB` | Replace `chainOfCustody` with fixed placeholder values (golden-friendly). |

## Non-goals

- No PKI, no timestamping authority, no TUF, no transparency log.
- No verification service: consumers may recompute SHA-256 over canonical bytes locally.

See also: `docs/khr/KHR_EVIDENCE_ARTIFACT_MANIFEST.md`, `pkg/khr/evidence/integrity/`, `examples/khr/evidence-integrity/`.
