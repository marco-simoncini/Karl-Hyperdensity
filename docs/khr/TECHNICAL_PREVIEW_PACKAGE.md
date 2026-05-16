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
| **Karl-Dashboard** | Read-only KHR projection documentation | `docs/khr/TECHNICAL_PREVIEW_DASHBOARD_GUIDE.md`, `docs/hyperdensity/KHR_PROJECTION_V1.md` |
| **Karl-Inventory** | Posture/observation schema + guides | `docs/khr/TECHNICAL_PREVIEW_INVENTORY_GUIDE.md`, `docs/contracts/khr/runtime-posture.schema.json` |
| **Karl-OS-ISO** | CRD foundation + host-runtime preview (disabled) | `docs/khr/KHR_TECHNICAL_PREVIEW_ISO_GUIDE.md`, `docs/khr/KHR_TECHNICAL_PREVIEW_PROFILE.md` |

| Validation (Hyperdensity) | Command |
|---------------------------|---------|
| Docs scope guard | `./scripts/guard_khr_docs_scope.sh` |
| TP package check | `./scripts/khr_tp_package_check.sh` |
| TP operator bundle | `./scripts/khr_tp_operator_bundle.sh` |
| Full validate hook | `./scripts/validate.sh` |
| ISO boundaries (ISO repo) | `./scripts/guard_khr_iso_boundaries.sh` |

---

## What the TP package contains

| Category | Contents |
|----------|----------|
| **Documentation** | Readiness audit (AA), package map (AB), per-repo guides, infrastructure scope (Z) |
| **Contracts** | Shell/Cell/ResourcePort/ResourceLease, projection v1alpha1-readonly-y (Dashboard doc) |
| **CRDs** | Host, Shell, Cell, ResourcePort, ResourceLease (install via ISO or Hyperdensity manifests) |
| **CLIs** | `karl-host-runtime`, `khr-cert-registry`, `khr-action-approval`, `khr-control-graph`, `khr-provenance-validate` |
| **Evidence (committed)** | Native-live certification, certification registry, provenance validation summaries |
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
| **Beta-1** | Wire read-only `tpReadinessSummary` in Dashboard API (contract bump) |
| **Beta-2** | Optional Inventory posture export job (non-enforcing) |
| **Beta-3** | ISO TP profile install script (still **disabled** systemd) |
| **Beta-4** | Operator runbook PDF/index bundling all four guides |

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
