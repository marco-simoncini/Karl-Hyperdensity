# Hyperdensity / KHR — network / SDN semantics (canonical, Sprint 88)

## Summary

Canonizes **KARL-native network** primitives and **OVN/SDN Dashboard compatibility mapping**. If OVN implementation files are not present in a given clone, OVN/SDN is documented as **user-provided existing Dashboard capability** with a future inventory task.

---

## Primitives

| Primitive | Purpose |
|-----------|---------|
| **KARLNetwork** | Top-level network definition |
| **CellNetwork** | Network slice for Cell execution |
| **ShellNetwork** | User-visible network attachment for Shell |
| **NetworkAttachment** | Logical port binding (OVN LogicalPort analog) |
| **NetworkLease** | Leased connectivity / DHCP / address scope |
| **NetworkPolicy** | ACL / isolation rules |
| **NetworkSegment** | Tenant / logical switch segment |
| **NetworkGateway** | Router / route domain |
| **ServiceExposure** | Controlled ingress/egress service publish |
| **ExternalEndpoint** | Floating IP / external reachability |

---

## OVN/SDN mapping (Dashboard compatibility)

| OVN / SDN concept | KHR primitive |
|-------------------|---------------|
| **LogicalSwitch** | NetworkSegment / TenantNetwork |
| **LogicalRouter** | NetworkGateway / RouteDomain |
| **LogicalPort** | NetworkAttachment |
| **ACL** | NetworkPolicy |
| **NAT / SNAT / DNAT** | EgressPolicy / IngressPolicy / ServiceExposure |
| **DHCPOptions** | NetworkLeaseConfig |
| **FloatingIP** | ExternalEndpoint |

**OVN/SDN compatibility mapping** is mandatory — Dashboard semantics remain source-of-truth until KARL-native models supersede them.

---

## Provider model

| Provider ID | Role |
|-------------|------|
| `khr.native.ovn` | Native OVN on KHR host |
| `kubevirt.legacy.ovn` | KubeVirt VM via legacy OVN attachment |
| `kubernetes.cni` | CNI-backed pods/shells |
| `cloud.vpc.aws` | AWS VPC |
| `cloud.vnet.azure` | Azure VNet |
| `cloud.vpc.gcp` | GCP VPC |
| `baremetal.bridge` | L2 bridge |
| `baremetal.vlan` | VLAN segment |

---

## Windows DaaS network target

| Requirement | Target |
|-------------|--------|
| Tenant isolation | **strict** |
| Ingress | **rdp-GW** / session gateway |
| Egress | **controlled** |
| Private app network | optional |
| Direct public exposure | **false** by default |

---

## Inventory note

If repository search does not locate OVN controller/config files in the current clone, document:

> OVN/SDN is a **preserved Dashboard capability**; Sprint 89+ may add `OVN_INVENTORY.md` without blocking Sprint 88 validation.

---

## Related

- `HYPERDENSITY_KHR_ARCHITECTURE_MEMORY.md`
- `HYPERDENSITY_KHR_RESOURCELEASE_DIRECTION.md`
