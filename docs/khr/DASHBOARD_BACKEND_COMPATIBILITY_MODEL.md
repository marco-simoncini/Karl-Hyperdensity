# Dashboard Backend Compatibility Model (KHR-BH)

| Field | Value |
|-------|-------|
| **Sprint** | KHR-BH / KHR-BI |
| **Scope** | Formal semantics for Dashboard KHR-first migration |
| **Runtime / CRD** | **No changes** |

---

## Shell / Cell-first worldview

Hyperdensity authoritative model:

| Concept | Role |
|---------|------|
| **Shell** | User-facing workload identity (desktop, app host, session container) |
| **Cell** | Runnable unit bound to a provider (one primary VMI per kubevirt-legacy cell) |
| **ShellLease** | Operator-scoped entitlement to observe or (future sprint) apply envelope changes |
| **GatewayRoute** | rdp-GW compatibility projection over ShellSession |
| **ProviderBinding** | Declares how a Cell is realized (`kubevirt.compatibility`, native-live, etc.) |

Dashboard **projects** this model over Parent Fabric state — it does not own reconcile or apply (KHR projection contract `khr-projection-v1alpha1-readonly-y`).

---

## ProviderBinding semantics

| Provider ID | Class | Meaning |
|-------------|-------|---------|
| `kubevirt.compatibility` | compatibility | VM/VMI-backed workloads; KubeVirt is implementation detail |
| `multus.legacy.transitional` | transitional | NAD/Multus network path; not long-term Shell fabric |
| `parent-fabric.observed` | observed | Non-KubeVirt objects discovered live |
| `windows.host-runtime` | native (Windows lane) | Windows host runtime projection |

**Rules:**

- Compatibility providers are **read-only** in Dashboard KHR-BH skeleton.
- `kubevirt.compatibility` must not be described as production GA or autonomous orchestration target.
- Provider binding does **not** imply CRD creation in KHR-BH.

---

## Compatibility provider semantics

| Legacy signal | Compatibility behavior |
|---------------|-------------------------|
| `linux-kubevirt-vm` / `windows-kubevirt-vm` object class | Map to Shell+Cell; badge "KubeVirt legacy" |
| `VirtualMachine` / `VirtualMachineInstance` kind | Force `kubevirt.compatibility` when no `karl.io/runtime-provider` label |
| `NetworkAttachmentDefinition` | `multus.legacy.transitional`; TP lists multus-target-fabric as **unsupported** |
| Windows pool replica | Map to ShellSession + GatewaySession (rdp-GW alignment) |

---

## Legacy VM projection semantics

| Field | Semantics |
|-------|-----------|
| `legacyKind` | Source K8s kind (`VirtualMachine`, `VirtualMachineInstance`, …) |
| `legacyRef` | `namespace/name` stable ref for technical panel |
| `technicalView` | Always expose VM/VMI namespace, name, uid (see Dashboard mapping fixture) |
| `compatibilityLayer` | `true` on all KHR-BH backend envelopes |

Projection functions in Dashboard `internal/khrcompat` mirror Hyperdensity `hyperdensityKHRProviderForObject` rules.

---

## Dashboard integration boundary

| Layer | Owner |
|-------|-------|
| CRDs / host-runtime apply | Karl-Hyperdensity + operator sprints |
| Parent Fabric discovery | Karl-Dashboard `pkg/server` (legacy, frozen routes) |
| KHR backend skeleton | Karl-Dashboard `internal/khrbackend` (KHR-BH) |
| KHR backend projection API | `GET /api/hyperdensity/khr-backend/projection` (KHR-BI) |
| Contract docs | Both repos; Hyperdensity is normative for Shell/Cell/Lease |

---

## KHR-BI: expected projection fields (Dashboard API)

When `HYPERDENSITY_KHR_BACKEND_PROJECTION_ENABLED=true` (reference/dev only):

| Field | Expected |
|-------|----------|
| `readOnly` | `true` |
| `backendModel` | `khr-first` |
| `compatibilityLayer` | `true` |
| `providerBindings[]` | Includes `kubevirt.compatibility` with `compatibility=true` |
| `shells[]` / `cells[]` | From Parent Fabric projection; kubevirt workloads use compatibility provider |
| `shellLeases[]` | Read-only; no autonomous apply |
| `gatewayRoutes[]` | Read-only; `noDisconnect` / `noRevoke` preserved |
| `resourcePorts[]` / `resourceLeases[]` / `resourceFuture[]` | Observation/simulation only |
| `accessGraphSummary` | Read-only |
| `scopeReadiness` | Scope gates; `resourceLeaseApplyEnabled=false` |
| `tpReadinessSummary` | `productionReady=false`, `autonomousOrchestration=false` |

Fixture: Karl-Dashboard `examples/khr-dashboard/khr-backend-projection-api.json`

---

## Compatibility boundary

| Inside boundary | Outside boundary |
|-----------------|------------------|
| Read-only JSON projection | CRD reconcile, host-runtime apply |
| KubeVirt as `kubevirt.compatibility` | Production namespace enablement |
| Multus/NAD as `multus.legacy.transitional` | NAD-first target fabric |
| Operator scope readiness display | Dashboard action buttons / apply UI |

---

## No action semantics

Hyperdensity expects the Dashboard KHR backend projection API to **never** expose:

- Top-level `actions` or orchestration triggers
- `applyEnabled=true` at envelope level
- `autonomousApply` or production enablement flags
- Mutating approval execution

`resourceLeases[].dryRunOnly` may be `true`; `applyState` must remain observation-oriented in TP.

---

## Related

- Karl-Dashboard `DASHBOARD_BACKEND_KHR_MIGRATION_PLAN.md`
- Karl-Dashboard `DASHBOARD_KHR_BACKEND_PROJECTION_API.md`
- `KHR_PROJECTION_V1.md` (Dashboard docs/hyperdensity)
- `RUNTIME_OBSERVATION_FEDERATION.md`
