# ResourceLease — JSON Schema vs CRD OpenAPI delta

| Field | Value |
|-------|-------|
| **Audit date** | 2026-05-16 |
| **Sprint** | KHR-A |
| **Schema** | `docs/contracts/khr/resourcelease.schema.json` (Sprint 91 sketch) |
| **CRD** | `api/crds/hyperdensity.karl.io/resourcelease.yaml` (v1alpha1 deployable) |

---

## 1. Executive summary

The two artifacts describe **different contract generations**:

| Aspect | JSON Schema (Sprint 91) | CRD OpenAPI |
|--------|-------------------------|-------------|
| **Mental model** | Unified workload lease (Shell + Cell + provider + storage + network) | Donor → receiver resource transfer between Shell/Cell |
| **Primary spec keys** | `shell`, `cell`, `provider`, `resources`, `storage`, `network`, `policy` | `donor`, `receiver`, `resource`, `mode` |
| **Provider** | Top-level enum string (`khr.native`, `kubevirt.compatibility`, …) | Not present |
| **Status** | Structured (`phase`, `providerBinding`, …) | Opaque (`preserve-unknown-fields`) |
| **Applied to cluster** | **No** (explicit in schema description) | **Yes** (CRD manifest) |

**Resolved (KHR-B):** **Unified ResourceLease** per `docs/adr/ADR-0005-resourcelease-unification.md`. Schema and CRD aligned; `spec.leaseKind` = `runtime` | `transfer`.

---

## 2. Fields present in JSON Schema, missing in CRD

| Schema path | Notes |
|-------------|-------|
| `spec.shell` | Shell ref, experience, kind |
| `spec.cell` | Cell ref, hostSelector, runtimeClass |
| `spec.provider` | Provider enum (11 values) |
| `spec.resources` | cpu/memory request/limit |
| `spec.storage` | disks[], ephemeral modes, promoteToImage |
| `spec.network` | attachments, exposure, providerNetwork, networkLease |
| `spec.policy` | Required in schema; no CRD equivalent |
| `spec.evidence` | Optional object |
| `spec.rollback` | Optional object (CRD has `rollbackPlanRef` only) |
| `spec.expiration` | Optional |
| `spec.promotion` | Optional |
| `status.phase` | CRD status is opaque |
| `status.providerBinding` | Missing |
| `status.effectiveResources` | Missing |
| `status.conditions` | Missing |
| `status.observedAt` | Missing |

---

## 3. Fields present in CRD, missing in JSON Schema

| CRD path | Notes |
|----------|-------|
| `spec.donor` | kind Shell\|Cell, name, namespace, apiGroup |
| `spec.receiver` | same shape as donor |
| `spec.resource` | enum: cpu, memory, disk, network, gpu |
| `spec.amount` | opaque object |
| `spec.mode` | string (unconstrained in CRD) |
| `spec.durationSeconds` | int64 |
| `spec.ttlSeconds` | int64 |
| `spec.rollbackPlanRef` | name/namespace ref |
| `spec.rollbackRequired` | boolean |
| `spec.verificationHooks` | opaque |
| `spec.noRestart` | boolean intent |
| `spec.guestVisible` | boolean |
| `spec.telemetryConvergedRequired` | boolean |
| `spec.dryRunOnly` | boolean (Grande Padre gate) |

---

## 4. Naming mismatches

| Concept | JSON Schema | CRD |
|---------|-------------|-----|
| Lease parties | `spec.shell` + `spec.cell` | `spec.donor` + `spec.receiver` (each Shell or Cell) |
| Provider | `spec.provider` (string enum) | *(absent — implied by RuntimeProvider CRD elsewhere)* |
| Rollback | `spec.rollback` (object) | `spec.rollbackPlanRef` + `spec.rollbackRequired` |
| Dry-run | *(not in schema)* | `spec.dryRunOnly` |
| Resource quantity | `spec.resources.cpu/memory` | `spec.resource` + `spec.amount` |
| Network | `spec.network.attachments[]` | `spec.resource: network` only |
| Storage disks | `spec.storage.disks[]` with `diskMode` enum | `spec.resource: disk` only |

---

## 5. Status / spec lifecycle mismatches

| Lifecycle aspect | JSON Schema | CRD |
|------------------|-------------|-----|
| **Status subresource** | Documented fields | Enabled; schema is free-form |
| **Required spec** | 7 top-level groups | 4 fields: donor, receiver, resource, mode |
| **Phase values** | `phase` string (unenum'd) | Not defined |
| **Reconciliation** | Implied via `providerBinding` | Implied via `telemetryConvergedRequired`, `dryRunOnly` |
| **Completion** | Not specified | Not specified |

**Recommended unified lifecycle (P1):**

```
Pending → DryRunValidated → Bound → Active → Completing → Completed | Failed | RolledBack
```

Map CRD `dryRunOnly=true` to terminal **DryRunValidated** without promotion.

---

## 6. Enum / type mismatches

| Field | JSON Schema | CRD |
|-------|-------------|-----|
| Provider | 11-value enum | N/A |
| `diskMode` | ephemeralOverlay, ephemeralClone, scratch, readonly, persistent | N/A |
| `sourceType` | pvc, image, snapshot, volume, goldenImage | N/A |
| `discardPolicy` | deleteOnStop, keepOnFailure, promoteOnRequest | N/A |
| `spec.resource` | N/A | cpu, memory, disk, network, gpu |
| donor/receiver `kind` | N/A | Shell, Cell |

---

## 7. TODO backlog

### P0

| ID | Action |
|----|--------|
| P0-1 | **Architecture decision:** merge models (A: extend CRD with schema sections, B: split ResourceTransfer vs ResourceLease) |
| P0-2 | Add `spec.dryRunOnly` to JSON Schema or document mapping from CRD |
| P0-3 | Align `apiVersion`/`kind` const in schema with CRD group `hyperdensity.karl.io/v1alpha1` |

### P1

| ID | Action |
|----|--------|
| P1-1 | Generate OpenAPI from schema (or vice versa) in CI `scripts/validate_resourcelease_schema.sh` |
| P1-2 | Define shared `status.phase` enum in both artifacts |
| P1-3 | Map `pkg/khr/resourcelease/dryrun.go` input to chosen canonical spec |

### P2

| ID | Action |
|----|--------|
| P2-1 | Deprecate donor/receiver-only CRD if unified model wins |
| P2-2 | Provider binding object in status aligned with RuntimeProvider CRD |

---

## 8. References

- `docs/contracts/khr/resourcelease.schema.json`
- `api/crds/hyperdensity.karl.io/resourcelease.yaml`
- `docs/extraction/HYPERDENSITY_KHR_RESOURCELEASE_MINIMAL_CONTRACT.md`
- `pkg/khr/resourcelease/dryrun.go`
