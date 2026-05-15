# ResourcePort Contract (v1alpha1)

**Kind:** `ResourcePort` (`resourceports.runtime.karl.io`)  
**Scope:** Cluster

## Role

`ResourcePort` is the **capability truth matrix** for a class of Cells. Hyperdensity (Grande Padre) and future KHR apply paths must treat advertised modes as **upper bounds** of what may be attempted—not guarantees that every host, guest, or cloud region can deliver them simultaneously.

## Schema summary

`spec.ports` defines supported **mode sets** per resource dimension:

| Dimension | Example `modes` values |
|-----------|-------------------------|
| `cpu` | `hotAdd`, `envelope`, `static` |
| `memory` | `hotAdd`, `balloon`, `virtioMem`, `envelope`, `static` |
| `disk` | `hotplug`, `static` |
| `network` | `hotplug`, `static` |
| `gpu` | `coldAttach`, `warmAttach`, `liveIfSupported` |

`cpu` and `memory` are **required** keys in v1alpha1 to force explicit statements for Hyperdensity’s primary economic surface.

## Relationship to RuntimeProvider

`spec.appliesToRuntimeProviderIds` optionally narrows the profile to specific `RuntimeProvider.spec.id` values. If empty, the profile is interpreted as **generic** (controllers may still restrict usage).

## Non-goals (explicit)

- This CRD **does not** implement enforcement.
- It does **not** promise universal hotplug; `notes` exists for human-readable caveats only.
- Cloud adaptive providers may expose **strict subsets**; see `RESOURCELEASE_RESOURCEFUTURE_CONTRACT.md`.

## Evolution

v1alpha1 uses discrete string enums for common modes. Additional modes may require CRD schema revision or migration to a `v1beta1` enum superset with conversion.
