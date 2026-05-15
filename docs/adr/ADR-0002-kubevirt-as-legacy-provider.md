# ADR-0002 — KubeVirt remains supported as a legacy / compatibility RuntimeProvider

| Field | Value |
|-------|-------|
| **Status** | Accepted (architecture foundation — Sprint 1) |
| **Date** | 2026-05-15 |
| **Applies to** | KARL OS, installer, dashboard, gateways, migration tooling |

---

## Context

KubeVirt is embedded in the current shipping path:

- `marco-simoncini/Karl-OS-ISO` provisioning documentation references KubeVirt, virtctl, CDI, and related manifests.
- `marco-simoncini/Karl-Installer` documents operational install/removal of the `kubevirt` namespace and APIs.
- `marco-simoncini/Karl-Dashboard` Hyperdensity control-plane notes reference KubeVirt VM labels for discovery/execution paths.

Removing KubeVirt immediately would break installations, customer VMs, and Hyperdensity evidence baselines that depend on VM Linux shells.

---

## Decision

1. **Retain KubeVirt** across supported releases until Shell/Cell coverage and migration tooling satisfy enterprise continuity requirements.
2. **Reclassify KubeVirt** as a **legacy / compatibility RuntimeProvider** (`kubevirt.legacy.v1` naming is illustrative) behind the Shell abstraction.
3. **De-emphasize KubeVirt in product narrative and default UX**, without removing technical access for operators.
4. **Do not** promise feature parity between KubeVirt Cells and future native VM-like Cells unless explicitly proven per ResourcePort matrix.

---

## Why not remove it now

- ISO and installer matrices are **truth** for field deployments.
- Migration Factory and customer contracts assume VM continuity.
- Hyperdensity Linux VM references in `Karl-Hyperdensity` README explicitly describe live resource behavior tied to current runtime implementation—needs phased port to KHR without losing proof chains.

---

## How to hide KubeVirt behind Shell / Cell

- Dashboard list views show **Shell** name, class, lease state, gateway—**not** `VirtualMachine` name as primary key.
- CRDs (future sprint) express desired Shell; a controller materializes KubeVirt VM only when `runtimeProvider: kubevirt.legacy.v1` is selected.
- Observability joins KubeVirt metrics under **Cell provider handle** fields.

---

## Impact by repository

### `marco-simoncini/Karl-OS-ISO`

- **Today:** KubeVirt remains a core engine artifact path (`install_kubevirt.sh`, vendored manifests).
- **Future:** add KHR packages alongside; adjust **runtime authority** documentation to show KHR + providers; keep KubeVirt manifests until ADR superseded by migration completion milestone.

### `marco-simoncini/Karl-Installer`

- **Today:** no change to cleanup/install docs.
- **Future:** optional KHR agent install hooks; KubeVirt steps become “legacy provider stack” section.

### `marco-simoncini/Karl-Dashboard`

- **Today:** Hyperdensity parent fabric and VM discovery reference KubeVirt labels—acceptable as implementation detail.
- **Future:** discovery listens for Shell owner labels; KubeVirt labels become an implementation of Cell identity.

---

## Migration risks

| Risk | Mitigation |
|------|------------|
| Label/schema drift between KubeVirt and Shell controllers | Versioned label mapping table; Migration Factory attestation. |
| Performance regression if extra indirection layers | Thin controllers first; measure reconciliation latency. |
| Operator confusion during transition | “Technical view” toggle; training docs. |
| Hyperdensity evidence breaks | Replay harnesses in `Karl-Hyperdensity` must include legacy provider fixtures alongside KHR fixtures. |

---

## Consequences

- Product can honestly claim **continuity** while changing **positioning**.
- Engineering must carry **temporary complexity** (dual models) until migration thresholds are met.

---

## Alternatives considered

1. **Hard-remove KubeVirt in a single release** — rejected (unacceptable customer breakage).
2. **Keep KubeVirt as hero forever** — rejected (strategic ceiling).
3. **Fork KubeVirt silently without strategy** — rejected; `FluidVirt` exists as an explicit fork/lab with its own ADR implications.

---

## References

- `docs/architecture/KARL_HOST_RUNTIME_VISION.md`
- `docs/architecture/KARL_SHELL_CELL_MODEL.md`
- `Karl-OS-ISO` / `Karl-Installer` KubeVirt references (repo-inspected)
