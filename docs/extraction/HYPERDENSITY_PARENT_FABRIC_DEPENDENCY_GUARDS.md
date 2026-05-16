# Hyperdensity Parent Fabric — dependency guard requirements (Sprint 44)

Future **`pkg/hyperdensity/parentfabric/**`** pure packages must remain **extractable** and **testable** without the OpenShift console runtime. This document lists **intended** static guard rules.

**Sprint 45:** **`scripts/validate_parentfabric_pure_deps.sh`** runs from **`scripts/validate.sh`** and **fails** on forbidden import strings under `pkg/hyperdensity/parentfabric` (grep-based; extend with `go list -deps` in a later sprint).

---

## Parentfabric pure packages — default deny list

Pure trees (`parentfabric`, `summary`, `governance`, `evidence`, `recommendation` — see **`HYPERDENSITY_PARENT_FABRIC_EXTRACTION_BOUNDARY.md`**) **must not** import:

| Pattern | Rationale |
|---------|-----------|
| **`k8s.io/*`** | Kubernetes client / API machinery belongs in Dashboard adapters or future operator code, not pure-core. |
| **`kubevirt.io/*`** | KubeVirt provider coupling; **no removal** of KubeVirt from product — only avoid entangling **pure** packages. |
| **`github.com/gorilla/*`** | HTTP / websocket stack tied to console server. |
| **`net/http`** | Unless isolated in a clearly named **`adapter/http`** package behind an interface (explicit sprint). |
| **Dashboard repo paths** | e.g. `github.com/openshift/console/pkg/server` — pure packages must not call back into the console. |
| **`Karl-Dashboard` string** | Sprint 45 guard rejects this literal in `parentfabric` sources to avoid accidental cross-repo path pastes. |

---

## Dashboard consumption rules

- Dashboard may import Hyperdensity `parentfabric` pure packages **only** with an **allowlist + sprint** (same discipline as contractkit runtime expansion).
- **contractkit** `claimpolicy` / `contracts` **runtime freeze** stays **independent** — see Dashboard **`HYPERDENSITY_CONTRACTKIT_RUNTIME_IMPORT_FREEZE_M17.md`**; do not conflate with Parent Fabric extraction gates.

---

## Enforcement status

| Guard | Mechanism (Sprint 45) |
|-------|---------------------|
| **Import deny (static)** | `scripts/validate_parentfabric_pure_deps.sh` — patterns above; requires `executiontypes/`, `workload/`, `primitives/` subpackages |
| **Import deny (deps)** | Future: `go list -deps` on `pkg/hyperdensity/parentfabric/...` |

---

## Related

- `HYPERDENSITY_PARENT_FABRIC_EXTRACTION_BOUNDARY.md`
- `HYPERDENSITY_PARENT_FABRIC_EXTRACTION_PHASES.md`
- `HYPERDENSITY_PARENT_FABRIC_PURE_PACKAGE_SKELETON.md`
- `scripts/validate_parentfabric_pure_deps.sh`
- Dashboard `docs/hyperdensity/HYPERDENSITY_PARENT_FABRIC_EXTRACTION_BOUNDARY_M28.md`
