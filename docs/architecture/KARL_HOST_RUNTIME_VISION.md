# KARL Host Runtime (KHR) — Platform Vision

**Status:** Architecture foundation (Sprint 1) — documentation only; no runtime or install behavior changes in this milestone.

**Repository map (authoritative operational remotes):** `marco-simoncini/*` as listed in the KARL program charter. This document is versioned in `marco-simoncini/Karl-Hyperdensity` alongside Hyperdensity contracts because that repository is the declared canonical home for resource-control contracts, evidence models, and extraction planning toward the full product core.

---

## Executive thesis

KARL is not “Kubernetes with VMs.” KARL is an **infrastructure operating layer** (infrastructure OS / infrastructure control plane): a coherent runtime, scheduling, lifecycle, storage, network, telemetry, policy, access, and governance surface that happens to use Kubernetes as a **declarative control plane and reconciliation substrate**, not as the product identity. Datacenter deployments are a major target environment, not the sole product definition.

**KHR (KARL Host Runtime)** is the **host-native daemon/runtime** on every KARL node—analogous to the kubelet for Kubernetes, or to libvirt’s operational role for traditional Linux virtualization—but scoped to **KARL’s execution semantics**, not to re-implementing KVM in kernel space.

Kubernetes remains responsible for:

- Declarative desired state (CRDs, API machinery).
- Durable state store and watch semantics.
- Controller reconciliation patterns and RBAC integration.

KHR becomes responsible for:

- **Runtime plane** execution on the node.
- **Cell materialization** (see `KARL_SHELL_CELL_MODEL.md`).
- **RuntimeProvider** selection and local orchestration.
- Host-level telemetry convergence and policy application where the host is owned by KARL.

**Hyperdensity** remains the economically decisive capability: a **live resource market** with evidence-backed movement, not “autoscaling with extra steps.” KHR is the architectural prerequisite that makes Hyperdensity’s decisions **physically enforceable on hosts KARL owns**, with explicit degradation modes on public cloud.

---

## Repo-aware current state (high signal)

| Area | Repository (remote) | Current truth relevant to KHR |
|------|---------------------|-------------------------------|
| Hyperdensity contracts, schemas, executors library, “Grande Padre” framing | `marco-simoncini/Karl-Hyperdensity` | Canonical documentation and Go packages for market ticks, leases slate, shell passport, install gates, canary cohort evidence, Windows FluidVirt **boundary** packages (not active Windows product claims per repo README). |
| Live Hyperdensity runtime / parent fabric / APIs | `marco-simoncini/Karl-Dashboard` (e.g. branch `The-Father`) | README in this repo’s Hyperdensity tree states implementation ownership today; extraction toward `Karl-Hyperdensity` is planned. |
| OS image / engine stack | `marco-simoncini/Karl-OS-ISO` | Provisioning matrix documents KubeVirt, virtctl, CDI as **installed core** virtualization paths. |
| Cluster installer | `marco-simoncini/Karl-Installer` | Documents KubeVirt install/cleanup flows—operational dependency remains. |
| RDP / gateway | `marco-simoncini/rdp-GW` | Go RDP file model already exposes RemoteApp-related fields and alternate shell (strategic for Shell/App gateway evolution). |
| VM-like experimentation / KubeVirt fork | `marco-simoncini/FluidVirt` | `go.mod` identifies `kubevirt.io/kubevirt`; tree includes FluidVirt-specific hyperdensity accounting and live resource control experiments—**not** KHR, but a **candidate Native VM-like RuntimeProvider** lab asset. |
| Enforcement / pool / agent direction | `marco-simoncini/Karl-Warden` | Existing architecture docs describe Warden 2.0; future alignment with KHR node registration and policy export is expected. |
| Identity, inventory, DLP, directory, licensing, migration | `marco-simoncini/karl-directoryservice`, `Karl-Inventory`, `Karl-DLP`, `Karl-Licenziatore`, `Karl-Migration-Factory`, `KARL-APP` | Product adjacency: Shell leases, catalog, compliance, and hybrid migration **consume** KHR/Hyperdensity primitives; no change in this sprint. |

---

## Why KubeVirt must be de-emphasized (without removing it)

KubeVirt correctly solved “VMs on Kubernetes” for early KARL. Strategically it **over-centers the product domain** in the customer mental model (“KARL = KubeVirt + UI”) and couples economic differentiation (Hyperdensity) to a **third-party virtualization personality** instead of to **KARL-owned runtime semantics**.

**Decommissioning KubeVirt immediately is explicitly out of scope** (compatibility, migration, ISO/install matrices, existing VMs). The architectural correction is:

