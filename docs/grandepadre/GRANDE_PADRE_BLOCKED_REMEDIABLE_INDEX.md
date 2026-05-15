# Grande Padre — blocked / remediable index (Sprint 12)

## BlockedRemediableIndex

`BuildBlockedRemediableIndex` groups **blocked-like** evidence rows by `cellRef` (namespace + name). Rows without a cell are grouped under an internal placeholder key for the skeleton.

Each entry contains:

- `cellRef` — API reference to the Cell when present.
- `blockedReasons` — de-duplicated reasons seen across rows for that cell.
- `remediationHints` — **heuristic** strings derived from reasons (e.g. cgroup-related copy). Not machine-verified remediation.
- `lastArtifactId` — artifact id from the **newest** row (by `indexedAt`) in the group.
- `confidence` / `trustTier` — copied from that newest row for quick triage.

## Usage

Hyperdensity planners can surface this structure for **dry-run** and **recommendation** flows: it highlights where evidence is not ready and suggests next checks, without implying that ingest or indexing authorizes apply.

## Limitations

- Hints are static pattern matching, not policy-engine output.
- No correlation with live cluster state; evidence may be stale relative to the cell.
