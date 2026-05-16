# KHR Installer Profile Expectations (KHR-AH / KHR-AI)

| Field | Value |
|-------|-------|
| **Sprint** | KHR-AH / KHR-AI / KHR-AJ / **KHR-AK** |
| **Contract set** | `khr-tp-contract-v1` — `docs/contracts/khr/khr-contract-manifest.yaml` |
| **Installer selector** | `KARL_INSTALLER_PROFILE` or `-profile` (Karl-Installer) |
| **ISO manifest** | `profile-manifest.yaml` — `crdAssetPath`, `expectedCrds` |
| **Production** | **NOT production ready** |

What Karl-Hyperdensity expects when operators select an installer profile after ISO provision. **No new runtime capabilities** in KHR-AI.

---

## Default (unchanged)

| Condition | Installer behavior | Hyperdensity expectation |
|-----------|-------------------|------------------------|
| `KARL_INSTALLER_PROFILE` **unset** | Same as `karl1-kubevirt-legacy` (manifest + virtctl + KubeVirt compatibility + CDI) | VM compatibility path available; KHR CRDs may exist from ISO `install_khr_crds` |
| Explicit `karl1-kubevirt-legacy` | Identical steps to default | Same as default |

---

## `karl2-khr-technical-preview` (KHR-AI)

| Topic | Expectation |
|-------|-------------|
| KHR CRDs | Installer **applies** CRD foundation from `crdAssetPath` assets and **verifies** `expectedCrds` in cluster |
| KubeVirt / CDI / virtctl | **Not** invoked by installer in this profile |
| Host runtime | **Must remain disabled** — installer does not enable systemd |
| Dry-run | `KARL_INSTALLER_KHR_CRD_DRY_RUN=true` — apply dry-run only, skip verify |
| Asset path | `KARL_INSTALLER_KHR_CRD_PATH` or `/opt/karl/karl-engine/khr/crds` or ISO tree path |
| Hyperdensity | Sandbox preflight + evidence on `khr-runtime-sandbox`; native-live does not require KubeVirt |
| Evidence | `khr_crd_foundation_evidence.sh` — includes manifest snapshot + `crd-skew.json` |
| Verify | `contractSetId` + `crdDiffEmpty: true`; `kubevirtCalledInKarl2Profile: false` |
| Skew guard | ISO `profile-manifest.yaml` must match Hyperdensity contract manifest `expectedCrds` |

---

## `hybrid-transition` (KHR-AI)

| Topic | Expectation |
|-------|-------------|
| Order | KHR CRD apply + verify **then** KubeVirt **compatibility provider** + CDI + virtctl |
| KubeVirt | Deployed as **compatibility provider** — not KHR core |
| KHR CRDs | Installer applies same foundation as karl2 before compatibility path |
| Host runtime | **Disabled** |
| Use case | VM workloads + KHR TP CRD foundation in parallel |

---

## Profile vs ISO manifest

| Profile | ISO manifest | Installer KHR-AI |
|---------|--------------|------------------|
| `karl1-kubevirt-legacy` | `defaultInstallerProfile` | Full legacy compile steps (no KHR CRD apply step) |
| `karl2-khr-technical-preview` | `kubevirtCompatibility: excluded` | CRD apply + verify only |
| `hybrid-transition` | KV + CRD required | CRD apply + compatibility provider + CDI |

---

## Operator checklist

```bash
# ISO tree
cd Karl-OS-ISO && ./scripts/khr_iso_tp_verify.sh

# Installer TP (dry-run safe)
export KARL_INSTALLER_PROFILE=karl2-khr-technical-preview
export KARL_INSTALLER_KHR_CRD_DRY_RUN=true
export KARL_INSTALLER_KHR_CRD_PATH=/opt/karl/karl-engine/khr/crds

# Hyperdensity
cd Karl-Hyperdensity && ./scripts/khr_runtime_sandbox_preflight.sh
kubectl get crd hosts.runtime.karl.io resourceleases.hyperdensity.karl.io
```

---

## Beta blockers

| ID | Blocker |
|----|---------|
| B-PROF-AI-01 | Shell `02_deploy-kubevirt.sh` not profile-aware |
| B-PROF-AI-02 | Beta still requires evidence bundle beyond CRD presence |

**Related:** `KHR_BOOTSTRAP_CONSUMER_EXPECTATIONS.md`, `TECHNICAL_PREVIEW_OPERATOR_RUNBOOK.md`.
