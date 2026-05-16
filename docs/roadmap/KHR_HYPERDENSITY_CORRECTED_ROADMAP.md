# KHR + Hyperdensity — corrected roadmap (post Sprint 16)

This roadmap **supersedes sequencing intent** in `docs/roadmap/KHR_HYPERDENSITY_SPRINT_ROADMAP.md` for **extraction-first** work: **real Hyperdensity currently lives in Karl-Dashboard**; Hyperdensity repo is the **canonical contracts + extractable libraries** home.

---

## Principles

1. **Extract real Hyperdensity first** — pure functions, DTOs, blocker catalogs, cohort math, evidence merge rules from `Karl-Dashboard/kubernetes-console/pkg/server/hyperdensity_*.go` into `Karl-Hyperdensity/pkg/hyperdensity/*` per `docs/extraction/HYPERDENSITY_PACKAGE_TARGET_PLAN.md`.
2. **Dashboard imports Hyperdensity module** — `go.mod` dependency; handlers become thin wrappers calling extracted libs; **no behavior change** in the first integration PRs.
3. **KHR adapter** — `khr-linux-agent` consumes `contracts` + evidence APIs; does not synthesize parent-fabric.
4. **KHR apply gate** — explicit policy boundary before any host mutation; aligns with Dashboard `no_production_mutation` style gates.
5. **Lab-only apply** — Linux reference cells only; dual-write telemetry; rollback mandatory.
6. **Bare metal OS integration** — after lab apply stable on reference fleet.
7. **Public cloud adaptive mode** — policy packs + blast-radius contracts; optional slower track.
8. **Dashboard migration KubeVirt → Shell/KHR** — **last** major UX shift; KubeVirt remains supported per ADR-0002 until cutover criteria met.
9. **Windows lane** — **planning- and evidence-only** until an explicit promotion milestone; keep `no_windows_lane` default posture in live collectors; Hyperdensity `windowsfluidvirt` remains non-production.

---

## Host runtime milestone (KHR-F / G / I)

| Item | State |
|------|-------|
| `karl-host-runtime` Linux MVP skeleton | Shipped (KHR-F) |
| Sandbox execution on `karl-metal-01@ovh` | **PASS** — [`docs/evidence/khr-runtime-sandbox/summary.json`](../evidence/khr-runtime-sandbox/summary.json) |
| Host CR + status JSON (`runtime.karl.io/Host`) | Shipped (KHR-I) — evidence [`docs/evidence/khr-host-registration/summary.json`](../evidence/khr-host-registration/summary.json) |
| Production host mutation | **Unsupported** |
| ISO `karl-host-runtime` | Preview / **disabled** by default |

---

## Near-term milestones (engineering order)

| Step | Outcome |
|------|---------|
| **H+2** | **ResourcePort controller loop** (observe/report candidates; sandbox namespace only; no production apply) |
| **H+3** | Host status **cluster apply** + heartbeat (sandbox only; behind explicit gate) |
| M1 | Golden JSON: Dashboard parent-fabric summary → Hyperdensity `contracts` testdata |
| M2 | `pkg/hyperdensity/blockers` exported IDs == Dashboard gate strings |
| M3 | First pure merge function ported with `_test.go` from Dashboard |
| M4 | Dashboard `go.mod` pins `Karl-Hyperdensity` pseudo-version / replace directive in dev |
| M5 | CI job: `go test ./pkg/hyperdensity/...` on both repos |

---

## Dependencies

- CTO/architect sign-off on **package boundaries** (this audit + package plan).
- Release process for tagging `Karl-Hyperdensity` modules consumed by Dashboard.

---

## Related documents

- `docs/extraction/HYPERDENSITY_REAL_EXTRACTION_AUDIT.md`
- `docs/extraction/HYPERDENSITY_PACKAGE_TARGET_PLAN.md`
- `docs/extraction/HYPERDENSITY_KHR_DUPLICATION_REPORT.md`
- `docs/migration/dashboard-to-hyperdensity-extraction-plan.md`
- `Karl-Dashboard`: `docs/hyperdensity/HYPERDENSITY_RUNTIME_DEPENDENCY_MAP.md`
