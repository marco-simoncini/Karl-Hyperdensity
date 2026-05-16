# KHR Technical Preview Readiness Audit (KHR-AA)

| Field | Value |
|-------|-------|
| **Sprint** | KHR-AA |
| **Audit date** | 2026-05-16 |
| **Program status** | **Technical Preview candidate (sandbox evidence)** |
| **Production** | **NOT production ready** — no production enable path, no GA |

This document is the **single Technical Preview (TP) readiness scorecard** for the KHR program. It supersedes informal readiness notes in `KHR_RELEASE_READINESS_MAP.md` for TP boundary decisions.

---

## Executive scorecard

| Dimension | TP verdict | Notes |
|-----------|------------|-------|
| **Overall TP readiness** | **Candidate** | Read-only APIs + sandbox evidence complete; external TP requires operator acceptance of boundaries below |
| **Production readiness** | **NOT ready** | Explicit: **NOT production ready** |
| **GA / certification claims** | **Forbidden** | `certified-preview` is sandbox evidence only |
| **Autonomous orchestration** | **Absent** | Operator CLI + local approval evidence only |
| **Hidden production enable** | **None found** | ISO systemd disabled; `sandboxApplyEnabled: false` default; production namespace blocklist |
| **Infrastructure framing** | **Aligned (KHR-Z)** | Infrastructure OS / operating layer — not datacenter-only |

---

## Capability audit (implemented / preview / experimental / unsupported)

| Capability | Sprint | Maturity | TP class | Evidence / contract |
|------------|--------|----------|----------|---------------------|
| Host / Shell / Cell / ResourcePort / ResourceLease CRDs | A–E | **Implemented** (API) | Safe for TP | CRDs + `SHELL_CELL_CONTRACT.md` |
| Parent Fabric KHR projection (read-only) | E–Y | **Preview** | Safe for TP | Dashboard `khr-projection-v1alpha1-readonly-y` |
| ResourcePort reporting loop | J | **Preview** | Sandbox only | JSON default; no CR apply on ISO |
| ResourcePort CR preview | K | **Preview** | Sandbox only | `--emit-cr=true`; apply disabled on ISO |
| ResourceLease dry-run vs ResourcePort | L | **Implemented** (read) | Sandbox only | `khr-runtime-sandbox` only |
| ResourceLease guarded apply | M | **Experimental** | Sandbox only | Marker file only in KHR-F; opt-in config |
| karl-host-runtime MVP | F | **Preview** | Sandbox only | ISO unit **disabled** |
| Host heartbeat / runtime session | N | **Preview** | Sandbox only | Observation; no production writes |
| RAM / CPU live scale (cgroup) | O | **Experimental** | Sandbox only | Linux sandbox lane |
| Windows live-scale lane contract | P | **Experimental** | Docs only | Observation contract; no apply |
| Multi-lane discovery | Q | **Preview** | Safe for TP | Read-only classification |
| ResourceFuture simulation | R | **Preview** | Safe for TP | No apply; eligibility only |
| Native-live lane prototype | S | **Experimental** | Sandbox only | `khr-native-live-*` workloads |
| Native-live certification | T | **Preview** | Sandbox only | `certified-preview` — not GA |
| Shell / app / user continuity | U | **Preview** | Sandbox only | Semantics + certification proofs |
| Certification registry | V | **Preview** | Sandbox only | Local registry + policy gate simulation |
| Operator action approval | W | **Preview** | Docs only | Local evidence; no autonomous apply |
| Unified control graph | X | **Preview** | Safe for TP | Export + lineage verify (sandbox) |
| Trust / provenance | Y | **Preview** | Safe for TP | Fingerprints; no production trust enforcement |
| Infrastructure scope language | Z | **Implemented** (docs) | Docs only | `KARL_INFRASTRUCTURE_SCOPE.md` |
| TP readiness audit | AA | **Implemented** (docs) | Docs only | This document |
| KubeVirt provider path | — | **Preview** (compat) | Unsupported as target | `kubevirt.compatibility` |
| Multus / NAD networking | — | **Unsupported** (target) | Not ready | Legacy / transitional only |
| Production namespace mutation | — | **Unsupported** | Not ready | Blocklist enforced in scripts |
| Autonomous apply / orchestration | — | **Unsupported** | Not ready | Forbidden by program guardrails |
| ISO default host-runtime enable | — | **Unsupported** | Not ready | `install_karl_host_runtime` not in `page_install` |
| Dashboard approval apply UI | — | **Unsupported** | Not ready | Projection read-only only |
| Inventory live cluster ingest | — | **Unsupported** | Not ready | Schema + manual observation stubs |

