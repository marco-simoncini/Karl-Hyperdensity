# KHR Linux sandbox safety (karl-host-runtime)

| Field | Value |
|-------|-------|
| **Applies to** | `karl-host-runtime` KHR-F |
| **Production** | **Never** in this sprint |

---

## Defaults

| Setting | Default | Effect |
|---------|---------|--------|
| `sandboxApplyEnabled` | `false` | All `apply-lease` attempts **blocked** |
| `sandboxMode` | must be `true` | Refuses non-sandbox configs |
| `linuxOnly` | must be `true` | Refuses non-Linux paths |

---

## Allowlists

- **Namespaces:** only listed namespaces may pass `SandboxApplyAllowed`.
- **Labels:** every key in `allowedLabels` must match exactly on the workload.

---

## What guarded apply does (when enabled)

Writes **only** a marker file under `-sandbox-dir` (e.g. `apply-marker.txt`).  
Does **not** write cgroup `cpu.max` / `memory.max` on the host in KHR-F.

---

## Rollback

`CaptureBaseline` + `RollbackBaseline` restore or remove the sandbox marker file.

---

## Blocked surfaces

- KubeVirt / libvirt / QMP
- Windows guests
- Autonomous / production apply without operator CLI + config

---

## ISO

`install_karl_host_runtime` ships unit **disabled** and does not enable apply by default.
