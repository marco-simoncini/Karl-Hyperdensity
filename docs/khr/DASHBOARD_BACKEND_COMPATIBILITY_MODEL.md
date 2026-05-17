# Dashboard Backend Compatibility Model (KHR-BH)

| Field | Value |
|-------|-------|
| **Sprint** | KHR-BH … KHR-CG / **KHR-CH** |
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

## KHR-BJ: source-of-truth expectations

| Consumer | Must use |
|----------|----------|
| TP readiness fields | `buildHyperdensityKHRReadinessContractV1` (Dashboard) |
| Federation / Hyperdensity preflight | Fixtures under `examples/khr-dashboard/khr-backend-projection-*.json` |
| Reference env | `khr-backend-projection-reference-env-enabled.json` |

Hyperdensity validation scripts MUST NOT treat legacy VM-first APIs as authoritative when `khr-backend/projection` is enabled in reference env.

---

## KHR-BJ: anti-regression requirements

Dashboard CI enforces via `scripts/validate_khr_backend_projection_contract.sh`:

1. **No action semantics** — no top-level mutation/orchestration fields.
2. **KubeVirt compatibility** — any `kubevirt` provider id requires `compatibility=true`.
3. **Lease apply disabled** — `scopeReadiness.resourceLeaseApplyEnabled=false` and `tpReadinessSummary.resourceLeaseApplyEnabled=false`.
4. **Feature flag default** — `HYPERDENSITY_KHR_BACKEND_PROJECTION_ENABLED=false` unless reference/dev.
5. **Golden stability** — `khr_backend_projection_golden.json` updated only with `KHR_UPDATE_GOLDEN=true`.

Regression in any of the above is a **contract freeze violation** for TP.

---

## KHR-BK: provider profile expectations

Dashboard deployment profiles (docs/fixtures only — **no Hyperdensity CRD/runtime change**):

| Dashboard profile | Installer / ISO profile | KubeVirt boundary |
|-------------------|----------------------|-------------------|
| `khr-native` | `karl2-khr-technical-preview` | **Not required** for projection; optional `kubevirt.compatibility` only |
| `public-cloud-kubevirt-compatibility` | `karl1-kubevirt-legacy` | **Required binding** `kubevirt.compatibility` with `compatibility=true` |
| `hybrid-transition` | `hybrid-transition` | Same as public cloud + KHR CRD foundation |

### Native vs public cloud semantics

| Dimension | `khr-native` | `public-cloud-kubevirt-compatibility` |
|-----------|--------------|----------------------------------------|
| Primary worldview | Shell / Cell / Lease | Shell / Cell / Lease (KHR-first) |
| KubeVirt role | Optional compatibility provider | Compatibility provider for VM fleets |
| Future default | **Yes** (greenfield) | **No** — preserves existing installs |
| ISO `defaultInstallerProfile` | Unchanged (`karl1-kubevirt-legacy`) | Unchanged |

### KubeVirt compatibility boundary

1. Hyperdensity and Dashboard **never** treat KubeVirt as the primary Shell model.
2. All kubevirt provider ids in projection JSON require `compatibility=true` (KHR-BJ guard).
3. Public cloud / hybrid profiles **may** list `kubevirt.compatibility` in required bindings; native profile **must not**.
4. Retaining KubeVirt in cluster/installer does **not** imply production enablement or autonomous orchestration.

Fixtures: Karl-Dashboard `examples/khr-dashboard/provider-profile-*.json`  
Normative Dashboard doc: `DASHBOARD_PROVIDER_PROFILE_MODEL.md`

---

## KHR-BL: expected provider profile projection fields

Dashboard `GET /api/hyperdensity/khr-backend/projection` exposes read-only provider profile readiness (no Hyperdensity CRD/runtime change):

