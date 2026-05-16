# Karl Host Runtime — Linux MVP skeleton (KHR-F / G)

| Field | Value |
|-------|-------|
| **Binary** | `karl-host-runtime` |
| **Status** | Preview / sandbox only — **production unsupported** |
| **ISO** | Packaged disabled — see Karl-OS-ISO `KHR_HOST_RUNTIME_PREVIEW.md` |
| **Cluster proof** | **PASS** on `karl-metal-01@ovh` (KHR-G sandbox, KHR-I host status) |
| **Host contract** | `docs/khr/HOST_CONTRACT.md` |

---

## Host registration status (KHR-I)

```bash
go run ./cmd/karl-host-runtime -mode=host-status \
  -config=examples/khr/runtime-sandbox/karl-host-runtime-config-default.yaml \
  -node-name=karl-metal-01
```

| Item | Value |
|------|-------|
| CRD | `runtime.karl.io/Host` (cluster) |
| Evidence | [`docs/evidence/khr-host-registration/summary.json`](../evidence/khr-host-registration/summary.json) |
| Apply to cluster | **No** — JSON generation only; no controller |

---

## Sandbox validation (KHR-G)

Real cluster execution (not ISO provision):

| Item | Value |
|------|-------|
| Cluster | `karl-metal-01@ovh` |
| Namespace | `khr-runtime-sandbox` |
| Label | `khr.karl.io/sandbox=true` |
| Evidence (latest) | [`docs/evidence/khr-runtime-sandbox/summary.json`](../evidence/khr-runtime-sandbox/summary.json) |
| Run artifacts | `docs/evidence/khr-runtime-sandbox/20260516T130810Z/` |
| Procedure | `docs/khr/KHR_RUNTIME_SANDBOX_EXECUTION.md` |

**Scope remains sandbox only:** guarded apply writes a local marker file only; no production cgroup mutation; `sandboxApplyEnabled` defaults to **false**.

---

## Purpose

Host-side daemon skeleton for KHR Linux MVP:

1. Register host identity (local JSON, no kube apply)
2. Emit **Host** status JSON (`-mode=host-status`) — sandbox/read-only
3. Report capabilities (cgroup v2, runtime providers)
4. Emit **ResourcePort** candidate documents
5. **Dry-run** ResourceLease (reuses `pkg/khr/resourcelease`)
6. **Guarded apply** — sandbox marker file only when `sandboxApplyEnabled: true`
7. **Rollback baseline** for sandbox marker restore
8. **Flight recorder** — in-memory event trace

---

## Constraints

| Rule | Enforcement |
|------|-------------|
| Linux / cgroup only | `linuxOnly`, `sandboxMode` required |
| Namespace allowlist | `allowedNamespaces` (e.g. `khr-runtime-sandbox`) |
| Label allowlist | `allowedLabels` (e.g. `khr.karl.io/sandbox=true`) |
| No production mutation | `sandboxApplyEnabled` default **false** |
| Production unsupported | No GA host daemon on ISO; no autonomous apply |
| No VM/QMP/libvirt | blocked surfaces in capabilities report |
| No Windows | CellContext + platform checks |
| No autonomous apply | apply requires explicit config + CLI mode |

---

## Modes

```bash
go run ./cmd/karl-host-runtime -mode=register-host -config=examples/khr/karl-host-runtime-config.yaml
go run ./cmd/karl-host-runtime -mode=report-capabilities -config=...
go run ./cmd/karl-host-runtime -mode=emit-resourceport -config=...
go run ./cmd/karl-host-runtime -mode=dry-run-lease -config=... \
  -lease-input=examples/khr/runtime-sandbox/resourcelease-dry-run.json \
  -resource-port-input=examples/khr/runtime-sandbox/resourceport-dry-run.json
```

---

## Packages

| Path | Role |
|------|------|
| `cmd/karl-host-runtime` | CLI entry |
| `pkg/khr/host` | Registration, capabilities, sandbox gates |
| `pkg/khr/resourceport` | ResourcePort candidate emission |
| `pkg/khr/resourcelease` | Dry-run, guarded apply, rollback |
| `pkg/khr/cgroup` | Detection + envelope plan (existing) |
| `pkg/khr/flightrecorder` | In-memory trace |

---

## ResourcePort reporting loop (KHR-J)

| Item | Value |
|------|-------|
| Mode | `-mode=resourceport-loop` |
| Default | JSON only; `resourcePortLoopEnabled: false` |
| Docs | `docs/khr/RESOURCEPORT_REPORTING_LOOP.md` |
| Evidence | `docs/evidence/khr-resourceport-loop/summary.json` |

No automatic ResourceLease apply.

---

## Roadmap — next engineering steps (post KHR-I)

| Step | Outcome |
|------|---------|
| **H+1** | Host CR + local status generation (**done** KHR-I); optional cluster apply behind flag |
| **H+2** | **ResourcePort controller loop** (observe → reconcile candidates; sandbox namespace only) |
| **H+3** | ResourceLease apply gate integration with controller (still no production mutation) |

ISO boundary: KubeVirt/CDI remain **compatibility providers**; OVN-native fabric remains **target** network — see Karl-OS-ISO `docs/adr/ADR-0001-iso-khr-boundaries.md`.

---

## Related

- `docs/khr/KHR_LINUX_SANDBOX_SAFETY.md`
- `docs/khr/KHR_RUNTIME_SANDBOX_EXECUTION.md`
- `docs/khr/KHR_LINUX_AGENT_RUNBOOK.md` (separate binary, still non-production)
- `docs/roadmap/KHR_HYPERDENSITY_CORRECTED_ROADMAP.md`