### Maturity legend

| Label | Meaning |
|-------|---------|
| **Implemented** | Contract + code path exists; may still be sandbox-gated |
| **Preview** | Stable read-only or opt-in; not production-enabled |
| **Experimental** | Prototype behavior; evidence required per run |
| **Unsupported** | Out of TP scope or explicitly blocked |

### TP classification legend

| Class | Operator may… |
|-------|----------------|
| **Safe for TP** | Use read-only APIs, docs, simulation, observation exports |
| **Docs only** | Read contracts and runbooks; no cluster mutation |
| **Sandbox only** | Run evidence scripts on `karl-metal-01@ovh` / `khr-runtime-sandbox` with explicit flags |
| **Not ready** | Do not use for production or autonomous workflows |

---

## TP onboarding map

### Cluster prerequisites

| Requirement | Value |
|-------------|-------|
| Kubernetes cluster | Reference: `karl-metal-01@ovh` |
| KHR CRDs | Installed (`install_khr_crds` on ISO; or Hyperdensity manifests) |
| KubeVirt | Present as **compatibility provider** (not removed for TP) |
| Hyperdensity repo | Built `karl-host-runtime`, CLIs, evidence scripts |
| Operator access | `kubectl` context matching `KHR_RUNTIME_CLUSTER_CONTEXT` |

### Required namespaces

| Namespace | Role | TP use |
|-----------|------|--------|
| `khr-runtime-sandbox` | **Required** for all mutation/evidence paths | Allowlisted sandbox |
| `karl-system` | Platform | **Blocked** for KHR apply |
| `kube-system` | Platform | **Blocked** |
| `default` | Legacy | **Blocked** for KHR apply |

### Required labels

| Label | Value | Required on |
|-------|-------|-------------|
| `khr.karl.io/sandbox=true` | `true` | Namespace + workloads in sandbox paths |
| `khr.karl.io/native-live=true` | `true` | Native-live lane workloads (optional lane) |
| `karl.io/sandbox-namespace` | `<ns>` | Cluster-scoped ResourcePort CRs (preview) |

### Workload naming (native-live)

| Pattern | Lane |
|---------|------|
| Prefix `khr-native-live-` in `khr-runtime-sandbox` | `native-live` |

### Supported lanes (TP)

| Lane | Classification | Provider | TP notes |
|------|----------------|----------|----------|
| `native-live` | `native-live` | `khr.native` | Linux container cgroup; certified-preview in sandbox |
| `live-in-place-capable` | Simulation | varies | ResourceFuture may show eligible; apply still sandbox |
| `compatibility-fallback` | `compatibility-fallback` | `kubevirt.compatibility` | Restart/rollout may be required |
| `observation-only` | `observation-only` | varies | Read-only discovery |

### Unsupported lanes (TP)

| Lane / path | Reason |
|-------------|--------|
| Production namespaces | Blocklist; `productionUnsupported` |
| Windows native-live certification parity | Experimental observation only (KHR-P) |
| Autonomous ResourceLease reconcile | Not implemented |
| ISO-shipped default `karl-host-runtime` systemd | Disabled; manual enable is operator risk |
| Multus as target network fabric | Unsupported architecture |
| Unqualified production cgroup apply | Blocked by `sandboxApplyEnabled: false` |

