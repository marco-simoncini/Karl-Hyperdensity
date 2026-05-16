# ADR-0004 — KHR foundation source of truth (branch audit)

| Field | Value |
|-------|-------|
| **Status** | Accepted (foundation — post Sprint 91 audit) |
| **Date** | 2026-05-16 |
| **Supersedes** | None (complements ADR-0001, ADR-0002, ADR-0003) |

---

## Context

Five repositories are checked out on branch **`KHR`**. Product code still ships KubeVirt + kube-ovn/multus paths while KHR native runtime matures.

---

## Decision

1. **Source of truth for target architecture**: `Karl-Hyperdensity` docs under `docs/khr/`, `docs/architecture/`, `docs/adr/`, `docs/contracts/khr/`, and `api/crds/`.
2. **Source of truth for production behavior today**: `Karl-Dashboard` Parent Fabric runtime + `FluidVirt` + `Karl-OS-ISO` provisioning.
3. **KHR Engine runtime code** lives in `Karl-Hyperdensity/pkg/khr/` (agent, cgroup, telemetry, resourcelease dry-run, evidence) — **not** a separate repo in this phase.
4. **Do not create a new KHR repository** until Phase 3 deliverables require independent release cadence (recommendation only).

---

## Consequences

- Dashboard must converge toward Shell/Cell/ResourceLease contracts without breaking VM discovery paths.
- Multus/NAD remain **transitional** only; OVN-native KARL Network Fabric is the network target.
- KubeVirt remains **compatibility provider** per ADR-0002.

---

## Related

- `docs/KHR_BRANCH_AUDIT_AND_FOUNDATION_PLAN.md`
- `docs/adr/ADR-0001-khr-shell-cell-runtime-model.md`
- `docs/adr/ADR-0002-kubevirt-as-legacy-provider.md`
