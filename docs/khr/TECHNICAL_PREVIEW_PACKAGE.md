# KHR Technical Preview Package (KHR-AB)

| Field | Value |
|-------|-------|
| **Sprint** | KHR-AB |
| **Package ID** | `khr-technical-preview-v1` |
| **Status** | **Consumable by operators/dev** — sandbox/manual only |
| **Production** | **NOT production ready** |

This document defines the **Technical Preview (TP) package** consumable without production enable, autonomous orchestration, or systemd enable on ISO.

---

## Cross-repo TP package map

| Repo | Role in TP package | Primary artifacts |
|------|-------------------|-------------------|
| **Karl-Hyperdensity** | Source of truth: CLIs, CRDs, evidence, sandbox scripts | `TECHNICAL_PREVIEW_READINESS.md`, `TECHNICAL_PREVIEW_PACKAGE.md`, `TECHNICAL_PREVIEW_OPERATOR_RUNBOOK.md`, `docs/evidence/khr-*`, `scripts/khr_*` |
| **Karl-Dashboard** | Read-only KHR projection documentation | `docs/khr/TECHNICAL_PREVIEW_DASHBOARD_GUIDE.md`, `DASHBOARD_REFERENCE_ENV_ACTIVATION_PROFILE.md` (KHR-BN), `docs/hyperdensity/KHR_PROJECTION_V1.md`, `DASHBOARD_BACKEND_KHR_MIGRATION_PLAN.md` (consumer) |
| **Hyperdensity** | Dashboard backend compatibility model | `DASHBOARD_BACKEND_COMPATIBILITY_MODEL.md` (KHR-BH) |
| **Karl-Inventory** | Posture/observation schema + guides | `docs/khr/TECHNICAL_PREVIEW_INVENTORY_GUIDE.md`, `docs/contracts/khr/runtime-posture.schema.json` |
| **Karl-OS-ISO** | CRD foundation + host-runtime preview (disabled) | `docs/khr/KHR_TECHNICAL_PREVIEW_ISO_GUIDE.md`, `docs/khr/KHR_TECHNICAL_PREVIEW_PROFILE.md`, `docs/khr/KHR_BOOTSTRAP_FLOW.md` |
| **Karl-OS-ISO_subiquity** | Host install UI — TP wording alignment only | `docs/khr/KHR_SUBIQUITY_ALIGNMENT.md` |
| **Karl-Installer** | Profile selector + CRD foundation apply (AI) | `KARL_INSTALLER_PROFILE`, `KARL_INSTALLER_KHR_CRD_PATH`, karl2 applies `expectedCrds` |
| **rdp-GW** | Read-only ShellLease / GatewayRoute compatibility (KHR-AM) | `docs/khr/RDPGW_KHR_ALIGNMENT_PLAN.md`, `GET /karl-gw/v1/shell/resolve`, `GET /karl-gw/v1/gatewayroute/resolve` |

| Validation (Hyperdensity) | Command |
|---------------------------|---------|
| Docs scope guard | `./scripts/guard_khr_docs_scope.sh` |
| TP package check | `./scripts/khr_tp_package_check.sh` |
| TP operator bundle | `./scripts/khr_tp_operator_bundle.sh` |
| Full validate hook | `./scripts/validate.sh` |
| Scope-4 certification check (read-only) | `./scripts/khr_scope4_certification_check.sh` |
| Scope-4 governance bundle (read-only) | `./scripts/khr_scope4_governance_bundle.sh` |
| ISO boundaries (ISO repo) | `./scripts/guard_khr_iso_boundaries.sh` |
| ISO TP bootstrap verify (ISO repo, read-only) | `./scripts/khr_iso_tp_verify.sh` |
| Installer karl2 CRD evidence (optional) | `khr_crd_foundation_evidence.sh` → `contractSetId` + `crdDiffEmpty` |
| Canonical contract manifest | `docs/contracts/khr/khr-contract-manifest.yaml` (`khr-tp-contract-v1`) |
| Contract manifest check | `./scripts/khr_contract_manifest_check.sh` |
| TP post-install bundle | `./scripts/khr_tp_post_install_bundle_check.sh` |
| ISO post-install verify | `Karl-OS-ISO/scripts/khr_post_install_verify.sh` |
| Hybrid transition evidence | `Karl-Installer/scripts/khr_hybrid_transition_evidence.sh` |
| Bootstrap consumer contract (this repo) | `docs/khr/KHR_BOOTSTRAP_CONSUMER_EXPECTATIONS.md` |