### Operator bootstrap (sandbox)

```bash
export KHR_RUNTIME_CLUSTER_CONTEXT=karl-metal-01@ovh
export KHR_RUNTIME_NAMESPACE=khr-runtime-sandbox
# Preflight + evidence (see KHR_RUNTIME_SANDBOX_EXECUTION.md)
KHR_RUNTIME_SANDBOX_LIVE=1 ./scripts/khr_runtime_sandbox_execute.sh
```

---

## Verification checklist (KHR-AA)

| Check | Result |
|-------|--------|
| No hidden production enable in ISO default provision | **PASS** — `install_karl_host_runtime` not in `page_install` |
| `systemctl disable karl-host-runtime` in install stub | **PASS** (ISO guard) |
| Default `sandboxApplyEnabled: false` | **PASS** |
| No autonomous orchestration product path | **PASS** |
| No unqualified GA wording in KHR docs guard | **PASS** (`guard_khr_docs_scope.sh`) |
| ISO boundary guard (GA / autonomous / datacenter-only) | **PASS** (`guard_khr_iso_boundaries.sh`) |

---

## Blocker list

| ID | Severity | Blocker | TP impact |
|----|----------|---------|-----------|
| B-01 | **P0** | **NOT production ready** — no production apply | External TP must not claim production |
| B-02 | **P0** | No autonomous orchestration / apply on approval | Approval is evidence-only |
| B-03 | **P0** | ISO host-runtime disabled by default | TP uses Hyperdensity sandbox, not ISO enable |
| B-04 | **P1** | Dashboard: no TP approval UI (projection doc only) | Operators use CLI + docs |
| B-05 | **P1** | Inventory: no live cluster observation feed | Manual / file-based posture only |
| B-06 | **P1** | Windows lane certification gap | Linux native-live only for TP cert |
| B-07 | **P2** | Multi-cluster registry federation | Single-cluster TP |
| B-08 | **P2** | Provenance trust store hardening | Sandbox fingerprints sufficient for TP |

---

## Recommended next milestones

| Milestone | Target | Depends on |
|-----------|--------|------------|
| **TP-1** | Publish per-repo TP boundary docs (Dashboard, Inventory, ISO) | KHR-AA |
| **TP-2** | Operator runbook bundle (single PDF/README index) | TP-1 |
| **TP-3** | Optional ISO TP profile install (still **disabled** systemd) | ISO guard unchanged |
| **Beta-1** | Dashboard read-only TP summary in API (future contract bump) | No runtime in AA |
| **Beta-2** | Inventory periodic posture export job | Agent design |
| **Beta-3** | Gated apply UX (non-autonomous) | Approval workflow hardening |

---

## Cross-repo deliverables (KHR-AA)

| Repo | Document |
|------|----------|
| Karl-Hyperdensity | This file (canonical scorecard) |
| Karl-Dashboard | `docs/khr/TECHNICAL_PREVIEW_READINESS_SUMMARY.md` |
| Karl-Inventory | `docs/khr/TECHNICAL_PREVIEW_READINESS_OBSERVATION.md` |
| Karl-OS-ISO | `docs/khr/KHR_TECHNICAL_PREVIEW_PROFILE.md` |

---

## Explicit statement

**KARL KHR Technical Preview is NOT production ready.**

Technical Preview permits **read-only observation**, **documentation**, and **named-sandbox evidence** only. Any production enable, autonomous orchestration, or GA certification claim is **out of scope** and **unsupported** for this program phase.

---

## Related

- `docs/khr/KHR_RELEASE_READINESS_MAP.md` — program roadmap and beta/production gaps
- `docs/khr/KARL_INFRASTRUCTURE_SCOPE.md` — product positioning
- `scripts/guard_khr_docs_scope.sh` — docs language guard
