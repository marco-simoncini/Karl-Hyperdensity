# Hyperdensity and KHR — Fusion Architecture

**Status:** Architecture foundation (Sprint 1) — aligns product narrative with repository facts in `marco-simoncini/Karl-Hyperdensity` and related operational repos.

**Related:** `KARL_HOST_RUNTIME_VISION.md`, `KARL_SHELL_CELL_MODEL.md`, ADR-0003.

---

## Strategic framing

**Hyperdensity is the economic advantage:** a **live resource market** that continuously maintains donor/receiver posture, risk posture, and **prevalidated** movement options so pressure is handled **before** it becomes outage.

**KHR is the architectural advantage:** a **host-native execution plane** that makes Hyperdensity’s commitments **physically real** on nodes KARL controls, with explicit capability truth and rollback discipline.

**Grande Padre** (already used in `Karl-Hyperdensity` documentation) is the **decisioning brain**: reconciliation loops, policy packs, gates, canary cohorts, and evidence requirements.

**KHR** is the **arm**: apply leases, mutate Cells, converge telemetry, and refuse illegal operations when **ResourcePort** forbids them.

Guiding statements (engineering, not marketing):

- Hyperdensity does not “wait to scale”; it maintains a **market of already-eligible** transfers bounded by evidence.
- With KHR, Hyperdensity becomes **host-native market making** instead of a controller that only knows API objects without uniform node semantics.
- Sub-second **decisioning** and seconds-level **transfer** remain **targets**; actual bounds are provider- and guest-dependent and must be surfaced honestly in SLAs.

---

## Repository integration: `marco-simoncini/Karl-Hyperdensity`

Today this repository contains (non-exhaustive, repo-inspected):

- **Go packages** for market controller reconciliation, live loops, durable store, install gates, admission dry-run, leases slate validation, shell passport validation, kernel boundary checks, canary cohort graduation evidence, and Windows FluidVirt **boundary/replay** test harnesses.
- **JSON schemas and examples** for canary cohorts, execution leases, rollback windows, fluidvirt invocation references, etc.—including **donorShellId** / **receiverShellId** patterns in reference surfaces.
- **Documentation** under `docs/architecture/` describing runtime overlay (Declared vs Runtime vs Observed) and Linux shell compliance semantics (live CPU/RAM, no reboot, no VMI recreate, no rollout, rollback proof).

**Gap to close in later sprints (not implemented now):**

- First-class **ResourceLease** and **ResourceFuture** CRDs shared with KHR apply path (today, lease-like concepts appear in contracts/examples and packages such as `pkg/leaseslate`).
- Uniform **ResourcePort** declaration emitted by every RuntimeProvider and consumed by Hyperdensity prevalidation.

This document names those primitives as the **contractual bridge** between Grande Padre and KHR.

---

## Grande Padre — responsibilities

- Maintain **donor index** and **receiver index** per pool/class/risk band.
- Compute **risk index**, **priority index**, and **blocked/remediable** classifications using evidence bundles.
- Produce **action slate** candidates (what could move, under which gates).
- Run **dry-run** / simulation paths (`pkg/installadmission` and related patterns exist toward dry-run discipline).
- Orchestrate **rollback** and **verification** windows (see canary cohort reference schemas for rollback gates and execution leases).

---

## Donor / receiver / risk / priority / action slate

| Construct | Meaning |
|-----------|---------|
| **Donor** | Shell/Cell aggregate eligible to surrender resources without violating its own floor / SLO / lease. |
| **Receiver** | Shell/Cell eligible to receive resources with headroom and compatible ResourcePort modes. |
| **Risk** | Quantified exposure if a lease executes or fails mid-flight (blast radius, co-tenant impact). |
| **Priority** | Business and safety ordering (tenant tier, admin override, starvation avoidance). |
| **Action slate** | Ordered, gate-checked set of candidate ResourceLeases / ResourceFutures ready for operator or policy-automated promotion. |

---

## ResourceLease

A **ResourceLease** is a **temporary contract** for moving or envelope-shaping resources between donor and receiver (or between declared/runtime planes of the same Shell where applicable).

Minimum fields (conceptual):

