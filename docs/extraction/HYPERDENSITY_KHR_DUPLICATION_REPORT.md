# KHR vs Dashboard Hyperdensity — duplication & conflict report (Sprint 16)

## Executive summary

`Karl-Hyperdensity:KHR` today holds **contracts, CRDs, Linux KHR agent skeleton, Grande Padre local evidence/recommendation (dry-run), and WindowsFluidVirt lab packages**.  
`Karl-Dashboard:KHR` holds the **authoritative live Hyperdensity runtime** (parent-fabric APIs, cohort builders, collectors, UI).

Overlap is **conceptual** (blockers, slates, donor/receiver language) — not a line-by-line duplicate codebase. Risk is **divergent semantics** if Hyperdensity skeletons evolve without pinning to Dashboard JSON.

---

## KHR Sprints 1–13 — what stays useful as vNext

| Asset | Verdict |
|-------|---------|
| `pkg/khr/*` Linux cgroup telemetry, discovery, evidence bundle | **Keep** — orthogonal host lane; feeds **evidence** contracts parent-fabric may ingest later. |
| `pkg/khr/evidenceingest` | **Keep** — aligns with `EvidenceIngestRequest` / manifest model; future bridge to Dashboard ingestion pipeline. |
| `pkg/grandepadre/evidence` | **vNext useful** — local index/query; **freeze** scope creep; map indices to Dashboard “readonly observation” language. |
| `pkg/grandepadre/recommendation` | **vNext useful** — dry-run slate only; **do not** market as replacement for Dashboard action slate / execution engine. |
| `api/crds`, `schemas`, `docs/contracts` | **Keep as SoT** for shapes; Dashboard snapshots should **converge** here over time. |
| `pkg/windowsfluidvirt` | **vNext useful** — deepest Windows FluidVirt logic already in Hyperdensity; Dashboard Windows server files should **delegate** or import this package later (facade in `windowslane`). |

---

## What duplicates Dashboard “Father” behavior

- **Naming:** “action slate”, “donor/receiver”, “equilibrium” appear in both repos with **different code paths** — Hyperdensity versions are **skeleton / dry-run**.
- **Blocker strings:** partial overlap with Dashboard gate IDs — must be **catalogued** (`pkg/hyperdensity/blockers`) to avoid forked spellings.

---

## What to freeze (no expansion until extraction anchor exists)

- Grande Padre **recommendation engine** feature set — only maintenance + alignment PRs.
- New Hyperdensity HTTP services — **not** until Dashboard imports libraries.
- Parallel “mini parent-fabric” in Go — **forbidden**; use extraction matrix instead.

---

## What must connect to Dashboard reality

- Any new **contract** field in `schemas/` or CRDs → **diff** against `HyperdensitySnapshot` TypeScript (`hyperdensity-model.ts`) before merge.
- Windows lane → must respect collector gates (`no_windows_lane`) and Hyperdensity `pkg/windowsfluidvirt/safety.go` planning-only policy.

---

## What not to expand now

- Second implementation of VM cohort merge logic in Hyperdensity **before** porting tests from `hyperdensity_parent_fabric_vm_linux_*_test.go`.
- Auto-apply, autonomous market controller, or “production readiness” scoring in Hyperdensity without Dashboard parity tests.

---

## Resolution strategy

1. **Single write path** for live Hyperdensity: Dashboard until explicit cutover.  
2. **Single read model** aspiration: contracts in Hyperdensity generated or verified from Dashboard golden outputs.  
3. **KHR agent** remains consumer of contracts, not producer of parent-fabric snapshots.
