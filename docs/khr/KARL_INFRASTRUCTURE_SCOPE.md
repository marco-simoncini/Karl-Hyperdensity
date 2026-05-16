# KARL infrastructure scope (KHR-Z)

| Field | Value |
|-------|-------|
| **Sprint** | KHR-Z |
| **Purpose** | Align product language across KHR program |

---

## Positioning

**KARL is not only a datacenter OS.** KARL is an **infrastructure operating layer** — also described as **infrastructure OS** and **infrastructure control plane** — that unifies host runtime, resource markets, certification, policy, and operator workflows across deployment environments.

| Term | Meaning |
|------|---------|
| **Infrastructure operating layer** | Product identity: runtime + policy + evidence + governance above raw orchestrators |
| **Infrastructure OS** | Shorthand for KARL-owned host/runtime semantics (KHR, Shell/Cell, providers) |
| **Infrastructure control plane** | Declarative APIs, CRDs, Dashboard projection, Inventory observation — not autonomous production orchestration |

**Datacenter** remains a valid **deployment environment** (on-prem metal, private cloud, edge racks) — not the sole product definition.

---

## Target environments (non-exhaustive)

| Environment | KARL posture |
|-------------|--------------|
| On-prem / datacenter metal | Primary certification sandbox (`karl-metal-01@ovh`) |
| Private cloud / hybrid | Adaptive enforcement via cloud providers |
| Edge / regional | Same Shell/Cell model; provider-dependent ports |
| Public cloud | Constraint-aware market; no hypervisor hotplug promise |

---

## What KARL is not (language guardrails)

| Avoid as sole framing | Use instead |
|----------------------|-------------|
| Sole datacenter product framing | Infrastructure OS / infrastructure operating layer |
| Unqualified GA or production-ready claims (KHR program) | Certified-preview, sandbox, technical preview |
| Autonomous orchestration | Operator-gated, read-only simulation paths |

---

## KHR program components (reference)

See `TECHNICAL_PREVIEW_READINESS.md` (TP scorecard) and `KHR_RELEASE_READINESS_MAP.md` for readiness by area: API foundation, native-live lane, certification, policy gates, operator approval, control graph, provenance, and cross-repo status (Dashboard, Inventory, ISO).

---

## Related docs

| Doc | Repo |
|-----|------|
| `docs/architecture/KARL_HOST_RUNTIME_VISION.md` | Karl-Hyperdensity |
| `docs/adr/ADR-0001-khr-shell-cell-runtime-model.md` | Karl-Hyperdensity |
| `KHR_PROJECTION_V1.md` | Karl-Dashboard |
| `INVENTORY_RUNTIME_POSTURE.md` | Karl-Inventory |
| `KHR_HOST_RUNTIME_PREVIEW.md` | Karl-OS-ISO |
