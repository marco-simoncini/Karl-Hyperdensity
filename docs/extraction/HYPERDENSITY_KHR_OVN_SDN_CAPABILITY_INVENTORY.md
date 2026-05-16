# Hyperdensity / KHR — OVN/SDN capability inventory (Sprint 90)

## Summary

Read-only inventory of network/SDN capabilities. OVN/SDN remains a **preserved Dashboard capability**; indexed implementation files are **sparse** in primary clones.

---

## 1. Scope

Inventory and mapping only — no network runtime changes.

---

## 2. Repositories scanned

Karl-Hyperdensity, Karl-Dashboard, FluidVirt (read-only), Karl-OS-ISO docs (spot-check).

---

## 3. Search terms

See `HYPERDENSITY_KHR_INVENTORY_SCAN_METHOD.md`.

---

## 4. OVN/SDN findings

### Indexed in Karl-Dashboard (runtime/docs — not OVN controller source)

| Path | Finding |
|------|---------|
| `pkg/server/hyperdensity_parent_fabric_live.go` | `KubeOVNNotReady` launcher failure classification |
| `pkg/server/hyperdensity_parent_fabric_vm_linux_cpu_burst.go` | Blocker text: kube-ovn / multus sandbox failure |
| `pkg/server/hyperdensity_parent_fabric_vm_linux_memory_runtime.go` | `kube_ovn_ovs_interface_not_ready` |
| `deployment/hyperdensity/control-plane/README.md` | multus / kube-ovn pod sandbox notes |
| `docs/hyperdensity/vm-linux-memory-support-contract-v1.md` | kube-ovn migration prerequisites |
| `artifacts/**` | CNI network-status JSON: `"name": "kube-ovn"`, gateway, IPAM |

### Karl-Hyperdensity

- **No** OVN controller Go package in repo.
- Network semantics in Sprint 88–89 **docs only** (`HYPERDENSITY_KHR_NETWORK_SDN_SEMANTICS.md`).

### FluidVirt (read-only)

| Path | Finding |
|------|---------|
| `pkg/virt-controller/services/template.go` | `multus.NetworkToResource` for VMI network attachment |

### Not found in primary clone

- Dedicated OVN schema repos (LogicalSwitch/Router CRD definitions as first-class KARL CRDs)
- `NetworkAttachmentDefinition` YAML templates in Karl-Hyperdensity examples

**Statement**: OVN/SDN is **user-provided existing Dashboard capability to preserve**. No indexed local OVN implementation file in Karl-Hyperdensity. Operational evidence exists in Dashboard **artifacts** and **runtime error strings** referencing **kube-ovn** + **multus**.

**Future inventory**: run in clone containing full SDN implementation (if separate repo exists outside workspace).

---

## 5. NetworkAttachment findings

| Source | Mapping |
|--------|---------|
| multus CNI + NAD (operational) | `NetworkAttachment`, `providerBinding` |
| kube-ovn default network in artifacts | `ShellNetwork` / tenant segment |
| VMI `interfaces` + `networks` in KubeVirt YAML | compatibility layer → `NetworkAttachment` |

---

## 6. Tenant/gateway/policy findings

| Concept | Evidence |
|---------|----------|
| Tenant isolation | Sprint 88 Windows DaaS profile; subnet IPs in kube-ovn CNI status (10.16.x.x) |
| Gateway | `gateway` in CNI network-status artifacts |
| Ingress/egress | Parent Fabric exposure model in docs; rdp-GW in architecture memory |
| ACL/NAT | **Documental** OVN mapping in Sprint 88–89; no ACL YAML indexed in Hyperdensity |

---

## 7. ResourceLease network mapping

| OVN / operational | ResourceLease |
|-----------------|---------------|
| LogicalSwitch | NetworkSegment / KARLNetwork |
| LogicalRouter | NetworkGateway |
| LogicalPort | NetworkAttachment |
| ACL | NetworkPolicy |
| NAT/SNAT/DNAT | ServiceExposure / ingress / egress |
| DHCPOptions | NetworkLease / NetworkLeaseConfig |
| FloatingIP | ExternalEndpoint |
| NAD + multus | NetworkAttachment + providerBinding |

---

## 8. Provider mapping

| Operational stack | Provider ID |
|-------------------|-------------|
| kube-ovn on cluster | `khr.native.ovn` or `kubevirt.legacy.ovn` |
| multus | `kubernetes.cni` + attachment binding |
| masquerade pod network (winsrv VM) | default pod network → map to `NetworkAttachment` role=primary |

---

## 9. Gaps

- No first-class OVN CRD inventory in Karl-Hyperdensity.
- LogicalSwitch/Router not grep-visible as schema in Dashboard `pkg/server` (runtime handles symptoms only).
- Full SDN repo may be outside scanned workspace.

---

## 10. Risks

- Assuming kube-ovn artifacts equal complete OVN product model.
- Skipping multus/NAD when designing NetworkAttachment contract.

---

## 11. Recommended next sprint

Dedicated **OVN inventory** in SDN-owning repo when available; add non-applied `examples/compatibility/ovn-tenant-network.yaml` sketch.

---

## Related

- `HYPERDENSITY_KHR_RESOURCELEASE_INVENTORY_MAPPING.md`
- `HYPERDENSITY_KHR_NETWORK_SDN_SEMANTICS.md`


---

## Sprint 91 (ResourceLease JSON Schema)

Sprint 91 adds non-applied JSON Schema and example fixtures under docs/contracts/khr/. No CRD, no controller, no runtime. See HYPERDENSITY_KHR_RESOURCELEASE_JSON_SCHEMA.md.
