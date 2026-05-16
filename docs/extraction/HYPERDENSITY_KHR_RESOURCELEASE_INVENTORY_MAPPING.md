# Hyperdensity / KHR — ResourceLease inventory mapping (Sprint 90)

## Summary

Maps Sprint 90 inventory findings to Sprint 89 ResourceLease contract sections.

---

## Storage mapping

| Found capability | ResourceLease field |
|------------------|-------------------|
| KubeVirt ephemeral PVC | `storage.disks[].mode=ephemeralOverlay`, `source.type=pvc` |
| PVC / claimName | `source.type=pvc`, `sourceRef` |
| containerDisk / image | `source.type=image` |
| snapshot | `source.type=snapshot` |
| golden image | `source.type=goldenImage` |
| persistent user data volume | `mode=persistent` |
| scratch / temp | `mode=scratch`, `discardPolicy=deleteOnStop` |

**Verified reference**: `5.winsrv_vm.yaml` os volume → `role=os`, `ephemeralOverlay`, `pvc`, `karl-os-nfs`.

---

## Network mapping

| Found / OVN concept | ResourceLease field |
|---------------------|---------------------|
| OVN LogicalSwitch | NetworkSegment / KARLNetwork |
| OVN LogicalRouter | NetworkGateway |
| OVN LogicalPort | NetworkAttachment |
| ACL | NetworkPolicy |
| NAT / SNAT / DNAT | exposure / ServiceExposure / ingress / egress |
| DHCPOptions | NetworkLease / NetworkLeaseConfig |
| FloatingIP | ExternalEndpoint |
| Multus / NAD | NetworkAttachment + providerBinding |

**Operational**: kube-ovn CNI in artifacts → `providerNetwork.provider=kubevirt.legacy.ovn` or `khr.native.ovn`.

---

## Provider mapping

| Found capability | ResourceLease provider |
|------------------|------------------------|
| KubeVirt VM / VMI | `kubevirt.compatibility` |
| KubeVirt VMPool | `kubevirt.compatibility` |
| Public cloud (KHR unavailable) | `kubevirt.public-cloud-fallback` |
| kube-ovn + multus | `kubevirt.legacy.ovn` / `ovn.compatibility` |
| baremetal L2 | `baremetal.bridge` / `baremetal.vlan` |

---

## Gap handling

| Gap | Rule |
|-----|------|
| File missing in primary clone | Requirement **preserved** (user-provided) |
| OVN source not indexed | Mapping **preserved** from Sprint 88–89 semantics |
| No CRD apply | Contract sketch only |

---

## Related

- `HYPERDENSITY_KHR_KUBEVIRT_CAPABILITY_INVENTORY.md`
- `HYPERDENSITY_KHR_OVN_SDN_CAPABILITY_INVENTORY.md`
- `HYPERDENSITY_KHR_RESOURCELEASE_MINIMAL_CONTRACT.md`


---

## Sprint 91 (ResourceLease JSON Schema)

Sprint 91 adds non-applied JSON Schema and example fixtures under docs/contracts/khr/. No CRD, no controller, no runtime. See HYPERDENSITY_KHR_RESOURCELEASE_JSON_SCHEMA.md.
