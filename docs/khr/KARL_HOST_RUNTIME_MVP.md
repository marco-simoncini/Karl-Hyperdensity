# Karl Host Runtime — Linux MVP skeleton (KHR-F)

| Field | Value |
|-------|-------|
| **Binary** | `karl-host-runtime` |
| **Status** | Preview / sandbox only — **not production** |
| **ISO** | Packaged disabled — see Karl-OS-ISO `KHR_HOST_RUNTIME_PREVIEW.md` |

---

## Purpose

Host-side daemon skeleton for KHR Linux MVP:

1. Register host identity (local JSON, no kube apply)
2. Report capabilities (cgroup v2, runtime providers)
3. Emit **ResourcePort** candidate documents
4. **Dry-run** ResourceLease (reuses `pkg/khr/resourcelease`)
5. **Guarded apply** — sandbox marker file only when `sandboxApplyEnabled: true`
6. **Rollback baseline** for sandbox marker restore
7. **Flight recorder** — in-memory event trace

---

## Constraints

| Rule | Enforcement |
|------|-------------|
| Linux / cgroup only | `linuxOnly`, `sandboxMode` required |
| Namespace allowlist | `allowedNamespaces` |
| Label allowlist | `allowedLabels` |
| No production mutation | `sandboxApplyEnabled` default **false** |
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
  -lease-input=examples/khr/resourcelease-linux-envelope-dry-run.json \
  -resource-port-input=examples/khr/resourceport-linux-envelope-for-dryrun.json
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

## Related

- `docs/khr/KHR_LINUX_SANDBOX_SAFETY.md`
- `docs/khr/KHR_LINUX_AGENT_RUNBOOK.md` (separate binary, still non-production)
