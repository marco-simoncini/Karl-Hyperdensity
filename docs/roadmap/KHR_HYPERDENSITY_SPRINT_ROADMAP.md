# KHR and Hyperdensity — Sprint Roadmap

**Status:** Planning artifact (Sprint 1). Dates are sequencing, not calendar commitments.

**Principles:** No destructive refactors; KubeVirt remains until explicitly superseded; Kubernetes stays the control plane; documentation and contracts precede daemons.

---

## Sprint 1 — Architecture foundation *(current)*

**Goal:** Lock the strategic contract across `marco-simoncini/*` without runtime changes.

**Deliverables:**

- `docs/architecture/KARL_HOST_RUNTIME_VISION.md`
- `docs/architecture/KARL_SHELL_CELL_MODEL.md`
- `docs/architecture/KARL_HYPERDENSITY_KHR_FUSION.md`
- `docs/adr/ADR-0001-khr-shell-cell-runtime-model.md`
- `docs/adr/ADR-0002-kubevirt-as-legacy-provider.md`
- `docs/adr/ADR-0003-hyperdensity-resourcelease-resourcefuture.md`
- `docs/roadmap/KHR_HYPERDENSITY_SPRINT_ROADMAP.md` *(this file)*

**Exit criteria:** CTO/investor/lead-dev can answer “what is KARL?” without Kubernetes-or-VM centrism; engineering knows Shell vs Cell vs KHR vs RuntimeProvider.

---

## Sprint 2 — CRD contract sketch

**Goal:** Publish **OpenAPI-compatible** CRD drafts (v1alpha1) for Shell, Cell, ShellLease, ResourcePort, ResourceLease, ResourceFuture—**schema + validation markers only** in a chosen repo (likely `Karl-Hyperdensity` for contracts, mirrored later).

**Non-goals:** No production controller behavior change.

---

## Sprint 3 — KubeVirt legacy wrapper (control plane only)

**Goal:** Controller that maps a Shell with `kubevirt.legacy.v1` to existing KubeVirt objects **without** user-facing removal of VMs.

**Repos:** `Karl-Dashboard` and/or new controller under `Karl-Warden` (TBD in Sprint 2 ownership note).

---

## Sprint 4 — Dashboard Shell model

**Goal:** UI/API reads primarily **Shell**; VM views become secondary/operator mode.

**Repos:** `marco-simoncini/Karl-Dashboard`, `marco-simoncini/KARL-APP` for catalog alignment.

---

## Sprint 5 — KHR Linux MVP

**Goal:** Host agent skeleton on Linux nodes: register host, report ResourcePort, apply **approved** cgroup envelope leases for Linux Cells (narrow scope).

**Repos:** new `Karl-KHR` or `Karl-Warden` extension (decision in Sprint 2); ISO packaging hooks deferred until MVP stable in lab.

---

## Sprint 6 — Hyperdensity ResourceLease engine

**Goal:** End-to-end **lease approval → KHR apply → verification** for Linux reference class; reuse `Karl-Hyperdensity` validators and evidence JSON patterns.

**Repos:** `Karl-Hyperdensity`, `Karl-Dashboard` extraction touchpoints.

---

## Sprint 7 — rdp-GW Shell Gateway

**Goal:** Map **ShellLease** + gateway route to RDP file issuance; leverage existing RemoteApp / alternate shell fields in `marco-simoncini/rdp-GW`.

**Repos:** `rdp-GW`, `Karl-Dashboard` API facade.

---

## Sprint 8 — Windows Host Agent

**Goal:** Privileged agent on Windows farm nodes; session/app materialization; telemetry hooks; cooperation with Hyperdensity risk indices.

**Repos:** new Windows agent repo or `Karl-Warden/Win-Side` evolution—decision gated on security review.

---

## Sprint 9 — Native VM-like Cell provider

**Goal:** Evaluate `FluidVirt` as provider backend; define ResourcePort truth; **lab** hardening before GA claims.

**Repos:** `FluidVirt`, `Karl-OS-ISO` (optional preview flags only).

---

## Sprint 10 — Public cloud adaptive mode

**Goal:** Capability map ingestion, adaptive Hyperdensity policies, cloud provider adapter in KHR; honest UX for limits vs bare metal.

**Repos:** cloud adapter (likely new), `Karl-Installer` profiles, `Karl-Migration-Factory` for hybrid.

---

## Dependency graph (simplified)

```
Sprint 1 docs → Sprint 2 CRDs → Sprint 3 legacy wrapper + Sprint 5 KHR MVP
                     ↓
              Sprint 6 lease engine
                     ↓
       Sprint 7 gateway ∥ Sprint 8 Windows agent
                     ↓
              Sprint 9 native VM-like → Sprint 10 cloud adaptive
```

---

## Recommended next sprint

**Sprint 2 (CRD contract)** immediately after merging Sprint 1 documentation: it turns ADRs into machine-verifiable shapes without touching production runtime paths.
