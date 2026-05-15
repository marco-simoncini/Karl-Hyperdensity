# EvidenceIngestRequest CRD (`hyperdensity.karl.io/v1alpha1`)

## Kind

- **API group:** `hyperdensity.karl.io`
- **Version:** `v1alpha1`
- **Kind:** `EvidenceIngestRequest`
- **Resource:** `evidenceingestrequests` (short name: `kringest`)

## Spec

| Field | Description |
|-------|-------------|
| `artifactId` | Optional id echoed from manifest. |
| `bundle` | Evidence bundle JSON (`collect-evidence` shape). |
| `manifest` | Artifact manifest JSON (Sprint 10). |
| `digest` | Lowercase hex SHA-256 line (trimmed). |
| `source` | `agentId` (required), optional `nodeName`, `hostId`, `tenant`. |
| `policy` | `requireDigestMatch`, `allowUnsigned`, `allowLocalDevSignature`, `maxBundleBytes`. |
| `dryRunOnly` | When `true`, consumer must not promote to execution. |
| `ttlSeconds` | Retention / processing window hint. |

## Status (expected controller fields)

| Field | Description |
|-------|-------------|
| `phase` | `Pending` \| `Validating` \| `Accepted` \| `Rejected` \| `Stored` \| `Indexed` \| `Failed` |
| `digestMatch` | Whether digest matched recomputed canonical bundle. |
| `signatureStatus` | e.g. `None`, `DevOnly`, `Valid`, `Invalid`, … |
| `rejectionReasons` | Human-readable reject list. |
| `evidenceBundleRef` | Name/namespace of resulting `EvidenceBundle`, if any. |

The CRD OpenAPI keeps `status` as `x-kubernetes-preserve-unknown-fields` in Sprint 11 to allow controller evolution without frequent CRD churn.

## Local preparation

`prepare-ingest-request` may populate `metadata.annotations`:

- `khr.karl.io/preparation-warnings` — JSON array of mismatch / sanity warnings.
- `khr.karl.io/signature-trust-tier: DevOnly` — when manifest `signingMode` is `local-dev`.
