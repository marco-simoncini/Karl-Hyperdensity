# KHR Installer Profile Expectations (KHR-AH)

| Field | Value |
|-------|-------|
| **Sprint** | KHR-AH |
| **Installer selector** | `KARL_INSTALLER_PROFILE` or `-profile` (Karl-Installer) |
| **ISO manifest** | `profile-manifest.yaml` / `KHR_BOOTSTRAP_VERSION_MANIFEST.md` |
| **Production** | **NOT production ready** |

What Karl-Hyperdensity expects when operators select an installer profile after ISO provision. **No new runtime capabilities** in KHR-AH.

---

## Default (unchanged)

| Condition | Installer behavior | Hyperdensity expectation |
|-----------|-------------------|------------------------|
| `KARL_INSTALLER_PROFILE` **unset** | Same as `karl1-kubevirt-legacy` (manifest + virtctl + KubeVirt compatibility + CDI) | VM compatibility path available; KHR CRDs from ISO if provision ran `install_khr_crds` |
| Explicit `karl1-kubevirt-legacy` | Identical steps to default | Same as default |

---

## `karl2-khr-technical-preview`

| Topic | Expectation |
|-------|-------------|
| KubeVirt | **Not** deployed as KHR core by installer; optional on cluster from ISO only |
| KHR CRDs | **Required** on cluster via ISO `install_khr_crds` (installer does read-only inventory/log only) |
| Host runtime | **Must remain disabled** — installer does not enable systemd |
| virtctl / CDI | **Not** run by installer in this profile |
| Hyperdensity | Run sandbox preflight + evidence scripts on `khr-runtime-sandbox`; native-live lane does not require KubeVirt |
| Verify | `Karl-OS-ISO/scripts/khr_iso_tp_verify.sh` on ISO tree |

---

## `hybrid-transition`

| Topic | Expectation |
|-------|-------------|
| Order | KHR TP foundation stub (read-only) **then** KubeVirt **compatibility provider** + CDI |
| KubeVirt | Deployed as **compatibility provider** for VM workloads — not KHR core |
| KHR CRDs | ISO foundation required; installer does not apply KHR CRDs |
| Host runtime | **Disabled** — same as other profiles |
| Hyperdensity | Supports both VM compatibility evidence and KHR sandbox lanes (operator-initiated) |
| Use case | Fleets migrating VM workloads while standing up KHR TP CRD foundation |

---

## Profile vs ISO manifest

| Profile | ISO `profile-manifest.yaml` | Installer KHR-AH |
|---------|----------------------------|------------------|
| `karl1-kubevirt-legacy` | `defaultInstallerProfile` | Full legacy compile steps |
| `karl2-khr-technical-preview` | `khrCrdFoundation: required` | Manifest + KHR TP stub only |
| `hybrid-transition` | KV + CRD required | KHR stub + compatibility provider + CDI |

---

## Operator checklist (read-only)

```bash
# ISO tree
cd Karl-OS-ISO && ./scripts/khr_iso_tp_verify.sh

# Installer profile (example TP — does not enable host-runtime)
export KARL_INSTALLER_PROFILE=karl2-khr-technical-preview
# run installer only when cluster context is intentional

# Hyperdensity
cd Karl-Hyperdensity && ./scripts/khr_runtime_sandbox_preflight.sh
```

---

## Beta blockers (profile layer)

| ID | Blocker |
|----|---------|
| B-PROF-AH-01 | Installer KHR stub does not cluster-apply CRDs — ISO provision still required |
| B-PROF-AH-02 | Shell `02_deploy-kubevirt.sh` not profile-aware (Go path only) |
| B-PROF-AH-03 | Beta sign-off still requires evidence bundle, not profile selection alone |

**Related:** `KHR_BOOTSTRAP_CONSUMER_EXPECTATIONS.md`, `TECHNICAL_PREVIEW_OPERATOR_RUNBOOK.md`.