---

## What the TP package contains

| Category | Contents |
|----------|----------|
| **Documentation** | Readiness audit (AA), package map (AB), per-repo guides, infrastructure scope (Z) |
| **Contracts** | Shell/Cell/ResourcePort/ResourceLease, projection v1alpha1-readonly-y (Dashboard doc) |
| **CRDs** | Host, Shell, Cell, ResourcePort, ResourceLease (install via ISO or Hyperdensity manifests) |
| **CLIs** | `karl-host-runtime`, `khr-cert-registry`, `khr-action-approval`, `khr-control-graph`, `khr-provenance-validate` |
| **Evidence (committed)** | Native-live certification, certification registry, provenance validation summaries, **Scope-4 guarded-apply certification (KHR-BF)**, **Scope-4 operational governance bundle (KHR-BG)** |
| **Sandbox scripts** | `khr_runtime_sandbox_*.sh`, lane/cert/registry/provenance evidence scripts |
| **Examples** | `examples/khr/runtime-sandbox/` manifests and fixtures |

---

## What the TP package does NOT contain

| Excluded | Reason |
|----------|--------|
| Production enable path | **NOT production ready** |
| Autonomous orchestration / apply on approval | Forbidden |
| ISO default `karl-host-runtime` systemd enable | Runtime disabled by default |
| Dashboard mutating action buttons | Cockpit unchanged; projection read-only |
| Inventory enforcement agent (required) | Observation only |
| GA / production-ready certification claims | Forbidden |
| Multus as target architecture | Legacy/transitional only |
| Windows native-live TP certification parity | Experimental observation only |

---

## Cluster prerequisites

| Requirement | Value |
|-------------|-------|
| Cluster context | `karl-metal-01@ovh` (reference) |
| Reference KHR-native activation | Dashboard env: `HYPERDENSITY_KHR_PROVIDER_PROFILE=khr-native` + reference flags — **not** global default (KHR-BN) |
| Legacy / default profile | `public-cloud-kubevirt-compatibility` when activation env unset |
| KHR CRDs | Installed |
| Hyperdensity binaries | Built from KHR branch |
| Operator | `kubectl` + bash; explicit sandbox flags |

---

## Namespace and labels

| Item | Value |
|------|-------|
| Sandbox namespace | `khr-runtime-sandbox` |
| Required label | `khr.karl.io/sandbox=true` |
| Native-live label | `khr.karl.io/native-live=true` (lane workloads) |
| ResourcePort label | `karl.io/sandbox-namespace=<ns>` |
| Blocked namespaces | `karl-system`, `kube-system`, `default`, … |

---

## Supported lanes (TP)

| Lane | Classification | Notes |
|------|----------------|-------|
| `native-live` | `native-live` | Linux cgroup; sandbox certification evidence |
| `live-in-place-capable` | Simulation | ResourceFuture eligibility |
| `compatibility-fallback` | `kubevirt.compatibility` | May require restart/rollout |
| `observation-only` | Read-only discovery | No apply |

---

## Unsupported lanes (TP)

| Path | Reason |
|------|--------|
| Production namespaces | Blocklist |
| Autonomous ResourceLease reconcile | Not implemented |
| ISO auto-enabled host-runtime | Disabled by default |
| Windows native-live cert parity | Gap (P2) |
| Unqualified production cgroup apply | `sandboxApplyEnabled: false` default |

---

## Required evidence (committed anchors)

| Evidence | Path | Purpose |
|----------|------|---------|
| Native-live certification | `docs/evidence/khr-native-live-lane/certification-summary.json` | TP native-live anchor |
| Certification registry | `docs/evidence/khr-certification-registry/summary.json` | Registry + policy gate evidence |
| Provenance validation | `docs/evidence/khr-provenance/summary.json` | Trust/lineage evidence |

Optional live regeneration: `scripts/khr_native_live_certification.sh`, `khr_cert_registry_policy_gates.sh`, `khr_provenance_evidence.sh` (sandbox only).

---

## Rollback / provenance requirements

| Requirement | TP rule |
|-------------|---------|
| **Rollback** | Every sandbox apply path must run `khr_runtime_sandbox_rollback.sh` or `resourcelease-rollback` in evidence pipelines |
| **Provenance** | Registry and approval bundles carry fingerprints; `khr-provenance-validate` must PASS before citing evidence |
| **Mutation check** | Evidence runs include `mutation-check.txt` proving no production namespace writes |
| **Flight recorder** | Required during sandbox dry-run (`-mode=flight-recorder`) |

