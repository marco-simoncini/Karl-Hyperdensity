# Grande Padre — evidence indexes (Sprint 12)

## EvidenceIndex

Each ingested request produces one **EvidenceIndex** row (after optional `DeduplicateBySha256` compaction):

| Field | Meaning |
|-------|---------|
| `artifactId` | Operator or manifest artifact id. |
| `bundleSha256` | SHA-256 of **canonical** collect-evidence JSON (see `pkg/khr/evidence/integrity`). |
| `cellRef` | Cell (or empty) inferred from bundle `cellRef` / telemetry `cellRef`. |
| `confidence` | Aggregated evidence confidence (`low` \| `medium` \| `high`). |
| `readyForGrandePadre` | From bundle `evidenceSummary` (KHR summarizer). |
| `blockedReasons` | From bundle `evidenceSummary` (may be empty). |
| `warnings` | From bundle `evidenceSummary`. |
| `trustTier` | Integrity-oriented tier (see below). |
| `indexedAt` | RFC3339 timestamp from injectable `NowFunc` (tests pin time). |

## TrustTier

| Tier | When |
|------|------|
| `IntegrityFailed` | Digest / manifest SHA alignment failed. |
| `DevOnly` | `signingMode=local-dev` or annotation `khr.karl.io/signature-trust-tier: DevOnly`. |
| `Unknown` | Unsupported signed mode in this skeleton (placeholder). |
| `IntegrityVerified` | Default for digest-only match (`signingMode` none / absent) — means **canonical bytes matched declared digest**, not production PKI. |
| `Unsigned` | Optional label when `-unsigned-digest-trust=unsigned` is used: same integrity as verified, but explicitly **no signature trust path** (see `trust.go` + runbook). |

**Important:** trust tiers are **not** apply authorization. They inform indexing, dry-run, and recommendations only.

## Queries

- **Ready:** `readyForGrandePadre`, no `blockedReasons`, trust ≠ `IntegrityFailed`.
- **Blocked:** not ready, or blocked reasons present, or `IntegrityFailed`.
- **By confidence / by cell:** filter on `EvidenceIndex` fields.

## Deduplication

`DeduplicateBySha256` keeps the **latest** `indexedAt` row per `bundleSha256`. `duplicateCount` in CLI output reflects duplicate ingest observations on the store before compaction logic runs (see `store.go`).
