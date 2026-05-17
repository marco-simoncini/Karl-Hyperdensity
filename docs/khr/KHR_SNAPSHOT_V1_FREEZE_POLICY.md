# KHR Snapshot v1 freeze policy (KHR-BU)

Policy for **KHR TP Reference Snapshot v1** (`committed-khr-bt-v1`) after sprint KHR-BT.

---

## Frozen artifacts

| Item | Path / id |
|------|-----------|
| Snapshot runId | `committed-khr-bt-v1` |
| Master JSON | `docs/evidence/khr-tp-reference-snapshot-v1/committed-khr-bt-v1/snapshot-summary.json` |
| Cross-repo index | `.../cross-repo-evidence-index.json` |
| Contract | `docs/contracts/khr/khr-tp-reference-snapshot-v1.json` |

---

## Freeze rules

1. **No in-place edits** to committed evidence directories referenced by the snapshot index without bumping snapshot version (v2).
2. **Offline validate** must pass using frozen paths — not current cluster state.
3. **New live runs** produce new `runId` directories; promotion to snapshot requires explicit sprint + aggregator re-run.
4. **globalDefaultsChanged** must remain `false` in snapshot and linked Dashboard evidence.
5. **productionReady** and **autonomousOrchestration** remain `false` for all snapshot-linked artifacts.

---

## Allowed changes post-freeze

- Documentation clarifying snapshot (no evidence field changes)
- Offline validation / CI wiring (KHR-BU)
- New evidence runs stored alongside (not replacing committed ids)

---

## Forbidden without new snapshot version

- Overwriting `committed-khr-bt-v1/` evidence files
- Changing `contractSetId` or scope readiness labels in snapshot without v2
- Promoting live-only FAIL verify summaries over committed PASS execute summaries

---

## Validation

```bash
./scripts/khr_validate_reference_snapshot.sh
```