| Field | `khr-native` | `public-cloud-kubevirt-compatibility` | `hybrid-transition` |
|-------|--------------|---------------------------------------|----------------------|
| `kubevirtRequired` | `false` | `true` | `true` |
| `supportsKhrNative` | `true` | `false` | `true` |
| `supportsKubeVirtCompatibility` | `true` | `true` | `true` |
| `providerProfileCompatibilityMode` | `khr-native` | `kubevirt.compatibility` | `kubevirt.compatibility` |
| `kubevirtProviderBinding` | `kubevirt.compatibility` (optional) | `kubevirt.compatibility` | `kubevirt.compatibility` |

Same fields are mirrored on `tpReadinessSummary`. Resolution is Dashboard-side (`HYPERDENSITY_KHR_PROVIDER_PROFILE` / `KARL_INSTALLER_PROFILE` / documented default).

Anti-regression: Karl-Dashboard `assertKHRProviderProfileProjectionContractV1` + `validate_khr_backend_projection_contract.sh`.

---

## KHR-BM: profile propagation chain expectations

Canonical chain (read-only; no Hyperdensity runtime change):

`KARL_INSTALLER_PROFILE` → ISO `profile-manifest.yaml` (`dashboardProviderProfile`) → `HYPERDENSITY_KHR_PROVIDER_PROFILE` → Dashboard `khr-backend/projection`.

| Checkpoint | Expected field / rule |
|------------|----------------------|
| Profile source | `providerProfileSource` always explicit (`env:…`, `default-documented`, `fallback-invalid-profile`) |
| Compatibility boundary | Public cloud / hybrid: `kubevirt.compatibility` with `compatibility=true` in bindings |
| No KubeVirt primary model | Shell/Cell-first; KubeVirt is ProviderBinding only |
| Unsafe fallback | Invalid dashboard env → `public-cloud-kubevirt-compatibility` — **never** `khr-native` |
| Unset default | `default-documented` → `public-cloud-kubevirt-compatibility` (existing installs unchanged) |

### KHR-BM checklist

- [ ] `providerProfileSource` present on projection + `tpReadinessSummary`
- [ ] `fallback-invalid-profile` never yields `khr-native`
- [ ] `public-cloud-kubevirt-compatibility` has `kubevirtRequired=true` and `kubevirt.compatibility` mode
- [ ] Dashboard env wins over installer env when both set
- [ ] No action buttons / mutation / production enable in projection JSON

Normative Dashboard doc: `DASHBOARD_PROVIDER_PROFILE_PROPAGATION.md`  
Fixtures: `provider-profile-propagation-*.json`

---

## KHR-BY: Dashboard UI consumption expectations

Dashboard remains the **same cockpit**; Hyperdensity does not change CRDs or runtime in BY. UI consumption is a **read-only projection consumer** in Karl-Dashboard only.

### Source of truth for UI (when enabled)

| Layer | Source | Flag (Dashboard) | Default |
|-------|--------|------------------|---------|
| Backend envelope | `GET /api/hyperdensity/khr-backend/projection` | `HYPERDENSITY_KHR_BACKEND_PROJECTION_ENABLED` | `false` |
| Cockpit view model | `internal/khrui.AdaptBackendProjection` | `HYPERDENSITY_KHR_UI_PROJECTION_ENABLED` | `false` |

Hyperdensity normative model (Shell / Cell / ProviderBinding) is projected by Dashboard; the UI adapter maps that envelope to existing column contracts (`shell-list-view-model.json`) without new Dashboard product or layout rewrite.

### Compatibility provider views

| Provider ID | UI badge / row semantics |
|-------------|--------------------------|
| `kubevirt.compatibility` | `legacyProviderBadge`: "KubeVirt legacy"; `runtimeProviderId` = binding id |
| `multus.legacy.transitional` | "Multus transitional" |
| `parent-fabric.observed` | "Compatibility provider" |
| Windows RDP lane | `accessSummary.protocol=RDP` when legacy class indicates Windows |

Cells attach to shells via `primaryCellRef` (namespace/name). Provider profile fields on the projection (`providerProfile`, `kubevirtRequired`, etc.) are **display-only** in BY — no installer or cluster mutation.

### No action semantics (UI)

KHR-BY requires:

