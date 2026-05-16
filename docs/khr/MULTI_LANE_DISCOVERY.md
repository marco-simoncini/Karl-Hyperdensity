# Multi-lane runtime discovery (KHR-Q)

KHR-Q performs **read-only** discovery of Linux and Windows workloads on the cluster and classifies runtime lanes for live scale posture. **No apply, patch, restart, rollout, or recreate.**

## Mode

```bash
karl-host-runtime -mode=lane-discovery \
  -config=examples/khr/runtime-sandbox/karl-host-runtime-config-lane-discovery.yaml \
  -cluster-context=karl-metal-01@ovh
```

Requires `spec.laneDiscoveryEnabled: true` and `spec.sandboxMode: true`.

## Lanes

| Lane | Description |
|------|-------------|
| `linux-container-cgroup` | Sandbox/container pods — cgroup live-in-place candidate |
| `linux-vm-compatibility` | Linux KubeVirt VMs — compatibility-fallback |
| `windows-vm-session` | Windows KubeVirt VMs / session — observation + blocked restart paths |
| `kubevirt-compatibility` | Unknown guest on KubeVirt |

## Classification

| Value | Meaning |
|-------|---------|
| `live-in-place-capable` | Linux sandbox cgroup lane (KHR-O) |
| `observation-only` | Discovered but no live apply asserted |
| `compatibility-fallback` | KubeVirt / Windows compatibility path |
| `unsupported` | Stopped or unsupported posture |

## Output JSON

| Field | Content |
|-------|---------|
| `discoveredHosts[]` | Cluster nodes |
| `discoveredShells[]` | Projected shells |
| `discoveredCells[]` | Projected cells with `vmType`, `osFamily`, `sessionType` |
| `discoveredResourcePorts[]` | Lane + provider + classification |
| `laneCapabilities[]` | Per-lane summary counts |
| `blockedStates[]` | `requiresRestart`, `providerUnsupported`, etc. |
| `safety` | `readOnly`, `noPatch`, `noApply`, `noRestart`, `noRollout`, `noRecreate` |

## Evidence

```bash
./scripts/khr_lane_discovery_evidence.sh
```

Artifacts under `docs/evidence/khr-lane-discovery/`.

## Cluster sources (read-only)

- `kubectl get nodes`
- `kubectl get virtualmachines.kubevirt.io -A`
- `kubectl get pods -A` (sandbox + virt-launcher correlation)
- `kubectl get resourceports -A` (when present)

No `kubectl apply`, `patch`, or workload restart.
