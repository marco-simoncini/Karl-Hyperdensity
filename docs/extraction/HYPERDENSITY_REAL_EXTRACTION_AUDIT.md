# Hyperdensity real extraction audit (Sprint 16)

**Repos:** `marco-simoncini/Karl-Dashboard` → `marco-simoncini/Karl-Hyperdensity`  
**Branches:** `KHR` / `KHR`  
**Purpose:** Map **today’s live Hyperdensity** in the Dashboard console (Go server + Next UI) to **future extractable packages** in Hyperdensity **without** adding parallel product logic or migrating runtime in this sprint.

**Sources reviewed (non-exhaustive):** `kubernetes-console/pkg/server/server.go` (route constants), `pkg/server/hyperdensity_*.go` (~350+ files), `frontend-next/src/features/hyperdensity/**` (routes under `/karl-hyperdensity`), `deployment/hyperdensity/**` (manifests, annotations, canary-style evidence namespaces).

---

## Fase 1 — Karl-Dashboard:KHR (findings)

### 1. API / runtime

| Surface | Location (indicative) | Notes |
|--------|------------------------|--------|
| `GET /api/hyperdensity/parent-fabric` | `server.go` + `hyperdensity_parent_fabric_*.go` | Full snapshot; `?view=summary` compact poll (`hyperdensity-client.ts`). |
| `POST /api/hyperdensity/parent-fabric/execution` | Same | Execution envelope; unauthenticated → **401** (documented in cohort copy). |
| `GET/PUT /api/hyperdensity/auto-scope` (+ audit) | `server.go`, client types | Policy engine, read-only reasons, audit hash metadata. |
| `GET /api/hyperdensity/executive-resource-market-cockpit` | `server.go`, client | Executive cockpit payload. |
| Auth / read-only | OpenShift console `auth` stack + handler guards | Parent-fabric VM collectors include gates such as `no_windows_lane`, `no_production_mutation`. |

### 2. Parent Fabric

- **Inventory live:** numerous `hyperdensity_parent_fabric_vm_*_cohort*.go` builders merge live vs bootstrap signals (`live_parent_fabric_observed`, inventory channels).
- **Fallback bootstrap:** cohort modules seed policy/evidence when live inventory partial.
- **Governance:** decision/arbitration summaries in snapshot model; auto-scope policy; execution engine flags.
- **Freeze / approval / budget / priority / deny:** expressed in snapshot `decisionEngine`, `arbitrationEngine`, auto-scope rules, and per-lane “blocker” narratives (strings + structured gate lists).

### 3. Linux container lanes

- Deployment / StatefulSet **CPU** reclaim, burst, rollback paths in parent-fabric merge layers (reference workloads in `karl-hyperdensity-evidence`).
- **Memory** burst / rollback symmetry.
- **Pods resize** and **no-rollout entitlement** semantics appear in annotations on reference manifests under `deployment/hyperdensity/control-plane/*.yaml` (floor/baseline/burst/ceiling, `no-autonomous-mutation`, authority modes).

### 4. VM Linux lanes

- **Narrow cohorts:** `hyperdensity_parent_fabric_vm_linux_cpu_cohort*.go`, `*_memory_cohort*.go` (+ guest-assisted executors).
- **VMI / runtime / guest overlay:** evidence collectors, readiness overlays, QGA bootstrap profiles (rhel-family).
- **KubeVirt LiveUpdate / migration blockers:** gates in runtime collector (`workload-update`, failed mount remediation refs).
- **Runtime-authority annotations:** guest-assisted control surface strings (`karl_parent_fabric_*`).
- **Second candidate rejection:** arbitration / overlap summaries in model + tests.

### 5. Hyperdensity “market”

- Donor/receiver, parent pool, overlap, execution engine, savings/value, SLO-style gates, durable history: carried in `hyperdensity-model.ts` types and parent-fabric merge helpers (`hyperdensity-runtime-parent-fabric.ts`, `hyperdensity-parent-fabric-summary-merge.ts`, execution helpers).

### 6. UI / cockpit

- **Route root:** `HYPERDENSITY_ROUTE_PATH = "/karl-hyperdensity"` (`hyperdensity-model.ts`).
- Sub-routes: overview, executive cockpit, savings-value, auto-scope, readiness, resource-movements, resize-certification, evidence, settings, exchange/equilibrium, live-authority, shells, governance (`hyperdensity-control-room-types.ts`).
- **Claim policy / forbidden claims / removed surfaces:** enforced in model tests (`__tests__/*`, forbidden word scans in artifacts pipelines referenced by tests).

### 7. Windows lane

- **Explicit exclusion:** runtime collector contributions include gate `no_windows_lane` → `windows_disabled` / `keep_windows_lane_disabled` (`hyperdensity_parent_fabric_vm_runtime_evidence_collector_v1.go`).
- **Where Windows / FluidVirt appears:** dedicated `hyperdensity_parent_fabric_windows_*.go`, `hyperdensity_windows_fluid_evidence.go`, DaaS performance / rollback certification tracks.
- **Planning vs evidence:** Windows FluidVirt lab paths remain **gated**; production mutation flags expected **false** in collector gates alongside Linux.

---

## Fase 2 — Karl-Hyperdensity:KHR (findings)

