# KHR Technical Preview — Post-Install Verification (KHR-AL)

| Field | Value |
|-------|-------|
| **Sprint** | KHR-AL |
| **Production** | **NOT production ready** |
| **Cluster** | `karl-metal-01@ovh` |

Single operator view for **post-install TP verification** aggregating installer evidence, ISO verify, contract manifest, and operator bundle checks.

---

## Verification stack

| Layer | Script / doc | Mode |
|-------|--------------|------|
| Contract manifest | `scripts/khr_contract_manifest_check.sh` | Read-only |
| ISO tree + cluster | `Karl-OS-ISO/scripts/khr_post_install_verify.sh` | Read-only |
| Installer karl2 | `Karl-Installer/scripts/khr_crd_foundation_evidence.sh` | CRD apply gated |
| Installer hybrid | `Karl-Installer/scripts/khr_hybrid_transition_evidence.sh` | CRD apply + KV attestation |
| TP operator bundle | `scripts/khr_tp_operator_bundle.sh` | Read-only |
| **Aggregate** | `scripts/khr_tp_post_install_bundle_check.sh` | Read-only |

---

## Run aggregate check

```bash
cd Karl-Hyperdensity
./scripts/khr_tp_post_install_bundle_check.sh
```

---

## Expected evidence fields (installer hybrid)

| Field | Expected |
|-------|----------|
| `contractSetId` | `khr-tp-contract-v1` |
| `crdDiffEmpty` | `true` (after apply) |
| `kubevirtCompatibility` | `true` |
| `kubevirtAsKhrCore` | `false` |
| `hostRuntimeEnabled` | `false` |
| `noAutonomousOrchestration` | `true` |

---

## Profiles (unchanged defaults)

| Profile | Default when unset |
|---------|-------------------|
| *(empty)* | `karl1-kubevirt-legacy` |
| `karl2-khr-technical-preview` | CRD only — no KV/CDI/virtctl |
| `hybrid-transition` | CRD + compatibility path (KV not KHR core) |

Verification layer does not enable runtime, production paths, or orchestration without operator intent.