- **Promote** Shell / Cell / KHR / RuntimeProvider as the conceptual API surface.
- **Retain** KubeVirt as a **legacy / compatibility RuntimeProvider** behind the same Shell abstraction.

See ADR-0002.

---

## KHR architecture (target)

```
KARL Dashboard / product APIs
        ↓
KARL Shell API / Cell API (product-facing)
        ↓
Kubernetes CRDs + controllers (desired state, durable store)
        ↓
KARL controllers / Grande Padre / schedulers / Hyperdensity
        ↓
KHR (per-node daemon)
        ↓
RuntimeProvider implementations
  - Linux Pod / systemd provider
  - Windows Session / App providers (+ Windows Host Agent)
  - Native VM-like provider (FluidVirt or successor where applicable)
  - Cloud adapter (adaptive)
  - KubeVirt legacy provider
```

KHR **does not replace Kubernetes**; it **specializes node execution** for KARL’s model the same way specialized agents extend generic orchestration in other platforms.

---

## Shell vs Cell (summary)

- **Shell:** user/product-visible workload identity (desktop, app, Linux environment, DaaS session, legacy VM-like experience).
- **Cell:** node-materialized execution unit (pod, container, session, VM-like object, cloud-backed instance, legacy KubeVirt VM).

Full definitions: `KARL_SHELL_CELL_MODEL.md`.

---

## RuntimeProvider model

A **RuntimeProvider** is a bounded plugin contract:

- Declares **ResourcePort** capabilities (CPU/RAM/disk/network/GPU modes).
- Accepts **ResourceLease** / **ResourceFuture** instructions when Hyperdensity or controllers commit market outcomes.
- Reports observed runtime and blocker states into the evidence plane.

Providers differ by **host ownership** (bare metal KARL vs guest on public cloud).

---

## Bare metal vs public cloud vs hybrid

1. **KARL OS — Bare metal / on-prem datacenter mode**  
   Full KHR authority where KARL owns the host stack (including datacenter metal). Hyperdensity can pursue aggressive live movement subject to **ResourcePort** truth and guest cooperation.

2. **KARL — Public cloud mode**  
   KHR operates in **adaptive** posture: Shell API stays uniform; enforcement defers to cloud APIs and provider limits; Hyperdensity becomes **constraint-aware market making**, not a promise of hypervisor-grade hotplug.

3. **KARL — Hybrid / migration mode**  
   Single control plane across KARL metal, cloud, and legacy; `marco-simoncini/Karl-Migration-Factory` remains the strategic bridge; KubeVirt provider stays until workloads are expressed as Shell/Cell with verified providers.

---

## Relationship to Hyperdensity

Hyperdensity decides **who should receive or surrender resources** under risk, priority, and evidence rules. KHR is the **apply path** that turns committed decisions into **node-local reality**, honoring `noRestart`, rollback, verification, and telemetry convergence described in `KARL_HYPERDENSITY_KHR_FUSION.md`.

---

## Relationship to rdp-GW

`marco-simoncini/rdp-GW` already models RDP parameters including **RemoteApp** and **alternate shell** fields in `cmd/rdpgw/rdp/rdp.go`. The target is **Shell Gateway / App Gateway**: issuance of RDP artifacts and tokens from **ShellLease** and Windows placement decisions, coordinated with Windows Host Agent.

---

## Windows Host Agent

Privileged component on Windows hosts implementing **Windows Shell** materialization (sessions, profiles, apps, policy, telemetry). It is the **kubelet analog for Windows Shells**—not “Windows in a pod.”

---

## FluidVirt relevance

`FluidVirt` is a **KubeVirt-derived codebase** with additional FluidVirt and hyperdensity-adjacent packages. It is **not** KHR. It is a **candidate lab/foundation** for a future **Native VM-like RuntimeProvider** where KARL requires deeper VM-class behavior than the legacy KubeVirt path alone—subject to separate hardening and productization gates.

---

## Synthetic roadmap

| Horizon | Outcome |
|---------|---------|
| Near | CRD/API contracts for Shell/Cell; KubeVirt legacy wrapper behind Shell. |
| Mid | KHR Linux MVP; Hyperdensity ResourceLease engine on KHR apply path. |
| Mid+ | rdp-GW as Shell/App gateway; Windows Host Agent. |
| Long | Native VM-like provider maturity; public cloud adaptive mode; progressive reduction of KubeVirt centrality in UX and defaults. |

Sprint-level plan: `docs/roadmap/KHR_HYPERDENSITY_SPRINT_ROADMAP.md`.

---

## Explicit non-goals (vision level)

- Rewriting KVM or replacing cloud provider hypervisors.
- Promising universal hotplug, universal GPU live attach, or universal live migration across Shell kinds.
- Removing KubeVirt in a single release.

These constraints are product-credible and engineering-honest.
