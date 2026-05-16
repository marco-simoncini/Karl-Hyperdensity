# ResourceLease lifecycle (unified v1alpha1)

| Field | Value |
|-------|-------|
| **ADR** | ADR-0005 |
| **Controller** | Not implemented (contract-only) |

---

## Phases

| Phase | Meaning |
|-------|---------|
| `Pending` | Accepted, not yet validated |
| `DryRunValidated` | Dry-run / Grande Padre gates passed; `spec.governance.dryRunOnly=true` may stop here |
| `Bound` | Provider binding resolved (`status.providerBinding`) |
| `Active` | Resources applied per lease kind |
| `Completing` | Rollback or completion in progress |
| `Completed` | Terminal success |
| `Failed` | Terminal failure |
| `RolledBack` | Terminal after rollback |

---

## By `spec.leaseKind`

### `runtime`

1. Validate Shell/Cell refs and provider enum.
2. Validate storage disk modes and network attachments (no NAD-first target architecture).
3. Bind provider → `Bound`.
4. Apply resources/storage/network → `Active` (future KHR agent).
5. Complete or rollback → `Completed` / `RolledBack`.

### `transfer`

1. Validate `spec.transfer` donor/receiver (Shell or Cell).
2. ResourcePort compatibility check (dry-run).
3. Apply envelope mutation or block if `dryRunOnly`.
4. Telemetry gate if `telemetryConvergedRequired`.

---

## Non-goals (KHR-B)

- No reconciler
- No automatic promotion from `DryRunValidated` to `Active`
- No cluster apply from Dashboard projection
