# KHR Beta Readiness Gap Analysis (KHR-AD)

| Field | Value |
|-------|-------|
| **Sprint** | KHR-AD |
| **Scope** | Documentation / contract analysis only |
| **Production** | **NOT production ready** — explicitly out of scope for beta program phase |

---

## Executive summary

Technical Preview (TP) packages read-only observation, sandbox evidence, and documented contracts. **Beta** requires contract freeze, broader observation ingest, and operator ergonomics — **without** production enablement, autonomous orchestration, or ISO systemd auto-enable.

| Dimension | TP (now) | Beta target | Production |
|-----------|----------|-------------|------------|
| Read-only APIs | Yes | Hardened + stable | Out of scope |
| Sandbox apply | Opt-in manual | Still manual only | Forbidden |
| Autonomous orchestration | Absent | **Still forbidden** | Forbidden |
| ISO host-runtime | Disabled | Still disabled default | Forbidden |
| GA claims | Forbidden | Forbidden | Forbidden |

---

## Contract inventory (summary)

Full machine-readable inventory: `./scripts/khr_contract_inventory.sh` → `docs/evidence/khr-contract-inventory/summary.json`

| Contract | Version | Stability | Tests | Evidence |
|----------|---------|-----------|-------|----------|
| KHR projection | `khr-projection-v1alpha1-readonly-y` | **Freeze TP** | Dashboard golden | N/A |
| TP readiness | `khr-tp-readiness-summary-v1alpha1` | **Freeze TP** | Dashboard golden | N/A |
| Host | `runtime.karl.io/v1alpha1` | Freeze TP | `pkg/khr/host/*` | heartbeat/registration |
| Shell / Cell | `runtime.karl.io/v1alpha1` | Freeze TP | crdv1alpha1, projection | runtime-sandbox |
| ResourcePort | `runtime.karl.io/v1alpha1` | **TP freeze candidate (KHR-AF)** | resourceport/* | loop, cr-preview, native-live cr-preview |
| Native-live lane | `native-live` / `khr.native` | **TP freeze candidate (KHR-AF)** | nativelive/* | certification-summary, continuity 1.0 |
| ResourceLease | `hyperdensity.karl.io/v1alpha1` | **TP freeze candidate (KHR-AE)** | resourcelease/* | dryrun, guarded-apply, rollback |
| ResourceLease sandbox apply | evidence scripts | **Experimental** | guarded_apply_sandbox | cgroup apply only |
| ResourceFuture | `hyperdensity.karl.io/v1alpha1` | Freeze TP | resourcefuture/* | resourcefuture |
| ShellLease / GatewayRoute | v1alpha1 | Freeze TP | projection | partial |
| Certification registry | `khr-cert-registry-v1` | **Experimental** | certregistry/* | summary.json |
| Policy gates | `khr-policy-gates-v1` | **Experimental** | policygates/* | cert registry runs |
| Action approval | `khr-action-approval-v1` | **Experimental** | actionapproval/* | summary.json |
| Control graph | `khr-control-graph-v1` | **Experimental** | controlgraph/* | summary.json |
| Provenance | `khr-provenance-v1` | **Experimental** | provenance/* | summary.json |

---

## Schema stability levels

| Level | Meaning | Examples |
|-------|---------|----------|
| **Frozen (TP)** | Additive JSON only; no field removals | Projection readonly-y, CRD v1alpha1 core types |
| **TP freeze candidate** | Contract/schema/projection frozen; apply execution sandbox-only | ResourceLease (KHR-AE), ResourcePort + native-live (KHR-AF) |
| **Experimental** | May change with evidence revision | Lease cgroup apply execution, cert registry enforcement semantics |
| **Beta candidate** | Freeze after beta-1 sign-off | TP observation export (Inventory), tp-readiness API defaults |
| **Deprecated** | Compatibility only | KubeVirt path, Multus/NAD |

---

## Test and evidence coverage

| Area | Unit tests | Evidence scripts | Gap |
|------|------------|-------------------|-----|
| CRD parse | Yes | N/A | — |
| Native-live cert | Yes | `khr_native_live_certify.sh` | Windows parity |
| Registry / gates | Yes | `khr_cert_registry_policy_gates.sh` | No scheduled refresh |
| Approval | Yes | `khr_action_approval_evidence.sh` | No Dashboard apply |
| Control graph | Yes | `khr_control_graph_evidence.sh` | No multi-cluster |
| Provenance | Yes | `khr_provenance_evidence.sh` | Trust store hardening |
| Operator bundle | N/A | `khr_tp_operator_bundle.sh` | — |
| Contract inventory | N/A | `khr_contract_inventory.sh` | — |

---

## Beta blockers

### P0 (program boundaries — never auto-unblock)

| ID | Blocker |
|----|---------|
| B-P0-01 | **NOT production ready** |
| B-P0-02 | **No autonomous orchestration** |
| B-P0-03 | **No production enable** / namespace mutation |
| B-P0-04 | **No ISO systemd auto-enable** |

### P1 (beta delivery gaps)

| ID | Blocker |
|----|---------|
| B-P1-01 | Inventory live posture ingest (file or agent) |
| B-P1-02 | Certification registry refresh automation (read-only export) |
| B-P1-03 | TP readiness API flag still off by default (documented) |
| B-P1-04 | Windows native-live certification parity |
| B-P1-05 | Contract bump process (`v1beta1`) not executed |

### P2 (hardening)

| ID | Blocker |
|----|---------|
| B-P2-01 | Multi-cluster evidence federation |
| B-P2-02 | Provenance trust store HA |
| B-P2-03 | ISO optional post-install read-only health script |

---

## Automation explicitly forbidden (beta)

| Automation | Status |
|------------|--------|
| Autonomous ResourceLease reconcile | **Forbidden** |
| Auto-apply on operator approval | **Forbidden** |
| Production namespace controllers | **Forbidden** |
| systemd enable on ISO install | **Forbidden** |
| Dashboard mutating action buttons | **Forbidden** |
| Self-healing remediation from Inventory | **Forbidden** |

Allowed: scheduled **read-only** exports, evidence regeneration in sandbox, guard scripts in CI.

---

## Prioritized beta backlog

| Priority | Item | Repo |
|----------|------|------|
| 1 | Contract freeze sign-off (`KHR_CONTRACT_FREEZE_PLAN.md`) | Hyperdensity |
| 2 | Inventory periodic posture export (read-only) | Inventory |
| 3 | Enable tp-readiness in reference deployments (flag) | Dashboard |
| 4 | Certification registry export job (no apply) | Hyperdensity |
| 5 | ISO post-install read-only verify script | ISO |
| 6 | Windows lane certification evidence | Hyperdensity |

---

## Non-goals (explicit)

- **Production ready** / GA certification
- **Autonomous orchestration**
- **systemd auto-enable** on ISO
- **Dashboard frontend rewrite**
- **Mutating action buttons** on approval/provenance surfaces
- Removing KubeVirt compatibility path in beta

---

## Related

- `KHR_CONTRACT_FREEZE_PLAN.md`
- `TECHNICAL_PREVIEW_READINESS.md`
- `scripts/khr_contract_inventory.sh`
