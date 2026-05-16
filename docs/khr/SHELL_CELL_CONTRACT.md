# Shell / Cell / ShellClass API foundation (KHR-D)

| Field | Value |
|-------|-------|
| **API group** | `runtime.karl.io/v1alpha1` |
| **Kinds** | `Shell`, `Cell`, `ShellClass` |
| **Controllers** | None (contract + CRD install only) |

---

## Purpose

Define **installable** Shell/Cell primitives as the KHR API foundation. Production VM/KubeVirt paths remain **compatibility** surfaces until controllers ship.

---

## Shared spec fields (Shell + Cell)

| Field | Description |
|-------|-------------|
| `providerBinding` | Runtime provider enum (`khr.native`, `kubevirt.compatibility`, …) |
| `runtimeClass` | Execution profile (`kubernetes-workload`, `kubevirt-vm`, …) |
| `owner` | `{ user, tenant }` tenancy |
| `resources` | CPU/memory request/limit strings |
| `networkRefs` | Named network attachments |
| `storageRefs` | Named storage bindings + optional `mode` |

## Shared status fields (Shell + Cell)

| Field | Description |
|-------|-------------|
| `phase` | `Pending`, `Provisioning`, `Ready`, `Degraded`, `Terminating`, `Observed` |
| `conditions` | Standard condition list |
| `observedResourcePorts` | Refs to observed ResourcePort objects |
| `observedResourceLeases` | Refs to observed ResourceLease objects |

---

## ShellClass

Cluster-scoped class template: `id`, `complianceFamily`, `defaultProviderBinding`, `defaultRuntimeClass`, optional default resource/network/storage templates.

---

## Artifacts

| Path | Role |
|------|------|
| `api/crds/runtime.karl.io/shell.yaml` | Shell CRD |
| `api/crds/runtime.karl.io/cell.yaml` | Cell CRD |
| `api/crds/runtime.karl.io/shellclass.yaml` | ShellClass CRD |
| `docs/contracts/khr/shell.schema.json` | JSON schema |
| `docs/contracts/khr/cell.schema.json` | JSON schema |
| `docs/contracts/khr/shellclass.schema.json` | JSON schema |
| `docs/contracts/khr/examples/shell-linux-dev.json` | Example |
| `docs/contracts/khr/examples/cell-linux-container.json` | Example |
| `docs/contracts/khr/examples/shellclass-linux-dev.json` | Example |

---

## ISO packaging

`Karl-OS-ISO` ships CRDs under `/opt/karl/karl-engine/khr/crds/` via `install_khr_crds.sh` (read-only apply, no controller).
