# ResourceLease TP Freeze Candidate (KHR-AE)

| Field | Value |
|-------|-------|
| **Sprint** | KHR-AE |
| **Contract** | `hyperdensity.karl.io/v1alpha1` / `ResourceLease` |
| **Cluster reference** | `karl-metal-01@ovh` |
| **Verdict** | **APPROVED — TP freeze candidate** (contract + read-only projection) |
| **Apply execution** | **Remains sandbox experimental** — not production-enabled |

**NOT production ready.** No autonomous orchestration. No ISO systemd enable.

---

## Freeze candidate decision

| Layer | TP freeze candidate? | Notes |
|-------|----------------------|-------|
| CRD OpenAPI (`resourceleases.hyperdensity.karl.io`) | **Yes** | Additive-only through beta-1 |
| JSON schema + contract examples | **Yes** | `docs/contracts/khr/resourcelease.schema.json` |
| Lifecycle documentation | **Yes** | `RESOURCELEASE_LIFECYCLE.md`, ADR-0005 |
| Dry-run semantics | **Yes** | Evidence on `karl-metal-01@ovh` |
| Dashboard `ResourceLeaseSummary` projection | **Yes** | Parent contract `khr-projection-v1alpha1-readonly-y` (≥ readonly-m) |
| Guarded cgroup apply | **No** | Sandbox/manual only; behavior frozen in docs, not promoted |
| ResourceLease reconciler | **No** | Not implemented |
| Production namespace apply | **No** | Explicitly forbidden |

Validation: `./scripts/khr_resourcelease_freeze_check.sh`

---

## Current fields (CRD spec)

| Field | Purpose | Freeze |
|-------|---------|--------|
| `spec.leaseKind` | `runtime` \| `transfer` | **Frozen** |
| `spec.shell` / `spec.cell` | Workload refs | **Frozen** |
| `spec.provider` | Provider enum (khr.native, kubevirt.compatibility, …) | **Frozen** |
| `spec.resources` | CPU/RAM/disk envelope (open schema) | **Frozen** shape |
| `spec.storage` / `spec.network` | Attachment hints | **Frozen** shape |
| `spec.transfer` | Donor/receiver transfer envelope | **Frozen** |
| `spec.governance` | dryRunOnly, rollback, verification, noRestart | **Frozen** |
| `spec.evidence` / `spec.rollback` / `spec.promotion` | Evidence and plans | **Frozen** shape |
| `spec.donor` / `spec.receiver` (top-level) | Deprecated aliases | **Experimental** — do not use in new docs |
| `status.phase` | Lifecycle phase | **Frozen** vocabulary |
| `status.providerBinding` | Resolved provider | **Frozen** |

Schema manifest: `docs/contracts/khr/resourcelease.schema.manifest.json` (`v0.1.0-khr-b-unified`).

---

## Lifecycle (contract-only)

| Phase | Meaning |
|-------|---------|
| `Pending` | Accepted |
| `DryRunValidated` | Dry-run gates passed |
| `Bound` | Provider resolved |
| `Active` | Apply semantics (future / sandbox only today) |
| `Completing` / `Completed` / `Failed` / `RolledBack` | Terminal |

No reconciler promotes phases automatically (KHR-B non-goal).

---

## `leaseKind` semantics

### `runtime`

Shell/Cell runtime envelope: resources, storage, network attachments. Native-live lane uses `khr.native` provider with cgroup paths in sandbox.

### `transfer`

Donor/receiver (`Shell` \| `Cell`) + `resource` enum (`cpu`, `memory`, `disk`, `network`, `gpu`) + `mode`. Dry-run against ResourcePort compatibility before any apply.

---

## CPU / RAM scale (native-live constraints)

| Operation | Sandbox evidence | Production |
|-----------|------------------|------------|
| CPU scale up | `resourcelease-native-live-cpu.json`, guarded-apply CPU ≤500m | **Forbidden** |
| RAM scale up/down | `resourcelease-native-live-memory-up/down.json`, certification runs | **Sandbox only** |
| Compatibility VM path | Dry-run observed; restart may be required | Read-only projection |

