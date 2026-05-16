# Hyperdensity / KHR — ResourceLease network contract (Sprint 89)

## Summary

Details `spec.network` for ResourceLease minimal contract. Anchored to `HYPERDENSITY_KHR_NETWORK_SDN_SEMANTICS.md`.

---

## network top-level

```yaml
network:
  attachments: []
  policies: []
  exposure: { ingress, egress, directPublic }
  providerNetwork: { provider, domainRef }
  networkLease: { leaseRef, dhcp, addresses }
```

---

## attachment fields

| Field | Description |
|-------|-------------|
| `name` | Attachment name |
| `networkRef` | KARLNetwork / segment ref |
| `role` | primary, secondary, mgmt, app |
| `ipam` | static / dhcp / pool |
| `macPolicy` | random / preserved |
| `isolation` | strict / shared |
| `providerBinding` | ovn.logicalPort, cni, vpc, … |

---

## Required primitives (lease must reference)

KARLNetwork, ShellNetwork, CellNetwork, NetworkAttachment, NetworkLease, NetworkPolicy, NetworkSegment, NetworkGateway, ServiceExposure, ExternalEndpoint

---

## OVN/SDN compatibility

| OVN | KHR |
|-----|-----|
| LogicalSwitch | NetworkSegment / TenantNetwork |
| LogicalRouter | NetworkGateway / RouteDomain |
| LogicalPort | NetworkAttachment |
| ACL | NetworkPolicy |
| NAT/SNAT/DNAT | EgressPolicy / IngressPolicy / ServiceExposure |
| DHCPOptions | NetworkLeaseConfig |
| FloatingIP | ExternalEndpoint |

---

## Windows DaaS network

| Setting | Value |
|---------|-------|
| tenantIsolation | **strict** |
| ingress | **rdp-GW** / session-gateway |
| egress | **controlled** |
| directPublicExposureDefault | **false** |

---

## Related

- `HYPERDENSITY_KHR_RESOURCELEASE_MINIMAL_CONTRACT.md`
- `HYPERDENSITY_KHR_NETWORK_SDN_SEMANTICS.md`


---

## Sprint 90 (inventory facts)

Sprint 90 adds read-only KubeVirt and OVN/SDN capability inventory mapped to ResourceLease contract. No CRD, no controller, no runtime. See HYPERDENSITY_KHR_KUBEVIRT_CAPABILITY_INVENTORY.md and related Sprint 90 docs.


---

## Sprint 91 (ResourceLease JSON Schema)

Sprint 91 adds non-applied JSON Schema and example fixtures under docs/contracts/khr/. No CRD, no controller, no runtime. See HYPERDENSITY_KHR_RESOURCELEASE_JSON_SCHEMA.md.
