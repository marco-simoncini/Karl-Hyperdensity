# Host contract (runtime.karl.io/Host) — KHR-I

| Field | Value |
|-------|-------|
| **API group** | `runtime.karl.io` |
| **Kind** | `Host` |
| **Scope** | Cluster |
| **Controller** | **None** in KHR-I (status from `karl-host-runtime` only) |

---

## Spec (minimal)

| Field | Type | Description |
|-------|------|-------------|
| `hostId` | string | Stable host identity |
| `nodeName` | string | Kubernetes node name |
| `provider` | string | e.g. `khr.native` |
| `runtimeMode` | enum | `sandbox` \| `preview` \| `disabled` |
| `labels` | map | Host labels |
| `taints` | array | Optional taints |

---

## Status (minimal)

| Field | Type | Description |
|-------|------|-------------|
| `phase` | string | e.g. `Observed` |
| `conditions` | array | Ready / SandboxOnly |
| `capabilities` | object | Cgroup + provider capabilities JSON |
| `observedResourcePorts` | array | `{name, namespace}` refs |
| `lastHeartbeatTime` | RFC3339 | Last local heartbeat |
| `runtimeVersion` | string | `karl-host-runtime` version |
| `safetyMode` | string | `sandbox` or `production-blocked` |

---

## karl-host-runtime

```bash
go run ./cmd/karl-host-runtime \
  -mode=host-status \
  -config=examples/khr/runtime-sandbox/karl-host-runtime-config-default.yaml \
  -node-name=karl-metal-01 \
  -namespace=khr-runtime-sandbox \
  -port-name=khr-runtime-sandbox-port
```

**Sandbox only.** No kube apply, no controller loop, no production mutation.

---

## Evidence

Cluster-generated status (read-only): `docs/evidence/khr-host-registration/`

---

## Related

- `docs/khr/KARL_HOST_RUNTIME_MVP.md`
- `docs/contracts/khr/examples/host-karl-metal-01.json`
- Karl-Dashboard `khrProjection.hosts[]` (read-only)
