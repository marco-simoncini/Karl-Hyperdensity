# ResourcePort contract (KHR-C observation)

| Field | Value |
|-------|-------|
| **API group** | `runtime.karl.io/v1alpha1` |
| **Kind** | `ResourcePort` |
| **Scope** | Cluster (capability profiles) + observation binding via refs |
| **Status** | Contract-only — no controller |

---

## Purpose

**ResourcePort** is capability truth for a Shell/Cell binding:

- Which **provider** fulfills the port
- **Hotplug** posture per resource class (cpu, memory, disk, network)
- **Constraints** (policy gates, platform limits)
- **Conditions** + **observedAt** for read-only inventory

ResourceLease dry-run uses the legacy `spec.ports` matrix; observation adds `shellRef` / `cellRef` for Dashboard and Inventory alignment.

---

## Spec (observation)

| Field | Required | Description |
|-------|----------|-------------|
| `provider` | yes | e.g. `khr.native`, `kubevirt.compatibility`, `parent-fabric.observed` |
| `shellRef` | yes | `namespace/Shell/name` or `Shell/name` |
| `cellRef` | yes | `namespace/Cell/name` |
| `capabilities` | yes | Capability tokens (`cpu.envelope`, `memory.static`, …) |
| `hotplug` | yes | `{ cpu, memory, disk, network }` booleans |
| `constraints` | no | Opaque policy object |
| `ports` | yes* | Legacy matrix for dry-run (*required on cluster profiles) |

---

## Status (observation)

| Field | Description |
|-------|-------------|
| `observedAt` | RFC3339 timestamp |
| `conditions` | `{ type, status, reason, message }[]` |

---

## Inventory candidacy

**Karl-Inventory** is the candidate **source of ResourcePort facts** (host posture, guest telemetry, FluidShell evidence) — see `Karl-Inventory/docs/khr/INVENTORY_OBSERVATION_CONTRACT.md`. No agent changes in KHR-C.

---

## Artifacts

- `docs/contracts/khr/resourceport.schema.json`
- `docs/contracts/khr/examples/resourceport-linux-container.json`
- `api/crds/runtime.karl.io/resourceport.yaml`
