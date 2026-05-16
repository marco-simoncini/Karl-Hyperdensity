# Hyperdensity / KHR — ResourceLease provider contract (Sprint 89)

## Summary

Defines allowed `spec.provider` values and selection rules for ResourceLease.

---

## Provider values

| Provider ID | Role |
|-------------|------|
| `khr.native` | Native KHR host runtime (preferred) |
| `kubevirt.compatibility` | KubeVirt VM/VMPool mapping |
| `kubevirt.public-cloud-fallback` | KHR unavailable on public cloud |
| `khr.native.ovn` | Native OVN on KHR |
| `kubevirt.legacy.ovn` | KubeVirt + legacy OVN attachment |
| `kubernetes.cni` | CNI-backed workloads |
| `cloud.vpc.aws` | AWS VPC |
| `cloud.vnet.azure` | Azure VNet |
| `cloud.vpc.gcp` | GCP VPC |
| `baremetal.bridge` | L2 bridge |
| `baremetal.vlan` | VLAN segment |

---

## Rules

1. **KHR native preferred** where host/runtime available.
2. **Public cloud**: when KHR cannot run, use `kubevirt.public-cloud-fallback` for VM/VMPool generation.
3. **KubeVirt is provider/fallback**, not product model — **Shell/Cell** remain product model.
4. **Provider selection must be explicit** in every ResourceLease `spec.provider` — no implicit default to KubeVirt in sketch validation.

---

## Network provider binding

`spec.network.providerNetwork.provider` may refine storage-level provider (e.g. `khr.native.ovn` vs `kubevirt.legacy.ovn`).

---

## Related

- `HYPERDENSITY_KHR_RESOURCELEASE_MINIMAL_CONTRACT.md`
- `HYPERDENSITY_KHR_ARCHITECTURE_MEMORY.md`


---

## Sprint 90 (inventory facts)

Sprint 90 adds read-only KubeVirt and OVN/SDN capability inventory mapped to ResourceLease contract. No CRD, no controller, no runtime. See HYPERDENSITY_KHR_KUBEVIRT_CAPABILITY_INVENTORY.md and related Sprint 90 docs.
