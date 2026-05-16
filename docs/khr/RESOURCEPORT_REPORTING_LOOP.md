# ResourcePort reporting loop (KHR-J)

| Field | Value |
|-------|-------|
| **Mode** | `resourceport-loop` |
| **Default** | JSON only — **no CR apply**, no mutation |
| **ResourceLease** | **Not** applied automatically |

---

## Behavior

1. Load sandbox `KarlHostRuntimeConfig` (`resourcePortLoopEnabled: true` required).
2. Validate namespace / label allowlists and production namespace blocklist.
3. Optionally verify cluster context (`karl-metal-01@ovh`).
4. Discover sandbox pods (kubectl) or use explicit targets (tests).
5. Emit periodic **ResourcePort** status JSON per iteration.
6. With `--emit-cr=true`: write local CR preview files only (`cr-preview`) — **never** `kubectl apply` in KHR-J.

---

## CLI

```bash
go run ./cmd/karl-host-runtime \
  -mode=resourceport-loop \
  -config=examples/khr/runtime-sandbox/karl-host-runtime-config-loop.yaml \
  -namespace=khr-runtime-sandbox \
  -cluster-context=karl-metal-01@ovh \
  -loop-iterations=2 \
  -loop-interval-ms=500
```

| Flag | Default | Purpose |
|------|---------|---------|
| `--emit-cr` | `false` | Local CR preview files under `--loop-output-dir` |
| `--loop-iterations` | `1` | Observation cycles |
| `--loop-interval-ms` | `0` | Delay between cycles |
| `--cluster-context` | current | Discovery guard |

---

## Safety

| Guard | Enforcement |
|-------|-------------|
| `resourcePortLoopEnabled` | `false` by default — loop blocked |
| Namespace allowlist | `allowedNamespaces` |
| Label allowlist | `allowedLabels` |
| Production namespaces | Blocked (`karl-system`, `kube-system`, …) |
| Cluster context | Must match `karl-metal-01@ovh` when discovery runs |
| Flight recorder | Events on start / iteration / complete / block |

---

## Evidence

`docs/evidence/khr-resourceport-loop/` — cluster JSON-only run on `karl-metal-01@ovh`.

---

## Related

- `docs/khr/KARL_HOST_RUNTIME_MVP.md`
- `docs/khr/HOST_CONTRACT.md`
