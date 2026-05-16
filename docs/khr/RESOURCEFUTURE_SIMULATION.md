# ResourceFuture simulation planner (KHR-R)

KHR-R introduces **read-only ResourceFuture planning** via `resourcefuture-simulate`. The planner consumes lane discovery, host status, ResourcePort/Lease observations, and posture summaries to emit **simulated** forecasts only.

## Scope

| In scope | Out of scope |
|----------|----------------|
| Candidate scale plans (simulation) | Automatic apply |
| Saturation / compatibility / restart forecasts | Autonomous orchestration |
| CPU / RAM / mixed / Windows compatibility simulations | Production mutation |
| Safety attestation (`noApply`, `noMutation`, …) | ResourceFuture CR apply |

## Mode

```bash
karl-host-runtime -mode=resourcefuture-simulate \
  -config=examples/khr/runtime-sandbox/karl-host-runtime-config-resourcefuture-simulate.yaml \
  -cluster-context=karl-metal-01@ovh
```

Requires `laneDiscoveryEnabled` and `resourceFutureSimulationEnabled`.

## Inputs (aggregated)

| Source | Use |
|--------|-----|
| Host status | `host.BuildHostStatus` from config + observed ports |
| Lane discovery | Full `lane-discovery` result (cells, ports, lanes) |
| ResourcePorts | From discovery + cluster |
| ResourceLeases | Synthesized dry-run lease refs per cell |
| Posture summary | Sandbox/preview/blocker snapshot |

## Outputs

| Field | Description |
|-------|-------------|
| `candidateScalePlans[]` | Simulated CPU/RAM scale proposals |
| `saturationForecast[]` | Pressure risk by target/resource |
| `blockedConstraints[]` | Predicted blockers |
| `compatibilityFallbackPrediction[]` | KubeVirt/Windows fallback likelihood |
| `liveInPlaceEligibility[]` | Live-in-place eligibility per target |
| `restartRequiredPrediction[]` | Restart risk per target |
| `forecasts` | `cpuScale`, `ramScale`, `mixedLane`, `windowsCompatibility` |

## Safety

`safety` object always reports: `readOnly`, `noMutation`, `noApply`, `noRestart`, `noRollout`, `noRecreate`, `noAutonomousOrchestration`, `simulationOnly`.

## Evidence

```bash
./scripts/khr_resourcefuture_evidence.sh
```

Artifacts: `docs/evidence/khr-resourcefuture/`
