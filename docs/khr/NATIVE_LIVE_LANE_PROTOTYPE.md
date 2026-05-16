# Native live lane prototype (KHR-S)

| Field | Value |
|-------|-------|
| **Lane** | `native-live` |
| **Classification** | `native-live` |
| **Provider** | `khr.native` (Linux container cgroup — not KubeVirt) |
| **Status** | Prototype / sandbox only |

---

## Purpose

Demonstrate a **real** Linux/container live-capable lane end-to-end on `karl-metal-01@ovh`:

1. Lane discovery (`lane-discovery`)
2. ResourceFuture simulation (`resourcefuture-simulate`) with `liveInPlaceEligible=true`
3. ResourceLease dry-run
4. Guarded apply (CPU + RAM)
5. Verify (no restart / rollout / recreate)
6. Rollback

Compare against `kubevirt-compatibility` lanes where live-in-place is not asserted.

---

## Workload identification

| Signal | Meaning |
|--------|---------|
| Label `khr.karl.io/native-live=true` | Native-live lane |
| Name prefix `khr-native-live-` in `khr-runtime-sandbox` | Native-live lane |
| Sandbox label `khr.karl.io/sandbox=true` | Required |
| Not `virt-launcher` / not VM | Excluded from VM compatibility paths |

Example manifest: `examples/khr/runtime-sandbox/native-live-workload.yaml`

---

## Config

| File | Modes |
|------|-------|
| `karl-host-runtime-config-native-live.yaml` | `lane-discovery`, `resourcefuture-simulate` |
| `karl-host-runtime-config-guarded-apply.yaml` | `resourcelease-dryrun`, `resourcelease-guarded-apply`, `resourcelease-rollback` |

---

## Evidence

```bash
export KHR_RUNTIME_CLUSTER_CONTEXT=karl-metal-01@ovh
./scripts/khr_native_live_lane_evidence.sh
```

Output: `docs/evidence/khr-native-live-lane/<runId>/`

| Artifact | Content |
|----------|---------|
| `lane-discovery.json` | Includes `native-live` vs `kubevirt-compatibility` |
| `simulation.json` | `liveInPlaceEligibility` with `lane=native-live` |
| `apply-*.json` | CPU/RAM guarded apply + verification |
| `summary.json` | Metrics: no restart, no rollout, latency, rollback PASS |

---

## Metrics (KHR-S)

| Metric | Expected |
|--------|----------|
| Pod restart count | Unchanged across apply |
| Deployment generation | Unchanged (no rollout) |
| Session continuity | No interruption |
| `verification.noRestart` | `true` |
| `verification.noRollout` | `true` |
| `verification.noRecreate` | `true` |

---

## Related

- `docs/khr/MULTI_LANE_DISCOVERY.md` (KHR-Q)
- `docs/khr/RESOURCEFUTURE_SIMULATION.md` (KHR-R)
- `docs/khr/RAM_LIVE_SCALE_SANDBOX.md` (KHR-O)
