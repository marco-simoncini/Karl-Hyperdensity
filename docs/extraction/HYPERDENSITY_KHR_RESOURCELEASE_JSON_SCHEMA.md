# Hyperdensity / KHR — ResourceLease JSON Schema (Sprint 91)

## Summary

**Sprint 91** publishes a **non-applied** JSON Schema and three example fixtures for KHR ResourceLease, anchored to Sprint 88–90 work. No CRD, no controller, no runtime.

---

## 1. Scope

| In scope | Out of scope |
|----------|--------------|
| `resourcelease.schema.json` | Kubernetes CRD apply |
| Schema manifest + 3 JSON examples | Controller reconciliation |
| `scripts/validate_resourcelease_schema.sh` | API / Parent Fabric runtime changes |
| Parity goldens in Dashboard | Broad observation / adapter drift |

---

## 2. Non-goals

- Deployable `CustomResourceDefinition`
- Cluster apply of examples
- JSON Schema validation library dependency (stdlib structural checks only)
- Copying `vm/karl_instances/5.winsrv_vm.yaml` into repo (documental citation only)

---

## 3. Schema location

`docs/contracts/khr/resourcelease.schema.json`

- `$schema`: draft 2020-12
- `$id`: `https://karl.local/schemas/khr/resourcelease.schema.json`
- Top-level required: `apiVersion`, `kind`, `metadata`, `spec`
- Optional: `status`

---

## 4. Manifest location

`docs/contracts/khr/resourcelease.schema.manifest.json`

Declares `schemaOnly: true`, `crdApplied: false`, anchors to Sprint 89 contract docs.

---

## 5. Example fixtures

| File | Purpose |
|------|---------|
| `examples/resourcelease-windows-daas-khr-native.json` | `khr.native`, goldenImage OS, persistent profile, scratch |
| `examples/resourcelease-public-cloud-kubevirt-fallback.json` | `kubevirt.public-cloud-fallback`, pvc `karl-os-nfs`, `kubevirt.legacy.ovn` |
| `examples/resourcelease-baremetal-native.json` | `khr.native`, `ephemeralClone`, `baremetal.vlan` |

---

## 6. Validation method

```bash
./scripts/validate_resourcelease_schema.sh
```

Integrated into `./scripts/validate.sh`. Uses Python 3 stdlib: JSON parse, enum extraction from `$defs`, structural checks on examples.

---

## 7. Storage coverage

Disk modes: `ephemeralOverlay`, `ephemeralClone`, `scratch`, `readonly`, `persistent`.

Source types: `pvc`, `image`, `snapshot`, `volume`, `goldenImage`.

Discard policies: `deleteOnStop`, `keepOnFailure`, `promoteOnRequest`.

---

## 8. Network coverage

`attachments[]`, `policies`, `exposure`, `providerNetwork`, `networkLease`.

Attachment required: `name`, `networkRef`, `role`.

---

## 9. Provider coverage

Enum: `khr.native`, `kubevirt.compatibility`, `kubevirt.public-cloud-fallback`, OVN/CNI/cloud/baremetal variants.

---

## 10. KubeVirt compatibility

Public-cloud example: `ephemeralOverlay` + `source.type=pvc` + `ref=karl-os-nfs` maps Sprint 90 ephemeral PVC inventory.

---

## 11. OVN/SDN compatibility

Public-cloud example: `providerNetwork.provider=kubevirt.legacy.ovn`, `providerBinding=ovn.logicalPort`.

---

## 12. Risks

- Treating schema file as applied CRD
- `additionalProperties: false` on root may reject forward-compatible fields — intentional for Sprint 91 sketch

---

## 13. Recommended next sprint

**Sprint 92** — schema-to-CRD planning doc or compatibility sample copy (still non-applied); optional `jsonschema` dev dependency for full validation in CI.

---

## Related

- `HYPERDENSITY_KHR_RESOURCELEASE_MINIMAL_CONTRACT.md`
- `HYPERDENSITY_KHR_KUBEVIRT_CAPABILITY_INVENTORY.md`
- `HYPERDENSITY_KHR_OVN_SDN_CAPABILITY_INVENTORY.md`