- No cockpit **action buttons** from KHR adapter (`actionCount=0`, `ExposedActions()` empty).
- No `apply`, `disconnect`, `revoke`, or orchestration triggers in UI-adapted JSON.
- `readOnly=true` and `noMutation=true` on adapter output.
- When `HYPERDENSITY_KHR_UI_PROJECTION_ENABLED=false`, legacy VM/VMI/Pool discovery path is **unchanged** (`legacyPathUsed=true`).

Inventory observation on projection (`inventoryObservationSummary`) remains read-only; `uiConsumption=false` in committed ingest evidence until a future UI sprint.

### Hyperdensity validation (BY)

- No runtime / CRD / operator change in this sprint.
- `scripts/validate.sh` offline mode continues to pass with updated compatibility doc only.

Normative Dashboard doc: `DASHBOARD_UI_KHR_PROJECTION_CONSUMPTION_PLAN.md`  
Adapter package: `kubernetes-console/internal/khrui/`

---

## KHR-CB: UI projection preview LIVE_PASS (reference env)

| Artifact | Expected |
|----------|----------|
| `evidenceStatus` | `LIVE_PASS` |
| `source` | `live-readonly` |
| `dataSource` | `khr-backend-projection` |
| `legacyPathUsed` | `false` |
| `providerProfile` | `khr-native` |
| `actionCount` | `0` |
| `rollbackVerified` | `true` (mandatory rollback after evidence) |

Committed: `docs/evidence/khr-dashboard-ui-projection-preview/committed-khr-cb-v1/`  
Orchestrator: `khr_dashboard_ui_projection_preview_live_rollout.sh` (KHR-BS rollout pattern).

---

## KHR-CA: Dashboard evidence classification (backend vs UI preview)

| Evidence type | Endpoint | Karl-Dashboard path | Typical `evidenceStatus` |
|---------------|----------|---------------------|--------------------------|
| **Backend projection live** | `GET /api/hyperdensity/khr-backend/projection` | `docs/evidence/khr-dashboard-reference-env-activation/` | `LIVE_PASS` (KHR-BS) |
| **UI projection preview live** | `GET /api/hyperdensity/khr-ui/projection-preview` | `docs/evidence/khr-dashboard-ui-projection-preview/` | `LIVE_PASS` / `REMEDIATION_PASS` (KHR-CA) |
| **UI preview fixture** | (offline) | `committed-khr-ca-v1` | `FIXTURE_PASS` |

Rules:

- Backend evidence proves Shell/Cell/Lease **backend envelope** and `khr-native` activation on reference console.
- UI preview evidence proves **cockpit shell-list adapter** output (`shellRows`) with dual-flag gating; does **not** migrate workload pages.
- `REMEDIATION_PASS` is valid when live image/env is not yet on dual-flag preview — remediation plan only, no Hyperdensity runtime change.
- Both evidence classes require `readOnly=true`, `actionCount=0`, no top-level mutation fields.

Normative Dashboard doc: `DASHBOARD_UI_KHR_PROJECTION_PREVIEW_EVIDENCE.md`

---

## KHR-BZ: UI projection preview semantics

Dashboard exposes read-only UI consumption preview without cockpit layout changes:

| Endpoint | Role |
|----------|------|
| `GET /api/hyperdensity/khr-backend/projection` | Backend Shell/Cell/Lease envelope (KHR-BI+) |
| `GET /api/hyperdensity/khr-ui/projection-preview` | Cockpit shell-list rows via `internal/khrui` (KHR-BZ) |

### Dual-flag gating (normative)

Preview adapter path is active **only** when:

1. `HYPERDENSITY_KHR_BACKEND_PROJECTION_ENABLED=true`
2. `HYPERDENSITY_KHR_UI_PROJECTION_ENABLED=true`

If either flag is `false` (default), response uses `dataSource=legacy-parent-fabric`, `legacyPathUsed=true`, and **empty** `shellRows` — legacy VM/VMI/Pool UI paths are unchanged.

### UI consumption preview

