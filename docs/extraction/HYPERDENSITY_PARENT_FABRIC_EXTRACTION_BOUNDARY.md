# Hyperdensity Parent Fabric — extraction boundary & target package design (Sprint 44)

## Purpose

Define **where** real Hyperdensity / Parent Fabric logic may eventually live in **Karl-Hyperdensity** when extracted from **Karl-Dashboard**, **without** moving runtime code yet. This is **design + audit scope only** — no API response change, no Parent Fabric runtime behavior change, no JSON ordering change, and no execution/apply path edits.

**Karl-Dashboard** remains the **runtime owner** until an explicit sprint authorizes a move.

---

## Design rule: pure-core first

Future `pkg/hyperdensity/parentfabric/...` trees must start as **pure-core** (stdlib + domain types + deterministic transforms) **before** any “runtime extraction” sprint wires Kubernetes, HTTP, or apply paths.

---

## Candidate target packages (future)

| Package | Responsibility | Pure / testable slice | Must **not** import (initial phase) |
|---------|----------------|------------------------|-------------------------------------|
| **`pkg/hyperdensity/parentfabric`** | Root types, shared constants, cross-surface orchestration **interfaces** (no handlers). | Version strings, category enums, pure validation of DTOs, golden-test helpers. | `k8s.io/*`, `kubevirt.io/*`, `net/http`, Dashboard `pkg/server`, gorilla, mutable apply executors. |
| **`pkg/hyperdensity/parentfabric/summary`** | Summary / performance / redaction **pure transforms** aligned with parity manifests. | Mapping from internal structs → redacted summary DTOs; fixture-driven tests. | HTTP routers, `client-go`, browser/frontend paths. |
| **`pkg/hyperdensity/parentfabric/governance`** | Policy pack, support matrix, limitation IDs — **catalog-aligned** pure rules. | Rule tables, string ID validation, “dry run only” category checks without cluster I/O. | Live cluster mutation, Windows enablement toggles, apply pipelines. |
| **`pkg/hyperdensity/parentfabric/evidence`** | Evidence sufficiency, bundle shape, witness collection **logic** (pure). | Evidence envelope math, sufficiency predicates, demo scenario packs as data + pure functions. | KubeVirt runtime hooks, SSH/exec adapters (until an explicit adapter sprint). |
| **`pkg/hyperdensity/parentfabric/recommendation`** | Action slate / futures / market **pure scoring** (if extracted). | Deterministic ranking from inputs; snapshot tests. | Cockpit HTTP, session store, OpenShift console deps. |

---

## Forbidden dependency classes (conceptual)

Until an **adapter sprint** explicitly documents exceptions, pure packages **must not** depend on:

| Class | Examples |
|-------|----------|
| **Kubernetes clients** | `k8s.io/client-go/...`, `k8s.io/apimachinery/...` controllers |
| **KubeVirt** | `kubevirt.io/client-go/...`, KubeVirt API types in hot paths |
| **HTTP handlers** | `net/http` handlers, `github.com/gorilla/*` mux, GraphQL resolvers |
| **Dashboard server** | `github.com/openshift/console/pkg/server/...` (or any Dashboard-internal path) |
| **Browser / frontend** | JS bundles, asset pipelines |
| **Mutable runtime / apply** | Code paths that perform production mutation; keep in Dashboard until a gated migration |

Detailed guard text: **`HYPERDENSITY_PARENT_FABRIC_DEPENDENCY_GUARDS.md`**.

---

## Ownership model

| Phase | Runtime owner | Pure-core owner |
|-------|---------------|-----------------|
| **Now – pre-extraction** | **Karl-Dashboard** (`pkg/server/hyperdensity_parent_fabric_*`) | **Karl-Hyperdensity** (contractkit, schemas, docs only) |
| **After pure helper extraction** | **Still Dashboard** (handlers, wiring, I/O) | **Hyperdensity** (moved pure packages + tests) |
| **After dedicated runtime sprint** | TBD per ADR | TBD |

---

## Related

- `HYPERDENSITY_PARENT_FABRIC_EXTRACTION_PHASES.md`
- `HYPERDENSITY_PARENT_FABRIC_DEPENDENCY_GUARDS.md`
- `HYPERDENSITY_CONTRACTKIT_RUNTIME_IMPORT_FREEZE` narrative on Dashboard (`HYPERDENSITY_CONTRACTKIT_RUNTIME_IMPORT_FREEZE_M17.md`) — **contractkit** governance stays **separate** from Parent Fabric extraction.
- Dashboard `docs/hyperdensity/HYPERDENSITY_PARENT_FABRIC_RUNTIME_FILE_INVENTORY_M27.md`
- Dashboard `docs/hyperdensity/HYPERDENSITY_PARENT_FABRIC_EXTRACTION_BOUNDARY_M28.md`
