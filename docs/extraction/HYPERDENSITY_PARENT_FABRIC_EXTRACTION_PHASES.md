# Hyperdensity Parent Fabric — extraction phase matrix (Sprint 44)

Phased approach for **real** extraction of Parent Fabric / Hyperdensity surfaces from **Karl-Dashboard** → **Karl-Hyperdensity**. **Sprint 44** documents only; **no** runtime move.

**Sprint 45 (Phase 1–2):** mechanical CSV inventory + heuristic guesses; Hyperdensity **`parentfabric` stdlib skeleton** + dependency guard.

**Sprint 46 (Phase 3 minimal):** first candidate `hyperdensity_parent_fabric_execution_types.go` → **`parentfabric/executiontypes`** copy-contract + golden; Dashboard source **unchanged**; **no** production import wiring.

**Sprint 47:** **`SourceManifest`** anti-drift guard + Dashboard `go/ast` parity test; **no** new type copy; **no** production import wiring.

**Sprint 48:** **`workload_helpers.go`** audit-first → **`copy-deferred`**; `parentfabric/workload` placeholder only. Phase 3 must **not** become bulk copy.

**Sprint 49:** **`parentfabric/primitives`** stdlib contract (nested map + quantity); Dashboard primitive loci audit; **workload_helpers still deferred**.

**Sprint 50:** workload **adapter boundary** docs + Dashboard **classification fixture** (46 functions); **no** adapter code; **copy-deferred** unchanged.

**Sprint 51:** Dashboard **test-only adapter stubs** + golden manifest; Hyperdensity **stub readiness** doc; **no** production wiring; **copy-deferred** unchanged.

**Sprint 52:** **three pure-candidates** copy-contract in `parentfabric/workload` + golden + Dashboard parity test; **full workload_helpers still `copy-deferred`**; **no** production Dashboard import.

**Sprint 53:** Dashboard **production-internal adapter v1** (path + observation) + tests; **not wired** to handlers; Hyperdensity hardening doc only.

**Sprint 54:** Dashboard **shadow tests** (legacy vs adapter v1 on path + observed-state); **still not wired**; Hyperdensity shadow-test doc only.

**Sprint 55:** Wiring **proposal** + call-site inventory audit (51 production sites).

**Sprint 56:** **Path-only** wiring on 6 approved non-apply files via `hyperdensityWorkloadPath*V1` wrappers.

**Sprint 57:** **Pilot-only** observed-state via `hyperdensityWorkloadPilotObservedStateV1`; apply/live/VM/rollback observation excluded.

**Sprint 58:** **Pilot observation hardening** (end-to-end tests + live inventory proposal); **no new wiring**.

---

## Phase 1 — Audit real Dashboard files

| | |
|--|--|
| **Repos** | Karl-Dashboard (inventory), Karl-Hyperdensity (boundary docs) |
| **Risk** | Low — read-only |
| **PASS** | Complete categorized inventory (see **`HYPERDENSITY_PARENT_FABRIC_RUNTIME_FILE_INVENTORY_M27.md`**) + script listing (`list_parent_fabric_runtime_files.sh` optional) + **CSV** `parent_fabric_runtime_inventory_m29.csv` (full file list) |
| **Rollback** | N/A (docs only) |
| **Forbidden** | Deleting or renaming production `.go`; changing handlers |

---

## Phase 2 — Identify pure helpers

| | |
|--|--|
| **Repos** | Karl-Dashboard (annotate), Karl-Hyperdensity (target package mapping) |
| **Risk** | Medium — misclassification could move the wrong code later |
| **PASS** | Each candidate tagged: **pure candidate** vs **adapter needed** vs **runtime-bound** (see M27 columns); Sprint 45 adds **M29** audit doc + CSV heuristics |
| **Rollback** | Revert doc annotations |
| **Forbidden** | Moving code; adding Hyperdensity imports to Dashboard runtime without allowlist sprint |

---

## Phase 3 — Move **only** pure helpers

| | |
|--|--|
| **Repos** | Karl-Hyperdensity (new `pkg/hyperdensity/parentfabric/...`), Karl-Dashboard (thin wrappers) |
| **Risk** | Medium — API drift if public types change |
| **PASS** | `go test` green in both repos; **no** handler signature change; JSON ordering unchanged for stable fixtures |
| **Rollback** | Revert commits; restore Dashboard-local copies |
| **Forbidden** | Moving HTTP handlers, `client-go` calls, or apply paths in this phase |

**Status (Sprint 46):** **partial** — `executiontypes` copy-contract only (summary + engine spine). Full engine + nested surfaces remain in Dashboard. Production import wiring still **forbidden**.

---

## Phase 4 — Test adapter in Dashboard

| | |
|--|--|
| **Repos** | Karl-Dashboard |
| **Risk** | Low–medium — import cycle risk if boundaries wrong |
| **PASS** | Dashboard tests call Hyperdensity pure packages; parity scripts green |
| **Rollback** | `replace` or version pin rollback in `go.mod` |
| **Forbidden** | Behavior change in production responses |

---

## Phase 5 — Evaluate **runtime** import (later)

| | |
|--|--|
| **Repos** | Karl-Dashboard, Karl-Hyperdensity |
| **Risk** | **High** — operational coupling |
| **PASS** | Dedicated ADR + security review + canary plan |
| **Rollback** | Feature flag / revert to Dashboard-local runtime |
| **Forbidden** | Doing this before Phases 1–4 are stable |

---

## Phase 6 — Runtime ownership move (much later)

| | |
|--|--|
| **Repos** | Karl-Dashboard, Karl-Hyperdensity, possibly operators/agents |
| **Risk** | **Very high** — production blast radius |
| **PASS** | Explicit multi-sprint program; KHR / ResourceLease alignment **after** Parent Fabric pure extraction stabilizes (per product roadmap) |
| **Rollback** | Blue/green deploy; pin to previous module version |
| **Forbidden** | Big-bang cutover without adapter phase |

---

## Related

- `HYPERDENSITY_PARENT_FABRIC_EXTRACTION_BOUNDARY.md`
- `HYPERDENSITY_PARENT_FABRIC_DEPENDENCY_GUARDS.md`
- `HYPERDENSITY_PARENT_FABRIC_PURE_PACKAGE_SKELETON.md`
- `HYPERDENSITY_PARENT_FABRIC_PURE_CANDIDATE_AUDIT.md`
- Dashboard `docs/hyperdensity/HYPERDENSITY_PARENT_FABRIC_EXTRACTION_BOUNDARY_M28.md`
- Dashboard `docs/hyperdensity/HYPERDENSITY_PARENT_FABRIC_PURE_CANDIDATE_AUDIT_M29.md`
- Dashboard `docs/hyperdensity/HYPERDENSITY_PARENT_FABRIC_EXTRACTION_STATUS_M30.md`
