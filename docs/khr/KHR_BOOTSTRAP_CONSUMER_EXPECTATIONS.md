# KHR Bootstrap Consumer Expectations (KHR-AG)

| Field | Value |
|-------|-------|
| **Sprint** | KHR-AG |
| **Consumer** | Karl-Hyperdensity CLIs, evidence scripts, sandbox lanes |
| **Cluster (reference)** | `karl-metal-01@ovh` |
| **Production** | **NOT production ready** |

Defines what Hyperdensity **expects** from ISO / installer / subiquity bootstrap before operators run Technical Preview evidence — **no new capabilities** in KHR-AG.

---

## Bootstrap inputs (upstream)

| Source | Delivers | Hyperdensity assumes |
|--------|----------|----------------------|
| **Karl-OS-ISO** | `install_khr_crds`, disabled `karl-host-runtime`, boundary docs | CRD APIs exist after provision apply |
| **Karl-OS-ISO_subiquity** | Host OS only | No KHR runtime enable at autoinstall |
| **Karl-Installer** | Profile-dependent KubeVirt/CDI (legacy default today) | VM path optional for TP; not required for native-live lane |
| **Operator** | Sandbox NS + labels + manual script invocation | `khr-runtime-sandbox`, `khr.karl.io/sandbox=true` |

---

## Required CRDs (API foundation)

After ISO provision (or equivalent manual apply from Hyperdensity manifests):

| CRD | API group | Used by |
|-----|-----------|---------|
| Host | `khr.karl.io` | Runtime discovery / posture |
| Shell, Cell, ShellClass | `khr.karl.io` | Shell factory / sandbox fixtures |
| ResourcePort | `khr.karl.io` | Native-live / port freeze evidence |
| ShellLease | `khr.karl.io` | Lease-adjacent sandbox tests |
| GatewayRoute | `khr.karl.io` | Gateway access foundation (read-only TP) |
| ResourceLease | `hyperdensity.karl.io` | ResourceLease freeze / observation |

**Not required for TP tree-verify:** controllers, webhooks, or `karl-host-runtime` daemon running.

---

## Required labels and namespaces

| Resource | Requirement |
|----------|-------------|
| Namespace | `khr-runtime-sandbox` (create via sandbox scripts if missing) |
| Namespace label | `khr.karl.io/sandbox=true` |
| Workloads in sandbox | `khr.karl.io/sandbox=true` on objects touched by evidence scripts |
| Context | `KHR_RUNTIME_CLUSTER_CONTEXT=karl-metal-01@ovh` (reference) |

Installer **does not** create sandbox namespace — Hyperdensity preflight / execute scripts do.

---

## TP evidence prerequisites

| Prerequisite | Verification |
|--------------|--------------|
| ISO CRD foundation applied on cluster | `kubectl get crd hosts.khr.karl.io` (operator OOB) |
| ISO tree contract | `Karl-OS-ISO/scripts/khr_iso_tp_verify.sh` |
| ISO boundaries | `Karl-OS-ISO/scripts/guard_khr_iso_boundaries.sh` |
| Runtime **not** enabled on ISO path | No `install_karl_host_runtime` in `page_install` |
| Hyperdensity repo validate | `./scripts/validate.sh` |
| Committed evidence anchors | `docs/evidence/khr-native-live-lane/`, registry, provenance dirs per runbook |

---

## Supported TP lanes (consumer)

| Lane | Bootstrap dependency | Notes |
|------|---------------------|-------|
| **Native-live** | ResourcePort CRD + sandbox NS | Certification in `docs/evidence/khr-native-live-lane/` |
| **ResourceLease** | ResourceLease CRD | Freeze candidate docs; apply experimental |
| **ResourceFuture / lane discovery** | CRD + sandbox labels | Evidence scripts only |
| **Certification registry / provenance** | Evidence files + CLIs | No production admission |
| **KubeVirt VM workloads** | KubeVirt compatibility path (ISO/installer) | **Optional** for TP; required for `karl1-kubevirt-legacy` |

**Not supported as default bootstrap:** production enable, autonomous orchestration, Dashboard rewrite, runtime mutation on ISO install path.

---

## Read-only verify flow (operator)

```bash
# ISO repo (tree)
cd Karl-OS-ISO && ./scripts/khr_iso_tp_verify.sh && ./scripts/guard_khr_iso_boundaries.sh

# Hyperdensity
cd Karl-Hyperdensity && ./scripts/guard_khr_docs_scope.sh && ./scripts/validate.sh
export KHR_RUNTIME_NAMESPACE=khr-runtime-sandbox
./scripts/khr_runtime_sandbox_preflight.sh
```

---

## Beta bootstrap blockers (Hyperdensity consumer)

| ID | Blocker |
|----|---------|
| B-HD-01 | Cluster CRD presence not automated in ISO CI verify |
| B-HD-02 | Installer always deploys KubeVirt — confuses TP “optional KV” story |
| B-HD-03 | No unified bootstrap version manifest across four repos |
| B-HD-04 | Subiquity autoinstall does not surface TP profile selection |
| B-HD-05 | Beta sign-off requires operator attestation + evidence bundle, not bootstrap alone |

---

## Cross-repo map

| Document | Repo |
|----------|------|
| `KHR_BOOTSTRAP_FLOW.md` | Karl-OS-ISO |
| `KHR_INSTALLER_PROFILE_MATRIX.md` | Karl-Installer |
| `KHR_SUBIQUITY_ALIGNMENT.md` | Karl-OS-ISO_subiquity |
| `TECHNICAL_PREVIEW_OPERATOR_RUNBOOK.md` | Karl-Hyperdensity (this repo) |
