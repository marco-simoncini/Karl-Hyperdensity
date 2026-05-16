# Hyperdensity / KHR — ResourceLease direction (Sprint 88)

## Summary

Defines the **future ResourceLease** contract shape for KARL/Hyperdensity. Direction-only — no CRD or runtime implementation in Sprint 88.

---

## ResourceLease includes

| Dimension | Content |
|-----------|---------|
| **cpu** | Request/limit envelope for Cell |
| **memory** | Request/limit envelope |
| **storage** | Ephemeral + persistent + scratch disks; promote eligibility |
| **network** | NetworkLease, attachments, policies, tenant isolation |
| **policy** | Admission, blast-radius, guarded apply |
| **provider** | khr.native, kubevirt.compatibility, public-cloud fallback |
| **rollback** | Observed-state rollback hooks (separate surface from Sprint 78–87 resource_exchange) |
| **evidence** | Evidence bundle references |
| **expiration** | Lease TTL / renewal |
| **promotion** | promote-to-image, scale-up, donor/receiver promotion actions |

---

## Storage section

Must include:

- **ephemeralDisks** (EphemeralDisk / ephemeralOverlay / ephemeralClone)
- **persistentDisks** (profile, data volumes)
- **scratchDisks** (deleteOnStop)
- **promoteToImageEligibility** (golden image maintenance)

---

## Network section

Must include:

- **NetworkLease**
- **NetworkAttachment**
- **NetworkPolicy**
- **tenantIsolation**
- **gatewayExposure** (session gateway, not direct public by default)

---

## Provider section

| Provider | When |
|----------|------|
| **khr.native** | KHR available on host/baremetal/private cloud |
| **kubevirt.compatibility** | VM/VMPool mapping on clusters with KubeVirt |
| **publicCloudKubeVirtFallback** | Public cloud where KHR cannot run directly |

Cloud (`aws`, `azure`, `gcp`) and baremetal providers are **future** explicit entries in the same lease dimension.

---

## Related

- `HYPERDENSITY_KHR_ARCHITECTURE_MEMORY.md`
- `HYPERDENSITY_KHR_STORAGE_SEMANTICS.md`
- `HYPERDENSITY_KHR_NETWORK_SDN_SEMANTICS.md`
- `docs/contracts/resource-equilibrium-v1.md` (existing equilibrium concepts)


---

## Sprint 89 (ResourceLease minimal contract)

Sprint 89 adds ResourceLease minimal contract sketch (storage/network/provider/examples). No CRD, no controller, no runtime. See HYPERDENSITY_KHR_RESOURCELEASE_MINIMAL_CONTRACT.md and related Sprint 89 docs.