- `donorRef`, `receiverRef` (Shell and/or Cell refs as defined by policy phase).
- `resource` (cpu, memory, disk, network, gpu-class…).
- `amount` / `step` / `envelopeDelta`.
- `mode` (hotAdd, balloon, virtio-mem, cgroup envelope, static, provider resize…).
- `duration` / `ttl` / `renewalPolicy`.
- `rollbackPlanRef` / `rollbackRequired`.
- `verificationHooks` (post-apply probes, SLO checks).
- `noRestart` (boolean intent; must be reconciled with ResourcePort truth).
- `guestVisible` (does the guest OS see a change vs host-only envelope).
- `telemetryConverged` (required evidence before lease closes).

**Important:** A lease is **invalid** if any party’s **ResourcePort** forbids the requested `mode`.

---

## ResourceFuture

A **ResourceFuture** is a **scheduled or contingent** lease: activates when predicates hold (predictive cushion, forecast pressure, scheduled maintenance window, spot price band in cloud adaptive mode).

This is how Hyperdensity maintains **prevalidated** liquidity without constantly applying mutations.

---

## ResourcePort

**ResourcePort** is the **capability declaration** for a Cell (and optionally aggregated at Shell level). Examples:

- `cpu.mode`: `hotAdd` | `envelope` | `static`
- `memory.mode`: `hotAdd` | `balloon` | `virtioMem` | `envelope` | `static`
- `disk.mode`: `hotplug` | `static`
- `network.mode`: `hotplug` | `static`
- `gpu.mode`: `coldAttach` | `warmAttach` | `liveIfSupported`

Hyperdensity **indexes the market** only in the space allowed by the union of honest ResourcePorts.

---

## Capability-aware hotplug and no-rollout strategy

**Promise:** capability-aware runtime with honest reporting.  
**Do not promise:** universal hotplug.

**No-rollout / no-reboot** applies when:

- Linux Cells use cgroup envelope paths already described in existing Hyperdensity Linux compliance semantics.
- Windows Sessions use envelope / priority / placement semantics—not vCPU hotplug per session.
- VM-like Cells depend on guest drivers, device models, and provider implementation (KubeVirt limits remain **legacy truth**).

---

## What changes when KHR lands

- **Uniform apply API** from Hyperdensity decisions to node execution.
- **Portable evidence**: the same lease verbs can be verified against KHR-reported observed state.
- **Provider isolation**: KubeVirt quirks stop leaking into core market vocabulary.

---

## What stays limited on public cloud

- KHR runs **adaptive**: many mutations become **provider API resize** or **replacement** with different blast radius.
- Hyperdensity must ingest **cloud capability maps** (quota, allowed burst, regional stockouts).
- **ResourceFuture** liquidity may be **cash/credit metaphor** (prebuy capacity) rather than host-level donor slicing.

---

## Commercial and technical targets

**Commercial:** position Hyperdensity as measurable **unit economics** (reclaim, burst headroom, idle compression) with audit trails—aligned with existing realized savings / idle compression validator packages in-repo.

**Technical:** deterministic gates, dry-run, rollback proof, and **sub-second** decision loops where data is local; **seconds-level** apply where hardware/guest paths allow.

---

## Explicit non-promises

See vision document; repeated here for emphasis in market context:

- No universal GPU hotplug.
- No guarantee that public cloud matches bare metal live behavior.
- No “Hyperdensity fixes all tail latency.”
- KubeVirt remains until migration completes—Hyperdensity must tolerate **legacy provider** Cells.

---

## Next implementation dependencies (forward pointer)

1. Promote **ResourcePort** schema per RuntimeProvider (Sprint 2–3 contract work).
2. Wire **ResourceLease** / **ResourceFuture** to KHR gRPC/HTTP apply and Linux cgroup paths (Sprint 5–6).
3. Keep **Grande Padre** reconciliation in Kubernetes controllers; avoid splitting brain between CRD status and KHR local state without durable store discipline (see existing `durable-controller-state-kubernetes-reconciler` contract direction in this repo).

This fusion paper is the **bridge document** between Hyperdensity repository work and the KHR program.
