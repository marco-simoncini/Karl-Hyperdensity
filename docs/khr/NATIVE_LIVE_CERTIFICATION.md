# Native live lane certification (KHR-T)

| Field | Value |
|-------|-------|
| **Lane** | `native-live` (unchanged — no new lanes) |
| **Certification ID** | `khr-native-live-certification-v1` |
| **Automation** | **None** — manual evidence scripts only |
| **Orchestration** | **None** — read-only certification semantics |

---

## Purpose

Consolidate the KHR-S native-live prototype into a **verifiable, repeatable baseline**:

- Multiple evidence runs with deterministic fingerprints
- Regression guard (fail on restart / rollout / recreate / interruption)
- Baseline compare against `examples/khr/native-live/baseline-certification.json`
- Certification summary JSON for Inventory/Dashboard projection

---

## Pipeline

| Step | Script / artifact |
|------|-------------------|
| Single run | `scripts/khr_native_live_lane_run.sh <run-dir>` → `run-metrics.json` |
| Multi-run certify | `scripts/khr_native_live_certify.sh` → `certification-summary.json` |
| Go aggregate + guard | `cmd/khr-native-live-certify` |
| Unit tests | `pkg/khr/nativelive` |

```bash
export KHR_RUNTIME_CLUSTER_CONTEXT=karl-metal-01@ovh
./scripts/khr_native_live_certify.sh
```

Environment:

| Variable | Default | Meaning |
|----------|---------|---------|
| `KHR_NATIVE_LIVE_CERT_RUNS` | `2` | Repeatable run count |
| `KHR_NATIVE_LIVE_CERT_ID` | UTC timestamp | Certification bundle id |

Output: `docs/evidence/khr-native-live-lane/certification/<id>/`

---

## Metrics (per run)

| Metric | Description |
|--------|-------------|
| `applyLatencyMs` | CPU / RAM up / RAM down guarded apply |
| `rollbackLatencyMs` | Per-step rollback |
| `interruptionWindowMs` | Window when restart/rollout/recreate detected (0 = pass) |
| `restartCountDelta` | Pod restart count change |
| `rolloutCount` | Deployment generation drift |
| `recreateDetected` | Pod UID change |

---

## Regression guard

Certification **fails** when any run reports:

- `restartCountDelta > 0`
- `rolloutDetected` or `rolloutCount > 0`
- `recreateDetected`
- `interruptionDetected` or `interruptionWindowMs > 0`

Go: `nativelive.CheckRegression(summary)`

---

## Scores (read-only)

| Field | Range / values |
|-------|----------------|
| `continuityScore` | `0.0` – `1.0` |
| `liveScaleConfidence` | `high` \| `medium` \| `low` |

---

## Related

- `docs/khr/NATIVE_LIVE_LANE_PROTOTYPE.md` (KHR-S)
- `examples/khr/native-live/baseline-certification.json`