- Maps `khr-backend/projection` shells/cells to cockpit `shellRows` (column contract in `shell-list-view-model.json`).
- `providerProfile` (e.g. `khr-native`) is display-only on preview JSON.
- KubeVirt workloads use `kubevirt.compatibility` with `legacyProviderBadge` "KubeVirt legacy".

### No action semantics (preview)

- `actionCount=0`, `noMutation=true`, `readOnly=true` on every response.
- No `actions`, apply, disconnect, revoke, or orchestration fields.
- Hyperdensity does not add runtime apply or CRD changes for BZ.

Normative Dashboard doc: `DASHBOARD_UI_KHR_PROJECTION_PREVIEW.md`  
Fixture: `ui-projection-preview-khr-native.json`

---

## KHR-CG: Cockpit component migration phase 1 (read-only)

First **cockpit component** consumption path — shell/workload list only; no action semantics.

| Endpoint | Role |
|----------|------|
| `GET /api/hyperdensity/khr-ui/cockpit-component-preview` | Shell/workload list rows + read-only badges (KHR-CG) |

### Triple-flag gating (normative)

Component preview is active **only** when all are `true`:

1. `HYPERDENSITY_KHR_BACKEND_PROJECTION_ENABLED`
2. `HYPERDENSITY_KHR_UI_PROJECTION_ENABLED`
3. `HYPERDENSITY_KHR_UI_COMPONENT_PREVIEW_ENABLED` (default **false**)

If component flag is `false`, response uses `dataSource=legacy-cockpit-component`, `legacyPathUsed=true`, empty `rows` — **legacy cockpit component path unchanged**.

### Read-only component consumption

| Field | Semantics |
|-------|-----------|
| `component` | `shell-workload-list` (phase 1 target) |
| `rows[].providerProfileBadge` | e.g. `KHR-native` |
| `rows[].compatibilityBadge` | e.g. `KubeVirt legacy` |
| `rows[].readinessBadge` | e.g. `TP-reference-not-production` |
| `actionCount` | Always `0` |
| `noMutation` | Always `true` |

No cockpit layout rewrite; adapter maps projection → existing list view model shape.

Normative Dashboard doc: `DASHBOARD_COCKPIT_COMPONENT_MIGRATION_PLAN.md`  
Fixture: `cockpit-component-preview.json`

---

## KHR-CH: Cockpit component preview LIVE_PASS (reference env)

| Item | Detail |
|------|--------|
| Evidence | Karl-Dashboard `docs/evidence/khr-dashboard-cockpit-component-preview/committed-khr-ch-v1/` |
| Endpoint | `GET /api/hyperdensity/khr-ui/cockpit-component-preview` |
| `evidenceStatus` | `LIVE_PASS` on `karl-metal-01@ovh` with mandatory rollback |
| `uiPageMigrationPending` | **true** — API/adapter proof only; cockpit page still legacy |
| `frontendRewrite` | **false** |

Live proof: `source=live-readonly`, `component=shell-workload-list`, `dataSource=khr-backend-projection`, badges present, `actionCount=0`, `rollbackVerified=true`.

---

## Related

- Karl-Dashboard `DASHBOARD_BACKEND_KHR_MIGRATION_PLAN.md`
- Karl-Dashboard `DASHBOARD_UI_KHR_PROJECTION_CONSUMPTION_PLAN.md`
- Karl-Dashboard `DASHBOARD_UI_KHR_PROJECTION_PREVIEW.md`
- Karl-Dashboard `DASHBOARD_COCKPIT_COMPONENT_MIGRATION_PLAN.md`
- Karl-Dashboard `DASHBOARD_PROVIDER_PROFILE_MODEL.md`
- Karl-Dashboard `DASHBOARD_PROVIDER_PROFILE_PROPAGATION.md`
- Karl-Dashboard `DASHBOARD_KHR_BACKEND_PROJECTION_API.md`
- `KHR_PROJECTION_V1.md` (Dashboard docs/hyperdensity)
- `RUNTIME_OBSERVATION_FEDERATION.md`
