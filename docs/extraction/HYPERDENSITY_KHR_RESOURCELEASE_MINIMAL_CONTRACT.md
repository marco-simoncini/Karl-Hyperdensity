# Hyperdensity / KHR — ResourceLease minimal contract (Sprint 89)

## Summary

**Sprint 89** defines the first **minimal ResourceLease contract** anchored to Sprint 88 architecture memory. **Contract-sketch / docs / parity-anchor only** — no CRD applied, no controller, no runtime consumer.

---

## 1. Scope

| In scope | Out of scope |
|----------|--------------|
| Documental contract shape | Kubernetes CRD activation |
| Storage / network / provider sections in lease | Controllers, KHR agent |
| Parity goldens in Dashboard | API / Parent Fabric runtime changes |
| Validation invariants | Broad observation, adapter drift |

---

## 2. Non-goals

- Apply CRD to cluster
- Implement reconciliation controller
- Activate providers operationally
- Modify `ObservationWiredV1` / `ProductionWiredV1` / resource_exchange state
- Parent Fabric runtime import from Hyperdensity `parentfabric`

---

## 3. Contract status

| Property | Value |
|----------|-------|
| Status | **direction-only / schema-sketch** |
| CRD applied | **no** |
| Runtime consumer | **no** |
| Provider activated | **no** |

---

## 4. ResourceLease top-level shape

```yaml
apiVersion: hyperdensity.karl.io/v1alpha1  # sketch only
kind: ResourceLease
metadata: { ... }
spec: { ... }
status: { ... }
```

---

## 5. Identity section

`metadata`: name, namespace, labels, annotations, `resourceVersion`, lease UID references for evidence chains.

---

## 6. Shell section

`spec.shell`: user-facing workload identity.

| Field | Purpose |
|-------|---------|
| `kind` | e.g. `windowsDesktop`, `vmLike`, `linuxShell`, `containerSession` |
| `ref` | Shell instance reference |
| `experience` | UX profile hints |

**Shell** is the product model — not KubeVirt VM.

---

## 7. Cell section

`spec.cell`: executable unit on host.

| Field | Purpose |
|-------|---------|
| `ref` | Cell materialization reference |
| `hostSelector` | Placement hints |
| `runtimeClass` | KHR envelope class |

---

## 8. Provider section

`spec.provider`: explicit backend selection. See `HYPERDENSITY_KHR_RESOURCELEASE_PROVIDER_CONTRACT.md`.

---

## 9. CPU / memory section

`spec.resources`:

```yaml
resources:
  cpu: { request, limit }
  memory: { request, limit }
```

Maps to Cell cgroup envelope; observed via existing equilibrium paths — **not** wired in Sprint 89.

---

## 10. Storage section

`spec.storage`: disks[], defaultDiscardPolicy, promoteToImage. See `HYPERDENSITY_KHR_RESOURCELEASE_STORAGE_CONTRACT.md`.

---

## 11. Network section

`spec.network`: attachments[], policies[], exposure, providerNetwork, networkLease. See `HYPERDENSITY_KHR_RESOURCELEASE_NETWORK_CONTRACT.md`.

---

## 12. Policy section

`spec.policy`: admission references, blast-radius, guarded-apply flags, network/storage policy bundles.

---

## 13. Evidence section

`spec.evidence`: EvidenceBundle refs, ingest request IDs, integrity hashes (aligns with KHR evidence model).

---

## 14. Rollback section

`spec.rollback`: observed-state rollback hooks — **separate** from Parent Fabric rollback.go legacy surface (Sprint 87 classification).

---

## 15. Expiration section

`spec.expiration`: TTL, renewal policy, grace period.

---

## 16. Promotion section

`spec.promotion`: promote-to-image, scale-up, donor/receiver promotion actions.

---

## 17. Public-cloud fallback rule

When `spec.provider` = `kubevirt.public-cloud-fallback`: KARL generates VM/VMPool via KubeVirt where KHR cannot run. **Shell/Cell/ResourceLease** remain product model; KubeVirt is backend only.

---

## 18. Validation invariants

- All Sprint 88 storage modes and sources must be expressible in `spec.storage.disks[]`
- All Sprint 88 network primitives must be expressible in `spec.network`
- `ObservationWiredV1` / `ProductionWiredV1` remain **false**
- `workload_helpers.go` verdict remains **copy-deferred**
- No `pkg/hyperdensity/parentfabric` import in Dashboard runtime

---

## 19. Risks

- Treating sketch YAML as deployed CRD
- Omitting explicit provider field → ambiguous KubeVirt vs KHR binding
- Collapsing rollback lease section with legacy rollback.go wiring

---

## 20. Recommended next sprint

**Sprint 90** — inventory facts / OVN file inventory, or JSON Schema file under `docs/contracts/` without cluster apply.

---

## status shape

```yaml
status:
  phase: Pending|Bound|Active|Expired|Failed
  providerBinding: { provider, backendRef, observedGeneration }
  effectiveResources: { cpu, memory, storageSummary, networkSummary }
  conditions: []
  observedAt: RFC3339 timestamp
```

---

## Related

- `HYPERDENSITY_KHR_RESOURCELEASE_STORAGE_CONTRACT.md`
- `HYPERDENSITY_KHR_RESOURCELEASE_NETWORK_CONTRACT.md`
- `HYPERDENSITY_KHR_RESOURCELEASE_PROVIDER_CONTRACT.md`
- `HYPERDENSITY_KHR_RESOURCELEASE_EXAMPLES.md`
- `HYPERDENSITY_KHR_ARCHITECTURE_MEMORY.md`


---

## Sprint 90 (inventory facts)

Sprint 90 adds read-only KubeVirt and OVN/SDN capability inventory mapped to ResourceLease contract. No CRD, no controller, no runtime. See HYPERDENSITY_KHR_KUBEVIRT_CAPABILITY_INVENTORY.md and related Sprint 90 docs.


---

## Sprint 91 (ResourceLease JSON Schema)

Sprint 91 adds non-applied JSON Schema and example fixtures under docs/contracts/khr/. No CRD, no controller, no runtime. See HYPERDENSITY_KHR_RESOURCELEASE_JSON_SCHEMA.md.
