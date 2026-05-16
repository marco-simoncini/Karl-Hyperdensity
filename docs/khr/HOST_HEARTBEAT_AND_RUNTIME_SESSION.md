# Host heartbeat and runtime session (KHR-N)

Stabilizes the runtime contract with **periodic Host status JSON**, **session identity**, and **flight-recorder correlation**. **No new mutation paths.**

## Mode: `host-heartbeat`

```bash
karl-host-runtime -mode=host-heartbeat \
  -config=examples/khr/runtime-sandbox/karl-host-runtime-config-loop.yaml \
  -namespace=khr-runtime-sandbox \
  -cluster-context=karl-metal-01@ovh \
  -heartbeat-iterations=3 -heartbeat-interval-ms=500 \
  -heartbeat-output=/tmp/khr-host-status.json
```

## Host status fields (each tick)

| Field | Description |
|-------|-------------|
| `lastHeartbeatTime` | RFC3339 UTC |
| `runtimeVersion` | `host.RuntimeVersion` |
| `safetyMode` | `sandbox` \| `production-blocked` |
| `activeResourcePorts` | Read-only list from cluster (`kubectl get resourceports`) |
| `activeResourceLeases` | Read-only hints from sandbox apply evidence |
| `lastApplyState` | `idle` \| `applied` (from sandbox baseline/evidence) |
| `runtimeSessionId` | Stable per process |
| `hostRuntimeInstanceId` | Stable per process |
| `correlationId` | Per-operation id |

## Runtime session

- `runtimeSessionId`: `khr-session-<hex>`
- `hostRuntimeInstanceId`: `khr-inst-<hex>`
- `correlationId`: advances per `CurrentRuntimeSession()` or explicit `SetCorrelationID`

Flight recorder events include all three fields for apply / rollback / heartbeat.

## Stale heartbeat

Pass `--prior-heartbeat-at=<RFC3339>` to simulate stale detection (default threshold 2m).

## Safety

- `noMutation: true` on all heartbeat output
- No Host CR `kubectl apply`
- No production namespace writes

## Evidence

`./scripts/khr_host_heartbeat_evidence.sh`

## Related

- Karl-Inventory `docs/khr/INVENTORY_RUNTIME_POSTURE.md` (read-only posture projection)
- Dashboard KHR projection `hosts[].heartbeat`, `runtimeSession`, `postureSummary`
