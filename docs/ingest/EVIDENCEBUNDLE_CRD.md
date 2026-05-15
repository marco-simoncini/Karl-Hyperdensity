# EvidenceBundle CRD (`hyperdensity.karl.io/v1alpha1`)

## Kind

- **API group:** `hyperdensity.karl.io`
- **Version:** `v1alpha1`
- **Kind:** `EvidenceBundle`
- **Resource:** `evidencebundles` (short name: `krbundle`)

## Spec (immutable snapshot)

| Field | Description |
|-------|-------------|
| `artifactId` | Optional operator id (may be empty). |
| `agentId` | Agent that produced the bundle. |
| `sourceMode` | Pipeline name (e.g. `collect-evidence`). |
| `bundleSha256` | Lowercase hex SHA-256 over **canonical** bundle JSON. |
| `bundleBytes` | Byte length of canonical bundle JSON. |
| `signingMode` | Integrity signing mode from manifest (`none`, `local-dev`, …). |
| `signaturePresent` | Whether a signature was present on the manifest. |
| `collectedAt` | RFC3339 from bundle. |
| `cellRef` | Optional Cell correlation. |
| `confidence` | Aggregated confidence string. |
| `readyForGrandePadre` | Readiness flag from bundle summary. |
| `blockedReasons` | List of blockers. |
| `warnings` | List of warnings. |
| `manifestRef` | Optional ref to stored manifest. |
| `storageRef` | Optional ref to stored bundle bytes. |

## Status (controller-owned)

Typical fields (preserved as unknown-friendly in the CRD stub; refine in later revisions):

| Field | Purpose |
|-------|---------|
| `integrityStatus` | `Pending` \| `Verified` \| `Failed` \| `Unsigned` \| `DevOnly` |
| `phase` | `Pending` \| `Accepted` \| `Rejected` \| `Indexed` \| `Expired` |
| `receivedAt` | When the control plane first recorded the bundle. |

## Relationship to `EvidenceIngestRequest`

An accepted ingest may materialize or reference an `EvidenceBundle` via `EvidenceIngestRequest.status.evidenceBundleRef`.