Native-live requires: namespace `khr-runtime-sandbox`, labels `khr.karl.io/sandbox=true`, lane `khr.karl.io/native-live=true`.

---

## Safety gates (frozen vocabulary)

| Gate | Source |
|------|--------|
| `dry-run-only-projection` | Dashboard projection |
| `no-autonomous-apply` | Projection + evidence |
| `sandbox-namespace-only` | Scripts / config |
| `cpu-cap-500m` | Guarded apply sandbox (KHR-M) |
| `rollback-required` | `spec.governance.rollbackRequired` |
| `noRestart` | `spec.governance.noRestart` |

---

## Rollback requirements

| Requirement | Evidence |
|-------------|----------|
| `rollbackPlanRef` + plan body | Guarded apply blocked without |
| Baseline capture before apply | `sandbox/baseline-*.json` |
| `resourcelease-rollback` mode | `rollback.json` in guarded-apply evidence |
| Operator-initiated only | No autonomous rollback |

---

## Provenance requirements

| Requirement | Link |
|-------------|------|
| Apply evidence refs | `applyEvidenceRef`, `baselineRef` on projection |
| Registry lineage | Certification registry → approval → lease refs |
| Fingerprint validation | `docs/evidence/khr-provenance/summary.json` |
| Provenance mismatch blocks approval | KHR-Y evidence (not lease apply) |

ResourceLease **does not** enforce provenance in-cluster; validation is evidence/CLI only.

---

## Breaking-change policy (TP freeze)

| Allowed | Blocked |
|---------|---------|
| Add optional spec/status fields | Remove or rename frozen fields |
| Add enum value with ADR note | Change `readOnly` projection to mutating |
| Add projection fields on Dashboard | Imply production-ready apply |
| New examples under `docs/contracts/khr/examples/` | Autonomous reconcile without sandbox gates |

Pin: `hyperdensity.karl.io/v1alpha1` until beta-1 ADR for `v1beta1`.

---

## Known experimental fields / behaviors

| Item | Status |
|------|--------|
| Guarded cgroup apply (`cpu.max` write) | Sandbox only |
| ResourceLease `kubectl apply` to cluster | Not in evidence path |
| `spec.promotion` auto-promotion | Not implemented |
| Top-level `spec.donor` / `spec.receiver` | Deprecated |
| `verificationHooks` open schema | Experimental |
| Windows RAM hot-adjust | Observation / blocked dry-run |
| Controller reconciliation | **Not implemented** |

---

## Blockers before beta

| ID | Severity | Blocker |
|----|----------|---------|
| RL-B-01 | P0 | **NOT production ready** — apply remains sandbox |
| RL-B-02 | P0 | No autonomous orchestration / reconcile |
| RL-B-03 | P1 | RAM guarded apply parity vs CPU evidence |
| RL-B-04 | P1 | Inventory observation of lease posture (no apply source) |
| RL-B-05 | P2 | ResourceLease status phase automation (if ever) requires ADR |

---

## Explicit non-goals

- Production enablement or GA claims
- Autonomous orchestration or approval→apply
- ISO `karl-host-runtime` systemd enable
- Dashboard frontend rewrite or mutating action buttons
- Removing KubeVirt compatibility provider

---

## Related

- `RESOURCELEASE_LIFECYCLE.md`, `RESOURCELEASE_DRYRUN_AGAINST_RESOURCEPORT.md`, `RESOURCELEASE_GUARDED_APPLY_SANDBOX.md`
- `KHR_CONTRACT_FREEZE_PLAN.md`
- Dashboard `RESOURCELEASE_PROJECTION_FREEZE_CANDIDATE.md`
- Inventory `RESOURCELEASE_OBSERVATION_FREEZE_CANDIDATE.md`
- ISO `RESOURCELEASE_ISO_FREEZE_CANDIDATE.md`
