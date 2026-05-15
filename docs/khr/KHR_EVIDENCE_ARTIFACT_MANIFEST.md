# KHR evidence artifact manifest (Sprint 10)

## Purpose

The artifact manifest is a **JSON sidecar** emitted next to (or near) the human-readable evidence bundle. It summarizes what was hashed and how it was signed (if at all).

## Required fields

| Field | Type | Notes |
|-------|------|--------|
| `artifactId` | string | Optional operator-supplied id; may be empty. |
| `agentId` | string | From agent config `spec.agentId`. |
| `generatedAt` | string | RFC3339 UTC; overridable in tests via `KHR_TEST_INTEGRITY_NOW`. |
| `bundleSha256` | string | Lowercase hex SHA-256 of **canonical** bundle JSON. |
| `bundleBytes` | int | Length in bytes of canonical bundle JSON. |
| `signingMode` | string | `none` or `local-dev`. |
| `signaturePresent` | bool | `true` iff an Ed25519 signature was written under `signature`. |
| `sourceMode` | string | Always `collect-evidence` for this manifest. |
| `mutationsForbidden` | bool | Always `true`; never implies apply authority. |
| `chainOfCustody` | object | Best-effort local process metadata (hostname, user, pid, paths). Stubbed in tests via `KHR_TEST_INTEGRITY_CHAIN_STUB=1`. |

## Optional fields

| Field | When |
|-------|------|
| `signatureAlgorithm` | `local-dev` signing succeeded (`ED25519-local-dev`). |
| `signature` | Base64-encoded Ed25519 signature over canonical bundle bytes. |
| `integrityNotes` | Static reminders that this is local-only metadata. |

## Relationship to the bundle

- **Stdout / `-evidence-output`** remain **pretty-printed** JSON for humans.
- **Digest** is always computed from **canonical** JSON bytes (may differ from the pretty-printed file on disk).

## Example fixtures

See `examples/khr/evidence-integrity/manifest-none.json`, `digest.txt`, `manifest-local-dev.json`, and `local-dev-key.example`.