---

## Validation commands

```bash
# Hyperdensity (required before citing TP package)
./scripts/guard_khr_docs_scope.sh
./scripts/khr_tp_package_check.sh
./scripts/validate.sh

# ISO repo (boundary guard)
cd ../Karl-OS-ISO && ./scripts/guard_khr_iso_boundaries.sh
```

Live sandbox (optional, operator cluster):

```bash
export KHR_RUNTIME_CLUSTER_CONTEXT=karl-metal-01@ovh
KHR_RUNTIME_SANDBOX_LIVE=1 ./scripts/khr_runtime_sandbox_execute.sh
```

---

## rdp-GW gateway consumer (KHR-AM)

Wave 2 alignment starts in **rdp-GW** without changing Hyperdensity CRDs or controllers:

| Item | KHR-AM behavior |
|------|-----------------|
| Contract source | This repo: `SHELLLEASE_GATEWAYROUTE_CONTRACT.md`, CRDs, JSON schemas |
| rdp-GW resolvers | Stub read-only; `contractSetId: khr-tp-contract-v1` |
| Legacy poolId | Compatibility mapping only; not the long-term routing key |
| Dashboard | Doc projection only (`DASHBOARD_GATEWAYROUTE_RDPGW_ALIGNMENT.md`) — no new UI |

Validation: rdp-GW `go test ./cmd/rdpgw/...` including `cmd/rdpgw/khr`.

### Access graph continuity evidence (KHR-AT)

| Trust | `source` | When |
|-------|----------|------|
| **live-readonly** (preferred) | `live-readonly` | Sandbox rdp-GW reachable (`RDP_GW_BASE_URL`) |
| **fixture-readonly** | `fixture-readonly` | CI/offline; golden fallback |

Hyperdensity: `./scripts/khr_access_graph_continuity_bundle_check.sh` — accepts both; ranks live above fixture.  
Docs: `ACCESS_GRAPH_CONTINUITY_EVIDENCE.md`, rdp-GW `RDPGW_SANDBOX_LIVE_EVIDENCE.md`.

### Dashboard reference activation evidence (KHR-BO / KHR-BP)

| Level | When |
|-------|------|
| **live-readonly** (`LIVE_PASS`) | KHR console image deployed + activation env on `karl-console-next-oidc` + live `GET .../khr-backend/projection` |
| **fixture-readonly** | CI/offline fixture script |
| **remediation-readonly** (`REMEDIATION_PASS`) | Live port-forward OK but env/route gap — operator plan, no auto-patch |

Script: `khr_dashboard_reference_env_live_evidence.sh` (port-forward documented).  
Artifact: `docs/evidence/khr-dashboard-reference-env-activation/<runId>/summary.json`.

Reference env checklist: `KHR_TP_LIVE_REFERENCE_ENVIRONMENT.md`.

---

## Residual blockers (TP package)

| ID | Severity | Blocker |
|----|----------|---------|
| AB-01 | P0 | **NOT production ready** — package is sandbox/manual only |
| AB-02 | P0 | No autonomous orchestration in package |
| AB-03 | P1 | Dashboard `tpReadinessSummary` not in API (doc only) |
| AB-04 | P1 | Inventory live ingest not in package |
| AB-05 | P2 | Multi-cluster evidence federation |

---

## Next steps toward beta

| Step | Description |
|------|-------------|
| **Beta-1** | Contract freeze sign-off (`KHR_CONTRACT_FREEZE_PLAN.md`) |
| **Beta-2** | TP readiness enabled in reference deployments (flag) |
| **Beta-3** | Optional Inventory posture export job (non-enforcing) |
| **Beta-4** | ISO TP profile verify script (still **disabled** systemd) |
| **Beta-5** | Operator runbook index across four repos |

---

## Related docs

| Doc | Repo |
|-----|------|
| `TECHNICAL_PREVIEW_READINESS.md` | Hyperdensity (scorecard) |
| `TECHNICAL_PREVIEW_DASHBOARD_GUIDE.md` | Dashboard |
| `TECHNICAL_PREVIEW_INVENTORY_GUIDE.md` | Inventory |
| `TECHNICAL_PREVIEW_ISO_GUIDE.md` | ISO |

---

## Explicit statement

**This TP package is NOT production ready.** Consumption is limited to documentation, read-only observation, and **named-sandbox manual evidence** on `khr-runtime-sandbox`. No systemd enable, no production enable, no autonomous orchestration.
