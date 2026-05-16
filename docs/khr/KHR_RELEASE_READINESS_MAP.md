# KHR release readiness map (KHR-Z)

| Field | Value |
|-------|-------|
| **Sprint** | KHR-Z |
| **Status** | Planning / sandbox evidence — **not GA**, **not production-ready** |

Readiness labels: **Done (sandbox)** | **Preview** | **Gap** | **Blocked**

---

## Capability matrix

| Area | Sprint | Status | Notes |
|------|--------|--------|-------|
| API foundation (Host, Shell, Cell, ResourcePort, ResourceLease) | A–E | Done (sandbox) | CRDs + contracts; read-only Dashboard projection |
| ResourceFuture simulation | R | Done (sandbox) | `resourcefuture-simulate`; no apply |
| Multi-lane discovery | Q | Done (sandbox) | Native-live + compatibility lanes |
| Native-live lane | S | Done (sandbox) | Cgroup live scale; `khr-runtime-sandbox` |
| Native-live certification | T | Done (sandbox) | Multi-run certify; cluster evidence |
| Shell continuity | U | Done (sandbox) | Semantics + certification proofs |
| Certification registry | V | Done (sandbox) | Registry + policy gates on ResourceFuture |
| Operator action approval | W | Done (sandbox) | Local evidence only; no apply on approve |
| Control graph | X | Done (sandbox) | Shell/Cell-first graph export |
| Trust / provenance | Y | Done (sandbox) | Fingerprints + lineage validation |
| Infrastructure scope language | Z | Preview | This sprint — doc alignment |

---

## Cross-repo status

| Repo | Role | Status | Gap for technical preview |
|------|------|--------|---------------------------|
| **Karl-Hyperdensity** | Source of truth, evidence, CLIs | Preview | Production apply path still sandbox-only |
| **Karl-Dashboard** | Read-only KHR projection (`readonly-y`) | Preview | Cockpit unchanged; no approval UI workflow |
| **Karl-Inventory** | Runtime posture observation schemas | Preview | No live agent ingest from cluster |
| **Karl-OS-ISO** | Host-runtime preview; systemd disabled | Preview | No default enable; no production enforcement |

---

## Blockers

| ID | Severity | Item |
|----|----------|------|
| P0 | Blocker | No production enable path for karl-host-runtime / guarded apply |
| P0 | Blocker | No autonomous orchestration — operator approval remains local evidence |
| P0 | Blocker | KubeVirt/Multus remain compatibility paths — not removed |
| P1 | Gap | Dashboard approval UX (read-only projection only) |
| P1 | Gap | Inventory ↔ cluster live sync for certification registry |
| P1 | Gap | ISO install still preview/disabled by default |
| P2 | Gap | Windows native-live lane certification parity |
| P2 | Gap | Multi-cluster registry federation |
| P2 | Gap | Beta hardening of provenance trust store |

---

## What is missing by release stage

### Technical preview (next external milestone)

| Required | Status |
|----------|--------|
| Stable read-only APIs + projection contracts | **Mostly done** |
| Documented infrastructure scope (not datacenter-only) | **KHR-Z** |
| Native-live certification evidence on reference cluster | **Done** |
| Policy gates + approval workflow (evidence) | **Done** |
| Operator-facing readiness map | **This doc** |
| Limited apply in named sandbox only | **Done** |
| Published technical preview boundary doc per repo | **KHR-AA** (`TECHNICAL_PREVIEW_READINESS.md` + per-repo TP docs) |

### Beta (future)

| Required | Status |
|----------|--------|
| Dashboard operator approval UI (non-mutating → gated apply) | **Gap** |
| Inventory observation from periodic cluster export | **Gap** |
| ISO optional enable with hardened defaults | **Gap** |
| Certification registry refresh automation | **Gap** |
| Broader lane certification (Windows, hybrid) | **Gap** |
| HA / multi-node KHR coordination | **Gap** |

### Production (explicitly out of scope for KHR program today)

| Required | Status |
|----------|--------|
| Production enable on ISO | **Not started** |
| GA certification claims | **Forbidden** |
| Autonomous apply orchestration | **Forbidden** |
| Zero-trust production guarantees | **Forbidden** |
| Datacenter-only product positioning | **Forbidden** |
| Full KubeVirt decommission | **Not planned near-term** |

---

## Evidence anchors

| Evidence | Path |
|----------|------|
| Native-live certification | `docs/evidence/khr-native-live-lane/certification/` |
| Certification registry | `docs/evidence/khr-certification-registry/` |
| Action approval | `docs/evidence/khr-action-approval/` |
| Control graph | `docs/evidence/khr-control-graph/` |
| Provenance validation | `docs/evidence/khr-provenance/` |

---

## Guardrails (program-wide)

- **No GA claim** for KHR-native-live or shell continuity on ISO.
- **No production-ready claim** for guarded apply or certification registry enforcement.
- **No datacenter-only framing** — use infrastructure operating layer / infrastructure OS.
- **No autonomous orchestration claim** — operator gates and local approval evidence only.
