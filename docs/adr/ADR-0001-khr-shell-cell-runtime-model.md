# ADR-0001 — Adopt Shell / Cell / KHR as the native KARL runtime model

| Field | Value |
|-------|-------|
| **Status** | Accepted (architecture foundation — Sprint 1) |
| **Date** | 2026-05-15 |
| **Decision makers** | KARL CTO / platform architecture |
| **Applies to** | KARL product platform (`marco-simoncini/*` operational repositories) |

---

## Context

KARL ships as an **infrastructure operating layer** (infrastructure OS / infrastructure control plane) and cloud-adaptive platform — not only a datacenter product. On-prem datacenter metal is one deployment environment among others (cloud, hybrid, edge). Today, customer and internal mental models often collapse to **“Kubernetes + KubeVirt + Dashboard”**, which under-describes KARL’s differentiation and over-couples the product to a specific virtualization implementation.

Kubernetes remains valuable as a **declarative API and reconciliation ecosystem**, but it must not define KARL’s identity.

Operational repositories already show:

- Hyperdensity resource market and evidence discipline concentrated in `marco-simoncini/Karl-Hyperdensity` and runtime-adjacent code in `marco-simoncini/Karl-Dashboard`.
- KubeVirt as an installed engine dependency in `marco-simoncini/Karl-OS-ISO` and `marco-simoncini/Karl-Installer`.
- RDP gateway code in `marco-simoncini/rdp-GW` with RemoteApp-capable RDP field model.

There is no single first-class **host execution daemon** vocabulary spanning Linux, Windows, VM-like, and cloud-backed workloads—this ADR introduces that vocabulary.

---

## Decision

1. Adopt **Shell** as the **primary product abstraction** for workloads and experiences.
2. Adopt **Cell** as the **node-materialized execution primitive** produced by **KHR (KARL Host Runtime)** or legacy paths until migration completes.
3. Adopt **RuntimeProvider** as the **only** supported extension point for execution backends (Linux pod, Windows agent, cloud adapter, **kubevirt legacy**, FluidVirt-class native VM-like when ready).
4. Maintain Kubernetes as the **control plane substrate** (CRDs, RBAC, etcd, controller patterns) without positioning it as the product definition.

---

## Consequences

### Positive

- Clear separation of **customer semantics** (Shell) from **infrastructure reality** (Cell).
- Hyperdensity can target **uniform lease verbs** across providers where ResourcePort allows.
- Gradual KubeVirt de-emphasis without a risky “big bang” removal.

### Negative / costs

- Additional abstraction layer and migration work across Dashboard, APIs, and controllers.
- Risk of “dual write” periods where both raw KubeVirt objects and Shell/Cell must stay consistent—requires disciplined ownership and feature flags.

### Neutral

- Documentation and ADR surface area increases (intentional for enterprise buyers and auditors).

---

## Alternatives considered

1. **Stay VM-centric (KubeVirt as hero)**  
   Rejected: caps strategic narrative and ties economics to upstream roadmap.

2. **Replace Kubernetes entirely**  
   Rejected: discards a mature declarative ecosystem; cost is incompatible with near-term delivery.

3. **Use only OAM / Crossplane / other generic workload model**  
   Rejected: insufficiently specific to DaaS, Windows session semantics, and Hyperdensity market making.

---

## Repository impact (documentation and future work)

| Repository | Impact |
|------------|--------|
| `Karl-Hyperdensity` | Contracts, ADRs, market verbs aligned to Shell/Cell; KHR apply path integration in later sprints. |
| `Karl-Dashboard` | UI/API language migration toward Shell; technical views for Cells. |
| `Karl-OS-ISO` / `Karl-Installer` | No install behavior change in Sprint 1; future manifests may bundle KHR. |
| `rdp-GW` | GatewayRoute target; Shell/App issuance semantics. |
| `KARL-APP` / `karl-directoryservice` / `Karl-Inventory` / `Karl-DLP` / `Karl-Licenziatore` | Attach policy and catalog to Shell / ShellLease identifiers. |
| `Karl-Migration-Factory` | Mapping legacy VMs to Shells during hybrid mode. |
| `Karl-Warden` | Host integrity, pool policy, agent attestation alongside KHR registration. |
| `FluidVirt` | Optional Native VM-like RuntimeProvider — not KHR itself. |

No runtime code changes are mandated by this ADR in Sprint 1.

---

## References

- `docs/architecture/KARL_HOST_RUNTIME_VISION.md`
- `docs/architecture/KARL_SHELL_CELL_MODEL.md`
- `docs/architecture/hyperdensity-overview.md` (existing)