| Area | Path / artifact | Role today |
|------|-------------------|------------|
| KHR Linux skeleton | `pkg/khr/*`, `cmd/khr-linux-agent` | cgroup discovery/telemetry, evidence bundle, dry-run lease — **Linux host lane**, not Dashboard parent-fabric. |
| Grande Padre | `pkg/grandepadre/evidence`, `pkg/grandepadre/recommendation` | Local evidence store + **dry-run** recommendation slate (Sprint 12–15); **not** Dashboard execution engine. |
| Windows FluidVirt | `pkg/windowsfluidvirt/*` | Contracts, admission replay, governance, planning-only safety (Sprint 15) — **extractable** cousin to Dashboard Windows server modules; naming overlap risk. |
| CRDs / contracts | `api/crds/**`, `docs/contracts/**`, `schemas/*.json` | Canonical **vNext** shapes; align with Dashboard snapshot fields over time. |
| Migration doc | `docs/migration/dashboard-to-hyperdensity-extraction-plan.md` | High-level phases; **this audit** refines targets. |

**Duplication vs Dashboard:** conceptual overlap on “action slate”, “blockers”, “donor/receiver”, “equilibrium” — Dashboard has **full runtime**; Hyperdensity has **typed skeletons** (Grande Padre recommendation, KHR evidence). Mark Grande Padre recommendation as **vNext alignment** to Dashboard decision surfaces, not a second product.

---

## Fase 3 — Extraction matrix

| Capability | Current source in Dashboard | Current status | Runtime critical? | Extraction target in Hyperdensity | Extraction phase | Risks | Tests needed |
|------------|----------------------------|----------------|--------------------|-----------------------------------|------------------|-------|--------------|
| Parent fabric GET snapshot | `pkg/server/hyperdensity_parent_fabric*.go` | Production path | **Yes** | `pkg/hyperdensity/parentfabric` (read-only builders first) | Phase B | Drift vs UI types | Contract tests vs golden JSON subset |
| Parent fabric POST execution | Same + execution helpers | Production | **Yes** | `pkg/hyperdensity/execution` | Phase C | Auth semantics | Handler integration tests (later) |
| Executive cockpit API | `hyperdensity_executive_*` (server) + UI | Production | **Yes** | `pkg/hyperdensity/contracts` + thin `cockpit` view model | Phase B | Payload size | Snapshot diff tests |
| Auto-scope policy | Server + ConfigMap assumptions | Production | **Yes** | `pkg/hyperdensity/claimpolicy` | Phase C | Cluster write boundaries | Policy round-trip tests |
| Linux container burst / rollback | Parent-fabric merge + manifests | Production | **Yes** | `pkg/hyperdensity/equilibrium` + `blockers` | Phase B–C | Annotation coupling | Fixture tests from YAML |
| VM Linux cohort CPU/memory | `*_vm_linux_*_cohort*.go` | Production | **Yes** | `pkg/hyperdensity/cohorts` | Phase B | Guest-assisted vs readonly | Cohort unit tests (ported) |
| VM readonly observation | `*_vm_readonly_runtime_observation*.go` | Production | Medium | `pkg/hyperdensity/evidence` | Phase B | Channel auth matrix | Existing `_test.go` ports |
| Arbitration / overlap | Merge helpers + model | Production | Medium | `pkg/hyperdensity/equilibrium` | Phase B | Second candidate logic | Model tests |
| Resource futures / savings ledger | Model + server slices | Mixed | Medium | `pkg/hyperdensity/contracts` then `history` | Phase B–D | Financial wording guard | Forbidden-word + schema tests |
| SLO guardrails | Gates in collectors | Production | Medium | `pkg/hyperdensity/blockers` | Phase B | False positives | Gate catalog tests |
| Durable history | Server stores (`hyperdensityUsageHistory` in `Server` struct) | Production | **Yes** | `pkg/hyperdensity/history` | Phase D | Storage backend | Replay tests |
| Windows lane / FluidVirt | `hyperdensity_parent_fabric_windows_*.go`, `hyperdensity_windows_fluid_evidence.go` | Lab + gated | Medium (no GA) | `pkg/hyperdensity/windowslane` **or** reuse `pkg/windowsfluidvirt` | Phase E | **Over-claim** | Planning-only parity tests |
| Windows exclusion gate | `vm_runtime_evidence_collector` | Production guard | Low | `pkg/hyperdensity/blockers` (constant export) | Phase A | None | Assert gate id stable |
| UI routes `/karl-hyperdensity` | `frontend-next/src/features/hyperdensity` | Production | **Yes (UX)** | **Stays in Dashboard** until module import | Phase F | Routing | Vitest (stay in Dashboard) |
| KHR Linux agent | N/A (Hyperdensity repo) | Skeleton | N/A (different lane) | `pkg/khr/*` **stays**; link via contracts only | Phase A | Confuse with parent-fabric | Agent golden tests |
| Grande Padre recommendation | Hyperdensity `pkg/grandepadre/recommendation` | Skeleton dry-run | No | Align types to `contracts`; **freeze** feature creep | Phase A | Duplicates Dashboard “action slate” UX | Keep dry-run only tests |
| CRD ResourceLease/Future | `api/crds` | Canonical | Medium | `pkg/hyperdensity/contracts` | Phase A | Version skew | CRD round-trip + kubectl dry-run |

**Legend — phases:** A = contracts/constants; B = read-only pure functions from Dashboard; C = execution request validation (still no cluster write in lib); D = history/adapters; E = Windows lane; F = Dashboard imports Go module.

---

## Out of scope (Sprint 16)

- No runtime or handler edits in Dashboard.
- No new Hyperdensity controllers or HTTP servers.
- No removal of KubeVirt paths.
- No Windows apply enablement.

**Next sprint (suggested):** pick **one** vertical for Phase A→B (e.g. `blockers` + `parentfabric` gate IDs only), open a compatibility table Dashboard JSON ↔ Hyperdensity `contracts`, and add a single golden fixture in Hyperdensity generated from a redacted Dashboard snapshot.
