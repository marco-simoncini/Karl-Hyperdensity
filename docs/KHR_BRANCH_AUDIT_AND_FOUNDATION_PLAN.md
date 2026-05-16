# KHR Branch Audit and Foundation Plan

| Field | Value |
|-------|-------|
| **Audit date** | 2026-05-16 |
| **Auditor** | Codex (read-only branch audit) |
| **Scope** | Five local repos on branch `KHR` only |

---

## 1. Repositories and branches verified

| Repository | Branch | HEAD (short) | Clean vs origin |
|------------|--------|--------------|-----------------|
| **Karl-Hyperdensity** | `KHR` | `528d852` | clean |
| **Karl-Dashboard** | `KHR` | `427e56c15` | untracked artifacts only |
| **FluidVirt** | `KHR` | `486e485063` | untracked artifacts only |
| **Karl-Inventory** | `KHR` | `f32ed24` | untracked artifacts only |
| **Karl-OS-ISO** | `KHR` | `352d91e8` | clean |

**Not audited / not modified:** KARL-APP, karl-directoryservice, Karl-DLP, Karl-Genesi, Karl-Installer, Karl-Licenziatore, Karl-Migration-Factory, Karl-OS-ISO_subiquity, Karl-Warden, Karl-technology/*.

**Validation spot-check (KHR Hyperdensity + Dashboard):**

- `./scripts/validate.sh` — PASS (Hyperdensity)
- `bash scripts/hyperdensity/test_hyperdensity_parity.sh` — PASS (Dashboard console)

---

## 2. Executive summary

| Dimension | Real state on KHR |
|-----------|-------------------|
| **Contract / memory** | Strong in Hyperdensity (Sprints 88–91): architecture memory, ResourceLease JSON schema, inventories |
| **KHR runtime (native)** | Partial in `pkg/khr/` — Linux agent, cgroup, telemetry, evidence, ResourceLease **dry-run** library |
| **CRD Shell/Cell/ResourceLease** | CRD YAML exists under `api/crds/` — **contract-only**, no in-repo controllers |
| **Dashboard product path** | Still **VM/KubeVirt-first** in `hyperdensity_parent_fabric_live.go` and VM runtime files; Shell/Cell via **parity tests + docs**, not dominant runtime |
| **FluidVirt** | Upstream KubeVirt fork; **multus** in virt-controller — legacy path |
| **Karl-OS-ISO** | Provisions **KubeVirt + kube-ovn + multus** — platform bootstrap, not KHR-native |
| **Karl-Inventory** | Agent/inventory docs; FluidShell module — peripheral to KHR engine |

**Recommendation: new KHR-only repository?** **No (for now).** Keep `pkg/khr/` in Karl-Hyperdensity until release boundaries force a split. Document in ADR-0004.

---

## 3. Per-repository audit

### 3.1 Karl-Hyperdensity (`KHR` @ `528d852`)

**Role:** Contract kit, Hyperdensity governance, KHR Linux MVP libraries, CRD definitions, extraction/parity documentation.

**Main areas**

| Path | Purpose |
|------|---------|
| `api/crds/runtime.karl.io/` | Shell, Cell, ShellPool, ShellLease, ResourcePort, RuntimeProvider, Host, HostPool |
| `api/crds/hyperdensity.karl.io/` | ResourceLease, ResourceFuture, Evidence* |
| `pkg/khr/` | KHR agent, cgroup, discovery, telemetry, evidence, **resourcelease/dryrun**, safety, audit |
| `pkg/grandepadre/` | Evidence store, recommendation, **Action Slate** model |
| `pkg/hyperdensity/contractkit/` | Parity / contractkit module |
| `docs/khr/` | KHR Linux MVP, dry-run, cgroup, evidence runbooks |
| `docs/contracts/khr/` | ResourceLease JSON Schema + 3 examples (Sprint 91) |
| `docs/extraction/` | Parent Fabric adapter extraction Sprints 78–91 |
| `docs/adr/` | ADR-0001 Shell/Cell, ADR-0002 KubeVirt legacy, ADR-0003 ResourceLease |
| `examples/providers/kubevirt/` | Legacy provider examples |

**KHR / Shell / Cell / ResourceLease / ResourcePort**

- **Implemented (library/contract):** `pkg/khr/crdv1alpha1/{cell,shell,resourcelease,resourceport,runtimeprovider}.go`, `pkg/khr/resourcelease/dryrun.go` (dry-run + rollback plan fields), CRD manifests.
- **Documented:** Sprint 88–91 architecture + ResourceLease schema; `docs/architecture/KARL_SHELL_CELL_MODEL.md`.
- **Not implemented:** In-cluster controllers reconciling Shell→Cell; KHR host runtime replacing virt-launcher.

**KubeVirt references**

- Explicitly **legacy provider** in ADR-0002, provider docs, examples — **aligned** with architecture.
- Risk: newcomers may read `examples/providers/kubevirt/` before `docs/khr/` — mitigate with ADR-0004 and this plan.

**Multus / NAD**

- Almost absent in Hyperdensity (good). OVN mentioned as **compatibility mapping** in Sprint 88–90 docs only.

**Tests**

- `go test ./...` — used by `validate.sh`
- `pkg/khr/**_test.go` — agent, cgroup, dryrun, evidence
- Contractkit parity tests

**Gap vs roadmap**

| Phase | Gap |
|-------|-----|
| 0 | Extraction complete for apply/resource_exchange; broad observation off |
| 1 | CRDs exist; need controller sketch + OpenAPI alignment with `resourcelease.schema.json` |
| 2 | No unified Shell API server in Hyperdensity |
| 3 | KHR agent exists; needs production host integration path |
| 4 | Provider contract docs exist; no automated provider adapter in Hyperdensity |
| 5 | OVN-native fabric not in code — docs only |
| 6 | Storage semantics in schema/docs; no Storage Fabric controller |
| 7+ | Windows/AppShell via Dashboard + FluidVirt, not KHR-native |

---

### 3.2 Karl-Dashboard (`KHR` @ `427e56c15`)

**Role:** Parent Fabric runtime (OpenShift console), live discovery, VM lanes, rollback, Windows DaaS vertical slices, Hyperdensity parity tests.

**Main areas**

| Path | Purpose |
|------|---------|
| `pkg/server/hyperdensity_parent_fabric_live.go` | **KubeVirt VM/VMI discovery** (`/apis/kubevirt.io/v1/...`) |
| `pkg/server/hyperdensity_parent_fabric_vm_linux_*.go` | VM CPU/memory runtime (legacy surface) |
| `pkg/server/hyperdensity_parent_fabric_windows_*.go` | Windows pool/VMPool governance |
| `pkg/server/hyperdensity_parent_fabric_rollback.go` | Rollback observed-state (legacy) |
| `pkg/server/hyperdensity_parent_fabric_resource_exchange_*.go` | Resource exchange wrappers (Sprint 83–86 activated) |
| `pkg/server/hyperdensity_khr_*.go` | **Sprint 88–91 parity tests only** (no production KHR runtime) |
| `scripts/hyperdensity/test_hyperdensity_parity.sh` | Large parity gate |

**KHR / Shell / Cell**

- **Tests:** `hyperdensity_khr_architecture_memory_test.go`, ResourceLease contract tests, inventory tests — **guardrails**, not product API.
- **Runtime:** Shell claim templates (`hyperdensity_parent_fabric_shell_claim_*.go`) — partial Shell vocabulary.
- **Gap:** No first-class “create Shell → bind Cell via KHR” API response path replacing VM discovery.

**KubeVirt**

- **Heavy:** `live.go`, VM runtime files, workload adapter VM paths — **expected for Phase 0** but **architectural regression** if treated as target engine.
- Parity tests correctly label KubeVirt as compatibility/fallback.

**Multus / kube-ovn**

- **Operational reality:** `hyperdensity_parent_fabric_vm_linux_cpu_burst.go` references **kube-ovn / multus sandbox failures**; `live.go` has `KubeOVNNotReady`.
- **Regression if:** documented as target network architecture (currently transitional in KHR docs only).

**Tests**

- Hyperdensity parity: PASS on KHR
- Hundreds of `hyperdensity_parent_fabric_*_test.go` files

**Gap**

- Wire ResourceLease schema to Parent Fabric summary (read-only) before mutating VM paths.
- Separate “Shell/Cell observation” track from VM runtime track explicitly in code comments/docs.

---

### 3.3 FluidVirt (`KHR` @ `486e485063`)

**Role:** KubeVirt-derived virtualization stack (fork).

**Main areas**

| Path | Notes |
|------|-------|
| `pkg/virt-controller/services/template.go` | **multus** `NetworkToResource` |
| `pkg/virt-launcher/`, `pkg/virt-api/` | VM execution |
| `vendor/` | Large upstream tree |

**KHR / Shell / Cell**

- **No KHR engine code.** FluidVirt is compatibility layer infrastructure.

**KubeVirt**

- **Is** the codebase — expected for Phase 4 legacy provider.

**Multus**

- **Present** in production code paths — **transitional**, must not become architecture target.

**Tests**

- Full KubeVirt test suite (heavy); not run in this audit.

**Gap**

- Document FluidVirt role as **provider backend only** in foundation plan.
- Future: thin adapter interface toward KARL Network Fabric (no refactor in this sprint).

---

### 3.4 Karl-Inventory (`KHR` @ `f32ed24`)

**Role:** Inventory agents, FluidShell witness module docs, benchmark artifacts.

**Main areas**

| Path | Notes |
|------|-------|
| `inventory/agent-windows/` | Windows agent, FluidShell module |
| `docs/fluid-shell-module.md` | Shell module documentation |
| `artifacts/` | Benchmark / witness artifacts |

**KHR / ResourceLease**

- Peripheral mentions only; no ResourceLease implementation.

**Gap**

- Align inventory telemetry with KHR evidence bundle format (`pkg/khr/evidence`).

---

### 3.5 Karl-OS-ISO (`KHR` @ `352d91e8`)

**Role:** OS image / provisioning — installs Kubernetes, KubeVirt, kube-ovn, multus, hostpath, etc.

**Main areas**

| Path | Notes |
|------|-------|
| `Stage-2/custom-files/provision/.../install_kubevirt.sh` | KubeVirt install |
| `Stage-2/.../install_kube_ovn.sh` | kube-ovn |
| `scripts/vendor_multus_resources.sh` | Multus vendoring |
| `KUBE_OVN_VM_NETWORKING.md` | VM networking doc |
| `Stage-2/.../karl-engine/kubevirt/` | Patched kubevirt YAML |

**KHR**

- “KARL Engine” path in ISO tree is **KubeVirt packaging**, not `pkg/khr` agent — **naming regression risk**.

**Gap**

- Rename/document `karl-engine` directory as **platform bootstrap** vs **KHR native runtime**.
- Phase 5+: optional KHR agent install hook without replacing KubeVirt install yet.

---

## 4. Roadmap phase assessment (real advancement)

| Phase | Description | Status on KHR | Evidence |
|-------|-------------|---------------|----------|
| **0** | Do not break existing product | **In progress** | Dashboard VM paths intact; extraction boundaries closed; parity PASS |
| **1** | CRD Shell/Cell | **Partial** | CRD YAML + Go types; no controllers |
| **2** | Shell API above all | **Minimal** | Dashboard Shell claim APIs; no unified Shell control plane |
| **3** | KHR MVP Linux | **Partial** | `pkg/khr/agent`, cgroup, telemetry, dry-run; not wired to Dashboard production |
| **4** | KubeVirt legacy provider | **Documented + operational** | ADR-0002, FluidVirt, ISO, Dashboard live discovery |
| **5** | Network Fabric OVN-native | **Docs only** | Sprint 88–90 OVN mapping; Dashboard still kube-ovn operational |
| **6** | Storage Fabric | **Docs + schema** | ephemeralOverlay etc. in JSON schema; no fabric controller |
| **7+** | Windows/AppShell/Hyperdensity integration | **Partial** | Windows vertical slices in Dashboard; not AppShell-native storage/network |

---

## 5. File-level evidence (high signal)

| File | Repo | Signal |
|------|------|--------|
| `api/crds/runtime.karl.io/shell.yaml` | Hyperdensity | Shell CRD contract-only |
| `api/crds/hyperdensity.karl.io/resourcelease.yaml` | Hyperdensity | ResourceLease CRD |
| `pkg/khr/resourcelease/dryrun.go` | Hyperdensity | Dry-run + rollback plan structure |
| `pkg/grandepadre/recommendation/action_slate.go` | Hyperdensity | Action Slate model |
| `docs/contracts/khr/resourcelease.schema.json` | Hyperdensity | Non-applied schema Sprint 91 |
| `hyperdensity_parent_fabric_live.go` | Dashboard | KubeVirt API discovery |
| `hyperdensity_parent_fabric_vm_linux_cpu_burst.go` | Dashboard | multus/kube-ovn blocker strings |
| `pkg/virt-controller/services/template.go` | FluidVirt | multus integration |
| `install_kube_ovn.sh`, `vendor_multus_resources.sh` | OS-ISO | Legacy network bootstrap |

---

## 6. Architectural regressions detected

| ID | Severity | Finding |
|----|----------|---------|
| **AR-1** | P1 | Dashboard production path is **VM/KubeVirt-first** while KHR docs claim Shell/Cell product model — risk of building “KHR as KubeVirt wrapper” |
| **AR-2** | P1 | Karl-OS-ISO `karl-engine/kubevirt` naming implies KHR == KubeVirt |
| **AR-3** | P2 | multus/kube-ovn treated as **operational default** in Dashboard errors; must stay labeled transitional in all new docs/code comments |
| **AR-4** | P2 | ResourceLease JSON schema not yet wired to CRD OpenAPI or Dashboard read models |
| **AR-5** | P3 | Windows DaaS performance artifacts in Inventory/Dashboard — ensure claims are benchmark-contract labeled (not production SLAs) |

**Not regressions (correct):**

- KubeVirt as **compatibility provider** in ADR-0002 and Sprint 88–91 docs.
- Broad observation flags remain false in Dashboard tests.

---

## 7. Gap backlog

### P0 (foundation blockers)

| ID | Gap | Owner repo |
|----|-----|------------|
| P0-1 | Single written **source of truth** linking CRD ↔ JSON schema ↔ Dashboard (ADR-0004 + this doc) | Hyperdensity |
| P0-2 | Dashboard **Shell/Cell read path** design (even read-only) before more VM features | Dashboard |
| P0-3 | Explicit **network target** doc in OS-ISO/FluidVirt: kube-ovn/multus transitional | All |

### P1 (next quarter)

| ID | Gap | Owner repo |
|----|-----|------------|
| P1-1 | Align `resourcelease.schema.json` with CRD OpenAPI in Hyperdensity | Hyperdensity |
| P1-2 | KHR agent install/runbook on bare metal (from `docs/khr/`) | Hyperdensity + OS-ISO |
| P1-3 | Parent Fabric ResourceLease **read-only** projection in API summary | Dashboard |
| P1-4 | Storage ephemeral semantics in Windows vertical slice docs | Dashboard |

### P2 (later)

| ID | Gap |
|----|-----|
| P2-1 | OVN-native Network Fabric controller design |
| P2-2 | Storage Fabric promote-to-image workflow |
| P2-3 | Split `pkg/khr` to separate repo (only if release cadence requires) |

---

## 8. New KHR repository decision

| Option | Recommendation |
|--------|----------------|
| Create new repo now | **No** |
| Rationale | `pkg/khr/` cohesive in Hyperdensity; CRDs and contracts co-located; splitting now adds release friction without controllers to ship |

**Revisit when:** KHR agent binaries ship on independent cadence from Hyperdensity CRDs (Phase 3 GA).

---

## 9. Next Codex sprint — atomic tasks

**Sprint A — Contract alignment (Hyperdensity + Dashboard docs/tests)**

1. Diff `resourcelease.schema.json` vs `api/crds/hyperdensity.karl.io/resourcelease.yaml` OpenAPI; document deltas.
2. Add Dashboard read-only fixture test: Parent Fabric summary includes ResourceLease-shaped stub.
3. Update `docs/khr/KHR_LINUX_AGENT_RUNBOOK.md` with explicit “not production in Dashboard yet”.

**Sprint B — Shell/Cell observation track (Dashboard)**

1. Design `hyperdensity_parent_fabric_shell_observation_*.go` stub package (flags off).
2. Inventory VM→Shell mapping in `live.go` (read-only table).

**Sprint C — Network transitional labeling (OS-ISO + FluidVirt docs only)**

1. Add `docs/NETWORK_TRANSITIONAL_KUBE_OVN.md` in Hyperdensity referencing ISO + FluidVirt paths.
2. No code changes in FluidVirt.

**Sprint D — KHR agent smoke (Hyperdensity)**

1. Extend `cmd/` or document `khr-agent` smoke test on single host.
2. Wire dry-run CLI example using `pkg/khr/resourcelease/dryrun.go`.

---

## 10. Related artifacts

- `docs/adr/ADR-0004-khr-foundation-source-of-truth.md`
- `docs/adr/ADR-0001-khr-shell-cell-runtime-model.md`
- `docs/adr/ADR-0002-kubevirt-as-legacy-provider.md`
- `docs/extraction/HYPERDENSITY_KHR_ARCHITECTURE_MEMORY.md`
- `docs/contracts/khr/resourcelease.schema.json`
